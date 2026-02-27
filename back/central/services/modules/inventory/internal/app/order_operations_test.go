package app

import (
	"context"
	"errors"
	"testing"

	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/dtos"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/mocks"
)

// buildUseCaseWithEvents crea un UseCase con publisher de eventos no-nil para tests que usan
// operaciones de orden. Las operaciones de orden (reserve, confirm, release, return) publican
// eventos en goroutines concurrentes, por lo que requieren un eventPublisher concreto.
func buildUseCaseWithEvents(repo *mocks.RepositoryMock) IUseCase {
	logger := &mocks.LoggerMock{}
	eventPub := &mocks.InventoryEventPublisherMock{}
	return New(repo, nil, eventPub, logger)
}

// -----------------------------------------------------------------------
// ReserveStockForOrder
// -----------------------------------------------------------------------

func TestReserveStockForOrder_SinWarehouseExplicito_UsaDefaultWarehouse(t *testing.T) {
	// Arrange
	warehouseIDUsada := uint(0)
	repo := &mocks.RepositoryMock{
		GetDefaultWarehouseIDFn: func(ctx context.Context, businessID uint) (uint, error) {
			return 5, nil // bodega por defecto
		},
		GetMovementTypeIDByCodeFn: func(ctx context.Context, code string) (uint, error) {
			return 1, nil
		},
		GetProductByIDFn: func(ctx context.Context, productID string, businessID uint) (string, string, bool, error) {
			return "Producto", "SKU-001", true, nil
		},
		ReserveStockTxFn: func(ctx context.Context, params dtos.ReserveStockTxParams) (*dtos.ReserveStockTxResult, error) {
			warehouseIDUsada = params.WarehouseID
			return &dtos.ReserveStockTxResult{
				Reserved:   params.Quantity,
				Sufficient: true,
			}, nil
		},
	}
	uc := buildUseCaseWithEvents(repo)

	items := []dtos.OrderInventoryItem{
		{ProductID: "prod-001", SKU: "SKU-001", Quantity: 2},
	}

	// Act
	result, err := uc.ReserveStockForOrder(context.Background(), "ORDER-001", 10, nil, items)

	// Assert
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if warehouseIDUsada != 5 {
		t.Errorf("bodega default esperada 5, se uso %d", warehouseIDUsada)
	}
	if !result.Success {
		t.Error("resultado esperado exitoso")
	}
}

func TestReserveStockForOrder_SinBodegaDefault_RetornaErrNoDefaultWarehouse(t *testing.T) {
	// Arrange
	// Este test retorna antes de llamar a publishEvent (falla en resolveWarehouse),
	// por lo que no hay riesgo de goroutines huerfanas. Usamos buildUseCase directamente.
	repo := &mocks.RepositoryMock{
		GetDefaultWarehouseIDFn: func(ctx context.Context, businessID uint) (uint, error) {
			return 0, errors.New("no default warehouse")
		},
		GetMovementTypeIDByCodeFn: func(ctx context.Context, code string) (uint, error) {
			return 1, nil
		},
	}
	uc := buildUseCaseWithEvents(repo)

	items := []dtos.OrderInventoryItem{
		{ProductID: "prod-001", SKU: "SKU-001", Quantity: 2},
	}

	// Act
	result, err := uc.ReserveStockForOrder(context.Background(), "ORDER-001", 10, nil, items)

	// Assert
	if !errors.Is(err, domainerrors.ErrNoDefaultWarehouse) {
		t.Errorf("error esperado ErrNoDefaultWarehouse, se obtuvo: %v", err)
	}
	if result != nil {
		t.Errorf("resultado esperado nil, se obtuvo: %v", result)
	}
}

