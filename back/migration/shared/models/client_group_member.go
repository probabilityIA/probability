package models

import "gorm.io/gorm"

type ClientGroupMember struct {
	gorm.Model
	BusinessID    uint `gorm:"not null;index"`
	ClientGroupID uint `gorm:"not null;index"`
	ClientID      uint `gorm:"not null;index;uniqueIndex:idx_client_group_member_client"`

	Business    Business    `gorm:"foreignKey:BusinessID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	ClientGroup ClientGroup `gorm:"foreignKey:ClientGroupID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Client      Client      `gorm:"foreignKey:ClientID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
