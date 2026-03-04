package ports

import (
	"context"

	"github.com/secamc93/probability/back/testing/modules/orders/internal/domain/dtos"
	"github.com/secamc93/probability/back/testing/modules/orders/internal/domain/entities"
	sharedtypes "github.com/secamc93/probability/back/testing/shared/types"
)

type IRepository interface {
	GetProducts(ctx context.Context, businessID uint) ([]entities.Product, error)
	GetIntegrations(ctx context.Context, businessID uint) ([]entities.Integration, error)
	GetPaymentMethods(ctx context.Context) ([]entities.PaymentMethod, error)
	GetOrderStatuses(ctx context.Context) ([]entities.OrderStatus, error)
	GetIntegrationTypeCode(ctx context.Context, integrationID uint) (string, error)
	GetIntegrationCategoryID(ctx context.Context, integrationID uint) (uint, error)
	DeleteAllOrders(ctx context.Context, businessID uint) (int64, error)
}

type ICentralClient interface {
	CreateOrder(ctx context.Context, token string, orderPayload map[string]interface{}) (*entities.CreatedOrder, *entities.APICallLog, error)
	GetBaseURL() string
}

// IWebhookSimulator simulates ecommerce platform webhooks
type IWebhookSimulator interface {
	SimulateOrder(topic string) error
	BuildWebhookPayload(topic string, baseURL string) (*sharedtypes.WebhookPayload, error)
	GetWebhookTopics() []string
}

type IUseCase interface {
	GetReferenceData(ctx context.Context, businessID uint) (*entities.ReferenceData, error)
	GenerateOrders(ctx context.Context, businessID uint, dto *dtos.GenerateOrdersDTO, token string) (*entities.GenerateResult, error)
	DeleteAllOrders(ctx context.Context, businessID uint) (int64, error)
}
