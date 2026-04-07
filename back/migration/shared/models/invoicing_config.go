package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

//
//	INVOICING CONFIGS - Configuraciones de facturación por integración
//

// InvoicingConfig define qué integraciones deben facturar automáticamente
// y con qué proveedor de facturación.
// Ejemplo: "Órdenes de Shopify y Plataforma se facturan con Softpymes"
type InvoicingConfig struct {
	gorm.Model

	// Relación con Business
	BusinessID uint     `gorm:"not null;index"`
	Business   Business `gorm:"foreignKey:BusinessID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	// NOTA: La relación con las integraciones de e-commerce (fuentes de órdenes)
	// ahora se gestiona via tabla pivote invoicing_config_integrations.
	ConfigIntegrations []InvoicingConfigIntegration `gorm:"foreignKey:ConfigID"`

	// Relación con InvoicingProvider (DEPRECATED - mantener temporalmente para dual-read)
	InvoicingProviderID *uint             `gorm:"index"`
	InvoicingProvider   InvoicingProvider `gorm:"foreignKey:InvoicingProviderID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	// Relación con Integration (proveedor de facturación desde integrations/)
	// Unique: solo una config por negocio + proveedor de facturación (aplicado via índice parcial en migración)
	InvoicingIntegrationID *uint       `gorm:"index"`
	InvoicingIntegration   Integration `gorm:"foreignKey:InvoicingIntegrationID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	// Estado
	Enabled      bool `gorm:"default:true;index"`  // Si la configuración está habilitada
	AutoInvoice  bool `gorm:"default:false;index"` // Si factura automáticamente al crear orden

	// Filtros (JSON - define qué órdenes deben facturarse)
	Filters datatypes.JSON `gorm:"type:jsonb"`

	// Configuración adicional de facturación
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