func TestReserveStockForOrder_ProductoSinTracking_SeProcessaSinReserva(t *testing.T) {
	// Arrange
	warehouseID := uint(1)
	reservado := false
	repo := &mocks.RepositoryMock{
		GetMovementTypeIDByCodeFn: func(ctx context.Context, code string) (uint, error) {
			return 1, nil
		},
		GetProductByIDFn: func(ctx context.Context, productID string, businessID uint) (string, string, bool, error) {
			return "Producto sin tracking", "SKU-002", false, nil // sin tracking
		},
		ReserveStockTxFn: func(ctx context.Context, params dtos.ReserveStockTxParams) (*dtos.ReserveStockTxResult, error) {
			reservado = true
			return nil, nil
		},
	}
	uc := buildUseCaseWithEvents(repo)

	items := []dtos.OrderInventoryItem{
		{ProductID: "prod-sin-tracking", SKU: "SKU-002", Quantity: 5},
	}

	// Act
	result, err := uc.ReserveStockForOrder(context.Background(), "ORDER-002", 10, &warehouseID, items)

	// Assert
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if reservado {
		t.Error("no deberia haberse ejecutado ReserveStockTx para producto sin tracking")
	}
	if len(result.ItemResults) != 1 {
		t.Errorf("cantidad de resultados esperada 1, se obtuvo %d", len(result.ItemResults))
	}
	if !result.ItemResults[0].Sufficient {
		t.Error("item sin tracking debe marcarse como sufficient=true (no afecta inventario)")
	}
}

func TestReserveStockForOrder_StockInsuficiente_ResultadoConSufficientFalse(t *testing.T) {
	// Arrange
	warehouseID := uint(1)
	repo := &mocks.RepositoryMock{
		GetMovementTypeIDByCodeFn: func(ctx context.Context, code string) (uint, error) {
			return 1, nil
		},
		GetProductByIDFn: func(ctx context.Context, productID string, businessID uint) (string, string, bool, error) {
			return "Producto", "SKU-001", true, nil
		},
		ReserveStockTxFn: func(ctx context.Context, params dtos.ReserveStockTxParams) (*dtos.ReserveStockTxResult, error) {
			return &dtos.ReserveStockTxResult{
				Reserved:   2, // solo pudo reservar 2 de 10 solicitados
				Sufficient: false,
			}, nil
		},
	}
	uc := buildUseCaseWithEvents(repo)

	items := []dtos.OrderInventoryItem{
		{ProductID: "prod-001", SKU: "SKU-001", Quantity: 10},
	}

	// Act
	result, err := uc.ReserveStockForOrder(context.Background(), "ORDER-003", 10, &warehouseID, items)

	// Assert
	if err != nil {
		t.Fatalf("error inesperado (stock insuficiente no es error, es resultado parcial): %v", err)
	}
	if len(result.ItemResults) != 1 {
		t.Fatalf("cantidad de resultados esperada 1, se obtuvo %d", len(result.ItemResults))
	}
	if result.ItemResults[0].Sufficient {
		t.Error("item con stock insuficiente debe tener Sufficient=false")
	}
}

func TestReserveStockForOrder_ConWarehouseExplicito_UsaWarehouseProvisto(t *testing.T) {
	// Arrange
	warehouseID := uint(7)
	warehouseIDUsada := uint(0)
	repo := &mocks.RepositoryMock{
		GetMovementTypeIDByCodeFn: func(ctx context.Context, code string) (uint, error) {
			return 1, nil
		},
		GetProductByIDFn: func(ctx context.Context, productID string, businessID uint) (string, string, bool, error) {
			return "Producto", "SKU-001", true, nil
		},
		ReserveStockTxFn: func(ctx context.Context, params dtos.ReserveStockTxParams) (*dtos.ReserveStockTxResult, error) {
			warehouseIDUsada = params.WarehouseID
			return &dtos.ReserveStockTxResult{Reserved: params.Quantity, Sufficient: true}, nil
		},
	}
	uc := buildUseCaseWithEvents(repo)

	items := []dtos.OrderInventoryItem{
		{ProductID: "prod-001", SKU: "SKU-001", Quantity: 3},
	}

	// Act
	_, err := uc.ReserveStockForOrder(context.Background(), "ORDER-004", 10, &warehouseID, items)

	// Assert
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if warehouseIDUsada != 7 {
		t.Errorf("bodega esperada 7, se uso %d", warehouseIDUsada)
	}
}

