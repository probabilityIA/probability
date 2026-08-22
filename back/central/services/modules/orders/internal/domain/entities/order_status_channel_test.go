package entities

import "testing"

func TestIsOperated(t *testing.T) {
	casos := []struct {
		estado   OrderStatus
		esperado bool
	}{
		{OrderStatusPending, false},
		{OrderStatusOnHold, false},
		{OrderStatusProcessing, false},
		{OrderStatusPicking, true},
		{OrderStatusPacking, true},
		{OrderStatusInTransit, true},
		{OrderStatusDelivered, true},
		{OrderStatusReturned, true},
		{OrderStatusCancelled, true},
		{OrderStatus("paid"), false},
		{OrderStatus(""), false},
	}

	for _, caso := range casos {
		t.Run(string(caso.estado), func(t *testing.T) {
			if got := caso.estado.IsOperated(); got != caso.esperado {
				t.Fatalf("IsOperated(%q) = %v, se esperaba %v", caso.estado, got, caso.esperado)
			}
		})
	}
}

func TestCanChannelOverride(t *testing.T) {
	casos := []struct {
		nombre   string
		actual   OrderStatus
		delCanal OrderStatus
		esperado bool
	}{
		{"orden nueva la manda el canal", OrderStatusPending, OrderStatusProcessing, true},
		{"orden en espera la manda el canal", OrderStatusOnHold, OrderStatus("paid"), true},
		{"estado del canal fuera del catalogo sobre orden nueva", OrderStatus("paid"), OrderStatus("processing"), true},
		{"en bodega no retrocede", OrderStatusPicking, OrderStatus("paid"), false},
		{"despachada no retrocede", OrderStatusInTransit, OrderStatusProcessing, false},
		{"entregada no retrocede", OrderStatusDelivered, OrderStatus("paid"), false},
		{"cancelacion pisa en bodega", OrderStatusPicking, OrderStatusCancelled, true},
		{"cancelacion pisa en transito", OrderStatusInTransit, OrderStatusCancelled, true},
		{"reembolso pisa entregada", OrderStatusDelivered, OrderStatusRefunded, true},
		{"entrega no pisa una cancelada", OrderStatusCancelled, OrderStatusDelivered, false},
	}

	for _, caso := range casos {
		t.Run(caso.nombre, func(t *testing.T) {
			if got := caso.delCanal.CanChannelOverride(caso.actual); got != caso.esperado {
				t.Fatalf("CanChannelOverride(actual=%q, canal=%q) = %v, se esperaba %v", caso.actual, caso.delCanal, got, caso.esperado)
			}
		})
	}
}
