package core

import (
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/app/usecaseintegrations"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/domain"
)

func mapDomainToPublicIntegration(useCase usecaseintegrations.IIntegrationUseCase, integration *domain.Integration) *Integration {
	useCaseImpl, ok := useCase.(*usecaseintegrations.IntegrationUseCase)
	if !ok {
		return &Integration{
			ID:         integration.ID,
			BusinessID: integration.BusinessID,
			Name:       integration.Name,
			Config:     nil,
		}
	}

	publicIntegration := useCaseImpl.MapToPublicIntegration(integration)
	return &Integration{
		ID:              publicIntegration.ID,
		BusinessID:      publicIntegration.BusinessID,
		Name:            publicIntegration.Name,
		IntegrationType: publicIntegration.IntegrationType,
		Config:          publicIntegration.Config,
	}
}
