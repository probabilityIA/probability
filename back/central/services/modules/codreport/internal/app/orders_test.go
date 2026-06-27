package app

import (
	"context"
	"testing"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/codreport/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/codreport/internal/domain/entities"
	"github.com/secamc93/probability/back/central/shared/log"
)

type repoMock struct {
	listCodOrdersFn func(ctx context.Context, f dtos.OrdersFilter) ([]entities.CodOrder, int64, error)
}

func (m *repoMock) ListCodOrders(ctx context.Context, f dtos.OrdersFilter) ([]entities.CodOrder, int64, error) {
	return m.listCodOrdersFn(ctx, f)
}

func (m *repoMock) AggregateByCarrier(_ context.Context, _ dtos.ReportFilter, _ bool) ([]entities.CarrierAggregate, error) {
	return nil, nil
}
func (m *repoMock) MonthlyHistory(_ context.Context, _ uint, _ int) ([]entities.MonthlyPoint, error) {
	return nil, nil
}
func (m *repoMock) WeeklyAggregates(_ context.Context, _ uint, _ int) ([]entities.WeekAggregate, error) {
	return nil, nil
}
func (m *repoMock) CarrierConfigs(_ context.Context, _ uint) ([]entities.CarrierConfig, error) {
	return nil, nil
}
func (m *repoMock) DiscoveredCarriers(_ context.Context, _ uint) ([]string, error) {
	return nil, nil
}
func (m *repoMock) SaveCarrierConfig(_ context.Context, _ dtos.SaveCarrierConfigDTO) (*entities.CarrierConfig, error) {
	return &entities.CarrierConfig{}, nil
}
func (m *repoMock) ConfirmedCuts(_ context.Context, _ uint) ([]entities.PaymentCut, error) {
	return nil, nil
}
func (m *repoMock) SaveConfirmedCut(_ context.Context, _ entities.PaymentCut, _ uint, _ string) (*entities.PaymentCut, error) {
	return &entities.PaymentCut{}, nil
}
func (m *repoMock) UserName(_ context.Context, _ uint) string { return "" }
func (m *repoMock) CutPeriodOrders(_ context.Context, _ uint, _, _ time.Time) ([]entities.CarrierAggregate, error) {
	return nil, nil
}

func TestListOrders_ForwardsHasGuideFilter(t *testing.T) {
	guide := true
	var captured dtos.OrdersFilter
	repo := &repoMock{
		listCodOrdersFn: func(_ context.Context, f dtos.OrdersFilter) ([]entities.CodOrder, int64, error) {
			captured = f
			return []entities.CodOrder{}, 0, nil
		},
	}
	uc := New(repo, log.New())

	_, _, err := uc.ListOrders(context.Background(), dtos.OrdersFilter{BusinessID: 10, HasGuide: &guide})
	if err != nil {
		t.Fatalf("se esperaba nil error, se obtuvo: %v", err)
	}
	if captured.HasGuide == nil || *captured.HasGuide != true {
		t.Fatalf("HasGuide: el filtro no se reenvio al repositorio: %v", captured.HasGuide)
	}
}

func TestListOrders_PreservesHasGuide(t *testing.T) {
	repo := &repoMock{
		listCodOrdersFn: func(_ context.Context, _ dtos.OrdersFilter) ([]entities.CodOrder, int64, error) {
			return []entities.CodOrder{
				{OrderID: "a", Carrier: "SERVIENTREGA", HasGuide: true},
				{OrderID: "b", Carrier: "COORDINADORA", HasGuide: false},
			}, 2, nil
		},
	}
	uc := New(repo, log.New())

	orders, total, err := uc.ListOrders(context.Background(), dtos.OrdersFilter{BusinessID: 10})
	if err != nil {
		t.Fatalf("se esperaba nil error, se obtuvo: %v", err)
	}
	if total != 2 {
		t.Fatalf("total: esperado 2, obtenido %d", total)
	}
	if !orders[0].HasGuide {
		t.Errorf("orden 0: esperado HasGuide true")
	}
	if orders[1].HasGuide {
		t.Errorf("orden 1: esperado HasGuide false")
	}
}
