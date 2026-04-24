package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// ProductFamily representa el producto padre o familia a la que pertenecen varias variantes.
type ProductFamily struct {
	gorm.Model

	BusinessID   uint           `gorm:"not null;index"`
	Name         string         `gorm:"size:255;not null"`
	Title        string         `gorm:"size:500"`
	Description  string         `gorm:"type:text"`
	Slug         string         `gorm:"size:255;index"`
	Category     string         `gorm:"size:255;index"`
	Brand        string         `gorm:"size:255;index"`
	ImageURL     string         `gorm:"size:500"`
	Status       string         `gorm:"size:50;default:'active';index"`
	IsActive     bool           `gorm:"default:true;index"`
	VariantAxes  datatypes.JSON `gorm:"type:jsonb"`
	Metadata     datatypes.JSON `gorm:"type:jsonb"`
	VariantCount int64          `gorm:"column:variant_count;->;-:migration" json:"variant_count,omitempty"`

	Business Business  `gorm:"foreignKey:BusinessID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Products []Product `gorm:"foreignKey:FamilyID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

// TableName especifica el nombre de la tabla.
func (ProductFamily) TableName() string {
	return "product_families"
}
