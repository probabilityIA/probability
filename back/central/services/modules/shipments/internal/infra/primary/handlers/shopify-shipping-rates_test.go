package handlers

import (
	"testing"

	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/domain"
)

func TestBuildShopifyQuotePayload(t *testing.T) {
	var req shopifyRateRequest
	req.Rate.Destination = shopifyRateAddress{
		City:     "Bogota",
		Province: "Bogota D.C.",
		Name:     "Juan Perez Gomez",
		Address1: "Calle 100 # 15-20",
		Address2: "Apto 301",
		Phone:    "3001234567",
	}
	req.Rate.Items = []shopifyRateItem{
		{Name: "Camiseta", Quantity: 2, Grams: 600, Price: 5000000},
		{Name: "Gorra", Quantity: 1, Grams: 300, Price: 2000000},
	}

	origin := &domain.OriginAddress{
		Company:      "Mi Tienda",
		FirstName:    "Bodega",
		LastName:     "Central",
		Email:        "bodega@tienda.com",
		Phone:        "6041234567",
		Street:       "Carrera 48 # 10-20",
		Suburb:       "Sabaneta",
		CityDaneCode: "05631000",
	}

	payload := buildShopifyQuotePayload(req, origin, "11001")

	originMap := payload["origin"].(map[string]interface{})
	if originMap["daneCode"] != "05631000" {
		t.Fatalf("origin daneCode = %v, want 05631000", originMap["daneCode"])
	}

	destMap := payload["destination"].(map[string]interface{})
	if destMap["daneCode"] != "11001" {
		t.Fatalf("destination daneCode = %v, want 11001", destMap["daneCode"])
	}
	if destMap["firstName"] != "Juan" {
		t.Fatalf("destination firstName = %v, want Juan", destMap["firstName"])
	}
	if destMap["lastName"] != "Perez Gomez" {
		t.Fatalf("destination lastName = %v, want 'Perez Gomez'", destMap["lastName"])
	}
	if destMap["address"] != "Calle 100 # 15-20 Apto 301" {
		t.Fatalf("destination address = %v", destMap["address"])
	}

	contentValue := payload["contentValue"].(float64)
	if contentValue != 120000.0 {
		t.Fatalf("contentValue = %v, want 120000 (2*50000 + 1*20000)", contentValue)
	}

	pkgs := payload["packages"].([]interface{})
	if len(pkgs) != 1 {
		t.Fatalf("packages len = %d, want 1", len(pkgs))
	}
	pkg := pkgs[0].(map[string]interface{})
	if pkg["weight"].(float64) != 1.5 {
		t.Fatalf("package weight = %v, want 1.5 ((2*600 + 1*300)/1000)", pkg["weight"])
	}
}

func TestBuildShopifyQuotePayloadDefaultsWeight(t *testing.T) {
	var req shopifyRateRequest
	req.Rate.Items = []shopifyRateItem{{Name: "Digital", Quantity: 1, Grams: 0, Price: 1000}}

	payload := buildShopifyQuotePayload(req, &domain.OriginAddress{CityDaneCode: "11001"}, "05001")
	pkg := payload["packages"].([]interface{})[0].(map[string]interface{})
	if pkg["weight"].(float64) != 1 {
		t.Fatalf("weight = %v, want default 1", pkg["weight"])
	}
}

func TestMapQuoteRatesToShopify(t *testing.T) {
	rates := []interface{}{
		map[string]interface{}{
			"carrier":      "COORDINADORA",
			"product":      "Express",
			"flete":        float64(25000),
			"deliveryDays": float64(2),
		},
		map[string]interface{}{
			"carrier":      "ENVIA",
			"product":      "",
			"flete":        float64(18000),
			"deliveryDays": float64(0),
		},
		map[string]interface{}{
			"carrier": "INVALIDO",
			"flete":   float64(0),
		},
	}

	out := mapQuoteRatesToShopify(rates, "COP")

	if len(out) != 2 {
		t.Fatalf("got %d rates, want 2 (zero-flete dropped)", len(out))
	}

	if out[0].ServiceName != "COORDINADORA - Express" {
		t.Fatalf("service_name = %q", out[0].ServiceName)
	}
	if out[0].ServiceCode != "coordinadora_express_0" {
		t.Fatalf("service_code = %q", out[0].ServiceCode)
	}
	if out[0].TotalPrice != "2500000" {
		t.Fatalf("total_price = %q, want 2500000 (25000 * 100)", out[0].TotalPrice)
	}
	if out[0].Currency != "COP" {
		t.Fatalf("currency = %q", out[0].Currency)
	}
	if out[0].MinDeliveryDate == "" || out[0].MaxDeliveryDate == "" {
		t.Fatalf("expected delivery dates set when deliveryDays > 0")
	}

	if out[1].ServiceName != "ENVIA" {
		t.Fatalf("service_name = %q, want ENVIA (no product)", out[1].ServiceName)
	}
	if out[1].Description != "" || out[1].MinDeliveryDate != "" {
		t.Fatalf("expected no delivery info when deliveryDays = 0")
	}
}

func TestMapQuoteRatesToShopifyEmpty(t *testing.T) {
	if got := mapQuoteRatesToShopify(nil, "COP"); len(got) != 0 {
		t.Fatalf("nil rates -> %d, want 0", len(got))
	}
	if got := mapQuoteRatesToShopify("not-a-list", "COP"); len(got) != 0 {
		t.Fatalf("bad type -> %d, want 0", len(got))
	}
}
