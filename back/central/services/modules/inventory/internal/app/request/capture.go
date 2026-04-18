package request

type CreateLPNDTO struct {
	BusinessID uint
	Code       string
	LpnType    string
	LocationID *uint
}

type UpdateLPNDTO struct {
	ID         uint
	BusinessID uint
	Code       string
	LpnType    string
	LocationID *uint
	Status     string
}

type AddToLPNDTO struct {
	BusinessID uint
	LpnID      uint
	ProductID  string
	LotID      *uint
	SerialID   *uint
	Qty        int
}

type MoveLPNDTO struct {
	BusinessID    uint
	LpnID         uint
	NewLocationID uint
}

type DissolveLPNDTO struct {
	BusinessID uint
	LpnID      uint
}

type MergeLPNDTO struct {
	BusinessID  uint
	SourceLpnID uint
	TargetLpnID uint
}

type ScanDTO struct {
	BusinessID uint
	Code       string
	DeviceID   string
	UserID     *uint
	Action     string
}

type InboundSyncDTO struct {
	BusinessID    uint
	IntegrationID uint
	Payload       map[string]any
}
