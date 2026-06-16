package models

import (
	"time"

	"gorm.io/datatypes"
)

type WalletKPISelection struct {
	ID                  uint           `gorm:"primaryKey"`
	SelectedBusinessIDs datatypes.JSON `gorm:"type:jsonb;not null"`
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

func (WalletKPISelection) TableName() string { return "wallet_kpi_selection" }
