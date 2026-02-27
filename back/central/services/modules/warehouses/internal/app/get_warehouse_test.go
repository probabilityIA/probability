package app

import (
	"context"
	"errors"
	"testing"

	domainerrors "github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/mocks"
)

func TestGetWarehouse_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	expected := &entities.Warehouse{
		ID:         7,
		BusinessID: 10,
		Name:       "Bodega Central",
		Code:       "BC-001",
		IsActive:   true,
	}
	repo := &mocks.RepositoryMock{
		GetByIDFn: func(_ context.Context, businessID, warehouseID uint) (*entities.Warehouse, error) {
			return expected, nil
		},
	}
	uc := New(repo)

	// Act
	result, err := uc.GetWarehouse(ctx, 10, 7)

	// Assert
	if err != nil {
		t.Fatalf("se esperaba nil error, se obtuvo: %v", err)
	}
	if result == nil {
		t.Fatal("se esperaba un warehouse, se obtuvo nil")
	}
	if result.ID != expected.ID {
		t.Errorf("ID: esperado %d, obtenido %d", expected.ID, result.ID)
	}
	if result.Name != expected.Name {
		t.Errorf("Name: esperado %q, obtenido %q", expected.Name, result.Name)
	}
}

func TestGetWarehouse_NotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	repo := &mocks.RepositoryMock{
		GetByIDFn: func(_ context.Context, _, _ uint) (*entities.Warehouse, error) {
			return nil, domainerrors.ErrWarehouseNotFound
		},
	}
	uc := New(repo)

	// Act
	result, err := uc.GetWarehouse(ctx, 10, 999)

	// Assert
	if result != nil {
		t.Errorf("se esperaba nil result, se obtuvo: %+v", result)
	}
	if !errors.Is(err, domainerrors.ErrWarehouseNotFound) {
		t.Errorf("se esperaba ErrWarehouseNotFound, se obtuvo: %v", err)
	}
}

func TestGetWarehouse_RepositoryError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	dbErr := errors.New("database unavailable")
	repo := &mocks.RepositoryMock{
		GetByIDFn: func(_ context.Context, _, _ uint) (*entities.Warehouse, error) {
			return nil, dbErr
		},
	}
	uc := New(repo)

	// Act
	result, err := uc.GetWarehouse(ctx, 10, 1)

	// Assert
	if result != nil {
		t.Errorf("se esperaba nil result, se obtuvo: %+v", result)
	}
	if !errors.Is(err, dbErr) {
		t.Errorf("se esperaba error %v, se obtuvo: %v", dbErr, err)
	}
}

func TestGetWarehouse_PropagatesCorrectIDs(t *testing.T) {
	// Arrange
	ctx := context.Background()
	var receivedBusinessID, receivedWarehouseID uint

	repo := &mocks.RepositoryMock{
		GetByIDFn: func(_ context.Context, bID, wID uint) (*entities.Warehouse, error) {
			receivedBusinessID = bID
			receivedWarehouseID = wID
			return &entities.Warehouse{ID: wID, BusinessID: bID}, nil
		},
	}
	uc := New(repo)

	// Act
	_, err := uc.GetWarehouse(ctx, 42, 17)

	// Assert
	if err != nil {
		t.Fatalf("se esperaba nil error, se obtuvo: %v", err)
	}
	if receivedBusinessID != 42 {
		t.Errorf("businessID: esperado 42, obtenido %d", receivedBusinessID)
	}
	if receivedWarehouseID != 17 {
		t.Errorf("warehouseID: esperado 17, obtenido %d", receivedWarehouseID)
	}
}
