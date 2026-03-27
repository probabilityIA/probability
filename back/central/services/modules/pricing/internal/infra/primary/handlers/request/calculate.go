package request

// CalculatePriceRequest payload para calcular un precio
type CalculatePriceRequest struct {
	ClientID  uint    `json:"client_id" binding:"required"`
	ProductID string  `json:"product_id" binding:"required"`
	BasePrice float64 `json:"base_price" binding:"required,gt=0"`
	Quantity  int     `json:"quantity" binding:"required,min=1"`
}
