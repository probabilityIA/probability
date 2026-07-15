package app

import (
	"context"
	"testing"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/codreport/internal/domain"
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
func (m *repoMock) UpsertCutOrders(_ context.Context, _ entities.PaymentCut, _ []entities.PayoutOrder, _ uint, _ string) (uint, error) {
	return 1, nil
}
func (m *repoMock) PaidAggregatesForCut(_ context.Context, _ uint) ([]entities.CarrierAggregate, error) {
	return nil, nil
}
func (m *repoMock) UpdateCutTotals(_ context.Context, _ entities.PaymentCut) error {
	return nil
}
func (m *repoMock) UserName(_ context.Context, _ uint) string { return "" }
func (m *repoMock) CutPeriodOrders(_ context.Context, _ uint, _, _ time.Time) ([]entities.CarrierAggregate, error) {
	return nil, nil
}
func (m *repoMock) SelectableCutOrders(_ context.Context, _ dtos.SelectableOrdersFilter) ([]entities.CodOrder, error) {
	return nil, nil
}
func (m *repoMock) PayoutOrders(_ context.Context, _ uint, _ []string) ([]entities.PayoutOrder, error) {
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

func TestListOrders_SoloPagadaCuentaComoRecaudada(t *testing.T) {
	repo := &repoMock{
		listCodOrdersFn: func(_ context.Context, _ dtos.OrdersFilter) ([]entities.CodOrder, int64, error) {
			return []entities.CodOrder{
				{OrderID: "a", Status: "in_transit", Collected: true},
				{OrderID: "b", Status: "picked_up"},
				{OrderID: "c", Status: "pending"},
				{OrderID: "d", Status: "delivered"},
				{OrderID: "e", Status: "cancelled", Collected: true},
				{OrderID: "f", Status: "delivered", Paid: true},
			}, 6, nil
		},
	}
	uc := New(repo, log.New())

	orders, _, err := uc.ListOrders(context.Background(), dtos.OrdersFilter{BusinessID: 10})
	if err != nil {
		t.Fatalf("se esperaba nil error, se obtuvo: %v", err)
	}

	want := []struct {
		state     string
		collected bool
	}{
		{domain.CodStateInProgress, false},
		{domain.CodStateInProgress, false},
		{domain.CodStatePending, false},
		{domain.CodStatePendingPayment, false},
		{domain.CodStateNotCollectable, false},
		{domain.CodStateCollected, true},
	}
	for i := range want {
		if orders[i].CodState != want[i].state {
			t.Errorf("orden %d: CodState esperado %s, obtenido %s", i, want[i].state, orders[i].CodState)
		}
		if orders[i].Collected != want[i].collected {
			t.Errorf("orden %d (%s): Collected esperado %v, obtenido %v", i, orders[i].Status, want[i].collected, orders[i].Collected)
		}
	}
}
