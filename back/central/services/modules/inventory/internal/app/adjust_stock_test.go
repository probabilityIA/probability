package app

import (
	"context"
	"errors"
	"testing"

	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/dtos"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/mocks"
)

func TestAdjustStock_CantidadCero_RetornaErrInvalidQuantity(t *testing.T) {
	// Arrange
	repo := &mocks.RepositoryMock{}
	uc := buildUseCase(repo, nil, nil)

	dto := dtos.AdjustStockDTO{
		ProductID:   "prod-001",
		WarehouseID: 1,
		BusinessID:  10,
		Quantity:    0,
		Reason:      "ajuste manual",
	}

	// Act
	result, err := uc.AdjustStock(context.Background(), dto)

	// Assert
	if err == nil {
		t.Fatal("se esperaba error pero se obtuvo nil")
	}
	if !errors.Is(err, domainerrors.ErrInvalidQuantity) {
		t.Errorf("error esperado ErrInvalidQuantity, se obtuvo: %v", err)
	}
	if result != nil {
		t.Errorf("resultado esperado nil, se obtuvo: %v", result)
	}
}

func TestAdjustStock_ProductoNoEncontrado_RetornaErrProductNotFound(t *testing.T) {
	// Arrange
	repo := &mocks.RepositoryMock{
		GetProductByIDFn: func(ctx context.Context, productID string, businessID uint) (string, string, bool, error) {
			return "", "", false, errors.New("not found")
		},
	}
	uc := buildUseCase(repo, nil, nil)

	dto := dtos.AdjustStockDTO{
		ProductID:   "prod-inexistente",
		WarehouseID: 1,
		BusinessID:  10,
		Quantity:    5,
		Reason:      "ajuste manual",
	}

	// Act
	result, err := uc.AdjustStock(context.Background(), dto)

	// Assert
	if !errors.Is(err, domainerrors.ErrProductNotFound) {
		t.Errorf("error esperado ErrProductNotFound, se obtuvo: %v", err)
	}
	if result != nil {
		t.Errorf("resultado esperado nil, se obtuvo: %v", result)
	}
}

func TestAdjustStock_ProductoSinTracking_RetornaErrProductNoTracking(t *testing.T) {
	// Arrange
	repo := &mocks.RepositoryMock{
		GetProductByIDFn: func(ctx context.Context, productID string, businessID uint) (string, string, bool, error) {
			return "Producto", "SKU-001", false, nil // trackInventory = false
		},
	}
	uc := buildUseCase(repo, nil, nil)

	dto := dtos.AdjustStockDTO{
		ProductID:   "prod-001",
		WarehouseID: 1,
		BusinessID:  10,
		Quantity:    5,
		Reason:      "ajuste manual",
	}

	// Act
	result, err := uc.AdjustStock(context.Background(), dto)

	// Assert
	if !errors.Is(err, domainerrors.ErrProductNoTracking) {
		t.Errorf("error esperado ErrProductNoTracking, se obtuvo: %v", err)
	}
	if result != nil {
		t.Errorf("resultado esperado nil, se obtuvo: %v", result)
	}
}

func TestAdjustStock_BodegaNoExiste_RetornaErrWarehouseNotFound(t *testing.T) {
	// Arrange
	repo := &mocks.RepositoryMock{
		GetProductByIDFn: func(ctx context.Context, productID string, businessID uint) (string, string, bool, error) {
			return "Producto", "SKU-001", true, nil
		},
		WarehouseExistsFn: func(ctx context.Context, warehouseID uint, businessID uint) (bool, error) {
			return false, nil
		},
	}
	uc := buildUseCase(repo, nil, nil)

	dto := dtos.AdjustStockDTO{
		ProductID:   "prod-001",
		WarehouseID: 99,
		BusinessID:  10,
		Quantity:    5,
		Reason:      "ajuste manual",
	}

	// Act
	result, err := uc.AdjustStock(context.Background(), dto)

	// Assert
	if !errors.Is(err, domainerrors.ErrWarehouseNotFound) {
		t.Errorf("error esperado ErrWarehouseNotFound, se obtuvo: %v", err)
	}
	if result != nil {
		t.Errorf("resultado esperado nil, se obtuvo: %v", result)
	}
}

func TestAdjustStock_ErrorVerificandoBodega_PropagaError(t *testing.T) {
	// Arrange
	expectedErr := errors.New("error de conexion a base de datos")
	repo := &mocks.RepositoryMock{
		GetProductByIDFn: func(ctx context.Context, productID string, businessID uint) (string, string, bool, error) {
			return "Producto", "SKU-001", true, nil
		},
		WarehouseExistsFn: func(ctx context.Context, warehouseID uint, businessID uint) (bool, error) {
			return false, expectedErr
		},
	}
	uc := buildUseCase(repo, nil, nil)

	dto := dtos.AdjustStockDTO{
		ProductID:   "prod-001",
		WarehouseID: 1,
		BusinessID:  10,
		Quantity:    5,
		Reason:      "ajuste manual",
	}

	// Act
	_, err := uc.AdjustStock(context.Background(), dto)

	// Assert
	if !errors.Is(err, expectedErr) {
		t.Errorf("error esperado %v, se obtuvo: %v", expectedErr, err)
	}
}

