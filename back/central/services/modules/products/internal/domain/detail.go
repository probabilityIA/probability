package domain

import "time"

type ProductChannel struct {
	IntegrationID       uint       `json:"integration_id"`
	IntegrationName     string     `json:"integration_name"`
	IntegrationCode     string     `json:"integration_code"`
	IntegrationCategory string     `json:"integration_category"`
	IntegrationImageURL string     `json:"integration_image_url,omitempty"`
	IsActive            bool       `json:"is_active"`
	ExternalProductID   string     `json:"external_product_id"`
	ExternalVariantID   *string    `json:"external_variant_id,omitempty"`
	ExternalSKU         *string    `json:"external_sku,omitempty"`
	ExternalBarcode     *string    `json:"external_barcode,omitempty"`
	LogisticType        *string    `json:"logistic_type,omitempty"`
	LastPushedQty       *int       `json:"last_pushed_qty,omitempty"`
	LinkedAt            time.Time  `json:"linked_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
	InventorySyncAt     *time.Time `json:"inventory_sync_at,omitempty"`
	InventorySyncStatus string     `json:"inventory_sync_status,omitempty"`
	CompareAt           *time.Time `json:"compare_at,omitempty"`
	CompareStatus       string     `json:"compare_status,omitempty"`
	SKUMatches          bool       `json:"sku_matches"`
}

type ProductStockLevel struct {
	WarehouseID   uint      `json:"warehouse_id"`
	WarehouseName string    `json:"warehouse_name"`
	WarehouseCode string    `json:"warehouse_code"`
	Quantity      int       `json:"quantity"`
	ReservedQty   int       `json:"reserved_qty"`
	AvailableQty  int       `json:"available_qty"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type ProductMovement struct {
	CreatedAt       time.Time `json:"created_at"`
	TypeName        string    `json:"type_name"`
	Direction       string    `json:"direction"`
	Reason          string    `json:"reason"`
	Quantity        int       `json:"quantity"`
	PreviousQty     int       `json:"previous_qty"`
	NewQty          int       `json:"new_qty"`
	WarehouseName   string    `json:"warehouse_name,omitempty"`
	IntegrationName string    `json:"integration_name,omitempty"`
}

type ProductDetailResponse struct {
	ProductResponse

	Origins   []FieldOriginView   `json:"origins"`
	Channels  []ProductChannel    `json:"channels"`
	Stock     []ProductStockLevel `json:"stock_levels"`
	Movements []ProductMovement   `json:"movements"`

	StockTotal     int `json:"stock_total"`
	ReservedTotal  int `json:"reserved_total"`
	AvailableTotal int `json:"available_total"`

	SiblingVariants int64 `json:"sibling_variants"`
}
