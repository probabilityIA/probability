package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/customers/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/customers/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/customers/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

type IUseCase interface {
	CreateClient(ctx context.Context, dto dtos.CreateClientDTO) (*entities.Client, error)
	GetClient(ctx context.Context, businessID, clientID uint) (*entities.Client, error)
	ListClients(ctx context.Context, params dtos.ListClientsParams) ([]entities.Client, int64, error)
	UpdateClient(ctx context.Context, dto dtos.UpdateClientDTO) (*entities.Client, error)
	DeleteClient(ctx context.Context, businessID, clientID uint) error

	GetCustomerSummary(ctx context.Context, businessID, customerID uint) (*entities.CustomerSummary, error)
	ListCustomerAddresses(ctx context.Context, params dtos.ListCustomerAddressesParams) ([]entities.CustomerAddress, int64, error)
	ListCustomerProducts(ctx context.Context, params dtos.ListCustomerProductsParams) ([]entities.CustomerProductHistory, int64, error)
	ListCustomerOrderItems(ctx context.Context, params dtos.ListCustomerOrderItemsParams) ([]entities.CustomerOrderItem, int64, error)
	ProcessOrderEvent(ctx context.Context, event dtos.OrderEventDTO) error
}

type UseCase struct {
	repo ports.IRepository
	log  log.ILogger
}

func New(repo ports.IRepository, logger log.ILogger) IUseCase {
	return &UseCase{repo: repo, log: logger}
}
