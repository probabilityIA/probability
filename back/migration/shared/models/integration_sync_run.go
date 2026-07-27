package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type IntegrationSyncRun struct {
	gorm.Model

	BusinessID    uint `gorm:"not null;index:idx_sync_run_lookup,priority:1"`
	IntegrationID uint `gorm:"not null;index:idx_sync_run_lookup,priority:2"`

	Kind string `gorm:"size:32;not null;index:idx_sync_run_lookup,priority:3"`

	CorrelationID string     `gorm:"size:64;index"`
	StartedAt     time.Time  `gorm:"index"`
	FinishedAt    *time.Time `gorm:"index"`

	Total     int `gorm:"default:0"`
	Updated   int `gorm:"default:0"`
	Unchanged int `gorm:"default:0"`
	Skipped   int `gorm:"default:0"`
	Failed    int `gorm:"default:0"`

	Matched           int `gorm:"default:0"`
	NotAssociated     int `gorm:"default:0"`
	OnlyInProbability int `gorm:"default:0"`
	OnlyInChannel     int `gorm:"default:0"`

	Status  string `gorm:"size:32;not null;default:'completed';index"`
	Message string `gorm:"size:500"`

	Detail datatypes.JSON `gorm:"type:jsonb"`
}

func (IntegrationSyncRun) TableName() string {
	return "integration_sync_runs"
}
