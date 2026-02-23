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
// TestConnectionRaw
// ============================================

func TestTestConnectionRaw_ConProvider_Exitoso(t *testing.T) {
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
	config := map[string]interface{}{"shop_domain": "mi-tienda.myshopify.com"}
	credentials := map[string]interface{}{"api_key": "key_123"}

	provider.On("TestConnection", mock.Anything, config, credentials).Return(nil)

	// Act
	err := uc.TestConnectionRaw(ctx, "shopify", config, credentials)

	// Assert
	assert.NoError(t, err)
	provider.AssertCalled(t, "TestConnection", mock.Anything, config, credentials)
}

func TestTestConnectionRaw_SinProvider_ValidaAccessToken(t *testing.T) {
	// Arrange — sin provider registrado, fallback a validateBasicCredentials
	repo := new(mocks.RepositoryMock)
	enc := new(mocks.EncryptionMock)
	cache := new(mocks.CacheMock)
	logger := new(mocks.LoggerMock)
	cfg := new(mocks.ConfigMock)
	configurarLoggerPermisivo(logger)

	uc := newTestUseCase(repo, enc, cache, logger, cfg)
	ctx := context.Background()

	credentials := map[string]interface{}{"access_token": "tok_valido"}

	// Act
	err := uc.TestConnectionRaw(ctx, "shopify", nil, credentials)

	// Assert
	assert.NoError(t, err)
}

func TestTestConnectionRaw_SinProvider_SinAccessToken(t *testing.T) {
	// Arrange
	repo := new(mocks.RepositoryMock)
	enc := new(mocks.EncryptionMock)
	cache := new(mocks.CacheMock)
	logger := new(mocks.LoggerMock)
	cfg := new(mocks.ConfigMock)
	configurarLoggerPermisivo(logger)

	uc := newTestUseCase(repo, enc, cache, logger, cfg)
	ctx := context.Background()

	credentials := map[string]interface{}{} // Sin access_token

	// Act
	err := uc.TestConnectionRaw(ctx, "shopify", nil, credentials)

	// Assert
	assert.ErrorIs(t, err, domain.ErrIntegrationAccessTokenNotFound)
}

func TestTestConnectionRaw_ConProvider_TestFalla(t *testing.T) {
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
	config := map[string]interface{}{}
	credentials := map[string]interface{}{"api_key": "invalida"}

	provider.On("TestConnection", mock.Anything, config, credentials).Return(errors.New("credenciales inválidas"))

	// Act
	err := uc.TestConnectionRaw(ctx, "shopify", config, credentials)

	// Assert
	assert.ErrorIs(t, err, domain.ErrIntegrationTestFailed)
}

// ============================================
// TestConnectionFromConfig
// ============================================

func TestTestConnectionFromConfig_TipoComoString(t *testing.T) {
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
	config := map[string]interface{}{"integration_type": "shopify"}
	credentials := map[string]interface{}{"api_key": "key"}

	provider.On("TestConnection", mock.Anything, config, credentials).Return(nil)

	// Act
	err := uc.TestConnectionFromConfig(ctx, config, credentials)

	// Assert
	assert.NoError(t, err)
}

func TestTestConnectionFromConfig_TipoComoInt(t *testing.T) {
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
	config := map[string]interface{}{"integration_type": 1} // int
	credentials := map[string]interface{}{}

	provider.On("TestConnection", mock.Anything, config, credentials).Return(nil)

	// Act
	err := uc.TestConnectionFromConfig(ctx, config, credentials)

	// Assert
	assert.NoError(t, err)
}

func TestTestConnectionFromConfig_SinTipoEnConfig(t *testing.T) {
	// Arrange
	repo := new(mocks.RepositoryMock)
	enc := new(mocks.EncryptionMock)
	cache := new(mocks.CacheMock)
	logger := new(mocks.LoggerMock)
	cfg := new(mocks.ConfigMock)
	configurarLoggerPermisivo(logger)

	uc := newTestUseCase(repo, enc, cache, logger, cfg)
	ctx := context.Background()

	config := map[string]interface{}{} // Sin integration_type

	// Act
	err := uc.TestConnectionFromConfig(ctx, config, nil)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "integration_type is required")
}

