package mapper

import (
	"testing"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/woocommerce/internal/domain"
)

// ─────────────────────────────────────────────────────────────
// MapWooOrderToProbability — orden completa
// ─────────────────────────────────────────────────────────────

func TestMapWooOrderToProbability_FullOrder(t *testing.T) {
	now := time.Now()
	paidAt := now.Add(-1 * time.Hour)

	order := &domain.WooCommerceOrder{
		ID:            12345,
		Number:        "WOO-5001",
		Status:        "processing",
		Currency:      "COP",
		Total:         "250000.00",
		TotalTax:      "39000.00",
		DiscountTotal: "10000.00",
		ShippingTotal: "15000.00",
		DateCreated:   now,
		CustomerNote:  "Entregar en la mañana",
		PaymentMethod: "bacs",
		PaymentMethodTitle: "Transferencia bancaria",
		DatePaid:      &paidAt,
		Billing: domain.WooCommerceBilling{
			FirstName: "María",
			LastName:  "García",
			Email:     "maria@example.com",
			Phone:     "+573009876543",
			Address1:  "Carrera 50 #30-10",
			Address2:  "Oficina 201",
			City:      "Medellín",
			State:     "ANT",
			Postcode:  "050001",
			Country:   "CO",
			Company:   "García SAS",
		},
		Shipping: domain.WooCommerceShipping{
			FirstName: "María",
			LastName:  "García",
			Address1:  "Carrera 50 #30-10",
			City:      "Medellín",
			State:     "ANT",
			Postcode:  "050001",
			Country:   "CO",
		},
		LineItems: []domain.WooCommerceLineItem{
			{
				ID:          1,
				Name:        "Zapatos deportivos",
				ProductID:   100,
				VariationID: 55,
				Quantity:    1,
				Total:       "200000.00",
				TotalTax:    "32000.00",
				Subtotal:    "200000.00",
				SKU:         "ZAP-DEP-001",
				Price:       200000,
				ImageURL:    "https://mitienda.com/zapatos.jpg",
			},
			{
				ID:        2,
				Name:      "Medias deportivas",
				ProductID: 101,
				Quantity:  2,
				Total:     "30000.00",
				TotalTax:  "4800.00",
				Subtotal:  "30000.00",
				SKU:       "MED-001",
				Price:     15000,
			},
		},
		ShippingLines: []domain.WooCommerceShippingLine{
			{
				ID:          1,
				MethodTitle: "Envío express",
				MethodID:    "flat_rate",
				Total:       "15000.00",
			},
		},
		CouponLines: []domain.WooCommerceCouponLine{
			{Code: "PROMO10", Discount: "10000.00"},
		},
	}

	rawJSON := []byte(`{"id":12345,"number":"WOO-5001"}`)

	dto := MapWooOrderToProbability(order, rawJSON)

	// Basic fields
	if dto.ExternalID != "12345" {
		t.Errorf("ExternalID esperado '12345', recibí '%s'", dto.ExternalID)
	}
	if dto.OrderNumber != "WOO-5001" {
		t.Errorf("OrderNumber esperado 'WOO-5001', recibí '%s'", dto.OrderNumber)
	}
	if dto.IntegrationType != "woocommerce" {
		t.Errorf("IntegrationType esperado 'woocommerce', recibí '%s'", dto.IntegrationType)
	}
	if dto.Platform != "woocommerce" {
		t.Errorf("Platform esperado 'woocommerce', recibí '%s'", dto.Platform)
	}

	// Status mapping
	if dto.Status != "paid" {
		t.Errorf("Status esperado 'paid' (processing→paid), recibí '%s'", dto.Status)
	}
	if dto.OriginalStatus != "processing" {
		t.Errorf("OriginalStatus esperado 'processing', recibí '%s'", dto.OriginalStatus)
	}

	// Amounts
	if dto.TotalAmount != 250000 {
		t.Errorf("TotalAmount esperado 250000, recibí %f", dto.TotalAmount)
	}
	if dto.Tax != 39000 {
		t.Errorf("Tax esperado 39000, recibí %f", dto.Tax)
	}
	if dto.Discount != 10000 {
		t.Errorf("Discount esperado 10000, recibí %f", dto.Discount)
	}
	if dto.ShippingCost != 15000 {
		t.Errorf("ShippingCost esperado 15000, recibí %f", dto.ShippingCost)
	}

	// Customer
	if dto.CustomerName != "María García" {
		t.Errorf("CustomerName esperado 'María García', recibí '%s'", dto.CustomerName)
	}
	if dto.CustomerEmail != "maria@example.com" {
		t.Errorf("CustomerEmail esperado 'maria@example.com', recibí '%s'", dto.CustomerEmail)
	}
	if dto.CustomerPhone != "+573009876543" {
		t.Errorf("CustomerPhone esperado '+573009876543', recibí '%s'", dto.CustomerPhone)
	}

	// Notes
	if dto.Notes == nil {
		t.Fatal("Notes no debe ser nil")
	}
	if *dto.Notes != "Entregar en la mañana" {
		t.Errorf("Notes esperado 'Entregar en la mañana', recibí '%s'", *dto.Notes)
	}

	// Coupon
	if dto.Coupon == nil {
		t.Fatal("Coupon no debe ser nil")
	}
	if *dto.Coupon != "PROMO10" {
		t.Errorf("Coupon esperado 'PROMO10', recibí '%s'", *dto.Coupon)
	}

	// Order items
	if len(dto.OrderItems) != 2 {
		t.Fatalf("esperaba 2 order items, recibí %d", len(dto.OrderItems))
	}
	item1 := dto.OrderItems[0]
	if item1.ProductName != "Zapatos deportivos" {
		t.Errorf("Item1.ProductName esperado 'Zapatos deportivos', recibí '%s'", item1.ProductName)
	}
	if item1.Quantity != 1 {
		t.Errorf("Item1.Quantity esperado 1, recibí %d", item1.Quantity)
	}
	if item1.UnitPrice != 200000 {
		t.Errorf("Item1.UnitPrice esperado 200000, recibí %f", item1.UnitPrice)
	}
	if item1.ProductSKU != "ZAP-DEP-001" {
		t.Errorf("Item1.ProductSKU esperado 'ZAP-DEP-001', recibí '%s'", item1.ProductSKU)
	}
	if item1.VariantID == nil {
		t.Fatal("Item1.VariantID no debe ser nil (VariationID=55)")
	}
	if *item1.VariantID != "55" {
		t.Errorf("Item1.VariantID esperado '55', recibí '%s'", *item1.VariantID)
	}
	if item1.ImageURL == nil {
		t.Fatal("Item1.ImageURL no debe ser nil")
	}

	item2 := dto.OrderItems[1]
	if item2.VariantID != nil {
		t.Errorf("Item2.VariantID debe ser nil (VariationID=0)")
	}
	if item2.ImageURL != nil {
		t.Errorf("Item2.ImageURL debe ser nil (no tiene imagen)")
	}

	// Addresses
	if len(dto.Addresses) != 2 {
		t.Fatalf("esperaba 2 direcciones, recibí %d", len(dto.Addresses))
	}
	billing := dto.Addresses[0]
	if billing.Type != "billing" {
		t.Errorf("primera dirección debe ser 'billing', recibí '%s'", billing.Type)
	}
	if billing.City != "Medellín" {
		t.Errorf("billing.City esperado 'Medellín', recibí '%s'", billing.City)
	}
	if billing.Company != "García SAS" {
		t.Errorf("billing.Company esperado 'García SAS', recibí '%s'", billing.Company)
	}

	shipping := dto.Addresses[1]
	if shipping.Type != "shipping" {
		t.Errorf("segunda dirección debe ser 'shipping', recibí '%s'", shipping.Type)
	}

	// Payments
	if len(dto.Payments) != 1 {
		t.Fatalf("esperaba 1 pago, recibí %d", len(dto.Payments))
	}
	payment := dto.Payments[0]
	if payment.Amount != 250000 {
		t.Errorf("payment.Amount esperado 250000, recibí %f", payment.Amount)
	}
	if payment.Gateway == nil || *payment.Gateway != "bacs" {
		t.Errorf("payment.Gateway esperado 'bacs'")
	}
	if payment.Status != "paid" {
		t.Errorf("payment.Status esperado 'paid' (tiene DatePaid), recibí '%s'", payment.Status)
	}
	if payment.PaidAt == nil {
		t.Error("payment.PaidAt no debe ser nil cuando DatePaid está presente")
	}

	// Shipments
	if len(dto.Shipments) != 1 {
		t.Fatalf("esperaba 1 shipment, recibí %d", len(dto.Shipments))
	}
	shipment := dto.Shipments[0]
	if shipment.Carrier == nil || *shipment.Carrier != "Envío express" {
		t.Errorf("shipment.Carrier esperado 'Envío express'")
	}

	// Channel metadata
	if dto.ChannelMetadata == nil {
		t.Fatal("ChannelMetadata no debe ser nil")
	}
	if dto.ChannelMetadata.ChannelSource != "woocommerce" {
		t.Errorf("ChannelSource esperado 'woocommerce', recibí '%s'", dto.ChannelMetadata.ChannelSource)
	}
	if dto.ChannelMetadata.Version != "v3" {
		t.Errorf("Version esperado 'v3', recibí '%s'", dto.ChannelMetadata.Version)
	}
	if !dto.ChannelMetadata.IsLatest {
		t.Error("IsLatest debe ser true")
	}
	if dto.ChannelMetadata.SyncStatus != "synced" {
		t.Errorf("SyncStatus esperado 'synced', recibí '%s'", dto.ChannelMetadata.SyncStatus)
	}
}

