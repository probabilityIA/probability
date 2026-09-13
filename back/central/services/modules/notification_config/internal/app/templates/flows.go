package templates

import (
	"context"
	"fmt"
	"strings"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/entities"
)

func (uc *useCase) ListFlows(ctx context.Context, sourceTemplateID, businessID uint) ([]entities.TemplateFlow, error) {
	if uc.flowRepository == nil {
		return nil, fmt.Errorf("los flujos de plantillas no estan disponibles")
	}

	source, err := uc.GetByID(ctx, sourceTemplateID, businessID)
	if err != nil {
		return nil, err
	}
	if source == nil {
		return nil, fmt.Errorf("plantilla no encontrada")
	}

	return uc.flowRepository.ListBySource(ctx, businessID, sourceTemplateID)
}

func (uc *useCase) ListBusinessFlows(ctx context.Context, businessID uint) ([]entities.TemplateFlow, error) {
	if uc.flowRepository == nil {
		return nil, fmt.Errorf("los flujos de plantillas no estan disponibles")
	}

	return uc.flowRepository.ListByBusiness(ctx, businessID)
}

func (uc *useCase) ReplaceFlows(ctx context.Context, dto dtos.ReplaceTemplateFlowsDTO) ([]entities.TemplateFlow, error) {
	if uc.flowRepository == nil {
		return nil, fmt.Errorf("los flujos de plantillas no estan disponibles")
	}

	source, err := uc.GetByID(ctx, dto.SourceTemplateID, dto.BusinessID)
	if err != nil {
		return nil, err
	}
	if source == nil {
		return nil, fmt.Errorf("plantilla no encontrada")
	}

	available := make(map[string]bool, len(source.Buttons))
	for _, button := range source.Buttons {
		if strings.EqualFold(strings.TrimSpace(button.Type), entities.TemplateButtonTypeURL) {
			continue
		}
		available[strings.ToLower(strings.TrimSpace(button.Text))] = true
	}

	seen := make(map[string]bool, len(dto.Flows))
	flows := make([]entities.TemplateFlow, 0, len(dto.Flows))

	for _, item := range dto.Flows {
		buttonText := strings.TrimSpace(item.ButtonText)
		if buttonText == "" {
			return nil, fmt.Errorf("una respuesta encadenada necesita el texto del boton")
		}

		key := strings.ToLower(buttonText)
		if !available[key] {
			return nil, fmt.Errorf("la plantilla no tiene un boton de respuesta llamado %q", buttonText)
		}
		if seen[key] {
			return nil, fmt.Errorf("el boton %q esta encadenado dos veces", buttonText)
		}
		seen[key] = true

		if item.TargetTemplateID == 0 {
			return nil, fmt.Errorf("el boton %q no tiene plantilla de respuesta", buttonText)
		}
		if item.TargetTemplateID == dto.SourceTemplateID {
			return nil, fmt.Errorf("el boton %q apunta a su propia plantilla", buttonText)
		}

		target, err := uc.GetByID(ctx, item.TargetTemplateID, dto.BusinessID)
		if err != nil {
			return nil, err
		}
		if target == nil {
			return nil, fmt.Errorf("la plantilla de respuesta del boton %q no existe", buttonText)
		}
		for _, variable := range target.Variables {
			if entities.IsFlowBlockedVariable(variable.Source) {
				return nil, fmt.Errorf(
					"la plantilla %s usa la variable %q, que no existe al responder un boton",
					target.Name,
					variable.Label,
				)
			}
		}

		enabled := true
		if item.Enabled != nil {
			enabled = *item.Enabled
		}

		flows = append(flows, entities.TemplateFlow{
			BusinessID:       dto.BusinessID,
			FlowID:           dto.FlowID,
			SourceTemplateID: dto.SourceTemplateID,
			ButtonText:       buttonText,
			TargetTemplateID: item.TargetTemplateID,
			Enabled:          enabled,
		})
	}

	if err := uc.validateGraph(ctx, dto.BusinessID, dto.SourceTemplateID, dto.FlowID, flows); err != nil {
		return nil, err
	}

	if err := uc.flowRepository.ReplaceForSource(ctx, dto.BusinessID, dto.SourceTemplateID, flows); err != nil {
		return nil, err
	}

	return uc.flowRepository.ListBySource(ctx, dto.BusinessID, dto.SourceTemplateID)
}

