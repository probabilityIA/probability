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

const (
	DefaultPageSize = 50
	MaxPageSize     = 200
)

type DetailItem struct {
	SKU          string `json:"sku"`
	Label        string `json:"label"`
	Tone         string `json:"tone"`
	Group        string `json:"group,omitempty"`
	ParentRef    string `json:"parent_ref,omitempty"`
	ParentLabel  string `json:"parent_label,omitempty"`
	VariantLabel string `json:"variant_label,omitempty"`

	CounterpartSKU  string `json:"counterpart_sku,omitempty"`
	CounterpartName string `json:"counterpart_name,omitempty"`
	ChannelQty      *int   `json:"channel_qty,omitempty"`
	OwnQty          *int   `json:"own_qty,omitempty"`
	FixSide         string `json:"fix_side,omitempty"`
	Pattern         string `json:"pattern,omitempty"`
}

type DetailQuery struct {
	BusinessID    uint
	IntegrationID uint
	Kind          string
	Group         string
	Search        string
	Page          int
	PageSize      int
}

func (q *DetailQuery) Normalize() {
	if q.Kind != KindProducts {
		q.Kind = KindInventory
	}
	if q.Page < 1 {
		q.Page = 1
	}
	if q.PageSize < 1 {
		q.PageSize = DefaultPageSize
	}
	if q.PageSize > MaxPageSize {
		q.PageSize = MaxPageSize
	}
}

func (q *DetailQuery) Offset() int {
	return (q.Page - 1) * q.PageSize
}

type DetailPage struct {
	Items      []DetailItem
	Total      int64
	Page       int
	PageSize   int
	TotalPages int
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
	ChannelNoSKU      int
	SKUChanged        int
	SKUTypo           int
	SKUSpacing        int

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
	if len(r.Message) > 500 {
		r.Message = r.Message[:500]
	}
}
