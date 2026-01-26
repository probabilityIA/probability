package app

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/whatsApp/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
)

// IConversationManager define la interfaz para el gestor de conversaciones
type IConversationManager interface {
	TransitionState(ctx context.Context, conversation *domain.Conversation, userResponse string) (*StateTransition, error)
	GetInitialState() domain.ConversationState
	IsTerminalState(state domain.ConversationState) bool
}

// StateTransition representa el resultado de una transición de estado
type StateTransition struct {
	NextState     domain.ConversationState
	TemplateName  string
	Variables     map[string]string
	PublishEvent  bool
	EventType     string
	EventMetadata map[string]interface{}
}

// ConversationManager implementa la máquina de estados conversacional
type ConversationManager struct {
	conversationRepo domain.IConversationRepository
	log              log.ILogger
}

// NewConversationManager crea una nueva instancia del gestor
func NewConversationManager(
	conversationRepo domain.IConversationRepository,
	logger log.ILogger,
) IConversationManager {
	return &ConversationManager{
		conversationRepo: conversationRepo,
		log:              logger.WithModule("whatsapp-conversation-manager"),
	}
}

// GetInitialState retorna el estado inicial de una conversación
func (m *ConversationManager) GetInitialState() domain.ConversationState {
	return domain.StateAwaitingConfirmation
}

// IsTerminalState verifica si un estado es terminal (conversación finalizada)
func (m *ConversationManager) IsTerminalState(state domain.ConversationState) bool {
	return state == domain.StateCompleted || state == domain.StateHandoffToHuman
}

// TransitionState evalúa la respuesta del usuario y determina la siguiente transición
func (m *ConversationManager) TransitionState(
	ctx context.Context,
	conversation *domain.Conversation,
	userResponse string,
) (*StateTransition, error) {
	m.log.Info(ctx).
		Str("conversation_id", conversation.ID).
		Str("current_state", string(conversation.CurrentState)).
		Str("user_response", userResponse).
		Msg("[Conversation Manager] - evaluando transición de estado")

	var transition *StateTransition
	var err error

	// Evaluar según el estado actual
	switch conversation.CurrentState {
	case domain.StateStart:
		transition = m.handleStartState(ctx, conversation)
	case domain.StateAwaitingConfirmation:
		transition = m.handleAwaitingConfirmationState(ctx, conversation, userResponse)
	case domain.StateAwaitingMenuSelection:
		transition = m.handleAwaitingMenuSelectionState(ctx, conversation, userResponse)
	case domain.StateAwaitingNoveltyType:
		transition = m.handleAwaitingNoveltyTypeState(ctx, conversation, userResponse)
	case domain.StateAwaitingCancelConfirm:
		transition = m.handleAwaitingCancelConfirmState(ctx, conversation, userResponse)
	case domain.StateAwaitingCancelReason:
		transition = m.handleAwaitingCancelReasonState(ctx, conversation, userResponse)
	default:
		m.log.Error(ctx).
			Str("state", string(conversation.CurrentState)).
			Msg("[Conversation Manager] - estado no reconocido")
		return nil, &domain.ErrInvalidStateTransition{
			CurrentState: conversation.CurrentState,
			UserResponse: userResponse,
		}
	}

	if transition == nil {
		m.log.Warn(ctx).
			Str("state", string(conversation.CurrentState)).
			Str("response", userResponse).
			Msg("[Conversation Manager] - transición no válida")
		return nil, &domain.ErrInvalidStateTransition{
			CurrentState: conversation.CurrentState,
			UserResponse: userResponse,
		}
	}

	m.log.Info(ctx).
		Str("current_state", string(conversation.CurrentState)).
		Str("next_state", string(transition.NextState)).
		Str("template", transition.TemplateName).
		Bool("publish_event", transition.PublishEvent).
		Msg("[Conversation Manager] - transición determinada")

	return transition, err
}

// handleStartState maneja el estado inicial (envío de primera plantilla)
func (m *ConversationManager) handleStartState(ctx context.Context, conversation *domain.Conversation) *StateTransition {
	// Este estado se usa cuando se crea una conversación nueva
	// Normalmente se envía la plantilla inicial desde SendTemplate directamente
	return &StateTransition{
		NextState:    domain.StateAwaitingConfirmation,
		TemplateName: "confirmacion_pedido_contraentrega",
		Variables: map[string]string{
			// Las variables deben venir del contexto de creación
			"1": conversation.Metadata["nombre"].(string),
			"2": conversation.Metadata["tienda"].(string),
			"3": conversation.OrderNumber,
			"4": conversation.Metadata["direccion"].(string),
			"5": conversation.Metadata["productos"].(string),
		},
		PublishEvent: false,
	}
}

