package request

type ManualReplyRequest struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
	BusinessID  uint   `json:"business_id"`
	Text        string `json:"text"         binding:"required,min=1,max=4096"`
}
