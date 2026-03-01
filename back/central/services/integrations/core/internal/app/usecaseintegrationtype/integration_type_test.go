package usecaseintegrationtype

import (
	"context"
	"errors"
	"testing"

	"github.com/secamc93/probability/back/central/services/integrations/core/internal/domain"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// newTestIntegrationTypeUseCase construye el use case con mocks para tests
func newTestIntegrationTypeUseCase(
	repo *mocks.RepositoryMock,
	s3 *mocks.S3Mock,
	logger *mocks.LoggerMock,
	cfg *mocks.ConfigMock,
) IIntegrationTypeUseCase {
	cache := new(mocks.CacheMock)
	enc := new(mocks.EncryptionMock)
	return New(repo, s3, cache, logger, cfg, enc)
}

// configurarLoggerPermisivoType permite cualquier llamada al logger sin fallo de test
func configurarLoggerPermisivoType(logger *mocks.LoggerMock) {
	logger.On("Info", mock.Anything).Maybe()
	logger.On("Error", mock.Anything).Maybe()
	logger.On("Warn", mock.Anything).Maybe()
	logger.On("Debug", mock.Anything).Maybe()
}

// ============================================
// CreateIntegrationType
// ============================================

func TestCreateIntegrationType_ExitoSinImagen(t *testing.T) {
	// Arrange
	repo := new(mocks.RepositoryMock)
	s3 := new(mocks.S3Mock)
	logger := new(mocks.LoggerMock)
	cfg := new(mocks.ConfigMock)
	configurarLoggerPermisivoType(logger)

	uc := newTestIntegrationTypeUseCase(repo, s3, logger, cfg)
	ctx := context.Background()

	dto := domain.CreateIntegrationTypeDTO{
		Name:       "Factus Facturación",
		Code:       "factus",
		CategoryID: 3,
		IsActive:   true,
	}

	// Nombre no existe
	repo.On("GetIntegrationTypeByName", mock.Anything, dto.Name).Return(nil, errors.New("no encontrado"))
	// Código no existe
	repo.On("GetIntegrationTypeByCode", mock.Anything, dto.Code).Return(nil, errors.New("no encontrado"))
	// Guardar exitoso
	repo.On("CreateIntegrationType", mock.Anything, mock.AnythingOfType("*domain.IntegrationType")).Return(nil)

	// Act
	resultado, err := uc.CreateIntegrationType(ctx, dto)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, resultado)
	assert.Equal(t, dto.Name, resultado.Name)
	assert.Equal(t, dto.Code, resultado.Code)
	assert.Equal(t, dto.CategoryID, resultado.CategoryID)
	repo.AssertExpectations(t)
}

func TestCreateIntegrationType_NombreYaExiste(t *testing.T) {
	// Arrange
	repo := new(mocks.RepositoryMock)
	s3 := new(mocks.S3Mock)
	logger := new(mocks.LoggerMock)
	cfg := new(mocks.ConfigMock)
	configurarLoggerPermisivoType(logger)

	uc := newTestIntegrationTypeUseCase(repo, s3, logger, cfg)
	ctx := context.Background()

	dto := domain.CreateIntegrationTypeDTO{
		Name: "Shopify",
		Code: "shopify",
	}

	tipoExistente := &domain.IntegrationType{ID: 1, Name: "Shopify"}
	repo.On("GetIntegrationTypeByName", mock.Anything, dto.Name).Return(tipoExistente, nil)

	// Act
	resultado, err := uc.CreateIntegrationType(ctx, dto)

	// Assert
	assert.ErrorIs(t, err, domain.ErrIntegrationTypeNameExists)
	assert.Nil(t, resultado)
	repo.AssertNotCalled(t, "CreateIntegrationType", mock.Anything, mock.Anything)
}

func TestCreateIntegrationType_CodigoYaExiste(t *testing.T) {
	// Arrange
	repo := new(mocks.RepositoryMock)
	s3 := new(mocks.S3Mock)
	logger := new(mocks.LoggerMock)
	cfg := new(mocks.ConfigMock)
	configurarLoggerPermisivoType(logger)

	uc := newTestIntegrationTypeUseCase(repo, s3, logger, cfg)
	ctx := context.Background()

	dto := domain.CreateIntegrationTypeDTO{
		Name: "Factus v2",
		Code: "factus", // Código duplicado
	}

	// Nombre no existe
	repo.On("GetIntegrationTypeByName", mock.Anything, dto.Name).Return(nil, errors.New("no encontrado"))
	// Código sí existe
	tipoExistente := &domain.IntegrationType{ID: 7, Code: "factus"}
	repo.On("GetIntegrationTypeByCode", mock.Anything, dto.Code).Return(tipoExistente, nil)

	// Act
	resultado, err := uc.CreateIntegrationType(ctx, dto)

	// Assert
	assert.ErrorIs(t, err, domain.ErrIntegrationTypeCodeExists)
	assert.Nil(t, resultado)
}

