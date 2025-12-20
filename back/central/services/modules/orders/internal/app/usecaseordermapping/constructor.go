package usecaseordermapping

import (
	fulfillmentstatusdomain "github.com/secamc93/probability/back/central/services/modules/fulfillmentstatus/domain"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/app/usecaseorderscore"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain"
	orderstatusdomain "github.com/secamc93/probability/back/central/services/modules/orderstatus/domain"
	paymentstatusdomain "github.com/secamc93/probability/back/central/services/modules/paymentstatus/domain"
	"github.com/secamc93/probability/back/central/shared/log"
)

type UseCaseOrderMapping struct {
	repo                        domain.IRepository
	logger                      log.ILogger
	eventPublisher              domain.IOrderEventPublisher
	scoreUseCase                domain.IOrderScoreUseCase
	orderStatusRepository       orderstatusdomain.IRepository
	paymentStatusRepository     paymentstatusdomain.IRepository
	fulfillmentStatusRepository fulfillmentstatusdomain.IRepository
}

func New(repo domain.IRepository, logger log.ILogger, eventPublisher domain.IOrderEventPublisher, orderStatusRepo orderstatusdomain.IRepository, paymentStatusRepo paymentstatusdomain.IRepository, fulfillmentStatusRepo fulfillmentstatusdomain.IRepository) domain.IOrderMappingUseCase {
	return &UseCaseOrderMapping{
		repo:                        repo,
		logger:                      logger,
		eventPublisher:              eventPublisher,
		scoreUseCase:                usecaseorderscore.New(repo),
		orderStatusRepository:       orderStatusRepo,
		paymentStatusRepository:     paymentStatusRepo,
		fulfillmentStatusRepository: fulfillmentStatusRepo,
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

// mapShopifyFinancialStatusToPaymentStatus mapea el financial_status de Shopify al código de PaymentStatus de Probability
func mapShopifyFinancialStatusToPaymentStatus(financialStatus string) string {
	switch financialStatus {
	case "pending":
		return "pending"
	case "authorized":
		return "authorized"
	case "paid":
		return "paid"
	case "partially_paid":
		return "partially_paid"
	case "refunded":
		return "refunded"
	case "partially_refunded":
		return "partially_refunded"
	case "voided":
		return "voided"
	case "unpaid":
		return "unpaid"
	default:
		return "pending" // Valor por defecto
	}
}

// mapShopifyFulfillmentStatusToFulfillmentStatus mapea el fulfillment_status de Shopify al código de FulfillmentStatus de Probability
func mapShopifyFulfillmentStatusToFulfillmentStatus(fulfillmentStatus *string) string {
	if fulfillmentStatus == nil || *fulfillmentStatus == "" {
		return "unfulfilled"
	}
	switch *fulfillmentStatus {
	case "unfulfilled":
		return "unfulfilled"
	case "partial":
		return "partial"
	case "fulfilled":
		return "fulfilled"
	case "shipped":
		return "shipped"
	default:
		return "unfulfilled" // Valor por defecto
	}
}
