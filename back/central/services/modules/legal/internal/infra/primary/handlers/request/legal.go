package request

type AcceptDocuments struct {
	DocumentIDs []uint `json:"document_ids" binding:"required,min=1"`
}
