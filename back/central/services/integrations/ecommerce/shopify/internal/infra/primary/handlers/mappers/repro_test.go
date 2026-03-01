package mapper

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/shopify/internal/infra/secondary/client/response"
)

func TestReproMapping(t *testing.T) {
	// JSON provided by user (truncated for brevity, keeping relevant parts)
	jsonStr := `
	{
		"id": 5166661337297,
		"name": "#65954",
		"currency": "USD",
		"total_price": "88.76",
		"total_price_set": {
			"shop_money": {
				"amount": "88.76",
				"currency_code": "USD"
			},
			"presentment_money": {
				"amount": "334300.00",
				"currency_code": "COP"
			}
		},
		"shipping_address": {
			"address1": "Cll 118 # 52b-03",
			"address2": "Apt 104 edificio andino",
			"city": "Bogota",
			"country": "Colombia"
		}
	}
	`

	// 1. Unmarshal into response.Order (like webhook_mapper does via MapWebhookPayloadToOrderResponse logic)
	// webhook_mapper takes map, marshals to bytes, unmarshals to struct.
	// We simulate the byte unmarshal directly.
	var orderResp response.Order
	if err := json.Unmarshal([]byte(jsonStr), &orderResp); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	// 2. Verify Address2
	if orderResp.ShippingAddress == nil {
		t.Fatal("ShippingAddress is nil")
	}
	if orderResp.ShippingAddress.Address2 == nil {
		t.Error("Address2 is nil")
	} else {
		fmt.Printf("Address2 encoded: '%s'\n", *orderResp.ShippingAddress.Address2)
		if *orderResp.ShippingAddress.Address2 != "Apt 104 edificio andino" {
			t.Errorf("Address2 mismatch. Got '%s', want 'Apt 104 edificio andino'", *orderResp.ShippingAddress.Address2)
		}
	}

	// 3. Verify Presentment Money
	if orderResp.TotalPriceSet == nil {
		t.Fatal("TotalPriceSet is nil")
	}
	fmt.Printf("Presentment Amount: '%s'\n", orderResp.TotalPriceSet.PresentmentMoney.Amount)
	fmt.Printf("Presentment Currency: '%s'\n", orderResp.TotalPriceSet.PresentmentMoney.CurrencyCode)

	if orderResp.TotalPriceSet.PresentmentMoney.Amount != "334300.00" {
		t.Errorf("Presentment Amount mismatch. Got '%s'", orderResp.TotalPriceSet.PresentmentMoney.Amount)
	}
}
