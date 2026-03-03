package app

import (
	"context"
	"errors"
	"testing"

	domainerrors "github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/mocks"
)

func TestListLocations_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	cap := 100
	expected := []entities.WarehouseLocation{
		{ID: 1, WarehouseID: 5, Name: "Zona A", Code: "ZA", Type: "storage", Capacity: &cap},
		{ID: 2, WarehouseID: 5, Name: "Zona B", Code: "ZB", Type: "picking"},
	}

	repo := &mocks.RepositoryMock{
		GetByIDFn: func(_ context.Context, _, _ uint) (*entities.Warehouse, error) {
			return &entities.Warehouse{ID: 5, BusinessID: 10}, nil
		},
		ListLocationsFn: func(_ context.Context, _ dtos.ListLocationsParams) ([]entities.WarehouseLocation, error) {
			return expected, nil
		},
	}
	uc := New(repo)

	params := dtos.ListLocationsParams{
		WarehouseID: 5,
		BusinessID:  10,
	}

	// Act
	result, err := uc.ListLocations(ctx, params)

	// Assert
	if err != nil {
		t.Fatalf("se esperaba nil error, se obtuvo: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("len(result): esperado 2, obtenido %d", len(result))
	}
	if result[0].Code != "ZA" {
		t.Errorf("result[0].Code: esperado %q, obtenido %q", "ZA", result[0].Code)
	}
}

func TestListLocations_EmptyResult(t *testing.T) {
	// Arrange
	ctx := context.Background()
	repo := &mocks.RepositoryMock{
		GetByIDFn: func(_ context.Context, _, _ uint) (*entities.Warehouse, error) {
			return &entities.Warehouse{ID: 5, BusinessID: 10}, nil
		},
		ListLocationsFn: func(_ context.Context, _ dtos.ListLocationsParams) ([]entities.WarehouseLocation, error) {
			return []entities.WarehouseLocation{}, nil
		},
	}
	uc := New(repo)

	params := dtos.ListLocationsParams{
		WarehouseID: 5,
		BusinessID:  10,
	}

	// Act
	result, err := uc.ListLocations(ctx, params)

	// Assert
	if err != nil {
		t.Fatalf("se esperaba nil error, se obtuvo: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("len(result): esperado 0, obtenido %d", len(result))
	}
}

func TestListLocations_WarehouseNotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	repo := &mocks.RepositoryMock{
		GetByIDFn: func(_ context.Context, _, _ uint) (*entities.Warehouse, error) {
			return nil, domainerrors.ErrWarehouseNotFound
		},
	}
	uc := New(repo)

	params := dtos.ListLocationsParams{
		WarehouseID: 999,
		BusinessID:  10,
	}

	// Act
	result, err := uc.ListLocations(ctx, params)

	// Assert
	if result != nil {
		t.Errorf("se esperaba nil result, se obtuvo: %v", result)
	}
	if !errors.Is(err, domainerrors.ErrWarehouseNotFound) {
		t.Errorf("se esperaba ErrWarehouseNotFound, se obtuvo: %v", err)
	}
}

func TestListLocations_GetByIDRepositoryError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	dbErr := errors.New("db error on get warehouse")

	repo := &mocks.RepositoryMock{
		GetByIDFn: func(_ context.Context, _, _ uint) (*entities.Warehouse, error) {
			return nil, dbErr
		},
	}
	uc := New(repo)

	params := dtos.ListLocationsParams{
		WarehouseID: 5,
		BusinessID:  10,
	}

	// Act
	result, err := uc.ListLocations(ctx, params)

	// Assert
	if result != nil {
		t.Errorf("se esperaba nil result, se obtuvo: %v", result)
	}
	if !errors.Is(err, dbErr) {
		t.Errorf("se esperaba error %v, se obtuvo: %v", dbErr, err)
	}
}

func TestListLocations_ListLocationsRepositoryError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	dbErr := errors.New("list locations failed")

	repo := &mocks.RepositoryMock{
		GetByIDFn: func(_ context.Context, _, _ uint) (*entities.Warehouse, error) {
			return &entities.Warehouse{ID: 5, BusinessID: 10}, nil
		},
		ListLocationsFn: func(_ context.Context, _ dtos.ListLocationsParams) ([]entities.WarehouseLocation, error) {
			return nil, dbErr
		},
	}
	uc := New(repo)

	params := dtos.ListLocationsParams{
		WarehouseID: 5,
		BusinessID:  10,
	}

	// Act
	result, err := uc.ListLocations(ctx, params)

	// Assert
	if result != nil {
		t.Errorf("se esperaba nil result, se obtuvo: %v", result)
	}
	if !errors.Is(err, dbErr) {
		t.Errorf("se esperaba error %v, se obtuvo: %v", dbErr, err)
	}
}

func TestListLocations_PropagatesParamsToRepository(t *testing.T) {
	// Arrange — verifica que los params lleguen correctamente al repositorio
	ctx := context.Background()
	var receivedParams dtos.ListLocationsParams
	var receivedGetBusinessID, receivedGetWarehouseID uint

	repo := &mocks.RepositoryMock{
		GetByIDFn: func(_ context.Context, bID, wID uint) (*entities.Warehouse, error) {
			receivedGetBusinessID = bID
			receivedGetWarehouseID = wID
			return &entities.Warehouse{ID: wID, BusinessID: bID}, nil
		},
		ListLocationsFn: func(_ context.Context, params dtos.ListLocationsParams) ([]entities.WarehouseLocation, error) {
			receivedParams = params
			return []entities.WarehouseLocation{}, nil
		},
	}
	uc := New(repo)

	params := dtos.ListLocationsParams{
		WarehouseID: 7,
		BusinessID:  42,
	}

	// Act
	_, err := uc.ListLocations(ctx, params)

	// Assert
	if err != nil {
		t.Fatalf("se esperaba nil error, se obtuvo: %v", err)
	}
	if receivedGetBusinessID != 42 {
		t.Errorf("GetByID businessID: esperado 42, obtenido %d", receivedGetBusinessID)
	}
	if receivedGetWarehouseID != 7 {
		t.Errorf("GetByID warehouseID: esperado 7, obtenido %d", receivedGetWarehouseID)
	}
	if receivedParams.WarehouseID != 7 {
		t.Errorf("ListLocations params.WarehouseID: esperado 7, obtenido %d", receivedParams.WarehouseID)
	}
	if receivedParams.BusinessID != 42 {
		t.Errorf("ListLocations params.BusinessID: esperado 42, obtenido %d", receivedParams.BusinessID)
	}
}
