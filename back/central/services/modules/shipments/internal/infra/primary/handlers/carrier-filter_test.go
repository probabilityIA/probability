package handlers

import (
	"context"
	"testing"

	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/domain"
	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/mocks"
)

func ratesFixture() []map[string]interface{} {
	return []map[string]interface{}{
		{"carrier": "ENVIA", "total": 12000.0},
		{"carrier": "TCC", "total": 11800.0},
		{"carrier": "Servientrega", "total": 13200.0},
	}
}

func carriersOf(rates []map[string]interface{}) []string {
	out := make([]string, 0, len(rates))
	for _, r := range rates {
		out = append(out, toStr(r["carrier"]))
	}
	return out
}

func newCarrierHandlers(settings []domain.CarrierSetting) *Handlers {
	repo := &mocks.RepositoryMock{
		GetBusinessCarrierSettingsFn: func(ctx context.Context, businessID uint, warehouseID *uint) ([]domain.CarrierSetting, error) {
			return settings, nil
		},
	}
	return newActorHandlers(repo)
}

func TestFilterRatesByBusinessCarriers_QuitaLaApagada(t *testing.T) {
	h := newCarrierHandlers([]domain.CarrierSetting{
		{Code: "ENVIA", Enabled: true, AllowCOD: true, AllowPrepaid: true},
		{Code: "TCC", Enabled: false, AllowCOD: true, AllowPrepaid: true},
		{Code: "SERVIENTREGA", Enabled: true, AllowCOD: true, AllowPrepaid: true},
	})

	out := h.filterRatesByBusinessCarriers(context.Background(), 26, nil, false, ratesFixture())

	got := carriersOf(out)
	if len(got) != 2 || got[0] != "ENVIA" || got[1] != "Servientrega" {
		t.Fatalf("TCC apagada debe salir de la lista, quedo %v", got)
	}
}

func TestFilterRatesByBusinessCarriers_RespetaContraEntrega(t *testing.T) {
	h := newCarrierHandlers([]domain.CarrierSetting{
		{Code: "ENVIA", Enabled: true, AllowCOD: true, AllowPrepaid: true},
		{Code: "TCC", Enabled: true, AllowCOD: true, AllowPrepaid: true},
		{Code: "SERVIENTREGA", Enabled: true, AllowCOD: false, AllowPrepaid: true},
	})

	cod := carriersOf(h.filterRatesByBusinessCarriers(context.Background(), 26, nil, true, ratesFixture()))
	if len(cod) != 2 {
		t.Fatalf("en contra entrega servientrega no debe aparecer, quedo %v", cod)
	}

	prepaid := carriersOf(h.filterRatesByBusinessCarriers(context.Background(), 26, nil, false, ratesFixture()))
	if len(prepaid) != 3 {
		t.Fatalf("en prepago deben quedar las tres, quedo %v", prepaid)
	}
}

func TestFilterRatesByBusinessCarriers_SinConfiguracionNoFiltra(t *testing.T) {
	h := newCarrierHandlers(nil)

	out := h.filterRatesByBusinessCarriers(context.Background(), 26, nil, false, ratesFixture())
	if len(out) != 3 {
		t.Fatalf("sin configuracion se deben mantener las tarifas, quedaron %d", len(out))
	}

	sinNegocio := h.filterRatesByBusinessCarriers(context.Background(), 0, nil, false, ratesFixture())
	if len(sinNegocio) != 3 {
		t.Fatalf("sin business_id no se filtra, quedaron %d", len(sinNegocio))
	}
}
