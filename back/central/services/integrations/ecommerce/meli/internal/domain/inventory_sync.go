package domain

import "context"

type MappedItem struct {
	ProductID      string
	SKU            string
	ExternalItemID string
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
}
