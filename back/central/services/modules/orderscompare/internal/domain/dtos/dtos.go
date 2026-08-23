package dtos

import (
	"time"

	"github.com/secamc93/probability/back/central/shared/orderscompare"
)

type CompareQuery struct {
	BusinessID    uint
	IntegrationID uint
	From          *time.Time
	To            *time.Time
	Limit         int
	Page          int
	PageSize      int
	OnlyDiff      bool
	Search        string
}

type ComparePage struct {
	Rows       []orderscompare.Row   `json:"rows"`
	Totals     orderscompare.Totals  `json:"totals"`
	Page       int                   `json:"page"`
	PageSize   int                   `json:"page_size"`
	Total      int                   `json:"total"`
	TotalPages int                   `json:"total_pages"`
	CheckedAt  time.Time             `json:"checked_at"`
	Channel    ChannelInfo           `json:"channel"`
}

type ChannelInfo struct {
	IntegrationID   uint   `json:"integration_id"`
	IntegrationName string `json:"integration_name"`
	IntegrationType uint   `json:"integration_type_id"`
	Supported       bool   `json:"supported"`
}

type ApplyCommand struct {
	BusinessID    uint
	IntegrationID uint
	ExternalIDs   []string
}

type ApplyResult struct {
	Queued           []string          `json:"queued"`
	Skipped          []string          `json:"skipped"`
	Failed           map[string]string `json:"failed,omitempty"`
	WithoutInventory []string          `json:"without_inventory"`
	Note             string            `json:"note,omitempty"`
}

type LocalOrder = orderscompare.LocalOrder
