package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/shipping_margins/internal/domain/dtos"
)

func (uc *UseCase) ProfitReport(ctx context.Context, params dtos.ProfitReportParams) (*dtos.ProfitReportResponse, error) {
	return uc.repo.ProfitReport(ctx, params)
}
