package domain

import (
	"fmt"
	"strings"
	"time"
)

// SyncOrdersParams representa los parámetros para sincronizar órdenes
type SyncOrdersParams struct {
	CreatedAtMin      *time.Time
	CreatedAtMax      *time.Time
	Status            string
	FinancialStatus   string
	FulfillmentStatus string
}

type GetOrdersParams struct {
	Status            string
	Limit             int
	SinceID           *int64
	CreatedAtMin      *time.Time
	CreatedAtMax      *time.Time
	UpdatedAtMin      *time.Time
	UpdatedAtMax      *time.Time
	ProcessedAtMin    *time.Time
	ProcessedAtMax    *time.Time
	Fields            []string
	FinancialStatus   string
	FulfillmentStatus string
}

func (p *GetOrdersParams) ToQueryString() map[string]string {
	params := make(map[string]string)

	if p.Status != "" {
		params["status"] = p.Status
	} else {
		params["status"] = "any"
	}

	if p.Limit > 0 {
		if p.Limit > 250 {
			params["limit"] = "250"
		} else {
			params["limit"] = fmt.Sprintf("%d", p.Limit)
		}
	} else {
		params["limit"] = "250"
	}

	if p.SinceID != nil {
		params["since_id"] = fmt.Sprintf("%d", *p.SinceID)
	}

	if p.CreatedAtMin != nil {
		params["created_at_min"] = p.CreatedAtMin.Format(time.RFC3339)
	}
	if p.CreatedAtMax != nil {
		params["created_at_max"] = p.CreatedAtMax.Format(time.RFC3339)
	}
	if p.UpdatedAtMin != nil {
		params["updated_at_min"] = p.UpdatedAtMin.Format(time.RFC3339)
	}
	if p.UpdatedAtMax != nil {
		params["updated_at_max"] = p.UpdatedAtMax.Format(time.RFC3339)
	}
	if p.ProcessedAtMin != nil {
		params["processed_at_min"] = p.ProcessedAtMin.Format(time.RFC3339)
	}
	if p.ProcessedAtMax != nil {
		params["processed_at_max"] = p.ProcessedAtMax.Format(time.RFC3339)
	}

	if len(p.Fields) > 0 {
		params["fields"] = strings.Join(p.Fields, ",")
	}

	if p.FinancialStatus != "" {
		params["financial_status"] = p.FinancialStatus
	}

	if p.FulfillmentStatus != "" {
		params["fulfillment_status"] = p.FulfillmentStatus
	}

	return params
}
