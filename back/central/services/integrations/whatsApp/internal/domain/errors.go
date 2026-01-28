package domain

import "fmt"

// ErrTemplateNotFound se retorna cuando se solicita una plantilla que no existe
type ErrTemplateNotFound struct {
	TemplateName string
}

func (e *ErrTemplateNotFound) Error() string {
	return fmt.Sprintf("plantilla no encontrada: %s", e.TemplateName)
}

// ErrMissingVariable se retorna cuando falta una variable requerida
type ErrMissingVariable struct {
	TemplateName string
	VariableName string
	VariableKey  string
}

func (e *ErrMissingVariable) Error() string {
	return fmt.Sprintf("variable requerida faltante para plantilla '%s': %s ({{%s}})",
		e.TemplateName, e.VariableName, e.VariableKey)
}

// ErrConversationNotFound se retorna cuando no se encuentra una conversación
type ErrConversationNotFound struct {
	PhoneNumber string
	OrderNumber string
}

func (e *ErrConversationNotFound) Error() string {
	return fmt.Sprintf("conversación no encontrada para phone=%s, order=%s",
		e.PhoneNumber, e.OrderNumber)
}

// ErrConversationExpired se retorna cuando una conversación ha expirado (>24h)
type ErrConversationExpired struct {
	ConversationID string
}

func (e *ErrConversationExpired) Error() string {
	return fmt.Sprintf("conversación expirada: %s", e.ConversationID)
}

// ErrInvalidStateTransition se retorna cuando se intenta una transición de estado inválida
type ErrInvalidStateTransition struct {
	CurrentState ConversationState
	UserResponse string
}

func (e *ErrInvalidStateTransition) Error() string {
	return fmt.Sprintf("transición de estado inválida desde '%s' con respuesta '%s'",
		e.CurrentState, e.UserResponse)
}

// ErrWebhookSignatureInvalid se retorna cuando la firma del webhook es inválida
type ErrWebhookSignatureInvalid struct {
	Message string
}

func (e *ErrWebhookSignatureInvalid) Error() string {
	return fmt.Sprintf("firma de webhook inválida: %s", e.Message)
}
