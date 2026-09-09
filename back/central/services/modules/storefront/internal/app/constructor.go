package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/storefront/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/storefront/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/storefront/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

const tiendaIntegrationTypeID = 30

type IUseCase interface {
	ListCatalog(ctx context.Context, businessID uint, filters dtos.CatalogFilters) ([]entities.StorefrontProduct, int64, error)
	GetCatalogFilters(ctx context.Context, businessID uint) (entities.StorefrontFilters, error)
	UpdateCatalogLayout(ctx context.Context, businessID, requesterUserID uint, columns, rows int) (entities.StorefrontCatalogLayout, error)
	UpdateCatalogBanner(ctx context.Context, businessID, requesterUserID uint, enabled *bool, imageURL *string) (entities.StorefrontBanner, error)
	GetProduct(ctx context.Context, businessID uint, productID string) (*entities.StorefrontProduct, error)
	CreateOrder(ctx context.Context, businessID, userID uint, dto *dtos.StorefrontCreateOrderDTO) error
	ListMyOrders(ctx context.Context, businessID, userID uint, page, pageSize int) ([]entities.StorefrontOrder, int64, error)
	GetMyOrder(ctx context.Context, orderID string, businessID, userID uint) (*entities.StorefrontOrder, error)
	CreateClient(ctx context.Context, businessID, requesterUserID uint, dto *dtos.CreateClientDTO) (*entities.StorefrontClient, string, error)
	ListClients(ctx context.Context, businessID uint, page, pageSize int) ([]entities.StorefrontClient, int64, error)
	Register(ctx context.Context, dto *dtos.RegisterDTO) error
}

type UseCase struct {
	repo      ports.IRepository
	publisher ports.IStorefrontPublisher
	logger    log.ILogger
}

func New(repo ports.IRepository, logger log.ILogger, publisher ports.IStorefrontPublisher) IUseCase {
	return &UseCase{
		repo:      repo,
		publisher: publisher,
		logger:    logger,
	}
}
