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

func TestCreateMovementType_Exitoso_RetornaTipoCreado(t *testing.T) {
	// Arrange
	tipoCapturado := (*entities.StockMovementType)(nil)
	repo := &mocks.RepositoryMock{
		CreateMovementTypeFn: func(ctx context.Context, movType *entities.StockMovementType) (*entities.StockMovementType, error) {
			tipoCapturado = movType
			movType.ID = 10
			return movType, nil
		},
	}
	uc := buildUseCase(repo, nil, nil)

	dto := dtos.CreateStockMovementTypeDTO{
		Code:        "devoluciones",
		Name:        "Devoluciones de clientes",
		Description: "Stock devuelto por clientes",
		Direction:   "in",
	}

	// Act
	result, err := uc.CreateMovementType(context.Background(), dto)

	// Assert
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if result == nil {
		t.Fatal("resultado esperado no nil")
	}
	if result.ID != 10 {
		t.Errorf("ID esperado 10, se obtuvo %d", result.ID)
	}
	if tipoCapturado.Code != "devoluciones" {
		t.Errorf("code esperado 'devoluciones', se obtuvo '%s'", tipoCapturado.Code)
	}
	if tipoCapturado.IsActive != true {
		t.Error("IsActive esperado true en tipos recien creados")
	}
}

func TestCreateMovementType_CodigoYaExiste_RetornaError(t *testing.T) {
	// Arrange
	repo := &mocks.RepositoryMock{
		CreateMovementTypeFn: func(ctx context.Context, movType *entities.StockMovementType) (*entities.StockMovementType, error) {
			return nil, domainerrors.ErrMovementTypeCodeExists
		},
	}
	uc := buildUseCase(repo, nil, nil)

	dto := dtos.CreateStockMovementTypeDTO{
		Code:      "inbound",
		Name:      "Entrada duplicada",
		Direction: "in",
	}

	// Act
	result, err := uc.CreateMovementType(context.Background(), dto)

	// Assert
	if !errors.Is(err, domainerrors.ErrMovementTypeCodeExists) {
		t.Errorf("error esperado ErrMovementTypeCodeExists, se obtuvo: %v", err)
	}
	if result != nil {
		t.Errorf("resultado esperado nil, se obtuvo: %v", result)
	}
}

func TestUpdateMovementType_TipoExiste_ActualizaCampos(t *testing.T) {
	// Arrange
	tipoExistente := &entities.StockMovementType{
		ID:          5,
		Code:        "ajuste",
		Name:        "Ajuste antiguo",
		Description: "Descripcion antigua",
		IsActive:    true,
		Direction:   "neutral",
	}

	repo := &mocks.RepositoryMock{
		GetMovementTypeByIDFn: func(ctx context.Context, id uint) (*entities.StockMovementType, error) {
			return tipoExistente, nil
		},
		UpdateMovementTypeFn: func(ctx context.Context, movType *entities.StockMovementType) error {
			return nil
		},
	}
	uc := buildUseCase(repo, nil, nil)

	isActive := false
	dto := dtos.UpdateStockMovementTypeDTO{
		ID:          5,
		Name:        "Ajuste actualizado",
		Description: "Nueva descripcion",
		IsActive:    &isActive,
		Direction:   "out",
	}

	// Act
	result, err := uc.UpdateMovementType(context.Background(), dto)

	// Assert
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if result.Name != "Ajuste actualizado" {
		t.Errorf("nombre esperado 'Ajuste actualizado', se obtuvo '%s'", result.Name)
	}
	if result.Description != "Nueva descripcion" {
		t.Errorf("descripcion esperada 'Nueva descripcion', se obtuvo '%s'", result.Description)
	}
	if result.IsActive != false {
		t.Errorf("IsActive esperado false, se obtuvo %v", result.IsActive)
	}
	if result.Direction != "out" {
		t.Errorf("direction esperada 'out', se obtuvo '%s'", result.Direction)
	}
}

func TestUpdateMovementType_CamposVacios_NoModifica(t *testing.T) {
	// Arrange
	tipoExistente := &entities.StockMovementType{
		ID:          5,
		Code:        "ajuste",
		Name:        "Ajuste original",
		Description: "Descripcion original",
		Direction:   "neutral",
		IsActive:    true,
	}

	repo := &mocks.RepositoryMock{
		GetMovementTypeByIDFn: func(ctx context.Context, id uint) (*entities.StockMovementType, error) {
			return tipoExistente, nil
		},
		UpdateMovementTypeFn: func(ctx context.Context, movType *entities.StockMovementType) error {
			return nil
		},
	}
	uc := buildUseCase(repo, nil, nil)

	// DTO con campos vacíos: no debe modificar los valores existentes
	dto := dtos.UpdateStockMovementTypeDTO{
		ID:       5,
		Name:     "",        // vacío -> no cambiar
		IsActive: nil,       // nil -> no cambiar
		Direction: "",       // vacío -> no cambiar
	}

	// Act
	result, err := uc.UpdateMovementType(context.Background(), dto)

	// Assert
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if result.Name != "Ajuste original" {
		t.Errorf("nombre no deberia haberse modificado, se obtuvo '%s'", result.Name)
	}
	if result.Direction != "neutral" {
		t.Errorf("direction no deberia haberse modificado, se obtuvo '%s'", result.Direction)
	}
}

func TestUpdateMovementType_TipoNoExiste_RetornaErrMovementTypeNotFound(t *testing.T) {
	// Arrange
	repo := &mocks.RepositoryMock{
		GetMovementTypeByIDFn: func(ctx context.Context, id uint) (*entities.StockMovementType, error) {
			return nil, errors.New("record not found")
		},
	}
	uc := buildUseCase(repo, nil, nil)

	dto := dtos.UpdateStockMovementTypeDTO{
		ID:   999,
		Name: "Inexistente",
	}

	// Act
	result, err := uc.UpdateMovementType(context.Background(), dto)

	// Assert
	if !errors.Is(err, domainerrors.ErrMovementTypeNotFound) {
		t.Errorf("error esperado ErrMovementTypeNotFound, se obtuvo: %v", err)
	}
	if result != nil {
		t.Errorf("resultado esperado nil, se obtuvo: %v", result)
	}
}

func TestDeleteMovementType_Exitoso_RetornaNil(t *testing.T) {
	// Arrange
	idCapturado := uint(0)
	repo := &mocks.RepositoryMock{
		DeleteMovementTypeFn: func(ctx context.Context, id uint) error {
			idCapturado = id
			return nil
		},
	}
	uc := buildUseCase(repo, nil, nil)

	// Act
	err := uc.DeleteMovementType(context.Background(), 7)

	// Assert
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if idCapturado != 7 {
		t.Errorf("ID esperado 7, se obtuvo %d", idCapturado)
	}
}

func TestDeleteMovementType_ErrorEnRepositorio_PropagaError(t *testing.T) {
	// Arrange
	expectedErr := errors.New("no se puede eliminar tipo de movimiento con movimientos asociados")
	repo := &mocks.RepositoryMock{
		DeleteMovementTypeFn: func(ctx context.Context, id uint) error {
			return expectedErr
		},
	}
	uc := buildUseCase(repo, nil, nil)

	// Act
	err := uc.DeleteMovementType(context.Background(), 1)

	// Assert
	if !errors.Is(err, expectedErr) {
		t.Errorf("error esperado %v, se obtuvo: %v", expectedErr, err)
	}
}
