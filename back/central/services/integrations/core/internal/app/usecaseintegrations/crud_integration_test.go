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
// UpdateIntegration
// ============================================

func TestUpdateIntegration_ActualizaNombreYDescripcion(t *testing.T) {
	// Arrange
	repo := new(mocks.RepositoryMock)
	enc := new(mocks.EncryptionMock)
	cache := new(mocks.CacheMock)
	logger := new(mocks.LoggerMock)
	cfg := new(mocks.ConfigMock)
	configurarLoggerPermisivo(logger)

	uc := newTestUseCase(repo, enc, cache, logger, cfg)
	ctx := context.Background()

	integracionExistente := &domain.Integration{
		ID:                1,
		Name:              "Nombre viejo",
		Code:              "shopify_001",
		IntegrationTypeID: 1,
		IsActive:          true,
	}

	nuevoNombre := "Nombre nuevo"
	nuevaDesc := "Descripción actualizada"
	dto := domain.UpdateIntegrationDTO{
		Name:        &nuevoNombre,
		Description: &nuevaDesc,
		UpdatedByID: 5,
	}

	repo.On("GetIntegrationByID", mock.Anything, uint(1)).Return(integracionExistente, nil)
	cache.On("InvalidateIntegration", mock.Anything, uint(1)).Return(nil)
	repo.On("UpdateIntegration", mock.Anything, uint(1), mock.AnythingOfType("*domain.Integration")).Return(nil)
	repo.On("GetIntegrationTypeByID", mock.Anything, uint(1)).Return(nil, nil)
	cache.On("SetIntegration", mock.Anything, mock.Anything).Return(nil)

	// Act
	resultado, err := uc.UpdateIntegration(ctx, 1, dto)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, resultado)
	assert.Equal(t, nuevoNombre, resultado.Name)
	assert.Equal(t, nuevaDesc, resultado.Description)
	repo.AssertExpectations(t)
}

func TestUpdateIntegration_NoEncontrado(t *testing.T) {
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

	dto := domain.UpdateIntegrationDTO{}

	// Act
	resultado, err := uc.UpdateIntegration(ctx, 999, dto)

	// Assert
	assert.ErrorIs(t, err, domain.ErrIntegrationNotFound)
	assert.Nil(t, resultado)
}

func TestUpdateIntegration_ErrorAlGuardar(t *testing.T) {
	// Arrange
	repo := new(mocks.RepositoryMock)
	enc := new(mocks.EncryptionMock)
	cache := new(mocks.CacheMock)
	logger := new(mocks.LoggerMock)
	cfg := new(mocks.ConfigMock)
	configurarLoggerPermisivo(logger)

	uc := newTestUseCase(repo, enc, cache, logger, cfg)
	ctx := context.Background()

	integracionExistente := &domain.Integration{
		ID:                2,
		Name:              "Factus",
		Code:              "factus_001",
		IntegrationTypeID: 7,
	}

	nuevoNombre := "Factus actualizado"
	dto := domain.UpdateIntegrationDTO{Name: &nuevoNombre}

	repo.On("GetIntegrationByID", mock.Anything, uint(2)).Return(integracionExistente, nil)
	cache.On("InvalidateIntegration", mock.Anything, uint(2)).Return(nil)
	repo.On("UpdateIntegration", mock.Anything, uint(2), mock.AnythingOfType("*domain.Integration")).Return(errors.New("db error"))

	// Act
	resultado, err := uc.UpdateIntegration(ctx, 2, dto)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, resultado)
	assert.Contains(t, err.Error(), "error al actualizar integración")
}

