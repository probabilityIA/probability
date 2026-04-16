package models

import "gorm.io/gorm"

// ProductBusinessIntegration representa la relación entre un producto y una integración
// dentro del contexto de un negocio. Un producto puede estar asociado a múltiples
// integraciones, pero todas deben pertenecer al mismo negocio del producto.
type ProductBusinessIntegration struct {
	gorm.Model

	// ID del producto (hash alfanumérico)
	ProductID string `gorm:"type:varchar(64);not null;index;uniqueIndex:idx_product_integration,priority:1"`

	// ID del negocio (desnormalizado para optimización de queries y validación)
	// Este campo DEBE coincidir con el BusinessID del producto y de la integración
	BusinessID uint `gorm:"not null;index"`

	// ID de la integración asociada
	IntegrationID uint `gorm:"not null;index;uniqueIndex:idx_product_integration,priority:2"`

	// ID del producto en el sistema externo de la integración
	// Por ejemplo: el product_id de Shopify, el item_id de Mercado Libre, etc.
	ExternalProductID string `gorm:"size:255;not null;index"`

	// Relaciones
	Product     Product     `gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Business    Business    `gorm:"foreignKey:BusinessID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Integration Integration `gorm:"foreignKey:IntegrationID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

// TableName especifica el nombre de la tabla
func (ProductBusinessIntegration) TableName() string {
	return "product_business_integrations"
}
