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
// SyncOrdersByIntegrationID
// ============================================

func TestSyncOrdersByIntegrationID_Exitoso(t *testing.T) {
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

	businessID := uint(1)
	integracion := &domain.Integration{
		ID:                5,
		IntegrationTypeID: uint(domain.IntegrationTypeShopify),
		BusinessID:        &businessID,
	}

	repo.On("GetIntegrationByID", mock.Anything, uint(5)).Return(integracion, nil)
	provider.On("SyncOrdersByIntegrationID", mock.Anything, "5").Return(nil)

	// Act
	err := uc.SyncOrdersByIntegrationID(ctx, "5")

	// Assert
	assert.NoError(t, err)
	provider.AssertCalled(t, "SyncOrdersByIntegrationID", mock.Anything, "5")
}

func TestSyncOrdersByIntegrationID_IntegracionNoEncontrada(t *testing.T) {
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
	err := uc.SyncOrdersByIntegrationID(ctx, "999")

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error al obtener integración")
}

func TestSyncOrdersByIntegrationID_TipoSinProviderRegistrado(t *testing.T) {
	// Arrange
	repo := new(mocks.RepositoryMock)
	enc := new(mocks.EncryptionMock)
	cache := new(mocks.CacheMock)
	logger := new(mocks.LoggerMock)
	cfg := new(mocks.ConfigMock)
	configurarLoggerPermisivo(logger)

	uc := newTestUseCase(repo, enc, cache, logger, cfg)
	// No se registra ningún provider para el tipo 77
	ctx := context.Background()

	businessID := uint(1)
	integracion := &domain.Integration{
		ID:                7,
		IntegrationTypeID: 77,
		BusinessID:        &businessID,
	}

	repo.On("GetIntegrationByID", mock.Anything, uint(7)).Return(integracion, nil)

	// Act
	err := uc.SyncOrdersByIntegrationID(ctx, "7")

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "integración no registrada para tipo 77")
}

// ============================================
// SyncOrdersByIntegrationIDWithParams
// ============================================

func TestSyncOrdersByIntegrationIDWithParams_Exitoso(t *testing.T) {
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

	businessID := uint(1)
	integracion := &domain.Integration{
		ID:                5,
		IntegrationTypeID: uint(domain.IntegrationTypeShopify),
		BusinessID:        &businessID,
	}

	params := map[string]interface{}{"status": "paid"}

	repo.On("GetIntegrationByID", mock.Anything, uint(5)).Return(integracion, nil)
	provider.On("SyncOrdersByIntegrationIDWithParams", mock.Anything, "5", params).Return(nil)

	// Act
	err := uc.SyncOrdersByIntegrationIDWithParams(ctx, "5", params)

	// Assert
	assert.NoError(t, err)
	provider.AssertCalled(t, "SyncOrdersByIntegrationIDWithParams", mock.Anything, "5", params)
}

func TestSyncOrdersByIntegrationIDWithParams_FallbackSiNoSoportado(t *testing.T) {
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

	businessID := uint(1)
	integracion := &domain.Integration{
		ID:                5,
		IntegrationTypeID: uint(domain.IntegrationTypeShopify),
		BusinessID:        &businessID,
	}

	params := map[string]interface{}{"filter": "test"}

	repo.On("GetIntegrationByID", mock.Anything, uint(5)).Return(integracion, nil)
	// El provider no soporta params — retorna ErrNotSupported
	provider.On("SyncOrdersByIntegrationIDWithParams", mock.Anything, "5", params).Return(domain.ErrNotSupported)
	// Fallback al método sin params
	provider.On("SyncOrdersByIntegrationID", mock.Anything, "5").Return(nil)

	// Act
	err := uc.SyncOrdersByIntegrationIDWithParams(ctx, "5", params)

	// Assert
	assert.NoError(t, err)
	provider.AssertCalled(t, "SyncOrdersByIntegrationID", mock.Anything, "5")
}

// ============================================
// SyncOrdersByBusiness
// ============================================

