package dtos

import "github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/entities"

// StateTransitionDTO representa el resultado de una transición de estado (DTO puro de dominio)
type StateTransitionDTO struct {
	NextState     entities.ConversationState
	TemplateName  string
	Variables     map[string]string
	PublishEvent  bool
	EventType     string
	EventMetadata map[string]interface{}
}
