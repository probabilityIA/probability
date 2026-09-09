package ports

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/storefront/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/storefront/internal/domain/entities"
)

type IRepository interface {
	ListActiveProducts(ctx context.Context, businessID uint, filters dtos.CatalogFilters) ([]entities.StorefrontProduct, int64, error)
	GetProductByID(ctx context.Context, businessID uint, productID string) (*entities.StorefrontProduct, error)
	GetCatalogFilters(ctx context.Context, businessID uint) (entities.StorefrontFilters, error)

	ListOrdersByUserID(ctx context.Context, businessID, userID uint, page, pageSize int) ([]entities.StorefrontOrder, int64, error)
	GetOrderByIDAndUserID(ctx context.Context, orderID string, businessID, userID uint) (*entities.StorefrontOrder, error)

	GetClientByUserID(ctx context.Context, businessID, userID uint) (*entities.StorefrontClient, error)
	ListClientsByBusiness(ctx context.Context, businessID uint, page, pageSize int) ([]entities.StorefrontClient, int64, error)

	GetBusinessByCode(ctx context.Context, code string) (*entities.StorefrontBusiness, error)
	CreateUser(ctx context.Context, user *entities.NewUser) (uint, error)
	CreateBusinessStaff(ctx context.Context, userID, businessID, roleID uint) error
	CreateClient(ctx context.Context, client *entities.StorefrontClient) error
	GetClienteFinalRoleID(ctx context.Context) (uint, error)
	UserExistsByEmail(ctx context.Context, email string) (bool, error)
	GetUserAuthByEmail(ctx context.Context, email string) (uint, string, error)
	VerifyPassword(hash, password string) bool
	StaffExists(ctx context.Context, userID, businessID uint) (bool, error)
	GetClientByBusinessAndEmail(ctx context.Context, businessID uint, email string) (*entities.StorefrontClient, error)
	LinkClientUser(ctx context.Context, clientID, userID uint) error

	GetRoleLevelByUserAndBusiness(ctx context.Context, userID, businessID uint) (int, error)

	GetPlatformIntegrationID(ctx context.Context, businessID uint) (uint, error)

	IsIntegrationActiveOrMissing(ctx context.Context, businessID uint, integrationTypeID uint) (bool, error)
	GetCatalogLayout(ctx context.Context, businessID, integrationTypeID uint) (entities.StorefrontCatalogLayout, error)
	UpdateCatalogLayout(ctx context.Context, businessID, integrationTypeID, requesterUserID uint, layout entities.StorefrontCatalogLayout) error
	GetCatalogBanner(ctx context.Context, businessID, integrationTypeID uint) (entities.StorefrontBanner, error)
	UpdateCatalogBanner(ctx context.Context, businessID, integrationTypeID, requesterUserID uint, banner entities.StorefrontBanner) error
}

type IStorefrontPublisher interface {
	PublishOrder(ctx context.Context, order []byte) error
}
