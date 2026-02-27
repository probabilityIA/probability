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

func TestCreateWarehouse_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	repo := &mocks.RepositoryMock{
		ExistsByCodeFn: func(_ context.Context, _ uint, _ string, _ *uint) (bool, error) {
			return false, nil
		},
		CreateFn: func(_ context.Context, w *entities.Warehouse) (*entities.Warehouse, error) {
			w.ID = 1
			return w, nil
		},
	}
	uc := New(repo)

	dto := dtos.CreateWarehouseDTO{
		BusinessID: 10,
		Name:       "Bodega Principal",
		Code:       "BOD-001",
		Address:    "Calle 1 # 2-3",
		City:       "Bogotá",
		IsActive:   true,
		IsDefault:  false,
	}

	// Act
	result, err := uc.CreateWarehouse(ctx, dto)

	// Assert
	if err != nil {
		t.Fatalf("se esperaba nil error, se obtuvo: %v", err)
	}
	if result == nil {
		t.Fatal("se esperaba un warehouse, se obtuvo nil")
	}
	if result.ID != 1 {
		t.Errorf("se esperaba ID=1, se obtuvo ID=%d", result.ID)
	}
	if result.Code != dto.Code {
		t.Errorf("se esperaba Code=%q, se obtuvo Code=%q", dto.Code, result.Code)
	}
	if result.BusinessID != dto.BusinessID {
		t.Errorf("se esperaba BusinessID=%d, se obtuvo BusinessID=%d", dto.BusinessID, result.BusinessID)
	}
}

func TestCreateWarehouse_DuplicateCode(t *testing.T) {
	// Arrange
	ctx := context.Background()
	repo := &mocks.RepositoryMock{
		ExistsByCodeFn: func(_ context.Context, _ uint, _ string, _ *uint) (bool, error) {
			return true, nil // código ya existe
		},
	}
	uc := New(repo)

	dto := dtos.CreateWarehouseDTO{
		BusinessID: 10,
		Name:       "Bodega Duplicada",
		Code:       "BOD-001",
	}

	// Act
	result, err := uc.CreateWarehouse(ctx, dto)

	// Assert
	if result != nil {
		t.Errorf("se esperaba nil result, se obtuvo: %+v", result)
	}
	if !errors.Is(err, domainerrors.ErrDuplicateCode) {
		t.Errorf("se esperaba ErrDuplicateCode, se obtuvo: %v", err)
	}
}

func TestCreateWarehouse_ExistsByCodeRepositoryError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	dbErr := errors.New("connection timeout")
	repo := &mocks.RepositoryMock{
		ExistsByCodeFn: func(_ context.Context, _ uint, _ string, _ *uint) (bool, error) {
			return false, dbErr
		},
	}
	uc := New(repo)

	dto := dtos.CreateWarehouseDTO{
		BusinessID: 10,
		Code:       "BOD-001",
	}

	// Act
	result, err := uc.CreateWarehouse(ctx, dto)

	// Assert
	if result != nil {
		t.Errorf("se esperaba nil result, se obtuvo: %+v", result)
	}
	if !errors.Is(err, dbErr) {
		t.Errorf("se esperaba error %v, se obtuvo: %v", dbErr, err)
	}
}