// handleAwaitingConfirmationState maneja respuestas en estado AWAITING_CONFIRMATION
func (m *ConversationManager) handleAwaitingConfirmationState(
	ctx context.Context,
	conversation *domain.Conversation,
	userResponse string,
) *StateTransition {
	switch userResponse {
	case "Confirmar pedido":
		return &StateTransition{
			NextState:    domain.StateCompleted,
			TemplateName: "pedido_confirmado",
			Variables: map[string]string{
				"1": conversation.OrderNumber,
			},
			PublishEvent: true,
			EventType:    "confirmed",
		}

	case "No confirmar":
		return &StateTransition{
			NextState:    domain.StateAwaitingMenuSelection,
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

// handleAwaitingMenuSelectionState maneja respuestas en estado AWAITING_MENU_SELECTION
func (m *ConversationManager) handleAwaitingMenuSelectionState(
	ctx context.Context,
	conversation *domain.Conversation,
	userResponse string,
) *StateTransition {
	switch userResponse {
	case "Presentar novedad":
		return &StateTransition{
			NextState:    domain.StateAwaitingNoveltyType,
			TemplateName: "tipo_novedad_pedido",
			Variables:    map[string]string{},
			PublishEvent: false,
		}

	case "Cancelar pedido":
		return &StateTransition{
			NextState:    domain.StateAwaitingCancelConfirm,
			TemplateName: "confirmar_cancelacion_pedido",
			Variables: map[string]string{
				"1": conversation.OrderNumber,
			},
			PublishEvent: false,
		}

	case "Asesor":
		return &StateTransition{
			NextState:    domain.StateHandoffToHuman,
			TemplateName: "handoff_asesor",
			Variables:    map[string]string{},
			PublishEvent: true,
			EventType:    "handoff",
		}

	default:
		return nil
	}
}

// handleAwaitingNoveltyTypeState maneja respuestas en estado AWAITING_NOVELTY_TYPE
func (m *ConversationManager) handleAwaitingNoveltyTypeState(
	ctx context.Context,
	conversation *domain.Conversation,
	userResponse string,
) *StateTransition {
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

	return &StateTransition{
		NextState:    domain.StateCompleted,
		TemplateName: templateName,
		Variables:    map[string]string{},
		PublishEvent: true,
		EventType:    "novelty",
		EventMetadata: map[string]interface{}{
			"novelty_type": noveltyType,
		},
	}
}

// handleAwaitingCancelConfirmState maneja respuestas en estado AWAITING_CANCEL_CONFIRM
func (m *ConversationManager) handleAwaitingCancelConfirmState(
	ctx context.Context,
	conversation *domain.Conversation,
	userResponse string,
) *StateTransition {
	switch userResponse {
	case "Sí, cancelar":
		return &StateTransition{
			NextState:    domain.StateAwaitingCancelReason,
			TemplateName: "motivo_cancelacion_pedido",
			Variables:    map[string]string{},
			PublishEvent: false,
		}

	case "No, volver":
		return &StateTransition{
			NextState:    domain.StateAwaitingMenuSelection,
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

// handleAwaitingCancelReasonState maneja respuestas en estado AWAITING_CANCEL_REASON
func (m *ConversationManager) handleAwaitingCancelReasonState(
	ctx context.Context,
	conversation *domain.Conversation,
	userResponse string,
) *StateTransition {
	// El usuario envió texto libre con el motivo
	// Cualquier texto es válido aquí
	return &StateTransition{
		NextState:    domain.StateCompleted,
		TemplateName: "pedido_cancelado",
		Variables: map[string]string{
			"1": conversation.OrderNumber,
		},
		PublishEvent: true,
		EventType:    "cancelled",
		EventMetadata: map[string]interface{}{
			"cancellation_reason": userResponse,
		},
	}
}

// GetAvailableResponses retorna las respuestas válidas para un estado dado
func (m *ConversationManager) GetAvailableResponses(state domain.ConversationState) []string {
	responses := map[domain.ConversationState][]string{
		domain.StateAwaitingConfirmation: {
			"Confirmar pedido",
			"No confirmar",
		},
		domain.StateAwaitingMenuSelection: {
			"Presentar novedad",
			"Cancelar pedido",
			"Asesor",
		},
		domain.StateAwaitingNoveltyType: {
			"Cambio de dirección",
			"Cambio de productos",
			"Cambio medio de pago",
		},
		domain.StateAwaitingCancelConfirm: {
			"Sí, cancelar",
			"No, volver",
		},
		domain.StateAwaitingCancelReason: {
			"[Texto libre]",
		},
	}

	return responses[state]
}

// ValidateTransition verifica si una transición es válida antes de ejecutarla
func (m *ConversationManager) ValidateTransition(
	ctx context.Context,
	currentState domain.ConversationState,
	nextState domain.ConversationState,
) error {
	// Definir transiciones válidas
	validTransitions := map[domain.ConversationState][]domain.ConversationState{
		domain.StateStart: {
			domain.StateAwaitingConfirmation,
		},
		domain.StateAwaitingConfirmation: {
			domain.StateCompleted,
			domain.StateAwaitingMenuSelection,
		},
		domain.StateAwaitingMenuSelection: {
			domain.StateAwaitingNoveltyType,
			domain.StateAwaitingCancelConfirm,
			domain.StateHandoffToHuman,
		},
		domain.StateAwaitingNoveltyType: {
			domain.StateCompleted,
		},
		domain.StateAwaitingCancelConfirm: {
			domain.StateAwaitingCancelReason,
			domain.StateAwaitingMenuSelection,
		},
		domain.StateAwaitingCancelReason: {
			domain.StateCompleted,
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
