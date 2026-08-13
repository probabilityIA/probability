package domain

import (
	"context"

	"github.com/secamc93/probability/back/central/shared/productmatch"
)

type CreateProductInput struct {
	Name          string
	SKU           string
	Price         float64
	Description   string
	StockQuantity int
	SiteID        string
	CurrencyID    string
	ListingTypeID string
	Brand         string
	ImageURL      string
	CategoryID    string
	VariantAttrs  map[string]string
}

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
	Brand          string
	Category       string
	MeliCategoryID string
	VariantAttrs   map[string]string
}

type MeliProduct struct {
	ID            string
	VariationID   string
	ParentName    string
	VariantLabel  string
	LogisticType  string
	SKU           string
	Barcode       string
	Name          string
	Brand         string
	CategoryID    string
	Category      string
	Color         string
	Size          string
	Price         float64
	StockQuantity int
	ImageURL      string
}

func (p ProductForSync) MatchItem() productmatch.Item {
	return productmatch.Item{
		SKU:        p.SKU,
		Barcode:    p.Barcode,
		ExternalID: p.ExternalID,
		Name:       p.Name,
	}
}

func (m MeliProduct) MatchItem() productmatch.Item {
	return productmatch.Item{
		SKU:        m.SKU,
		Barcode:    m.Barcode,
		ExternalID: m.ID,
		VariantID:  m.VariationID,
		Name:       m.Name,
	}
}

type ProductBrief struct {
	SKU          string
	Name         string
	MatchedBy    string
	MatchedValue string
	ParentRef    string
	ParentLabel  string
	VariantLabel string
}

type ReconcileResult struct {
	Matched              int
	MatchedItems         []ProductBrief
	MatchedNotAssociated []ProductBrief
	OnlyInProbability    []ProductBrief
	OnlyInMeli           []ProductBrief
	ProbabilityNoSKU     int
	MeliNoSKU            int
	MeliNoSKUItems       []ProductBrief
	SKUChangedItems      []ProductBrief
	TypoSuspects         []productmatch.TypoSuspect
	MatchRules           []productmatch.Rule
}

type IProductRepository interface {
	ListProductsByBusiness(ctx context.Context, businessID uint) ([]ProductForSync, error)
	GetExternalProductID(ctx context.Context, productID string, integrationID uint) (string, bool, error)
	UpsertProductIntegrationMapping(ctx context.Context, productID string, businessID, integrationID uint, refs productmatch.ExternalRefs) error
	SaveChannelSnapshots(ctx context.Context, businessID, integrationID uint, entries []productmatch.SnapshotEntry) error
	ERPFeedsInventory(ctx context.Context, businessID uint) (bool, error)
}