func TestTestConnectionFromConfig_ProviderNoRegistrado(t *testing.T) {
	// Arrange
	repo := new(mocks.RepositoryMock)
	enc := new(mocks.EncryptionMock)
	cache := new(mocks.CacheMock)
	logger := new(mocks.LoggerMock)
	cfg := new(mocks.ConfigMock)
	configurarLoggerPermisivo(logger)

	uc := newTestUseCase(repo, enc, cache, logger, cfg)
	ctx := context.Background()

	// Tipo 99 no registrado
	config := map[string]interface{}{"integration_type": 99}

	// Act
	err := uc.TestConnectionFromConfig(ctx, config, nil)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no registrada para tipo")
}

// ============================================
// SyncOrdersByIntegrationID
// ============================================

func TestSyncOrdersByIntegrationID_ProviderNoRegistrado(t *testing.T) {
	// Arrange
	repo := new(mocks.RepositoryMock)
	enc := new(mocks.EncryptionMock)
	cache := new(mocks.CacheMock)
	logger := new(mocks.LoggerMock)
	cfg := new(mocks.ConfigMock)
	configurarLoggerPermisivo(logger)

	uc := newTestUseCase(repo, enc, cache, logger, cfg)
	ctx := context.Background()

	// La integración tiene tipo 99 sin provider registrado
	integracion := &domain.Integration{
		ID:                8,
		IntegrationTypeID: 99,
	}
	repo.On("GetIntegrationByID", mock.Anything, uint(8)).Return(integracion, nil)

	// Act
	err := uc.SyncOrdersByIntegrationID(ctx, "8")

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "integración no registrada")
}

func TestSyncOrdersByIntegrationIDWithParams_FallbackCuandoNoSoportado(t *testing.T) {
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
	integracion := &domain.Integration{
		ID:                11,
		IntegrationTypeID: domain.IntegrationTypeShopify,
	}
	repo.On("GetIntegrationByID", mock.Anything, uint(11)).Return(integracion, nil)
	// SyncOrdersByIntegrationIDWithParams retorna ErrNotSupported → debe hacer fallback
	provider.On("SyncOrdersByIntegrationIDWithParams", mock.Anything, "11", mock.Anything).Return(domain.ErrNotSupported)
	provider.On("SyncOrdersByIntegrationID", mock.Anything, "11").Return(nil)

	// Act
	err := uc.SyncOrdersByIntegrationIDWithParams(ctx, "11", nil)

	// Assert
	assert.NoError(t, err)
	provider.AssertCalled(t, "SyncOrdersByIntegrationID", mock.Anything, "11")
}

// ============================================
// ProviderRegistry
// ============================================

func TestProviderRegistry_RegisterYGet(t *testing.T) {
	// Arrange
	repo := new(mocks.RepositoryMock)
	enc := new(mocks.EncryptionMock)
	cache := new(mocks.CacheMock)
	logger := new(mocks.LoggerMock)
	cfg := new(mocks.ConfigMock)
	provider := new(mocks.ProviderMock)
	configurarLoggerPermisivo(logger)

	uc := newTestUseCase(repo, enc, cache, logger, cfg)

	// Act — registrar y recuperar
	uc.RegisterProvider(domain.IntegrationTypeShopify, provider)
	recuperado, existe := uc.GetProvider(domain.IntegrationTypeShopify)

	// Assert
	assert.True(t, existe)
	assert.Equal(t, provider, recuperado)
}

func TestProviderRegistry_NoRegistrado_RetornaFalse(t *testing.T) {
	// Arrange
	repo := new(mocks.RepositoryMock)
	enc := new(mocks.EncryptionMock)
	cache := new(mocks.CacheMock)
	logger := new(mocks.LoggerMock)
	cfg := new(mocks.ConfigMock)
	configurarLoggerPermisivo(logger)

	uc := newTestUseCase(repo, enc, cache, logger, cfg)

	// Act
	recuperado, existe := uc.GetProvider(999)

	// Assert
	assert.False(t, existe)
	assert.Nil(t, recuperado)
}

func TestProviderRegistry_IgnoraNilProvider(t *testing.T) {
	// Arrange
	registro := newProviderRegistry()

	// Act — registrar nil no debe causar panic ni guardar nada
	registro.Register(1, nil)
	proveedor, existe := registro.Get(1)

	// Assert
	assert.False(t, existe)
	assert.Nil(t, proveedor)
}

func TestProviderRegistry_IgnoraTipoCero(t *testing.T) {
	// Arrange
	registro := newProviderRegistry()
	provider := new(mocks.ProviderMock)

	// Act — tipo 0 no debe registrarse
	registro.Register(0, provider)
	proveedor, existe := registro.Get(0)

	// Assert
	assert.False(t, existe)
	assert.Nil(t, proveedor)
}