func TestUpdateIntegration_MarcarComoDefault(t *testing.T) {
	// Arrange
	repo := new(mocks.RepositoryMock)
	enc := new(mocks.EncryptionMock)
	cache := new(mocks.CacheMock)
	logger := new(mocks.LoggerMock)
	cfg := new(mocks.ConfigMock)
	configurarLoggerPermisivo(logger)

	uc := newTestUseCase(repo, enc, cache, logger, cfg)
	ctx := context.Background()

	integracionExistente := &domain.Integration{
		ID:                3,
		Name:              "Siigo",
		Code:              "siigo_001",
		IntegrationTypeID: 8,
		IsDefault:         false,
	}

	esDefault := true
	dto := domain.UpdateIntegrationDTO{IsDefault: &esDefault}

	repo.On("GetIntegrationByID", mock.Anything, uint(3)).Return(integracionExistente, nil)
	repo.On("SetIntegrationAsDefault", mock.Anything, uint(3)).Return(nil)
	cache.On("InvalidateIntegration", mock.Anything, uint(3)).Return(nil)
	repo.On("UpdateIntegration", mock.Anything, uint(3), mock.AnythingOfType("*domain.Integration")).Return(nil)
	repo.On("GetIntegrationTypeByID", mock.Anything, uint(8)).Return(nil, nil)
	cache.On("SetIntegration", mock.Anything, mock.Anything).Return(nil)

	// Act
	resultado, err := uc.UpdateIntegration(ctx, 3, dto)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, resultado)
	assert.True(t, resultado.IsDefault)
	repo.AssertCalled(t, "SetIntegrationAsDefault", mock.Anything, uint(3))
}

// ============================================
// DeleteIntegration
// ============================================

func TestDeleteIntegration_Exitoso(t *testing.T) {
	// Arrange
	repo := new(mocks.RepositoryMock)
	enc := new(mocks.EncryptionMock)
	cache := new(mocks.CacheMock)
	logger := new(mocks.LoggerMock)
	cfg := new(mocks.ConfigMock)
	configurarLoggerPermisivo(logger)

	uc := newTestUseCase(repo, enc, cache, logger, cfg)
	ctx := context.Background()

	tipoShopify := &domain.IntegrationType{ID: 1, Code: "shopify"}
	integracion := &domain.Integration{
		ID:              10,
		Name:            "Shopify Store",
		IntegrationType: tipoShopify,
	}

	repo.On("GetIntegrationByID", mock.Anything, uint(10)).Return(integracion, nil)
	repo.On("DeleteIntegration", mock.Anything, uint(10)).Return(nil)
	cache.On("InvalidateIntegration", mock.Anything, uint(10)).Return(nil)

	// Act
	err := uc.DeleteIntegration(ctx, 10)

	// Assert
	assert.NoError(t, err)
	repo.AssertExpectations(t)
	cache.AssertCalled(t, "InvalidateIntegration", mock.Anything, uint(10))
}

func TestDeleteIntegration_NoEncontrada(t *testing.T) {
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
	err := uc.DeleteIntegration(ctx, 999)

	// Assert
	assert.ErrorIs(t, err, domain.ErrIntegrationNotFound)
	repo.AssertNotCalled(t, "DeleteIntegration", mock.Anything, mock.Anything)
}

func TestDeleteIntegration_NoSePuedeEliminarWhatsApp(t *testing.T) {
	// Arrange
	repo := new(mocks.RepositoryMock)
	enc := new(mocks.EncryptionMock)
	cache := new(mocks.CacheMock)
	logger := new(mocks.LoggerMock)
	cfg := new(mocks.ConfigMock)
	configurarLoggerPermisivo(logger)

	uc := newTestUseCase(repo, enc, cache, logger, cfg)
	ctx := context.Background()

	tipoWhatsapp := &domain.IntegrationType{ID: 2, Code: "whatsapp"}
	integracion := &domain.Integration{
		ID:              5,
		Name:            "WhatsApp Principal",
		IntegrationType: tipoWhatsapp,
	}

	repo.On("GetIntegrationByID", mock.Anything, uint(5)).Return(integracion, nil)

	// Act
	err := uc.DeleteIntegration(ctx, 5)

	// Assert
	assert.ErrorIs(t, err, domain.ErrIntegrationCannotDeleteWhatsApp)
	repo.AssertNotCalled(t, "DeleteIntegration", mock.Anything, mock.Anything)
}

// ============================================
// ActivateIntegration / DeactivateIntegration
// ============================================

func TestActivateIntegration_IntegracionYaActiva(t *testing.T) {
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
		ID:       1,
		IsActive: true, // Ya está activa
	}

	repo.On("GetIntegrationByID", mock.Anything, uint(1)).Return(integracion, nil)

	// Act
	err := uc.ActivateIntegration(ctx, 1)

	// Assert
	assert.NoError(t, err)
	// No debe llamar UpdateIntegration porque ya estaba activa
	repo.AssertNotCalled(t, "UpdateIntegration", mock.Anything, mock.Anything, mock.Anything)
}

