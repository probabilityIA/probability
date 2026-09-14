package entities

import "time"

type MessageAuditLog struct {
	ID             string
	ConversationID string
	MessageID      string
	Direction      string
	TemplateName   string
	Content        string
	Status         string
	DeliveredAt    *time.Time
	ReadAt         *time.Time
	CreatedAt      time.Time

	PhoneNumber string
	OrderNumber string
	BusinessID  uint
}

type MessageAuditStats struct {
	TotalSent      int64
	TotalDelivered int64
	TotalRead      int64
	TotalFailed    int64
	SuccessRate    float64
}

type ConversationSummary struct {
	ID                   string
	PhoneNumber          string
	OrderNumber          string
	OrderID              string
	CampaignID           *uint
	CampaignName         string
	CustomerName         string
	UnreadCount          int
	OptedOut             bool
	ConversationType     string
	BusinessID           uint
	CurrentState         string
	MessageCount         int
	LastMessageContent   string
	LastMessageDirection string
	LastMessageStatus    string
	LastActivity         time.Time
	CreatedAt            time.Time
}

const ChatRetentionDays = 365

type ChatPurgeResult struct {
	Messages      int64
	Conversations int64
	Reads         int64
}

type MessageMedia struct {
	Type     string
	Key      string
	Mime     string
	Filename string
	Size     int64
}

type MessageButton struct {
	Text string
	Type string
}

type ConversationMessage struct {
	ID           string
	Direction    string
	MessageID    string
	TemplateName string
	Content      string
	Buttons      []MessageButton
	Media        *MessageMedia
	Status       string
	DeliveredAt  *time.Time
	ReadAt       *time.Time
	CreatedAt    time.Time
}
