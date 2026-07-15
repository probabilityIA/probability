package mapper

import (
	"testing"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/meli/internal/domain"
)

func baseOrder() *domain.MeliOrder {
	sku := "SKU-1"
	return &domain.MeliOrder{
		ID:          123,
		Status:      "paid",
		TotalAmount: 100,
		CurrencyID:  "COP",
		Buyer: domain.MeliBuyer{
			FirstName: "Ana",
			LastName:  "Perez",
			BillingInfo: &domain.MeliBillingInfo{
				DocType:   "CC",
				DocNumber: "999",
			},
		},
		OrderItems: []domain.MeliOrderItem{
			{Item: domain.MeliItem{ID: "MCO1", Title: "Item", SellerSKU: &sku}, Quantity: 2, UnitPrice: 50, Currency: "COP"},
		},
	}
}

func TestDeriveOrderType(t *testing.T) {
	cases := map[string]string{
		"fulfillment":   "full",
		"self_service":  "flex",
		"cross_docking": "cross_docking",
		"drop_off":      "mercado_envios",
	}
	for logistic, want := range cases {
		got := deriveOrderType(&domain.MeliShippingDetail{LogisticType: logistic})
		if got != want {
			t.Errorf("logistic %q: got %q want %q", logistic, got, want)
		}
	}
	if got := deriveOrderType(nil); got != "marketplace" {
		t.Errorf("nil shipping: got %q want marketplace", got)
	}
}

func TestMapCustomerAndGuide(t *testing.T) {
	order := baseOrder()
	shipping := &domain.MeliShippingDetail{
		ID:             777,
		Status:         "shipped",
		LogisticType:   "fulfillment",
		TrackingNumber: "TRACK9",
		ReceiverAddress: &domain.MeliReceiverAddress{
			StreetName:   "Calle 1",
			StreetNumber: "23",
			City:         domain.MeliLocation{Name: "Suba"},
			State:        domain.MeliLocation{Name: "Bogota D.C."},
			Country:      domain.MeliLocation{Name: "Colombia"},
		},
	}

	dto := MapMeliOrderToProbability(order, shipping, []byte(`{"raw":1}`))

	if dto.CustomerName != "Ana Perez" {
		t.Errorf("customer name: got %q", dto.CustomerName)
	}
	if dto.CustomerDNI != "999" {
		t.Errorf("dni: got %q", dto.CustomerDNI)
	}
	if dto.OrderTypeName != "full" {
		t.Errorf("order type: got %q", dto.OrderTypeName)
	}
	if len(dto.Shipments) != 1 {
		t.Fatalf("expected 1 shipment, got %d", len(dto.Shipments))
	}
	sh := dto.Shipments[0]
	if sh.GuideID == nil || *sh.GuideID != "777" {
		t.Errorf("guide id not set to shipment id")
	}
	if sh.TrackingNumber == nil || *sh.TrackingNumber != "TRACK9" {
		t.Errorf("tracking not mapped")
	}
	if len(dto.Addresses) != 1 || dto.Addresses[0].City != "Bogota D.C." {
		t.Errorf("bogota special-case not applied: %+v", dto.Addresses)
	}
}
