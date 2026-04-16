package app

import (
	"testing"

	"github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes/internal/domain/ports"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes/internal/mocks"
)

// TestNew verifica que el constructor crea correctamente el use case
// y retorna una instancia que satisface la interfaz ports.IInvoiceUseCase.
func TestNew_ReturnsIInvoiceUseCase(t *testing.T) {
	// Arrange
	mockClient := &mocks.SoftpymesClientMock{}
	mockLogger := mocks.NewLoggerMock()

	// Act
	uc := New(mockClient, mockLogger)

	// Assert
	if uc == nil {
		t.Fatal("New() retornó nil, se esperaba una instancia válida")
	}

	// Verificar que implementa la interfaz
	var _ ports.IInvoiceUseCase = uc
}

// TestNew_WithNilClient verifica que el constructor no paniquea con un cliente nil.
// (Algunos entornos de test pueden pasar nil; el panic se detectaría en tiempo de ejecución
// al invocar métodos, no en la construcción.)
func TestNew_WithNilClientDoesNotPanic(t *testing.T) {
	// Arrange
	mockLogger := mocks.NewLoggerMock()

	// Act / Assert - no debe hacer panic
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("New() provocó un panic inesperado: %v", r)
		}
	}()

	uc := New(nil, mockLogger)
	if uc == nil {
		t.Fatal("New() retornó nil")
	}
}
