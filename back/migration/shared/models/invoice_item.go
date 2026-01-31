package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// ───────────────────────────────────────────
//
//	INVOICE ITEMS - Líneas/Items de factura
//
// ───────────────────────────────────────────

// InvoiceItem representa un item/línea de una factura
type InvoiceItem struct {
	gorm.Model

	// Relación con Invoice
	InvoiceID uint    `gorm:"not null;index"`
	Invoice   Invoice `gorm:"foreignKey:InvoiceID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	// Relación con Product (opcional - puede ser null si el producto se eliminó)
	ProductID *string `gorm:"type:varchar(64);index"`
	Product   Product `gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`

	// Información del item (copiada de OrderItem en el momento de facturación)
	SKU         string  `gorm:"size:255"`                    // SKU del producto
	Name        string  `gorm:"size:255;not null"`           // Nombre del producto
	Description *string `gorm:"type:text"`                   // Descripción
	Quantity    int     `gorm:"not null;default:1"`          // Cantidad
	UnitPrice   float64 `gorm:"type:decimal(12,2);not null"` // Precio unitario
	TotalPrice  float64 `gorm:"type:decimal(12,2);not null"` // Precio total
	Currency    string  `gorm:"size:10;default:'COP'"`       // Moneda

	// Impuestos y descuentos
	Tax      float64  `gorm:"type:decimal(12,2);default:0"` // Impuesto
	TaxRate  *float64 `gorm:"type:decimal(5,4)"`            // Tasa de impuesto (ej: 0.19)
	Discount float64  `gorm:"type:decimal(12,2);default:0"` // Descuento aplicado

	// Información adicional del proveedor
	ProviderItemID *string        `gorm:"size:255"` // ID del item en el sistema del proveedor
	Metadata       datatypes.JSON `gorm:"type:jsonb"` // Metadata adicional
}

// TableName especifica el nombre de la tabla para InvoiceItem
func (InvoiceItem) TableName() string {
	return "invoice_items"
}
