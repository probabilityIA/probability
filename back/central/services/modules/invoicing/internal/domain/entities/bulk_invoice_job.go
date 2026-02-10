package entities

import "time"

// BulkInvoiceJob representa un job de facturación masiva (sin tags GORM)
type BulkInvoiceJob struct {
	ID              string
	BusinessID      uint
	CreatedByUserID *uint
	TotalOrders     int
	Processed       int
	Successful      int
	Failed          int
	Status          string
	StartedAt       *time.Time
	CompletedAt     *time.Time
	ErrorMessage    *string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// BulkInvoiceJobItem representa un item individual de un job
type BulkInvoiceJobItem struct {
	ID           uint
	JobID        string
	OrderID      string
	InvoiceID    *uint
	Status       string
	ErrorMessage *string
	ProcessedAt  *time.Time
	CreatedAt    time.Time
}

// Estados de Job
const (
	JobStatusPending    = "pending"
	JobStatusProcessing = "processing"
	JobStatusCompleted  = "completed"
	JobStatusFailed     = "failed"
)

// Estados de JobItem
const (
	JobItemStatusPending    = "pending"
	JobItemStatusProcessing = "processing"
	JobItemStatusSuccess    = "success"
	JobItemStatusFailed     = "failed"
)

// IsCompleted verifica si el job está completado
func (j *BulkInvoiceJob) IsCompleted() bool {
	return j.Processed >= j.TotalOrders
}

// GetProgress calcula el progreso del job (0-100)
func (j *BulkInvoiceJob) GetProgress() float64 {
	if j.TotalOrders == 0 {
		return 0
	}
	return (float64(j.Processed) / float64(j.TotalOrders)) * 100
}

// GetSuccessRate calcula la tasa de éxito (0-100)
func (j *BulkInvoiceJob) GetSuccessRate() float64 {
	if j.Processed == 0 {
		return 0
	}
	return (float64(j.Successful) / float64(j.Processed)) * 100
}
