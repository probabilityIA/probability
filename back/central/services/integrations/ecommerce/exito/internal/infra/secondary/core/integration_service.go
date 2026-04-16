package core

import (
	"context"

	integrationcore "github.com/secamc93/probability/back/central/services/integrations/core"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/exito/internal/domain"
)

// integrationServiceAdapter adapta core.IIntegrationService -> domain.IIntegrationService.
type integrationServiceAdapter struct {
	core integrationcore.IIntegrationService
}

// NewIntegrationService crea el adaptador de servicio de integracion para Exito.
func NewIntegrationService(core integrationcore.IIntegrationService) domain.IIntegrationService {
	return &integrationServiceAdapter{core: core}
}

func (a *integrationServiceAdapter) GetIntegrationByID(ctx context.Context, integrationID string) (*domain.Integration, error) {
	pub, err := a.core.GetIntegrationByID(ctx, integrationID)
	if err != nil {
		return nil, err
	}
	if pub == nil {
		return nil, nil
	}
	return &domain.Integration{
		ID:              pub.ID,
		BusinessID:      pub.BusinessID,
		Name:            pub.Name,
		StoreID:         pub.StoreID,
		IntegrationType: pub.IntegrationType,
		Config:          pub.Config,
	}, nil
}

func (a *integrationServiceAdapter) DecryptCredential(ctx context.Context, integrationID string, fieldName string) (string, error) {
	return a.core.DecryptCredential(ctx, integrationID, fieldName)
}

func (a *integrationServiceAdapter) UpdateIntegrationConfig(ctx context.Context, integrationID string, config map[string]interface{}) error {
	return a.core.UpdateIntegrationConfig(ctx, integrationID, config)
}
