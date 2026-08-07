package usecaseintegrations

import (
	"context"
	"testing"

	"github.com/secamc93/probability/back/central/services/integrations/core/internal/domain"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateIntegration_RechazaStoreIDDeOtroNegocio(t *testing.T) {
	repo := new(mocks.RepositoryMock)
	enc := new(mocks.EncryptionMock)
	cache := new(mocks.CacheMock)
	logger := new(mocks.LoggerMock)
	cfg := new(mocks.ConfigMock)
	configurarLoggerPermisivo(logger)

	uc := newTestUseCase(repo, enc, cache, logger, cfg)

	ctx := context.Background()
	businessID := uint(46)
	otroNegocio := uint(37)

	dto := domain.CreateIntegrationDTO{
		Name:              "viga",
		Code:              "mercado_libre_viga",
		IntegrationTypeID: 3,
		BusinessID:        &businessID,
		StoreID:           "706343696",
	}

	repo.On("ExistsIntegrationByCode", mock.Anything, dto.Code, dto.BusinessID).Return(false, nil)
	repo.On("FindStoreIDOwner", mock.Anything, dto.StoreID, uint(3), uint(0)).
		Return(&domain.StoreIDOwner{IntegrationID: 250, BusinessID: &otroNegocio}, nil)

	resultado, err := uc.CreateIntegration(ctx, dto)

	assert.Error(t, err)
	assert.Nil(t, resultado)
	assert.ErrorIs(t, err, domain.ErrIntegrationStoreIDInUse)
	repo.AssertNotCalled(t, "CreateIntegration", mock.Anything, mock.Anything)
}

func TestCreateIntegration_RechazaStoreIDDelMismoNegocio(t *testing.T) {
	repo := new(mocks.RepositoryMock)
	enc := new(mocks.EncryptionMock)
	cache := new(mocks.CacheMock)
	logger := new(mocks.LoggerMock)
	cfg := new(mocks.ConfigMock)
	configurarLoggerPermisivo(logger)

	uc := newTestUseCase(repo, enc, cache, logger, cfg)

	ctx := context.Background()
	businessID := uint(26)

	dto := domain.CreateIntegrationDTO{
		Name:              "Meli Demo 2",
		Code:              "meli_demo_2",
		IntegrationTypeID: 3,
		BusinessID:        &businessID,
		StoreID:           "188190196",
	}

	repo.On("ExistsIntegrationByCode", mock.Anything, dto.Code, dto.BusinessID).Return(false, nil)
	repo.On("FindStoreIDOwner", mock.Anything, dto.StoreID, uint(3), uint(0)).
		Return(&domain.StoreIDOwner{IntegrationID: 248, BusinessID: &businessID}, nil)

	resultado, err := uc.CreateIntegration(ctx, dto)

	assert.Error(t, err)
	assert.Nil(t, resultado)
	assert.ErrorIs(t, err, domain.ErrIntegrationStoreIDSameBusiness)
	repo.AssertNotCalled(t, "CreateIntegration", mock.Anything, mock.Anything)
}

func TestCreateIntegration_PermiteStoreIDLibre(t *testing.T) {
	repo := new(mocks.RepositoryMock)
	enc := new(mocks.EncryptionMock)
	cache := new(mocks.CacheMock)
	logger := new(mocks.LoggerMock)
	cfg := new(mocks.ConfigMock)
	configurarLoggerPermisivo(logger)

	uc := newTestUseCase(repo, enc, cache, logger, cfg)

	ctx := context.Background()
	businessID := uint(26)

	categoria := &domain.IntegrationCategory{ID: 10, Code: "external"}
	tipoIntegracion := &domain.IntegrationType{
		ID:         5,
		Code:       "factus",
		CategoryID: 10,
		Category:   categoria,
	}

	dto := domain.CreateIntegrationDTO{
		Name:              "Cuenta nueva",
		Code:              "factus_002",
		IntegrationTypeID: 5,
		BusinessID:        &businessID,
		StoreID:           "999999999",
		Credentials:       map[string]interface{}{"access_token": "tok_123"},
	}

	repo.On("ExistsIntegrationByCode", mock.Anything, dto.Code, dto.BusinessID).Return(false, nil)
	repo.On("FindStoreIDOwner", mock.Anything, dto.StoreID, uint(5), uint(0)).Return(nil, nil)
	repo.On("GetIntegrationTypeByID", mock.Anything, uint(5)).Return(tipoIntegracion, nil)
	repo.On("CreateIntegration", mock.Anything, mock.AnythingOfType("*domain.Integration")).Return(nil)
	cache.On("SetIntegration", mock.Anything, mock.Anything).Return(nil)
	cache.On("SetCredentials", mock.Anything, mock.Anything).Return(nil)

	resultado, err := uc.CreateIntegration(ctx, dto)

	assert.NoError(t, err)
	assert.NotNil(t, resultado)
}
