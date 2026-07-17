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
	Weight        *float64
	Height        *float64
	Width         *float64
	Length        *float64
}

type UpdateProductInput struct {
	Name        string
	Price       *float64
	Description string
	Weight      *float64
	Height      *float64
	Width       *float64
	Length      *float64
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
	Weight         *float64
	WeightUnit     string
	Length         *float64
	Width          *float64
	Height         *float64
	DimensionUnit  string
}

type ProductBrief struct {
	SKU  string
	Name string
}

type ReconcileResult struct {
	Matched              int
	MatchedNotAssociated []ProductBrief
	OnlyInProbability    []ProductBrief
	OnlyInJumpseller     []ProductBrief
	ProbabilityNoSKU     int
	JumpsellerNoSKU      int
}

type MappedItem struct {
	ProductID      string
	SKU            string
	ExternalItemID string
}

const (
	InventoryModeSingle = "single"
	InventoryModeMapped = "mapped"
)

type InventoryConfig struct {
	Enabled           bool
	Mode              string
	SingleWarehouseID uint
	DefaultLocationID int64
	LocationMappings  []LocationMapping
}

func (c InventoryConfig) LocationGroups() map[int64][]uint {
	groups := make(map[int64][]uint)
	for _, m := range c.LocationMappings {
		groups[m.JumpsellerLocationID] = append(groups[m.JumpsellerLocationID], m.InternalWarehouseID)
	}
	return groups
}

type LocationMapping struct {
	InternalWarehouseID  uint
	JumpsellerLocationID int64
}

type IProductRepository interface {
	ListProductsByBusiness(ctx context.Context, businessID uint) ([]ProductForSync, error)
	GetExternalProductID(ctx context.Context, productID string, integrationID uint) (string, bool, error)
	UpsertProductIntegrationMapping(ctx context.Context, productID string, businessID, integrationID uint, externalProductID string) error
	ListMappedItems(ctx context.Context, integrationID uint) ([]MappedItem, error)
	GetStockForProducts(ctx context.Context, productIDs []string, warehouseIDs []uint) (map[string]int, error)
}
