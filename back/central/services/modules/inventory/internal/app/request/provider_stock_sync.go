package request

type ProviderStockSyncDTO struct {
	CorrelationID string                  `json:"correlation_id"`
	BusinessID    uint                    `json:"business_id"`
	IntegrationID uint                    `json:"integration_id"`
	Provider      string                  `json:"provider"`
	Items         []ProviderStockSyncItem `json:"items"`
}

type ProviderStockSyncItem struct {
	SKU         string `json:"sku"`
	WarehouseID uint   `json:"warehouse_id"`
	Quantity    int    `json:"quantity"`
}
