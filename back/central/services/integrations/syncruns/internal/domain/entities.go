package domain

import "time"

const (
	KindInventory = "inventory"
	KindProducts  = "products"
)

const (
	StatusCompleted = "completed"
	StatusFailed    = "failed"
)

const MaxDetailItems = 200

type DetailItem struct {
	SKU   string `json:"sku"`
	Label string `json:"label"`
	Tone  string `json:"tone"`
	Group string `json:"group,omitempty"`
}

type SyncRun struct {
	ID            uint
	BusinessID    uint
	IntegrationID uint
	Kind          string
	CorrelationID string
	StartedAt     time.Time
	FinishedAt    *time.Time

	Total     int
	Updated   int
	Unchanged int
	Skipped   int
	Failed    int

	Matched           int
	NotAssociated     int
	OnlyInProbability int
	OnlyInChannel     int

	Status  string
	Message string
	Detail  []DetailItem
}

func (r *SyncRun) Normalize() {
	if r.Kind != KindProducts {
		r.Kind = KindInventory
	}
	if r.Status != StatusFailed {
		r.Status = StatusCompleted
	}
	if r.StartedAt.IsZero() {
		r.StartedAt = time.Now()
	}
	if r.FinishedAt == nil {
		now := time.Now()
		r.FinishedAt = &now
	}
	if len(r.Detail) > MaxDetailItems {
		r.Detail = r.Detail[:MaxDetailItems]
	}
	if len(r.Message) > 500 {
		r.Message = r.Message[:500]
	}
}
