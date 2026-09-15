package request

type AIControlRequest struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
	BusinessID  uint   `json:"business_id"`
}