func (uc *useCase) validateGraph(
	ctx context.Context,
	businessID, sourceTemplateID uint,
	flowID *uint,
	incoming []entities.TemplateFlow,
) error {
	var existing []entities.TemplateFlow
	var err error

	if flowID != nil && *flowID > 0 {
		existing, err = uc.flowRepository.ListByFlow(ctx, businessID, *flowID)
	} else {
		existing, err = uc.flowRepository.ListByBusiness(ctx, businessID)
	}
	if err != nil {
		return err
	}

	edges := make(map[uint][]uint)
	for _, flow := range existing {
		if flow.SourceTemplateID == sourceTemplateID {
			continue
		}
		edges[flow.SourceTemplateID] = append(edges[flow.SourceTemplateID], flow.TargetTemplateID)
	}
	for _, flow := range incoming {
		edges[flow.SourceTemplateID] = append(edges[flow.SourceTemplateID], flow.TargetTemplateID)
	}

	onStack := make(map[uint]bool)

	var walk func(node uint, depth int) error
	walk = func(node uint, depth int) error {
		if depth > entities.MaxFlowDepth {
			return fmt.Errorf("el flujo encadena mas de %d plantillas seguidas", entities.MaxFlowDepth)
		}
		if onStack[node] {
			return fmt.Errorf("el flujo se cierra en un circulo y se repetiria sin fin")
		}

		onStack[node] = true
		for _, next := range edges[node] {
			if err := walk(next, depth+1); err != nil {
				return err
			}
		}
		onStack[node] = false

		return nil
	}

	for node := range edges {
		if err := walk(node, 1); err != nil {
			return err
		}
	}

	return nil
}

func (uc *useCase) resolveFlowParameters(
	ctx context.Context,
	event dtos.ButtonReplyEvent,
	targetTemplateID uint,
) []string {
	target, err := uc.repository.GetTemplateByID(ctx, targetTemplateID)
	if err != nil || target == nil || len(target.Variables) == 0 {
		return nil
	}

	candidate := entities.SegmentCandidate{BusinessID: event.BusinessID}

	if uc.segments != nil {
		found, err := uc.segments.FindCandidateByPhone(ctx, event.BusinessID, event.PhoneNumber)
		if err != nil {
			uc.logger.Warn().Err(err).
				Uint("business_id", event.BusinessID).
				Msg("No se pudo resolver el cliente por telefono: la respuesta usa los valores por defecto")
		} else if found != nil {
			candidate = *found
		}
	}

	return entities.BuildTemplateParameters(target, candidate)
}

func (uc *useCase) HandleButtonReply(ctx context.Context, event dtos.ButtonReplyEvent) error {
	if uc.flowRepository == nil || uc.flowPublisher == nil {
		return nil
	}
	if event.BusinessID == 0 || event.PhoneNumber == "" || event.ButtonText == "" {
		return nil
	}

	sourceName, err := uc.flowRepository.TemplateNameByOutboundMessageID(ctx, event.ContextMessageID)
	if err != nil {
		return err
	}
	if sourceName == "" {
		return nil
	}

	source, err := uc.repository.GetTemplateByName(ctx, event.BusinessID, sourceName, "")
	if err != nil {
		return err
	}
	if source == nil {
		return nil
	}

	flow, err := uc.flowRepository.Resolve(ctx, event.BusinessID, source.ID, event.ButtonText)
	if err != nil {
		return err
	}
	if flow == nil {
		return nil
	}

	if flow.TargetStatus != entities.TemplateStatusApproved {
		uc.logger.Warn().
			Uint("business_id", event.BusinessID).
			Str("button", event.ButtonText).
			Str("target", flow.TargetName).
			Str("status", flow.TargetStatus).
			Msg("El boton tiene respuesta encadenada pero Meta no aprobo la plantilla destino")
		return nil
	}

	message := dtos.FlowSendMessage{
		BusinessID:     event.BusinessID,
		Phone:          event.PhoneNumber,
		TemplateName:   flow.TargetName,
		Language:       flow.TargetLanguage,
		Parameters:     uc.resolveFlowParameters(ctx, event, flow.TargetTemplateID),
		HeaderImageURL: flow.TargetHeaderMediaURL,
	}

	if err := uc.flowPublisher.PublishFlowSend(ctx, message); err != nil {
		return err
	}

	uc.logger.Info().
		Uint("business_id", event.BusinessID).
		Str("button", event.ButtonText).
		Str("source", source.Name).
		Str("target", flow.TargetName).
		Msg("Respuesta encadenada del flujo encolada")

	return nil
}
