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

func TestUpdateLocation_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	cap := 200
	updated := &entities.WarehouseLocation{ID: 3, WarehouseID: 5, Code: "LOC-NEW", Name: "Zona Actualizada"}

	repo := &mocks.RepositoryMock{
		GetByIDFn: func(_ context.Context, _, _ uint) (*entities.Warehouse, error) {
			return &entities.Warehouse{ID: 5, BusinessID: 10}, nil
		},
		GetLocationByIDFn: func(_ context.Context, _, _ uint) (*entities.WarehouseLocation, error) {
			return &entities.WarehouseLocation{ID: 3, WarehouseID: 5}, nil
		},
		LocationExistsByCodeFn: func(_ context.Context, _ uint, _ string, _ *uint) (bool, error) {
			return false, nil
		},
		UpdateLocationFn: func(_ context.Context, _ *entities.WarehouseLocation) (*entities.WarehouseLocation, error) {
			return updated, nil
		},
	}
	uc := New(repo)

	dto := dtos.UpdateLocationDTO{
		ID:          3,
		WarehouseID: 5,
		BusinessID:  10,
		Name:        "Zona Actualizada",
		Code:        "LOC-NEW",
		Type:        "packing",
		IsActive:    true,
		Capacity:    &cap,
	}

	// Act
	result, err := uc.UpdateLocation(ctx, dto)

	// Assert
	if err != nil {
		t.Fatalf("se esperaba nil error, se obtuvo: %v", err)
	}
	if result == nil {
		t.Fatal("se esperaba una ubicación, se obtuvo nil")
	}
	if result.Code != updated.Code {
		t.Errorf("Code: esperado %q, obtenido %q", updated.Code, result.Code)
	}
}

func TestUpdateLocation_WarehouseNotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	repo := &mocks.RepositoryMock{
		GetByIDFn: func(_ context.Context, _, _ uint) (*entities.Warehouse, error) {
			return nil, domainerrors.ErrWarehouseNotFound
		},
	}
	uc := New(repo)

	dto := dtos.UpdateLocationDTO{
		ID:          3,
		WarehouseID: 999,
		BusinessID:  10,
		Code:        "LOC-001",
	}

	// Act
	result, err := uc.UpdateLocation(ctx, dto)

	// Assert
	if result != nil {
		t.Errorf("se esperaba nil result, se obtuvo: %+v", result)
	}
	if !errors.Is(err, domainerrors.ErrWarehouseNotFound) {
		t.Errorf("se esperaba ErrWarehouseNotFound, se obtuvo: %v", err)
	}
}

func TestUpdateLocation_GetByIDRepositoryError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	dbErr := errors.New("get warehouse failed")

	repo := &mocks.RepositoryMock{
		GetByIDFn: func(_ context.Context, _, _ uint) (*entities.Warehouse, error) {
			return nil, dbErr
		},
	}
	uc := New(repo)

	dto := dtos.UpdateLocationDTO{
		ID:          3,
		WarehouseID: 5,
		BusinessID:  10,
		Code:        "LOC-001",
	}

	// Act
	result, err := uc.UpdateLocation(ctx, dto)

	// Assert
	if result != nil {
		t.Errorf("se esperaba nil result, se obtuvo: %+v", result)
	}
	if !errors.Is(err, dbErr) {
		t.Errorf("se esperaba error %v, se obtuvo: %v", dbErr, err)
	}
}

func TestUpdateLocation_LocationNotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	repo := &mocks.RepositoryMock{
		GetByIDFn: func(_ context.Context, _, _ uint) (*entities.Warehouse, error) {
			return &entities.Warehouse{ID: 5, BusinessID: 10}, nil
		},
		GetLocationByIDFn: func(_ context.Context, _, _ uint) (*entities.WarehouseLocation, error) {
			return nil, domainerrors.ErrLocationNotFound
		},
	}
	uc := New(repo)

	dto := dtos.UpdateLocationDTO{
		ID:          999,
		WarehouseID: 5,
		BusinessID:  10,
		Code:        "LOC-GHOST",
	}

	// Act
	result, err := uc.UpdateLocation(ctx, dto)

	// Assert
	if result != nil {
		t.Errorf("se esperaba nil result, se obtuvo: %+v", result)
	}
	if !errors.Is(err, domainerrors.ErrLocationNotFound) {
		t.Errorf("se esperaba ErrLocationNotFound, se obtuvo: %v", err)
	}
}

func TestUpdateLocation_GetLocationByIDRepositoryError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	dbErr := errors.New("get location failed")

	repo := &mocks.RepositoryMock{
		GetByIDFn: func(_ context.Context, _, _ uint) (*entities.Warehouse, error) {
			return &entities.Warehouse{ID: 5, BusinessID: 10}, nil
		},
		GetLocationByIDFn: func(_ context.Context, _, _ uint) (*entities.WarehouseLocation, error) {
			return nil, dbErr
		},
	}
	uc := New(repo)

	dto := dtos.UpdateLocationDTO{
		ID:          3,
		WarehouseID: 5,
		BusinessID:  10,
		Code:        "LOC-001",
	}

	// Act
	result, err := uc.UpdateLocation(ctx, dto)

	// Assert
	if result != nil {
		t.Errorf("se esperaba nil result, se obtuvo: %+v", result)
	}
	if !errors.Is(err, dbErr) {
		t.Errorf("se esperaba error %v, se obtuvo: %v", dbErr, err)
	}
}

