package dtos

// ---- Client Pricing Rules ----

// CreateClientPricingRuleDTO datos para crear una regla de precio
type CreateClientPricingRuleDTO struct {
	BusinessID      uint
	ClientID        uint
	ProductID       *string
	AdjustmentType  string
	AdjustmentValue float64
	IsActive        bool
	Priority        int
	Description     string
}

// UpdateClientPricingRuleDTO datos para actualizar una regla de precio
type UpdateClientPricingRuleDTO struct {
	ID              uint
	BusinessID      uint
	ClientID        uint
	ProductID       *string
	AdjustmentType  string
	AdjustmentValue float64
	IsActive        bool
	Priority        int
	Description     string
}

// ListClientPricingRulesParams parámetros de búsqueda y paginación
type ListClientPricingRulesParams struct {
	BusinessID uint
	ClientID   *uint   // filtro opcional por cliente
	ProductID  *string // filtro opcional por producto
	Page       int
	PageSize   int
}

// Offset calcula el offset para paginación
func (p ListClientPricingRulesParams) Offset() int {
	if p.Page < 1 {
		p.Page = 1
	}
	return (p.Page - 1) * p.PageSize
}

// ---- Quantity Discounts ----

// CreateQuantityDiscountDTO datos para crear un descuento por cantidad
type CreateQuantityDiscountDTO struct {
	BusinessID      uint
	ProductID       *string
	MinQuantity     int
	DiscountPercent float64
	IsActive        bool
	Description     string
}

// UpdateQuantityDiscountDTO datos para actualizar un descuento por cantidad
type UpdateQuantityDiscountDTO struct {
	ID              uint
	BusinessID      uint
	ProductID       *string
	MinQuantity     int
	DiscountPercent float64
	IsActive        bool
	Description     string
}

// ListQuantityDiscountsParams parámetros de búsqueda y paginación
type ListQuantityDiscountsParams struct {
	BusinessID uint
	ProductID  *string // filtro opcional por producto
	Page       int
	PageSize   int
}

// Offset calcula el offset para paginación
func (p ListQuantityDiscountsParams) Offset() int {
	if p.Page < 1 {
		p.Page = 1
	}
	return (p.Page - 1) * p.PageSize
}

// ---- Price Calculator ----

// CalculatePriceRequest datos para calcular un precio
type CalculatePriceRequest struct {
	BusinessID uint
	ClientID   uint
	ProductID  string
	BasePrice  float64
	Quantity   int
}
