package usecaseintegrations

import (
	"context"
	"errors"
	"testing"

	"github.com/secamc93/probability/back/central/services/integrations/core/internal/domain"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// newTestUseCase construye un IntegrationUseCase con mocks para todos los tests.
func newTestUseCase(
	repo *mocks.RepositoryMock,
	enc *mocks.EncryptionMock,
	cache *mocks.CacheMock,
	logger *mocks.LoggerMock,
	cfg *mocks.ConfigMock,
) *IntegrationUseCase {
	return New(repo, enc, cache, logger, cfg)
}

// configurarLoggerPermisivo configura el logger para que acepte cualquier llamada
// sin fallar el test — el logger es llamado internamente por el use case.
func configurarLoggerPermisivo(logger *mocks.LoggerMock) {
	logger.On("Info", mock.Anything).Maybe()
	logger.On("Error", mock.Anything).Maybe()
	logger.On("Warn", mock.Anything).Maybe()
	logger.On("Debug", mock.Anything).Maybe()
}

func TestCreateIntegration_ExitoSinProvider(t *testing.T) {
	// Arrange
	repo := new(mocks.RepositoryMock)
	enc := new(mocks.EncryptionMock)
	cache := new(mocks.CacheMock)
	logger := new(mocks.LoggerMock)
	cfg := new(mocks.ConfigMock)
	configurarLoggerPermisivo(logger)

	uc := newTestUseCase(repo, enc, cache, logger, cfg)

	ctx := context.Background()
	businessID := uint(1)

	categoria := &domain.IntegrationCategory{ID: 10, Code: "external"}
	tipoIntegracion := &domain.IntegrationType{
		ID:         5,
		Code:       "factus",
		CategoryID: 10,
		Category:   categoria,
	}

	dto := domain.CreateIntegrationDTO{
		Name:              "Mi integración Factus",
		Code:              "factus_001",
		IntegrationTypeID: 5,
		BusinessID:        &businessID,
		Credentials:       map[string]interface{}{"access_token": "tok_123"},
	}

	repo.On("ExistsIntegrationByCode", mock.Anything, dto.Code, dto.BusinessID).Return(false, nil)
	repo.On("GetIntegrationTypeByID", mock.Anything, uint(5)).Return(tipoIntegracion, nil)
	repo.On("CreateIntegration", mock.Anything, mock.AnythingOfType("*domain.Integration")).Return(nil)
	cache.On("SetIntegration", mock.Anything, mock.Anything).Return(nil)
	cache.On("SetCredentials", mock.Anything, mock.Anything).Return(nil)

	// Act
	resultado, err := uc.CreateIntegration(ctx, dto)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, resultado)
	assert.Equal(t, dto.Name, resultado.Name)
	assert.Equal(t, dto.Code, resultado.Code)
	assert.Equal(t, "external", resultado.Category)
	repo.AssertExpectations(t)
	cache.AssertCalled(t, "SetIntegration", mock.Anything, mock.Anything)
	cache.AssertCalled(t, "SetCredentials", mock.Anything, mock.Anything)
}

func TestCreateIntegration_ConProvider_TestConnectionExitoso(t *testing.T) {
	// Arrange
	repo := new(mocks.RepositoryMock)
	enc := new(mocks.EncryptionMock)
	cache := new(mocks.CacheMock)
	logger := new(mocks.LoggerMock)
	cfg := new(mocks.ConfigMock)
	provider := new(mocks.ProviderMock)
	configurarLoggerPermisivo(logger)

	uc := newTestUseCase(repo, enc, cache, logger, cfg)
	uc.RegisterProvider(domain.IntegrationTypeShopify, provider)

	ctx := context.Background()
	businessID := uint(2)

	categoria := &domain.IntegrationCategory{ID: 20, Code: "ecommerce"}
	tipoIntegracion := &domain.IntegrationType{
		ID:         domain.IntegrationTypeShopify,
		Code:       "shopify",
		CategoryID: 20,
		Category:   categoria,
	}

	dto := domain.CreateIntegrationDTO{
		Name:              "Mi Shopify",
		Code:              "shopify_tienda",
		IntegrationTypeID: domain.IntegrationTypeShopify,
		BusinessID:        &businessID,
		Credentials:       map[string]interface{}{"api_key": "key_abc"},
	}

	repo.On("ExistsIntegrationByCode", mock.Anything, dto.Code, dto.BusinessID).Return(false, nil)
	repo.On("GetIntegrationTypeByID", mock.Anything, uint(domain.IntegrationTypeShopify)).Return(tipoIntegracion, nil)
	provider.On("TestConnection", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	repo.On("CreateIntegration", mock.Anything, mock.AnythingOfType("*domain.Integration")).Return(nil)
	cache.On("SetIntegration", mock.Anything, mock.Anything).Return(nil)
	cache.On("SetCredentials", mock.Anything, mock.Anything).Return(nil)
	// CreateWebhookForIntegration se llama en goroutine — usar Maybe para evitar race
	cache.On("GetIntegration", mock.Anything, mock.Anything).Maybe().Return(nil, errors.New("cache miss"))
	repo.On("GetIntegrationByID", mock.Anything, mock.Anything).Maybe().Return(nil, errors.New("not called sync"))

	// Act
	resultado, err := uc.CreateIntegration(ctx, dto)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, resultado)
	provider.AssertCalled(t, "TestConnection", mock.Anything, mock.Anything, mock.Anything)
}

