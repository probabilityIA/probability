package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// ───────────────────────────────────────────
//
//	INVOICING CONFIGS - Configuraciones de facturación por integración
//
// ───────────────────────────────────────────

// InvoicingConfig define qué integraciones deben facturar automáticamente
// y con qué proveedor de facturación.
// Ejemplo: "Órdenes de Shopify se facturan con Softpymes"
type InvoicingConfig struct {
	gorm.Model

	// Relación con Business
	BusinessID uint     `gorm:"not null;index;uniqueIndex:idx_business_integration_config,priority:1"`
	Business   Business `gorm:"foreignKey:BusinessID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	// Relación con Integration (fuente de órdenes)
	IntegrationID uint        `gorm:"not null;index;uniqueIndex:idx_business_integration_config,priority:2"`
	Integration   Integration `gorm:"foreignKey:IntegrationID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	// Relación con InvoicingProvider (DEPRECATED - mantener temporalmente para dual-read)
	InvoicingProviderID *uint             `gorm:"index"`
	InvoicingProvider   InvoicingProvider `gorm:"foreignKey:InvoicingProviderID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	// Relación con Integration (nuevo - provider de facturación desde integrations/)
	InvoicingIntegrationID *uint       `gorm:"index"`
	InvoicingIntegration   Integration `gorm:"foreignKey:InvoicingIntegrationID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	// Estado
	Enabled      bool `gorm:"default:true;index"`  // Si la configuración está habilitada
	AutoInvoice  bool `gorm:"default:false;index"` // Si factura automáticamente al crear orden

	// Filtros (JSON - define qué órdenes deben facturarse)
	// Estructura:
	//   {
	//     "min_amount": 50000,                    // Monto mínimo para facturar
	//     "payment_status": "paid",               // Solo facturar si está pagada
	//     "payment_methods": [1, 3, 5],           // IDs de métodos de pago permitidos (opcional)
	//     "order_types": ["delivery", "pickup"],  // Tipos de orden permitidos (opcional)
	//     "exclude_statuses": ["cancelled"]       // Estados a excluir
	//   }
	Filters datatypes.JSON `gorm:"type:jsonb"`

	// Configuración adicional de facturación
	// Ejemplo:
	//   {
	//     "include_shipping": true,           // Si incluye costo de envío
	//     "apply_discount": true,             // Si aplica descuentos
	//     "default_tax_rate": 0.19,           // Tasa de impuesto por defecto
	//     "invoice_type": "electronic",       // Tipo de factura
	//     "notes": "Gracias por su compra"    // Notas por defecto
	//   }
	InvoiceConfig datatypes.JSON `gorm:"type:jsonb"`

	// Metadata
	Description string `gorm:"size:500"`
	CreatedByID uint   `gorm:"index"`
	UpdatedByID *uint  `gorm:"index"`

	// Relaciones
	CreatedBy User  `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	UpdatedBy *User `gorm:"foreignKey:UpdatedByID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

// TableName especifica el nombre de la tabla para InvoicingConfig
func (InvoicingConfig) TableName() string {
	return "invoicing_configs"
}