func TestActivateIntegration_ActivaIntegracionInactiva(t *testing.T) {
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
		ID:                2,
		Name:              "Factus",
		Code:              "factus_001",
		IntegrationTypeID: 7,
		IsActive:          false, // Inactiva
	}

	repo.On("GetIntegrationByID", mock.Anything, uint(2)).Return(integracion, nil)
	cache.On("InvalidateIntegration", mock.Anything, uint(2)).Return(nil)
	repo.On("UpdateIntegration", mock.Anything, uint(2), mock.AnythingOfType("*domain.Integration")).Return(nil)
	repo.On("GetIntegrationTypeByID", mock.Anything, uint(7)).Return(nil, nil)
	cache.On("SetIntegration", mock.Anything, mock.Anything).Return(nil)

	// Act
	err := uc.ActivateIntegration(ctx, 2)

	// Assert
	assert.NoError(t, err)
	repo.AssertCalled(t, "UpdateIntegration", mock.Anything, uint(2), mock.AnythingOfType("*domain.Integration"))
}

func TestDeactivateIntegration_IntegracionYaInactiva(t *testing.T) {
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
		ID:       3,
		IsActive: false, // Ya inactiva
	}

	repo.On("GetIntegrationByID", mock.Anything, uint(3)).Return(integracion, nil)

	// Act
	err := uc.DeactivateIntegration(ctx, 3)

	// Assert
	assert.NoError(t, err)
	repo.AssertNotCalled(t, "UpdateIntegration", mock.Anything, mock.Anything, mock.Anything)
}

func TestDeactivateIntegration_DesactivaIntegracionActiva(t *testing.T) {
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
		ID:                4,
		Name:              "Siigo",
		Code:              "siigo_001",
		IntegrationTypeID: 8,
		IsActive:          true,
	}

	repo.On("GetIntegrationByID", mock.Anything, uint(4)).Return(integracion, nil)
	cache.On("InvalidateIntegration", mock.Anything, uint(4)).Return(nil)
	repo.On("UpdateIntegration", mock.Anything, uint(4), mock.AnythingOfType("*domain.Integration")).Return(nil)
	repo.On("GetIntegrationTypeByID", mock.Anything, uint(8)).Return(nil, nil)
	cache.On("SetIntegration", mock.Anything, mock.Anything).Return(nil)

	// Act
	err := uc.DeactivateIntegration(ctx, 4)

	// Assert
	assert.NoError(t, err)
	repo.AssertCalled(t, "UpdateIntegration", mock.Anything, uint(4), mock.AnythingOfType("*domain.Integration"))
}

// ============================================
// ListIntegrations
// ============================================

func TestListIntegrations_AplicaValoresPorDefecto(t *testing.T) {
	// Arrange
	repo := new(mocks.RepositoryMock)
	enc := new(mocks.EncryptionMock)
	cache := new(mocks.CacheMock)
	logger := new(mocks.LoggerMock)
	cfg := new(mocks.ConfigMock)
	configurarLoggerPermisivo(logger)

	uc := newTestUseCase(repo, enc, cache, logger, cfg)
	ctx := context.Background()

	integraciones := []*domain.Integration{
		{ID: 1, Name: "Shopify"},
		{ID: 2, Name: "Factus"},
	}

	// Page y PageSize de 0 — deben tomar valores por defecto
	filtros := domain.IntegrationFilters{Page: 0, PageSize: 0}
	repo.On("ListIntegrations", mock.Anything, mock.MatchedBy(func(f domain.IntegrationFilters) bool {
		return f.Page == 1 && f.PageSize == 10
	})).Return(integraciones, int64(2), nil)

	// Act
	resultado, total, err := uc.ListIntegrations(ctx, filtros)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, resultado, 2)
	assert.Equal(t, int64(2), total)
}

