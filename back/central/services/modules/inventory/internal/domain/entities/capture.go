package entities

import "time"

type LicensePlate struct {
	ID                uint
	BusinessID        uint
	Code              string
	LpnType           string
	CurrentLocationID *uint
	Status            string
	CreatedAt         time.Time
	UpdatedAt         time.Time
	Lines             []LicensePlateLine
}

type LicensePlateLine struct {
	ID         uint
	LpnID      uint
	BusinessID uint
	ProductID  string
	LotID      *uint
	SerialID   *uint
	Qty        int
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type ScanEvent struct {
	ID          uint
	BusinessID  uint
	UserID      *uint
	DeviceID    string
	ScannedCode string
	CodeType    string
	Action      string
	ScannedAt   time.Time
	CreatedAt   time.Time
}

type ScanResolution struct {
	Code       string
	CodeType   string
	MatchedID  *uint
	ProductID  string
	LocationID *uint
	LotID      *uint
	SerialID   *uint
	LpnID      *uint
	Suggested  string
	Data       map[string]any
}

type InventorySyncLog struct {
	ID            uint
	BusinessID    uint
	IntegrationID *uint
	Direction     string
	PayloadHash   string
	Status        string
	Error         string
	SyncedAt      *time.Time
	CreatedAt     time.Time
}
