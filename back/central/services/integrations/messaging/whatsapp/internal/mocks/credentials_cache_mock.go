package mocks

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/ports"
)

type CredentialsCacheMock struct {
	GetWhatsAppConfigFn              func(ctx context.Context, businessID uint) (*ports.WhatsAppConfig, error)
	GetWhatsAppConfigByIntegrationFn func(ctx context.Context, integrationID uint) (*ports.WhatsAppConfig, error)
	GetWhatsAppDefaultConfigFn       func(ctx context.Context) (*ports.WhatsAppConfig, error)
	ResolveByPhoneNumberIDFn         func(ctx context.Context, phoneNumberID string) (*ports.IntegrationOwner, error)
}

func (m *CredentialsCacheMock) GetWhatsAppConfig(ctx context.Context, businessID uint) (*ports.WhatsAppConfig, error) {
	if m.GetWhatsAppConfigFn != nil {
		return m.GetWhatsAppConfigFn(ctx, businessID)
	}
	return &ports.WhatsAppConfig{PhoneNumberID: 123456, AccessToken: "mock-token"}, nil
}

func (m *CredentialsCacheMock) GetWhatsAppConfigByIntegration(ctx context.Context, integrationID uint) (*ports.WhatsAppConfig, error) {
	if m.GetWhatsAppConfigByIntegrationFn != nil {
		return m.GetWhatsAppConfigByIntegrationFn(ctx, integrationID)
	}
	return &ports.WhatsAppConfig{PhoneNumberID: 123456, AccessToken: "mock-token", IntegrationID: integrationID}, nil
}

func (m *CredentialsCacheMock) GetWhatsAppDefaultConfig(ctx context.Context) (*ports.WhatsAppConfig, error) {
	if m.GetWhatsAppDefaultConfigFn != nil {
		return m.GetWhatsAppDefaultConfigFn(ctx)
	}
	return &ports.WhatsAppConfig{PhoneNumberID: 999999, AccessToken: "mock-default-token"}, nil
}

func (m *CredentialsCacheMock) ResolveByPhoneNumberID(ctx context.Context, phoneNumberID string) (*ports.IntegrationOwner, error) {
	if m.ResolveByPhoneNumberIDFn != nil {
		return m.ResolveByPhoneNumberIDFn(ctx, phoneNumberID)
	}
	return nil, nil
}
