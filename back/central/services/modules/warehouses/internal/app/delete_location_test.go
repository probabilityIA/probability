package app

import (
	"context"
	"errors"
	"testing"

	domainerrors "github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/mocks"
)

func TestDeleteLocation_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	repo := &mocks.RepositoryMock{
		GetByIDFn: func(_ context.Context, _, _ uint) (*entities.Warehouse, error) {
			return &entities.Warehouse{ID: 5, BusinessID: 10}, nil
		},
		DeleteLocationFn: func(_ context.Context, _, _ uint) error {
			return nil
		},
	}
	uc := New(repo)

	// Act
	err := uc.DeleteLocation(ctx, 5, 3, 10)

	// Assert
	if err != nil {
		t.Errorf("se esperaba nil error, se obtuvo: %v", err)
	}
}

func TestDeleteLocation_WarehouseNotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	repo := &mocks.RepositoryMock{
		GetByIDFn: func(_ context.Context, _, _ uint) (*entities.Warehouse, error) {
			return nil, domainerrors.ErrWarehouseNotFound
		},
	}
	uc := New(repo)

	// Act
	err := uc.DeleteLocation(ctx, 999, 3, 10)

	// Assert
	if !errors.Is(err, domainerrors.ErrWarehouseNotFound) {
		t.Errorf("se esperaba ErrWarehouseNotFound, se obtuvo: %v", err)
	}
}

func TestDeleteLocation_GetByIDRepositoryError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	dbErr := errors.New("db unavailable")

	repo := &mocks.RepositoryMock{
		GetByIDFn: func(_ context.Context, _, _ uint) (*entities.Warehouse, error) {
			return nil, dbErr
		},
	}
	uc := New(repo)

	// Act
	err := uc.DeleteLocation(ctx, 5, 3, 10)

	// Assert
	if !errors.Is(err, dbErr) {
		t.Errorf("se esperaba error %v, se obtuvo: %v", dbErr, err)
	}
}

func TestDeleteLocation_LocationHasStock(t *testing.T) {
	// Arrange
	ctx := context.Background()
	repo := &mocks.RepositoryMock{
		GetByIDFn: func(_ context.Context, _, _ uint) (*entities.Warehouse, error) {
			return &entities.Warehouse{ID: 5, BusinessID: 10}, nil
		},
		DeleteLocationFn: func(_ context.Context, _, _ uint) error {
			return domainerrors.ErrLocationHasStock
		},
	}
	uc := New(repo)

	// Act
	err := uc.DeleteLocation(ctx, 5, 3, 10)

	// Assert
	if !errors.Is(err, domainerrors.ErrLocationHasStock) {
		t.Errorf("se esperaba ErrLocationHasStock, se obtuvo: %v", err)
	}
}

func TestDeleteLocation_LocationNotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	repo := &mocks.RepositoryMock{
		GetByIDFn: func(_ context.Context, _, _ uint) (*entities.Warehouse, error) {
			return &entities.Warehouse{ID: 5, BusinessID: 10}, nil
		},
		DeleteLocationFn: func(_ context.Context, _, _ uint) error {
			return domainerrors.ErrLocationNotFound
		},
	}
	uc := New(repo)

	// Act
	err := uc.DeleteLocation(ctx, 5, 999, 10)

	// Assert
	if !errors.Is(err, domainerrors.ErrLocationNotFound) {
		t.Errorf("se esperaba ErrLocationNotFound, se obtuvo: %v", err)
	}
}

func TestDeleteLocation_DeleteRepositoryError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	dbErr := errors.New("delete location failed")

	repo := &mocks.RepositoryMock{
		GetByIDFn: func(_ context.Context, _, _ uint) (*entities.Warehouse, error) {
			return &entities.Warehouse{ID: 5, BusinessID: 10}, nil
		},
		DeleteLocationFn: func(_ context.Context, _, _ uint) error {
			return dbErr
		},
	}
	uc := New(repo)

	// Act
	err := uc.DeleteLocation(ctx, 5, 3, 10)

	// Assert
	if !errors.Is(err, dbErr) {
		t.Errorf("se esperaba error %v, se obtuvo: %v", dbErr, err)
	}
}

func TestDeleteLocation_PropagatesCorrectIDs(t *testing.T) {
	// Arrange — verifica que los IDs se pasen correctamente a GetByID y DeleteLocation
	ctx := context.Background()
	var receivedGetBusinessID, receivedGetWarehouseID uint
	var receivedDelWarehouseID, receivedDelLocationID uint

	repo := &mocks.RepositoryMock{
		GetByIDFn: func(_ context.Context, bID, wID uint) (*entities.Warehouse, error) {
			receivedGetBusinessID = bID
			receivedGetWarehouseID = wID
			return &entities.Warehouse{ID: wID, BusinessID: bID}, nil
		},
		DeleteLocationFn: func(_ context.Context, wID, lID uint) error {
			receivedDelWarehouseID = wID
			receivedDelLocationID = lID
			return nil
		},
	}
	uc := New(repo)

	// Act: DeleteLocation(ctx, warehouseID, locationID, businessID)
	err := uc.DeleteLocation(ctx, 5, 3, 10)

	// Assert
	if err != nil {
		t.Fatalf("se esperaba nil error, se obtuvo: %v", err)
	}
	if receivedGetBusinessID != 10 {
		t.Errorf("GetByID businessID: esperado 10, obtenido %d", receivedGetBusinessID)
	}
	if receivedGetWarehouseID != 5 {
		t.Errorf("GetByID warehouseID: esperado 5, obtenido %d", receivedGetWarehouseID)
	}
	if receivedDelWarehouseID != 5 {
		t.Errorf("DeleteLocation warehouseID: esperado 5, obtenido %d", receivedDelWarehouseID)
	}
	if receivedDelLocationID != 3 {
		t.Errorf("DeleteLocation locationID: esperado 3, obtenido %d", receivedDelLocationID)
	}
}