func TestSyncOrdersByBusiness_SincronizaIntegracionesActivas(t *testing.T) {
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

	businessID := uint(10)
	tipoShopify := &domain.IntegrationType{ID: uint(domain.IntegrationTypeShopify), Code: "shopify"}
	integraciones := []*domain.Integration{
		{
			ID:                1,
			IntegrationTypeID: uint(domain.IntegrationTypeShopify),
			BusinessID:        &businessID,
			IsActive:          true,
			IntegrationType:   tipoShopify,
		},
	}

	repo.On("ListIntegrations", mock.Anything, mock.MatchedBy(func(f domain.IntegrationFilters) bool {
		return f.BusinessID != nil && *f.BusinessID == 10 && f.IsActive != nil && *f.IsActive
	})).Return(integraciones, int64(1), nil)

	provider.On("SyncOrdersByIntegrationID", mock.Anything, "1").Return(nil)

	// Act
	err := uc.SyncOrdersByBusiness(ctx, 10)

	// Assert
	assert.NoError(t, err)
}

func TestSyncOrdersByBusiness_ErrorDeRepositorio(t *testing.T) {
	// Arrange
	repo := new(mocks.RepositoryMock)
	enc := new(mocks.EncryptionMock)
	cache := new(mocks.CacheMock)
	logger := new(mocks.LoggerMock)
	cfg := new(mocks.ConfigMock)
	configurarLoggerPermisivo(logger)

	uc := newTestUseCase(repo, enc, cache, logger, cfg)
	ctx := context.Background()

	repo.On("ListIntegrations", mock.Anything, mock.Anything).Return(nil, int64(0), errors.New("db error"))

	// Act
	err := uc.SyncOrdersByBusiness(ctx, 5)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error al obtener integraciones")
}

func TestSyncOrdersByBusiness_OmiteIntegracionesSinTipo(t *testing.T) {
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
	integraciones := []*domain.Integration{
		{
			ID:              99,
			BusinessID:      &businessID,
			IsActive:        true,
			IntegrationType: nil, // Sin tipo
		},
	}

	repo.On("ListIntegrations", mock.Anything, mock.Anything).Return(integraciones, int64(1), nil)

	// Act
	err := uc.SyncOrdersByBusiness(ctx, 1)

	// Assert — no debe fallar aunque la integración no tenga tipo
	assert.NoError(t, err)
}

// ============================================
// GetWebhookURL
// ============================================

func TestGetWebhookURL_Exitoso(t *testing.T) {
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

	businessID := uint(1)
	integracion := &domain.Integration{
		ID:                5,
		IntegrationTypeID: uint(domain.IntegrationTypeShopify),
		BusinessID:        &businessID,
	}

	expectedWebhookInfo := &domain.WebhookInfo{
		URL:    "https://api.example.com/webhook",
		Method: "POST",
	}

	cfg.On("Get", "WEBHOOK_BASE_URL").Return("https://api.example.com")
	repo.On("GetIntegrationByID", mock.Anything, uint(5)).Return(integracion, nil)
	provider.On("GetWebhookURL", mock.Anything, "https://api.example.com", uint(5)).Return(expectedWebhookInfo, nil)

	// Act
	result, err := uc.GetWebhookURL(ctx, 5)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "https://api.example.com/webhook", result.URL)
}

func TestGetWebhookURL_FaltaBaseURL(t *testing.T) {
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

	businessID := uint(1)
	integracion := &domain.Integration{
		ID:                5,
		IntegrationTypeID: uint(domain.IntegrationTypeShopify),
		BusinessID:        &businessID,
	}

	// Ninguna URL configurada
	cfg.On("Get", "WEBHOOK_BASE_URL").Return("")
	cfg.On("Get", "URL_BASE_SWAGGER").Return("")
	repo.On("GetIntegrationByID", mock.Anything, uint(5)).Return(integracion, nil)

	// Act
	result, err := uc.GetWebhookURL(ctx, 5)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "WEBHOOK_BASE_URL")
}

// ============================================
// ListWebhooks
// ============================================

