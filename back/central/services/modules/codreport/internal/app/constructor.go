package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/codreport/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/codreport/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/codreport/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

type Iapp interface {
	Summary(ctx context.Context, f dtos.ReportFilter) (*entities.CodSummary, error)
	ListOrders(ctx context.Context, f dtos.OrdersFilter) ([]entities.CodOrder, int64, error)
	ListCuts(ctx context.Context, businessID uint, isAdmin bool) ([]entities.PaymentCut, error)
	SelectableOrders(ctx context.Context, f dtos.SelectableOrdersFilter) ([]entities.CodOrder, error)
	CutOrders(ctx context.Context, businessID uint, cutID uint) ([]entities.CodOrder, error)
	DeleteCut(ctx context.Context, businessID uint, cutID uint) error
	CreateDraft(ctx context.Context, d dtos.ConfirmCutDTO) (*entities.PaymentCut, error)
	ConfirmCut(ctx context.Context, businessID uint, cutID uint, userID uint, userName string) error
	CarrierConfigs(ctx context.Context, businessID uint) ([]entities.CarrierConfig, error)
	SaveCarrierConfig(ctx context.Context, d dtos.SaveCarrierConfigDTO) (*entities.CarrierConfig, error)
}

type UseCase struct {
	repo ports.IRepository
	log  log.ILogger
}

func New(repo ports.IRepository, logger log.ILogger) Iapp {
	return &UseCase{repo: repo, log: logger}
}

func (uc *UseCase) discountMap(_ context.Context, _ uint) map[string]float64 {
	return map[string]float64{}
}
