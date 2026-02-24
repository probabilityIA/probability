package usecases

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/canonical"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/woocommerce/internal/mocks"
)

// ─────────────────────────────────────────────────────────────
// ProcessWebhookOrder — camino feliz order.created
// ─────────────────────────────────────────────────────────────

func TestProcessWebhookOrder_Success_OrderCreated(t *testing.T) {
	ctx := context.Background()
	publisher := &mocks.OrderPublisherMock{}

	uc := New(&mocks.WooClientMock{}, &mocks.IntegrationServiceMock{}, publisher, mocks.NewLoggerMock())

	rawBody := buildSampleOrderJSON(t)

	err := uc.ProcessWebhookOrder(ctx, "order.created", "https://mitienda.com", rawBody)

	if err != nil {
		t.Fatalf("esperaba sin error, recibí: %v", err)
	}
	if len(publisher.Published) != 1 {
		t.Fatalf("esperaba 1 orden publicada, recibí %d", len(publisher.Published))
	}

	dto := publisher.Published[0]
	if dto.ExternalID != "123" {
		t.Errorf("ExternalID esperado '123', recibí '%s'", dto.ExternalID)
	}
	if dto.OrderNumber != "1001" {
		t.Errorf("OrderNumber esperado '1001', recibí '%s'", dto.OrderNumber)
	}
	if dto.IntegrationType != "woocommerce" {
		t.Errorf("IntegrationType esperado 'woocommerce', recibí '%s'", dto.IntegrationType)
	}
	if dto.Status != "paid" {
		t.Errorf("Status esperado 'paid' (processing→paid), recibí '%s'", dto.Status)
	}
	if dto.OriginalStatus != "processing" {
		t.Errorf("OriginalStatus esperado 'processing', recibí '%s'", dto.OriginalStatus)
	}
}

// ─────────────────────────────────────────────────────────────
// ProcessWebhookOrder — camino feliz order.updated
// ─────────────────────────────────────────────────────────────

func TestProcessWebhookOrder_Success_OrderUpdated(t *testing.T) {
	ctx := context.Background()
	publisher := &mocks.OrderPublisherMock{}

	uc := New(&mocks.WooClientMock{}, &mocks.IntegrationServiceMock{}, publisher, mocks.NewLoggerMock())

	rawBody := buildSampleOrderJSON(t)

	err := uc.ProcessWebhookOrder(ctx, "order.updated", "https://mitienda.com", rawBody)

	if err != nil {
		t.Fatalf("esperaba sin error, recibí: %v", err)
	}
	if len(publisher.Published) != 1 {
		t.Fatalf("esperaba 1 orden publicada, recibí %d", len(publisher.Published))
	}
}

// ─────────────────────────────────────────────────────────────
// ProcessWebhookOrder — order.deleted es ignorado (no publica)
// ─────────────────────────────────────────────────────────────

func TestProcessWebhookOrder_OrderDeletedSkipped(t *testing.T) {
	ctx := context.Background()
	publisher := &mocks.OrderPublisherMock{}

	uc := New(&mocks.WooClientMock{}, &mocks.IntegrationServiceMock{}, publisher, mocks.NewLoggerMock())

	rawBody := buildSampleOrderJSON(t)

	err := uc.ProcessWebhookOrder(ctx, "order.deleted", "https://mitienda.com", rawBody)

	if err != nil {
		t.Fatalf("esperaba sin error para order.deleted, recibí: %v", err)
	}
	if len(publisher.Published) != 0 {
		t.Fatalf("order.deleted NO debe publicar, se publicaron %d órdenes", len(publisher.Published))
	}
}

// ─────────────────────────────────────────────────────────────
// ProcessWebhookOrder — JSON inválido retorna error
// ─────────────────────────────────────────────────────────────

