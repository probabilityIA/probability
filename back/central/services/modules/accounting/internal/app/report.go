package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/accounting/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/accounting/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/accounting/internal/domain/errors"
)

func (uc *UseCase) Report(ctx context.Context, params dtos.ReportParams) (*dtos.ReportResponse, error) {
	if params.From.IsZero() || params.To.IsZero() || params.To.Before(params.From) {
		return nil, domainerrors.ErrInvalidPeriod
	}
	byConcept, byTax, err := uc.repo.Report(ctx, params)
	if err != nil {
		return nil, err
	}
	totals := dtos.ReportTotals{}
	for _, row := range byConcept {
		totals.TaxTotal += row.TaxTotal
		if row.Kind == entities.KindIncome {
			totals.CashIn += row.Amount
			if row.IsRealIncome {
				totals.RealIncome += row.Amount
			}
		} else {
			totals.CashOut += row.Amount
			if row.IsRealIncome {
				totals.RealExpense += row.Amount
			}
		}
	}
	totals.Net = round2(totals.RealIncome - totals.RealExpense - totals.TaxTotal)
	totals.RealIncome = round2(totals.RealIncome)
	totals.CashIn = round2(totals.CashIn)
	totals.RealExpense = round2(totals.RealExpense)
	totals.CashOut = round2(totals.CashOut)
	totals.TaxTotal = round2(totals.TaxTotal)

	return &dtos.ReportResponse{
		From:      params.From.Format("2006-01-02"),
		To:        params.To.Format("2006-01-02"),
		Totals:    totals,
		ByConcept: byConcept,
		ByTax:     byTax,
	}, nil
}
