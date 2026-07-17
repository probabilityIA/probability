package client

import (
	"context"
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/jumpseller/internal/app/usecases/mapper"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/jumpseller/internal/domain"
)

const storeInfoBody = `{"store":{"code":"probability-mock","name":"Tienda Mock","url":"http://jumpseller-mock.local","country":"CO","currency":"COP","hooks_token":"mock-hooks-token","weight_unit":"kg"}}`

const productsBody = `[{"product":{
  "id":100,
  "name":"Camiseta",
  "sku":"JS-001",
  "price":75000,
  "stock":10,
  "status":"available",
  "weight":1.5,
  "height":12,
  "width":20,
  "length":30,
  "diameter":5,
  "package_format":"box",
  "variants":[]
}}]`

const ordersBody = `[{"order":{
  "id":5001,
  "created_at":"2026-07-16 13:56:35 UTC",
  "status":"Paid",
  "currency":"COP",
  "subtotal":75000.0,
  "tax":14250.0,
  "shipping":8000.0,
  "shipping_required":true,
  "total":97250.0,
  "shipment_status":"requested",
  "shipping_method_name":"Envio estandar",
  "payment_method_name":"Contra entrega",
  "payment_method_type":"cod",
  "customer":{"id":"900","name":"Cliente Prueba","email":"cliente@example.com","phone":"3001234567"},
  "billing_address":{"name":"Cliente","surname":"Prueba","taxid":"1020304050","address":"Calle 123","street_number":"45","city":"Bogota","region":"Bogota D.C.","country":"Colombia"},
  "shipping_address":{"name":"Cliente","surname":"Prueba","address":"Calle 123","city":"Bogota","country":"Colombia"},
  "products":[{"id":100,"variant_id":0,"sku":"JS-MOCK-001","name":"Producto","qty":2,"price":75000.0,"tax":14250.0,"discount":0.0,"weight":1.5}]
}}]`

