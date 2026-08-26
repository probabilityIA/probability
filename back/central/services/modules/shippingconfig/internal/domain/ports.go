package domain

import "context"

type IRepository interface {
	GetConfig(ctx context.Context, businessID uint, warehouseID *uint) (*ShippingConfig, error)
	ListConfigs(ctx context.Context, businessID uint) ([]ShippingConfig, error)
	UpsertConfig(ctx context.Context, cfg *ShippingConfig) error
	DeleteConfig(ctx context.Context, businessID uint, warehouseID uint) error
	ListWarehouseOrigins(ctx context.Context, businessID uint) ([]WarehouseOrigin, error)
	SetDefaultWarehouse(ctx context.Context, businessID uint, warehouseID uint) error
}

type IUseCase interface {
	GetOverview(ctx context.Context, businessID uint) (*Overview, error)
	SaveBusinessConfig(ctx context.Context, businessID uint, req SaveConfigRequest) (*ShippingConfig, error)
	SaveWarehouseConfig(ctx context.Context, businessID, warehouseID uint, req SaveConfigRequest) (*ShippingConfig, error)
	RemoveWarehouseConfig(ctx context.Context, businessID, warehouseID uint) error
	SetDefaultWarehouse(ctx context.Context, businessID, warehouseID uint) error
	Resolve(ctx context.Context, businessID uint, warehouseID *uint) (*ShippingConfig, error)
}

type Overview struct {
	Business   *ShippingConfig   `json:"business"`
	Warehouses []WarehouseOrigin `json:"warehouses"`
	Overrides  []ShippingConfig  `json:"overrides"`
	Carriers   []Carrier         `json:"carriers"`
}

type SaveConfigRequest struct {
	PackageStrategy string           `json:"package_strategy"`
	Boxes           []Box            `json:"boxes"`
	Carriers        []CarrierSetting `json:"carriers"`
	AlwaysInsure    bool             `json:"always_insure"`
	UpdatedBy       uint             `json:"-"`
	UpdatedByName   string           `json:"-"`
}
