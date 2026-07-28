package consumer

import (
	"testing"

	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/domain"
)

func strPtrTest(s string) *string { return &s }

func TestCarrierFromPlatform(t *testing.T) {
	cases := map[string]string{
		"mercadolibre": "MERCADO ENVIOS",
		"MercadoLibre": "MERCADO ENVIOS",
		"shopify":      "SHOPIFY",
		"":             "",
	}
	for in, want := range cases {
		if got := carrierFromPlatform(in); got != want {
			t.Fatalf("carrierFromPlatform(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestHasGuide(t *testing.T) {
	base := domain.OrderExternalGuide{TrackingNumber: "MEL123", Status: "shipped"}
	if !base.HasGuide() {
		t.Fatal("orden con tracking deberia tener guia")
	}

	cancelled := base
	cancelled.Status = "cancelled"
	if cancelled.HasGuide() {
		t.Fatal("orden cancelada no deberia materializar envio")
	}

	empty := domain.OrderExternalGuide{Status: "shipped"}
	if empty.HasGuide() {
		t.Fatal("orden sin datos de guia no deberia materializar envio")
	}
}

func TestShouldSkipExternalGuide(t *testing.T) {
	guide := &domain.OrderExternalGuide{
		Platform:       "mercadolibre",
		TrackingNumber: "MEL123",
		GuideID:        "47626453039",
		GuideURL:       "https://api.mercadolibre.com/shipment_labels?shipment_ids=47626453039",
	}

	tests := []struct {
		name     string
		existing []domain.Shipment
		want     bool
	}{
		{
			name:     "sin envios previos crea",
			existing: nil,
			want:     false,
		},
		{
			name:     "mismo tracking no duplica",
			existing: []domain.Shipment{{TrackingNumber: strPtrTest("MEL123")}},
			want:     true,
		},
		{
			name:     "envio con guia propia no se toca",
			existing: []domain.Shipment{{GuideURL: strPtrTest("https://envioclick/guia.pdf"), TrackingNumber: strPtrTest("034058029645")}},
			want:     true,
		},
		{
			name:     "envio pendiente de envioclick sin tracking no se toca",
			existing: []domain.Shipment{{Status: "pending"}},
			want:     true,
		},
		{
			name:     "mismo guide id no duplica",
			existing: []domain.Shipment{{GuideID: strPtrTest("47626453039")}},
			want:     true,
		},
		{
			name:     "envio cancelado de otra guia permite crear",
			existing: []domain.Shipment{{Status: "cancelled", TrackingNumber: strPtrTest("OTRO999")}},
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := shouldSkipExternalGuide(tt.existing, guide); got != tt.want {
				t.Fatalf("shouldSkipExternalGuide() = %v, want %v", got, tt.want)
			}
		})
	}
}
