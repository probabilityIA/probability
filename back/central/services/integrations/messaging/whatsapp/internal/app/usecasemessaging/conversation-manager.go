package usecasemessaging

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/errors"
)

func (u *usecases) GetInitialState() entities.ConversationState {
	return entities.StateAwaitingConfirmation
}

func (u *usecases) IsTerminalState(state entities.ConversationState) bool {
	return state == entities.StateCompleted || state == entities.StateHandoffToHuman
}

func (u *usecases) TransitionState(
	ctx context.Context,
	conversation *entities.Conversation,
	userResponse string,
) (*dtos.StateTransitionDTO, error) {
	u.log.Info(ctx).
		Str("conversation_id", conversation.ID).
		Str("current_state", string(conversation.CurrentState)).
		Str("user_response", userResponse).
		Msg("[Conversation Manager] - evaluando transición de estado")

	var transition *dtos.StateTransitionDTO
	var err error

	switch conversation.CurrentState {
	case entities.StateStart:
		transition = u.handleStartState(ctx, conversation)
	case entities.StateAwaitingConfirmation:
		transition = u.handleAwaitingConfirmationState(ctx, conversation, userResponse)
	case entities.StateAwaitingMenuSelection:
		transition = u.handleAwaitingMenuSelectionState(ctx, conversation, userResponse)
	case entities.StateAwaitingNoveltyType:
		transition = u.handleAwaitingNoveltyTypeState(ctx, conversation, userResponse)
	case entities.StateAwaitingCancelConfirm:
		transition = u.handleAwaitingCancelConfirmState(ctx, conversation, userResponse)
	case entities.StateAwaitingCancelReason:
		transition = u.handleAwaitingCancelReasonState(ctx, conversation, userResponse)
	default:
		u.log.Error(ctx).
			Str("state", string(conversation.CurrentState)).
			Msg("[Conversation Manager] - estado no reconocido")
		return nil, &errors.ErrInvalidStateTransition{
			CurrentState: string(conversation.CurrentState),
			UserResponse: userResponse,
		}
	}

	if transition == nil {
		u.log.Warn(ctx).
			Str("state", string(conversation.CurrentState)).
			Str("response", userResponse).
			Msg("[Conversation Manager] - transición no válida")
		return nil, &errors.ErrInvalidStateTransition{
			CurrentState: string(conversation.CurrentState),
			UserResponse: userResponse,
		}
	}

	u.log.Info(ctx).
		Str("current_state", string(conversation.CurrentState)).
		Str("next_state", string(transition.NextState)).
		Str("template", transition.TemplateName).
		Bool("publish_event", transition.PublishEvent).
		Msg("[Conversation Manager] - transición determinada")

	return transition, err
}

func (u *usecases) handleStartState(ctx context.Context, conversation *entities.Conversation) *dtos.StateTransitionDTO {
	getMetaStr := func(key string) string {
		if conversation.Metadata == nil {
			return ""
		}
		if v, ok := conversation.Metadata[key].(string); ok {
			return v
		}
		return ""
	}

	return &dtos.StateTransitionDTO{
		NextState:    entities.StateAwaitingConfirmation,
		TemplateName: "confirmacion_pedido_contraentrega",
		Variables: map[string]string{
			"1": getMetaStr("nombre"),
			"2": getMetaStr("tienda"),
			"3": conversation.OrderNumber,
			"4": getMetaStr("direccion"),
			"5": getMetaStr("productos"),
		},
		PublishEvent: false,
	}
}

func (u *usecases) handleAwaitingConfirmationState(
	ctx context.Context,
	conversation *entities.Conversation,
	userResponse string,
) *dtos.StateTransitionDTO {
	switch userResponse {
	case "Confirmar pedido":
		getMetaStr := func(key string) string {
			if conversation.Metadata == nil {
				return ""
			}
			if v, ok := conversation.Metadata[key].(string); ok {
				return v
			}
			return ""
		}

		return &dtos.StateTransitionDTO{
			NextState:    entities.StateCompleted,
			TemplateName: "pedido_confirmado_v2",
			Variables: map[string]string{
				"1": getMetaStr("nombre"),
				"2": conversation.OrderNumber,
				"3": getMetaStr("tienda"),
				"4": getMetaStr("direccion"),
				"5": getMetaStr("productos"),
			},
			PublishEvent: true,
			EventType:    "confirmed",
		}

	case "No confirmar":
		return &dtos.StateTransitionDTO{
			NextState:    entities.StateAwaitingMenuSelection,
			TemplateName: "menu_no_confirmacion",
			Variables: map[string]string{
				"1": conversation.OrderNumber,
			},
			PublishEvent: false,
		}

	default:
		return nil
	}
}

