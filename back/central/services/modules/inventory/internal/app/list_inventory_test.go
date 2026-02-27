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

func TestListWarehouseInventory_BodegaExiste_RetornaInventarioPaginado(t *testing.T) {
	// Arrange
	expectedLevels := []entities.InventoryLevel{
		{ID: 1, ProductID: "prod-001", WarehouseID: 1, Quantity: 100},
		{ID: 2, ProductID: "prod-002", WarehouseID: 1, Quantity: 50},
	}
	const expectedTotal int64 = 2

	repo := &mocks.RepositoryMock{
		WarehouseExistsFn: func(ctx context.Context, warehouseID uint, businessID uint) (bool, error) {
			return true, nil
		},
		ListWarehouseInventoryFn: func(ctx context.Context, params dtos.ListWarehouseInventoryParams) ([]entities.InventoryLevel, int64, error) {
			return expectedLevels, expectedTotal, nil
		},
	}
	uc := buildUseCase(repo, nil, nil)

	params := dtos.ListWarehouseInventoryParams{
		WarehouseID: 1,
		BusinessID:  10,
		Page:        1,
		PageSize:    20,
	}

	// Act
	result, total, err := uc.ListWarehouseInventory(context.Background(), params)

	// Assert
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if total != expectedTotal {
		t.Errorf("total esperado %d, se obtuvo %d", expectedTotal, total)
	}
	if len(result) != 2 {
		t.Errorf("cantidad de resultados esperada 2, se obtuvo %d", len(result))
	}
}

func TestListWarehouseInventory_BodegaNoExiste_RetornaErrWarehouseNotFound(t *testing.T) {
	// Arrange
	repo := &mocks.RepositoryMock{
		WarehouseExistsFn: func(ctx context.Context, warehouseID uint, businessID uint) (bool, error) {
			return false, nil
		},
	}
	uc := buildUseCase(repo, nil, nil)

	params := dtos.ListWarehouseInventoryParams{
		WarehouseID: 99,
		BusinessID:  10,
		Page:        1,
		PageSize:    20,
	}

	// Act
	result, total, err := uc.ListWarehouseInventory(context.Background(), params)

	// Assert
	if !errors.Is(err, domainerrors.ErrWarehouseNotFound) {
		t.Errorf("error esperado ErrWarehouseNotFound, se obtuvo: %v", err)
	}
	if result != nil {
		t.Errorf("resultado esperado nil, se obtuvo: %v", result)
	}
	if total != 0 {
		t.Errorf("total esperado 0, se obtuvo %d", total)
	}
}

func TestListWarehouseInventory_PaginaInvalida_UsaDefault1(t *testing.T) {
	// Arrange
	paramsCapturados := dtos.ListWarehouseInventoryParams{}
	repo := &mocks.RepositoryMock{
		WarehouseExistsFn: func(ctx context.Context, warehouseID uint, businessID uint) (bool, error) {
			return true, nil
		},
		ListWarehouseInventoryFn: func(ctx context.Context, params dtos.ListWarehouseInventoryParams) ([]entities.InventoryLevel, int64, error) {
			paramsCapturados = params
			return nil, 0, nil
		},
	}
	uc := buildUseCase(repo, nil, nil)

	// Enviar page = 0 (inválido)
	params := dtos.ListWarehouseInventoryParams{
		WarehouseID: 1,
		BusinessID:  10,
		Page:        0, // invalido
		PageSize:    20,
	}

	// Act
	_, _, err := uc.ListWarehouseInventory(context.Background(), params)

	// Assert
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if paramsCapturados.Page != 1 {
		t.Errorf("page corregida esperada 1, se obtuvo %d", paramsCapturados.Page)
	}
}

func TestListWarehouseInventory_PageSizeMuyGrande_UsaDefault20(t *testing.T) {
	// Arrange
	paramsCapturados := dtos.ListWarehouseInventoryParams{}
	repo := &mocks.RepositoryMock{
		WarehouseExistsFn: func(ctx context.Context, warehouseID uint, businessID uint) (bool, error) {
			return true, nil
		},
		ListWarehouseInventoryFn: func(ctx context.Context, params dtos.ListWarehouseInventoryParams) ([]entities.InventoryLevel, int64, error) {
			paramsCapturados = params
			return nil, 0, nil
		},
	}
	uc := buildUseCase(repo, nil, nil)

	params := dtos.ListWarehouseInventoryParams{
		WarehouseID: 1,
		BusinessID:  10,
		Page:        1,
		PageSize:    500, // mayor a 100 -> debe usar default 20
	}

	// Act
	_, _, err := uc.ListWarehouseInventory(context.Background(), params)

	// Assert
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if paramsCapturados.PageSize != 20 {
		t.Errorf("page_size corregida esperada 20, se obtuvo %d", paramsCapturados.PageSize)
	}
}

func TestListWarehouseInventory_ErrorVerificandoBodega_PropagaError(t *testing.T) {
	// Arrange
	expectedErr := errors.New("error de conexion")
	repo := &mocks.RepositoryMock{
		WarehouseExistsFn: func(ctx context.Context, warehouseID uint, businessID uint) (bool, error) {
			return false, expectedErr
		},
	}
	uc := buildUseCase(repo, nil, nil)

	params := dtos.ListWarehouseInventoryParams{
		WarehouseID: 1,
		BusinessID:  10,
		Page:        1,
		PageSize:    20,
	}

	// Act
	_, _, err := uc.ListWarehouseInventory(context.Background(), params)

	// Assert
	if !errors.Is(err, expectedErr) {
		t.Errorf("error esperado %v, se obtuvo: %v", expectedErr, err)
	}
}

func TestListMovements_PaginacionCorregida(t *testing.T) {
	// Arrange
	paramsCapturados := dtos.ListMovementsParams{}
	repo := &mocks.RepositoryMock{
		ListMovementsFn: func(ctx context.Context, params dtos.ListMovementsParams) ([]entities.StockMovement, int64, error) {
			paramsCapturados = params
			return nil, 0, nil
		},
	}
	uc := buildUseCase(repo, nil, nil)

	params := dtos.ListMovementsParams{
		BusinessID: 10,
		Page:       -1,   // inválido
		PageSize:   9999, // mayor a 100
	}

	// Act
	_, _, err := uc.ListMovements(context.Background(), params)

	// Assert
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if paramsCapturados.Page != 1 {
		t.Errorf("page esperada 1, se obtuvo %d", paramsCapturados.Page)
	}
	if paramsCapturados.PageSize != 20 {
		t.Errorf("page_size esperada 20, se obtuvo %d", paramsCapturados.PageSize)
	}
}

func TestListMovementTypes_PaginacionCorregida(t *testing.T) {
	// Arrange
	paramsCapturados := dtos.ListStockMovementTypesParams{}
	repo := &mocks.RepositoryMock{
		ListMovementTypesFn: func(ctx context.Context, params dtos.ListStockMovementTypesParams) ([]entities.StockMovementType, int64, error) {
			paramsCapturados = params
			return nil, 0, nil
		},
	}
	uc := buildUseCase(repo, nil, nil)

	params := dtos.ListStockMovementTypesParams{
		Page:     0,   // inválido
		PageSize: 200, // mayor a 100
	}

	// Act
	_, _, err := uc.ListMovementTypes(context.Background(), params)

	// Assert
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if paramsCapturados.Page != 1 {
		t.Errorf("page esperada 1, se obtuvo %d", paramsCapturados.Page)
	}
	if paramsCapturados.PageSize != 20 {
		t.Errorf("page_size esperada 20, se obtuvo %d", paramsCapturados.PageSize)
	}
}
