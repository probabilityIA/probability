package domain_test

import (
	"testing"

	"github.com/secamc93/probability/back/central/services/modules/codreport/internal/domain"
)

func TestCollectionState(t *testing.T) {
	cases := map[string]string{
		"delivered":        domain.CodStateCollected,
		"picked_up":        domain.CodStateInProgress,
		"in_transit":       domain.CodStateInProgress,
		"out_for_delivery": domain.CodStateInProgress,
		"on_hold":          domain.CodStateInProgress,
		"pending":          domain.CodStatePending,
		"cancelled":        domain.CodStateNotCollectable,
		"failed":           domain.CodStateNotCollectable,
		"returned":         domain.CodStateNotCollectable,
		"":                 domain.CodStateNotCollectable,
	}

	for status, want := range cases {
		if got := domain.CollectionState(status); got != want {
			t.Errorf("CollectionState(%q): esperado %s, obtenido %s", status, want, got)
		}
	}
}

func TestIsCollectedSoloEntregado(t *testing.T) {
	for _, status := range []string{"pending", "picked_up", "in_transit", "out_for_delivery", "on_hold", "failed", "returned", "cancelled"} {
		if domain.IsCollected(status) {
			t.Errorf("IsCollected(%q): una orden no entregada no puede contar como recaudada", status)
		}
	}
	if !domain.IsCollected("delivered") {
		t.Error("IsCollected(\"delivered\"): esperado true")
	}
}
