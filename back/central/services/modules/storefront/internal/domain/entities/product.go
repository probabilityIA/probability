package entities

import "time"

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
	TrackInventory   bool
	Category         string
	Brand            string
	Status           string
	IsFeatured       bool
	CreatedAt        time.Time
}
