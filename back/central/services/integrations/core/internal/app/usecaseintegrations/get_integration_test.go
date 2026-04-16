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

// ============================================
// GetIntegrationByID
// ============================================

func TestGetIntegrationByID_CacheHit(t *testing.T) {
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
	cachedInteg := &domain.CachedIntegration{
		ID:                42,
		Name:              "Mi Shopify",
		Code:              "shopify_001",
		Category:          "ecommerce",
		IntegrationTypeID: 1,
		BusinessID:        &businessID,
		IsActive:          true,
	}

	cache.On("GetIntegration", mock.Anything, uint(42)).Return(cachedInteg, nil)
	repo.On("GetIntegrationTypeByID", mock.Anything, uint(1)).Return(nil, nil)

	// Act
	resultado, err := uc.GetIntegrationByID(ctx, 42)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, resultado)
	assert.Equal(t, uint(42), resultado.ID)
	assert.Equal(t, "Mi Shopify", resultado.Name)
	assert.Equal(t, "ecommerce", resultado.Category)
	// No debe consultar la BD si hay cache hit
	repo.AssertNotCalled(t, "GetIntegrationByID", mock.Anything, mock.Anything)
}

func TestGetIntegrationByID_CacheMiss_CargaDesdeBD(t *testing.T) {
	// Arrange
	repo := new(mocks.RepositoryMock)
	enc := new(mocks.EncryptionMock)
	cache := new(mocks.CacheMock)
	logger := new(mocks.LoggerMock)
	cfg := new(mocks.ConfigMock)
	configurarLoggerPermisivo(logger)

	uc := newTestUseCase(repo, enc, cache, logger, cfg)
	ctx := context.Background()

	integracion := &domain.Integration{
		ID:                10,
		Name:              "Factus",
		Code:              "factus_001",
		Category:          "invoicing",
		IntegrationTypeID: 7,
		IsActive:          true,
	}

	cache.On("GetIntegration", mock.Anything, uint(10)).Return(nil, errors.New("cache miss"))
	repo.On("GetIntegrationByID", mock.Anything, uint(10)).Return(integracion, nil)
	cache.On("SetIntegration", mock.Anything, mock.Anything).Return(nil)

	// Act
	resultado, err := uc.GetIntegrationByID(ctx, 10)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, resultado)
	assert.Equal(t, uint(10), resultado.ID)
	assert.Equal(t, "Factus", resultado.Name)
	repo.AssertCalled(t, "GetIntegrationByID", mock.Anything, uint(10))
	cache.AssertCalled(t, "SetIntegration", mock.Anything, mock.Anything)
}

func TestGetIntegrationByID_ErrorEnBD(t *testing.T) {
	// Arrange
	repo := new(mocks.RepositoryMock)
	enc := new(mocks.EncryptionMock)
	cache := new(mocks.CacheMock)
	logger := new(mocks.LoggerMock)
	cfg := new(mocks.ConfigMock)
	configurarLoggerPermisivo(logger)

	uc := newTestUseCase(repo, enc, cache, logger, cfg)
	ctx := context.Background()

	cache.On("GetIntegration", mock.Anything, uint(99)).Return(nil, errors.New("cache miss"))
	repo.On("GetIntegrationByID", mock.Anything, uint(99)).Return(nil, errors.New("registro no encontrado"))

	// Act
	resultado, err := uc.GetIntegrationByID(ctx, 99)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, resultado)
	assert.Contains(t, err.Error(), "registro no encontrado")
}

// ============================================
// GetPublicIntegrationByID
// ============================================

func TestGetPublicIntegrationByID_Exitoso(t *testing.T) {
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
	integracion := &domain.Integration{
		ID:                5,
		Name:              "Shopify Store",
		BusinessID:        &businessID,
		IntegrationTypeID: 1,
		StoreID:           "mi-tienda.myshopify.com",
	}

	repo.On("GetIntegrationByID", mock.Anything, uint(5)).Return(integracion, nil)

	// Act
	resultado, err := uc.GetPublicIntegrationByID(ctx, "5")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, resultado)
	assert.Equal(t, uint(5), resultado.ID)
	assert.Equal(t, "Shopify Store", resultado.Name)
	assert.Equal(t, 1, resultado.IntegrationType)
	assert.Equal(t, "mi-tienda.myshopify.com", resultado.StoreID)
}

