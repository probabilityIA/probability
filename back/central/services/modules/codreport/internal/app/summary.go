package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/codreport/internal/domain"
	"github.com/secamc93/probability/back/central/services/modules/codreport/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/codreport/internal/domain/entities"
)

func (uc *UseCase) Summary(ctx context.Context, f dtos.ReportFilter) (*entities.CodSummary, error) {
	dm := uc.discountMap(ctx, f.BusinessID)

	collected, err := uc.repo.AggregateByCarrier(ctx, f, true)
	if err != nil {
		return nil, err
	}
	domain.EnrichCarrierAggregates(collected, dm)

	pending, err := uc.repo.AggregateByCarrier(ctx, f, false)
	if err != nil {
		return nil, err
	}

	ordersCollected, totalCollected, totalDiscount, totalNet := domain.SumAggregates(collected)
	ordersPending, totalPending, _, _ := domain.SumAggregates(pending)

	effRate := 0.0
	if totalCollected > 0 {
		effRate = totalDiscount / totalCollected * 100
	}

	monthly, err := uc.repo.MonthlyHistory(ctx, f.BusinessID, 6)
	if err != nil {
		return nil, err
	}
	for i := range monthly {
		d, n := domain.ApplyDiscount(monthly[i].Collected, effRate)
		monthly[i].Discount = d
		monthly[i].Net = n
	}

	detail, err := uc.repo.SummaryCarrierDetail(ctx, f)
	if err != nil {
		return nil, err
	}
	history, err := uc.repo.SummaryHistory(ctx, f)
	if err != nil {
		return nil, err
	}

	var enCursoTotal, entregadoTotal float64
	var enCursoOrders, entregadoOrders int
	for i := range detail {
		enCursoTotal += detail[i].EnCurso
		enCursoOrders += detail[i].EnCursoOrders
		entregadoTotal += detail[i].Entregado
		entregadoOrders += detail[i].EntregadoOrders
	}

	return &entities.CodSummary{
		TotalCollected:  totalCollected,
		TotalPending:    totalPending,
		TotalDiscount:   totalDiscount,
		TotalNet:        totalNet,
		OrdersCollected: ordersCollected,
		OrdersPending:   ordersPending,
		ByCarrier:       collected,
		Monthly:         monthly,
		EnCursoTotal:    enCursoTotal,
		EnCursoOrders:   enCursoOrders,
		EntregadoTotal:  entregadoTotal,
		EntregadoOrders: entregadoOrders,
		CarrierDetail:   detail,
		History:         history,
	}, nil
}