// ─────────────────────────────────────────────────────────────
// MapWooOrderToProbability — orden mínima (sin items opcionales)
// ─────────────────────────────────────────────────────────────

func TestMapWooOrderToProbability_MinimalOrder(t *testing.T) {
	order := &domain.WooCommerceOrder{
		ID:       1,
		Number:   "1",
		Status:   "pending",
		Currency: "USD",
		Total:    "50.00",
	}

	dto := MapWooOrderToProbability(order, nil)

	if dto.ExternalID != "1" {
		t.Errorf("ExternalID esperado '1', recibí '%s'", dto.ExternalID)
	}
	if dto.Status != "pending" {
		t.Errorf("Status esperado 'pending', recibí '%s'", dto.Status)
	}
	if len(dto.OrderItems) != 0 {
		t.Errorf("esperaba 0 items, recibí %d", len(dto.OrderItems))
	}
	if dto.Notes != nil {
		t.Errorf("Notes debe ser nil para orden sin nota")
	}
	if dto.Coupon != nil {
		t.Errorf("Coupon debe ser nil sin cupones")
	}
	if len(dto.Payments) != 0 {
		t.Errorf("esperaba 0 pagos (sin payment_method), recibí %d", len(dto.Payments))
	}
	if dto.ChannelMetadata != nil {
		t.Errorf("ChannelMetadata debe ser nil sin rawJSON")
	}
}

