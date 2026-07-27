package models

import "time"

type IntegrationSyncRunItem struct {
	ID uint `gorm:"primarykey"`

	RunID uint `gorm:"not null;index:idx_sync_run_item_lookup,priority:1"`

	GroupCode string `gorm:"column:group_code;size:32;not null;index:idx_sync_run_item_lookup,priority:2"`

	SKU   string `gorm:"column:sku;size:120;index"`
	Label string `gorm:"size:300"`
	Tone  string `gorm:"size:16"`

	CreatedAt time.Time
}

func (IntegrationSyncRunItem) TableName() string {
	return "integration_sync_run_items"
}
