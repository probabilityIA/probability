package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type ShipmentSyncLog struct {
	gorm.Model

	ShipmentID *uint     `gorm:"index"`
	Shipment   *Shipment `gorm:"foreignKey:ShipmentID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`

	OperationType string `gorm:"size:64;not null;index"`
	Provider      string `gorm:"size:64;not null;index"`
	Status        string `gorm:"size:32;not null;index"`

	RequestURL     string         `gorm:"size:1024"`
	RequestMethod  string         `gorm:"size:16"`
	RequestPayload datatypes.JSON `gorm:"type:jsonb"`

	ResponseStatus int            `gorm:"index"`
	ResponseBody   datatypes.JSON `gorm:"type:jsonb"`

	ErrorMessage *string `gorm:"type:text"`
	ErrorCode    *string `gorm:"size:128"`

	CorrelationID string `gorm:"size:128;index"`
	TriggeredBy   string `gorm:"size:32"`
	UserID        *uint  `gorm:"index"`

	StartedAt   time.Time  `gorm:"not null;index"`
	CompletedAt *time.Time `gorm:"index"`
	Duration    *int
}

func (ShipmentSyncLog) TableName() string {
	return "shipment_sync_logs"
}
