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
	MonthlyHistory(ctx context.Context, businessID uint, months int) ([]entities.MonthlyPoint, error)
	WeeklyAggregates(ctx context.Context, businessID uint, weeks int) ([]entities.WeekAggregate, error)
	CarrierConfigs(ctx context.Context, businessID uint) ([]entities.CarrierConfig, error)
	DiscoveredCarriers(ctx context.Context, businessID uint) ([]string, error)
	SaveCarrierConfig(ctx context.Context, d dtos.SaveCarrierConfigDTO) (*entities.CarrierConfig, error)
	ConfirmedCuts(ctx context.Context, businessID uint) ([]entities.PaymentCut, error)
	SaveConfirmedCut(ctx context.Context, cut entities.PaymentCut, userID uint, userName string) (*entities.PaymentCut, error)
	UserName(ctx context.Context, userID uint) string
	CutPeriodOrders(ctx context.Context, businessID uint, start, end time.Time) ([]entities.CarrierAggregate, error)
}
