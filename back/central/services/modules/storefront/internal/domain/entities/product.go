package entities

import "time"

// StorefrontProduct represents a product visible to storefront customers
type StorefrontProduct struct {
	ID               string
	Name             string
	Description      string
	ShortDescription string
	Price            float64
	CompareAtPrice   *float64
	Currency         string
	ImageURL         string
	Images           []byte
	SKU              string
	StockQuantity    int
	Category         string
	Brand            string
	Status           string
	IsFeatured       bool
	CreatedAt        time.Time
}
