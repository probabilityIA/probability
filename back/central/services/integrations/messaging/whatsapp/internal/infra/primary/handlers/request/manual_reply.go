package request

// ManualReplyRequest cuerpo del endpoint POST /whatsapp/conversations/:id/reply
type ManualReplyRequest struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
	BusinessID  uint   `json:"business_id"  binding:"required"`
	Text        string `json:"text"         binding:"required,min=1,max=4096"`
}
