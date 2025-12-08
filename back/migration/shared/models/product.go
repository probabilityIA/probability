package models

import "gorm.io/gorm"

// Product representa un producto único en el catálogo del negocio
type Product struct {
	gorm.Model
	BusinessID uint   `gorm:"not null;index;uniqueIndex:idx_business_product_sku,priority:1"`
	SKU        string `gorm:"size:128;not null;uniqueIndex:idx_business_product_sku,priority:2"`
	Name       string `gorm:"size:255;not null"`
	ExternalID string `gorm:"size:255;index"` // ID en la plataforma externa

	// Relación
	Business Business `gorm:"foreignKey:BusinessID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

// TableName especifica el nombre de la tabla
func (Product) TableName() string {
	return "products"
}
