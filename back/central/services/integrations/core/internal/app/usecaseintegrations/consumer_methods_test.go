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
// GetIntegrationByExternalID
// ============================================

func TestGetIntegrationByExternalID_Exitoso(t *testing.T) {
	// Arrange
	repo := new(mocks.RepositoryMock)
	enc := new(mocks.EncryptionMock)
	cache := new(mocks.CacheMock)
	logger := new(mocks.LoggerMock)
	cfg := new(mocks.ConfigMock)
	configurarLoggerPermisivo(logger)

	uc := newTestUseCase(repo, enc, cache, logger, cfg)
	ctx := context.Background()

	externalID := "mi-tienda.myshopify.com"
	integrationType := 1 // Shopify

	businessID := uint(5)
	integracion := &domain.Integration{
		ID:                20,
		Name:              "Shopify Store",
		BusinessID:        &businessID,
		IntegrationTypeID: 1,
		StoreID:           externalID,
	}

	repo.On("ListIntegrations", mock.Anything, mock.MatchedBy(func(f domain.IntegrationFilters) bool {
		return f.StoreID != nil && *f.StoreID == externalID && f.Page == 1 && f.PageSize == 1
	})).Return([]*domain.Integration{integracion}, int64(1), nil)

	// Act
	resultado, err := uc.GetIntegrationByExternalID(ctx, externalID, integrationType)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, resultado)
	assert.Equal(t, uint(20), resultado.ID)
	assert.Equal(t, "Shopify Store", resultado.Name)
	assert.Equal(t, externalID, resultado.StoreID)
}

func TestGetIntegrationByExternalID_NoEncontrado(t *testing.T) {
	// Arrange
	repo := new(mocks.RepositoryMock)
	enc := new(mocks.EncryptionMock)
	cache := new(mocks.CacheMock)
	logger := new(mocks.LoggerMock)
	cfg := new(mocks.ConfigMock)
	configurarLoggerPermisivo(logger)

	uc := newTestUseCase(repo, enc, cache, logger, cfg)
	ctx := context.Background()

	repo.On("ListIntegrations", mock.Anything, mock.Anything).Return([]*domain.Integration{}, int64(0), nil)

	// Act
	resultado, err := uc.GetIntegrationByExternalID(ctx, "tienda-no-existe.myshopify.com", 1)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, resultado)
	assert.Contains(t, err.Error(), "integration not found")
}

func TestGetIntegrationByExternalID_ErrorDeRepositorio(t *testing.T) {
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
	resultado, err := uc.GetIntegrationByExternalID(ctx, "tienda.myshopify.com", 1)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, resultado)
	assert.Contains(t, err.Error(), "error listing integrations")
}

// ============================================
// UpdateIntegrationConfig
// ============================================

func TestUpdateIntegrationConfig_MergeExitoso(t *testing.T) {
	// Arrange
	repo := new(mocks.RepositoryMock)
	enc := new(mocks.EncryptionMock)
	cache := new(mocks.CacheMock)
	logger := new(mocks.LoggerMock)
	cfg := new(mocks.ConfigMock)
	configurarLoggerPermisivo(logger)

	uc := newTestUseCase(repo, enc, cache, logger, cfg)
	ctx := context.Background()

	// Config existente en la integración
	configJSON := `{"tax_rate": 0.19, "currency": "COP"}`
	businessID := uint(1)
	integracion := &domain.Integration{
		ID:                9,
		IntegrationTypeID: 7,
		Config:            []byte(configJSON),
		BusinessID:        &businessID,
	}

	repo.On("GetIntegrationByID", mock.Anything, uint(9)).Return(integracion, nil)
	cache.On("InvalidateIntegration", mock.Anything, uint(9)).Return(nil)
	repo.On("UpdateIntegration", mock.Anything, uint(9), mock.AnythingOfType("*domain.Integration")).Return(nil)
	repo.On("GetIntegrationTypeByID", mock.Anything, uint(7)).Return(nil, nil)
	cache.On("SetIntegration", mock.Anything, mock.Anything).Return(nil)

	// Nuevos campos a mergear
	nuevoConfig := map[string]interface{}{"last_sync_token": "tok_abc123"}

	// Act
	err := uc.UpdateIntegrationConfig(ctx, "9", nuevoConfig)

	// Assert
	assert.NoError(t, err)
	repo.AssertCalled(t, "UpdateIntegration", mock.Anything, uint(9), mock.AnythingOfType("*domain.Integration"))
}

