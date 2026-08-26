package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type ShippingConfig struct {
	gorm.Model
	BusinessID      uint           `gorm:"not null;index;uniqueIndex:idx_shipping_config_scope,priority:1"`
	WarehouseID     *uint          `gorm:"index;uniqueIndex:idx_shipping_config_scope,priority:2"`
	PackageStrategy string         `gorm:"size:32;not null;default:'product_dimensions'"`
	Boxes           datatypes.JSON `gorm:"type:jsonb"`
	Carriers        datatypes.JSON `gorm:"type:jsonb"`
	AlwaysInsure    bool           `gorm:"not null;default:false"`
	CreatedBy       uint
	CreatedByName   string `gorm:"size:255"`
	UpdatedBy       uint
	UpdatedByName   string `gorm:"size:255"`
}

func (ShippingConfig) TableName() string {
	return "shipping_configs"
}
