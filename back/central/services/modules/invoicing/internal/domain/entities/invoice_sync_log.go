package entities

import "time"

// InvoiceSyncLog registra cada intento de sincronización con el proveedor
// Entidad PURA de dominio - SIN TAGS de infraestructura
type InvoiceSyncLog struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time

	// Relaciones (solo IDs)
	InvoiceID uint

	// Tipo de operación
	OperationType string

	// Estado del intento
	Status string

	// Detalles del request
	RequestPayload  map[string]interface{}
	RequestHeaders  map[string]interface{}
	RequestURL      string

	// Detalles del response
	ResponseStatus  int
	ResponseBody    map[string]interface{}
	ResponseHeaders map[string]interface{}

	// Detalles del recibo de caja (request/response separados de la factura)
	CashReceiptRequestURL     string
	CashReceiptRequestPayload map[string]interface{}
	CashReceiptResponseStatus int
	CashReceiptResponseBody   map[string]interface{}

	// Información de error
	ErrorMessage *string
	ErrorCode    *string
	ErrorDetails map[string]interface{}

	// Información de reintentos
	RetryCount  int
	NextRetryAt *time.Time
	MaxRetries  int
	RetriedAt   *time.Time

	// Timestamps
	StartedAt   time.Time
	CompletedAt *time.Time
	Duration    *int // Milisegundos

	// Metadata
	TriggeredBy string
	UserID      *uint
}