func TestUpdateIntegrationConfig_IDInvalido(t *testing.T) {
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
	err := uc.UpdateIntegrationConfig(ctx, "no_es_numero", map[string]interface{}{})

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ID de integración inválido")
}

// ============================================
// OnIntegrationCreated
// ============================================

func TestOnIntegrationCreated_ObservadorEjecutadoPorTipo(t *testing.T) {
	// Arrange
	repo := new(mocks.RepositoryMock)
	enc := new(mocks.EncryptionMock)
	cache := new(mocks.CacheMock)
	logger := new(mocks.LoggerMock)
	cfg := new(mocks.ConfigMock)
	configurarLoggerPermisivo(logger)

	uc := newTestUseCase(repo, enc, cache, logger, cfg)

	llamado := false
	uc.OnIntegrationCreated(domain.IntegrationTypeShopify, func(ctx context.Context, pi *domain.PublicIntegration) {
		llamado = true
	})

	// Simular notificación de integración creada
	tipoShopify := &domain.IntegrationType{ID: 1, Code: "shopify"}
	integracion := &domain.Integration{
		ID:              99,
		IntegrationType: tipoShopify,
	}

	// Invocar observadores directamente
	for _, obs := range uc.observers {
		obs(context.Background(), integracion)
	}

	// Assert
	assert.True(t, llamado, "El observador debería haberse ejecutado para tipo Shopify")
}

func TestOnIntegrationCreated_ObservadorNoEjecutadoPorOtroTipo(t *testing.T) {
	// Arrange
	repo := new(mocks.RepositoryMock)
	enc := new(mocks.EncryptionMock)
	cache := new(mocks.CacheMock)
	logger := new(mocks.LoggerMock)
	cfg := new(mocks.ConfigMock)
	configurarLoggerPermisivo(logger)

	uc := newTestUseCase(repo, enc, cache, logger, cfg)

	llamadoShopify := false
	uc.OnIntegrationCreated(domain.IntegrationTypeShopify, func(ctx context.Context, pi *domain.PublicIntegration) {
		llamadoShopify = true
	})

	// Integración de tipo Factus (no Shopify)
	tipoFactus := &domain.IntegrationType{ID: 7, Code: "factus"}
	integracion := &domain.Integration{
		ID:              50,
		IntegrationType: tipoFactus,
	}

	for _, obs := range uc.observers {
		obs(context.Background(), integracion)
	}

	// Assert
	assert.False(t, llamadoShopify, "El observador de Shopify no debe ejecutarse para tipo Factus")
}

// ============================================
// decodeEncryptedCredentials (función interna)
// ============================================

func TestDecodeEncryptedCredentials_FormatoWrapper(t *testing.T) {
	// Arrange — formato {"encrypted": "<base64>"}
	// "test" en base64 es "dGVzdA=="
	input := []byte(`{"encrypted": "dGVzdA=="}`)

	// Act
	resultado, err := decodeEncryptedCredentials(input)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, []byte("test"), resultado)
}

func TestDecodeEncryptedCredentials_Base64Directo(t *testing.T) {
	// Arrange — string base64 sin wrapper JSON
	// "hello" en base64 es "aGVsbG8="
	input := []byte("aGVsbG8=")

	// Act
	resultado, err := decodeEncryptedCredentials(input)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, []byte("hello"), resultado)
}

func TestDecodeEncryptedCredentials_SinCampoEncrypted(t *testing.T) {
	// Arrange — JSON válido pero sin campo "encrypted"
	input := []byte(`{"otro_campo": "valor"}`)

	// Act
	_, err := decodeEncryptedCredentials(input)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "encrypted")
}