func TestListIntegrations_LimitaPageSizeA100(t *testing.T) {
	// Arrange
	repo := new(mocks.RepositoryMock)
	enc := new(mocks.EncryptionMock)
	cache := new(mocks.CacheMock)
	logger := new(mocks.LoggerMock)
	cfg := new(mocks.ConfigMock)
	configurarLoggerPermisivo(logger)

	uc := newTestUseCase(repo, enc, cache, logger, cfg)
	ctx := context.Background()

	filtros := domain.IntegrationFilters{Page: 1, PageSize: 500} // Excede el máximo
	repo.On("ListIntegrations", mock.Anything, mock.MatchedBy(func(f domain.IntegrationFilters) bool {
		return f.PageSize == 100 // Debe limitarse a 100
	})).Return([]*domain.Integration{}, int64(0), nil)

	// Act
	_, _, err := uc.ListIntegrations(ctx, filtros)

	// Assert
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestListIntegrations_ErrorDeRepositorio(t *testing.T) {
	// Arrange
	repo := new(mocks.RepositoryMock)
	enc := new(mocks.EncryptionMock)
	cache := new(mocks.CacheMock)
	logger := new(mocks.LoggerMock)
	cfg := new(mocks.ConfigMock)
	configurarLoggerPermisivo(logger)

	uc := newTestUseCase(repo, enc, cache, logger, cfg)
	ctx := context.Background()

	filtros := domain.IntegrationFilters{Page: 1, PageSize: 10}
	repo.On("ListIntegrations", mock.Anything, mock.Anything).Return(nil, int64(0), errors.New("db error"))

	// Act
	resultado, total, err := uc.ListIntegrations(ctx, filtros)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, resultado)
	assert.Equal(t, int64(0), total)
}

// ============================================
// SetAsDefault
// ============================================

func TestSetAsDefault_Exitoso(t *testing.T) {
	// Arrange
	repo := new(mocks.RepositoryMock)
	enc := new(mocks.EncryptionMock)
	cache := new(mocks.CacheMock)
	logger := new(mocks.LoggerMock)
	cfg := new(mocks.ConfigMock)
	configurarLoggerPermisivo(logger)

	uc := newTestUseCase(repo, enc, cache, logger, cfg)
	ctx := context.Background()

	repo.On("SetIntegrationAsDefault", mock.Anything, uint(7)).Return(nil)

	// Act
	err := uc.SetAsDefault(ctx, 7)

	// Assert
	assert.NoError(t, err)
	repo.AssertCalled(t, "SetIntegrationAsDefault", mock.Anything, uint(7))
}

func TestSetAsDefault_ErrorDeRepositorio(t *testing.T) {
	// Arrange
	repo := new(mocks.RepositoryMock)
	enc := new(mocks.EncryptionMock)
	cache := new(mocks.CacheMock)
	logger := new(mocks.LoggerMock)
	cfg := new(mocks.ConfigMock)
	configurarLoggerPermisivo(logger)

	uc := newTestUseCase(repo, enc, cache, logger, cfg)
	ctx := context.Background()

	repo.On("SetIntegrationAsDefault", mock.Anything, uint(7)).Return(errors.New("constraint error"))

	// Act
	err := uc.SetAsDefault(ctx, 7)

	// Assert
	assert.Error(t, err)
}

// ============================================
// UpdateLastSync
// ============================================

func TestUpdateLastSync_Exitoso(t *testing.T) {
	// Arrange
	repo := new(mocks.RepositoryMock)
	enc := new(mocks.EncryptionMock)
	cache := new(mocks.CacheMock)
	logger := new(mocks.LoggerMock)
	cfg := new(mocks.ConfigMock)
	configurarLoggerPermisivo(logger)

	uc := newTestUseCase(repo, enc, cache, logger, cfg)
	ctx := context.Background()

	repo.On("UpdateLastSync", mock.Anything, uint(3), mock.Anything).Return(nil)

	// Act
	err := uc.UpdateLastSync(ctx, "3")

	// Assert
	assert.NoError(t, err)
	repo.AssertCalled(t, "UpdateLastSync", mock.Anything, uint(3), mock.Anything)
}

func TestUpdateLastSync_IDInvalido(t *testing.T) {
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
	err := uc.UpdateLastSync(ctx, "no_es_numero")

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid integration ID")
	repo.AssertNotCalled(t, "UpdateLastSync", mock.Anything, mock.Anything, mock.Anything)
}
