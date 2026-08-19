package domain

import (
	"context"
	"time"

	"github.com/secamc93/probability/back/central/shared/inventorycompare"
	"github.com/secamc93/probability/back/central/shared/productmatch"
)

type CreateProductInput struct {
	Name          string
	SKU           string
	Barcode       string
	Price         float64
	Description   string
	StockQuantity int
	ManageStock   bool
	ImageURL      string
	Weight        *float64
	Height        *float64
	Width         *float64
	Length        *float64
}

type UpdateProductInput struct {
	Name        string
	Description string
}

type UpdateVariantInput struct {
	Price   *float64
	Weight  *float64
	Height  *float64
	Width   *float64
	Depth   *float64
	Barcode string
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
	Weight         *float64
	WeightUnit     string
	Length         *float64
	Width          *float64
	Height         *float64
	DimensionUnit  string
}

type ProductBrief struct {
	SKU          string
	Name         string
	MatchedBy    string
	MatchedValue string
}

func (p ProductForSync) MatchItem() productmatch.Item {
	return productmatch.Item{
		SKU:        p.SKU,
		Barcode:    p.Barcode,
		ExternalID: p.ExternalID,
		Name:       p.Name,
	}
}

type ReconcileResult struct {
	Matched              int
	MatchedItems         []ProductBrief
	MatchedNotAssociated []ProductBrief
	OnlyInProbability    []ProductBrief
	OnlyInTiendanube     []ProductBrief
	ProbabilityNoSKU     int
	TiendanubeNoSKU      int
	MatchRules           []productmatch.Rule
}

type MappedItem struct {
	ProductID         string
	SKU               string
	Name              string
	ImageURL          string
	Barcode           string
	ExternalItemID    string
	ExternalVariantID string
	ExternalSKU       string
	ExternalBarcode   string
}

type InventoryConfig struct {
	Enabled           bool
	SingleWarehouseID uint
}

type ChannelStock struct {
	ExternalID  string
	Quantity    int
	ManageStock bool
	Found       bool
}

type IProductRepository interface {
	ListProductsByBusiness(ctx context.Context, businessID uint) ([]ProductForSync, error)
	GetExternalProductID(ctx context.Context, productID string, integrationID uint) (string, bool, error)
	UpsertProductIntegrationMapping(ctx context.Context, productID string, businessID, integrationID uint, refs productmatch.ExternalRefs) error
	ListMappedItems(ctx context.Context, integrationID uint) ([]MappedItem, error)
	GetStockForProducts(ctx context.Context, productIDs []string, warehouseIDs []uint) (map[string]int, error)
	SaveCompareSnapshot(ctx context.Context, businessID, integrationID uint, rows []inventorycompare.Row, checkedAt time.Time) error
	LoadCompareSnapshot(ctx context.Context, businessID, integrationID uint, opts inventorycompare.LoadOptions) (*inventorycompare.Page, error)
}
