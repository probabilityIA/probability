package models

import "gorm.io/gorm"

type CustomProductPrice struct {
	gorm.Model
	BusinessID    uint    `gorm:"not null;index"`
	ProductID     string  `gorm:"type:varchar(64);not null;index;uniqueIndex:idx_custom_price_group,priority:1;uniqueIndex:idx_custom_price_client,priority:1"`
	ClientGroupID *uint   `gorm:"index;uniqueIndex:idx_custom_price_group,priority:2"`
	ClientID      *uint   `gorm:"index;uniqueIndex:idx_custom_price_client,priority:2"`
	Price         float64 `gorm:"type:decimal(15,2);not null"`
	IsActive      bool    `gorm:"default:true;index"`

	Business    Business     `gorm:"foreignKey:BusinessID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Product     Product      `gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	ClientGroup *ClientGroup `gorm:"foreignKey:ClientGroupID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Client      *Client      `gorm:"foreignKey:ClientID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
