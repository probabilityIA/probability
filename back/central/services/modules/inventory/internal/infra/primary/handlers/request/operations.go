package request

type CreatePutawayRuleBody struct {
	ProductID    *string `json:"product_id"`
	CategoryID   *uint   `json:"category_id"`
	TargetZoneID uint    `json:"target_zone_id" binding:"required,min=1"`
	Priority     int     `json:"priority"`
	Strategy     string  `json:"strategy"`
	IsActive     *bool   `json:"is_active"`
}

type UpdatePutawayRuleBody struct {
	ProductID    *string `json:"product_id"`
	CategoryID   *uint   `json:"category_id"`
	TargetZoneID *uint   `json:"target_zone_id"`
	Priority     *int    `json:"priority"`
	Strategy     string  `json:"strategy"`
	IsActive     *bool   `json:"is_active"`
}

type SuggestPutawayItemBody struct {
	ProductID string `json:"product_id" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required,min=1"`
}

type SuggestPutawayBody struct {
	Items []SuggestPutawayItemBody `json:"items" binding:"required,min=1,dive"`
}

type ConfirmPutawayBody struct {
	ActualLocationID uint `json:"actual_location_id" binding:"required,min=1"`
}

type CreateReplenishmentTaskBody struct {
	ProductID      string `json:"product_id" binding:"required"`
	WarehouseID    uint   `json:"warehouse_id" binding:"required,min=1"`
	FromLocationID *uint  `json:"from_location_id"`
	ToLocationID   *uint  `json:"to_location_id"`
	Quantity       int    `json:"quantity" binding:"required,min=1"`
	TriggeredBy    string `json:"triggered_by"`
	Notes          string `json:"notes"`
}

type AssignReplenishmentBody struct {
	UserID uint `json:"user_id" binding:"required,min=1"`
}

type CompleteReplenishmentBody struct {
	Notes string `json:"notes"`
}

type CancelReplenishmentBody struct {
	Reason string `json:"reason"`
}

type CreateCrossDockLinkBody struct {
	InboundShipmentID *uint  `json:"inbound_shipment_id"`
	OutboundOrderID   string `json:"outbound_order_id" binding:"required"`
	ProductID         string `json:"product_id" binding:"required"`
	Quantity          int    `json:"quantity" binding:"required,min=1"`
}

type RunSlottingBody struct {
	WarehouseID uint   `json:"warehouse_id" binding:"required,min=1"`
	Period      string `json:"period"`
}
