package app

import (
	"context"
	"errors"
	"testing"

	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/dtos"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/mocks"
)

func TestTransferStock_CantidadCeroONegativa_RetornaErrTransferQtyNeg(t *testing.T) {
	tests := []struct {
		name     string
		cantidad int
	}{
		{"cantidad cero", 0},
		{"cantidad negativa", -5},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			repo := &mocks.RepositoryMock{}
			uc := buildUseCase(repo, nil, nil)

			dto := dtos.TransferStockDTO{
				ProductID:       "prod-001",
				FromWarehouseID: 1,
				ToWarehouseID:   2,
				BusinessID:      10,
				Quantity:        tc.cantidad,
			}

			// Act
			err := uc.TransferStock(context.Background(), dto)

			// Assert
			if !errors.Is(err, domainerrors.ErrTransferQtyNeg) {
				t.Errorf("error esperado ErrTransferQtyNeg, se obtuvo: %v", err)
			}
		})
	}
}

func TestTransferStock_MismaBodegaOrigenDestino_RetornaErrSameWarehouse(t *testing.T) {
	// Arrange
	repo := &mocks.RepositoryMock{}
	uc := buildUseCase(repo, nil, nil)

	dto := dtos.TransferStockDTO{
		ProductID:       "prod-001",
		FromWarehouseID: 1,
		ToWarehouseID:   1, // misma bodega
		BusinessID:      10,
		Quantity:        10,
	}

	// Act
	err := uc.TransferStock(context.Background(), dto)

	// Assert
	if !errors.Is(err, domainerrors.ErrSameWarehouse) {
		t.Errorf("error esperado ErrSameWarehouse, se obtuvo: %v", err)
	}
}

func TestTransferStock_ProductoNoExiste_RetornaErrProductNotFound(t *testing.T) {
	// Arrange
	repo := &mocks.RepositoryMock{
		GetProductByIDFn: func(ctx context.Context, productID string, businessID uint) (string, string, bool, error) {
			return "", "", false, errors.New("not found")
		},
	}
	uc := buildUseCase(repo, nil, nil)

	dto := dtos.TransferStockDTO{
		ProductID:       "prod-inexistente",
		FromWarehouseID: 1,
		ToWarehouseID:   2,
		BusinessID:      10,
		Quantity:        10,
	}

	// Act
	err := uc.TransferStock(context.Background(), dto)

	// Assert
	if !errors.Is(err, domainerrors.ErrProductNotFound) {
		t.Errorf("error esperado ErrProductNotFound, se obtuvo: %v", err)
	}
}

func TestTransferStock_ProductoSinTracking_RetornaErrProductNoTracking(t *testing.T) {
	// Arrange
	repo := &mocks.RepositoryMock{
		GetProductByIDFn: func(ctx context.Context, productID string, businessID uint) (string, string, bool, error) {
			return "Producto", "SKU-001", false, nil // sin tracking
		},
	}
	uc := buildUseCase(repo, nil, nil)

	dto := dtos.TransferStockDTO{
		ProductID:       "prod-001",
		FromWarehouseID: 1,
		ToWarehouseID:   2,
		BusinessID:      10,
		Quantity:        10,
	}

	// Act
	err := uc.TransferStock(context.Background(), dto)

	// Assert
	if !errors.Is(err, domainerrors.ErrProductNoTracking) {
		t.Errorf("error esperado ErrProductNoTracking, se obtuvo: %v", err)
	}
}

func TestTransferStock_BodegaOrigenNoExiste_RetornaErrWarehouseNotFound(t *testing.T) {
	// Arrange
	repo := &mocks.RepositoryMock{
		GetProductByIDFn: func(ctx context.Context, productID string, businessID uint) (string, string, bool, error) {
			return "Producto", "SKU-001", true, nil
		},
		WarehouseExistsFn: func(ctx context.Context, warehouseID uint, businessID uint) (bool, error) {
			// La bodega de origen (ID=99) no existe
			if warehouseID == 99 {
				return false, nil
			}
			return true, nil
		},
	}
	uc := buildUseCase(repo, nil, nil)

	dto := dtos.TransferStockDTO{
		ProductID:       "prod-001",
		FromWarehouseID: 99, // no existe
		ToWarehouseID:   2,
		BusinessID:      10,
		Quantity:        10,
	}

	// Act
	err := uc.TransferStock(context.Background(), dto)

	// Assert
	if !errors.Is(err, domainerrors.ErrWarehouseNotFound) {
		t.Errorf("error esperado ErrWarehouseNotFound, se obtuvo: %v", err)
	}
}

