package models

import "gorm.io/gorm"

type WarehouseZone struct {
	gorm.Model
	WarehouseID uint   `gorm:"not null;index;uniqueIndex:idx_zone_warehouse_code,priority:1"`
	BusinessID  uint   `gorm:"not null;index"`
	Code        string `gorm:"size:50;not null;uniqueIndex:idx_zone_warehouse_code,priority:2"`
	Name        string `gorm:"size:255;not null"`
	Purpose     string `gorm:"size:30;default:'storage';index"`
	IsActive    bool   `gorm:"default:true;index"`
	ColorHex    string `gorm:"size:7"`

	Warehouse Warehouse        `gorm:"foreignKey:WarehouseID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Aisles    []WarehouseAisle `gorm:"foreignKey:ZoneID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (WarehouseZone) TableName() string {
	return "warehouse_zones"
}

type WarehouseAisle struct {
	gorm.Model
	ZoneID     uint   `gorm:"not null;index;uniqueIndex:idx_aisle_zone_code,priority:1"`
	BusinessID uint   `gorm:"not null;index"`
	Code       string `gorm:"size:50;not null;uniqueIndex:idx_aisle_zone_code,priority:2"`
	Name       string `gorm:"size:255;not null"`
	IsActive   bool   `gorm:"default:true;index"`

	Zone  WarehouseZone   `gorm:"foreignKey:ZoneID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Racks []WarehouseRack `gorm:"foreignKey:AisleID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (WarehouseAisle) TableName() string {
	return "warehouse_aisles"
}

type WarehouseRack struct {
	gorm.Model
	AisleID     uint   `gorm:"not null;index;uniqueIndex:idx_rack_aisle_code,priority:1"`
	BusinessID  uint   `gorm:"not null;index"`
	Code        string `gorm:"size:50;not null;uniqueIndex:idx_rack_aisle_code,priority:2"`
	Name        string `gorm:"size:255;not null"`
	LevelsCount int    `gorm:"default:0"`
	IsActive    bool   `gorm:"default:true;index"`

	Aisle  WarehouseAisle         `gorm:"foreignKey:AisleID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Levels []WarehouseRackLevel `gorm:"foreignKey:RackID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (WarehouseRack) TableName() string {
	return "warehouse_racks"
}

type WarehouseRackLevel struct {
	gorm.Model
	RackID     uint   `gorm:"not null;index;uniqueIndex:idx_level_rack_code,priority:1"`
	BusinessID uint   `gorm:"not null;index"`
	Code       string `gorm:"size:50;not null;uniqueIndex:idx_level_rack_code,priority:2"`
	Ordinal    int    `gorm:"default:0"`
	IsActive   bool   `gorm:"default:true;index"`

	Rack      WarehouseRack       `gorm:"foreignKey:RackID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Positions []WarehouseLocation `gorm:"foreignKey:LevelID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

func (WarehouseRackLevel) TableName() string {
	return "warehouse_rack_levels"
}
