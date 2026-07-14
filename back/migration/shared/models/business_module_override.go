package models

import "gorm.io/gorm"

type BusinessModuleOverride struct {
	gorm.Model
	BusinessID      uint    `gorm:"not null;index;uniqueIndex:idx_business_module"`
	ModuleCode      string  `gorm:"size:60;not null;uniqueIndex:idx_business_module"`
	GrantedByUserID uint    `gorm:"not null"`
	Notes           *string `gorm:"type:text"`

	Business Business `gorm:"foreignKey:BusinessID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (BusinessModuleOverride) TableName() string {
	return "business_module_overrides"
}
