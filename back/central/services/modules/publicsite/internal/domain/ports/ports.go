package ports

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/publicsite/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/publicsite/internal/domain/entities"
)

type IRepository interface {
	GetBusinessBySlug(ctx context.Context, slug string) (*entities.BusinessPage, error)
	ListActiveProducts(ctx context.Context, businessID uint, filters dtos.CatalogFilters) ([]entities.PublicProduct, int64, error)
	GetProductByID(ctx context.Context, businessID uint, productID string) (*entities.PublicProduct, error)
	GetFeaturedProducts(ctx context.Context, businessID uint, limit int) ([]entities.PublicProduct, error)
	SaveContactSubmission(ctx context.Context, businessID uint, dto *dtos.ContactFormDTO) error

	// Integration gate — checks if an integration is active (or missing = backward compat)
	IsIntegrationActiveOrMissing(ctx context.Context, businessID uint, integrationTypeID uint) (bool, error)
}
