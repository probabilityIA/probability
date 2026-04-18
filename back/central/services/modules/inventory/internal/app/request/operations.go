package request

type CreatePutawayRuleDTO struct {
	BusinessID   uint
	ProductID    *string
	CategoryID   *uint
	TargetZoneID uint
	Priority     int
	Strategy     string
	IsActive     bool
}

type UpdatePutawayRuleDTO struct {
	ID           uint
	BusinessID   uint
	ProductID    *string
	CategoryID   *uint
	TargetZoneID *uint
	Priority     *int
	Strategy     string
	IsActive     *bool
}

type PutawaySuggestItem struct {
	ProductID string
	Quantity  int
}

type PutawaySuggestDTO struct {
	BusinessID uint
	Items      []PutawaySuggestItem
}

type ConfirmPutawayDTO struct {
	BusinessID       uint
	SuggestionID     uint
	ActualLocationID uint
	UserID           *uint
}

type CreateReplenishmentTaskDTO struct {
	BusinessID     uint
	ProductID      string
	WarehouseID    uint
	FromLocationID *uint
	ToLocationID   *uint
	Quantity       int
	TriggeredBy    string
	Notes          string
}

type AssignReplenishmentDTO struct {
	BusinessID uint
	TaskID     uint
	UserID     uint
}

type CompleteReplenishmentDTO struct {
	BusinessID uint
	TaskID     uint
	Notes      string
}

type CreateCrossDockLinkDTO struct {
	BusinessID        uint
	InboundShipmentID *uint
	OutboundOrderID   string
	ProductID         string
	Quantity          int
}

type ExecuteCrossDockDTO struct {
	BusinessID uint
	LinkID     uint
}

type RunSlottingDTO struct {
	BusinessID  uint
	WarehouseID uint
	Period      string
}