// ─────────────────────────────────────────────────────────────
// mapWooStatus — tabla de mapeos de estado
// ─────────────────────────────────────────────────────────────

func TestMapWooStatus_AllStatuses(t *testing.T) {
	tests := []struct {
		wooStatus     string
		expectedStatus string
	}{
		{"pending", "pending"},
		{"processing", "paid"},
		{"on-hold", "on_hold"},
		{"completed", "fulfilled"},
		{"cancelled", "cancelled"},
		{"refunded", "refunded"},
		{"failed", "failed"},
		{"trash", "deleted"},
		{"unknown_status", "unknown_status"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.wooStatus, func(t *testing.T) {
			result := mapWooStatus(tt.wooStatus)
			if result != tt.expectedStatus {
				t.Errorf("mapWooStatus(%s) = '%s', esperaba '%s'", tt.wooStatus, result, tt.expectedStatus)
			}
		})
	}
}

// ─────────────────────────────────────────────────────────────
// mapShipmentStatus — tabla de mapeos de estado de envío
// ─────────────────────────────────────────────────────────────

func TestMapShipmentStatus_AllStatuses(t *testing.T) {
	tests := []struct {
		wooStatus      string
		expectedStatus string
	}{
		{"completed", "delivered"},
		{"processing", "pending"},
		{"on-hold", "pending"},
		{"cancelled", "cancelled"},
		{"refunded", "cancelled"},
		{"failed", "cancelled"},
		{"pending", "pending"},
		{"unknown", "pending"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.wooStatus, func(t *testing.T) {
			result := mapShipmentStatus(tt.wooStatus)
			if result != tt.expectedStatus {
				t.Errorf("mapShipmentStatus(%s) = '%s', esperaba '%s'", tt.wooStatus, result, tt.expectedStatus)
			}
		})
	}
}

