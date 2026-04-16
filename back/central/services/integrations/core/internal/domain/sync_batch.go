package domain

import "time"

// SyncBatchMessage es el mensaje que se publica a la cola de lotes de sincronización.
// Cada mensaje representa un chunk de fecha que el consumer debe procesar.
type SyncBatchMessage struct {
	JobID             string    `json:"job_id"`
	IntegrationID     string    `json:"integration_id"`
	IntegrationTypeID int       `json:"integration_type_id"`
	BusinessID        *uint     `json:"business_id"`
	BatchIndex        int       `json:"batch_index"`
	TotalBatches      int       `json:"total_batches"`
	CreatedAtMin      time.Time `json:"created_at_min"`
	CreatedAtMax      time.Time `json:"created_at_max"`
	Status            string    `json:"status,omitempty"`
	FinancialStatus   string    `json:"financial_status,omitempty"`
	FulfillmentStatus string    `json:"fulfillment_status,omitempty"`
	EnqueuedAt        time.Time `json:"enqueued_at"`
}

// SyncBatchParams representa los parámetros para sincronización por lotes.
type SyncBatchParams struct {
	CreatedAtMin      *time.Time
	CreatedAtMax      *time.Time
	Status            string
	FinancialStatus   string
	FulfillmentStatus string
}

// ToGenericMap convierte los parámetros a map[string]interface{} para pasar al provider.
func (p *SyncBatchParams) ToGenericMap() map[string]interface{} {
	m := make(map[string]interface{})
	if p.CreatedAtMin != nil {
		m["created_at_min"] = *p.CreatedAtMin
	}
	if p.CreatedAtMax != nil {
		m["created_at_max"] = *p.CreatedAtMax
	}
	if p.Status != "" {
		m["status"] = p.Status
	}
	if p.FinancialStatus != "" {
		m["financial_status"] = p.FinancialStatus
	}
	if p.FulfillmentStatus != "" {
		m["fulfillment_status"] = p.FulfillmentStatus
	}
	return m
}
