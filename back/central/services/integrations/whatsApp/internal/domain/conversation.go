package domain

import "time"

// ConversationState representa el estado actual de una conversación
type ConversationState string

const (
	StateStart                 ConversationState = "START"
	StateAwaitingConfirmation  ConversationState = "AWAITING_CONFIRMATION"
	StateAwaitingMenuSelection ConversationState = "AWAITING_MENU_SELECTION"
	StateAwaitingNoveltyType   ConversationState = "AWAITING_NOVELTY_TYPE"
	StateAwaitingCancelConfirm ConversationState = "AWAITING_CANCEL_CONFIRM"
	StateAwaitingCancelReason  ConversationState = "AWAITING_CANCEL_REASON"
	StateCompleted             ConversationState = "COMPLETED"
	StateHandoffToHuman        ConversationState = "HANDOFF_TO_HUMAN"
)

// Conversation representa una conversación activa de WhatsApp (entidad pura del dominio)
type Conversation struct {
	ID             string
	PhoneNumber    string
	OrderNumber    string
	BusinessID     uint
	CurrentState   ConversationState
	LastMessageID  string
	LastTemplateID string
	Metadata       map[string]interface{}
	CreatedAt      time.Time
	UpdatedAt      time.Time
	ExpiresAt      time.Time
}

// IsExpired verifica si la conversación ha expirado (ventana de 24h)
func (c *Conversation) IsExpired() bool {
	return time.Now().After(c.ExpiresAt)
}

// IsActive verifica si la conversación está activa y no ha expirado
func (c *Conversation) IsActive() bool {
	return !c.IsExpired() && c.CurrentState != StateCompleted && c.CurrentState != StateHandoffToHuman
}

// MessageDirection indica la dirección del mensaje
type MessageDirection string

const (
	MessageDirectionOutbound MessageDirection = "outbound" // Enviado por el sistema
	MessageDirectionInbound  MessageDirection = "inbound"  // Recibido del usuario
)

// MessageStatus representa el estado de un mensaje
type MessageStatus string

const (
	MessageStatusSent      MessageStatus = "sent"      // Enviado
	MessageStatusDelivered MessageStatus = "delivered" // Entregado al dispositivo
	MessageStatusRead      MessageStatus = "read"      // Leído por el usuario
	MessageStatusFailed    MessageStatus = "failed"    // Falló el envío
)

// MessageLog representa el registro de un mensaje en una conversación (entidad pura del dominio)
type MessageLog struct {
	ID             string
	ConversationID string
	Direction      MessageDirection
	MessageID      string // WhatsApp message ID
	TemplateName   string
	Content        string
	Status         MessageStatus
	DeliveredAt    *time.Time
	ReadAt         *time.Time
	CreatedAt      time.Time
}
