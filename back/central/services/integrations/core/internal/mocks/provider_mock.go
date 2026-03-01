package mocks

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/core/internal/domain"
	"github.com/stretchr/testify/mock"
)

// ProviderMock es un mock de domain.IIntegrationContract usando testify/mock
type ProviderMock struct {
	mock.Mock
}

func (m *ProviderMock) TestConnection(ctx context.Context, config map[string]interface{}, credentials map[string]interface{}) error {
	args := m.Called(ctx, config, credentials)
	return args.Error(0)
}

func (m *ProviderMock) SyncOrdersByIntegrationID(ctx context.Context, integrationID string) error {
	args := m.Called(ctx, integrationID)
	return args.Error(0)
}

func (m *ProviderMock) SyncOrdersByIntegrationIDWithParams(ctx context.Context, integrationID string, params interface{}) error {
	args := m.Called(ctx, integrationID, params)
	return args.Error(0)
}

func (m *ProviderMock) GetWebhookURL(ctx context.Context, baseURL string, integrationID uint) (*domain.WebhookInfo, error) {
	args := m.Called(ctx, baseURL, integrationID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.WebhookInfo), args.Error(1)
}

func (m *ProviderMock) ListWebhooks(ctx context.Context, integrationID string) ([]interface{}, error) {
	args := m.Called(ctx, integrationID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]interface{}), args.Error(1)
}

func (m *ProviderMock) DeleteWebhook(ctx context.Context, integrationID, webhookID string) error {
	args := m.Called(ctx, integrationID, webhookID)
	return args.Error(0)
}

func (m *ProviderMock) VerifyWebhooksByURL(ctx context.Context, integrationID string, baseURL string) ([]interface{}, error) {
	args := m.Called(ctx, integrationID, baseURL)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]interface{}), args.Error(1)
}

func (m *ProviderMock) CreateWebhook(ctx context.Context, integrationID string, baseURL string) (interface{}, error) {
	args := m.Called(ctx, integrationID, baseURL)
	return args.Get(0), args.Error(1)
}

func (m *ProviderMock) UpdateInventory(ctx context.Context, integrationID string, productExternalID string, quantity int) error {
	args := m.Called(ctx, integrationID, productExternalID, quantity)
	return args.Error(0)
}
