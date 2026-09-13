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
	CampaignScheduleDaily    = "daily"
	CampaignScheduleInterval = "interval"
	CampaignScheduleDates    = "dates"
)

const (
	CampaignDeliveryDistribute = "distribute"
	CampaignDeliveryRepeat     = "repeat"
)

const (
	CampaignMaxDailySendCap = 1000
	CampaignMaxBatchSize    = 200
	CampaignMaxIntervalDays = 365
	CampaignMaxSendDates    = 120
	CampaignMaxOccurrences  = 365
)

type AudienceLocation struct {
	City    string
	State   string
	Clients uint
}

type CampaignAudienceParams struct {
	City             string
	State            string
	CreatedFromDays  int
	OnlyWithoutOrder bool
	ClientIDs        []uint

	ExcludeRecentDays int
	ExcludeRecentMax  int

	RegisteredBeforeDays   int
	MinOrders              int
	MinSpent               float64
	LastPurchaseBeforeDays int
}

type Campaign struct {
	ID            uint
	BusinessID    uint
	IntegrationID *uint

	WhatsappTemplateID *uint
	FlowID             *uint

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

	ScheduleMode string
	IntervalDays uint
	SendDates    []string
	DeliveryMode string
	Occurrences  uint
	CurrentRound uint

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
	Round          uint
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
	return c.Status == CampaignStatusDraft || c.Status == CampaignStatusScheduled || c.Status == CampaignStatusPaused
}

func (c *Campaign) IsRepeat() bool {
	return c.DeliveryMode == CampaignDeliveryRepeat
}

func (c *Campaign) IsFinished() bool {
	return c.Status == CampaignStatusCompleted || c.Status == CampaignStatusCancelled
}

func IsSupportedAudience(audienceType string) bool {
	return audienceType == CampaignAudienceAllClients || audienceType == CampaignAudienceFiltered
}
