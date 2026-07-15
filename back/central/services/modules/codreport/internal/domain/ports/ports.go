package ports

import (
	"context"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/codreport/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/codreport/internal/domain/entities"
)

type IRepository interface {
	ListCodOrders(ctx context.Context, f dtos.OrdersFilter) ([]entities.CodOrder, int64, error)
	AggregateByCarrier(ctx context.Context, f dtos.ReportFilter, collected bool) ([]entities.CarrierAggregate, error)
	SummaryCarrierDetail(ctx context.Context, f dtos.ReportFilter) ([]entities.CarrierDetail, error)
	SummaryHistory(ctx context.Context, f dtos.ReportFilter) ([]entities.HistoryPoint, error)
	MonthlyHistory(ctx context.Context, businessID uint, months int) ([]entities.MonthlyPoint, error)
	WeeklyAggregates(ctx context.Context, businessID uint, weeks int) ([]entities.WeekAggregate, error)
	CarrierConfigs(ctx context.Context, businessID uint) ([]entities.CarrierConfig, error)
	DiscoveredCarriers(ctx context.Context, businessID uint) ([]string, error)
	SaveCarrierConfig(ctx context.Context, d dtos.SaveCarrierConfigDTO) (*entities.CarrierConfig, error)
	ConfirmedCuts(ctx context.Context, businessID uint) ([]entities.PaymentCut, error)
	UpsertCutOrders(ctx context.Context, cut entities.PaymentCut, orders []entities.PayoutOrder, userID uint, userName string) (uint, error)
	PaidAggregatesForCut(ctx context.Context, cutID uint) ([]entities.CarrierAggregate, error)
	UpdateCutTotals(ctx context.Context, cut entities.PaymentCut) error
	UserName(ctx context.Context, userID uint) string
	CutPeriodOrders(ctx context.Context, businessID uint, start, end time.Time) ([]entities.CarrierAggregate, error)
	SelectableCutOrders(ctx context.Context, f dtos.SelectableOrdersFilter) ([]entities.CodOrder, error)
	PayoutOrders(ctx context.Context, businessID uint, orderIDs []string) ([]entities.PayoutOrder, error)
	CutOrders(ctx context.Context, businessID uint, cutID uint) ([]entities.CodOrder, error)
	ConfirmDraftCut(ctx context.Context, businessID uint, cutID uint, userID uint, userName string) error
	DeleteCut(ctx context.Context, businessID uint, cutID uint) error
}
