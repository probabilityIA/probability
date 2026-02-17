package dtos

import "time"

// InvoiceResponseMessage es el mensaje que los proveedores publican de vuelta a Invoicing Module
type InvoiceResponseMessage struct {
	InvoiceID      uint                   `json:"invoice_id"`
	Provider       string                 `json:"provider"` // "softpymes", "siigo", "factus"
	Status         string                 `json:"status"`   // "success", "error"
	InvoiceNumber  string                 `json:"invoice_number,omitempty"`
	ExternalID     string                 `json:"external_id,omitempty"`
	InvoiceURL     string                 `json:"invoice_url,omitempty"`
	PDFURL         string                 `json:"pdf_url,omitempty"`
	XMLURL         string                 `json:"xml_url,omitempty"`
	CUFE           string                 `json:"cufe,omitempty"`
	IssuedAt       *time.Time             `json:"issued_at,omitempty"`
	DocumentJSON   map[string]interface{} `json:"document_json,omitempty"` // Documento completo del proveedor
	Error          string                 `json:"error,omitempty"`
	ErrorCode      string                 `json:"error_code,omitempty"`
	ErrorDetails   map[string]interface{} `json:"error_details,omitempty"`
	CorrelationID  string                 `json:"correlation_id"` // Mismo UUID del request
	Timestamp      time.Time              `json:"timestamp"`
	ProcessingTime int64                  `json:"processing_time_ms"` // Tiempo de procesamiento en ms

	// Audit data del request/response HTTP al proveedor
	AuditRequestURL     string                 `json:"audit_request_url,omitempty"`
	AuditRequestPayload map[string]interface{} `json:"audit_request_payload,omitempty"`
	AuditResponseStatus int                    `json:"audit_response_status,omitempty"`
	AuditResponseBody   string                 `json:"audit_response_body,omitempty"`
}

// Response statuses
const (
	ResponseStatusSuccess = "success"
	ResponseStatusError   = "error"
)
