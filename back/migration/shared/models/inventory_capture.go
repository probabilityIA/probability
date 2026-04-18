package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type LicensePlate struct {
	gorm.Model
	BusinessID        uint   `gorm:"not null;index;uniqueIndex:idx_lpn_business_code,priority:1"`
	Code              string `gorm:"size:100;not null;uniqueIndex:idx_lpn_business_code,priority:2"`
	LpnType           string `gorm:"size:20;default:'pallet';index"`
	CurrentLocationID *uint  `gorm:"index"`
	Status            string `gorm:"size:20;default:'active';index"`
	Metadata          datatypes.JSON `gorm:"type:jsonb"`

	Business Business           `gorm:"foreignKey:BusinessID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Location *WarehouseLocation `gorm:"foreignKey:CurrentLocationID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Lines    []LicensePlateLine `gorm:"foreignKey:LpnID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (LicensePlate) TableName() string {
	return "license_plates"
}

type LicensePlateLine struct {
	gorm.Model
	LpnID      uint   `gorm:"not null;index"`
	BusinessID uint   `gorm:"not null;index"`
	ProductID  string `gorm:"type:varchar(64);not null;index"`
	LotID      *uint  `gorm:"index"`
	SerialID   *uint  `gorm:"index"`
	Qty        int    `gorm:"not null"`

	Lpn     LicensePlate `gorm:"foreignKey:LpnID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Product Product      `gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (LicensePlateLine) TableName() string {
	return "license_plate_lines"
}

type ScanEvent struct {
	gorm.Model
	BusinessID  uint           `gorm:"not null;index"`
	UserID      *uint          `gorm:"index"`
	DeviceID    string         `gorm:"size:100;index"`
	ScannedCode string         `gorm:"size:255;not null;index"`
	CodeType    string         `gorm:"size:20;index"`
	Action      string         `gorm:"size:30"`
	ContextJSON datatypes.JSON `gorm:"type:jsonb"`
	ScannedAt   time.Time      `gorm:"index"`

	Business Business `gorm:"foreignKey:BusinessID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (ScanEvent) TableName() string {
	return "scan_events"
}

type InventorySyncLog struct {
	gorm.Model
	BusinessID    uint   `gorm:"not null;index"`
	IntegrationID *uint  `gorm:"index"`
	Direction     string `gorm:"size:3;not null;index"`
	PayloadHash   string `gorm:"size:64;index;uniqueIndex:idx_sync_hash_direction,priority:1"`
	DirectionKey  string `gorm:"size:3;index;uniqueIndex:idx_sync_hash_direction,priority:2"`
	Status        string `gorm:"size:20;default:'pending';index"`
	Error         string `gorm:"type:text"`
	Payload       datatypes.JSON `gorm:"type:jsonb"`
	SyncedAt      *time.Time

	Business Business `gorm:"foreignKey:BusinessID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (InventorySyncLog) TableName() string {
	return "inventory_sync_logs"
}
