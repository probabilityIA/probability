package domain

import "context"

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

type ProductBrief struct {
	SKU  string
	Name string
}

type ReconcileResult struct {
	Matched              int
	MatchedNotAssociated []ProductBrief
	OnlyInProbability    []ProductBrief
	OnlyInShopify        []ProductBrief
	ProbabilityNoSKU     int
	ShopifyNoSKU         int
}

type ShopifyProductForSync struct {
	ProductID string
	SKU       string
	Name      string
}

type CreateProductInput struct {
	Name          string
	SKU           string
	Price         float64
	Description   string
	StockQuantity int
}

type IProductRepository interface {
	ListProductsByBusiness(ctx context.Context, businessID uint) ([]ProductForSync, error)
	GetExternalProductID(ctx context.Context, productID string, integrationID uint) (string, bool, error)
	UpsertProductIntegrationMapping(ctx context.Context, productID string, businessID, integrationID uint, externalProductID string) error
}
