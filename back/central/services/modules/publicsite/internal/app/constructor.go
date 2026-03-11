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
}

type UseCase struct {
	repo   ports.IRepository
	logger log.ILogger
}

func New(repo ports.IRepository, logger log.ILogger) IUseCase {
	return &UseCase{repo: repo, logger: logger}
}
