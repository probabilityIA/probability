package entities

import "testing"

func TestInventarioNoPisaEstadosFinales(t *testing.T) {
	tests := []struct {
		name         string
		actual       OrderStatus
		objetivo     OrderStatus
		puedeCambiar bool
	}{
		{"pendiente entra a picking", OrderStatusPending, OrderStatusPicking, true},
		{"en espera entra a picking", OrderStatusOnHold, OrderStatusPicking, true},
		{"novedad de inventario reintenta picking", OrderStatusInventoryIssue, OrderStatusPicking, true},
		{"picking pasa a novedad de inventario", OrderStatusPicking, OrderStatusInventoryIssue, true},

		{"completada no vuelve a picking", OrderStatusCompleted, OrderStatusPicking, false},
		{"cancelada no vuelve a picking", OrderStatusCancelled, OrderStatusPicking, false},
		{"reembolsada no vuelve a picking", OrderStatusRefunded, OrderStatusPicking, false},
		{"fallida no vuelve a picking", OrderStatusFailed, OrderStatusPicking, false},
		{"entregada no vuelve a picking", OrderStatusDelivered, OrderStatusPicking, false},
		{"en camino no vuelve a picking", OrderStatusInTransit, OrderStatusPicking, false},
		{"devuelta no vuelve a picking", OrderStatusReturned, OrderStatusPicking, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.actual.CanTransitionTo(tt.objetivo)
			if got != tt.puedeCambiar {
				t.Errorf("%s -> %s = %v, se esperaba %v", tt.actual, tt.objetivo, got, tt.puedeCambiar)
			}
		})
	}
}
