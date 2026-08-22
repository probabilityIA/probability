package mapper

import (
	"testing"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/tiendanube/internal/domain"
)

func baseOrder() *domain.TiendanubeOrder {
	return &domain.TiendanubeOrder{
		ID:                    "1234",
		Number:                "108",
		Currency:              "COP",
		Status:                "open",
		PaymentStatus:         "paid",
		ShippingStatus:        "unshipped",
		Subtotal:              50000,
		Discount:              5000,
		Total:                 58000,
		ShippingCost:          13000,
		ContactName:           "Ana Perez",
		ContactEmail:          "ana@example.com",
		ContactPhone:          "3001234567",
		ContactIdentification: "1020304050",
		Gateway:               "cash",
		GatewayName:           "Pago contra entrega",
		Items: []domain.TiendanubeOrderItem{
			{ProductID: "77", VariantID: "88", Name: "Camiseta", SKU: "CAM-01", Price: 25000, Quantity: 2},
		},
		ShippingAddress: domain.TiendanubeAddress{
			FirstName: "Ana", LastName: "Perez", Street: "Calle 80", Number: "12-34",
			City: "Bogota", Province: "Cundinamarca", Country: "CO", Zipcode: "110111",
		},
		CreatedAt: "2026-08-20T10:00:00-05:00",
		PaidAt:    "2026-08-20T10:05:00-05:00",
	}
}

func TestMapOrder_CamposBasicos(t *testing.T) {
	dto := MapTiendanubeOrderToProbability(baseOrder(), nil)

	if dto.ExternalID != "1234" {
		t.Errorf("external_id = %q, se esperaba 1234", dto.ExternalID)
	}
	if dto.OrderNumber != "108" {
		t.Errorf("order_number = %q, se esperaba 108", dto.OrderNumber)
	}
	if dto.Platform != "tiendanube" {
		t.Errorf("platform = %q", dto.Platform)
	}
	if dto.TotalAmount != 58000 {
		t.Errorf("total = %v, se esperaba 58000", dto.TotalAmount)
	}
	if dto.CustomerDNI != "1020304050" {
		t.Errorf("customer_dni = %q", dto.CustomerDNI)
	}
	if !dto.Invoiceable {
		t.Error("una orden en COP debe ser facturable")
	}
}

func TestMapOrder_ItemsCalculanTotalPorLinea(t *testing.T) {
	dto := MapTiendanubeOrderToProbability(baseOrder(), nil)

	if len(dto.OrderItems) != 1 {
		t.Fatalf("items = %d, se esperaba 1", len(dto.OrderItems))
	}
	item := dto.OrderItems[0]
	if item.ProductSKU != "CAM-01" {
		t.Errorf("sku = %q", item.ProductSKU)
	}
	if item.Quantity != 2 {
		t.Errorf("cantidad = %d", item.Quantity)
	}
	if item.TotalPrice != 50000 {
		t.Errorf("total de linea = %v, se esperaba 50000 (25000 x 2)", item.TotalPrice)
	}
	if item.VariantID == nil || *item.VariantID != "88" {
		t.Error("se esperaba variant_id 88")
	}
}

func TestMapOrder_ContraEntregaFijaCodTotal(t *testing.T) {
	dto := MapTiendanubeOrderToProbability(baseOrder(), nil)

	if len(dto.Payments) != 1 {
		t.Fatalf("pagos = %d, se esperaba 1", len(dto.Payments))
	}
	if dto.Payments[0].PaymentMethodID != paymentMethodCOD {
		t.Errorf("metodo de pago = %d, se esperaba COD por el gateway 'Pago contra entrega'", dto.Payments[0].PaymentMethodID)
	}
	if dto.CodTotal == nil || *dto.CodTotal != 58000 {
		t.Error("una orden contra entrega debe fijar cod_total con el total")
	}
	if dto.Payments[0].PaidAt == nil {
		t.Error("se esperaba paid_at parseado")
	}
}

func TestMapOrder_SinContraEntregaNoFijaCodTotal(t *testing.T) {
	order := baseOrder()
	order.Gateway = "mercadopago"
	order.GatewayName = "Mercado Pago"

	dto := MapTiendanubeOrderToProbability(order, nil)

	if dto.CodTotal != nil {
		t.Error("una orden que no es contra entrega no debe fijar cod_total")
	}
	if dto.Payments[0].PaymentMethodID != paymentMethodMercadoPago {
		t.Errorf("metodo de pago = %d, se esperaba Mercado Pago", dto.Payments[0].PaymentMethodID)
	}
}

func TestMapOrder_DireccionDeEnvio(t *testing.T) {
	dto := MapTiendanubeOrderToProbability(baseOrder(), nil)

	var shipping *string
	for _, addr := range dto.Addresses {
		if addr.Type == "shipping" {
			street := addr.Street
			shipping = &street
			if addr.City != "Bogota" {
				t.Errorf("ciudad = %q", addr.City)
			}
			if addr.State != "Cundinamarca" {
				t.Errorf("departamento = %q", addr.State)
			}
		}
	}
	if shipping == nil {
		t.Fatal("se esperaba una direccion de envio")
	}
	if *shipping != "Calle 80 12-34" {
		t.Errorf("calle = %q, se esperaba 'Calle 80 12-34' (calle + numero)", *shipping)
	}
}

func TestMapOrderStatus(t *testing.T) {
	cases := []struct {
		status, payment, shipping, want string
	}{
		{"open", "pending", "unshipped", "pending"},
		{"open", "paid", "unshipped", "paid"},
		{"open", "paid", "delivered", "completed"},
		{"cancelled", "paid", "unshipped", "cancelled"},
		{"closed", "paid", "delivered", "completed"},
		{"open", "refunded", "unshipped", "refunded"},
		{"open", "voided", "unshipped", "cancelled"},
	}

	for _, tc := range cases {
		got := MapOrderStatus(tc.status, tc.payment, tc.shipping)
		if got != tc.want {
			t.Errorf("MapOrderStatus(%q,%q,%q) = %q, se esperaba %q", tc.status, tc.payment, tc.shipping, got, tc.want)
		}
	}
}

func TestMapShipmentStatus(t *testing.T) {
	cases := map[string]string{
		"delivered":           "delivered",
		"shipped":             "in_transit",
		"partially_fulfilled": "in_transit",
		"unshipped":           "pending",
		"":                    "pending",
	}
	for input, want := range cases {
		if got := MapShipmentStatus(input); got != want {
			t.Errorf("MapShipmentStatus(%q) = %q, se esperaba %q", input, got, want)
		}
	}
}

func TestMapOrder_NombreDesdeDireccionSiNoHayContacto(t *testing.T) {
	order := baseOrder()
	order.ContactName = ""

	dto := MapTiendanubeOrderToProbability(order, nil)

	if dto.CustomerName != "Ana Perez" {
		t.Errorf("customer_name = %q, se esperaba caer a la direccion de envio", dto.CustomerName)
	}
}
