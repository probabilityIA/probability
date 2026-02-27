package app

import (
	"context"
	"errors"
	"testing"

	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/mocks"
)

func TestListWarehouses_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	expected := []entities.Warehouse{
		{ID: 1, BusinessID: 10, Name: "Bodega A"},
		{ID: 2, BusinessID: 10, Name: "Bodega B"},
	}
	repo := &mocks.RepositoryMock{
		ListFn: func(_ context.Context, params dtos.ListWarehousesParams) ([]entities.Warehouse, int64, error) {
			return expected, int64(len(expected)), nil
		},
	}
	uc := New(repo)

	params := dtos.ListWarehousesParams{
		BusinessID: 10,
		Page:       1,
		PageSize:   20,
	}

	// Act
	result, total, err := uc.ListWarehouses(ctx, params)

	// Assert
	if err != nil {
		t.Fatalf("se esperaba nil error, se obtuvo: %v", err)
	}
	if total != 2 {
		t.Errorf("total: esperado 2, obtenido %d", total)
	}
	if len(result) != 2 {
		t.Errorf("len(result): esperado 2, obtenido %d", len(result))
	}
}

func TestListWarehouses_EmptyResult(t *testing.T) {
	// Arrange
	ctx := context.Background()
	repo := &mocks.RepositoryMock{
		ListFn: func(_ context.Context, _ dtos.ListWarehousesParams) ([]entities.Warehouse, int64, error) {
			return []entities.Warehouse{}, 0, nil
		},
	}
	uc := New(repo)

	params := dtos.ListWarehousesParams{
		BusinessID: 99,
		Page:       1,
		PageSize:   20,
	}

	// Act
	result, total, err := uc.ListWarehouses(ctx, params)

	// Assert
	if err != nil {
		t.Fatalf("se esperaba nil error, se obtuvo: %v", err)
	}
	if total != 0 {
		t.Errorf("total: esperado 0, obtenido %d", total)
	}
	if len(result) != 0 {
		t.Errorf("len(result): esperado 0, obtenido %d", len(result))
	}
}

func TestListWarehouses_RepositoryError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	dbErr := errors.New("query failed")
	repo := &mocks.RepositoryMock{
		ListFn: func(_ context.Context, _ dtos.ListWarehousesParams) ([]entities.Warehouse, int64, error) {
			return nil, 0, dbErr
		},
	}
	uc := New(repo)

	params := dtos.ListWarehousesParams{
		BusinessID: 10,
		Page:       1,
		PageSize:   20,
	}

	// Act
	result, total, err := uc.ListWarehouses(ctx, params)

	// Assert
	if result != nil {
		t.Errorf("se esperaba nil result, se obtuvo: %v", result)
	}
	if total != 0 {
		t.Errorf("se esperaba total=0, se obtuvo: %d", total)
	}
	if !errors.Is(err, dbErr) {
		t.Errorf("se esperaba error %v, se obtuvo: %v", dbErr, err)
	}
}

func TestListWarehouses_PageLessThanOne_NormalizesToOne(t *testing.T) {
	// Arrange
	ctx := context.Background()
	var receivedParams dtos.ListWarehousesParams

	repo := &mocks.RepositoryMock{
		ListFn: func(_ context.Context, params dtos.ListWarehousesParams) ([]entities.Warehouse, int64, error) {
			receivedParams = params
			return []entities.Warehouse{}, 0, nil
		},
	}
	uc := New(repo)

	// Se envía page=0 (inválida)
	params := dtos.ListWarehousesParams{
		BusinessID: 10,
		Page:       0,
		PageSize:   20,
	}

	// Act
	_, _, err := uc.ListWarehouses(ctx, params)

	// Assert
	if err != nil {
		t.Fatalf("se esperaba nil error, se obtuvo: %v", err)
	}
	if receivedParams.Page != 1 {
		t.Errorf("Page debería normalizarse a 1, se obtuvo: %d", receivedParams.Page)
	}
}

func TestListWarehouses_PageSizeZero_NormalizesTo20(t *testing.T) {
	// Arrange
	ctx := context.Background()
	var receivedParams dtos.ListWarehousesParams

	repo := &mocks.RepositoryMock{
		ListFn: func(_ context.Context, params dtos.ListWarehousesParams) ([]entities.Warehouse, int64, error) {
			receivedParams = params
			return []entities.Warehouse{}, 0, nil
		},
	}
	uc := New(repo)

	params := dtos.ListWarehousesParams{
		BusinessID: 10,
		Page:       1,
		PageSize:   0, // inválido
	}

	// Act
	_, _, err := uc.ListWarehouses(ctx, params)

	// Assert
	if err != nil {
		t.Fatalf("se esperaba nil error, se obtuvo: %v", err)
	}
	if receivedParams.PageSize != 20 {
		t.Errorf("PageSize debería normalizarse a 20, se obtuvo: %d", receivedParams.PageSize)
	}
}

func TestListWarehouses_PageSizeAboveMax_NormalizesTo20(t *testing.T) {
	// Arrange
	ctx := context.Background()
	var receivedParams dtos.ListWarehousesParams

	repo := &mocks.RepositoryMock{
		ListFn: func(_ context.Context, params dtos.ListWarehousesParams) ([]entities.Warehouse, int64, error) {
			receivedParams = params
			return []entities.Warehouse{}, 0, nil
		},
	}
	uc := New(repo)

	params := dtos.ListWarehousesParams{
		BusinessID: 10,
		Page:       1,
		PageSize:   200, // supera el máximo de 100
	}

	// Act
	_, _, err := uc.ListWarehouses(ctx, params)

	// Assert
	if err != nil {
		t.Fatalf("se esperaba nil error, se obtuvo: %v", err)
	}
	if receivedParams.PageSize != 20 {
		t.Errorf("PageSize debería normalizarse a 20 cuando supera 100, se obtuvo: %d", receivedParams.PageSize)
	}
}

func TestListWarehouses_NegativePage_NormalizesToOne(t *testing.T) {
	// Arrange
	ctx := context.Background()
	var receivedParams dtos.ListWarehousesParams

	repo := &mocks.RepositoryMock{
		ListFn: func(_ context.Context, params dtos.ListWarehousesParams) ([]entities.Warehouse, int64, error) {
			receivedParams = params
			return []entities.Warehouse{}, 0, nil
		},
	}
	uc := New(repo)

	params := dtos.ListWarehousesParams{
		BusinessID: 10,
		Page:       -5, // negativo
		PageSize:   10,
	}

	// Act
	_, _, err := uc.ListWarehouses(ctx, params)

	// Assert
	if err != nil {
		t.Fatalf("se esperaba nil error, se obtuvo: %v", err)
	}
	if receivedParams.Page != 1 {
		t.Errorf("Page debería normalizarse a 1 cuando es negativo, se obtuvo: %d", receivedParams.Page)
	}
}
