package domain

import (
	"context"

	"github.com/secamc93/probability/back/migration/shared/models"
)

// ───────────────────────────────────────────
//
//	REPOSITORY INTERFACE
//
// ───────────────────────────────────────────

// IRepository define todos los métodos de repositorio del módulo payments
type IRepository interface {
	// Payment Methods
	CreatePaymentMethod(ctx context.Context, method *models.PaymentMethod) error
	GetPaymentMethodByID(ctx context.Context, id uint) (*models.PaymentMethod, error)
	GetPaymentMethodByCode(ctx context.Context, code string) (*models.PaymentMethod, error)
	ListPaymentMethods(ctx context.Context, page, pageSize int, filters map[string]interface{}) ([]models.PaymentMethod, int64, error)
	UpdatePaymentMethod(ctx context.Context, method *models.PaymentMethod) error
	DeletePaymentMethod(ctx context.Context, id uint) error
	TogglePaymentMethodActive(ctx context.Context, id uint) (*models.PaymentMethod, error)
	PaymentMethodExists(ctx context.Context, code string) (bool, error)
	PaymentMethodHasActiveMappings(ctx context.Context, id uint) (bool, error)

	// Payment Mappings
	CreatePaymentMapping(ctx context.Context, mapping *models.PaymentMethodMapping) error
	GetPaymentMappingByID(ctx context.Context, id uint) (*models.PaymentMethodMapping, error)
	GetPaymentMappingByIDWithMethod(ctx context.Context, id uint) (*models.PaymentMethodMapping, error)
	ListPaymentMappings(ctx context.Context, filters map[string]interface{}) ([]models.PaymentMethodMapping, int64, error)
	ListPaymentMappingsWithMethods(ctx context.Context, filters map[string]interface{}) ([]models.PaymentMethodMapping, int64, error)
	UpdatePaymentMapping(ctx context.Context, mapping *models.PaymentMethodMapping) error
	DeletePaymentMapping(ctx context.Context, id uint) error
	GetPaymentMappingsByIntegrationType(ctx context.Context, integrationType string) ([]models.PaymentMethodMapping, error)
	GetPaymentMappingsByIntegrationTypeWithMethods(ctx context.Context, integrationType string) ([]models.PaymentMethodMapping, error)
	TogglePaymentMappingActive(ctx context.Context, id uint) (*models.PaymentMethodMapping, error)
	PaymentMappingExists(ctx context.Context, integrationType, originalMethod string) (bool, error)
}
