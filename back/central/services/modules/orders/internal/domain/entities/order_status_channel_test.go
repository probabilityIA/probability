package entities

import "testing"

func TestIsChannelOverride(t *testing.T) {
	casos := []struct {
		estado   OrderStatus
		esperado bool
	}{
		{OrderStatusCancelled, true},
		{OrderStatusRefunded, true},
		{OrderStatusDelivered, false},
		{OrderStatusPicking, false},
		{OrderStatusPending, false},
		{OrderStatus("paid"), false},
		{OrderStatus(""), false},
	}

	for _, caso := range casos {
		t.Run(string(caso.estado), func(t *testing.T) {
			if got := caso.estado.IsChannelOverride(); got != caso.esperado {
				t.Fatalf("IsChannelOverride(%q) = %v, se esperaba %v", caso.estado, got, caso.esperado)
			}
		})
	}
}
