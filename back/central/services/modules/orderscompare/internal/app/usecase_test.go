package app

import (
	"context"
	"testing"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/orderscompare/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/orderscompare/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/orderscompare"
)

type repoMock struct {
	integration *ports.Integration
	locals      []orderscompare.LocalOrder
}

func (r *repoMock) GetIntegration(_ context.Context, _ uint) (*ports.Integration, error) {
	return r.integration, nil
}

func (r *repoMock) ListLocalOrders(_ context.Context, _ uint, _ uint, _, _ *time.Time, _ int) ([]dtos.LocalOrder, error) {
	return r.locals, nil
}

type registryMock struct {
	supports  bool
	orders    []orderscompare.ChannelOrder
	imported  []string
	resultado orderscompare.ImportResult
}

func (r *registryMock) Supports(_ uint) bool { return r.supports }

func (r *registryMock) ListOrders(_ context.Context, _ uint, _ string, _ orderscompare.ChannelFilters) ([]orderscompare.ChannelOrder, error) {
	return r.orders, nil
}

func (r *registryMock) ImportOrders(_ context.Context, _ uint, _ string, externalIDs []string) (orderscompare.ImportResult, error) {
	r.imported = externalIDs
	if r.resultado.Queued == nil {
		return orderscompare.ImportResult{Queued: externalIDs, Failed: map[string]string{}}, nil
	}
	return r.resultado, nil
}

func nuevoCaso(locals []orderscompare.LocalOrder, canal []orderscompare.ChannelOrder) (ports.IUseCase, *registryMock) {
	businessID := uint(26)
	repo := &repoMock{
		integration: &ports.Integration{ID: 259, Name: "Demo Tiendanube", BusinessID: &businessID, IntegrationType: 17, IsActive: true},
		locals:      locals,
	}
	registry := &registryMock{supports: true, orders: canal}
	return New(repo, registry, log.New()), registry
}

func TestCompareDevuelveFilasYTotales(t *testing.T) {
	uc, _ := nuevoCaso(
		[]orderscompare.LocalOrder{{OrderID: "u1", ExternalID: "1001", Status: "paid", Total: 100}},
		[]orderscompare.ChannelOrder{
			{ExternalID: "1001", Status: "paid", Total: 100},
			{ExternalID: "1002", Status: "completed", Total: 200},
		},
	)

	page, err := uc.Compare(context.Background(), dtos.CompareQuery{BusinessID: 26, IntegrationID: 259})
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if page.Totals.ToCreate != 1 || page.Totals.InSync != 1 {
		t.Fatalf("totales inesperados: %+v", page.Totals)
	}
	if page.Totals.WithoutInventory != 1 {
		t.Fatalf("esperaba 1 orden sin inventario, obtuve %d", page.Totals.WithoutInventory)
	}
}

func TestCompareRechazaIntegracionDeOtroNegocio(t *testing.T) {
	uc, _ := nuevoCaso(nil, nil)

	if _, err := uc.Compare(context.Background(), dtos.CompareQuery{BusinessID: 99, IntegrationID: 259}); err == nil {
		t.Fatal("esperaba error de aislamiento multi-tenant")
	}
}

func TestApplyOmiteLasQueYaExisten(t *testing.T) {
	uc, registry := nuevoCaso(
		[]orderscompare.LocalOrder{{OrderID: "u1", ExternalID: "1001"}},
		[]orderscompare.ChannelOrder{
			{ExternalID: "1001", Status: "paid"},
			{ExternalID: "1002", Status: "delivered"},
		},
	)

	result, err := uc.Apply(context.Background(), dtos.ApplyCommand{
		BusinessID:    26,
		IntegrationID: 259,
		ExternalIDs:   []string{"1001", "1002", "1002"},
	})
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if len(result.Skipped) != 1 || result.Skipped[0] != "1001" {
		t.Fatalf("esperaba omitir 1001, obtuve %+v", result.Skipped)
	}
	if len(registry.imported) != 1 || registry.imported[0] != "1002" {
		t.Fatalf("esperaba importar solo 1002, obtuve %+v", registry.imported)
	}
	if len(result.WithoutInventory) != 1 {
		t.Fatalf("esperaba avisar que 1002 no mueve inventario, obtuve %+v", result.WithoutInventory)
	}
	if result.Note == "" {
		t.Fatal("esperaba una nota explicando las ordenes historicas")
	}
}