func TestAdjustStock_Exitoso_CantidadPositiva_UsaMovimientoInbound(t *testing.T) {
	// Arrange
	codigoMovimientoCapturado := ""
	repo := &mocks.RepositoryMock{
		GetProductByIDFn: func(ctx context.Context, productID string, businessID uint) (string, string, bool, error) {
			return "Producto", "SKU-001", true, nil
		},
		WarehouseExistsFn: func(ctx context.Context, warehouseID uint, businessID uint) (bool, error) {
			return true, nil
		},
		GetMovementTypeIDByCodeFn: func(ctx context.Context, code string) (uint, error) {
			codigoMovimientoCapturado = code
			return 1, nil
		},
	}
	uc := buildUseCase(repo, nil, nil)

	dto := dtos.AdjustStockDTO{
		ProductID:   "prod-001",
		WarehouseID: 1,
		BusinessID:  10,
		Quantity:    10,
		Reason:      "reposicion de stock",
	}

	// Act
	result, err := uc.AdjustStock(context.Background(), dto)

	// Assert
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if result == nil {
		t.Fatal("resultado esperado no nil")
	}
	if codigoMovimientoCapturado != "inbound" {
		t.Errorf("codigo de movimiento esperado 'inbound', se obtuvo '%s'", codigoMovimientoCapturado)
	}
}

func TestAdjustStock_Exitoso_CantidadNegativa_UsaMovimientoOutbound(t *testing.T) {
	// Arrange
	codigoMovimientoCapturado := ""
	repo := &mocks.RepositoryMock{
		GetProductByIDFn: func(ctx context.Context, productID string, businessID uint) (string, string, bool, error) {
			return "Producto", "SKU-001", true, nil
		},
		WarehouseExistsFn: func(ctx context.Context, warehouseID uint, businessID uint) (bool, error) {
			return true, nil
		},
		GetMovementTypeIDByCodeFn: func(ctx context.Context, code string) (uint, error) {
			codigoMovimientoCapturado = code
			return 2, nil
		},
	}
	uc := buildUseCase(repo, nil, nil)

	dto := dtos.AdjustStockDTO{
		ProductID:   "prod-001",
		WarehouseID: 1,
		BusinessID:  10,
		Quantity:    -5,
		Reason:      "merma",
	}

	// Act
	result, err := uc.AdjustStock(context.Background(), dto)

	// Assert
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if result == nil {
		t.Fatal("resultado esperado no nil")
	}
	if codigoMovimientoCapturado != "outbound" {
		t.Errorf("codigo de movimiento esperado 'outbound', se obtuvo '%s'", codigoMovimientoCapturado)
	}
}

func TestAdjustStock_ErrorEnTransaccion_PropagaError(t *testing.T) {
	// Arrange
	expectedErr := errors.New("error en transaccion")
	repo := &mocks.RepositoryMock{
		GetProductByIDFn: func(ctx context.Context, productID string, businessID uint) (string, string, bool, error) {
			return "Producto", "SKU-001", true, nil
		},
		WarehouseExistsFn: func(ctx context.Context, warehouseID uint, businessID uint) (bool, error) {
			return true, nil
		},
		GetMovementTypeIDByCodeFn: func(ctx context.Context, code string) (uint, error) {
			return 1, nil
		},
		AdjustStockTxFn: func(ctx context.Context, params dtos.AdjustStockTxParams) (*dtos.AdjustStockTxResult, error) {
			return nil, expectedErr
		},
	}
	uc := buildUseCase(repo, nil, nil)

	dto := dtos.AdjustStockDTO{
		ProductID:   "prod-001",
		WarehouseID: 1,
		BusinessID:  10,
		Quantity:    5,
		Reason:      "ajuste",
	}

	// Act
	result, err := uc.AdjustStock(context.Background(), dto)

	// Assert
	if !errors.Is(err, expectedErr) {
		t.Errorf("error esperado %v, se obtuvo: %v", expectedErr, err)
	}
	if result != nil {
		t.Errorf("resultado esperado nil, se obtuvo: %v", result)
	}
}

func TestAdjustStock_PublisherNil_NoFalla(t *testing.T) {
	// Arrange: publisher nil debe ser manejado sin panic (fire-and-forget)
	repo := &mocks.RepositoryMock{
		GetProductByIDFn: func(ctx context.Context, productID string, businessID uint) (string, string, bool, error) {
			return "Producto", "SKU-001", true, nil
		},
		WarehouseExistsFn: func(ctx context.Context, warehouseID uint, businessID uint) (bool, error) {
			return true, nil
		},
		GetMovementTypeIDByCodeFn: func(ctx context.Context, code string) (uint, error) {
			return 1, nil
		},
	}

	// Publisher nil: el use case debe manejarlo sin panic
	uc := buildUseCase(repo, nil, nil)

	dto := dtos.AdjustStockDTO{
		ProductID:   "prod-001",
		WarehouseID: 1,
		BusinessID:  10,
		Quantity:    5,
		Reason:      "ajuste manual",
	}

	// Act
	result, err := uc.AdjustStock(context.Background(), dto)

	// Assert: debe funcionar aunque no haya publisher configurado
	if err != nil {
		t.Fatalf("error inesperado con publisher nil: %v", err)
	}
	if result == nil {
		t.Fatal("resultado esperado no nil")
	}
}

// buildUseCase es un helper que construye el UseCase con mocks para tests
func buildUseCase(repo *mocks.RepositoryMock, publisher *mocks.SyncPublisherMock, eventPub *mocks.InventoryEventPublisherMock) IUseCase {
	logger := &mocks.LoggerMock{}
	return New(repo, publisher, eventPub, logger)
}
