package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type WhatsAppConversation struct {
	ID               uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	PhoneNumber      string         `gorm:"type:varchar(20);not null;index:idx_whatsapp_phone_order,priority:1"`
	OrderNumber      string         `gorm:"type:varchar(100);index:idx_whatsapp_phone_order,priority:2"`
	ConversationType string         `gorm:"type:varchar(20);not null;default:'order';check:conversation_type IN ('order','system_alert')"`
	BusinessID       uint           `gorm:"not null;index:idx_whatsapp_business_id"`
	CurrentState     string         `gorm:"type:varchar(50);not null;index:idx_whatsapp_current_state"`
	LastMessageID    string         `gorm:"type:varchar(255)"`
	LastTemplateID   string         `gorm:"type:varchar(100)"`
	Metadata         datatypes.JSON `gorm:"type:jsonb;default:'{}'"`
	CreatedAt        time.Time      `gorm:"not null;default:now()"`
	UpdatedAt        time.Time      `gorm:"not null;default:now()"`
	ExpiresAt        time.Time      `gorm:"not null;index:idx_whatsapp_expires_at"`

	MessageLogs []WhatsAppMessageLog `gorm:"foreignKey:ConversationID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (WhatsAppConversation) TableName() string {
	return "whatsapp_conversations"
}

type WhatsAppMessageLog struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ConversationID uuid.UUID `gorm:"type:uuid;not null;index:idx_whatsapp_msg_conversation"`
	Direction      string    `gorm:"type:varchar(10);not null;check:direction IN ('outbound', 'inbound')"`
	MessageID      string    `gorm:"type:varchar(255);not null;uniqueIndex:idx_whatsapp_msg_id"`
	TemplateName   string    `gorm:"type:varchar(100)"`
	Content        string    `gorm:"type:text"`
	Status         string    `gorm:"type:varchar(20);not null;index:idx_whatsapp_msg_status;check:status IN ('sent', 'delivered', 'read', 'failed')"`
	DeliveredAt    *time.Time
	ReadAt         *time.Time
	CreatedAt      time.Time `gorm:"not null;default:now();index:idx_whatsapp_msg_created_at"`

	Conversation WhatsAppConversation `gorm:"foreignKey:ConversationID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (WhatsAppMessageLog) TableName() string {
	return "whatsapp_message_logs"
}
