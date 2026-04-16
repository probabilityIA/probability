package request

// CreateJournal es el request para crear un comprobante contable manualmente
type CreateJournal struct {
	OrderID string `json:"order_id" binding:"required"`
}