func TestGetPublicIntegrationByID_IDInvalido(t *testing.T) {
	// Arrange
	repo := new(mocks.RepositoryMock)
	enc := new(mocks.EncryptionMock)
	cache := new(mocks.CacheMock)
	logger := new(mocks.LoggerMock)
	cfg := new(mocks.ConfigMock)
	configurarLoggerPermisivo(logger)

	uc := newTestUseCase(repo, enc, cache, logger, cfg)
	ctx := context.Background()

	// Act — ID no es un número
	resultado, err := uc.GetPublicIntegrationByID(ctx, "abc")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, resultado)
	assert.Contains(t, err.Error(), "invalid integration ID")
}

func TestGetPublicIntegrationByID_NoEncontrado(t *testing.T) {
	// Arrange
	repo := new(mocks.RepositoryMock)
	enc := new(mocks.EncryptionMock)
	cache := new(mocks.CacheMock)
	logger := new(mocks.LoggerMock)
	cfg := new(mocks.ConfigMock)
	configurarLoggerPermisivo(logger)

	uc := newTestUseCase(repo, enc, cache, logger, cfg)
	ctx := context.Background()

	repo.On("GetIntegrationByID", mock.Anything, uint(999)).Return(nil, errors.New("not found"))

	// Act
	resultado, err := uc.GetPublicIntegrationByID(ctx, "999")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, resultado)
}

// ============================================
// GetIntegrationByType
// ============================================

func TestGetIntegrationByType_ExitoConCredenciales(t *testing.T) {
	// Arrange
	repo := new(mocks.RepositoryMock)
	enc := new(mocks.EncryptionMock)
	cache := new(mocks.CacheMock)
	logger := new(mocks.LoggerMock)
	cfg := new(mocks.ConfigMock)
	configurarLoggerPermisivo(logger)

	uc := newTestUseCase(repo, enc, cache, logger, cfg)
	ctx := context.Background()

	businessID := uint(10)
	tipoIntegracion := &domain.IntegrationType{ID: 1, Code: "shopify"}

	// Credenciales encriptadas en formato {"encrypted": "<base64>"}
	import64 := `{"encrypted": "dGVzdA=="}`
	integracion := &domain.Integration{
		ID:                3,
		Name:              "Shopify Store",
		IntegrationTypeID: 1,
		Credentials:       []byte(import64),
	}

	credencialesDesc := map[string]interface{}{"api_key": "key_real"}

	repo.On("GetIntegrationTypeByCode", mock.Anything, "shopify").Return(tipoIntegracion, nil)
	repo.On("GetActiveIntegrationByIntegrationTypeID", mock.Anything, uint(1), &businessID).Return(integracion, nil)
	enc.On("DecryptCredentials", mock.Anything, mock.Anything).Return(credencialesDesc, nil)

	// Act
	resultado, err := uc.GetIntegrationByType(ctx, "shopify", &businessID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, resultado)
	assert.Equal(t, uint(3), resultado.ID)
	assert.Equal(t, "key_real", resultado.DecryptedCredentials["api_key"])
}

func TestGetIntegrationByType_TipoNoEncontrado(t *testing.T) {
	// Arrange
	repo := new(mocks.RepositoryMock)
	enc := new(mocks.EncryptionMock)
	cache := new(mocks.CacheMock)
	logger := new(mocks.LoggerMock)
	cfg := new(mocks.ConfigMock)
	configurarLoggerPermisivo(logger)

	uc := newTestUseCase(repo, enc, cache, logger, cfg)
	ctx := context.Background()

	repo.On("GetIntegrationTypeByCode", mock.Anything, "inexistente").Return(nil, errors.New("no encontrado"))

	// Act
	resultado, err := uc.GetIntegrationByType(ctx, "inexistente", nil)

	// Assert
	assert.ErrorIs(t, err, domain.ErrIntegrationTypeNotFound)
	assert.Nil(t, resultado)
}

