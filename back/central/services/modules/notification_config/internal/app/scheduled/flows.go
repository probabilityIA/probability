package scheduled

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/entities"
)

func (uc *useCase) flowRoot(ctx context.Context, businessID, flowID uint) (*uint, error) {
	if uc.flows == nil {
		return nil, fmt.Errorf("los flujos no estan disponibles")
	}

	flow, err := uc.flows.GetByID(ctx, flowID, businessID)
	if err != nil {
		return nil, err
	}
	if flow == nil {
		return nil, fmt.Errorf("el flujo de la tarea programada ya no existe")
	}
	if flow.RootTemplateID == nil || *flow.RootTemplateID == 0 {
		return nil, fmt.Errorf("el flujo %q no tiene plantilla inicial", flow.Name)
	}

	return flow.RootTemplateID, nil
}

func (uc *useCase) assertFlowApproved(ctx context.Context, businessID, flowID uint) error {
	if uc.steps == nil {
		return nil
	}

	steps, err := uc.steps.ListByFlow(ctx, businessID, flowID)
	if err != nil {
		return err
	}

	for _, step := range steps {
		if step.TargetStatus != entities.TemplateStatusApproved {
			return fmt.Errorf(
				"la respuesta %q del flujo esta en estado %s: esa rama no contestaria, espera a que Meta la apruebe",
				step.TargetName,
				step.TargetStatus,
			)
		}
	}

	return nil
}
