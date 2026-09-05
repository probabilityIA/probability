package entities

import "time"

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

type ConversationType string

const (
	ConversationTypeOrder       ConversationType = "order"
	ConversationTypeSystemAlert ConversationType = "system_alert"
	ConversationTypeInbound     ConversationType = "inbound"
)

type Conversation struct {
	ID               string
	PhoneNumber      string
	OrderNumber      string
	ConversationType ConversationType
	BusinessID       uint
	CurrentState     ConversationState
	LastMessageID    string
	LastTemplateID   string
	Metadata         map[string]interface{}
	CreatedAt        time.Time
	UpdatedAt        time.Time
	ExpiresAt        time.Time
}

func (c *Conversation) IsExpired() bool {
	return time.Now().After(c.ExpiresAt)
}

func (c *Conversation) IsActive() bool {
	return !c.IsExpired() && c.CurrentState != StateCompleted && c.CurrentState != StateHandoffToHuman
}

type MessageDirection string

const (
	MessageDirectionOutbound MessageDirection = "outbound"
	MessageDirectionInbound  MessageDirection = "inbound"
)

type MessageStatus string

const (
	MessageStatusSent      MessageStatus = "sent"
	MessageStatusDelivered MessageStatus = "delivered"
	MessageStatusRead      MessageStatus = "read"
	MessageStatusFailed    MessageStatus = "failed"
)

type MessageLog struct {
	ID             string
	ConversationID string
	Direction      MessageDirection
	MessageID      string
	TemplateName   string
	Content        string
	Status         MessageStatus
	DeliveredAt    *time.Time
	ReadAt         *time.Time
	CreatedAt      time.Time
}
