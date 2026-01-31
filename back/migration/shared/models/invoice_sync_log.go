package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// ───────────────────────────────────────────
//
//	INVOICE SYNC LOGS - Logs de sincronización de facturas
//
// ───────────────────────────────────────────

// InvoiceSyncLog registra cada intento de sincronización con el proveedor
// Útil para debugging, reintentos y auditoría
type InvoiceSyncLog struct {
	gorm.Model

	// Relación con Invoice
	InvoiceID uint    `gorm:"not null;index"`
	Invoice   Invoice `gorm:"foreignKey:InvoiceID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	// Tipo de operación
	// "create" = crear factura
	// "cancel" = cancelar factura
	// "credit_note" = crear nota de crédito
	// "query" = consultar estado
	OperationType string `gorm:"size:64;not null;index"`

	// Estado del intento
	// "pending" = pendiente
	// "processing" = en proceso
	// "success" = exitoso
	// "failed" = fallido
	Status string `gorm:"size:64;not null;index"`

	// Detalles del request
	RequestPayload  datatypes.JSON `gorm:"type:jsonb"` // Payload enviado al proveedor
	RequestHeaders  datatypes.JSON `gorm:"type:jsonb"` // Headers del request
	RequestURL      string         `gorm:"size:512"`   // URL del endpoint

	// Detalles del response
	ResponseStatus  int            `gorm:"index"`      // HTTP status code
	ResponseBody    datatypes.JSON `gorm:"type:jsonb"` // Cuerpo de la respuesta
	ResponseHeaders datatypes.JSON `gorm:"type:jsonb"` // Headers de la respuesta

	// Información de error (si aplica)
	ErrorMessage *string `gorm:"type:text"` // Mensaje de error
	ErrorCode    *string `gorm:"size:64"`   // Código de error del proveedor
	ErrorDetails datatypes.JSON `gorm:"type:jsonb"` // Detalles adicionales del error

	// Información de reintentos
	RetryCount    int        `gorm:"default:0;index"` // Número de reintentos realizados
	NextRetryAt   *time.Time `gorm:"index"`           // Cuándo se debe reintentar
	MaxRetries    int        `gorm:"default:3"`       // Máximo de reintentos permitidos
	RetriedAt     *time.Time                          // Cuándo se reintentó por última vez

	// Timestamps
	StartedAt   time.Time  `gorm:"not null;index"` // Cuándo comenzó el intento
	CompletedAt *time.Time `gorm:"index"`          // Cuándo completó (éxito o fallo)
	Duration    *int       // Duración en milisegundos

	// Metadata
	TriggeredBy string `gorm:"size:64"` // "auto", "manual", "retry_job"
	UserID      *uint  `gorm:"index"`   // ID del usuario que triggereó (si es manual)
}

// TableName especifica el nombre de la tabla para InvoiceSyncLog
func (InvoiceSyncLog) TableName() string {
	return "invoice_sync_logs"
}
