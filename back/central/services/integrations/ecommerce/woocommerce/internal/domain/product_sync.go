package domain

import "context"

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
	Name           string
	Description    string
	Price          float64
	StockQuantity  int
	TrackInventory bool
	ImageURL       string
}

type WooProduct struct {
	ID            string
	ParentID      string
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
	OnlyInWoo         []ProductBrief
	ProbabilityNoSKU  int
	WooNoSKU          int
}

type MappedItem struct {
	ProductID      string
	SKU            string
	ExternalItemID string
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
	UpsertProductIntegrationMapping(ctx context.Context, productID string, businessID, integrationID uint, externalProductID string) error
	ListMappedItems(ctx context.Context, integrationID uint) ([]MappedItem, error)
	GetStockForProducts(ctx context.Context, productIDs []string, warehouseIDs []uint) (map[string]int, error)
}
