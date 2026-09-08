package entities

import "time"

const (
	CampaignStatusDraft     = "draft"
	CampaignStatusScheduled = "scheduled"
	CampaignStatusRunning   = "running"
	CampaignStatusPaused    = "paused"
	CampaignStatusCompleted = "completed"
	CampaignStatusCancelled = "cancelled"
)

const (
	CampaignSendStatusPending   = "pending"
	CampaignSendStatusQueued    = "queued"
	CampaignSendStatusSent      = "sent"
	CampaignSendStatusDelivered = "delivered"
	CampaignSendStatusRead      = "read"
	CampaignSendStatusReplied   = "replied"
	CampaignSendStatusFailed    = "failed"
	CampaignSendStatusSkipped   = "skipped"
)

const (
	CampaignAudienceAllClients = "all_clients"
	CampaignAudienceFiltered   = "filtered_clients"
)

const (
	CampaignMaxDailySendCap = 1000
	CampaignMaxBatchSize    = 200
)

type CampaignAudienceParams struct {
	City             string
	CreatedFromDays  int
	OnlyWithoutOrder bool
	ClientIDs        []uint
}

type Campaign struct {
	ID            uint
	BusinessID    uint
	IntegrationID *uint

	WhatsappTemplateID *uint

	Name        string
	Description string
	SenderName  string

	AudienceType   string
	AudienceParams CampaignAudienceParams
	VariableValues map[string]string

	Timezone        string
	SendWindowStart string
	SendWindowEnd   string

	ScheduledAt  *time.Time
	DailySendCap uint
	BatchSize    uint

	Status string

	AudienceCount uint
	QueuedCount   uint
	SentCount     uint
	FailedCount   uint
	SkippedCount  uint
	RepliedCount  uint

	StartedAt   *time.Time
	FinishedAt  *time.Time
	LastBatchAt *time.Time

	ErrorMessage string

	CreatedByID *uint
	CreatedAt   time.Time
	UpdatedAt   time.Time

	Template *WhatsappTemplate
}

type CampaignCandidate struct {
	ClientID   uint
	BusinessID uint
	Name       string
	FirstName  string
	Phone      string
	City       string
}

type CampaignSend struct {
	ID             uint
	CampaignID     uint
	ClientID       uint
	BusinessID     uint
	Phone          string
	ClientName     string
	ConversationID *string
	Status         string
	MessageID      string
	ErrorMessage   string
	QueuedAt       *time.Time
	SentAt         *time.Time
	DeliveredAt    *time.Time
	ReadAt         *time.Time
	RepliedAt      *time.Time
}

func (c *Campaign) IsEditable() bool {
	return c.Status == CampaignStatusDraft || c.Status == CampaignStatusScheduled
}

func (c *Campaign) IsFinished() bool {
	return c.Status == CampaignStatusCompleted || c.Status == CampaignStatusCancelled
}

func IsSupportedAudience(audienceType string) bool {
	return audienceType == CampaignAudienceAllClients || audienceType == CampaignAudienceFiltered
}
