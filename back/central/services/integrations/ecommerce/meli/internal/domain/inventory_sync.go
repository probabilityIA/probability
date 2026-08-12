package domain

import "context"

const LogisticFulfillment = "fulfillment"

type PushState struct {
	LogisticType  string
	LastPushedQty *int
}

type MappedItem struct {
	ProductID         string
	SKU               string
	Barcode           string
	ExternalItemID    string
	ExternalVariantID string
	ExternalSKU       string
	ExternalBarcode   string
	LogisticType      string
	LastPushedQty     *int
}

type WarehouseMapping struct {
	InternalWarehouseID uint
	MLStoreID           string
	MLNetworkNodeID     string
}

type InventoryConfig struct {
	Enabled           bool
	Mode              string
	SingleWarehouseID uint
	WarehouseIDs      []uint
	WarehouseMappings []WarehouseMapping
}

type IInventoryRepository interface {
	ListMappedItems(ctx context.Context, integrationID uint) ([]MappedItem, error)
	GetStockForProducts(ctx context.Context, productIDs []string, warehouseIDs []uint) (map[string]int, error)
	GetInventoryByWarehouses(ctx context.Context, productIDs []string, warehouseIDs []uint) (map[string]map[uint]int, error)
	UpdateLogisticType(ctx context.Context, integrationID uint, externalItemID, logisticType string) error
	MarkPushedQty(ctx context.Context, integrationID uint, productID, variantID string, quantity int) error
	GetPushState(ctx context.Context, integrationID uint, productID, variantID string) (*PushState, error)
}
