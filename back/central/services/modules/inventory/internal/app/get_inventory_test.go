package app

import (
	"context"
	"errors"
	"testing"

	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/mocks"
)

func TestGetProductInventory_ProductoExiste_RetornaInventario(t *testing.T) {
	// Arrange
	expectedLevels := []entities.InventoryLevel{
		{
			ID:           1,
			ProductID:    "prod-001",
			WarehouseID:  1,
			BusinessID:   10,
			Quantity:     50,
			AvailableQty: 45,
			ReservedQty:  5,
		},
		{
			ID:           2,
			ProductID:    "prod-001",
			WarehouseID:  2,
			BusinessID:   10,
			Quantity:     30,
			AvailableQty: 30,
			ReservedQty:  0,
		},
	}

	repo := &mocks.RepositoryMock{
		GetProductByIDFn: func(ctx context.Context, productID string, businessID uint) (string, string, bool, error) {
			return "Producto Test", "SKU-001", true, nil
		},
		GetProductInventoryFn: func(ctx context.Context, params dtos.GetProductInventoryParams) ([]entities.InventoryLevel, error) {
			return expectedLevels, nil
		},
	}
	uc := buildUseCase(repo, nil, nil)

	params := dtos.GetProductInventoryParams{
		ProductID:  "prod-001",
		BusinessID: 10,
	}

	// Act
	result, err := uc.GetProductInventory(context.Background(), params)

	// Assert
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("cantidad de niveles esperada 2, se obtuvo %d", len(result))
	}
	if result[0].Quantity != 50 {
		t.Errorf("cantidad esperada 50, se obtuvo %d", result[0].Quantity)
	}
}

func TestGetProductInventory_ProductoNoExiste_RetornaErrProductNotFound(t *testing.T) {
	// Arrange
	repo := &mocks.RepositoryMock{
		GetProductByIDFn: func(ctx context.Context, productID string, businessID uint) (string, string, bool, error) {
			return "", "", false, errors.New("not found")
		},
	}
	uc := buildUseCase(repo, nil, nil)

	params := dtos.GetProductInventoryParams{
		ProductID:  "prod-inexistente",
		BusinessID: 10,
	}

	// Act
	result, err := uc.GetProductInventory(context.Background(), params)

	// Assert
	if !errors.Is(err, domainerrors.ErrProductNotFound) {
		t.Errorf("error esperado ErrProductNotFound, se obtuvo: %v", err)
	}
	if result != nil {
		t.Errorf("resultado esperado nil, se obtuvo: %v", result)
	}
}

func TestGetProductInventory_ErrorEnRepositorio_PropagaError(t *testing.T) {
	// Arrange
	expectedErr := errors.New("error en base de datos")
	repo := &mocks.RepositoryMock{
		GetProductByIDFn: func(ctx context.Context, productID string, businessID uint) (string, string, bool, error) {
			return "Producto", "SKU-001", true, nil
		},
		GetProductInventoryFn: func(ctx context.Context, params dtos.GetProductInventoryParams) ([]entities.InventoryLevel, error) {
			return nil, expectedErr
		},
	}
	uc := buildUseCase(repo, nil, nil)

	params := dtos.GetProductInventoryParams{
		ProductID:  "prod-001",
		BusinessID: 10,
	}

	// Act
	result, err := uc.GetProductInventory(context.Background(), params)

	// Assert
	if !errors.Is(err, expectedErr) {
		t.Errorf("error esperado %v, se obtuvo: %v", expectedErr, err)
	}
	if result != nil {
		t.Errorf("resultado esperado nil, se obtuvo: %v", result)
	}
}

func TestGetProductInventory_SinInventario_RetornaSliceVacio(t *testing.T) {
	// Arrange
	repo := &mocks.RepositoryMock{
		GetProductByIDFn: func(ctx context.Context, productID string, businessID uint) (string, string, bool, error) {
			return "Producto Test", "SKU-001", true, nil
		},
		GetProductInventoryFn: func(ctx context.Context, params dtos.GetProductInventoryParams) ([]entities.InventoryLevel, error) {
			return []entities.InventoryLevel{}, nil
		},
	}
	uc := buildUseCase(repo, nil, nil)

	params := dtos.GetProductInventoryParams{
		ProductID:  "prod-001",
		BusinessID: 10,
	}

	// Act
	result, err := uc.GetProductInventory(context.Background(), params)

	// Assert
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("cantidad de niveles esperada 0, se obtuvo %d", len(result))
	}
}