func TestCreateIntegration_ErrorNombreVacio(t *testing.T) {
	// Arrange
	repo := new(mocks.RepositoryMock)
	enc := new(mocks.EncryptionMock)
	cache := new(mocks.CacheMock)
	logger := new(mocks.LoggerMock)
	cfg := new(mocks.ConfigMock)
	configurarLoggerPermisivo(logger)

	uc := newTestUseCase(repo, enc, cache, logger, cfg)
	ctx := context.Background()

	dto := domain.CreateIntegrationDTO{
		Name:              "", // Nombre vacío — debe fallar
		Code:              "codigo",
		IntegrationTypeID: 1,
	}

	// Act
	resultado, err := uc.CreateIntegration(ctx, dto)

	// Assert
	assert.ErrorIs(t, err, domain.ErrIntegrationNameRequired)
	assert.Nil(t, resultado)
	repo.AssertNotCalled(t, "ExistsIntegrationByCode", mock.Anything, mock.Anything, mock.Anything)
}

func TestCreateIntegration_ErrorCodigoVacio(t *testing.T) {
	// Arrange
	repo := new(mocks.RepositoryMock)
	enc := new(mocks.EncryptionMock)
	cache := new(mocks.CacheMock)
	logger := new(mocks.LoggerMock)
	cfg := new(mocks.ConfigMock)
	configurarLoggerPermisivo(logger)

	uc := newTestUseCase(repo, enc, cache, logger, cfg)
	ctx := context.Background()

	dto := domain.CreateIntegrationDTO{
		Name:              "Mi integración",
		Code:              "", // Código vacío — debe fallar
		IntegrationTypeID: 1,
	}

	// Act
	resultado, err := uc.CreateIntegration(ctx, dto)

	// Assert
	assert.ErrorIs(t, err, domain.ErrIntegrationCodeRequired)
	assert.Nil(t, resultado)
}

func TestCreateIntegration_ErrorTipoIntegracionCero(t *testing.T) {
	// Arrange
	repo := new(mocks.RepositoryMock)
	enc := new(mocks.EncryptionMock)
	cache := new(mocks.CacheMock)
	logger := new(mocks.LoggerMock)
	cfg := new(mocks.ConfigMock)
	configurarLoggerPermisivo(logger)

	uc := newTestUseCase(repo, enc, cache, logger, cfg)
	ctx := context.Background()

	dto := domain.CreateIntegrationDTO{
		Name:              "Mi integración",
		Code:              "mi_codigo",
		IntegrationTypeID: 0, // Tipo inválido — debe fallar
	}

	// Act
	resultado, err := uc.CreateIntegration(ctx, dto)

	// Assert
	assert.ErrorIs(t, err, domain.ErrIntegrationTypeRequired)
	assert.Nil(t, resultado)
}

func TestCreateIntegration_ErrorCodigoDuplicado(t *testing.T) {
	// Arrange
	repo := new(mocks.RepositoryMock)
	enc := new(mocks.EncryptionMock)
	cache := new(mocks.CacheMock)
	logger := new(mocks.LoggerMock)
	cfg := new(mocks.ConfigMock)
	configurarLoggerPermisivo(logger)

	uc := newTestUseCase(repo, enc, cache, logger, cfg)
	ctx := context.Background()

	dto := domain.CreateIntegrationDTO{
		Name:              "Mi integración",
		Code:              "codigo_existente",
		IntegrationTypeID: 1,
	}

	repo.On("ExistsIntegrationByCode", mock.Anything, dto.Code, dto.BusinessID).Return(true, nil)

	// Act
	resultado, err := uc.CreateIntegration(ctx, dto)

	// Assert
	assert.ErrorIs(t, err, domain.ErrIntegrationCodeExists)
	assert.Nil(t, resultado)
	repo.AssertExpectations(t)
}

