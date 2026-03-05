package request

type AddStopRequest struct {
	OrderID       *string  `json:"order_id"`
	Address       string   `json:"address" binding:"required,max=500"`
	City          string   `json:"city" binding:"omitempty,max=128"`
	Lat           *float64 `json:"lat"`
	Lng           *float64 `json:"lng"`
	CustomerName  string   `json:"customer_name" binding:"omitempty,max=255"`
	CustomerPhone string   `json:"customer_phone" binding:"omitempty,max=50"`
	DeliveryNotes *string  `json:"delivery_notes"`
}

type UpdateStopRequest struct {
	Address       string   `json:"address" binding:"required,max=500"`
	City          string   `json:"city" binding:"omitempty,max=128"`
	Lat           *float64 `json:"lat"`
	Lng           *float64 `json:"lng"`
	CustomerName  string   `json:"customer_name" binding:"omitempty,max=255"`
	CustomerPhone string   `json:"customer_phone" binding:"omitempty,max=50"`
	DeliveryNotes *string  `json:"delivery_notes"`
}

type UpdateStopStatusRequest struct {
	Status        string  `json:"status" binding:"required,oneof=arrived delivered failed skipped"`
	FailureReason *string `json:"failure_reason"`
	SignatureURL  string  `json:"signature_url" binding:"omitempty,max=512"`
	PhotoURL      string  `json:"photo_url" binding:"omitempty,max=512"`
}

type ReorderStopsRequest struct {
	StopIDs []uint `json:"stop_ids" binding:"required"`
}
