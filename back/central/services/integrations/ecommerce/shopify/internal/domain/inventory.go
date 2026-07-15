package domain

import "context"

type MappedItem struct {
	ProductID      string
	SKU            string
	ExternalItemID string
}

type WarehouseMapping struct {
	InternalWarehouseID uint
	ShopifyLocationID   int64
}

type InventoryConfig struct {
	Enabled           bool
	Mode              string
	SingleWarehouseID uint
	WarehouseIDs      []uint
	DefaultLocationID int64
	LocationMappings  []WarehouseMapping
}

type ShopifyLocation struct {
	ID   int64
	Name string
}

type ShopifyVariant struct {
	ID              int64
	SKU             string
	InventoryItemID int64
}

type ShopifyProduct struct {
	ID       int64
	Variants []ShopifyVariant
}

type IInventoryRepository interface {
	ListMappedItems(ctx context.Context, integrationID uint) ([]MappedItem, error)
	GetStockForProducts(ctx context.Context, productIDs []string, warehouseIDs []uint) (map[string]int, error)
	GetInventoryByWarehouses(ctx context.Context, productIDs []string, warehouseIDs []uint) (map[string]map[uint]int, error)
}
