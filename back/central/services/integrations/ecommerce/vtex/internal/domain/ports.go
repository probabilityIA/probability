package domain

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/canonical"
)

type Credential struct {
	AccountName string
	AppKey      string
	AppToken    string
}

type IVTEXClient interface {
	TestConnection(ctx context.Context, cred Credential) error

	GetOrders(ctx context.Context, cred Credential, page, perPage int, filters map[string]string) (*VTEXOrderListResponse, error)
	GetOrderByID(ctx context.Context, cred Credential, orderID string) (*VTEXOrder, []byte, error)

	ListSKUs(ctx context.Context, cred Credential) ([]VTEXSKU, error)
	GetSKUIDByRefID(ctx context.Context, cred Credential, refID string, isSeller bool) (string, error)

	GetWarehouses(ctx context.Context, cred Credential) ([]Warehouse, error)
	UpdateSKUInventory(ctx context.Context, cred Credential, skuID, warehouseID string, quantity int) error

	GetOrderHook(ctx context.Context, cred Credential) (*HookConfig, error)
	SetOrderHook(ctx context.Context, cred Credential, url, hookKey string) error
	DeleteOrderHook(ctx context.Context, cred Credential) error
}

type IIntegrationService interface {
	GetIntegrationByID(ctx context.Context, integrationID string) (*Integration, error)
	DecryptCredential(ctx context.Context, integrationID string, fieldName string) (string, error)
	UpdateIntegrationConfig(ctx context.Context, integrationID string, config map[string]interface{}) error
}

type OrderPublisher interface {
	Publish(ctx context.Context, order *canonical.ProbabilityOrderDTO) error
}
