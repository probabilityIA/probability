package domain

import "testing"

func TestMapStatusStepToProbability(t *testing.T) {
	tests := []struct {
		name        string
		statusStep  string
		incidence   bool
		wantStatus  ProbabilityShipmentStatus
		wantUnknown bool
	}{
		{"pendiente", "Pendiente", false, StatusPending, false},
		{"pendiente de recoleccion con acento", "Pendiente de Recolecci\u00f3n", false, StatusPending, false},
		{"pendiente de recoleccion sin acento", "Pendiente de Recoleccion", false, StatusPending, false},
		{"envio recolectado con acento", "Env\u00edo Recolectado", false, StatusPickedUp, false},
		{"envio recolectado sin acento", "Envio Recolectado", false, StatusPickedUp, false},
		{"recolectado solo", "Recolectado", false, StatusPickedUp, false},
		{"en transito con acento", "En tr\u00e1nsito", false, StatusInTransit, false},
		{"en transito mayuscula con acento", "En Tr\u00e1nsito", false, StatusInTransit, false},
		{"en transito sin acento", "En Transito", false, StatusInTransit, false},
		{"en transito minuscula", "en transito", false, StatusInTransit, false},
		{"en distribucion con acento", "En Distribuci\u00f3n", false, StatusOutForDelivery, false},
		{"en distribucion sin acento", "En distribucion", false, StatusOutForDelivery, false},
		{"entregado masculino", "Entregado", false, StatusDelivered, false},
		{"entregada femenino", "Entregada", false, StatusDelivered, false},
		{"novedad", "Novedad", false, StatusOnHold, false},
		{"incidencia por status step", "Incidencia", false, StatusOnHold, false},
		{"no entregado minuscula", "No entregado", false, StatusFailed, false},
		{"no entregado mayuscula", "No Entregado", false, StatusFailed, false},
		{"devuelto", "Devuelto", false, StatusReturned, false},
		{"devuelta", "Devuelta", false, StatusReturned, false},
		{"regresado", "Regresado", false, StatusReturned, false},
		{"cancelado", "Cancelado", false, StatusCancelled, false},
		{"incidencia flag con pendiente", "Pendiente", true, StatusOnHold, false},
		{"incidencia flag con entregado resuelve como delivered", "Entregado", true, StatusDelivered, false},
		{"incidencia flag con en transito", "En Transito", true, StatusOnHold, false},
		{"status desconocido retorna in_transit con flag unknown", "Estado Desconocido XYZ", false, StatusInTransit, true},
		{"status vacio retorna in_transit con flag unknown", "", false, StatusInTransit, true},
		{"espacios y mayusculas no afectan", "   EN TRANSITO   ", false, StatusInTransit, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStatus, gotUnknown := MapStatusStepToProbability(tt.statusStep, tt.incidence)
			if gotStatus != tt.wantStatus {
				t.Errorf("MapStatusStepToProbability(%q, %v) status = %q, want %q", tt.statusStep, tt.incidence, gotStatus, tt.wantStatus)
			}
			if gotUnknown != tt.wantUnknown {
				t.Errorf("MapStatusStepToProbability(%q, %v) unknown = %v, want %v", tt.statusStep, tt.incidence, gotUnknown, tt.wantUnknown)
			}
		})
	}
}

func TestWebhookPayloadToNormalizedUpdate(t *testing.T) {
	t.Run("entregado con fecha de entrega", func(t *testing.T) {
		p := &WebhookPayload{
			TrackingCode:     "034056642049",
			RealDeliveryDate: "2026-04-17 15:30:00",
			RealPickupDate:   "2026-04-15 10:00:00",
			Events: []WebhookEvent{
				{StatusStep: "Entregado", Incidence: false, Description: "Entregado al cliente"},
			},
		}
		u := p.ToNormalizedUpdate()
		if u == nil {
			t.Fatal("expected update, got nil")
		}
		if u.ProbabilityStatus != StatusDelivered {
			t.Errorf("status = %q, want %q", u.ProbabilityStatus, StatusDelivered)
		}
		if u.DeliveredAt == nil || *u.DeliveredAt != "2026-04-17 15:30:00" {
			t.Errorf("DeliveredAt = %v, want 2026-04-17 15:30:00", u.DeliveredAt)
		}
		if u.ShippedAt == nil || *u.ShippedAt != "2026-04-15 10:00:00" {
			t.Errorf("ShippedAt = %v, want 2026-04-15 10:00:00", u.ShippedAt)
		}
	})

	t.Run("pendiente no tiene shipped_at", func(t *testing.T) {
		p := &WebhookPayload{
			TrackingCode:   "ABC123",
			RealPickupDate: "2026-04-15 10:00:00",
			Events: []WebhookEvent{
				{StatusStep: "Pendiente", Incidence: false},
			},
		}
		u := p.ToNormalizedUpdate()
		if u.ProbabilityStatus != StatusPending {
			t.Errorf("status = %q, want %q", u.ProbabilityStatus, StatusPending)
		}
		if u.ShippedAt != nil {
			t.Errorf("ShippedAt should be nil for pending, got %v", *u.ShippedAt)
		}
	})

	t.Run("sin eventos retorna nil", func(t *testing.T) {
		p := &WebhookPayload{TrackingCode: "X", Events: []WebhookEvent{}}
		if u := p.ToNormalizedUpdate(); u != nil {
			t.Errorf("expected nil, got %+v", u)
		}
	})

	t.Run("status desconocido marca unknown", func(t *testing.T) {
		p := &WebhookPayload{
			TrackingCode: "X",
			Events:       []WebhookEvent{{StatusStep: "Estado Raro", Incidence: false}},
		}
		u := p.ToNormalizedUpdate()
		if !u.IsUnknownStatus {
			t.Error("expected IsUnknownStatus=true")
		}
		if u.ProbabilityStatus != StatusInTransit {
			t.Errorf("fallback should be in_transit, got %q", u.ProbabilityStatus)
		}
	})
}
