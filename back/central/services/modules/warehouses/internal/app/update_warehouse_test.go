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

func TestUpdateWarehouse_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	existing := &entities.Warehouse{ID: 3, BusinessID: 10, Name: "Bodega Vieja", Code: "BOD-OLD"}
	updated := &entities.Warehouse{ID: 3, BusinessID: 10, Name: "Bodega Nueva", Code: "BOD-NEW"}

	repo := &mocks.RepositoryMock{
		GetByIDFn: func(_ context.Context, _, _ uint) (*entities.Warehouse, error) {
			return existing, nil
		},
		ExistsByCodeFn: func(_ context.Context, _ uint, _ string, _ *uint) (bool, error) {
			return false, nil
		},
		UpdateFn: func(_ context.Context, w *entities.Warehouse) (*entities.Warehouse, error) {
			return updated, nil
		},
	}
	uc := New(repo)

	dto := dtos.UpdateWarehouseDTO{
		ID:         3,
		BusinessID: 10,
		Name:       "Bodega Nueva",
		Code:       "BOD-NEW",
		IsDefault:  false,
	}

	// Act
	result, err := uc.UpdateWarehouse(ctx, dto)

	// Assert
	if err != nil {
		t.Fatalf("se esperaba nil error, se obtuvo: %v", err)
	}
	if result == nil {
		t.Fatal("se esperaba un warehouse, se obtuvo nil")
	}
	if result.Name != "Bodega Nueva" {
		t.Errorf("Name: esperado %q, obtenido %q", "Bodega Nueva", result.Name)
	}
}

func TestUpdateWarehouse_WarehouseNotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	repo := &mocks.RepositoryMock{
		GetByIDFn: func(_ context.Context, _, _ uint) (*entities.Warehouse, error) {
			return nil, domainerrors.ErrWarehouseNotFound
		},
	}
	uc := New(repo)

	dto := dtos.UpdateWarehouseDTO{
		ID:         999,
		BusinessID: 10,
		Code:       "BOD-X",
	}

	// Act
	result, err := uc.UpdateWarehouse(ctx, dto)

	// Assert
	if result != nil {
		t.Errorf("se esperaba nil result, se obtuvo: %+v", result)
	}
	if !errors.Is(err, domainerrors.ErrWarehouseNotFound) {
		t.Errorf("se esperaba ErrWarehouseNotFound, se obtuvo: %v", err)
	}
}

func TestUpdateWarehouse_GetByIDRepositoryError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	dbErr := errors.New("db error on get")

	repo := &mocks.RepositoryMock{
		GetByIDFn: func(_ context.Context, _, _ uint) (*entities.Warehouse, error) {
			return nil, dbErr
		},
	}
	uc := New(repo)

	dto := dtos.UpdateWarehouseDTO{
		ID:         1,
		BusinessID: 10,
		Code:       "BOD-001",
	}

	// Act
	result, err := uc.UpdateWarehouse(ctx, dto)

	// Assert
	if result != nil {
		t.Errorf("se esperaba nil result, se obtuvo: %+v", result)
	}
	if !errors.Is(err, dbErr) {
		t.Errorf("se esperaba error %v, se obtuvo: %v", dbErr, err)
	}
}

func TestUpdateWarehouse_DuplicateCode(t *testing.T) {
	// Arrange
	ctx := context.Background()
	repo := &mocks.RepositoryMock{
		GetByIDFn: func(_ context.Context, _, _ uint) (*entities.Warehouse, error) {
			return &entities.Warehouse{ID: 3, BusinessID: 10}, nil
		},
		ExistsByCodeFn: func(_ context.Context, _ uint, _ string, _ *uint) (bool, error) {
			return true, nil // código ya está en uso por otra bodega
		},
	}
	uc := New(repo)

	dto := dtos.UpdateWarehouseDTO{
		ID:         3,
		BusinessID: 10,
		Code:       "BOD-DUPLICADO",
	}

	// Act
	result, err := uc.UpdateWarehouse(ctx, dto)

	// Assert
	if result != nil {
		t.Errorf("se esperaba nil result, se obtuvo: %+v", result)
	}
	if !errors.Is(err, domainerrors.ErrDuplicateCode) {
		t.Errorf("se esperaba ErrDuplicateCode, se obtuvo: %v", err)
	}
}

