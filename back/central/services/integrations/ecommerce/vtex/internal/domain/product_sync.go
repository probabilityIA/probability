package domain

import "context"

const (
	DirectionToVTEX        = "to_vtex"
	DirectionToProbability = "to_probability"
)

const (
	ModeCreate = "create"
	ModeUpdate = "update"
)

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
	OnlyInVTEX           []ProductBrief
	ProbabilityNoSKU     int
	VTEXNoSKU            int
}

type MappedItem struct {
	ProductID      string
	SKU            string
	ExternalItemID string
}

type VTEXSKU struct {
	ID            string
	ProductID     string
	Name          string
	RefID         string
	EAN           string
	Price         float64
	IsActive      bool
	Weight        *float64
	Length        *float64
	Width         *float64
	Height        *float64
	MeasurementUnit string
}

type Warehouse struct {
	ID   string
	Name string
	IsActive bool
}

type WarehousesInfo struct {
	Warehouses []Warehouse
}

type WarehouseMapping struct {
	InternalWarehouseID uint
	VTEXWarehouseID     string
}

type InventoryConfig struct {
	Enabled          bool
	IsSeller         bool
	WarehouseMappings []WarehouseMapping
}

func (c InventoryConfig) WarehouseGroups() map[string][]uint {
	groups := make(map[string][]uint)
	for _, m := range c.WarehouseMappings {
		if m.VTEXWarehouseID == "" || m.InternalWarehouseID == 0 {
			continue
		}
		groups[m.VTEXWarehouseID] = append(groups[m.VTEXWarehouseID], m.InternalWarehouseID)
	}
	return groups
}

func (c InventoryConfig) InternalWarehouseIDs() []uint {
	ids := make([]uint, 0, len(c.WarehouseMappings))
	seen := make(map[uint]bool)
	for _, m := range c.WarehouseMappings {
		if m.InternalWarehouseID == 0 || seen[m.InternalWarehouseID] {
			continue
		}
		seen[m.InternalWarehouseID] = true
		ids = append(ids, m.InternalWarehouseID)
	}
	return ids
}

type IProductRepository interface {
	ListProductsByBusiness(ctx context.Context, businessID uint) ([]ProductForSync, error)
	GetExternalProductID(ctx context.Context, productID string, integrationID uint) (string, bool, error)
	UpsertProductIntegrationMapping(ctx context.Context, productID string, businessID, integrationID uint, externalProductID string) error
	ListMappedItems(ctx context.Context, integrationID uint) ([]MappedItem, error)
	GetStockForProducts(ctx context.Context, productIDs []string, warehouseIDs []uint) (map[string]int, error)
}