func TestUpdateLocation_DuplicateCode(t *testing.T) {
	// Arrange
	ctx := context.Background()
	repo := &mocks.RepositoryMock{
		GetByIDFn: func(_ context.Context, _, _ uint) (*entities.Warehouse, error) {
			return &entities.Warehouse{ID: 5, BusinessID: 10}, nil
		},
		GetLocationByIDFn: func(_ context.Context, _, _ uint) (*entities.WarehouseLocation, error) {
			return &entities.WarehouseLocation{ID: 3, WarehouseID: 5}, nil
		},
		LocationExistsByCodeFn: func(_ context.Context, _ uint, _ string, _ *uint) (bool, error) {
			return true, nil // código ya existe en otra ubicación
		},
	}
	uc := New(repo)

	dto := dtos.UpdateLocationDTO{
		ID:          3,
		WarehouseID: 5,
		BusinessID:  10,
		Code:        "LOC-DUPLICADO",
	}

	// Act
	result, err := uc.UpdateLocation(ctx, dto)

	// Assert
	if result != nil {
		t.Errorf("se esperaba nil result, se obtuvo: %+v", result)
	}
	if !errors.Is(err, domainerrors.ErrDuplicateLocCode) {
		t.Errorf("se esperaba ErrDuplicateLocCode, se obtuvo: %v", err)
	}
}

func TestUpdateLocation_LocationExistsByCodeError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	dbErr := errors.New("exists check failed")

	repo := &mocks.RepositoryMock{
		GetByIDFn: func(_ context.Context, _, _ uint) (*entities.Warehouse, error) {
			return &entities.Warehouse{ID: 5, BusinessID: 10}, nil
		},
		GetLocationByIDFn: func(_ context.Context, _, _ uint) (*entities.WarehouseLocation, error) {
			return &entities.WarehouseLocation{ID: 3, WarehouseID: 5}, nil
		},
		LocationExistsByCodeFn: func(_ context.Context, _ uint, _ string, _ *uint) (bool, error) {
			return false, dbErr
		},
	}
	uc := New(repo)

	dto := dtos.UpdateLocationDTO{
		ID:          3,
		WarehouseID: 5,
		BusinessID:  10,
		Code:        "LOC-001",
	}

	// Act
	result, err := uc.UpdateLocation(ctx, dto)

	// Assert
	if result != nil {
		t.Errorf("se esperaba nil result, se obtuvo: %+v", result)
	}
	if !errors.Is(err, dbErr) {
		t.Errorf("se esperaba error %v, se obtuvo: %v", dbErr, err)
	}
}

func TestUpdateLocation_UpdateLocationRepositoryError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	dbErr := errors.New("update location failed")

	repo := &mocks.RepositoryMock{
		GetByIDFn: func(_ context.Context, _, _ uint) (*entities.Warehouse, error) {
			return &entities.Warehouse{ID: 5, BusinessID: 10}, nil
		},
		GetLocationByIDFn: func(_ context.Context, _, _ uint) (*entities.WarehouseLocation, error) {
			return &entities.WarehouseLocation{ID: 3, WarehouseID: 5}, nil
		},
		LocationExistsByCodeFn: func(_ context.Context, _ uint, _ string, _ *uint) (bool, error) {
			return false, nil
		},
		UpdateLocationFn: func(_ context.Context, _ *entities.WarehouseLocation) (*entities.WarehouseLocation, error) {
			return nil, dbErr
		},
	}
	uc := New(repo)

	dto := dtos.UpdateLocationDTO{
		ID:          3,
		WarehouseID: 5,
		BusinessID:  10,
		Code:        "LOC-001",
	}

	// Act
	result, err := uc.UpdateLocation(ctx, dto)

	// Assert
	if result != nil {
		t.Errorf("se esperaba nil result, se obtuvo: %+v", result)
	}
	if !errors.Is(err, dbErr) {
		t.Errorf("se esperaba error %v, se obtuvo: %v", dbErr, err)
	}
}

func TestUpdateLocation_LocationExistsByCode_ReceivesExcludeID(t *testing.T) {
	// Arrange — valida que LocationExistsByCode se llame con el ID correcto como excludeID
	ctx := context.Background()
	var receivedExcludeID *uint

	repo := &mocks.RepositoryMock{
		GetByIDFn: func(_ context.Context, _, _ uint) (*entities.Warehouse, error) {
			return &entities.Warehouse{ID: 5, BusinessID: 10}, nil
		},
		GetLocationByIDFn: func(_ context.Context, _, _ uint) (*entities.WarehouseLocation, error) {
			return &entities.WarehouseLocation{ID: 8, WarehouseID: 5}, nil
		},
		LocationExistsByCodeFn: func(_ context.Context, _ uint, _ string, excludeID *uint) (bool, error) {
			receivedExcludeID = excludeID
			return false, nil
		},
		UpdateLocationFn: func(_ context.Context, loc *entities.WarehouseLocation) (*entities.WarehouseLocation, error) {
			return loc, nil
		},
	}
	uc := New(repo)

	dto := dtos.UpdateLocationDTO{
		ID:          8,
		WarehouseID: 5,
		BusinessID:  10,
		Code:        "LOC-008",
	}

	// Act
	_, err := uc.UpdateLocation(ctx, dto)

	// Assert
	if err != nil {
		t.Fatalf("se esperaba nil error, se obtuvo: %v", err)
	}
	if receivedExcludeID == nil {
		t.Fatal("LocationExistsByCode debería recibir excludeID no nil")
	}
	if *receivedExcludeID != dto.ID {
		t.Errorf("excludeID: esperado %d, obtenido %d", dto.ID, *receivedExcludeID)
	}
}
