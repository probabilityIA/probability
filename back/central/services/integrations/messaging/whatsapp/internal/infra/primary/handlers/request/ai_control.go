package request

// AIControlRequest body para pause-ai y resume-ai
type AIControlRequest struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
	BusinessID  uint   `json:"business_id"  binding:"required"`
}