func TestCreateIntegrationType_CodigoGeneradoAutomaticamente(t *testing.T) {
	// Arrange — no se proporciona código, debe generarse del nombre
	repo := new(mocks.RepositoryMock)
	s3 := new(mocks.S3Mock)
	logger := new(mocks.LoggerMock)
	cfg := new(mocks.ConfigMock)
	configurarLoggerPermisivoType(logger)

	uc := newTestIntegrationTypeUseCase(repo, s3, logger, cfg)
	ctx := context.Background()

	dto := domain.CreateIntegrationTypeDTO{
		Name:       "Mi Integración Nueva",
		Code:       "", // Sin código — se genera automáticamente
		CategoryID: 1,
	}

	repo.On("GetIntegrationTypeByName", mock.Anything, dto.Name).Return(nil, errors.New("no encontrado"))
	repo.On("CreateIntegrationType", mock.Anything, mock.MatchedBy(func(it *domain.IntegrationType) bool {
		// El código generado debe empezar con "mi_integraci" (derivado del nombre)
		return it.Code != "" && it.Name == dto.Name
	})).Return(nil)

	// Act
	resultado, err := uc.CreateIntegrationType(ctx, dto)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, resultado)
	assert.NotEmpty(t, resultado.Code, "El código debe haberse generado automáticamente")
}

func TestCreateIntegrationType_ErrorAlGuardar(t *testing.T) {
	// Arrange
	repo := new(mocks.RepositoryMock)
	s3 := new(mocks.S3Mock)
	logger := new(mocks.LoggerMock)
	cfg := new(mocks.ConfigMock)
	configurarLoggerPermisivoType(logger)

	uc := newTestIntegrationTypeUseCase(repo, s3, logger, cfg)
	ctx := context.Background()

	dto := domain.CreateIntegrationTypeDTO{
		Name: "Nuevo Tipo",
		Code: "nuevo_tipo",
	}

	repo.On("GetIntegrationTypeByName", mock.Anything, dto.Name).Return(nil, errors.New("no encontrado"))
	repo.On("GetIntegrationTypeByCode", mock.Anything, dto.Code).Return(nil, errors.New("no encontrado"))
	repo.On("CreateIntegrationType", mock.Anything, mock.AnythingOfType("*domain.IntegrationType")).Return(errors.New("db error"))

	// Act
	resultado, err := uc.CreateIntegrationType(ctx, dto)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, resultado)
	assert.Contains(t, err.Error(), "error al guardar")
}

// ============================================
// GetIntegrationTypeByID
// ============================================

func TestGetIntegrationTypeByID_Exitoso(t *testing.T) {
	// Arrange
	repo := new(mocks.RepositoryMock)
	s3 := new(mocks.S3Mock)
	logger := new(mocks.LoggerMock)
	cfg := new(mocks.ConfigMock)
	configurarLoggerPermisivoType(logger)

	uc := newTestIntegrationTypeUseCase(repo, s3, logger, cfg)
	ctx := context.Background()

	tipoIntegracion := &domain.IntegrationType{
		ID:   1,
		Name: "Shopify",
		Code: "shopify",
	}
	repo.On("GetIntegrationTypeByID", mock.Anything, uint(1)).Return(tipoIntegracion, nil)

	// Act
	resultado, err := uc.GetIntegrationTypeByID(ctx, 1)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, resultado)
	assert.Equal(t, uint(1), resultado.ID)
	assert.Equal(t, "shopify", resultado.Code)
}

func TestGetIntegrationTypeByID_NoEncontrado(t *testing.T) {
	// Arrange
	repo := new(mocks.RepositoryMock)
	s3 := new(mocks.S3Mock)
	logger := new(mocks.LoggerMock)
	cfg := new(mocks.ConfigMock)
	configurarLoggerPermisivoType(logger)

	uc := newTestIntegrationTypeUseCase(repo, s3, logger, cfg)
	ctx := context.Background()

	repo.On("GetIntegrationTypeByID", mock.Anything, uint(99)).Return(nil, errors.New("not found"))

	// Act
	resultado, err := uc.GetIntegrationTypeByID(ctx, 99)

	// Assert
	assert.ErrorIs(t, err, domain.ErrIntegrationTypeNotFound)
	assert.Nil(t, resultado)
}

