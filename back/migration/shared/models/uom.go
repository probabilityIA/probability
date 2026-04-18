package models

import "gorm.io/gorm"

type UnitOfMeasure struct {
	gorm.Model
	Code     string `gorm:"size:20;not null;uniqueIndex"`
	Name     string `gorm:"size:100;not null"`
	Type     string `gorm:"size:20;not null;index"`
	IsActive bool   `gorm:"default:true;index"`
}

func (UnitOfMeasure) TableName() string {
	return "units_of_measure"
}

type ProductUoM struct {
	gorm.Model
	ProductID        string  `gorm:"type:varchar(64);not null;uniqueIndex:idx_product_uom,priority:1"`
	UomID            uint    `gorm:"not null;uniqueIndex:idx_product_uom,priority:2"`
	BusinessID       uint    `gorm:"not null;index"`
	ConversionFactor float64 `gorm:"not null;default:1"`
	IsBase           bool    `gorm:"default:false;index"`
	Barcode          string  `gorm:"size:100;index"`
	IsActive         bool    `gorm:"default:true;index"`

	Uom      UnitOfMeasure `gorm:"foreignKey:UomID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Product  Product       `gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Business Business      `gorm:"foreignKey:BusinessID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (ProductUoM) TableName() string {
	return "product_uoms"
}
