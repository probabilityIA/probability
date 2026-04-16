package request

// CreateQuantityDiscountRequest payload de creación de descuento por cantidad
type CreateQuantityDiscountRequest struct {
	ProductID       *string `json:"product_id" binding:"omitempty"`
	MinQuantity     int     `json:"min_quantity" binding:"required,min=1"`
	DiscountPercent float64 `json:"discount_percent" binding:"required,gt=0,lte=100"`
	IsActive        *bool   `json:"is_active"`
	Description     string  `json:"description" binding:"omitempty,max=255"`
}

// UpdateQuantityDiscountRequest payload de actualización de descuento por cantidad
type UpdateQuantityDiscountRequest struct {
	ProductID       *string `json:"product_id" binding:"omitempty"`
	MinQuantity     int     `json:"min_quantity" binding:"required,min=1"`
	DiscountPercent float64 `json:"discount_percent" binding:"required,gt=0,lte=100"`
	IsActive        *bool   `json:"is_active"`
	Description     string  `json:"description" binding:"omitempty,max=255"`
}
