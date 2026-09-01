package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/entities"
)

func (uc *useCase) GetCODAccountReport(ctx context.Context, businessID uint, startDate, endDate string) (*entities.CODAccountReport, error) {
	uc.log.Info(ctx).
		Uint("business_id", businessID).
		Str("start_date", startDate).
		Str("end_date", endDate).
		Msg("Getting COD account report")

	report, err := uc.repo.GetCODAccountReport(ctx, businessID, startDate, endDate)
	if err != nil {
		uc.log.Error(ctx).Err(err).Msg("Failed to get COD account report")
		return nil, err
	}

	report.Period = entities.PeriodInfo{
		Label: formatReportPeriodLabel(startDate, endDate),
	}

	uc.log.Info(ctx).
		Int("total_invoices", report.TotalInvoices).
		Int("cod_count", report.CODCount).
		Msg("COD account report retrieved successfully")

	return report, nil
}

func formatReportPeriodLabel(startDate, endDate string) string {
	if startDate == "" && endDate == "" {
		return "Todo el histórico"
	}
	if startDate == "" {
		return "Hasta " + endDate
	}
	if endDate == "" {
		return "Desde " + startDate
	}
	return startDate + " - " + endDate
}