// ============================================
// GetIntegrationTypeByCode
// ============================================

func TestGetIntegrationTypeByCode_Exitoso(t *testing.T) {
	// Arrange
	repo := new(mocks.RepositoryMock)
	s3 := new(mocks.S3Mock)
	logger := new(mocks.LoggerMock)
	cfg := new(mocks.ConfigMock)
	configurarLoggerPermisivoType(logger)

	uc := newTestIntegrationTypeUseCase(repo, s3, logger, cfg)
	ctx := context.Background()

	tipoIntegracion := &domain.IntegrationType{ID: 7, Name: "Factus", Code: "factus"}
	repo.On("GetIntegrationTypeByCode", mock.Anything, "factus").Return(tipoIntegracion, nil)

	// Act
	resultado, err := uc.GetIntegrationTypeByCode(ctx, "factus")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, resultado)
	assert.Equal(t, "factus", resultado.Code)
}

func TestGetIntegrationTypeByCode_NoEncontrado(t *testing.T) {
	// Arrange
	repo := new(mocks.RepositoryMock)
	s3 := new(mocks.S3Mock)
	logger := new(mocks.LoggerMock)
	cfg := new(mocks.ConfigMock)
	configurarLoggerPermisivoType(logger)

	uc := newTestIntegrationTypeUseCase(repo, s3, logger, cfg)
	ctx := context.Background()

	repo.On("GetIntegrationTypeByCode", mock.Anything, "inexistente").Return(nil, errors.New("no encontrado"))

	// Act
	resultado, err := uc.GetIntegrationTypeByCode(ctx, "inexistente")

	// Assert
	assert.ErrorIs(t, err, domain.ErrIntegrationTypeNotFound)
	assert.Nil(t, resultado)
}

// ============================================
// DeleteIntegrationType
// ============================================

func TestDeleteIntegrationType_Exitoso(t *testing.T) {
	// Arrange
	repo := new(mocks.RepositoryMock)
	s3 := new(mocks.S3Mock)
	logger := new(mocks.LoggerMock)
	cfg := new(mocks.ConfigMock)
	configurarLoggerPermisivoType(logger)

	uc := newTestIntegrationTypeUseCase(repo, s3, logger, cfg)
	ctx := context.Background()

	tipoIntegracion := &domain.IntegrationType{ID: 5, Name: "Siigo", Code: "siigo"}
	repo.On("GetIntegrationTypeByID", mock.Anything, uint(5)).Return(tipoIntegracion, nil)
	repo.On("DeleteIntegrationType", mock.Anything, uint(5)).Return(nil)

	// Act
	err := uc.DeleteIntegrationType(ctx, 5)

	// Assert
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestDeleteIntegrationType_NoEncontrado(t *testing.T) {
	// Arrange
	repo := new(mocks.RepositoryMock)
	s3 := new(mocks.S3Mock)
	logger := new(mocks.LoggerMock)
	cfg := new(mocks.ConfigMock)
	configurarLoggerPermisivoType(logger)

	uc := newTestIntegrationTypeUseCase(repo, s3, logger, cfg)
	ctx := context.Background()

	repo.On("GetIntegrationTypeByID", mock.Anything, uint(99)).Return(nil, errors.New("not found"))

	// Act
	err := uc.DeleteIntegrationType(ctx, 99)

	// Assert
	assert.ErrorIs(t, err, domain.ErrIntegrationTypeNotFound)
	repo.AssertNotCalled(t, "DeleteIntegrationType", mock.Anything, mock.Anything)
}

func TestDeleteIntegrationType_ErrorAlEliminar(t *testing.T) {
	// Arrange
	repo := new(mocks.RepositoryMock)
	s3 := new(mocks.S3Mock)
	logger := new(mocks.LoggerMock)
	cfg := new(mocks.ConfigMock)
	configurarLoggerPermisivoType(logger)

	uc := newTestIntegrationTypeUseCase(repo, s3, logger, cfg)
	ctx := context.Background()

	tipoIntegracion := &domain.IntegrationType{ID: 3}
	repo.On("GetIntegrationTypeByID", mock.Anything, uint(3)).Return(tipoIntegracion, nil)
	repo.On("DeleteIntegrationType", mock.Anything, uint(3)).Return(errors.New("foreign key constraint"))

	// Act
	err := uc.DeleteIntegrationType(ctx, 3)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error al eliminar tipo de integración")
}

// ============================================
// ListIntegrationTypes
// ============================================

func TestListIntegrationTypes_RetornaLista(t *testing.T) {
	// Arrange
	repo := new(mocks.RepositoryMock)
	s3 := new(mocks.S3Mock)
	logger := new(mocks.LoggerMock)
	cfg := new(mocks.ConfigMock)
	configurarLoggerPermisivoType(logger)

	uc := newTestIntegrationTypeUseCase(repo, s3, logger, cfg)
	ctx := context.Background()

	tipos := []*domain.IntegrationType{
		{ID: 1, Name: "Shopify", Code: "shopify"},
		{ID: 7, Name: "Factus", Code: "factus"},
		{ID: 8, Name: "Siigo", Code: "siigo"},
	}
	repo.On("ListIntegrationTypes", mock.Anything, mock.Anything).Return(tipos, nil)

	// Act
	resultado, err := uc.ListIntegrationTypes(ctx, nil)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, resultado, 3)
}

