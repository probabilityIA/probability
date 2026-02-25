package mocks

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/core"
)

// IntegrationCoreMock mock de core.IIntegrationService para tests unitarios.
// Permite simular las respuestas de GetIntegrationByID y DecryptCredential
// sin necesidad de conectarse a base de datos ni servicios externos.
type IntegrationCoreMock struct {
	GetIntegrationByIDFn         func(ctx context.Context, integrationID string) (*core.PublicIntegration, error)
	GetIntegrationByExternalIDFn func(ctx context.Context, externalID string, integrationType int) (*core.PublicIntegration, error)
	DecryptCredentialFn          func(ctx context.Context, integrationID string, fieldName string) (string, error)
	UpdateIntegrationConfigFn    func(ctx context.Context, integrationID string, newConfig map[string]interface{}) error
}

// Verificar en tiempo de compilación que implementa la interfaz.
var _ core.IIntegrationService = (*IntegrationCoreMock)(nil)

func (m *IntegrationCoreMock) GetIntegrationByID(ctx context.Context, integrationID string) (*core.PublicIntegration, error) {
	if m.GetIntegrationByIDFn != nil {
		return m.GetIntegrationByIDFn(ctx, integrationID)
	}
	return &core.PublicIntegration{
		ID:     1,
		Config: map[string]interface{}{},
	}, nil
}

func (m *IntegrationCoreMock) GetIntegrationByExternalID(ctx context.Context, externalID string, integrationType int) (*core.PublicIntegration, error) {
	if m.GetIntegrationByExternalIDFn != nil {
		return m.GetIntegrationByExternalIDFn(ctx, externalID, integrationType)
	}
	return &core.PublicIntegration{}, nil
}

func (m *IntegrationCoreMock) DecryptCredential(ctx context.Context, integrationID string, fieldName string) (string, error) {
	if m.DecryptCredentialFn != nil {
		return m.DecryptCredentialFn(ctx, integrationID, fieldName)
	}
	// Valores por defecto útiles para tests de camino feliz
	defaults := map[string]string{
		"client_id":     "test-client-id",
		"client_secret": "test-client-secret",
		"username":      "test@example.com",
		"password":      "test-password",
		"api_url":       "",
	}
	if v, ok := defaults[fieldName]; ok {
		return v, nil
	}
	return "", nil
}

func (m *IntegrationCoreMock) UpdateIntegrationConfig(ctx context.Context, integrationID string, newConfig map[string]interface{}) error {
	if m.UpdateIntegrationConfigFn != nil {
		return m.UpdateIntegrationConfigFn(ctx, integrationID, newConfig)
	}
	return nil
}

func (m *IntegrationCoreMock) GetIntegrationConfig(ctx context.Context, integrationID string) (map[string]interface{}, error) {
	return map[string]interface{}{}, nil
}
