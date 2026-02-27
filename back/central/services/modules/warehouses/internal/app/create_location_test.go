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

func TestCreateLocation_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	capacity := 100

	repo := &mocks.RepositoryMock{
		GetByIDFn: func(_ context.Context, _, _ uint) (*entities.Warehouse, error) {
			return &entities.Warehouse{ID: 5, BusinessID: 10}, nil
		},
		LocationExistsByCodeFn: func(_ context.Context, _ uint, _ string, _ *uint) (bool, error) {
			return false, nil
		},
		CreateLocationFn: func(_ context.Context, loc *entities.WarehouseLocation) (*entities.WarehouseLocation, error) {
			loc.ID = 1
			return loc, nil
		},
	}
	uc := New(repo)

	dto := dtos.CreateLocationDTO{
		WarehouseID:   5,
		BusinessID:    10,
		Name:          "Estante A1",
		Code:          "EST-A1",
		Type:          "storage",
		IsActive:      true,
		IsFulfillment: false,
		Capacity:      &capacity,
	}

	// Act
	result, err := uc.CreateLocation(ctx, dto)

	// Assert
	if err != nil {
		t.Fatalf("se esperaba nil error, se obtuvo: %v", err)
	}
	if result == nil {
		t.Fatal("se esperaba una ubicación, se obtuvo nil")
	}
	if result.ID != 1 {
		t.Errorf("ID: esperado 1, obtenido %d", result.ID)
	}
	if result.Code != dto.Code {
		t.Errorf("Code: esperado %q, obtenido %q", dto.Code, result.Code)
	}
}

func TestCreateLocation_WarehouseNotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	repo := &mocks.RepositoryMock{
		GetByIDFn: func(_ context.Context, _, _ uint) (*entities.Warehouse, error) {
			return nil, domainerrors.ErrWarehouseNotFound
		},
	}
	uc := New(repo)

	dto := dtos.CreateLocationDTO{
		WarehouseID: 999,
		BusinessID:  10,
		Code:        "LOC-001",
	}

	// Act
	result, err := uc.CreateLocation(ctx, dto)

	// Assert
	if result != nil {
		t.Errorf("se esperaba nil result, se obtuvo: %+v", result)
	}
	if !errors.Is(err, domainerrors.ErrWarehouseNotFound) {
		t.Errorf("se esperaba ErrWarehouseNotFound, se obtuvo: %v", err)
	}
}

func TestCreateLocation_GetByIDRepositoryError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	dbErr := errors.New("db error")

	repo := &mocks.RepositoryMock{
		GetByIDFn: func(_ context.Context, _, _ uint) (*entities.Warehouse, error) {
			return nil, dbErr
		},
	}
	uc := New(repo)

	dto := dtos.CreateLocationDTO{
		WarehouseID: 5,
		BusinessID:  10,
		Code:        "LOC-001",
	}

	// Act
	result, err := uc.CreateLocation(ctx, dto)

	// Assert
	if result != nil {
		t.Errorf("se esperaba nil result, se obtuvo: %+v", result)
	}
	if !errors.Is(err, dbErr) {
		t.Errorf("se esperaba error %v, se obtuvo: %v", dbErr, err)
	}
}

func TestCreateLocation_DuplicateCode(t *testing.T) {
	// Arrange
	ctx := context.Background()
	repo := &mocks.RepositoryMock{
		GetByIDFn: func(_ context.Context, _, _ uint) (*entities.Warehouse, error) {
			return &entities.Warehouse{ID: 5, BusinessID: 10}, nil
		},
		LocationExistsByCodeFn: func(_ context.Context, _ uint, _ string, _ *uint) (bool, error) {
			return true, nil // código de ubicación duplicado
		},
	}
	uc := New(repo)

	dto := dtos.CreateLocationDTO{
		WarehouseID: 5,
		BusinessID:  10,
		Code:        "EST-DUPLICADO",
	}

	// Act
	result, err := uc.CreateLocation(ctx, dto)

	// Assert
	if result != nil {
		t.Errorf("se esperaba nil result, se obtuvo: %+v", result)
	}
	if !errors.Is(err, domainerrors.ErrDuplicateLocCode) {
		t.Errorf("se esperaba ErrDuplicateLocCode, se obtuvo: %v", err)
	}
}

func TestCreateLocation_LocationExistsByCodeError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	dbErr := errors.New("exists check failed")

	repo := &mocks.RepositoryMock{
		GetByIDFn: func(_ context.Context, _, _ uint) (*entities.Warehouse, error) {
			return &entities.Warehouse{ID: 5, BusinessID: 10}, nil
		},
		LocationExistsByCodeFn: func(_ context.Context, _ uint, _ string, _ *uint) (bool, error) {
			return false, dbErr
		},
	}
	uc := New(repo)

	dto := dtos.CreateLocationDTO{
		WarehouseID: 5,
		BusinessID:  10,
		Code:        "LOC-001",
	}

	// Act
	result, err := uc.CreateLocation(ctx, dto)

	// Assert
	if result != nil {
		t.Errorf("se esperaba nil result, se obtuvo: %+v", result)
	}
	if !errors.Is(err, dbErr) {
		t.Errorf("se esperaba error %v, se obtuvo: %v", dbErr, err)
	}
}

