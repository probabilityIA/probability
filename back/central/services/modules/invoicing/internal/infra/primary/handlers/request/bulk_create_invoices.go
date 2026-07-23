package request

type BulkCreateInvoicesRequest struct {
	OrderIDs   []string `json:"order_ids" binding:"required,min=1,max=1000"`
	BusinessID *uint    `json:"business_id"`
}
