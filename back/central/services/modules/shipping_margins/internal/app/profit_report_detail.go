package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/shipping_margins/internal/domain/dtos"
)

func (uc *UseCase) ProfitReportDetail(ctx context.Context, params dtos.ProfitReportDetailParams) (*dtos.ProfitReportDetailResponse, error) {
	return uc.repo.ProfitReportDetail(ctx, params)
}