func TestListWebhooks_Exitoso(t *testing.T) {
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

	businessID := uint(1)
	integracion := &domain.Integration{
		ID:                5,
		IntegrationTypeID: uint(domain.IntegrationTypeShopify),
		BusinessID:        &businessID,
	}

	webhooks := []interface{}{map[string]interface{}{"id": "1", "topic": "orders/create"}}

	repo.On("GetIntegrationByID", mock.Anything, uint(5)).Return(integracion, nil)
	provider.On("ListWebhooks", mock.Anything, "5").Return(webhooks, nil)

	// Act
	result, err := uc.ListWebhooks(ctx, "5")

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestListWebhooks_ErrorDeProvider(t *testing.T) {
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

	businessID := uint(1)
	integracion := &domain.Integration{
		ID:                5,
		IntegrationTypeID: uint(domain.IntegrationTypeShopify),
		BusinessID:        &businessID,
	}

	repo.On("GetIntegrationByID", mock.Anything, uint(5)).Return(integracion, nil)
	provider.On("ListWebhooks", mock.Anything, "5").Return(nil, errors.New("api error"))

	// Act
	result, err := uc.ListWebhooks(ctx, "5")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
}

// ============================================
// DeleteWebhook
// ============================================

func TestDeleteWebhook_Exitoso(t *testing.T) {
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

	businessID := uint(1)
	integracion := &domain.Integration{
		ID:                5,
		IntegrationTypeID: uint(domain.IntegrationTypeShopify),
		BusinessID:        &businessID,
	}

	repo.On("GetIntegrationByID", mock.Anything, uint(5)).Return(integracion, nil)
	provider.On("DeleteWebhook", mock.Anything, "5", "webhook_123").Return(nil)

	// Act
	err := uc.DeleteWebhook(ctx, "5", "webhook_123")

	// Assert
	assert.NoError(t, err)
	provider.AssertCalled(t, "DeleteWebhook", mock.Anything, "5", "webhook_123")
}

func TestDeleteWebhook_ErrorDeProvider(t *testing.T) {
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

	businessID := uint(1)
	integracion := &domain.Integration{
		ID:                5,
		IntegrationTypeID: uint(domain.IntegrationTypeShopify),
		BusinessID:        &businessID,
	}

	repo.On("GetIntegrationByID", mock.Anything, uint(5)).Return(integracion, nil)
	provider.On("DeleteWebhook", mock.Anything, "5", "webhook_999").Return(errors.New("webhook not found"))

	// Act
	err := uc.DeleteWebhook(ctx, "5", "webhook_999")

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "webhook not found")
}

// ============================================
// VerifyWebhooksByURL
// ============================================

