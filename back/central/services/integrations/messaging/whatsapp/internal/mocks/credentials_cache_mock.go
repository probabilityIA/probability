package mocks

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/ports"
)

// CredentialsCacheMock implementa ports.ICredentialsCache para tests unitarios
type CredentialsCacheMock struct {
	GetWhatsAppConfigFn        func(ctx context.Context, businessID uint) (*ports.WhatsAppConfig, error)
	GetWhatsAppDefaultConfigFn func(ctx context.Context) (*ports.WhatsAppConfig, error)
}

func (m *CredentialsCacheMock) GetWhatsAppConfig(ctx context.Context, businessID uint) (*ports.WhatsAppConfig, error) {
	if m.GetWhatsAppConfigFn != nil {
		return m.GetWhatsAppConfigFn(ctx, businessID)
	}
	return &ports.WhatsAppConfig{PhoneNumberID: 123456, AccessToken: "mock-token"}, nil
}

func (m *CredentialsCacheMock) GetWhatsAppDefaultConfig(ctx context.Context) (*ports.WhatsAppConfig, error) {
	if m.GetWhatsAppDefaultConfigFn != nil {
		return m.GetWhatsAppDefaultConfigFn(ctx)
	}
	return &ports.WhatsAppConfig{PhoneNumberID: 999999, AccessToken: "mock-default-token"}, nil
}
