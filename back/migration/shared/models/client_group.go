package models

import "gorm.io/gorm"

type ClientGroup struct {
	gorm.Model
	BusinessID  uint   `gorm:"not null;index;uniqueIndex:idx_client_group_biz_name,priority:1"`
	Name        string `gorm:"size:120;not null;uniqueIndex:idx_client_group_biz_name,priority:2"`
	Description string `gorm:"size:500"`
	Color       string `gorm:"size:20"`
	IsActive    bool   `gorm:"default:true;index"`

	Business Business            `gorm:"foreignKey:BusinessID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Members  []ClientGroupMember `gorm:"foreignKey:ClientGroupID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