func TestProcessWebhookOrder_InvalidJSON(t *testing.T) {
	ctx := context.Background()
	publisher := &mocks.OrderPublisherMock{}

	uc := New(&mocks.WooClientMock{}, &mocks.IntegrationServiceMock{}, publisher, mocks.NewLoggerMock())

	err := uc.ProcessWebhookOrder(ctx, "order.created", "https://mitienda.com", []byte("invalid json{"))

	if err == nil {
		t.Fatal("esperaba error por JSON inválido, recibí nil")
	}
	if len(publisher.Published) != 0 {
		t.Fatalf("no debe publicar con JSON inválido, se publicaron %d", len(publisher.Published))
	}
}

// ─────────────────────────────────────────────────────────────
// ProcessWebhookOrder — error al publicar
// ─────────────────────────────────────────────────────────────

func TestProcessWebhookOrder_PublishError(t *testing.T) {
	ctx := context.Background()
	publishErr := errors.New("rabbitmq: connection refused")

	publisher := &mocks.OrderPublisherMock{
		PublishFn: func(_ context.Context, _ *canonical.ProbabilityOrderDTO) error {
			return publishErr
		},
	}

	uc := New(&mocks.WooClientMock{}, &mocks.IntegrationServiceMock{}, publisher, mocks.NewLoggerMock())

	rawBody := buildSampleOrderJSON(t)

	err := uc.ProcessWebhookOrder(ctx, "order.created", "https://mitienda.com", rawBody)

	if err == nil {
		t.Fatal("esperaba error de publicación, recibí nil")
	}
	if !errors.Is(err, publishErr) {
		t.Errorf("esperaba error wrapping '%v', recibí '%v'", publishErr, err)
	}
}

// ─────────────────────────────────────────────────────────────
// ProcessWebhookOrder — order.restored se procesa normalmente
// ─────────────────────────────────────────────────────────────

func TestProcessWebhookOrder_OrderRestored(t *testing.T) {
	ctx := context.Background()
	publisher := &mocks.OrderPublisherMock{}

	uc := New(&mocks.WooClientMock{}, &mocks.IntegrationServiceMock{}, publisher, mocks.NewLoggerMock())

	rawBody := buildSampleOrderJSON(t)

	err := uc.ProcessWebhookOrder(ctx, "order.restored", "https://mitienda.com", rawBody)

	if err != nil {
		t.Fatalf("esperaba sin error para order.restored, recibí: %v", err)
	}
	if len(publisher.Published) != 1 {
		t.Fatalf("order.restored debe publicar, se publicaron %d órdenes", len(publisher.Published))
	}
}

// ─────────────────────────────────────────────────────────────
// ProcessWebhookOrder — verificar mapeo de campos del DTO
// ─────────────────────────────────────────────────────────────

func TestProcessWebhookOrder_VerifyDTOMapping(t *testing.T) {
	ctx := context.Background()
	publisher := &mocks.OrderPublisherMock{}

	uc := New(&mocks.WooClientMock{}, &mocks.IntegrationServiceMock{}, publisher, mocks.NewLoggerMock())

	rawBody := buildSampleOrderJSON(t)

	err := uc.ProcessWebhookOrder(ctx, "order.created", "https://mitienda.com", rawBody)

	if err != nil {
		t.Fatalf("esperaba sin error, recibí: %v", err)
	}

	dto := publisher.Published[0]

	// Verificar customer mapping
	if dto.CustomerName != "Juan Pérez" {
		t.Errorf("CustomerName esperado 'Juan Pérez', recibí '%s'", dto.CustomerName)
	}
	if dto.CustomerEmail != "juan@example.com" {
		t.Errorf("CustomerEmail esperado 'juan@example.com', recibí '%s'", dto.CustomerEmail)
	}
	if dto.CustomerPhone != "+573001234567" {
		t.Errorf("CustomerPhone esperado '+573001234567', recibí '%s'", dto.CustomerPhone)
	}

	// Verificar currency
	if dto.Currency != "COP" {
		t.Errorf("Currency esperado 'COP', recibí '%s'", dto.Currency)
	}

	// Verificar platform
	if dto.Platform != "woocommerce" {
		t.Errorf("Platform esperado 'woocommerce', recibí '%s'", dto.Platform)
	}

	// Verificar total amount
	if dto.TotalAmount != 150000 {
		t.Errorf("TotalAmount esperado 150000, recibí %f", dto.TotalAmount)
	}

	// Verificar order items
	if len(dto.OrderItems) != 1 {
		t.Fatalf("esperaba 1 order item, recibí %d", len(dto.OrderItems))
	}
	if dto.OrderItems[0].ProductName != "Camiseta Test" {
		t.Errorf("ProductName esperado 'Camiseta Test', recibí '%s'", dto.OrderItems[0].ProductName)
	}

	// Verificar addresses
	if len(dto.Addresses) != 2 {
		t.Fatalf("esperaba 2 direcciones (billing+shipping), recibí %d", len(dto.Addresses))
	}

	// Verificar channel metadata
	if dto.ChannelMetadata == nil {
		t.Fatal("ChannelMetadata no debe ser nil cuando hay rawJSON")
	}
	if dto.ChannelMetadata.ChannelSource != "woocommerce" {
		t.Errorf("ChannelSource esperado 'woocommerce', recibí '%s'", dto.ChannelMetadata.ChannelSource)
	}
}

