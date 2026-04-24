package dtos

type ListLPNParams struct {
	BusinessID uint
	LpnType    string
	Status     string
	LocationID *uint
	Page       int
	PageSize   int
}

func (p ListLPNParams) Offset() int {
	if p.Page < 1 {
		p.Page = 1
	}
	return (p.Page - 1) * p.PageSize
}

type ListSyncLogsParams struct {
	BusinessID    uint
	IntegrationID *uint
	Direction     string
	Status        string
	Page          int
	PageSize      int
}

func (p ListSyncLogsParams) Offset() int {
	if p.Page < 1 {
		p.Page = 1
	}
	return (p.Page - 1) * p.PageSize
}

type LPNAddLine struct {
	ProductID string
	LotID     *uint
	SerialID  *uint
	Qty       int
}

type ScanResolveParams struct {
	BusinessID  uint
	Code        string
	DeviceID    string
	UserID      *uint
	Action      string
	ContextJSON map[string]any
}
