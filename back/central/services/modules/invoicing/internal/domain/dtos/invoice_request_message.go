package dtos

import "time"

// InvoiceRequestMessage es el mensaje que Invoicing Module publica a los proveedores
type InvoiceRequestMessage struct {
	InvoiceID     uint        `json:"invoice_id"`
	Provider      string      `json:"provider"`       // "softpymes", "siigo", "factus"
	Operation     string      `json:"operation"`      // "create", "retry", "cancel"
	InvoiceData   InvoiceData `json:"invoice_data"`   // Datos tipados de factura
	CorrelationID string      `json:"correlation_id"` // UUID para correlacionar request/response
	Timestamp     time.Time   `json:"timestamp"`
}

// Operations
const (
	OperationCreate = "create"
	OperationRetry  = "retry"
	OperationCancel = "cancel"
)

// Providers
const (
	ProviderSoftpymes = "softpymes"
	ProviderSiigo     = "siigo"
	ProviderFactus    = "factus"
)
