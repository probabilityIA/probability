package entities

import "time"

// PublicProduct represents a product visible on the public page
type PublicProduct struct {
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
	IsFeatured       bool
	CreatedAt        time.Time
}