func TestListIntegrationTypes_ErrorDeRepositorio(t *testing.T) {
	// Arrange
	repo := new(mocks.RepositoryMock)
	s3 := new(mocks.S3Mock)
	logger := new(mocks.LoggerMock)
	cfg := new(mocks.ConfigMock)
	configurarLoggerPermisivoType(logger)

	uc := newTestIntegrationTypeUseCase(repo, s3, logger, cfg)
	ctx := context.Background()

	repo.On("ListIntegrationTypes", mock.Anything, mock.Anything).Return(nil, errors.New("db error"))

	// Act
	resultado, err := uc.ListIntegrationTypes(ctx, nil)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, resultado)
}

func TestListActiveIntegrationTypes_SoloActivos(t *testing.T) {
	// Arrange
	repo := new(mocks.RepositoryMock)
	s3 := new(mocks.S3Mock)
	logger := new(mocks.LoggerMock)
	cfg := new(mocks.ConfigMock)
	configurarLoggerPermisivoType(logger)

	uc := newTestIntegrationTypeUseCase(repo, s3, logger, cfg)
	ctx := context.Background()

	tipos := []*domain.IntegrationType{
		{ID: 1, Name: "Shopify", Code: "shopify", IsActive: true},
	}
	repo.On("ListActiveIntegrationTypes", mock.Anything).Return(tipos, nil)

	// Act
	resultado, err := uc.ListActiveIntegrationTypes(ctx)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, resultado, 1)
	assert.True(t, resultado[0].IsActive)
}

// ============================================
// ListIntegrationCategories
// ============================================

func TestListIntegrationCategories_RetornaLista(t *testing.T) {
	// Arrange
	repo := new(mocks.RepositoryMock)
	s3 := new(mocks.S3Mock)
	logger := new(mocks.LoggerMock)
	cfg := new(mocks.ConfigMock)
	configurarLoggerPermisivoType(logger)

	uc := newTestIntegrationTypeUseCase(repo, s3, logger, cfg)
	ctx := context.Background()

	categorias := []*domain.IntegrationCategory{
		{ID: 1, Code: "ecommerce", Name: "E-Commerce"},
		{ID: 2, Code: "invoicing", Name: "Facturación"},
	}
	repo.On("ListIntegrationCategories", mock.Anything).Return(categorias, nil)

	// Act
	resultado, err := uc.ListIntegrationCategories(ctx)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, resultado, 2)
	assert.Equal(t, "ecommerce", resultado[0].Code)
}

func TestListIntegrationCategories_ErrorDeRepositorio(t *testing.T) {
	// Arrange
	repo := new(mocks.RepositoryMock)
	s3 := new(mocks.S3Mock)
	logger := new(mocks.LoggerMock)
	cfg := new(mocks.ConfigMock)
	configurarLoggerPermisivoType(logger)

	uc := newTestIntegrationTypeUseCase(repo, s3, logger, cfg)
	ctx := context.Background()

	repo.On("ListIntegrationCategories", mock.Anything).Return(nil, errors.New("db error"))

	// Act
	resultado, err := uc.ListIntegrationCategories(ctx)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, resultado)
}

