package models

import "gorm.io/gorm"

type TikTokShopSnapshot struct {
	gorm.Model
	BusinessID           uint   `gorm:"not null;index:idx_tiktok_snapshots_business_integration,priority:1"`
	IntegrationID        uint   `gorm:"not null;index:idx_tiktok_snapshots_business_integration,priority:2"`
	ShopID               int64  `gorm:"index"`
	ShopName             string `gorm:"size:255;index"`
	ShopRating           string `gorm:"size:16"`
	ShopReviewCount      int64  `gorm:"not null;default:0"`
	ItemSoldCount        int64  `gorm:"not null;default:0"`
	ShopPerformanceValue int64  `gorm:"not null;default:0"`
}
