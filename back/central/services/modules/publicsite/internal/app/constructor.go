package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/publicsite/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/publicsite/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/publicsite/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

const tiendaWebIntegrationTypeID = 31

type IUseCase interface {
	GetBusinessPage(ctx context.Context, slug string) (*entities.BusinessPage, error)
	ListCatalog(ctx context.Context, slug string, filters dtos.CatalogFilters) ([]entities.PublicProduct, int64, error)
	GetProduct(ctx context.Context, slug string, productID string) (*entities.PublicProduct, error)
	GetFeaturedProducts(ctx context.Context, slug string, limit int) ([]entities.PublicProduct, error)
	SubmitContact(ctx context.Context, slug string, dto *dtos.ContactFormDTO) error
	CreateCheckoutSession(ctx context.Context, slug string, dto *dtos.CreateCheckoutDTO, userID uint) (*dtos.CheckoutSessionDTO, error)
	GetCheckoutStatus(ctx context.Context, reference string) (string, error)
	GetSession(ctx context.Context, slug string, userID uint) (*entities.TiendaSession, error)
	GetMyPrices(ctx context.Context, slug string, userID uint) (map[string]float64, error)
	GetMyOrders(ctx context.Context, slug string, userID uint, page, pageSize int) ([]entities.TiendaOrder, int64, error)
	UpdateProfile(ctx context.Context, slug string, userID uint, name, phone, dni string) error
	SaveAvatar(ctx context.Context, slug string, userID uint, avatarURL string) error
	ListAddresses(ctx context.Context, slug string, userID uint) ([]entities.CustomerAddress, error)
	AddAddress(ctx context.Context, slug string, userID uint, addr *entities.CustomerAddress) error
	SetPrimaryAddress(ctx context.Context, slug string, userID uint, addressID uint) error
	DeleteAddress(ctx context.Context, slug string, userID uint, addressID uint) error
	GetCategories(ctx context.Context, slug string) ([]string, error)
}

type UseCase struct {
	repo   ports.IRepository
	orders ports.IBoldGateway
	logger log.ILogger
}

func New(repo ports.IRepository, orders ports.IBoldGateway, logger log.ILogger) IUseCase {
	return &UseCase{repo: repo, orders: orders, logger: logger}
}