func TestCreateWarehouse_IsDefault_ClearsOtherDefaults(t *testing.T) {
	// Arrange
	ctx := context.Background()
	clearDefaultCalled := false
	var clearDefaultBusinessID uint
	var clearDefaultExcludeID uint

	repo := &mocks.RepositoryMock{
		ExistsByCodeFn: func(_ context.Context, _ uint, _ string, _ *uint) (bool, error) {
			return false, nil
		},
		ClearDefaultFn: func(_ context.Context, bID uint, exID uint) error {
			clearDefaultCalled = true
			clearDefaultBusinessID = bID
			clearDefaultExcludeID = exID
			return nil
		},
		CreateFn: func(_ context.Context, w *entities.Warehouse) (*entities.Warehouse, error) {
			w.ID = 5
			return w, nil
		},
	}
	uc := New(repo)

	dto := dtos.CreateWarehouseDTO{
		BusinessID: 10,
		Code:       "BOD-002",
		IsDefault:  true, // se marca como default
	}

	// Act
	result, err := uc.CreateWarehouse(ctx, dto)

	// Assert
	if err != nil {
		t.Fatalf("se esperaba nil error, se obtuvo: %v", err)
	}
	if result == nil {
		t.Fatal("se esperaba un warehouse, se obtuvo nil")
	}
	if !clearDefaultCalled {
		t.Error("se esperaba que ClearDefault fuera llamado al crear bodega default")
	}
	if clearDefaultBusinessID != dto.BusinessID {
		t.Errorf("ClearDefault recibió businessID=%d, se esperaba %d", clearDefaultBusinessID, dto.BusinessID)
	}
	// Al crear (ID=0 aún), el excludeID debe ser 0
	if clearDefaultExcludeID != 0 {
		t.Errorf("ClearDefault recibió excludeID=%d, se esperaba 0", clearDefaultExcludeID)
	}
}

func TestCreateWarehouse_IsNotDefault_DoesNotClearDefaults(t *testing.T) {
	// Arrange
	ctx := context.Background()
	clearDefaultCalled := false

	repo := &mocks.RepositoryMock{
		ExistsByCodeFn: func(_ context.Context, _ uint, _ string, _ *uint) (bool, error) {
			return false, nil
		},
		ClearDefaultFn: func(_ context.Context, _ uint, _ uint) error {
			clearDefaultCalled = true
			return nil
		},
		CreateFn: func(_ context.Context, w *entities.Warehouse) (*entities.Warehouse, error) {
			return w, nil
		},
	}
	uc := New(repo)

	dto := dtos.CreateWarehouseDTO{
		BusinessID: 10,
		Code:       "BOD-003",
		IsDefault:  false, // NO es default
	}

	// Act
	_, err := uc.CreateWarehouse(ctx, dto)

	// Assert
	if err != nil {
		t.Fatalf("se esperaba nil error, se obtuvo: %v", err)
	}
	if clearDefaultCalled {
		t.Error("ClearDefault NO debería llamarse cuando IsDefault es false")
	}
}

func TestCreateWarehouse_ClearDefaultError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	dbErr := errors.New("failed to clear defaults")

	repo := &mocks.RepositoryMock{
		ExistsByCodeFn: func(_ context.Context, _ uint, _ string, _ *uint) (bool, error) {
			return false, nil
		},
		ClearDefaultFn: func(_ context.Context, _ uint, _ uint) error {
			return dbErr
		},
	}
	uc := New(repo)

	dto := dtos.CreateWarehouseDTO{
		BusinessID: 10,
		Code:       "BOD-004",
		IsDefault:  true,
	}

	// Act
	result, err := uc.CreateWarehouse(ctx, dto)

	// Assert
	if result != nil {
		t.Errorf("se esperaba nil result, se obtuvo: %+v", result)
	}
	if !errors.Is(err, dbErr) {
		t.Errorf("se esperaba error %v, se obtuvo: %v", dbErr, err)
	}
}

func TestCreateWarehouse_CreateRepositoryError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	dbErr := errors.New("insert failed")

	repo := &mocks.RepositoryMock{
		ExistsByCodeFn: func(_ context.Context, _ uint, _ string, _ *uint) (bool, error) {
			return false, nil
		},
		CreateFn: func(_ context.Context, _ *entities.Warehouse) (*entities.Warehouse, error) {
			return nil, dbErr
		},
	}
	uc := New(repo)

	dto := dtos.CreateWarehouseDTO{
		BusinessID: 10,
		Code:       "BOD-005",
	}

	// Act
	result, err := uc.CreateWarehouse(ctx, dto)

	// Assert
	if result != nil {
		t.Errorf("se esperaba nil result, se obtuvo: %+v", result)
	}
	if !errors.Is(err, dbErr) {
		t.Errorf("se esperaba error %v, se obtuvo: %v", dbErr, err)
	}
}

