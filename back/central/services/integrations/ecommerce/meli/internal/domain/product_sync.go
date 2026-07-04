package domain

import "context"

type CreateProductInput struct {
	Name          string
	SKU           string
	Price         float64
	Description   string
	StockQuantity int
}

type ProductForSync struct {
	ID             string
	SKU            string
	Name           string
	Description    string
	Price          float64
	StockQuantity  int
	TrackInventory bool
	ImageURL       string
}

type MeliProduct struct {
	ID            string
	SKU           string
	Name          string
	Price         float64
	StockQuantity int
}

type ProductBrief struct {
	SKU  string
	Name string
}

type ReconcileResult struct {
	Matched           int
	OnlyInProbability []ProductBrief
	OnlyInMeli        []ProductBrief
	ProbabilityNoSKU  int
	MeliNoSKU         int
}

type IProductRepository interface {
	ListProductsByBusiness(ctx context.Context, businessID uint) ([]ProductForSync, error)
	GetExternalProductID(ctx context.Context, productID string, integrationID uint) (string, bool, error)
	UpsertProductIntegrationMapping(ctx context.Context, productID string, businessID, integrationID uint, externalProductID string) error
}