func TestCreateIntegration_ErrorRepositorioAlVerificarCodigo(t *testing.T) {
	// Arrange
	repo := new(mocks.RepositoryMock)
	enc := new(mocks.EncryptionMock)
	cache := new(mocks.CacheMock)
	logger := new(mocks.LoggerMock)
	cfg := new(mocks.ConfigMock)
	configurarLoggerPermisivo(logger)

	uc := newTestUseCase(repo, enc, cache, logger, cfg)
	ctx := context.Background()

	dto := domain.CreateIntegrationDTO{
		Name:              "Mi integración",
		Code:              "codigo",
		IntegrationTypeID: 1,
	}

	errDB := errors.New("connection timeout")
	repo.On("ExistsIntegrationByCode", mock.Anything, dto.Code, dto.BusinessID).Return(false, errDB)

	// Act
	resultado, err := uc.CreateIntegration(ctx, dto)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, resultado)
	assert.Contains(t, err.Error(), "error al verificar código")
}

func TestCreateIntegration_ErrorAlObtenerTipoIntegracion(t *testing.T) {
	// Arrange
	repo := new(mocks.RepositoryMock)
	enc := new(mocks.EncryptionMock)
	cache := new(mocks.CacheMock)
	logger := new(mocks.LoggerMock)
	cfg := new(mocks.ConfigMock)
	configurarLoggerPermisivo(logger)

	uc := newTestUseCase(repo, enc, cache, logger, cfg)
	ctx := context.Background()

	dto := domain.CreateIntegrationDTO{
		Name:              "Mi integración",
		Code:              "codigo",
		IntegrationTypeID: 99,
	}

	repo.On("ExistsIntegrationByCode", mock.Anything, dto.Code, dto.BusinessID).Return(false, nil)
	repo.On("GetIntegrationTypeByID", mock.Anything, uint(99)).Return(nil, errors.New("tipo no encontrado"))

	// Act
	resultado, err := uc.CreateIntegration(ctx, dto)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, resultado)
	assert.Contains(t, err.Error(), "error al obtener tipo de integración")
}

func TestCreateIntegration_ErrorCategoriaSinCodigo(t *testing.T) {
	// Arrange
	repo := new(mocks.RepositoryMock)
	enc := new(mocks.EncryptionMock)
	cache := new(mocks.CacheMock)
	logger := new(mocks.LoggerMock)
	cfg := new(mocks.ConfigMock)
	configurarLoggerPermisivo(logger)

	uc := newTestUseCase(repo, enc, cache, logger, cfg)
	ctx := context.Background()

	// Tipo sin categoría cargada ni CategoryID — debe retornar ErrIntegrationCategoryInvalid
	tipoSinCategoria := &domain.IntegrationType{
		ID:         7,
		Code:       "factus",
		CategoryID: 0,
		Category:   nil,
	}

	dto := domain.CreateIntegrationDTO{
		Name:              "Mi integración",
		Code:              "codigo",
		IntegrationTypeID: 7,
	}

	repo.On("ExistsIntegrationByCode", mock.Anything, dto.Code, dto.BusinessID).Return(false, nil)
	repo.On("GetIntegrationTypeByID", mock.Anything, uint(7)).Return(tipoSinCategoria, nil)

	// Act
	resultado, err := uc.CreateIntegration(ctx, dto)

	// Assert
	assert.ErrorIs(t, err, domain.ErrIntegrationCategoryInvalid)
	assert.Nil(t, resultado)
}

