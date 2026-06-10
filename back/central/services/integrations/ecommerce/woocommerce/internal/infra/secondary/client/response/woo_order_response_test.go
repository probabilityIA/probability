package response

import (
	"encoding/json"
	"testing"
)

func TestWooLineItemImageIDFlexible(t *testing.T) {
	cases := []struct {
		name string
		raw  string
		want int64
	}{
		{"id como string vacio", `{"id":"","src":"http://x/img.jpg"}`, 0},
		{"id como string numerico", `{"id":"42","src":"http://x/img.jpg"}`, 42},
		{"id como numero", `{"id":42,"src":"http://x/img.jpg"}`, 42},
		{"id null", `{"id":null,"src":""}`, 0},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var img WooLineItemImage
			if err := json.Unmarshal([]byte(tc.raw), &img); err != nil {
				t.Fatalf("unmarshal: %v", err)
			}
			if int64(img.ID) != tc.want {
				t.Fatalf("ID = %d, want %d", img.ID, tc.want)
			}
		})
	}
}

func TestWooOrderResponseWithStringImageID(t *testing.T) {
	raw := `{
		"id": 12,
		"number": "12",
		"status": "processing",
		"total": "142000.00",
		"line_items": [{
			"id": 1,
			"name": "Camiseta Test Probability",
			"product_id": 11,
			"quantity": 2,
			"price": 65000,
			"image": {"id": "", "src": ""}
		}]
	}`
	var resp WooOrderResponse
	if err := json.Unmarshal([]byte(raw), &resp); err != nil {
		t.Fatalf("payload de webhook con image.id string debe deserializar: %v", err)
	}
	if len(resp.LineItems) != 1 || resp.LineItems[0].ProductID != 11 {
		t.Fatalf("line items mal parseados: %+v", resp.LineItems)
	}
}
