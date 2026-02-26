package dtos

import "time"

// PaymentRequestMessage es el mensaje publicado a la cola pay.requests
type PaymentRequestMessage struct {
	PaymentTransactionID uint                   `json:"payment_transaction_id"`
	BusinessID           uint                   `json:"business_id"`
	GatewayCode          string                 `json:"gateway_code"` // "nequi"
	Amount               float64                `json:"amount"`
	Currency             string                 `json:"currency"`
	Reference            string                 `json:"reference"`
	PaymentMethod        string                 `json:"payment_method"`
	Description          string                 `json:"description"`
	Metadata             map[string]interface{} `json:"metadata,omitempty"`
	CorrelationID        string                 `json:"correlation_id"`
	Timestamp            time.Time              `json:"timestamp"`
}

// PaymentResponseMessage es el mensaje publicado a la cola pay.responses
type PaymentResponseMessage struct {
	PaymentTransactionID uint                   `json:"payment_transaction_id"`
	GatewayCode          string                 `json:"gateway_code"`
	Status               string                 `json:"status"` // "success"|"error"
	ExternalID           *string                `json:"external_id,omitempty"`
	GatewayResponse      map[string]interface{} `json:"gateway_response,omitempty"`
	Error                string                 `json:"error,omitempty"`
	ErrorCode            string                 `json:"error_code,omitempty"`
	CorrelationID        string                 `json:"correlation_id"`
	Timestamp            time.Time              `json:"timestamp"`
	ProcessingTimeMs     int64                  `json:"processing_time_ms"`
}

// PaginatedPaymentsDTO resultado paginado de listado de pagos
type PaginatedPaymentsDTO struct {
	Data       interface{}
	Total      int64
	Page       int
	PageSize   int
	TotalPages int
}