// ============================================
// DecryptCredentialField
// ============================================

func TestDecryptCredentialField_CampoEncontrado(t *testing.T) {
	// Arrange
	repo := new(mocks.RepositoryMock)
	enc := new(mocks.EncryptionMock)
	cache := new(mocks.CacheMock)
	logger := new(mocks.LoggerMock)
	cfg := new(mocks.ConfigMock)
	configurarLoggerPermisivo(logger)

	uc := newTestUseCase(repo, enc, cache, logger, cfg)
	ctx := context.Background()

	import64 := `{"encrypted": "dGVzdA=="}`
	integracion := &domain.Integration{
		ID:          7,
		Credentials: []byte(import64),
	}

	credenciales := map[string]interface{}{"api_key": "valor_secreto"}
	cachedCreds := &domain.CachedCredentials{IntegrationID: 7, Credentials: credenciales}

	cache.On("GetIntegration", mock.Anything, uint(7)).Return(nil, errors.New("cache miss"))
	repo.On("GetIntegrationByID", mock.Anything, uint(7)).Return(integracion, nil)
	cache.On("SetIntegration", mock.Anything, mock.Anything).Return(nil)
	cache.On("GetCredentials", mock.Anything, uint(7)).Return(cachedCreds, nil)

	// Act
	valor, err := uc.DecryptCredentialField(ctx, "7", "api_key")

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "valor_secreto", valor)
}

func TestDecryptCredentialField_CampoNoExiste(t *testing.T) {
	// Arrange
	repo := new(mocks.RepositoryMock)
	enc := new(mocks.EncryptionMock)
	cache := new(mocks.CacheMock)
	logger := new(mocks.LoggerMock)
	cfg := new(mocks.ConfigMock)
	configurarLoggerPermisivo(logger)

	uc := newTestUseCase(repo, enc, cache, logger, cfg)
	ctx := context.Background()

	import64 := `{"encrypted": "dGVzdA=="}`
	integracion := &domain.Integration{
		ID:          7,
		Credentials: []byte(import64),
	}

	credenciales := map[string]interface{}{"api_key": "valor"}
	cachedCreds := &domain.CachedCredentials{IntegrationID: 7, Credentials: credenciales}

	cache.On("GetIntegration", mock.Anything, uint(7)).Return(nil, errors.New("miss"))
	repo.On("GetIntegrationByID", mock.Anything, uint(7)).Return(integracion, nil)
	cache.On("SetIntegration", mock.Anything, mock.Anything).Return(nil)
	cache.On("GetCredentials", mock.Anything, uint(7)).Return(cachedCreds, nil)

	// Act
	valor, err := uc.DecryptCredentialField(ctx, "7", "campo_inexistente")

	// Assert
	assert.Error(t, err)
	assert.Empty(t, valor)
	assert.Contains(t, err.Error(), "campo_inexistente")
}

func TestDecryptCredentialField_IDInvalido(t *testing.T) {
	// Arrange
	repo := new(mocks.RepositoryMock)
	enc := new(mocks.EncryptionMock)
	cache := new(mocks.CacheMock)
	logger := new(mocks.LoggerMock)
	cfg := new(mocks.ConfigMock)
	configurarLoggerPermisivo(logger)

	uc := newTestUseCase(repo, enc, cache, logger, cfg)
	ctx := context.Background()

	// Act
	valor, err := uc.DecryptCredentialField(ctx, "no_es_numero", "api_key")

	// Assert
	assert.Error(t, err)
	assert.Empty(t, valor)
	assert.Contains(t, err.Error(), "invalid integration ID")
}
