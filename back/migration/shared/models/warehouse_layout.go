package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type WarehouseLayout struct {
	gorm.Model
	WarehouseID  uint           `gorm:"not null;uniqueIndex:idx_layout_warehouse"`
	BusinessID   uint           `gorm:"not null;index"`
	CanvasWidth  float64        `gorm:"default:1200"`
	CanvasHeight float64        `gorm:"default:800"`
	GridSize     float64        `gorm:"default:20"`
	Nodes        datatypes.JSON `gorm:"type:jsonb"`

	Warehouse Warehouse `gorm:"foreignKey:WarehouseID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (WarehouseLayout) TableName() string {
	return "warehouse_layouts"
}
