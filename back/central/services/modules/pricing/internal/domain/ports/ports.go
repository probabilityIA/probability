package ports

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/pricing/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/pricing/internal/domain/entities"
)

// IRepository define los métodos del repositorio del módulo pricing
type IRepository interface {
	// Client Pricing Rules CRUD
	CreateClientPricingRule(ctx context.Context, rule *entities.ClientPricingRule) (*entities.ClientPricingRule, error)
	GetClientPricingRule(ctx context.Context, businessID, ruleID uint) (*entities.ClientPricingRule, error)
	ListClientPricingRules(ctx context.Context, params dtos.ListClientPricingRulesParams) ([]entities.ClientPricingRule, int64, error)
	UpdateClientPricingRule(ctx context.Context, rule *entities.ClientPricingRule) (*entities.ClientPricingRule, error)
	DeleteClientPricingRule(ctx context.Context, businessID, ruleID uint) error

	// Quantity Discounts CRUD
	CreateQuantityDiscount(ctx context.Context, discount *entities.QuantityDiscount) (*entities.QuantityDiscount, error)
	GetQuantityDiscount(ctx context.Context, businessID, discountID uint) (*entities.QuantityDiscount, error)
	ListQuantityDiscounts(ctx context.Context, params dtos.ListQuantityDiscountsParams) ([]entities.QuantityDiscount, int64, error)
	UpdateQuantityDiscount(ctx context.Context, discount *entities.QuantityDiscount) (*entities.QuantityDiscount, error)
	DeleteQuantityDiscount(ctx context.Context, businessID, discountID uint) error

	// Price Calculation queries
	GetApplicableClientRule(ctx context.Context, businessID, clientID uint, productID string) (*entities.ClientPricingRule, error)
	GetApplicableQuantityDiscount(ctx context.Context, businessID uint, productID string, quantity int) (*entities.QuantityDiscount, error)
}
