package app

import (
	"context"
	"strings"
	"testing"

	"github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/domain/ports"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/mocks"
)

// TestProcessOrderForInvoicing verifica que el método siempre retorna error
// indicando que Siigo no usa el flujo de eventos directos, sino la queue
// dedicada "invoicing.siigo.requests".
func TestProcessOrderForInvoicing_AlwaysReturnsError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	uc := New(
		&mocks.SiigoClientMock{},
		&mocks.IntegrationCoreMock{},
		&mocks.LoggerMock{},
	)

	event := &ports.OrderEventMessage{
		OrderID: "order-123",
	}

	// Act
	err := uc.ProcessOrderForInvoicing(ctx, event)

	// Assert
	if err == nil {
		t.Fatal("se esperaba un error porque Siigo no usa el flujo de eventos directos, se obtuvo nil")
	}
}

// TestProcessOrderForInvoicing_ErrorMessageContainsSiigo verifica que el mensaje de error
// sea descriptivo y mencione la queue correcta.
func TestProcessOrderForInvoicing_ErrorMessageContainsSiigo(t *testing.T) {
	// Arrange
	ctx := context.Background()
	uc := New(
		&mocks.SiigoClientMock{},
		&mocks.IntegrationCoreMock{},
		&mocks.LoggerMock{},
	)

	event := &ports.OrderEventMessage{
		OrderID:   "order-456",
		EventType: "order.paid",
	}

	// Act
	err := uc.ProcessOrderForInvoicing(ctx, event)

	// Assert
	if err == nil {
		t.Fatal("se esperaba un error, se obtuvo nil")
	}

	errMsg := err.Error()
	if !strings.Contains(errMsg, "siigo") {
		t.Errorf("el mensaje de error debería mencionar 'siigo', se obtuvo: %q", errMsg)
	}

	if !strings.Contains(errMsg, "invoicing.siigo.requests") {
		t.Errorf("el mensaje de error debería mencionar la queue 'invoicing.siigo.requests', se obtuvo: %q", errMsg)
	}
}

// TestProcessOrderForInvoicing_WithNilOrderSnapshot verifica que el método maneja
// eventos sin snapshot de orden sin entrar en pánico.
func TestProcessOrderForInvoicing_WithNilOrderSnapshot(t *testing.T) {
	// Arrange
	ctx := context.Background()
	uc := New(
		&mocks.SiigoClientMock{},
		&mocks.IntegrationCoreMock{},
		&mocks.LoggerMock{},
	)

	event := &ports.OrderEventMessage{
		OrderID: "order-789",
		Order:   nil, // Sin snapshot
	}

	// Act — no debe entrar en pánico
	err := uc.ProcessOrderForInvoicing(ctx, event)

	// Assert
	if err == nil {
		t.Fatal("se esperaba un error, se obtuvo nil")
	}
}