func TestUpdateWarehouse_ExistsByCodeError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	dbErr := errors.New("db error on exists check")

	repo := &mocks.RepositoryMock{
		GetByIDFn: func(_ context.Context, _, _ uint) (*entities.Warehouse, error) {
			return &entities.Warehouse{ID: 3, BusinessID: 10}, nil
		},
		ExistsByCodeFn: func(_ context.Context, _ uint, _ string, _ *uint) (bool, error) {
			return false, dbErr
		},
	}
	uc := New(repo)

	dto := dtos.UpdateWarehouseDTO{
		ID:         3,
		BusinessID: 10,
		Code:       "BOD-001",
	}

	// Act
	result, err := uc.UpdateWarehouse(ctx, dto)

	// Assert
	if result != nil {
		t.Errorf("se esperaba nil result, se obtuvo: %+v", result)
	}
	if !errors.Is(err, dbErr) {
		t.Errorf("se esperaba error %v, se obtuvo: %v", dbErr, err)
	}
}

func TestUpdateWarehouse_IsDefault_ClearsOtherDefaults(t *testing.T) {
	// Arrange
	ctx := context.Background()
	clearDefaultCalled := false
	var clearDefaultExcludeID uint

	repo := &mocks.RepositoryMock{
		GetByIDFn: func(_ context.Context, _, _ uint) (*entities.Warehouse, error) {
			return &entities.Warehouse{ID: 5, BusinessID: 10}, nil
		},
		ExistsByCodeFn: func(_ context.Context, _ uint, _ string, _ *uint) (bool, error) {
			return false, nil
		},
		ClearDefaultFn: func(_ context.Context, _ uint, exID uint) error {
			clearDefaultCalled = true
			clearDefaultExcludeID = exID
			return nil
		},
		UpdateFn: func(_ context.Context, w *entities.Warehouse) (*entities.Warehouse, error) {
			return w, nil
		},
	}
	uc := New(repo)

	dto := dtos.UpdateWarehouseDTO{
		ID:         5,
		BusinessID: 10,
		Code:       "BOD-001",
		IsDefault:  true, // se marca como default
	}

	// Act
	_, err := uc.UpdateWarehouse(ctx, dto)

	// Assert
	if err != nil {
		t.Fatalf("se esperaba nil error, se obtuvo: %v", err)
	}
	if !clearDefaultCalled {
		t.Error("se esperaba que ClearDefault fuera llamado")
	}
	// Al actualizar, el excludeID debe ser el ID de la bodega actualizada
	if clearDefaultExcludeID != dto.ID {
		t.Errorf("ClearDefault excludeID: esperado %d, obtenido %d", dto.ID, clearDefaultExcludeID)
	}
}

func TestUpdateWarehouse_IsNotDefault_DoesNotClearDefaults(t *testing.T) {
	// Arrange
	ctx := context.Background()
	clearDefaultCalled := false

	repo := &mocks.RepositoryMock{
		GetByIDFn: func(_ context.Context, _, _ uint) (*entities.Warehouse, error) {
			return &entities.Warehouse{ID: 5, BusinessID: 10}, nil
		},
		ExistsByCodeFn: func(_ context.Context, _ uint, _ string, _ *uint) (bool, error) {
			return false, nil
		},
		ClearDefaultFn: func(_ context.Context, _ uint, _ uint) error {
			clearDefaultCalled = true
			return nil
		},
		UpdateFn: func(_ context.Context, w *entities.Warehouse) (*entities.Warehouse, error) {
			return w, nil
		},
	}
	uc := New(repo)

	dto := dtos.UpdateWarehouseDTO{
		ID:         5,
		BusinessID: 10,
		Code:       "BOD-001",
		IsDefault:  false, // NO default
	}

	// Act
	_, err := uc.UpdateWarehouse(ctx, dto)

	// Assert
	if err != nil {
		t.Fatalf("se esperaba nil error, se obtuvo: %v", err)
	}
	if clearDefaultCalled {
		t.Error("ClearDefault NO debería llamarse cuando IsDefault es false")
	}
}