func newTestServer(t *testing.T) (*httptest.Server, *string) {
	t.Helper()
	var gotAuth string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		w.Header().Set("Content-Type", "application/json")

		switch {
		case r.URL.Path == "/store/info.json":
			_, _ = w.Write([]byte(storeInfoBody))
		case r.URL.Path == "/orders.json":
			_, _ = w.Write([]byte(ordersBody))
		case r.URL.Path == "/products.json":
			_, _ = w.Write([]byte(productsBody))
		case strings.HasPrefix(r.URL.Path, "/products/") && r.Method == http.MethodPut:
			_, _ = w.Write([]byte(`{"product":{"id":100,"stock":7}}`))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	t.Cleanup(server.Close)

	return server, &gotAuth
}

func testCred(server *httptest.Server) domain.Credential {
	return domain.Credential{APIKey: "login", APISecret: "token", BaseURL: server.URL}
}

func TestGetStoreInfoUsesBasicAuth(t *testing.T) {
	server, gotAuth := newTestServer(t)
	client := New()

	info, err := client.GetStoreInfo(context.Background(), testCred(server))
	if err != nil {
		t.Fatalf("GetStoreInfo: %v", err)
	}
	if info.Code != "probability-mock" || info.HooksToken != "mock-hooks-token" {
		t.Fatalf("store info mal parseado: %+v", info)
	}

	expected := "Basic " + base64.StdEncoding.EncodeToString([]byte("login:token"))
	if *gotAuth != expected {
		t.Fatalf("Authorization = %q, se esperaba %q", *gotAuth, expected)
	}
}

func TestGetOrdersParsesJumpsellerTimeAndEnvelope(t *testing.T) {
	server, _ := newTestServer(t)
	client := New()

	result, raw, err := client.GetOrders(context.Background(), testCred(server), &domain.GetOrdersParams{PerPage: 100})
	if err != nil {
		t.Fatalf("GetOrders: %v", err)
	}
	if len(result.Orders) != 1 || len(raw) != 1 {
		t.Fatalf("se esperaba 1 orden, hubo %d (raw %d)", len(result.Orders), len(raw))
	}

	order := result.Orders[0]
	if order.ID != 5001 {
		t.Fatalf("order.ID = %d", order.ID)
	}
	if order.CreatedAt.IsZero() {
		t.Fatal("created_at no se parseo: el formato de Jumpseller no es RFC3339")
	}
	if order.CreatedAt.Year() != 2026 || order.CreatedAt.Day() != 16 {
		t.Fatalf("created_at mal parseado: %v", order.CreatedAt)
	}
	if order.Customer.Email != "cliente@example.com" {
		t.Fatalf("customer mal parseado: %+v", order.Customer)
	}
}

func TestMapperProducesCanonicalOrder(t *testing.T) {
	server, _ := newTestServer(t)
	client := New()

	result, raw, err := client.GetOrders(context.Background(), testCred(server), &domain.GetOrdersParams{PerPage: 100})
	if err != nil {
		t.Fatalf("GetOrders: %v", err)
	}

	order := result.Orders[0]
	dto := mapper.MapJumpsellerOrderToProbability(&order, raw[0])

	if dto.Platform != "jumpseller" || dto.ExternalID != "5001" {
		t.Fatalf("dto mal mapeado: platform=%s external=%s", dto.Platform, dto.ExternalID)
	}
	if dto.Status != "paid" {
		t.Fatalf("status = %q, se esperaba paid", dto.Status)
	}
	if dto.TotalAmount != 97250.0 {
		t.Fatalf("total = %v", dto.TotalAmount)
	}
	if dto.Tax != 14250.0 {
		t.Fatalf("tax = %v (debe sumar tax + shipping_tax)", dto.Tax)
	}
	if dto.CustomerDNI != "1020304050" {
		t.Fatalf("customer dni = %q", dto.CustomerDNI)
	}
	if len(dto.OrderItems) != 1 {
		t.Fatalf("items = %d", len(dto.OrderItems))
	}
	item := dto.OrderItems[0]
	if item.ProductSKU != "JS-MOCK-001" || item.Quantity != 2 {
		t.Fatalf("item mal mapeado: %+v", item)
	}
	if item.Weight == nil || *item.Weight != 1.5 {
		t.Fatal("weight no mapeado")
	}
	if len(dto.Payments) != 1 || dto.Payments[0].PaymentMethodID != 6 {
		t.Fatalf("el pago contra entrega debe mapear a COD (6): %+v", dto.Payments)
	}
	if dto.CodTotal == nil || *dto.CodTotal != 97250.0 {
		t.Fatal("cod_total no se calculo para una orden contra entrega")
	}
	if len(dto.Addresses) != 2 {
		t.Fatalf("addresses = %d, se esperaban billing + shipping", len(dto.Addresses))
	}
	if dto.Addresses[0].Street != "Calle 123 45" {
		t.Fatalf("street = %q, se esperaba address + street_number", dto.Addresses[0].Street)
	}
	if dto.ChannelMetadata == nil || dto.ChannelMetadata.ChannelSource != "jumpseller" {
		t.Fatal("channel metadata ausente")
	}
	if !dto.Invoiceable {
		t.Fatal("una orden en COP debe ser facturable")
	}
}

func TestSetProductStock(t *testing.T) {
	server, _ := newTestServer(t)
	client := New()

	if err := client.SetProductStock(context.Background(), testCred(server), 100, 7); err != nil {
		t.Fatalf("SetProductStock: %v", err)
	}
}

func TestSinBaseURLFallaNoUsaFallback(t *testing.T) {
	client := New()

	_, err := client.GetStoreInfo(context.Background(), domain.Credential{APIKey: "login", APISecret: "token"})
	if err != domain.ErrMissingBaseURL {
		t.Fatalf("err = %v, se esperaba ErrMissingBaseURL: sin URL en la DB debe fallar, nunca caer a la API real", err)
	}
}

func TestInvalidCredentialsMapToDomainError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	client := New()
	_, err := client.GetStoreInfo(context.Background(), domain.Credential{APIKey: "x", APISecret: "y", BaseURL: server.URL})
	if err != domain.ErrInvalidCredentials {
		t.Fatalf("err = %v, se esperaba ErrInvalidCredentials", err)
	}
}

func TestGetProductsCapturaDimensiones(t *testing.T) {
	server, _ := newTestServer(t)
	client := New()

	products, err := client.GetProducts(context.Background(), testCred(server))
	if err != nil {
		t.Fatalf("GetProducts: %v", err)
	}
	if len(products) != 1 {
		t.Fatalf("productos = %d", len(products))
	}

	p := products[0]
	if p.Weight != 1.5 {
		t.Fatalf("weight = %v, se esperaba 1.5", p.Weight)
	}
	if p.Height != 12 || p.Width != 20 || p.Length != 30 {
		t.Fatalf("dimensiones perdidas: height=%v width=%v length=%v", p.Height, p.Width, p.Length)
	}
	if p.Diameter != 5 {
		t.Fatalf("diameter = %v, se esperaba 5", p.Diameter)
	}
	if p.PackageFormat != "box" {
		t.Fatalf("package_format = %q", p.PackageFormat)
	}
}

func TestGetStoreInfoTraeWeightUnit(t *testing.T) {
	server, _ := newTestServer(t)
	client := New()

	info, err := client.GetStoreInfo(context.Background(), testCred(server))
	if err != nil {
		t.Fatalf("GetStoreInfo: %v", err)
	}
	if info.WeightUnit != "kg" {
		t.Fatalf("weight_unit = %q, se esperaba kg: sin esto no podemos convertir el peso", info.WeightUnit)
	}
}