func TestVerifyWebhooksByURL_Exitoso(t *testing.T) {
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

	businessID := uint(1)
	integracion := &domain.Integration{
		ID:                5,
		IntegrationTypeID: uint(domain.IntegrationTypeShopify),
		BusinessID:        &businessID,
	}

	matchedWebhooks := []interface{}{map[string]interface{}{"id": "abc"}}

	cfg.On("Get", "WEBHOOK_BASE_URL").Return("https://api.example.com")
	repo.On("GetIntegrationByID", mock.Anything, uint(5)).Return(integracion, nil)
	provider.On("VerifyWebhooksByURL", mock.Anything, "5", "https://api.example.com").Return(matchedWebhooks, nil)

	// Act
	result, err := uc.VerifyWebhooksByURL(ctx, "5")

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestVerifyWebhooksByURL_FaltaBaseURL(t *testing.T) {
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

	businessID := uint(1)
	integracion := &domain.Integration{
		ID:                5,
		IntegrationTypeID: uint(domain.IntegrationTypeShopify),
		BusinessID:        &businessID,
	}

	cfg.On("Get", "WEBHOOK_BASE_URL").Return("")
	cfg.On("Get", "URL_BASE_SWAGGER").Return("")
	repo.On("GetIntegrationByID", mock.Anything, uint(5)).Return(integracion, nil)

	// Act
	result, err := uc.VerifyWebhooksByURL(ctx, "5")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
}

// ============================================
// CreateWebhookForIntegration
// ============================================

func TestCreateWebhookForIntegration_Exitoso(t *testing.T) {
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

	businessID := uint(1)
	integracion := &domain.Integration{
		ID:                5,
		IntegrationTypeID: uint(domain.IntegrationTypeShopify),
		BusinessID:        &businessID,
	}

	createdWebhook := map[string]interface{}{"id": "new_webhook_id"}

	cfg.On("Get", "WEBHOOK_BASE_URL").Return("https://api.example.com")
	repo.On("GetIntegrationByID", mock.Anything, uint(5)).Return(integracion, nil)
	provider.On("CreateWebhook", mock.Anything, "5", "https://api.example.com").Return(createdWebhook, nil)

	// Act
	result, err := uc.CreateWebhookForIntegration(ctx, "5")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
}

func TestCreateWebhookForIntegration_FaltaBaseURL(t *testing.T) {
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

	businessID := uint(1)
	integracion := &domain.Integration{
		ID:                5,
		IntegrationTypeID: uint(domain.IntegrationTypeShopify),
		BusinessID:        &businessID,
	}

	cfg.On("Get", "WEBHOOK_BASE_URL").Return("")
	cfg.On("Get", "URL_BASE_SWAGGER").Return("")
	repo.On("GetIntegrationByID", mock.Anything, uint(5)).Return(integracion, nil)

	// Act
	result, err := uc.CreateWebhookForIntegration(ctx, "5")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
}

// ============================================
// WarmCache
// ============================================

func TestWarmCache_SinIntegracionesActivas(t *testing.T) {
	// Arrange
	repo := new(mocks.RepositoryMock)
	enc := new(mocks.EncryptionMock)
	cache := new(mocks.CacheMock)
	logger := new(mocks.LoggerMock)
	cfg := new(mocks.ConfigMock)
	configurarLoggerPermisivo(logger)

	uc := newTestUseCase(repo, enc, cache, logger, cfg)
	ctx := context.Background()

	// Sin integraciones activas
	repo.On("ListIntegrations", mock.Anything, mock.Anything).Return([]*domain.Integration{}, int64(0), nil)

	// Act
	err := uc.WarmCache(ctx)

	// Assert
	assert.NoError(t, err)
	cache.AssertNotCalled(t, "SetIntegration", mock.Anything, mock.Anything)
}

func TestWarmCache_ErrorDeRepositorio(t *testing.T) {
	// Arrange
	repo := new(mocks.RepositoryMock)
	enc := new(mocks.EncryptionMock)
	cache := new(mocks.CacheMock)
	logger := new(mocks.LoggerMock)
	cfg := new(mocks.ConfigMock)
	configurarLoggerPermisivo(logger)

	uc := newTestUseCase(repo, enc, cache, logger, cfg)
	ctx := context.Background()

	repo.On("ListIntegrations", mock.Anything, mock.Anything).Return(nil, int64(0), errors.New("db error"))

	// Act
	err := uc.WarmCache(ctx)

	// Assert
	assert.Error(t, err)
}

func TestWarmCache_CacheaIntegracionesSinCredenciales(t *testing.T) {
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
		ID:          10,
		Name:        "Shopify Store",
		Code:        "shopify_001",
		BusinessID:  &businessID,
		IsActive:    true,
		Credentials: nil, // Sin credenciales
	}

	repo.On("ListIntegrations", mock.Anything, mock.Anything).Return([]*domain.Integration{integracion}, int64(1), nil)
	cache.On("SetIntegration", mock.Anything, mock.Anything).Return(nil)

	// Act
	err := uc.WarmCache(ctx)

	// Assert
	assert.NoError(t, err)
	cache.AssertCalled(t, "SetIntegration", mock.Anything, mock.Anything)
	// No debe intentar cachear credenciales si no hay
	cache.AssertNotCalled(t, "SetCredentials", mock.Anything, mock.Anything)
}

func TestWarmCache_ContinuaSiErrorAlCachearUnaIntegracion(t *testing.T) {
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
	integraciones := []*domain.Integration{
		{ID: 1, Name: "Shopify", Code: "shopify", BusinessID: &businessID, IsActive: true},
		{ID: 2, Name: "Factus", Code: "factus", BusinessID: &businessID, IsActive: true},
	}

	repo.On("ListIntegrations", mock.Anything, mock.Anything).Return(integraciones, int64(2), nil)
	// La primera falla al cachear, la segunda funciona
	cache.On("SetIntegration", mock.Anything, mock.MatchedBy(func(i *domain.CachedIntegration) bool {
		return i.ID == 1
	})).Return(errors.New("redis error"))
	cache.On("SetIntegration", mock.Anything, mock.MatchedBy(func(i *domain.CachedIntegration) bool {
		return i.ID == 2
	})).Return(nil)

	// Act
	err := uc.WarmCache(ctx)

	// Assert — no debe retornar error aunque una falle
	assert.NoError(t, err)
}

// ============================================
// GetIntegrationByIDWithCredentials
// ============================================

