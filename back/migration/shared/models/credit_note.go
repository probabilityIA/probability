package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// ───────────────────────────────────────────
//
//	CREDIT NOTES - Notas de crédito
//
// ───────────────────────────────────────────

// CreditNote representa una nota de crédito asociada a una factura
// Se usa para devoluciones, anulaciones parciales o correcciones
type CreditNote struct {
	gorm.Model

	// Relación con Invoice
	InvoiceID uint    `gorm:"not null;index"`
	Invoice   Invoice `gorm:"foreignKey:InvoiceID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	// Relación con Business
	BusinessID uint     `gorm:"not null;index"`
	Business   Business `gorm:"foreignKey:BusinessID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	// Identificadores
	CreditNoteNumber string  `gorm:"size:128;not null;index"` // Número de nota de crédito del proveedor
	ExternalID       *string `gorm:"size:255;index"`          // ID en el sistema del proveedor
	InternalNumber   string  `gorm:"size:128;unique;index"`   // Número interno Probability (auto-generado)

	// Tipo de nota de crédito
	// "full_refund" = devolución total
	// "partial_refund" = devolución parcial
	// "cancellation" = anulación
	// "correction" = corrección
	NoteType string `gorm:"size:64;not null;index"`

	// Información financiera
	Amount   float64 `gorm:"type:decimal(12,2);not null"` // Monto de la nota de crédito
	Currency string  `gorm:"size:10;default:'COP'"`       // Moneda

	// Razón y descripción
	Reason      string  `gorm:"size:255;not null"` // Razón de la nota de crédito
	Description *string `gorm:"type:text"`         // Descripción detallada

	// Estado
	// "draft" = borrador
	// "pending" = pendiente
	// "issued" = emitida
	// "cancelled" = cancelada
	// "failed" = falló la generación
	Status string `gorm:"size:64;not null;index;default:'pending'"`

	// Timestamps
	IssuedAt    *time.Time `gorm:"index"` // Cuándo se emitió
	CancelledAt *time.Time               // Cuándo se canceló

	// URLs y archivos
	NoteURL *string `gorm:"size:512"` // URL del PDF/XML de la nota
	PDFURL  *string `gorm:"size:512"` // URL del PDF
	XMLURL  *string `gorm:"size:512"` // URL del XML
	CUFE    *string `gorm:"size:255"` // CUFE (si aplica)

	// Metadata
	Metadata datatypes.JSON `gorm:"type:jsonb"` // Metadata adicional del proveedor

	// Respuesta del proveedor (JSON crudo)
	ProviderResponse datatypes.JSON `gorm:"type:jsonb"` // Respuesta completa del proveedor

	// Usuario que creó la nota
	CreatedByID uint `gorm:"not null;index"`
	CreatedBy   User `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

// TableName especifica el nombre de la tabla para CreditNote
func (CreditNote) TableName() string {
	return "credit_notes"
}

// BeforeCreate genera el número interno antes de crear
func (c *CreditNote) BeforeCreate(tx *gorm.DB) error {
	// Generar número interno si no existe
	if c.InternalNumber == "" {
		c.InternalNumber = generateCreditNoteNumber()
	}
	return nil
}

// generateCreditNoteNumber genera un número interno único para la nota de crédito
func generateCreditNoteNumber() string {
	// Formato: CN-YYYYMMDD-HHMMSS-RANDOM
	now := time.Now()
	return "CN-" + now.Format("20060102-150405") + "-" + generateRandomString(6)
}
