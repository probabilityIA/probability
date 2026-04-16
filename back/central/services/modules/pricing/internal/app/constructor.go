package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/pricing/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/pricing/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/pricing/internal/domain/ports"
)

// IUseCase define los casos de uso del módulo pricing
type IUseCase interface {
	// Client Pricing Rules
	CreateClientPricingRule(ctx context.Context, dto dtos.CreateClientPricingRuleDTO) (*entities.ClientPricingRule, error)
	GetClientPricingRule(ctx context.Context, businessID, ruleID uint) (*entities.ClientPricingRule, error)
	ListClientPricingRules(ctx context.Context, params dtos.ListClientPricingRulesParams) ([]entities.ClientPricingRule, int64, error)
	UpdateClientPricingRule(ctx context.Context, dto dtos.UpdateClientPricingRuleDTO) (*entities.ClientPricingRule, error)
	DeleteClientPricingRule(ctx context.Context, businessID, ruleID uint) error

	// Quantity Discounts
	CreateQuantityDiscount(ctx context.Context, dto dtos.CreateQuantityDiscountDTO) (*entities.QuantityDiscount, error)
	GetQuantityDiscount(ctx context.Context, businessID, discountID uint) (*entities.QuantityDiscount, error)
	ListQuantityDiscounts(ctx context.Context, params dtos.ListQuantityDiscountsParams) ([]entities.QuantityDiscount, int64, error)
	UpdateQuantityDiscount(ctx context.Context, dto dtos.UpdateQuantityDiscountDTO) (*entities.QuantityDiscount, error)
	DeleteQuantityDiscount(ctx context.Context, businessID, discountID uint) error

	// Price Calculation
	CalculatePrice(ctx context.Context, req dtos.CalculatePriceRequest) (*entities.PriceResult, error)
}

// UseCase implementa IUseCase
type UseCase struct {
	repo ports.IRepository
}

// New crea una nueva instancia del use case
func New(repo ports.IRepository) IUseCase {
	return &UseCase{repo: repo}
}
