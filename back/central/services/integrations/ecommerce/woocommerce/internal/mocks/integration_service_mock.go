package mocks

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/woocommerce/internal/domain"
)

// IntegrationServiceMock mock de domain.IIntegrationService para tests unitarios.
// Permite simular las respuestas de GetIntegrationByID y DecryptCredential
// sin necesidad de conectarse a base de datos ni servicios externos.
type IntegrationServiceMock struct {
	GetIntegrationByIDFn    func(ctx context.Context, integrationID string) (*domain.Integration, error)
	DecryptCredentialFn     func(ctx context.Context, integrationID string, fieldName string) (string, error)
	UpdateIntegrationConfigFn func(ctx context.Context, integrationID string, config map[string]interface{}) error
}

// Verificar en tiempo de compilación que implementa la interfaz.
var _ domain.IIntegrationService = (*IntegrationServiceMock)(nil)

func (m *IntegrationServiceMock) GetIntegrationByID(ctx context.Context, integrationID string) (*domain.Integration, error) {
	if m.GetIntegrationByIDFn != nil {
		return m.GetIntegrationByIDFn(ctx, integrationID)
	}
	return &domain.Integration{
		ID:   1,
		Name: "Test WooCommerce",
		Config: map[string]interface{}{
			"store_url": "https://test-store.com",
		},
	}, nil
}

func (m *IntegrationServiceMock) DecryptCredential(ctx context.Context, integrationID string, fieldName string) (string, error) {
	if m.DecryptCredentialFn != nil {
		return m.DecryptCredentialFn(ctx, integrationID, fieldName)
	}
	// Valores por defecto útiles para tests de camino feliz
	defaults := map[string]string{
		"consumer_key":    "ck_test_consumer_key",
		"consumer_secret": "cs_test_consumer_secret",
	}
	if v, ok := defaults[fieldName]; ok {
		return v, nil
	}
	return "", nil
}

func (m *IntegrationServiceMock) UpdateIntegrationConfig(ctx context.Context, integrationID string, config map[string]interface{}) error {
	if m.UpdateIntegrationConfigFn != nil {
		return m.UpdateIntegrationConfigFn(ctx, integrationID, config)
	}
	return nil
}
