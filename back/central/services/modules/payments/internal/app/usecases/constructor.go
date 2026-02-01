package usecases

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/payments/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/payments/internal/domain/ports"
)

// ═══════════════════════════════════════════
// INTERFACE
// ═══════════════════════════════════════════

// IUseCase define todos los casos de uso del módulo payments
type IUseCase interface {
	// Payment Methods
	ListPaymentMethods(ctx context.Context, page, pageSize int, filters map[string]interface{}) (*dtos.PaymentMethodsListResponse, error)
	GetPaymentMethodByID(ctx context.Context, id uint) (*dtos.PaymentMethodResponse, error)
	GetPaymentMethodByCode(ctx context.Context, code string) (*dtos.PaymentMethodResponse, error)
	CreatePaymentMethod(ctx context.Context, req *dtos.CreatePaymentMethod) (*dtos.PaymentMethodResponse, error)
	UpdatePaymentMethod(ctx context.Context, id uint, req *dtos.UpdatePaymentMethod) (*dtos.PaymentMethodResponse, error)
	DeletePaymentMethod(ctx context.Context, id uint) error
	TogglePaymentMethodActive(ctx context.Context, id uint) (*dtos.PaymentMethodResponse, error)

	// Payment Mappings
	ListPaymentMappings(ctx context.Context, filters map[string]interface{}) (*dtos.PaymentMappingsListResponse, error)
	GetPaymentMappingByID(ctx context.Context, id uint) (*dtos.PaymentMappingResponse, error)
	GetPaymentMappingsByIntegrationType(ctx context.Context, integrationType string) ([]dtos.PaymentMappingResponse, error)
	GetAllPaymentMappingsGroupedByIntegration(ctx context.Context) ([]dtos.PaymentMappingsByIntegrationResponse, error)
	CreatePaymentMapping(ctx context.Context, req *dtos.CreatePaymentMapping) (*dtos.PaymentMappingResponse, error)
	UpdatePaymentMapping(ctx context.Context, id uint, req *dtos.UpdatePaymentMapping) (*dtos.PaymentMappingResponse, error)
	DeletePaymentMapping(ctx context.Context, id uint) error
	TogglePaymentMappingActive(ctx context.Context, id uint) (*dtos.PaymentMappingResponse, error)
}

// ═══════════════════════════════════════════
// CONSTRUCTOR
// ═══════════════════════════════════════════

// UseCase contiene todos los casos de uso del módulo payments
type UseCase struct {
	repo ports.IRepository
}

// New crea una nueva instancia de todos los casos de uso
func New(repo ports.IRepository) IUseCase {
	return &UseCase{
		repo: repo,
	}
}
