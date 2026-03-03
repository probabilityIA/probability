package ports

import (
	"context"

	"github.com/secamc93/probability/back/testing/modules/orders/internal/domain/dtos"
	"github.com/secamc93/probability/back/testing/modules/orders/internal/domain/entities"
)

type IRepository interface {
	GetProducts(ctx context.Context, businessID uint) ([]entities.Product, error)
	GetIntegrations(ctx context.Context, businessID uint) ([]entities.Integration, error)
	GetPaymentMethods(ctx context.Context) ([]entities.PaymentMethod, error)
	GetOrderStatuses(ctx context.Context) ([]entities.OrderStatus, error)
}

type ICentralClient interface {
	CreateOrder(ctx context.Context, token string, orderPayload map[string]interface{}) (*entities.CreatedOrder, *entities.APICallLog, error)
}

type IUseCase interface {
	GetReferenceData(ctx context.Context, businessID uint) (*entities.ReferenceData, error)
	GenerateOrders(ctx context.Context, businessID uint, dto *dtos.GenerateOrdersDTO, token string) (*entities.GenerateResult, error)
}
