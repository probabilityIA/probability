package templates

import (
	"context"
	"fmt"
	"strings"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/entities"
)

const maxFlowNameLength = 120

func (uc *useCase) ListFlowGroups(ctx context.Context, businessID uint) ([]entities.Flow, error) {
	if uc.flowGroups == nil {
		return nil, fmt.Errorf("los flujos no estan disponibles")
	}

	return uc.flowGroups.List(ctx, businessID)
}

func (uc *useCase) CreateFlowGroup(ctx context.Context, dto dtos.CreateFlowDTO) (*entities.Flow, error) {
	if uc.flowGroups == nil {
		return nil, fmt.Errorf("los flujos no estan disponibles")
	}

	name := strings.TrimSpace(dto.Name)
	if name == "" {
		return nil, fmt.Errorf("el nombre del flujo es obligatorio")
	}
	if len(name) > maxFlowNameLength {
		return nil, fmt.Errorf("el nombre del flujo supera %d caracteres", maxFlowNameLength)
	}

	taken, err := uc.flowGroups.ExistsByName(ctx, dto.BusinessID, name, 0)
	if err != nil {
		return nil, err
	}
	if taken {
		return nil, fmt.Errorf("ya existe un flujo llamado %q", name)
	}

	if err := uc.assertOwnedTemplate(ctx, dto.RootTemplateID, dto.BusinessID); err != nil {
		return nil, err
	}

	flow := &entities.Flow{
		BusinessID:     dto.BusinessID,
		Name:           name,
		Description:    strings.TrimSpace(dto.Description),
		RootTemplateID: dto.RootTemplateID,
		Enabled:        true,
	}

	if err := uc.flowGroups.Create(ctx, flow); err != nil {
		return nil, err
	}

	return uc.flowGroups.GetByID(ctx, flow.ID, dto.BusinessID)
}

func (uc *useCase) UpdateFlowGroup(ctx context.Context, dto dtos.UpdateFlowDTO) (*entities.Flow, error) {
	if uc.flowGroups == nil {
		return nil, fmt.Errorf("los flujos no estan disponibles")
	}

	current, err := uc.flowGroups.GetByID(ctx, dto.ID, dto.BusinessID)
	if err != nil {
		return nil, err
	}
	if current == nil {
		return nil, fmt.Errorf("flujo no encontrado")
	}

	name := strings.TrimSpace(dto.Name)
	if name == "" {
		return nil, fmt.Errorf("el nombre del flujo es obligatorio")
	}
	if len(name) > maxFlowNameLength {
		return nil, fmt.Errorf("el nombre del flujo supera %d caracteres", maxFlowNameLength)
	}

	taken, err := uc.flowGroups.ExistsByName(ctx, dto.BusinessID, name, dto.ID)
	if err != nil {
		return nil, err
	}
	if taken {
		return nil, fmt.Errorf("ya existe un flujo llamado %q", name)
	}

	if err := uc.assertOwnedTemplate(ctx, dto.RootTemplateID, dto.BusinessID); err != nil {
		return nil, err
	}

	current.Name = name
	current.Description = strings.TrimSpace(dto.Description)
	current.RootTemplateID = dto.RootTemplateID
	if dto.Enabled != nil {
		current.Enabled = *dto.Enabled
	}

	if err := uc.flowGroups.Update(ctx, current); err != nil {
		return nil, err
	}

	return uc.flowGroups.GetByID(ctx, dto.ID, dto.BusinessID)
}

func (uc *useCase) DeleteFlowGroup(ctx context.Context, id, businessID uint) error {
	if uc.flowGroups == nil {
		return fmt.Errorf("los flujos no estan disponibles")
	}

	flow, err := uc.flowGroups.GetByID(ctx, id, businessID)
	if err != nil {
		return err
	}
	if flow == nil {
		return fmt.Errorf("flujo no encontrado")
	}

	return uc.flowGroups.Delete(ctx, id, businessID)
}

func (uc *useCase) ListFlowTransitions(ctx context.Context, flowID, businessID uint) ([]entities.TemplateFlow, error) {
	if uc.flowGroups == nil || uc.flowRepository == nil {
		return nil, fmt.Errorf("los flujos no estan disponibles")
	}

	flow, err := uc.flowGroups.GetByID(ctx, flowID, businessID)
	if err != nil {
		return nil, err
	}
	if flow == nil {
		return nil, fmt.Errorf("flujo no encontrado")
	}

	return uc.flowRepository.ListByFlow(ctx, businessID, flowID)
}

func (uc *useCase) assertOwnedTemplate(ctx context.Context, templateID *uint, businessID uint) error {
	if templateID == nil || *templateID == 0 {
		return nil
	}

	template, err := uc.GetByID(ctx, *templateID, businessID)
	if err != nil {
		return err
	}
	if template == nil {
		return fmt.Errorf("la plantilla inicial del flujo no existe")
	}

	return nil
}
