package mocks

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/core"
)

// IntegrationCoreMock es el mock de core.IIntegrationService.
// Implementa todos los métodos de la interfaz de servicio del core.
type IntegrationCoreMock struct {
	GetIntegrationByIDFn         func(ctx context.Context, integrationID string) (*core.PublicIntegration, error)
	GetIntegrationByExternalIDFn func(ctx context.Context, externalID string, integrationType int) (*core.PublicIntegration, error)
	DecryptCredentialFn          func(ctx context.Context, integrationID string, fieldName string) (string, error)
	UpdateIntegrationConfigFn    func(ctx context.Context, integrationID string, newConfig map[string]interface{}) error
	GetIntegrationConfigFn       func(ctx context.Context, integrationID string) (map[string]interface{}, error)
	GetPlatformCredentialFn      func(ctx context.Context, integrationID string, fieldName string) (string, error)
}

// GetIntegrationByID implementa core.IIntegrationService.
func (m *IntegrationCoreMock) GetIntegrationByID(ctx context.Context, integrationID string) (*core.PublicIntegration, error) {
	if m.GetIntegrationByIDFn != nil {
		return m.GetIntegrationByIDFn(ctx, integrationID)
	}
	return nil, nil
}

// GetIntegrationByExternalID implementa core.IIntegrationService.
func (m *IntegrationCoreMock) GetIntegrationByExternalID(ctx context.Context, externalID string, integrationType int) (*core.PublicIntegration, error) {
	if m.GetIntegrationByExternalIDFn != nil {
		return m.GetIntegrationByExternalIDFn(ctx, externalID, integrationType)
	}
	return nil, nil
}

// DecryptCredential implementa core.IIntegrationService.
func (m *IntegrationCoreMock) DecryptCredential(ctx context.Context, integrationID string, fieldName string) (string, error) {
	if m.DecryptCredentialFn != nil {
		return m.DecryptCredentialFn(ctx, integrationID, fieldName)
	}
	return "", nil
}

// UpdateIntegrationConfig implementa core.IIntegrationService.
func (m *IntegrationCoreMock) UpdateIntegrationConfig(ctx context.Context, integrationID string, newConfig map[string]interface{}) error {
	if m.UpdateIntegrationConfigFn != nil {
		return m.UpdateIntegrationConfigFn(ctx, integrationID, newConfig)
	}
	return nil
}

// GetIntegrationConfig implementa core.IIntegrationService.
func (m *IntegrationCoreMock) GetIntegrationConfig(ctx context.Context, integrationID string) (map[string]interface{}, error) {
	if m.GetIntegrationConfigFn != nil {
		return m.GetIntegrationConfigFn(ctx, integrationID)
	}
	return nil, nil
}

// GetPlatformCredential implementa core.IIntegrationService.
func (m *IntegrationCoreMock) GetPlatformCredential(ctx context.Context, integrationID string, fieldName string) (string, error) {
	if m.GetPlatformCredentialFn != nil {
		return m.GetPlatformCredentialFn(ctx, integrationID, fieldName)
	}
	return "", nil
}
