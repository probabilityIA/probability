package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/pricing/internal/domain"
	"github.com/secamc93/probability/back/central/services/modules/pricing/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/pricing/internal/domain/entities"
)

func (uc *UseCase) CalculatePrice(ctx context.Context, req dtos.CalculatePriceRequest) (*entities.PriceResult, error) {
	clientRule, _ := uc.repo.GetApplicableClientRule(ctx, req.BusinessID, req.ClientID, req.ProductID)
	qtyDiscount, _ := uc.repo.GetApplicableQuantityDiscount(ctx, req.BusinessID, req.ProductID, req.Quantity)

	result := domain.CalculatePrice(req.BasePrice, clientRule, qtyDiscount)
	return &result, nil
}
