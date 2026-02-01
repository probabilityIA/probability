package entities

import "time"

// ProbabilityOrderItem representa un item de la orden que se guarda en la base de datos
// ✅ ENTIDAD PURA - SIN TAGS
type ProbabilityOrderItem struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time

	OrderID string

	ProductID    *string
	ProductSKU   string
	ProductName  string
	ProductTitle string
	VariantID    *string

	Quantity   int
	UnitPrice  float64
	TotalPrice float64
	Currency   string

	Discount float64
	Tax      float64
	TaxRate  *float64

	// Precios en moneda presentment (presentment_money - moneda local)
	UnitPricePresentment  float64
	TotalPricePresentment float64
	DiscountPresentment   float64
	TaxPresentment        float64

	ImageURL          *string
	ProductURL        *string
	Weight            *float64
	RequiresShipping  bool
	IsGiftCard        bool
	FulfillmentStatus *string
	Metadata          []byte
}