// ─────────────────────────────────────────────────────────────
// Helper: construir JSON de una orden WooCommerce de ejemplo
// ─────────────────────────────────────────────────────────────

func buildSampleOrderJSON(t *testing.T) []byte {
	t.Helper()
	order := map[string]interface{}{
		"id":             123,
		"number":         "1001",
		"status":         "processing",
		"currency":       "COP",
		"date_created":   "2026-02-20T10:00:00",
		"date_modified":  "2026-02-20T10:05:00",
		"total":          "150000.00",
		"total_tax":      "0.00",
		"discount_total": "0.00",
		"discount_tax":   "0.00",
		"shipping_total": "10000.00",
		"shipping_tax":   "0.00",
		"cart_tax":       "0.00",
		"payment_method": "bacs",
		"payment_method_title": "Transferencia bancaria",
		"customer_note":  "Entregar en portería",
		"billing": map[string]interface{}{
			"first_name": "Juan",
			"last_name":  "Pérez",
			"email":      "juan@example.com",
			"phone":      "+573001234567",
			"address_1":  "Calle 100 #10-20",
			"address_2":  "Apto 301",
			"city":       "Bogotá",
			"state":      "DC",
			"postcode":   "110111",
			"country":    "CO",
			"company":    "",
		},
		"shipping": map[string]interface{}{
			"first_name": "Juan",
			"last_name":  "Pérez",
			"address_1":  "Calle 100 #10-20",
			"address_2":  "Apto 301",
			"city":       "Bogotá",
			"state":      "DC",
			"postcode":   "110111",
			"country":    "CO",
			"company":    "",
			"phone":      "+573001234567",
		},
		"line_items": []map[string]interface{}{
			{
				"id":           1,
				"name":         "Camiseta Test",
				"product_id":   42,
				"variation_id": 0,
				"quantity":     2,
				"subtotal":     "140000.00",
				"subtotal_tax": "0.00",
				"total":        "140000.00",
				"total_tax":    "0.00",
				"sku":          "CAM-001",
				"price":        70000,
				"image":        map[string]interface{}{"src": "https://mitienda.com/image.jpg"},
			},
		},
		"shipping_lines": []map[string]interface{}{
			{
				"id":           1,
				"method_title": "Envío estándar",
				"method_id":    "flat_rate",
				"total":        "10000.00",
				"total_tax":    "0.00",
			},
		},
		"fee_lines":    []interface{}{},
		"coupon_lines": []interface{}{},
		"meta_data":    []interface{}{},
	}

	data, err := json.Marshal(order)
	if err != nil {
		t.Fatalf("error marshaling sample order: %v", err)
	}
	return data
}
