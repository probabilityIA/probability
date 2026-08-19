package domain

import (
	"github.com/secamc93/probability/back/central/shared/productmatch"
)

type Integration struct {
	ID                uint
	BusinessID        *uint
	Name              string
	StoreID           string
	IntegrationType   int
	Config            map[string]interface{}
	IsTesting         bool
	BaseURL           string
	BaseURLTest       string
	ProductMatchRules []productmatch.Rule
}

type Credential struct {
	AccessToken string
	StoreID     string
	BaseURL     string
	UserAgent   string
}

type StoreInfo struct {
	ID       int64
	Name     string
	URL      string
	Country  string
	Currency string
	Language string
}

type TiendanubeVariant struct {
	ID              int64
	ProductID       int64
	SKU             string
	Barcode         string
	Price           float64
	Stock           int
	StockManagement bool
	Weight          float64
	Depth           float64
	Width           float64
	Height          float64
}

type TiendanubeProduct struct {
	ID          int64
	Name        string
	Description string
	ImageURL    string
	Published   bool
	Variants    []TiendanubeVariant
}

type StockTarget struct {
	ProductID int64
	VariantID int64
	Found     bool
}