// -----------------------------------------------------------------------
// ConfirmSaleForOrder
// -----------------------------------------------------------------------

func TestConfirmSaleForOrder_Exitoso_RetornaResultadoExitoso(t *testing.T) {
	// Arrange
	warehouseID := uint(1)
	confirmado := false
	repo := &mocks.RepositoryMock{
		GetMovementTypeIDByCodeFn: func(ctx context.Context, code string) (uint, error) {
			return 2, nil
		},
		GetProductByIDFn: func(ctx context.Context, productID string, businessID uint) (string, string, bool, error) {
			return "Producto", "SKU-001", true, nil
		},
		ConfirmSaleTxFn: func(ctx context.Context, params dtos.ConfirmSaleTxParams) error {
			confirmado = true
			return nil
		},
	}
	uc := buildUseCaseWithEvents(repo)

	items := []dtos.OrderInventoryItem{
		{ProductID: "prod-001", SKU: "SKU-001", Quantity: 2},
	}

	// Act
	result, err := uc.ConfirmSaleForOrder(context.Background(), "ORDER-005", 10, &warehouseID, items)

	// Assert
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if !confirmado {
		t.Error("se esperaba que se ejecutara ConfirmSaleTx")
	}
	if !result.Success {
		t.Error("resultado esperado exitoso")
	}
	if result.ItemResults[0].Processed != 2 {
		t.Errorf("procesados esperados 2, se obtuvo %d", result.ItemResults[0].Processed)
	}
}

func TestConfirmSaleForOrder_ErrorEnTx_MarcaItemComoFallido(t *testing.T) {
	// Arrange
	warehouseID := uint(1)
	repo := &mocks.RepositoryMock{
		GetMovementTypeIDByCodeFn: func(ctx context.Context, code string) (uint, error) {
			return 2, nil
		},
		GetProductByIDFn: func(ctx context.Context, productID string, businessID uint) (string, string, bool, error) {
			return "Producto", "SKU-001", true, nil
		},
		ConfirmSaleTxFn: func(ctx context.Context, params dtos.ConfirmSaleTxParams) error {
			return errors.New("error en confirmacion")
		},
	}
	uc := buildUseCaseWithEvents(repo)

	items := []dtos.OrderInventoryItem{
		{ProductID: "prod-001", SKU: "SKU-001", Quantity: 2},
	}

	// Act
	result, err := uc.ConfirmSaleForOrder(context.Background(), "ORDER-006", 10, &warehouseID, items)

	// Assert
	if err != nil {
		t.Fatalf("error inesperado (errores de item no son errores de la operacion): %v", err)
	}
	if result.Success {
		t.Error("resultado esperado no exitoso cuando hay errores en items")
	}
	if result.ItemResults[0].ErrorMessage == "" {
		t.Error("se esperaba mensaje de error en el item")
	}
}

// -----------------------------------------------------------------------
// ReleaseStockForOrder
// -----------------------------------------------------------------------

func TestReleaseStockForOrder_Exitoso_LiberaReserva(t *testing.T) {
	// Arrange
	warehouseID := uint(1)
	liberado := false
	repo := &mocks.RepositoryMock{
		GetMovementTypeIDByCodeFn: func(ctx context.Context, code string) (uint, error) {
			return 4, nil
		},
		GetProductByIDFn: func(ctx context.Context, productID string, businessID uint) (string, string, bool, error) {
			return "Producto", "SKU-001", true, nil
		},
		ReleaseStockTxFn: func(ctx context.Context, params dtos.ReleaseTxParams) error {
			liberado = true
			return nil
		},
	}
	uc := buildUseCaseWithEvents(repo)

	items := []dtos.OrderInventoryItem{
		{ProductID: "prod-001", SKU: "SKU-001", Quantity: 3},
	}

	// Act
	result, err := uc.ReleaseStockForOrder(context.Background(), "ORDER-007", 10, &warehouseID, items)

	// Assert
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if !liberado {
		t.Error("se esperaba que se ejecutara ReleaseStockTx")
	}
	if !result.Success {
		t.Error("resultado esperado exitoso")
	}
}

