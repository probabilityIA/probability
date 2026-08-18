package response

import (
	"encoding/json"
	"testing"
)

// JSON real devuelto por GET /shipments/{id} (seller 228589645, agosto 2026):
// la direccion llega en destination.shipping_address, no en receiver_address.
const shipmentConDestination = `{
  "id": 47791295504,
  "status": "delivered",
  "tracking_number": "MEL47791295504FMDOF01",
  "destination": {
    "type": "agency",
    "receiver_name": "Jose de Jesus Mejia Rodriguez",
    "receiver_phone": "XXXXXXX",
    "shipping_address": {
      "address_line": "CALLE 6 15-03",
      "street_name": "CALLE 6",
      "street_number": "15-03",
      "zip_code": "205010",
      "city": { "name": "Aguachica" },
      "state": { "id": "CO-CES", "name": "Cesar" },
      "country": { "id": "CO", "name": "Colombia" },
      "neighborhood": { "name": "Olaya Herrera" },
      "latitude": 8.308608,
      "longitude": -73.6192217,
      "comment": "Almacen local esquina"
    }
  }
}`

func TestToDomain_DireccionDesdeDestination(t *testing.T) {
	var resp MeliShippingDetailResponse
	if err := json.Unmarshal([]byte(shipmentConDestination), &resp); err != nil {
		t.Fatalf("no parseo el shipment: %v", err)
	}

	detail := resp.ToDomain()

	if detail.ReceiverAddress == nil {
		t.Fatal("ReceiverAddress quedo nil: no se leyo destination.shipping_address")
	}
	addr := detail.ReceiverAddress
	if addr.StreetName != "CALLE 6" {
		t.Errorf("StreetName: got %q, want %q", addr.StreetName, "CALLE 6")
	}
	if addr.StreetNumber != "15-03" {
		t.Errorf("StreetNumber: got %q, want %q", addr.StreetNumber, "15-03")
	}
	if addr.City.Name != "Aguachica" {
		t.Errorf("City: got %q, want %q", addr.City.Name, "Aguachica")
	}
	if addr.State.Name != "Cesar" {
		t.Errorf("State: got %q, want %q", addr.State.Name, "Cesar")
	}
	if addr.ZipCode != "205010" {
		t.Errorf("ZipCode: got %q, want %q", addr.ZipCode, "205010")
	}
	if addr.Latitude == nil || *addr.Latitude != 8.308608 {
		t.Errorf("Latitude: got %v, want 8.308608", addr.Latitude)
	}
	if addr.Longitude == nil || *addr.Longitude != -73.6192217 {
		t.Errorf("Longitude: got %v, want -73.6192217", addr.Longitude)
	}
	if addr.ReceiverName != "Jose de Jesus Mejia Rodriguez" {
		t.Errorf("ReceiverName: got %q", addr.ReceiverName)
	}
	if addr.Neighborhood == nil || addr.Neighborhood.Name != "Olaya Herrera" {
		t.Errorf("Neighborhood: got %v", addr.Neighborhood)
	}
}

func TestToDomain_ReceiverAddressTienePrioridad(t *testing.T) {
	raw := `{
      "id": 1,
      "receiver_address": { "street_name": "VIEJA", "city": { "name": "Bogota" } },
      "destination": { "shipping_address": { "street_name": "NUEVA", "city": { "name": "Cali" } } }
    }`

	var resp MeliShippingDetailResponse
	if err := json.Unmarshal([]byte(raw), &resp); err != nil {
		t.Fatalf("no parseo: %v", err)
	}

	detail := resp.ToDomain()

	if detail.ReceiverAddress == nil || detail.ReceiverAddress.StreetName != "VIEJA" {
		t.Errorf("se esperaba conservar receiver_address cuando viene, got %v", detail.ReceiverAddress)
	}
}
