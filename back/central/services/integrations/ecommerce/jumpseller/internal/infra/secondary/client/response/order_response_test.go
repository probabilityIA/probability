package response

import (
	"encoding/json"
	"testing"
)

const webhookOrderPayload = `{
  "order": {
    "id": 1012,
    "created_at": "2026-07-23 17:48:38 UTC",
    "completed_at": "2026-07-23 17:50:17 UTC",
    "currency": "CLP",
    "subtotal": 50,
    "tax": 0,
    "shipping_tax": 0,
    "shipping": 2877,
    "shipping_required": true,
    "total": 50,
    "discount": 2877,
    "shipping_discount": 2877,
    "gift_cards_discount": 0,
    "fulfillment_status": "unfulfilled",
    "shipping_method_id": 836621,
    "shipping_method_name": "BlueExpress",
    "payment_method_name": "Webpay Plus",
    "payment_method_type": "webpay_plus",
    "payment_information": "Tipo de Transaccion: Venta",
    "additional_information": null,
    "coupons": null,
    "promotions": [{ "id": 891331, "name": "SinEnvio", "code": null }],
    "customer": {
      "id": 20019293,
      "email": "rms.chile@gmail.com",
      "phone": null,
      "phone_prefix": null,
      "ip": "181.43.125.120",
      "fullname": "Rodrigo Morales"
    },
    "shipping_branch": { "id": null, "name": null },
    "shipping_address": {
      "name": "Rod",
      "surname": "Morales",
      "address": "Rodrigo de Quiroga 2823",
      "city": "Santiago",
      "postal": null,
      "region": "Region Metropolitana",
      "country": "Chile",
      "country_code": "CL",
      "region_code": "12",
      "street_number": null,
      "complement": "404",
      "latitude": null,
      "longitude": null,
      "municipality": "Vitacura"
    },
    "billing_address": {
      "name": "Rod",
      "surname": "Morales",
      "taxid": null,
      "address": "Rodrigo de Quiroga 2823",
      "city": "Santiago",
      "postal": null,
      "region": "Region Metropolitana",
      "country": "Chile",
      "country_code": "CL",
      "region_code": "12",
      "street_number": null,
      "complement": "404",
      "municipality": "Vitacura"
    },
    "pickup_address": null,
    "products": [
      {
        "id": 36335795,
        "variant_id": null,
        "sku": "100101999",
        "name": "Test eBox",
        "qty": 1,
        "price": 50.0,
        "tax": 0.0,
        "discount": 0.0,
        "weight": 0.5,
        "type": "physical",
        "taxes": [],
        "options": [],
        "stock_locations": [{ "location_id": 314455, "stock": 1 }],
        "files": []
      }
    ],
    "additional_fields": [
      { "value": "No", "label": "accepts_marketing", "id": null, "area": "contact" }
    ],
    "shipping_taxes": [],
    "status": "Paid",
    "status_name": "Paid",
    "status_enum": "paid",
    "tracking_url": null,
    "tracking_company": null,
    "tracking_number": null,
    "shipping_option": "delivery",
    "same_day_delivery": false,
    "shipment_status": "No Procesado",
    "shipment_status_enum": "unfulfilled",
    "recovered_from": null,
    "external_shipping_rate_id": null,
    "external_shipping_rate_description": null,
    "billing_information": null
  }
}`

func TestOrderEnvelopeUnmarshalRealWebhookPayload(t *testing.T) {
	var envelope OrderEnvelope
	if err := json.Unmarshal([]byte(webhookOrderPayload), &envelope); err != nil {
		t.Fatalf("unmarshal fallo: %v", err)
	}

	order := envelope.Order.ToDomain()

	if order.ID != 1012 {
		t.Errorf("ID = %d, esperado 1012", order.ID)
	}
	if order.Customer.ID != "20019293" {
		t.Errorf("Customer.ID = %q, esperado 20019293", order.Customer.ID)
	}
	if order.Customer.Name != "Rodrigo Morales" {
		t.Errorf("Customer.Name = %q, esperado Rodrigo Morales", order.Customer.Name)
	}
	if order.Customer.Phone != "" {
		t.Errorf("Customer.Phone = %q, esperado vacio", order.Customer.Phone)
	}
	if order.Status != "Paid" || order.StatusEnum != "paid" {
		t.Errorf("Status = %q / %q", order.Status, order.StatusEnum)
	}
	if order.ShipmentStatusEnum != "unfulfilled" {
		t.Errorf("ShipmentStatusEnum = %q, esperado unfulfilled", order.ShipmentStatusEnum)
	}
	if order.CompletedAt.IsZero() {
		t.Error("CompletedAt vacio, esperado 2026-07-23 17:50:17 UTC")
	}
	if order.ShippingAddress.Complement != "404" {
		t.Errorf("Complement = %q, esperado 404", order.ShippingAddress.Complement)
	}
	if order.ShippingAddress.Municipality != "Vitacura" {
		t.Errorf("Municipality = %q, esperado Vitacura", order.ShippingAddress.Municipality)
	}
	if len(order.Products) != 1 {
		t.Fatalf("Products = %d, esperado 1", len(order.Products))
	}
	if order.Products[0].SKU != "100101999" || order.Products[0].VariantID != 0 {
		t.Errorf("Producto SKU=%q VariantID=%d", order.Products[0].SKU, order.Products[0].VariantID)
	}
}

func TestOrderCustomerIDAsString(t *testing.T) {
	var envelope OrderEnvelope
	payload := `{"order":{"id":5,"customer":{"id":"abc-123","fullname":"Test"}}}`
	if err := json.Unmarshal([]byte(payload), &envelope); err != nil {
		t.Fatalf("unmarshal fallo: %v", err)
	}
	if string(envelope.Order.Customer.ID) != "abc-123" {
		t.Errorf("Customer.ID = %q, esperado abc-123", envelope.Order.Customer.ID)
	}
}
