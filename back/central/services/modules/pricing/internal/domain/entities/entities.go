package entities

import "time"

type ClientGroup struct {
	ID          uint
	BusinessID  uint
	Name        string
	Description string
	Color       string
	IsActive    bool
	MemberCount int64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type ClientSummary struct {
	ID        uint
	Name      string
	Email     string
	Phone     string
	Dni       string
	GroupID   *uint
	GroupName string
}

type CatalogPriceRow struct {
	ProductID   string
	ProductName string
	ProductSKU  string
	ImageURL    string
	Currency    string
	BasePrice   float64
	CustomPrice *float64
}

type EffectivePrice struct {
	ProductID  string
	BasePrice  float64
	FinalPrice float64
	Source     string
	GroupID    *uint
	GroupName  string
}
