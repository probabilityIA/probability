package ports

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/pay"
	"github.com/secamc93/probability/back/central/services/modules/publicsite/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/publicsite/internal/domain/entities"
)

type IRepository interface {
	GetBusinessBySlug(ctx context.Context, slug string) (*entities.BusinessPage, error)
	ListActiveProducts(ctx context.Context, businessID uint, filters dtos.CatalogFilters) ([]entities.PublicProduct, int64, error)
	GetProductByID(ctx context.Context, businessID uint, productID string) (*entities.PublicProduct, error)
	GetFeaturedProducts(ctx context.Context, businessID uint, limit int) ([]entities.PublicProduct, error)
	SaveContactSubmission(ctx context.Context, businessID uint, dto *dtos.ContactFormDTO) error

	IsIntegrationActive(ctx context.Context, businessID uint, integrationTypeID uint) (bool, error)

	GetTiendaIntegrationID(ctx context.Context, businessID uint) (uint, error)

	CreateCheckout(ctx context.Context, checkout *entities.PublicCheckout) error
	GetCheckoutByReference(ctx context.Context, reference string) (*entities.PublicCheckout, error)

	GetClientSession(ctx context.Context, businessID uint, userID uint) (*entities.TiendaSession, error)
	GetUserBasic(ctx context.Context, userID uint) (*entities.TiendaSession, error)
	UserBelongsToBusiness(ctx context.Context, businessID uint, userID uint) (bool, error)
	GetCustomerPrices(ctx context.Context, businessID uint, clientID uint) (map[string]float64, error)
	ListCustomerOrders(ctx context.Context, businessID, userID, customerID uint, page, pageSize int) ([]entities.TiendaOrder, int64, error)

	UpdateClientProfile(ctx context.Context, businessID, userID uint, name, phone, dni string) error
	UpdateUserAvatar(ctx context.Context, userID uint, avatarURL string) error
	GetUserAvatar(ctx context.Context, userID uint) (string, error)
	ListCustomerAddresses(ctx context.Context, businessID, customerID uint) ([]entities.CustomerAddress, error)
	CreateCustomerAddress(ctx context.Context, businessID, customerID uint, addr *entities.CustomerAddress) error
	SetPrimaryAddress(ctx context.Context, businessID, customerID, addressID uint) error
	DeleteCustomerAddress(ctx context.Context, businessID, customerID, addressID uint) error

	GetHiddenCategories(ctx context.Context, businessID uint) ([]string, error)
	ListPublicCategories(ctx context.Context, businessID uint, hidden []string) ([]string, error)
}

type IBoldGateway interface {
	BoldGenerateSignatureForReference(ctx context.Context, businessID uint, amount float64, currency, referencePrefix string) (*pay.BoldSignature, error)
	PublishAgreedStorefrontOrder(ctx context.Context, reference string) error
}