func TestUpdateWarehouse_ClearDefaultError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	dbErr := errors.New("error clearing defaults")

	repo := &mocks.RepositoryMock{
		GetByIDFn: func(_ context.Context, _, _ uint) (*entities.Warehouse, error) {
			return &entities.Warehouse{ID: 5, BusinessID: 10}, nil
		},
		ExistsByCodeFn: func(_ context.Context, _ uint, _ string, _ *uint) (bool, error) {
			return false, nil
		},
		ClearDefaultFn: func(_ context.Context, _ uint, _ uint) error {
			return dbErr
		},
	}
	uc := New(repo)

	dto := dtos.UpdateWarehouseDTO{
		ID:         5,
		BusinessID: 10,
		Code:       "BOD-001",
		IsDefault:  true,
	}

	// Act
	result, err := uc.UpdateWarehouse(ctx, dto)

	// Assert
	if result != nil {
		t.Errorf("se esperaba nil result, se obtuvo: %+v", result)
	}
	if !errors.Is(err, dbErr) {
		t.Errorf("se esperaba error %v, se obtuvo: %v", dbErr, err)
	}
}

func TestUpdateWarehouse_UpdateRepositoryError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	dbErr := errors.New("update failed")

	repo := &mocks.RepositoryMock{
		GetByIDFn: func(_ context.Context, _, _ uint) (*entities.Warehouse, error) {
			return &entities.Warehouse{ID: 3, BusinessID: 10}, nil
		},
		ExistsByCodeFn: func(_ context.Context, _ uint, _ string, _ *uint) (bool, error) {
			return false, nil
		},
		UpdateFn: func(_ context.Context, _ *entities.Warehouse) (*entities.Warehouse, error) {
			return nil, dbErr
		},
	}
	uc := New(repo)

	dto := dtos.UpdateWarehouseDTO{
		ID:         3,
		BusinessID: 10,
		Code:       "BOD-001",
		IsDefault:  false,
	}

	// Act
	result, err := uc.UpdateWarehouse(ctx, dto)

	// Assert
	if result != nil {
		t.Errorf("se esperaba nil result, se obtuvo: %+v", result)
	}
	if !errors.Is(err, dbErr) {
		t.Errorf("se esperaba error %v, se obtuvo: %v", dbErr, err)
	}
}

func TestUpdateWarehouse_ExistsByCode_ReceivesExcludeID(t *testing.T) {
	// Arrange — valida que ExistsByCode se llame con el ID correcto como excludeID
	ctx := context.Background()
	var receivedExcludeID *uint

	repo := &mocks.RepositoryMock{
		GetByIDFn: func(_ context.Context, _, _ uint) (*entities.Warehouse, error) {
			return &entities.Warehouse{ID: 7, BusinessID: 10}, nil
		},
		ExistsByCodeFn: func(_ context.Context, _ uint, _ string, excludeID *uint) (bool, error) {
			receivedExcludeID = excludeID
			return false, nil
		},
		UpdateFn: func(_ context.Context, w *entities.Warehouse) (*entities.Warehouse, error) {
			return w, nil
		},
	}
	uc := New(repo)

	dto := dtos.UpdateWarehouseDTO{
		ID:         7,
		BusinessID: 10,
		Code:       "BOD-007",
		IsDefault:  false,
	}

	// Act
	_, err := uc.UpdateWarehouse(ctx, dto)

	// Assert
	if err != nil {
		t.Fatalf("se esperaba nil error, se obtuvo: %v", err)
	}
	if receivedExcludeID == nil {
		t.Fatal("ExistsByCode debería recibir un excludeID no nil")
	}
	if *receivedExcludeID != dto.ID {
		t.Errorf("excludeID: esperado %d, obtenido %d", dto.ID, *receivedExcludeID)
	}
}
