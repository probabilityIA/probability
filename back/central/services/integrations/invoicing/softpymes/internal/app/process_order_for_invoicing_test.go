package app

import (
	"context"
	"testing"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes/internal/domain/ports"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes/internal/mocks"
)

// buildUseCase construye el use case con mocks.
func buildUseCase(client *mocks.SoftpymesClientMock) ports.IInvoiceUseCase {
	return New(client, mocks.NewLoggerMock())
}

// buildOrderEvent devuelve un OrderEventMessage mínimo válido para tests.
func buildOrderEvent() *ports.OrderEventMessage {
	integrationID := uint(42)
	businessID := uint(10)
	return &ports.OrderEventMessage{
		EventID:       "evt-001",
		EventType:     "order.paid",
		OrderID:       "order-123",
		BusinessID:    &businessID,
		IntegrationID: &integrationID,
		Timestamp:     time.Now(),
		Order: &ports.OrderSnapshot{
			ID:          "order-123",
			OrderNumber: "ORD-001",
			TotalAmount: 150_000,
			Currency:    "COP",
		},
	}
}

// TestProcessOrderForInvoicing_ReturnsNil verifica que el stub retorna nil.
// La implementación actual es un stub que delega el procesamiento real
// al consumer (InvoiceRequestConsumer) que tiene acceso directo al cliente HTTP.
func TestProcessOrderForInvoicing_ReturnsNil(t *testing.T) {
	// Arrange
	uc := buildUseCase(&mocks.SoftpymesClientMock{})
	ctx := context.Background()
	event := buildOrderEvent()

	// Act
	err := uc.ProcessOrderForInvoicing(ctx, event)

	// Assert
	if err != nil {
		t.Errorf("ProcessOrderForInvoicing() esperaba nil, obtuvo: %v", err)
	}
}

// TestProcessOrderForInvoicing_WithNilEvent verifica que el stub no paniquea
// cuando el evento es nil, ya que la implementación actual no lo consume.
func TestProcessOrderForInvoicing_WithNilEventDoesNotPanic(t *testing.T) {
	// Arrange
	uc := buildUseCase(&mocks.SoftpymesClientMock{})
	ctx := context.Background()

	// Act / Assert
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("ProcessOrderForInvoicing(nil) provocó panic inesperado: %v", r)
		}
	}()

	err := uc.ProcessOrderForInvoicing(ctx, nil)
	if err != nil {
		t.Errorf("ProcessOrderForInvoicing() esperaba nil, obtuvo: %v", err)
	}
}

// TestProcessOrderForInvoicing_DoesNotCallClient verifica que el stub
// nunca invoca el cliente HTTP; el procesamiento real ocurre en el consumer.
func TestProcessOrderForInvoicing_DoesNotCallClient(t *testing.T) {
	// Arrange
	clientCalled := false
	mockClient := &mocks.SoftpymesClientMock{
		CreateInvoiceFn: func(_ context.Context, _ *dtos.CreateInvoiceRequest) (*dtos.CreateInvoiceResult, error) {
			clientCalled = true
			return nil, nil
		},
	}
	uc := buildUseCase(mockClient)
	ctx := context.Background()

	// Act
	_ = uc.ProcessOrderForInvoicing(ctx, buildOrderEvent())

	// Assert
	if clientCalled {
		t.Error("ProcessOrderForInvoicing() no debería llamar al cliente HTTP (es un stub)")
	}
}
