package ports

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/storefront/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/storefront/internal/domain/entities"
)

// IRepository defines the storefront repository interface
type IRepository interface {
	// Catalog
	ListActiveProducts(ctx context.Context, businessID uint, filters dtos.CatalogFilters) ([]entities.StorefrontProduct, int64, error)
	GetProductByID(ctx context.Context, businessID uint, productID string) (*entities.StorefrontProduct, error)

	// Orders
	ListOrdersByUserID(ctx context.Context, businessID, userID uint, page, pageSize int) ([]entities.StorefrontOrder, int64, error)
	GetOrderByIDAndUserID(ctx context.Context, orderID string, businessID, userID uint) (*entities.StorefrontOrder, error)

	// Client
	GetClientByUserID(ctx context.Context, businessID, userID uint) (*entities.StorefrontClient, error)

	// Registration
	GetBusinessByCode(ctx context.Context, code string) (*entities.StorefrontBusiness, error)
	CreateUser(ctx context.Context, user *entities.NewUser) (uint, error)
	CreateBusinessStaff(ctx context.Context, userID, businessID, roleID uint) error
	CreateClient(ctx context.Context, client *entities.StorefrontClient) error
	GetClienteFinalRoleID(ctx context.Context) (uint, error)
	UserExistsByEmail(ctx context.Context, email string) (bool, error)

	// Validation
	GetRoleLevelByUserAndBusiness(ctx context.Context, userID, businessID uint) (int, error)

	// Integration
	GetPlatformIntegrationID(ctx context.Context, businessID uint) (uint, error)
}

// IStorefrontPublisher publishes storefront orders to RabbitMQ
type IStorefrontPublisher interface {
	PublishOrder(ctx context.Context, order []byte) error
}
