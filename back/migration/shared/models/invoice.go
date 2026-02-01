package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// ───────────────────────────────────────────
//
//	INVOICES - Facturas generadas
//
// ───────────────────────────────────────────

// Invoice representa una factura electrónica generada para una orden
type Invoice struct {
	gorm.Model

	// Relación con Order
	OrderID string `gorm:"type:varchar(36);not null;index;uniqueIndex:idx_order_provider,priority:1"`
	Order   Order  `gorm:"foreignKey:OrderID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	// Relación con Business
	BusinessID uint     `gorm:"not null;index"`
	Business   Business `gorm:"foreignKey:BusinessID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	// Relación con InvoicingProvider (DEPRECATED - mantener temporalmente para dual-read)
	InvoicingProviderID *uint             `gorm:"index"`
	InvoicingProvider   InvoicingProvider `gorm:"foreignKey:InvoicingProviderID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	// Relación con Integration (nuevo - provider de facturación desde integrations/)
	InvoicingIntegrationID *uint       `gorm:"index"`
	InvoicingIntegration   Integration `gorm:"foreignKey:InvoicingIntegrationID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	// Identificadores
	InvoiceNumber string  `gorm:"size:128;not null;index"` // Número de factura del proveedor
	ExternalID    *string `gorm:"size:255;index"`          // ID en el sistema del proveedor
	InternalNumber string `gorm:"size:128;unique;index"`   // Número interno Probability (auto-generado)

	// Información financiera (copiada de la orden en el momento de facturación)
	Subtotal     float64 `gorm:"type:decimal(12,2);not null"` // Subtotal
	Tax          float64 `gorm:"type:decimal(12,2);not null"` // Impuestos
	Discount     float64 `gorm:"type:decimal(12,2);not null"` // Descuentos
	ShippingCost float64 `gorm:"type:decimal(12,2);not null"` // Costo de envío
	TotalAmount  float64 `gorm:"type:decimal(12,2);not null"` // Total
	Currency     string  `gorm:"size:10;default:'COP'"`       // Moneda

	// Información del cliente (desnormalizada de la orden)
	CustomerName  string `gorm:"size:255;not null"` // Nombre del cliente
	CustomerEmail string `gorm:"size:255"`          // Email
	CustomerPhone string `gorm:"size:32"`           // Teléfono
	CustomerDNI   string `gorm:"size:64"`           // DNI/NIT/Identificación

	// Estado de la factura
	// "draft" = borrador (no enviado al proveedor)
	// "pending" = pendiente (enviado al proveedor, esperando respuesta)
	// "issued" = emitida (factura generada exitosamente)
	// "cancelled" = cancelada
	// "failed" = falló la generación
	Status string `gorm:"size:64;not null;index;default:'pending'"`

	// Timestamps
	IssuedAt   *time.Time `gorm:"index"` // Cuándo se emitió la factura
	CancelledAt *time.Time               // Cuándo se canceló
	ExpiresAt  *time.Time               // Cuándo expira (si aplica)

	// URLs y archivos
	InvoiceURL *string `gorm:"size:512"` // URL del PDF/XML de la factura
	PDFURL     *string `gorm:"size:512"` // URL del PDF
	XMLURL     *string `gorm:"size:512"` // URL del XML
	CUFE       *string `gorm:"size:255"` // CUFE (Código Único de Factura Electrónica - Colombia)

	// Información adicional
	Notes    *string        `gorm:"type:text"`  // Notas de la factura
	Metadata datatypes.JSON `gorm:"type:jsonb"` // Metadata adicional del proveedor

	// Respuesta del proveedor (JSON crudo)
	ProviderResponse datatypes.JSON `gorm:"type:jsonb"` // Respuesta completa del proveedor

	// Relaciones
	InvoiceItems []InvoiceItem    `gorm:"foreignKey:InvoiceID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	SyncLogs     []InvoiceSyncLog `gorm:"foreignKey:InvoiceID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CreditNotes  []CreditNote     `gorm:"foreignKey:InvoiceID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

// TableName especifica el nombre de la tabla para Invoice
func (Invoice) TableName() string {
	return "invoices"
}

// BeforeCreate genera el número interno antes de crear
func (i *Invoice) BeforeCreate(tx *gorm.DB) error {
	// Generar número interno si no existe
	if i.InternalNumber == "" {
		i.InternalNumber = generateInvoiceNumber()
	}
	return nil
}

// generateInvoiceNumber genera un número interno único para la factura
func generateInvoiceNumber() string {
	// Formato: INV-YYYYMMDD-HHMMSS-RANDOM
	now := time.Now()
	return "INV-" + now.Format("20060102-150405") + "-" + generateRandomString(6)
}
