package handlers

import (
	"testing"

	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestShipmentHasActiveGuide(t *testing.T) {
	cases := []struct {
		name     string
		shipment *domain.Shipment
		want     bool
	}{
		{
			name:     "nil no tiene guia",
			shipment: nil,
			want:     false,
		},
		{
			name:     "pending sin tracking ni guide es reusable",
			shipment: &domain.Shipment{Status: "pending"},
			want:     false,
		},
		{
			name:     "pending con tracking ya tiene guia activa",
			shipment: &domain.Shipment{Status: "pending", TrackingNumber: strPtr("034057376067")},
			want:     true,
		},
		{
			name:     "pending con guide_url sin tracking ya tiene guia activa",
			shipment: &domain.Shipment{Status: "pending", GuideURL: strPtr("https://s3/guia.pdf")},
			want:     true,
		},
		{
			name:     "in_transit con tracking tiene guia activa",
			shipment: &domain.Shipment{Status: "in_transit", TrackingNumber: strPtr("034057376085")},
			want:     true,
		},
		{
			name:     "cancelled con tracking permite regenerar",
			shipment: &domain.Shipment{Status: "cancelled", TrackingNumber: strPtr("034057376067")},
			want:     false,
		},
		{
			name:     "failed con guide_url permite regenerar",
			shipment: &domain.Shipment{Status: "failed", GuideURL: strPtr("https://s3/guia.pdf")},
			want:     false,
		},
		{
			name:     "tracking vacio no cuenta como guia",
			shipment: &domain.Shipment{Status: "pending", TrackingNumber: strPtr("")},
			want:     false,
		},
		{
			name:     "guide_url vacio no cuenta como guia",
			shipment: &domain.Shipment{Status: "pending", GuideURL: strPtr("")},
			want:     false,
		},
		{
			name:     "tracking y guide_url vacios no cuentan como guia",
			shipment: &domain.Shipment{Status: "pending", TrackingNumber: strPtr(""), GuideURL: strPtr("")},
			want:     false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, shipmentHasActiveGuide(tc.shipment))
		})
	}
}

func TestShipmentHasActiveGuide_EscenarioIncidenteDobleCobro(t *testing.T) {
	autogenerado := &domain.Shipment{
		ID:             34678,
		Status:         "pending",
		TrackingNumber: strPtr("034057376067"),
		GuideURL:       strPtr("https://s3/034057376067.pdf"),
	}

	assert.True(t, shipmentHasActiveGuide(autogenerado),
		"un shipment autogenerado con tracking y PDF debe bloquear la generacion manual")
}

func TestShipmentHasActiveGuide_TrasCancelarPermiteNuevaGuia(t *testing.T) {
	cancelado := &domain.Shipment{
		ID:             34678,
		Status:         "cancelled",
		TrackingNumber: strPtr("034057376067"),
		GuideURL:       strPtr("https://s3/034057376067.pdf"),
	}

	assert.False(t, shipmentHasActiveGuide(cancelado),
		"tras cancelar la guia el negocio debe poder generar una nueva")
}
