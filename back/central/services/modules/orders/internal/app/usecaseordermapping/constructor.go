package usecaseordermapping

import (
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/app/usecaseorderscore"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain"
	orderstatusdomain "github.com/secamc93/probability/back/central/services/modules/orderstatus/domain"
	"github.com/secamc93/probability/back/central/shared/log"
)

type UseCaseOrderMapping struct {
	repo                  domain.IRepository
	logger                log.ILogger
	eventPublisher        domain.IOrderEventPublisher
	scoreUseCase          domain.IOrderScoreUseCase
	orderStatusRepository orderstatusdomain.IRepository
}

func New(repo domain.IRepository, logger log.ILogger, eventPublisher domain.IOrderEventPublisher, orderStatusRepo orderstatusdomain.IRepository) domain.IOrderMappingUseCase {
	return &UseCaseOrderMapping{
		repo:                  repo,
		logger:                logger,
		eventPublisher:        eventPublisher,
		scoreUseCase:          usecaseorderscore.New(repo),
		orderStatusRepository: orderStatusRepo,
	}
}

// getIntegrationTypeID convierte el código de tipo de integración a ID numérico
func getIntegrationTypeID(integrationType string) uint {
	switch integrationType {
	case "shopify":
		return 1
	case "whatsapp", "whatsap", "whastap":
		return 2
	case "mercado_libre", "mercadolibre":
		return 3
	case "woocommerce", "woocormerce":
		return 4
	default:
		return 0
	}
}
