package domain

import (
	"context"
	"time"

	"github.com/secamc93/probability/back/central/shared/inventorycompare"
	"github.com/secamc93/probability/back/central/shared/productmatch"
)

type ChannelStock struct {
	ExternalID        string
	AvailableQuantity int
	ManageStock       bool
	Found             bool
}

type CreateProductInput struct {
	Name          string
	SKU           string
	Price         float64
	Description   string
	StockQuantity int
	ManageStock   bool
	ImageURL      string
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
}

type WooProduct struct {
	ID                string
	ParentID          string
	ParentName        string
	VariantLabel      string
	VariantAttributes map[string]string
	SKU               string
	Barcode           string
	Name              string
	Description       string
	Category          string
	Brand             string
	Weight            *float64
	Length            *float64
	Width             *float64
	Height            *float64
	Price             float64
	StockQuantity     int
	ImageURL          string
}

func (p ProductForSync) MatchItem() productmatch.Item {
	return productmatch.Item{
		SKU:        p.SKU,
		Barcode:    p.Barcode,
		ExternalID: p.ExternalID,
		Name:       p.Name,
	}
}

func (w WooProduct) MatchItem() productmatch.Item {
	externalID := w.ID
	variantID := ""
	if w.ParentID != "" {
		externalID = w.ParentID
		variantID = w.ID
	}
	return productmatch.Item{
		SKU:        w.SKU,
		Barcode:    w.Barcode,
		ExternalID: externalID,
		VariantID:  variantID,
		Name:       w.Name,
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
	OnlyInWoo            []ProductBrief
	ProbabilityNoSKU     int
	WooNoSKU             int
	TypoSuspects         []productmatch.TypoSuspect
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
	Mode              string
	SingleWarehouseID uint
	WarehouseIDs      []uint
}

type IProductRepository interface {
	ListProductsByBusiness(ctx context.Context, businessID uint) ([]ProductForSync, error)
	GetExternalProductID(ctx context.Context, productID string, integrationID uint) (string, bool, error)
	UpsertProductIntegrationMapping(ctx context.Context, productID string, businessID, integrationID uint, refs productmatch.ExternalRefs) error
	SaveChannelSnapshots(ctx context.Context, businessID, integrationID uint, entries []productmatch.SnapshotEntry) error
	ListMappedItems(ctx context.Context, integrationID uint) ([]MappedItem, error)
	GetStockForProducts(ctx context.Context, productIDs []string, warehouseIDs []uint) (map[string]int, error)
	ERPFeedsInventory(ctx context.Context, businessID uint) (bool, error)
	SaveCompareSnapshot(ctx context.Context, businessID, integrationID uint, rows []inventorycompare.Row, checkedAt time.Time) error
	LoadCompareSnapshot(ctx context.Context, businessID, integrationID uint, opts inventorycompare.LoadOptions) (*inventorycompare.Page, error)
}
