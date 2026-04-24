package models

import "gorm.io/gorm"

type InventoryState struct {
	gorm.Model
	Code        string `gorm:"size:30;not null;uniqueIndex"`
	Name        string `gorm:"size:100;not null"`
	Description string `gorm:"size:255"`
	IsTerminal  bool   `gorm:"default:false"`
	IsActive    bool   `gorm:"default:true;index"`
}

func (InventoryState) TableName() string {
	return "inventory_states"
}