func TestCreateWarehouse_MapsAllDTOFields(t *testing.T) {
	// Arrange
	ctx := context.Background()
	cap := 500
	_ = cap // no se usa para warehouse pero ilustra la intención

	var captured *entities.Warehouse
	repo := &mocks.RepositoryMock{
		ExistsByCodeFn: func(_ context.Context, _ uint, _ string, _ *uint) (bool, error) {
			return false, nil
		},
		CreateFn: func(_ context.Context, w *entities.Warehouse) (*entities.Warehouse, error) {
			captured = w
			w.ID = 99
			return w, nil
		},
	}
	uc := New(repo)

	dto := dtos.CreateWarehouseDTO{
		BusinessID:    42,
		Name:          "Bodega Test",
		Code:          "BT-001",
		Address:       "Av. Siempre Viva 742",
		City:          "Medellín",
		State:         "Antioquia",
		Country:       "Colombia",
		ZipCode:       "050001",
		Phone:         "+573001234567",
		ContactName:   "Homero Simpson",
		ContactEmail:  "homero@springfield.com",
		IsActive:      true,
		IsDefault:     false,
		IsFulfillment: true,
	}

	// Act
	_, err := uc.CreateWarehouse(ctx, dto)

	// Assert
	if err != nil {
		t.Fatalf("se esperaba nil error, se obtuvo: %v", err)
	}
	if captured == nil {
		t.Fatal("Create no fue llamado")
	}
	if captured.BusinessID != dto.BusinessID {
		t.Errorf("BusinessID: esperado %d, obtenido %d", dto.BusinessID, captured.BusinessID)
	}
	if captured.Name != dto.Name {
		t.Errorf("Name: esperado %q, obtenido %q", dto.Name, captured.Name)
	}
	if captured.Code != dto.Code {
		t.Errorf("Code: esperado %q, obtenido %q", dto.Code, captured.Code)
	}
	if captured.Address != dto.Address {
		t.Errorf("Address: esperado %q, obtenido %q", dto.Address, captured.Address)
	}
	if captured.City != dto.City {
		t.Errorf("City: esperado %q, obtenido %q", dto.City, captured.City)
	}
	if captured.State != dto.State {
		t.Errorf("State: esperado %q, obtenido %q", dto.State, captured.State)
	}
	if captured.Country != dto.Country {
		t.Errorf("Country: esperado %q, obtenido %q", dto.Country, captured.Country)
	}
	if captured.ZipCode != dto.ZipCode {
		t.Errorf("ZipCode: esperado %q, obtenido %q", dto.ZipCode, captured.ZipCode)
	}
	if captured.Phone != dto.Phone {
		t.Errorf("Phone: esperado %q, obtenido %q", dto.Phone, captured.Phone)
	}
	if captured.ContactName != dto.ContactName {
		t.Errorf("ContactName: esperado %q, obtenido %q", dto.ContactName, captured.ContactName)
	}
	if captured.ContactEmail != dto.ContactEmail {
		t.Errorf("ContactEmail: esperado %q, obtenido %q", dto.ContactEmail, captured.ContactEmail)
	}
	if captured.IsActive != dto.IsActive {
		t.Errorf("IsActive: esperado %v, obtenido %v", dto.IsActive, captured.IsActive)
	}
	if captured.IsDefault != dto.IsDefault {
		t.Errorf("IsDefault: esperado %v, obtenido %v", dto.IsDefault, captured.IsDefault)
	}
	if captured.IsFulfillment != dto.IsFulfillment {
		t.Errorf("IsFulfillment: esperado %v, obtenido %v", dto.IsFulfillment, captured.IsFulfillment)
	}
}
