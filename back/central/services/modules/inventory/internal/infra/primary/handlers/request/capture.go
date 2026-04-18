package request

type CreateLPNBody struct {
	Code       string `json:"code" binding:"required"`
	LpnType    string `json:"lpn_type"`
	LocationID *uint  `json:"location_id"`
}

type UpdateLPNBody struct {
	Code       string `json:"code"`
	LpnType    string `json:"lpn_type"`
	LocationID *uint  `json:"location_id"`
	Status     string `json:"status"`
}

type AddToLPNBody struct {
	ProductID string `json:"product_id" binding:"required"`
	LotID     *uint  `json:"lot_id"`
	SerialID  *uint  `json:"serial_id"`
	Qty       int    `json:"qty" binding:"required,min=1"`
}

type MoveLPNBody struct {
	NewLocationID uint `json:"new_location_id" binding:"required,min=1"`
}

type MergeLPNBody struct {
	TargetLpnID uint `json:"target_lpn_id" binding:"required,min=1"`
}

type ScanBody struct {
	Code     string `json:"code" binding:"required"`
	DeviceID string `json:"device_id"`
	Action   string `json:"action"`
}

type InboundSyncBody struct {
	Payload map[string]any `json:"payload" binding:"required"`
}
