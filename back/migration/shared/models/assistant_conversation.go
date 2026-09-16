package models

import "time"

type AssistantConversation struct {
	ID string `gorm:"type:uuid;primaryKey"`

	BusinessID *uint     `gorm:"index"`
	Business   *Business `gorm:"foreignKey:BusinessID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	UserID uint `gorm:"not null;index"`
	User   User `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	MessageCount  int       `gorm:"not null;default:0"`
	StartedAt     time.Time `gorm:"not null"`
	LastMessageAt time.Time `gorm:"not null;index"`

	CreatedAt time.Time
	UpdatedAt time.Time

	Messages []AssistantMessage `gorm:"foreignKey:ConversationID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (AssistantConversation) TableName() string {
	return "assistant_conversations"
}

type AssistantMessage struct {
	ID string `gorm:"type:uuid;primaryKey"`

	ConversationID string `gorm:"type:uuid;not null;index"`

	BusinessID *uint     `gorm:"index"`
	Business   *Business `gorm:"foreignKey:BusinessID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	UserID uint `gorm:"not null;index"`
	User   User `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	Pathname string `gorm:"size:255"`
	Question string `gorm:"type:text;not null"`
	Answer   string `gorm:"type:text"`

	DestinationKey   string `gorm:"size:80;index"`
	DestinationRoute string `gorm:"size:255"`
	ErrorCode        string `gorm:"size:40;index"`

	Model        string `gorm:"size:80"`
	InputTokens  int    `gorm:"not null;default:0"`
	OutputTokens int    `gorm:"not null;default:0"`
	LatencyMs    int    `gorm:"not null;default:0"`

	Feedback   int16 `gorm:"not null;default:0;index"`
	FeedbackAt *time.Time
	ClickedAt  *time.Time

	CreatedAt time.Time `gorm:"not null;index"`
}

func (AssistantMessage) TableName() string {
	return "assistant_messages"
}