func TestTransferStock_BodegaDestinoNoExiste_RetornaErrWarehouseNotFound(t *testing.T) {
	// Arrange
	repo := &mocks.RepositoryMock{
		GetProductByIDFn: func(ctx context.Context, productID string, businessID uint) (string, string, bool, error) {
			return "Producto", "SKU-001", true, nil
		},
		WarehouseExistsFn: func(ctx context.Context, warehouseID uint, businessID uint) (bool, error) {
			// La bodega destino (ID=88) no existe
			if warehouseID == 88 {
				return false, nil
			}
			return true, nil
		},
	}
	uc := buildUseCase(repo, nil, nil)

	dto := dtos.TransferStockDTO{
		ProductID:       "prod-001",
		FromWarehouseID: 1,
		ToWarehouseID:   88, // no existe
		BusinessID:      10,
		Quantity:        10,
	}

	// Act
	err := uc.TransferStock(context.Background(), dto)

	// Assert
	if !errors.Is(err, domainerrors.ErrWarehouseNotFound) {
		t.Errorf("error esperado ErrWarehouseNotFound, se obtuvo: %v", err)
	}
}

func TestTransferStock_Exitoso_RetornaNil(t *testing.T) {
	// Arrange
	txParams := dtos.TransferStockTxParams{}
	repo := &mocks.RepositoryMock{
		GetProductByIDFn: func(ctx context.Context, productID string, businessID uint) (string, string, bool, error) {
			return "Producto", "SKU-001", true, nil
		},
		WarehouseExistsFn: func(ctx context.Context, warehouseID uint, businessID uint) (bool, error) {
			return true, nil
		},
		GetMovementTypeIDByCodeFn: func(ctx context.Context, code string) (uint, error) {
			return 3, nil // ID del tipo "transfer"
		},
		TransferStockTxFn: func(ctx context.Context, params dtos.TransferStockTxParams) (*dtos.TransferStockTxResult, error) {
			txParams = params
			return &dtos.TransferStockTxResult{
				FromNewQty: 40,
				ToNewQty:   60,
			}, nil
		},
	}
	uc := buildUseCase(repo, nil, nil)

	dto := dtos.TransferStockDTO{
		ProductID:       "prod-001",
		FromWarehouseID: 1,
		ToWarehouseID:   2,
		BusinessID:      10,
		Quantity:        10,
		Reason:          "traslado entre bodegas",
	}

	// Act
	err := uc.TransferStock(context.Background(), dto)

	// Assert
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if txParams.Quantity != 10 {
		t.Errorf("cantidad en transaccion esperada 10, se obtuvo %d", txParams.Quantity)
	}
	if txParams.FromWarehouseID != 1 {
		t.Errorf("bodega origen esperada 1, se obtuvo %d", txParams.FromWarehouseID)
	}
	if txParams.ToWarehouseID != 2 {
		t.Errorf("bodega destino esperada 2, se obtuvo %d", txParams.ToWarehouseID)
	}
	if txParams.ReferenceType != "manual" {
		t.Errorf("reference_type esperado 'manual', se obtuvo '%s'", txParams.ReferenceType)
	}
}

func TestTransferStock_ErrorEnTransaccion_PropagaError(t *testing.T) {
	// Arrange
	expectedErr := domainerrors.ErrInsufficientStock
	repo := &mocks.RepositoryMock{
		GetProductByIDFn: func(ctx context.Context, productID string, businessID uint) (string, string, bool, error) {
			return "Producto", "SKU-001", true, nil
		},
		WarehouseExistsFn: func(ctx context.Context, warehouseID uint, businessID uint) (bool, error) {
			return true, nil
		},
		GetMovementTypeIDByCodeFn: func(ctx context.Context, code string) (uint, error) {
			return 3, nil
		},
		TransferStockTxFn: func(ctx context.Context, params dtos.TransferStockTxParams) (*dtos.TransferStockTxResult, error) {
			return nil, expectedErr
		},
	}
	uc := buildUseCase(repo, nil, nil)

	dto := dtos.TransferStockDTO{
		ProductID:       "prod-001",
		FromWarehouseID: 1,
		ToWarehouseID:   2,
		BusinessID:      10,
		Quantity:        1000, // más de lo disponible
	}

	// Act
	err := uc.TransferStock(context.Background(), dto)

	// Assert
	if !errors.Is(err, expectedErr) {
		t.Errorf("error esperado %v, se obtuvo: %v", expectedErr, err)
	}
}
