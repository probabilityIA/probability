package handlers

import (
	"testing"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestGuideInFlight(t *testing.T) {
	now := time.Now()

	cases := []struct {
		name     string
		shipment *domain.Shipment
		want     bool
	}{
		{"nil", nil, false},
		{"pending no esta en vuelo", &domain.Shipment{Status: domain.ShipmentStatusPending}, false},
		{"generating recien marcado esta en vuelo", &domain.Shipment{Status: domain.ShipmentStatusGenerating, UpdatedAt: now.Add(-2 * time.Second)}, true},
		{"generating a los 44 s sigue en vuelo", &domain.Shipment{Status: domain.ShipmentStatusGenerating, UpdatedAt: now.Add(-44 * time.Second)}, true},
		{"generating a los 2 min sigue en vuelo", &domain.Shipment{Status: domain.ShipmentStatusGenerating, UpdatedAt: now.Add(-2 * time.Minute)}, true},
		{"generating vencido deja reintentar", &domain.Shipment{Status: domain.ShipmentStatusGenerating, UpdatedAt: now.Add(-16 * time.Minute)}, false},
		{"failed no esta en vuelo", &domain.Shipment{Status: domain.ShipmentStatusFailed, UpdatedAt: now}, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, tc.shipment.GuideInFlight(now))
		})
	}
}

func TestGuideInFlight_EscenarioMYS0718(t *testing.T) {
	primerClic := time.Date(2026, 8, 18, 19, 30, 49, 0, time.UTC)
	segundoClic := time.Date(2026, 8, 18, 19, 32, 33, 0, time.UTC)

	shipment := &domain.Shipment{
		ID:        43991,
		Status:    domain.ShipmentStatusGenerating,
		UpdatedAt: primerClic,
	}

	assert.True(t, shipment.GuideInFlight(segundoClic),
		"el segundo clic 104 s despues debe encontrar la guia en vuelo y ser rechazado")
	assert.False(t, shipment.ReusableForGuide(segundoClic),
		"un shipment en vuelo no se puede reutilizar para otra generacion")
}

func TestReusableForGuide(t *testing.T) {
	now := time.Now()

	cases := []struct {
		name     string
		shipment *domain.Shipment
		want     bool
	}{
		{"nil", nil, false},
		{"pending limpio es reusable", &domain.Shipment{Status: domain.ShipmentStatusPending}, true},
		{"pending con tracking no es reusable", &domain.Shipment{Status: domain.ShipmentStatusPending, TrackingNumber: strPtr("2284436982")}, false},
		{"pending con guide_url no es reusable", &domain.Shipment{Status: domain.ShipmentStatusPending, GuideURL: strPtr("https://s3/g.pdf")}, false},
		{"generating en vuelo no es reusable", &domain.Shipment{Status: domain.ShipmentStatusGenerating, UpdatedAt: now.Add(-time.Minute)}, false},
		{"generating vencido es reusable", &domain.Shipment{Status: domain.ShipmentStatusGenerating, UpdatedAt: now.Add(-20 * time.Minute)}, true},
		{"failed no se reusa se crea uno nuevo", &domain.Shipment{Status: domain.ShipmentStatusFailed}, false},
		{"cancelled no se reusa", &domain.Shipment{Status: domain.ShipmentStatusCancelled}, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, tc.shipment.ReusableForGuide(now))
		})
	}
}