func TestCreateLocation_CreateLocationRepositoryError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	dbErr := errors.New("insert location failed")

	repo := &mocks.RepositoryMock{
		GetByIDFn: func(_ context.Context, _, _ uint) (*entities.Warehouse, error) {
			return &entities.Warehouse{ID: 5, BusinessID: 10}, nil
		},
		LocationExistsByCodeFn: func(_ context.Context, _ uint, _ string, _ *uint) (bool, error) {
			return false, nil
		},
		CreateLocationFn: func(_ context.Context, _ *entities.WarehouseLocation) (*entities.WarehouseLocation, error) {
			return nil, dbErr
		},
	}
	uc := New(repo)

	dto := dtos.CreateLocationDTO{
		WarehouseID: 5,
		BusinessID:  10,
		Code:        "LOC-001",
	}

	// Act
	result, err := uc.CreateLocation(ctx, dto)

	// Assert
	if result != nil {
		t.Errorf("se esperaba nil result, se obtuvo: %+v", result)
	}
	if !errors.Is(err, dbErr) {
		t.Errorf("se esperaba error %v, se obtuvo: %v", dbErr, err)
	}
}

func TestCreateLocation_MapsAllDTOFields(t *testing.T) {
	// Arrange
	ctx := context.Background()
	capacity := 250

	var captured *entities.WarehouseLocation

	repo := &mocks.RepositoryMock{
		GetByIDFn: func(_ context.Context, _, _ uint) (*entities.Warehouse, error) {
			return &entities.Warehouse{ID: 5, BusinessID: 10}, nil
		},
		LocationExistsByCodeFn: func(_ context.Context, _ uint, _ string, _ *uint) (bool, error) {
			return false, nil
		},
		CreateLocationFn: func(_ context.Context, loc *entities.WarehouseLocation) (*entities.WarehouseLocation, error) {
			captured = loc
			loc.ID = 10
			return loc, nil
		},
	}
	uc := New(repo)

	dto := dtos.CreateLocationDTO{
		WarehouseID:   5,
		BusinessID:    10,
		Name:          "Zona de Picking",
		Code:          "ZP-001",
		Type:          "picking",
		IsActive:      true,
		IsFulfillment: true,
		Capacity:      &capacity,
	}

	// Act
	_, err := uc.CreateLocation(ctx, dto)

	// Assert
	if err != nil {
		t.Fatalf("se esperaba nil error, se obtuvo: %v", err)
	}
	if captured == nil {
		t.Fatal("CreateLocation no fue llamado")
	}
	if captured.WarehouseID != dto.WarehouseID {
		t.Errorf("WarehouseID: esperado %d, obtenido %d", dto.WarehouseID, captured.WarehouseID)
	}
	if captured.Name != dto.Name {
		t.Errorf("Name: esperado %q, obtenido %q", dto.Name, captured.Name)
	}
	if captured.Code != dto.Code {
		t.Errorf("Code: esperado %q, obtenido %q", dto.Code, captured.Code)
	}
	if captured.Type != dto.Type {
		t.Errorf("Type: esperado %q, obtenido %q", dto.Type, captured.Type)
	}
	if captured.IsActive != dto.IsActive {
		t.Errorf("IsActive: esperado %v, obtenido %v", dto.IsActive, captured.IsActive)
	}
	if captured.IsFulfillment != dto.IsFulfillment {
		t.Errorf("IsFulfillment: esperado %v, obtenido %v", dto.IsFulfillment, captured.IsFulfillment)
	}
	if captured.Capacity == nil || *captured.Capacity != capacity {
		t.Errorf("Capacity: esperado %d, obtenido %v", capacity, captured.Capacity)
	}
}

func TestCreateLocation_NilCapacity(t *testing.T) {
	// Arrange — la capacidad puede ser nil (sin límite)
	ctx := context.Background()

	var captured *entities.WarehouseLocation

	repo := &mocks.RepositoryMock{
		GetByIDFn: func(_ context.Context, _, _ uint) (*entities.Warehouse, error) {
			return &entities.Warehouse{ID: 5, BusinessID: 10}, nil
		},
		LocationExistsByCodeFn: func(_ context.Context, _ uint, _ string, _ *uint) (bool, error) {
			return false, nil
		},
		CreateLocationFn: func(_ context.Context, loc *entities.WarehouseLocation) (*entities.WarehouseLocation, error) {
			captured = loc
			return loc, nil
		},
	}
	uc := New(repo)

	dto := dtos.CreateLocationDTO{
		WarehouseID: 5,
		BusinessID:  10,
		Code:        "LOC-UNLIMITED",
		Capacity:    nil, // sin límite de capacidad
	}

	// Act
	result, err := uc.CreateLocation(ctx, dto)

	// Assert
	if err != nil {
		t.Fatalf("se esperaba nil error, se obtuvo: %v", err)
	}
	if result == nil {
		t.Fatal("se esperaba una ubicación, se obtuvo nil")
	}
	if captured != nil && captured.Capacity != nil {
		t.Errorf("Capacity debería ser nil, se obtuvo: %v", *captured.Capacity)
	}
}
