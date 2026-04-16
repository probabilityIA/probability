package models

import "gorm.io/gorm"

// InvoicingConfigIntegration relaciona una configuración de facturación
// con una o más integraciones de e-commerce (muchos a muchos).
// Una config puede cubrir N fuentes de órdenes (Shopify, Plataforma, etc.)
type InvoicingConfigIntegration struct {
	gorm.Model

	ConfigID      uint            `gorm:"not null;index;uniqueIndex:idx_config_integration,priority:1"`
	IntegrationID uint            `gorm:"not null;index;uniqueIndex:idx_config_integration,priority:2"`

	// Relación con InvoicingConfig (cascade delete)
	Config        InvoicingConfig `gorm:"foreignKey:ConfigID;constraint:OnDelete:CASCADE"`
	// Relación con Integration
	Integration   Integration     `gorm:"foreignKey:IntegrationID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

// TableName especifica el nombre de la tabla
func (InvoicingConfigIntegration) TableName() string {
	return "invoicing_config_integrations"
}