// ─────────────────────────────────────────────────────────────
// MapWooOrderToProbability — pago sin DatePaid (status pending)
// ─────────────────────────────────────────────────────────────

func TestMapWooOrderToProbability_PaymentPending(t *testing.T) {
	order := &domain.WooCommerceOrder{
		ID:            1,
		Number:        "1",
		Status:        "pending",
		Currency:      "COP",
		Total:         "100000.00",
		PaymentMethod: "cod",
		// DatePaid es nil
	}

	dto := MapWooOrderToProbability(order, nil)

	if len(dto.Payments) != 1 {
		t.Fatalf("esperaba 1 pago, recibí %d", len(dto.Payments))
	}
	if dto.Payments[0].Status != "pending" {
		t.Errorf("payment.Status esperado 'pending' (sin DatePaid), recibí '%s'", dto.Payments[0].Status)
	}
	if dto.Payments[0].PaidAt != nil {
		t.Errorf("payment.PaidAt debe ser nil cuando DatePaid es nil")
	}
}

// ─────────────────────────────────────────────────────────────
// MapWooOrderToProbability — múltiples cupones
// ─────────────────────────────────────────────────────────────

func TestMapWooOrderToProbability_MultipleCoupons(t *testing.T) {
	order := &domain.WooCommerceOrder{
		ID:       1,
		Number:   "1",
		Status:   "processing",
		Currency: "COP",
		Total:    "80000.00",
		CouponLines: []domain.WooCommerceCouponLine{
			{Code: "DESCUENTO10"},
			{Code: "ENVIOGRATIS"},
		},
	}

	dto := MapWooOrderToProbability(order, nil)

	if dto.Coupon == nil {
		t.Fatal("Coupon no debe ser nil con 2 cupones")
	}
	if *dto.Coupon != "DESCUENTO10, ENVIOGRATIS" {
		t.Errorf("Coupon esperado 'DESCUENTO10, ENVIOGRATIS', recibí '%s'", *dto.Coupon)
	}
}

// ─────────────────────────────────────────────────────────────
// parseFloat — conversión de strings a float
// ─────────────────────────────────────────────────────────────

func TestParseFloat(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"100.50", 100.50},
		{"0.00", 0},
		{"", 0},
		{"abc", 0},
		{"1000000.99", 1000000.99},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.input, func(t *testing.T) {
			result := parseFloat(tt.input)
			if result != tt.expected {
				t.Errorf("parseFloat(%s) = %f, esperaba %f", tt.input, result, tt.expected)
			}
		})
	}
}
