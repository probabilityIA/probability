package mapper

import (
	"encoding/json"
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

	dto := MapMeliOrderToProbability(order, shipping, []byte(`{"raw":1}`), nil)

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

func TestMapMeliOrderToProbability_AdjuntaShipmentAlRaw(t *testing.T) {
	order := &domain.MeliOrder{
		ID:          2000014575744355,
		Status:      "paid",
		CurrencyID:  "COP",
		TotalAmount: 109400,
	}
	orderRaw := []byte(`{"id":2000014575744355,"shipping":{"id":47791295504}}`)
	shipmentRaw := []byte(`{"id":47791295504,"receiver_address":{"city":{"name":"BOGOTA"}}}`)

	dto := MapMeliOrderToProbability(order, nil, orderRaw, shipmentRaw)

	if dto.ChannelMetadata == nil {
		t.Fatal("no se genero channel metadata")
	}

	var got map[string]json.RawMessage
	if err := json.Unmarshal(dto.ChannelMetadata.RawData, &got); err != nil {
		t.Fatalf("raw_data no es JSON valido: %v", err)
	}
	if _, ok := got["shipment_detail"]; !ok {
		t.Fatal("raw_data no incluye shipment_detail")
	}
	if _, ok := got["shipping"]; !ok {
		t.Error("raw_data perdio las claves originales de la orden")
	}
}

func TestMapMeliOrderToProbability_SinShipmentRawDejaElOrderRaw(t *testing.T) {
	order := &domain.MeliOrder{ID: 123, Status: "paid", CurrencyID: "COP"}
	orderRaw := []byte(`{"id":123}`)

	dto := MapMeliOrderToProbability(order, nil, orderRaw, nil)

	if dto.ChannelMetadata == nil {
		t.Fatal("no se genero channel metadata")
	}
	if string(dto.ChannelMetadata.RawData) != string(orderRaw) {
		t.Errorf("raw_data modificado sin necesidad: got %s, want %s", dto.ChannelMetadata.RawData, orderRaw)
	}
}

func TestMapMeliOrderToProbability_GuardaLosIDsDelPack(t *testing.T) {
	packID := int64(2000014575744355)
	order := &domain.MeliOrder{
		ID:           packID,
		Status:       "paid",
		CurrencyID:   "COP",
		TotalAmount:  109400,
		PackID:       &packID,
		PackOrderIDs: []int64{2000017981110684, 2000017981120028},
	}
	orderRaw := []byte(`{"id":2000014575744355}`)

	dto := MapMeliOrderToProbability(order, nil, orderRaw, nil)

	if dto.ExternalID != "2000014575744355" {
		t.Errorf("ExternalID: got %q, want el pack_id", dto.ExternalID)
	}

	var got map[string]json.RawMessage
	if err := json.Unmarshal(dto.ChannelMetadata.RawData, &got); err != nil {
		t.Fatalf("raw_data invalido: %v", err)
	}
	ids, ok := got["pack_order_ids"]
	if !ok {
		t.Fatal("raw_data no incluye pack_order_ids")
	}
	if string(ids) != "[2000017981110684,2000017981120028]" {
		t.Errorf("pack_order_ids: got %s", ids)
	}
}

func TestMapMeliOrderToProbability_SinPackNoAgregaIDs(t *testing.T) {
	order := &domain.MeliOrder{ID: 123, Status: "paid", CurrencyID: "COP"}
	orderRaw := []byte(`{"id":123}`)

	dto := MapMeliOrderToProbability(order, nil, orderRaw, nil)

	var got map[string]json.RawMessage
	if err := json.Unmarshal(dto.ChannelMetadata.RawData, &got); err != nil {
		t.Fatalf("raw_data invalido: %v", err)
	}
	if _, ok := got["pack_order_ids"]; ok {
		t.Error("una orden sin pack no deberia traer pack_order_ids")
	}
}
