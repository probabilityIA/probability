package mapper

import (
	"encoding/json"
	"testing"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/shopify/internal/domain"
)

func TestMapShopifyOrderToProbability_BasicFields(t *testing.T) {
	order := &domain.ShopifyOrder{
		ExternalID:  "12345",
		OrderNumber: "#1001",
		TotalAmount: 150000,
		Currency:    "COP",
		Subtotal:    120000,
		Tax:         20000,
		Discount:    5000,
		ShippingCost: 15000,
		Metadata:    map[string]interface{}{},
	}

	result := MapShopifyOrderToProbability(order)

	if result.ExternalID != "12345" {
		t.Errorf("ExternalID: got %q, want %q", result.ExternalID, "12345")
	}
	if result.OrderNumber != "#1001" {
		t.Errorf("OrderNumber: got %q, want %q", result.OrderNumber, "#1001")
	}
	if result.TotalAmount != 150000 {
		t.Errorf("TotalAmount: got %f, want 150000", result.TotalAmount)
	}
	if result.Currency != "COP" {
		t.Errorf("Currency: got %q, want %q", result.Currency, "COP")
	}
	if result.Subtotal != 120000 {
		t.Errorf("Subtotal: got %f, want 120000", result.Subtotal)
	}
	if result.Tax != 20000 {
		t.Errorf("Tax: got %f, want 20000", result.Tax)
	}
}

func TestMapShopifyOrderToProbability_CustomerFields(t *testing.T) {
	order := &domain.ShopifyOrder{
		Customer: domain.ShopifyCustomer{
			Name:  "Juan Perez",
			Email: "juan@test.com",
			Phone: "+573001234567",
		},
		Metadata: map[string]interface{}{},
	}

	result := MapShopifyOrderToProbability(order)

	if result.CustomerName != "Juan Perez" {
		t.Errorf("CustomerName: got %q, want %q", result.CustomerName, "Juan Perez")
	}
	if result.CustomerEmail != "juan@test.com" {
		t.Errorf("CustomerEmail: got %q, want %q", result.CustomerEmail, "juan@test.com")
	}
	if result.CustomerPhone != "+573001234567" {
		t.Errorf("CustomerPhone: got %q, want %q", result.CustomerPhone, "+573001234567")
	}
}

func TestMapShopifyOrderToProbability_ShippingAddress(t *testing.T) {
	order := &domain.ShopifyOrder{
		ShippingAddress: domain.ShopifyAddress{
			Street:     "Calle 100 #15-25",
			Address2:   "Apto 301",
			City:       "Bogota",
			State:      "Cundinamarca",
			Country:    "Colombia",
			PostalCode: "110111",
		},
		Metadata: map[string]interface{}{},
	}

	result := MapShopifyOrderToProbability(order)

	if len(result.Addresses) != 1 {
		t.Fatalf("se esperaba 1 address, se obtuvieron %d", len(result.Addresses))
	}

	addr := result.Addresses[0]
	if addr.Type != "shipping" {
		t.Errorf("Type: got %q, want %q", addr.Type, "shipping")
	}
	if addr.Street != "Calle 100 #15-25" {
		t.Errorf("Street: got %q, want %q", addr.Street, "Calle 100 #15-25")
	}
	if addr.Street2 != "Apto 301" {
		t.Errorf("Street2: got %q, want %q", addr.Street2, "Apto 301")
	}
	if addr.City != "Bogota" {
		t.Errorf("City: got %q, want %q", addr.City, "Bogota")
	}
	if addr.Country != "Colombia" {
		t.Errorf("Country: got %q, want %q", addr.Country, "Colombia")
	}
}

func TestMapShopifyOrderToProbability_Items(t *testing.T) {
	productID := int64(100)
	variantID := int64(200)

	order := &domain.ShopifyOrder{
		Currency: "COP",
		Items: []domain.ShopifyOrderItem{
			{
				Name:      "Camiseta Azul",
				Title:     "Camiseta Azul M",
				SKU:       "CAM-AZUL-M",
				Quantity:  2,
				UnitPrice: 50000,
				Discount:  5000,
				Tax:       9000,
				ProductID: &productID,
				VariantID: &variantID,
			},
		},
		Metadata: map[string]interface{}{},
	}

	result := MapShopifyOrderToProbability(order)

	if len(result.OrderItems) != 1 {
		t.Fatalf("se esperaba 1 item, se obtuvieron %d", len(result.OrderItems))
	}

	item := result.OrderItems[0]
	if item.ProductSKU != "CAM-AZUL-M" {
		t.Errorf("ProductSKU: got %q, want %q", item.ProductSKU, "CAM-AZUL-M")
	}
	if item.Quantity != 2 {
		t.Errorf("Quantity: got %d, want 2", item.Quantity)
	}
	if item.UnitPrice != 50000 {
		t.Errorf("UnitPrice: got %f, want 50000", item.UnitPrice)
	}
	// TotalPrice = (50000 * 2) - 5000 = 95000
	if item.TotalPrice != 95000 {
		t.Errorf("TotalPrice: got %f, want 95000", item.TotalPrice)
	}
	if item.ProductID == nil || *item.ProductID != "100" {
		t.Errorf("ProductID: got %v, want '100'", item.ProductID)
	}
	if item.VariantID == nil || *item.VariantID != "200" {
		t.Errorf("VariantID: got %v, want '200'", item.VariantID)
	}

	// Verificar que Items (JSON) es valido
	var itemsJSON []domain.ProbabilityOrderItemDTO
	if err := json.Unmarshal(result.Items, &itemsJSON); err != nil {
		t.Errorf("Items JSON invalido: %v", err)
	}
}

