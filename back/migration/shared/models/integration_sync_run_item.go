package models

import "time"

type IntegrationSyncRunItem struct {
	ID uint `gorm:"primarykey"`

	RunID uint `gorm:"not null;index:idx_sync_run_item_lookup,priority:1"`

	GroupCode string `gorm:"column:group_code;size:32;not null;index:idx_sync_run_item_lookup,priority:2"`

	SKU   string `gorm:"column:sku;size:120;index"`
	Label string `gorm:"size:300"`
	Tone  string `gorm:"size:16"`

	ParentRef    string `gorm:"column:parent_ref;size:64;index"`
	ParentLabel  string `gorm:"column:parent_label;size:300"`
	VariantLabel string `gorm:"column:variant_label;size:200"`

	CounterpartSKU  string `gorm:"column:counterpart_sku;size:120"`
	CounterpartName string `gorm:"column:counterpart_name;size:300"`
	ChannelQty      *int   `gorm:"column:channel_qty"`
	OwnQty          *int   `gorm:"column:own_qty"`
	FixSide         string `gorm:"column:fix_side;size:16"`
	Pattern         string `gorm:"column:pattern;size:24"`

	CreatedAt time.Time
}

func (IntegrationSyncRunItem) TableName() string {
	return "integration_sync_run_items"
}