// ============================================
// UpdateIntegrationType
// ============================================

func TestUpdateIntegrationType_ActualizaNombreYDescripcion(t *testing.T) {
	// Arrange
	repo := new(mocks.RepositoryMock)
	s3 := new(mocks.S3Mock)
	logger := new(mocks.LoggerMock)
	cfg := new(mocks.ConfigMock)
	configurarLoggerPermisivoType(logger)

	uc := newTestIntegrationTypeUseCase(repo, s3, logger, cfg)
	ctx := context.Background()

	tipoExistente := &domain.IntegrationType{
		ID:          7,
		Name:        "Factus v1",
		Code:        "factus",
		Description: "Descripción antigua",
	}

	nuevoNombre := "Factus v2"
	nuevaDesc := "Descripción actualizada"
	dto := domain.UpdateIntegrationTypeDTO{
		Name:        &nuevoNombre,
		Description: &nuevaDesc,
	}

	repo.On("GetIntegrationTypeByID", mock.Anything, uint(7)).Return(tipoExistente, nil)
	// Verificar que el nuevo nombre no esté en uso
	repo.On("GetIntegrationTypeByName", mock.Anything, nuevoNombre).Return(nil, errors.New("no encontrado"))
	repo.On("UpdateIntegrationType", mock.Anything, uint(7), mock.AnythingOfType("*domain.IntegrationType")).Return(nil)
	// El use case invalida el caché de las integraciones que usan este tipo
	repo.On("ListIntegrationsByIntegrationTypeID", mock.Anything, uint(7)).Return([]*domain.Integration{}, nil).Maybe()

	// Act
	resultado, err := uc.UpdateIntegrationType(ctx, 7, dto)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, resultado)
	assert.Equal(t, nuevoNombre, resultado.Name)
	assert.Equal(t, nuevaDesc, resultado.Description)
}

func TestUpdateIntegrationType_NombreYaEnUsoPorOtro(t *testing.T) {
	// Arrange
	repo := new(mocks.RepositoryMock)
	s3 := new(mocks.S3Mock)
	logger := new(mocks.LoggerMock)
	cfg := new(mocks.ConfigMock)
	configurarLoggerPermisivoType(logger)

	uc := newTestIntegrationTypeUseCase(repo, s3, logger, cfg)
	ctx := context.Background()

	tipoExistente := &domain.IntegrationType{ID: 7, Name: "Factus", Code: "factus"}
	otroTipo := &domain.IntegrationType{ID: 1, Name: "Shopify"} // Tiene el nombre que queremos usar

	nuevoNombre := "Shopify"
	dto := domain.UpdateIntegrationTypeDTO{Name: &nuevoNombre}

	repo.On("GetIntegrationTypeByID", mock.Anything, uint(7)).Return(tipoExistente, nil)
	repo.On("GetIntegrationTypeByName", mock.Anything, nuevoNombre).Return(otroTipo, nil)
	// El use case puede intentar listar integraciones incluso en error si llega tan lejos
	repo.On("ListIntegrationsByIntegrationTypeID", mock.Anything, mock.Anything).Return([]*domain.Integration{}, nil).Maybe()

	// Act
	resultado, err := uc.UpdateIntegrationType(ctx, 7, dto)

	// Assert
	assert.ErrorIs(t, err, domain.ErrIntegrationTypeNameExists)
	assert.Nil(t, resultado)
}

func TestUpdateIntegrationType_NoEncontrado(t *testing.T) {
	// Arrange
	repo := new(mocks.RepositoryMock)
	s3 := new(mocks.S3Mock)
	logger := new(mocks.LoggerMock)
	cfg := new(mocks.ConfigMock)
	configurarLoggerPermisivoType(logger)

	uc := newTestIntegrationTypeUseCase(repo, s3, logger, cfg)
	ctx := context.Background()

	repo.On("GetIntegrationTypeByID", mock.Anything, uint(999)).Return(nil, errors.New("not found"))
	repo.On("ListIntegrationsByIntegrationTypeID", mock.Anything, mock.Anything).Return([]*domain.Integration{}, nil).Maybe()

	// Act
	resultado, err := uc.UpdateIntegrationType(ctx, 999, domain.UpdateIntegrationTypeDTO{})

	// Assert
	assert.ErrorIs(t, err, domain.ErrIntegrationTypeNotFound)
	assert.Nil(t, resultado)
}