func TestGetIntegrationByIDWithCredentials_CacheHitCredenciales(t *testing.T) {
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
		ID:                10,
		Name:              "Shopify",
		Code:              "shopify_001",
		IntegrationTypeID: 1,
		BusinessID:        &businessID,
	}

	credenciales := map[string]interface{}{"api_key": "key_secreto"}
	cachedCreds := &domain.CachedCredentials{
		IntegrationID: 10,
		Credentials:   credenciales,
	}

	cache.On("GetIntegration", mock.Anything, uint(10)).Return(cachedInteg, nil)
	repo.On("GetIntegrationTypeByID", mock.Anything, uint(1)).Return(nil, nil)
	cache.On("GetCredentials", mock.Anything, uint(10)).Return(cachedCreds, nil)

	// Act
	result, err := uc.GetIntegrationByIDWithCredentials(ctx, 10)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "key_secreto", result.DecryptedCredentials["api_key"])
	enc.AssertNotCalled(t, "DecryptCredentials", mock.Anything, mock.Anything)
}

func TestGetIntegrationByIDWithCredentials_DesencriptaDesdeBD(t *testing.T) {
	// Arrange
	repo := new(mocks.RepositoryMock)
	enc := new(mocks.EncryptionMock)
	cache := new(mocks.CacheMock)
	logger := new(mocks.LoggerMock)
	cfg := new(mocks.ConfigMock)
	configurarLoggerPermisivo(logger)

	uc := newTestUseCase(repo, enc, cache, logger, cfg)
	ctx := context.Background()

	credJSON := `{"encrypted": "dGVzdA=="}`
	integracion := &domain.Integration{
		ID:          10,
		Credentials: []byte(credJSON),
	}

	credenciales := map[string]interface{}{"api_key": "valor_real"}

	// Cache miss de integración
	cache.On("GetIntegration", mock.Anything, uint(10)).Return(nil, errors.New("miss"))
	repo.On("GetIntegrationByID", mock.Anything, uint(10)).Return(integracion, nil)
	cache.On("SetIntegration", mock.Anything, mock.Anything).Return(nil)
	// Cache miss de credenciales
	cache.On("GetCredentials", mock.Anything, uint(10)).Return(nil, errors.New("miss"))
	// Desencripta correctamente
	enc.On("DecryptCredentials", mock.Anything, mock.Anything).Return(credenciales, nil)
	// Cachea las credenciales nuevas
	cache.On("SetCredentials", mock.Anything, mock.Anything).Return(nil)

	// Act
	result, err := uc.GetIntegrationByIDWithCredentials(ctx, 10)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "valor_real", result.DecryptedCredentials["api_key"])
	enc.AssertCalled(t, "DecryptCredentials", mock.Anything, mock.Anything)
}

func TestGetIntegrationByIDWithCredentials_IntegracionNoEncontrada(t *testing.T) {
	// Arrange
	repo := new(mocks.RepositoryMock)
	enc := new(mocks.EncryptionMock)
	cache := new(mocks.CacheMock)
	logger := new(mocks.LoggerMock)
	cfg := new(mocks.ConfigMock)
	configurarLoggerPermisivo(logger)

	uc := newTestUseCase(repo, enc, cache, logger, cfg)
	ctx := context.Background()

	cache.On("GetIntegration", mock.Anything, uint(999)).Return(nil, errors.New("miss"))
	repo.On("GetIntegrationByID", mock.Anything, uint(999)).Return(nil, errors.New("not found"))

	// Act
	result, err := uc.GetIntegrationByIDWithCredentials(ctx, 999)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.ErrorIs(t, err, domain.ErrIntegrationNotFound)
}

func TestGetIntegrationByIDWithCredentials_SinCredenciales(t *testing.T) {
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
		ID:          7,
		Name:        "Factus",
		Credentials: nil, // Sin credenciales
	}

	cache.On("GetIntegration", mock.Anything, uint(7)).Return(nil, errors.New("miss"))
	repo.On("GetIntegrationByID", mock.Anything, uint(7)).Return(integracion, nil)
	cache.On("SetIntegration", mock.Anything, mock.Anything).Return(nil)
	cache.On("GetCredentials", mock.Anything, uint(7)).Return(nil, errors.New("miss"))

	// Act
	result, err := uc.GetIntegrationByIDWithCredentials(ctx, 7)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Empty(t, result.DecryptedCredentials)
	enc.AssertNotCalled(t, "DecryptCredentials", mock.Anything, mock.Anything)
}
