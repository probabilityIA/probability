package models

import "time"

type InventoryCompareSnapshot struct {
	ID uint `gorm:"primarykey"`

	BusinessID    uint `gorm:"not null;index:idx_inv_compare_lookup,priority:1"`
	IntegrationID uint `gorm:"not null;index:idx_inv_compare_lookup,priority:2"`

	ProductID         string `gorm:"column:product_id;size:64;not null"`
	SKU               string `gorm:"column:sku;size:120;index"`
	Name              string `gorm:"size:300"`
	ImageURL          string `gorm:"column:image_url;size:600"`
	ExternalItemID    string `gorm:"column:external_item_id;size:120"`
	ExternalVariantID string `gorm:"column:external_variant_id;size:120;not null;default:''"`

	ProbabilityQty *int   `gorm:"column:probability_qty"`
	ChannelQty     *int   `gorm:"column:channel_qty"`
	Delta          *int   `gorm:"column:delta"`
	Action         string `gorm:"size:16;not null;index:idx_inv_compare_lookup,priority:3"`
	Reason         string `gorm:"size:300"`

	CheckedAt time.Time `gorm:"not null;index"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (InventoryCompareSnapshot) TableName() string {
	return "inventory_compare_snapshots"
}
