package domain

import (
	"context"

	"github.com/secamc93/probability/back/central/shared/productmatch"
)

type ProductForSync struct {
	ID             string
	SKU            string
	Barcode        string
	ExternalID     string
	Name           string
	Description    string
	Price          float64
	StockQuantity  int
	TrackInventory bool
	ImageURL       string
}

func (p ProductForSync) MatchItem() productmatch.Item {
	return productmatch.Item{
		SKU:        p.SKU,
		Barcode:    p.Barcode,
		ExternalID: p.ExternalID,
		Name:       p.Name,
	}
}

type ProductBrief struct {
	SKU          string
	Name         string
	MatchedBy    string
	MatchedValue string
}

type ReconcileResult struct {
	Matched              int
	MatchedItems         []ProductBrief
	MatchedNotAssociated []ProductBrief
	OnlyInProbability    []ProductBrief
	OnlyInShopify        []ProductBrief
	ProbabilityNoSKU     int
	ShopifyNoSKU         int
	MatchRules           []productmatch.Rule
}

type ShopifyProductForSync struct {
	ProductID string
	VariantID string
	SKU       string
	Barcode   string
	Name      string
	ImageURL  string
}

func (p ShopifyProductForSync) MatchItem() productmatch.Item {
	return productmatch.Item{
		SKU:        p.SKU,
		Barcode:    p.Barcode,
		ExternalID: p.ProductID,
		VariantID:  p.VariantID,
		Name:       p.Name,
	}
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
	UpsertProductIntegrationMapping(ctx context.Context, productID string, businessID, integrationID uint, refs productmatch.ExternalRefs) error
}
