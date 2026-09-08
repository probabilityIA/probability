package entities

import "time"

const (
	SegmentCustomersInactive = "customers_inactive"
)

const (
	RunStatusRunning   = "running"
	RunStatusCompleted = "completed"
	RunStatusFailed    = "failed"
)

const (
	SendStatusQueued    = "queued"
	SendStatusSent      = "sent"
	SendStatusFailed    = "failed"
	SendStatusDiscarded = "discarded"
)

type SegmentParams struct {
	DaysWithoutPurchase int
	MinOrders           int
	MaxOrders           int
}

type ScheduledRule struct {
	ID            uint
	BusinessID    uint
	IntegrationID *uint

	NotificationTypeID uint
	WhatsappTemplateID *uint

	Name        string
	Description string

	SegmentType   string
	SegmentParams SegmentParams

	Timezone         string
	FrequencyMinutes uint
	SendWindowStart  string
	SendWindowEnd    string

	CooldownDays  uint
	DailySendCap  uint
	BatchSizeCap  uint
	RequiresOptIn bool

	Enabled   bool
	LastRunAt *time.Time
	NextRunAt *time.Time

	CreatedByID *uint
	CreatedAt   time.Time
	UpdatedAt   time.Time

	Template *WhatsappTemplate
}

type ScheduledRun struct {
	ID         uint
	RuleID     uint
	BusinessID uint

	StartedAt  time.Time
	FinishedAt *time.Time

	Status string

	MatchedCount   uint
	QueuedCount    uint
	SkippedCount   uint
	FailedCount    uint
	CapReachedFlag bool

	ErrorMessage string
}

type SegmentCandidate struct {
	ClientID     uint
	BusinessID   uint
	Name         string
	Phone        string
	DaysInactive int
	TotalOrders  int
	LastProduct  string
	BusinessName string
}

type ScheduledSend struct {
	ID           uint
	RuleID       uint
	ClientID     uint
	SendDate     string
	RunID        *uint
	BusinessID   uint
	Phone        string
	Status       string
	MessageID    string
	ErrorMessage string
	QueuedAt     time.Time
	SentAt       *time.Time
}

func IsSupportedSegment(segmentType string) bool {
	return segmentType == SegmentCustomersInactive
}
