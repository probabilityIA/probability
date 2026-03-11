package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/storefront/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/storefront/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/storefront/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

const tiendaIntegrationTypeID = 30

// IUseCase defines the storefront use cases
type IUseCase interface {
	ListCatalog(ctx context.Context, businessID uint, filters dtos.CatalogFilters) ([]entities.StorefrontProduct, int64, error)
	GetProduct(ctx context.Context, businessID uint, productID string) (*entities.StorefrontProduct, error)
	CreateOrder(ctx context.Context, businessID, userID uint, dto *dtos.StorefrontCreateOrderDTO) error
	ListMyOrders(ctx context.Context, businessID, userID uint, page, pageSize int) ([]entities.StorefrontOrder, int64, error)
	GetMyOrder(ctx context.Context, orderID string, businessID, userID uint) (*entities.StorefrontOrder, error)
	Register(ctx context.Context, dto *dtos.RegisterDTO) error
}

// UseCase implements IUseCase
type UseCase struct {
	repo      ports.IRepository
	publisher ports.IStorefrontPublisher
	logger    log.ILogger
}

// New creates a new storefront use case
func New(repo ports.IRepository, logger log.ILogger, publisher ports.IStorefrontPublisher) IUseCase {
	return &UseCase{
		repo:      repo,
		publisher: publisher,
		logger:    logger,
	}
}