func (u *usecases) handleAwaitingMenuSelectionState(
	ctx context.Context,
	conversation *entities.Conversation,
	userResponse string,
) *dtos.StateTransitionDTO {
	switch userResponse {
	case "Presentar novedad":
		return &dtos.StateTransitionDTO{
			NextState:    entities.StateAwaitingNoveltyType,
			TemplateName: "tipo_novedad_pedido",
			Variables:    map[string]string{},
			PublishEvent: false,
		}

	case "Cancelar pedido":
		return &dtos.StateTransitionDTO{
			NextState:    entities.StateAwaitingCancelConfirm,
			TemplateName: "confirmar_cancelacion_pedido",
			Variables: map[string]string{
				"1": conversation.OrderNumber,
			},
			PublishEvent: false,
		}

	case "Asesor":
		return &dtos.StateTransitionDTO{
			NextState:    entities.StateHandoffToHuman,
			TemplateName: "handoff_asesor",
			Variables:    map[string]string{},
			PublishEvent: true,
			EventType:    "handoff",
		}

	default:
		return nil
	}
}

func (u *usecases) handleAwaitingNoveltyTypeState(
	ctx context.Context,
	conversation *entities.Conversation,
	userResponse string,
) *dtos.StateTransitionDTO {
	var templateName string
	var noveltyType string

	switch userResponse {
	case "Cambio de dirección":
		templateName = "novedad_cambio_direccion"
		noveltyType = "change_address"
	case "Cambio de productos":
		templateName = "novedad_cambio_productos"
		noveltyType = "change_products"
	case "Cambio medio de pago":
		templateName = "novedad_cambio_medio_pago"
		noveltyType = "change_payment"
	default:
		return nil
	}

	return &dtos.StateTransitionDTO{
		NextState:    entities.StateCompleted,
		TemplateName: templateName,
		Variables:    map[string]string{},
		PublishEvent: true,
		EventType:    "novelty",
		EventMetadata: map[string]interface{}{
			"novelty_type": noveltyType,
		},
	}
}

func (u *usecases) handleAwaitingCancelConfirmState(
	ctx context.Context,
	conversation *entities.Conversation,
	userResponse string,
) *dtos.StateTransitionDTO {
	switch userResponse {
	case "Sí, cancelar":
		return &dtos.StateTransitionDTO{
			NextState:    entities.StateCompleted,
			Variables:    map[string]string{},
			PublishEvent: true,
			EventType:    "cancelled",
		}

	case "No, volver":
		return &dtos.StateTransitionDTO{
			NextState:    entities.StateAwaitingMenuSelection,
			TemplateName: "menu_no_confirmacion",
			Variables: map[string]string{
				"1": conversation.OrderNumber,
			},
			PublishEvent: false,
		}

	default:
		return nil
	}
}

func (u *usecases) handleAwaitingCancelReasonState(
	ctx context.Context,
	conversation *entities.Conversation,
	userResponse string,
) *dtos.StateTransitionDTO {
	return &dtos.StateTransitionDTO{
		NextState:    entities.StateCompleted,
		Variables:    map[string]string{},
		PublishEvent: true,
		EventType:    "cancelled",
		EventMetadata: map[string]interface{}{
			"cancellation_reason": userResponse,
		},
	}
}

func (u *usecases) GetAvailableResponses(state entities.ConversationState) []string {
	responses := map[entities.ConversationState][]string{
		entities.StateAwaitingConfirmation: {
			"Confirmar pedido",
			"No confirmar",
		},
		entities.StateAwaitingMenuSelection: {
			"Presentar novedad",
			"Cancelar pedido",
			"Asesor",
		},
		entities.StateAwaitingNoveltyType: {
			"Cambio de dirección",
			"Cambio de productos",
			"Cambio medio de pago",
		},
		entities.StateAwaitingCancelConfirm: {
			"Sí, cancelar",
			"No, volver",
		},
		entities.StateAwaitingCancelReason: {
			"[Texto libre]",
		},
	}

	return responses[state]
}

func (u *usecases) ValidateTransition(
	ctx context.Context,
	currentState entities.ConversationState,
	nextState entities.ConversationState,
) error {
	validTransitions := map[entities.ConversationState][]entities.ConversationState{
		entities.StateStart: {
			entities.StateAwaitingConfirmation,
		},
		entities.StateAwaitingConfirmation: {
			entities.StateCompleted,
			entities.StateAwaitingMenuSelection,
		},
		entities.StateAwaitingMenuSelection: {
			entities.StateAwaitingNoveltyType,
			entities.StateAwaitingCancelConfirm,
			entities.StateHandoffToHuman,
		},
		entities.StateAwaitingNoveltyType: {
			entities.StateCompleted,
		},
		entities.StateAwaitingCancelConfirm: {
			entities.StateAwaitingCancelReason,
			entities.StateAwaitingMenuSelection,
		},
		entities.StateAwaitingCancelReason: {
			entities.StateCompleted,
		},
	}

	allowed, exists := validTransitions[currentState]
	if !exists {
		return fmt.Errorf("estado actual no reconocido: %s", currentState)
	}

	for _, validNext := range allowed {
		if validNext == nextState {
			return nil
		}
	}

	return fmt.Errorf("transición no permitida de %s a %s", currentState, nextState)
}