func TestCreateIntegration_ErrorTestConnectionFalla(t *testing.T) {
	// Arrange
	repo := new(mocks.RepositoryMock)
	enc := new(mocks.EncryptionMock)
	cache := new(mocks.CacheMock)
	logger := new(mocks.LoggerMock)
	cfg := new(mocks.ConfigMock)
	provider := new(mocks.ProviderMock)
	configurarLoggerPermisivo(logger)

	uc := newTestUseCase(repo, enc, cache, logger, cfg)
	uc.RegisterProvider(domain.IntegrationTypeShopify, provider)

	ctx := context.Background()
	categoria := &domain.IntegrationCategory{ID: 20, Code: "ecommerce"}
	tipoIntegracion := &domain.IntegrationType{
		ID:         domain.IntegrationTypeShopify,
		Code:       "shopify",
		CategoryID: 20,
		Category:   categoria,
	}

	dto := domain.CreateIntegrationDTO{
		Name:              "Mi Shopify",
		Code:              "shopify_tienda",
		IntegrationTypeID: domain.IntegrationTypeShopify,
		Credentials:       map[string]interface{}{"api_key": "key_invalida"},
	}

	repo.On("ExistsIntegrationByCode", mock.Anything, dto.Code, dto.BusinessID).Return(false, nil)
	repo.On("GetIntegrationTypeByID", mock.Anything, uint(domain.IntegrationTypeShopify)).Return(tipoIntegracion, nil)
	provider.On("TestConnection", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("credenciales inválidas"))

	// Act
	resultado, err := uc.CreateIntegration(ctx, dto)

	// Assert
	assert.ErrorIs(t, err, domain.ErrIntegrationTestFailed)
	assert.Nil(t, resultado)
	repo.AssertNotCalled(t, "CreateIntegration", mock.Anything, mock.Anything)
}

func TestCreateIntegration_ErrorAlGuardarEnRepositorio(t *testing.T) {
	// Arrange
	repo := new(mocks.RepositoryMock)
	enc := new(mocks.EncryptionMock)
	cache := new(mocks.CacheMock)
	logger := new(mocks.LoggerMock)
	cfg := new(mocks.ConfigMock)
	configurarLoggerPermisivo(logger)

	uc := newTestUseCase(repo, enc, cache, logger, cfg)
	ctx := context.Background()

	categoria := &domain.IntegrationCategory{ID: 10, Code: "external"}
	tipoIntegracion := &domain.IntegrationType{
		ID:         7,
		Code:       "factus",
		CategoryID: 10,
		Category:   categoria,
	}

	dto := domain.CreateIntegrationDTO{
		Name:              "Mi integración",
		Code:              "factus_001",
		IntegrationTypeID: 7,
		Credentials:       map[string]interface{}{"access_token": "tok"},
	}

	repo.On("ExistsIntegrationByCode", mock.Anything, dto.Code, dto.BusinessID).Return(false, nil)
	repo.On("GetIntegrationTypeByID", mock.Anything, uint(7)).Return(tipoIntegracion, nil)
	repo.On("CreateIntegration", mock.Anything, mock.AnythingOfType("*domain.Integration")).Return(errors.New("constraint violation"))

	// Act
	resultado, err := uc.CreateIntegration(ctx, dto)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, resultado)
	assert.Contains(t, err.Error(), "error al crear integración")
}

func TestCreateIntegration_CategoriaObtenidaPorID(t *testing.T) {
	// Arrange: tipo de integración sin Category precargada, pero con CategoryID válido
	repo := new(mocks.RepositoryMock)
	enc := new(mocks.EncryptionMock)
	cache := new(mocks.CacheMock)
	logger := new(mocks.LoggerMock)
	cfg := new(mocks.ConfigMock)
	configurarLoggerPermisivo(logger)

	uc := newTestUseCase(repo, enc, cache, logger, cfg)
	ctx := context.Background()

	tipoIntegracion := &domain.IntegrationType{
		ID:         8,
		Code:       "siigo",
		CategoryID: 5, // Tiene CategoryID pero no Category precargada
		Category:   nil,
	}
	categoria := &domain.IntegrationCategory{ID: 5, Code: "invoicing"}

	dto := domain.CreateIntegrationDTO{
		Name:              "Siigo",
		Code:              "siigo_001",
		IntegrationTypeID: 8,
		Credentials:       map[string]interface{}{"access_token": "tok"},
	}

	repo.On("ExistsIntegrationByCode", mock.Anything, dto.Code, dto.BusinessID).Return(false, nil)
	repo.On("GetIntegrationTypeByID", mock.Anything, uint(8)).Return(tipoIntegracion, nil)
	repo.On("GetIntegrationCategoryByID", mock.Anything, uint(5)).Return(categoria, nil)
	repo.On("CreateIntegration", mock.Anything, mock.AnythingOfType("*domain.Integration")).Return(nil)
	cache.On("SetIntegration", mock.Anything, mock.Anything).Return(nil)
	cache.On("SetCredentials", mock.Anything, mock.Anything).Return(nil)

	// Act
	resultado, err := uc.CreateIntegration(ctx, dto)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, resultado)
	assert.Equal(t, "invoicing", resultado.Category)
	repo.AssertCalled(t, "GetIntegrationCategoryByID", mock.Anything, uint(5))
}
