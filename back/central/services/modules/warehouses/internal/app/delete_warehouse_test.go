package app

import (
	"context"
	"errors"
	"testing"

	domainerrors "github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/mocks"
)

func TestDeleteWarehouse_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	repo := &mocks.RepositoryMock{
		DeleteFn: func(_ context.Context, _, _ uint) error {
			return nil
		},
	}
	uc := New(repo)

	// Act
	err := uc.DeleteWarehouse(ctx, 10, 5)

	// Assert
	if err != nil {
		t.Errorf("se esperaba nil error, se obtuvo: %v", err)
	}
}

func TestDeleteWarehouse_NotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	repo := &mocks.RepositoryMock{
		DeleteFn: func(_ context.Context, _, _ uint) error {
			return domainerrors.ErrWarehouseNotFound
		},
	}
	uc := New(repo)

	// Act
	err := uc.DeleteWarehouse(ctx, 10, 999)

	// Assert
	if !errors.Is(err, domainerrors.ErrWarehouseNotFound) {
		t.Errorf("se esperaba ErrWarehouseNotFound, se obtuvo: %v", err)
	}
}

func TestDeleteWarehouse_HasStock(t *testing.T) {
	// Arrange
	ctx := context.Background()
	repo := &mocks.RepositoryMock{
		DeleteFn: func(_ context.Context, _, _ uint) error {
			return domainerrors.ErrWarehouseHasStock
		},
	}
	uc := New(repo)

	// Act
	err := uc.DeleteWarehouse(ctx, 10, 3)

	// Assert
	if !errors.Is(err, domainerrors.ErrWarehouseHasStock) {
		t.Errorf("se esperaba ErrWarehouseHasStock, se obtuvo: %v", err)
	}
}

func TestDeleteWarehouse_RepositoryError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	dbErr := errors.New("delete query failed")
	repo := &mocks.RepositoryMock{
		DeleteFn: func(_ context.Context, _, _ uint) error {
			return dbErr
		},
	}
	uc := New(repo)

	// Act
	err := uc.DeleteWarehouse(ctx, 10, 1)

	// Assert
	if !errors.Is(err, dbErr) {
		t.Errorf("se esperaba error %v, se obtuvo: %v", dbErr, err)
	}
}

func TestDeleteWarehouse_PropagatesCorrectIDs(t *testing.T) {
	// Arrange
	ctx := context.Background()
	var receivedBusinessID, receivedWarehouseID uint

	repo := &mocks.RepositoryMock{
		DeleteFn: func(_ context.Context, bID, wID uint) error {
			receivedBusinessID = bID
			receivedWarehouseID = wID
			return nil
		},
	}
	uc := New(repo)

	// Act
	err := uc.DeleteWarehouse(ctx, 42, 17)

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
