package mocks

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/ports"
)

// IntegrationRepositoryMock implementa ports.IIntegrationRepository para tests unitarios
type IntegrationRepositoryMock struct {
	GetWhatsAppConfigFn        func(ctx context.Context, businessID uint) (*ports.WhatsAppConfig, error)
	GetWhatsAppDefaultConfigFn func(ctx context.Context) (*ports.WhatsAppConfig, error)
}

func (m *IntegrationRepositoryMock) GetWhatsAppConfig(ctx context.Context, businessID uint) (*ports.WhatsAppConfig, error) {
	if m.GetWhatsAppConfigFn != nil {
		return m.GetWhatsAppConfigFn(ctx, businessID)
	}
	return &ports.WhatsAppConfig{
		PhoneNumberID: 99999,
		AccessToken:   "mock-access-token",
		IntegrationID: 1,
	}, nil
}

func (m *IntegrationRepositoryMock) GetWhatsAppDefaultConfig(ctx context.Context) (*ports.WhatsAppConfig, error) {
	if m.GetWhatsAppDefaultConfigFn != nil {
		return m.GetWhatsAppDefaultConfigFn(ctx)
	}
	return &ports.WhatsAppConfig{
		PhoneNumberID: 99999,
		AccessToken:   "mock-access-token",
		IntegrationID: 0,
	}, nil
}