func TestMapShopifyOrderToProbability_SKUFallbacks(t *testing.T) {
	variantID := int64(300)
	productID := int64(400)

	tests := []struct {
		name        string
		item        domain.ShopifyOrderItem
		expectedSKU string
	}{
		{
			name:        "SKU explicito",
			item:        domain.ShopifyOrderItem{SKU: "MY-SKU"},
			expectedSKU: "MY-SKU",
		},
		{
			name:        "Fallback a VAR-xxx",
			item:        domain.ShopifyOrderItem{SKU: "", VariantID: &variantID},
			expectedSKU: "VAR-300",
		},
		{
			name:        "Fallback a PROD-xxx",
			item:        domain.ShopifyOrderItem{SKU: "", ProductID: &productID},
			expectedSKU: "PROD-400",
		},
		{
			name:        "Fallback a ITEM-0",
			item:        domain.ShopifyOrderItem{SKU: ""},
			expectedSKU: "ITEM-0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			order := &domain.ShopifyOrder{
				Items:    []domain.ShopifyOrderItem{tt.item},
				Metadata: map[string]interface{}{},
			}

			result := MapShopifyOrderToProbability(order)

			if len(result.OrderItems) != 1 {
				t.Fatalf("se esperaba 1 item, se obtuvieron %d", len(result.OrderItems))
			}
			if result.OrderItems[0].ProductSKU != tt.expectedSKU {
				t.Errorf("SKU: got %q, want %q", result.OrderItems[0].ProductSKU, tt.expectedSKU)
			}
		})
	}
}

func TestMapShopifyOrderToProbability_NegativePriceClamped(t *testing.T) {
	order := &domain.ShopifyOrder{
		Items: []domain.ShopifyOrderItem{
			{
				SKU:       "TEST",
				Quantity:  1,
				UnitPrice: 100,
				Discount:  200, // descuento mayor que precio → total negativo
			},
		},
		Metadata: map[string]interface{}{},
	}

	result := MapShopifyOrderToProbability(order)

	if result.OrderItems[0].TotalPrice != 0 {
		t.Errorf("TotalPrice negativo debe ser 0: got %f", result.OrderItems[0].TotalPrice)
	}
}

func TestMapShopifyOrderToProbability_Invoiceable(t *testing.T) {
	// Orden en COP → facturable
	orderCOP := &domain.ShopifyOrder{
		ExternalID: "123",
		Currency:   "COP",
		Metadata:   map[string]interface{}{},
	}
	resultCOP := MapShopifyOrderToProbability(orderCOP)
	if !resultCOP.Invoiceable {
		t.Error("Invoiceable debe ser true para ordenes en COP")
	}

	// Orden en USD → no facturable
	orderUSD := &domain.ShopifyOrder{
		ExternalID: "456",
		Currency:   "USD",
		Metadata:   map[string]interface{}{},
	}
	resultUSD := MapShopifyOrderToProbability(orderUSD)
	if resultUSD.Invoiceable {
		t.Error("Invoiceable debe ser false para ordenes en USD")
	}

	// Orden sin moneda → no facturable
	orderEmpty := &domain.ShopifyOrder{
		ExternalID: "789",
		Metadata:   map[string]interface{}{},
	}
	resultEmpty := MapShopifyOrderToProbability(orderEmpty)
	if resultEmpty.Invoiceable {
		t.Error("Invoiceable debe ser false para ordenes sin moneda definida")
	}

	// Orden dual-currency con presentment COP → facturable
	orderDual := &domain.ShopifyOrder{
		ExternalID:             "101",
		Currency:               "USD",
		CurrencyPresentment:    "COP",
		TotalAmountPresentment: 100000,
		TotalAmount:            25,
		Metadata:               map[string]interface{}{},
	}
	resultDual := MapShopifyOrderToProbability(orderDual)
	if !resultDual.Invoiceable {
		t.Error("Invoiceable debe ser true para ordenes dual-currency con presentment COP")
	}
}

func TestEnrichOrderWithDetails_FinancialStatus(t *testing.T) {
	rawPayload := []byte(`{
		"financial_status": "paid",
		"fulfillment_status": "fulfilled",
		"subtotal_price": "150000.00",
		"total_tax": "28500.00",
		"total_discounts": "5000.00"
	}`)

	order := &domain.ProbabilityOrderDTO{}

	EnrichOrderWithDetails(order, rawPayload)

	// FinancialDetails debe contener informacion extraida
	if order.FinancialDetails == nil {
		t.Fatal("FinancialDetails no debe ser nil con payload valido")
	}

	var financialDetails map[string]interface{}
	if err := json.Unmarshal(order.FinancialDetails, &financialDetails); err != nil {
		t.Fatalf("FinancialDetails no es JSON valido: %v", err)
	}
	if financialDetails["subtotal_price"] != "150000.00" {
		t.Errorf("subtotal_price: got %v, want %q", financialDetails["subtotal_price"], "150000.00")
	}
}

func TestEnrichOrderWithDetails_EmptyPayload(t *testing.T) {
	order := &domain.ProbabilityOrderDTO{}

	// No debe entrar en panic con payload vacio
	EnrichOrderWithDetails(order, nil)
	EnrichOrderWithDetails(order, []byte{})

	if order.FinancialDetails != nil {
		t.Error("FinancialDetails debe ser nil con payload vacio")
	}
}

func TestEnrichOrderWithDetails_InvalidJSON(t *testing.T) {
	order := &domain.ProbabilityOrderDTO{}

	// No debe entrar en panic con JSON invalido
	EnrichOrderWithDetails(order, []byte("not-json{{{"))

	// Los extractors fallan silenciosamente, no asignan nada
	if order.FinancialDetails != nil {
		t.Error("FinancialDetails debe ser nil con JSON invalido")
	}
}
