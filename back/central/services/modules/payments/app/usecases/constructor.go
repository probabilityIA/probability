package usecases

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/payments/domain"
)

// ═══════════════════════════════════════════
// INTERFACE
// ═══════════════════════════════════════════

// IUseCase define todos los casos de uso del módulo payments
type IUseCase interface {
	// Payment Methods
	ListPaymentMethods(ctx context.Context, page, pageSize int, filters map[string]interface{}) (*domain.PaymentMethodsListResponse, error)
	GetPaymentMethodByID(ctx context.Context, id uint) (*domain.PaymentMethodResponse, error)
	GetPaymentMethodByCode(ctx context.Context, code string) (*domain.PaymentMethodResponse, error)
	CreatePaymentMethod(ctx context.Context, req *domain.CreatePaymentMethodRequest) (*domain.PaymentMethodResponse, error)
	UpdatePaymentMethod(ctx context.Context, id uint, req *domain.UpdatePaymentMethodRequest) (*domain.PaymentMethodResponse, error)
	DeletePaymentMethod(ctx context.Context, id uint) error
	TogglePaymentMethodActive(ctx context.Context, id uint) (*domain.PaymentMethodResponse, error)

	// Payment Mappings
	ListPaymentMappings(ctx context.Context, filters map[string]interface{}) (*domain.PaymentMappingsListResponse, error)
	GetPaymentMappingByID(ctx context.Context, id uint) (*domain.PaymentMappingResponse, error)
	GetPaymentMappingsByIntegrationType(ctx context.Context, integrationType string) ([]domain.PaymentMappingResponse, error)
	GetAllPaymentMappingsGroupedByIntegration(ctx context.Context) ([]domain.PaymentMappingsByIntegrationResponse, error)
	CreatePaymentMapping(ctx context.Context, req *domain.CreatePaymentMappingRequest) (*domain.PaymentMappingResponse, error)
	UpdatePaymentMapping(ctx context.Context, id uint, req *domain.UpdatePaymentMappingRequest) (*domain.PaymentMappingResponse, error)
	DeletePaymentMapping(ctx context.Context, id uint) error
	TogglePaymentMappingActive(ctx context.Context, id uint) (*domain.PaymentMappingResponse, error)
}

// ═══════════════════════════════════════════
// CONSTRUCTOR
// ═══════════════════════════════════════════

// UseCase contiene todos los casos de uso del módulo payments
type UseCase struct {
	repo domain.IRepository
}

// New crea una nueva instancia de todos los casos de uso
func New(repo domain.IRepository) IUseCase {
	return &UseCase{
		repo: repo,
	}
}