// -----------------------------------------------------------------------
// ReturnStockForOrder
// -----------------------------------------------------------------------

func TestReturnStockForOrder_Exitoso_DevuelveStock(t *testing.T) {
	// Arrange
	warehouseID := uint(1)
	devuelto := false
	repo := &mocks.RepositoryMock{
		GetMovementTypeIDByCodeFn: func(ctx context.Context, code string) (uint, error) {
			return 5, nil
		},
		GetProductByIDFn: func(ctx context.Context, productID string, businessID uint) (string, string, bool, error) {
			return "Producto", "SKU-001", true, nil
		},
		ReturnStockTxFn: func(ctx context.Context, params dtos.ReturnStockTxParams) error {
			devuelto = true
			return nil
		},
	}
	uc := buildUseCaseWithEvents(repo)

	items := []dtos.OrderInventoryItem{
		{ProductID: "prod-001", SKU: "SKU-001", Quantity: 1},
	}

	// Act
	result, err := uc.ReturnStockForOrder(context.Background(), "ORDER-008", 10, &warehouseID, items)

	// Assert
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if !devuelto {
		t.Error("se esperaba que se ejecutara ReturnStockTx")
	}
	if !result.Success {
		t.Error("resultado esperado exitoso")
	}
}

func TestReturnStockForOrder_ProductoSinTracking_SkipSinError(t *testing.T) {
	// Arrange
	warehouseID := uint(1)
	repo := &mocks.RepositoryMock{
		GetMovementTypeIDByCodeFn: func(ctx context.Context, code string) (uint, error) {
			return 5, nil
		},
		GetProductByIDFn: func(ctx context.Context, productID string, businessID uint) (string, string, bool, error) {
			return "Servicio", "SRV-001", false, nil // sin tracking -> skip
		},
	}
	uc := buildUseCaseWithEvents(repo)

	items := []dtos.OrderInventoryItem{
		{ProductID: "srv-001", SKU: "SRV-001", Quantity: 1},
	}

	// Act
	result, err := uc.ReturnStockForOrder(context.Background(), "ORDER-009", 10, &warehouseID, items)

	// Assert
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if !result.Success {
		t.Error("resultado esperado exitoso incluso para productos sin tracking")
	}
	if len(result.ItemResults) != 1 {
		t.Fatalf("cantidad de items esperada 1, se obtuvo %d", len(result.ItemResults))
	}
	if !result.ItemResults[0].Sufficient {
		t.Error("item sin tracking debe tener Sufficient=true")
	}
}

func TestReturnStockForOrder_ErrorEnTx_MarcaItemComoFallido(t *testing.T) {
	// Arrange
	warehouseID := uint(1)
	repo := &mocks.RepositoryMock{
		GetMovementTypeIDByCodeFn: func(ctx context.Context, code string) (uint, error) {
			return 5, nil
		},
		GetProductByIDFn: func(ctx context.Context, productID string, businessID uint) (string, string, bool, error) {
			return "Producto", "SKU-001", true, nil
		},
		ReturnStockTxFn: func(ctx context.Context, params dtos.ReturnStockTxParams) error {
			return errors.New("error al devolver stock")
		},
	}
	uc := buildUseCaseWithEvents(repo)

	items := []dtos.OrderInventoryItem{
		{ProductID: "prod-001", SKU: "SKU-001", Quantity: 1},
	}

	// Act
	result, err := uc.ReturnStockForOrder(context.Background(), "ORDER-010", 10, &warehouseID, items)

	// Assert
	if err != nil {
		t.Fatalf("error de item no debe ser error de la operacion: %v", err)
	}
	if result.Success {
		t.Error("resultado esperado no exitoso cuando hay errores en items")
	}
}
