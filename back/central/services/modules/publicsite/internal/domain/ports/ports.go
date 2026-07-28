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

	// Integration gate — checks if an integration is active (or missing = backward compat)
	IsIntegrationActiveOrMissing(ctx context.Context, businessID uint, integrationTypeID uint) (bool, error)

	// GetPlatformIntegrationID resuelve el integration_id "platform" interno del negocio,
	// usado para atribuir las ordenes creadas desde el checkout de la tienda publica.
	GetPlatformIntegrationID(ctx context.Context, businessID uint) (uint, error)

	CreateCheckout(ctx context.Context, checkout *entities.PublicCheckout) error
	GetCheckoutByReference(ctx context.Context, reference string) (*entities.PublicCheckout, error)
}

// IBoldGateway es el contrato minimo que el modulo pay expone para generar una firma
// de Bold para un pago que no es recarga de billetera (checkout de tienda publica).
// *pay.Bundle satisface esta interfaz directamente (mismo patron que IWalletDebiter
// en el modulo subscriptions, que recibe payBundle como su propia interfaz local).
type IBoldGateway interface {
	BoldGenerateSignatureForReference(ctx context.Context, businessID uint, amount float64, currency, referencePrefix string) (*pay.BoldSignature, error)
}
