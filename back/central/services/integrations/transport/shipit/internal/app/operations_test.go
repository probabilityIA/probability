package app

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/secamc93/probability/back/central/services/integrations/transport/shipit/internal/domain"
	clientinfra "github.com/secamc93/probability/back/central/services/integrations/transport/shipit/internal/infra/secondary/client"
	"github.com/secamc93/probability/back/central/shared/log"
)

func newTestLogger() log.ILogger {
	zerolog.SetGlobalLevel(zerolog.Disabled)
	return log.New()
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

func newShipitAPIStub(t *testing.T) *httptest.Server {
	t.Helper()

	mux := http.NewServeMux()

	requireAuth := func(w http.ResponseWriter, r *http.Request) bool {
		if r.Header.Get("X-Shipit-Email") == "" || r.Header.Get("X-Shipit-Access-Token") == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"state": "error", "message": "unauthorized"})
			return false
		}
		return true
	}

	mux.HandleFunc("/v/communes", func(w http.ResponseWriter, r *http.Request) {
		if !requireAuth(w, r) {
			return
		}
		writeJSON(w, []map[string]any{
			{"id": 308, "name": "LAS CONDES"},
			{"id": 309, "name": "SANTIAGO"},
		})
	})

	mux.HandleFunc("/v/rates", func(w http.ResponseWriter, r *http.Request) {
		if !requireAuth(w, r) {
			return
		}
		var req domain.RatesRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Parcel.OriginID == 0 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"state": "error", "message": "invalid parcel"})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]any{
			"algorithm": "1",
			"prices": []map[string]any{
				{"courier": map[string]string{"name": "starken"}, "name": "DIA HABIL SIGUIENTE", "price": 3803, "days": 1, "available_to_shipping": true},
				{"courier": map[string]string{"name": "chilexpress"}, "name": "EXPRESS", "price": 4200, "days": 1, "available_to_shipping": false},
			},
		})
	})

	mux.HandleFunc("/v/shipments", func(w http.ResponseWriter, r *http.Request) {
		if !requireAuth(w, r) {
			return
		}
		var req domain.ShipmentRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Shipment.Reference == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"state": "error", "message": "reference required"})
			return
		}
		if req.Shipment.Destiny.CommuneID == 0 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"state": "error", "message": "commune required"})
			return
		}
		tracking := "SHPT000123"
		ticket := "https://shipit.example/labels/123.pdf"
		writeJSON(w, map[string]any{
			"id":                    1544954,
			"reference":             req.Shipment.Reference,
			"status":                "created",
			"tracking_number":       tracking,
			"shipit_code":           "Shipit / " + req.Shipment.Reference,
			"courier_for_client":    "starken",
			"ticket_shipit_pdf_url": ticket,
		})
	})

	mux.HandleFunc("/v/trackings/number/SHPT000123", func(w http.ResponseWriter, r *http.Request) {
		if !requireAuth(w, r) {
			return
		}
		writeJSON(w, map[string]any{
			"data": map[string]any{
				"id":           1,
				"courier":      "starken",
				"number":       "SHPT000123",
				"is_delivered": false,
				"statuses": []map[string]any{
					{"id": 1, "generic_status": 0, "courier_status": "in_preparation", "created_at": "2026-07-29T10:00:00"},
					{"id": 2, "generic_status": 1, "courier_status": "in_transit", "created_at": "2026-07-29T12:00:00"},
				},
			},
		})
	})

	return httptest.NewServer(mux)
}

func testGuideRequest() domain.GuideRequest {
	return domain.GuideRequest{
		MyShipmentReference: "ORD-2026-000045",
		Description:         "Zapatos",
		ContentValue:        25000,
		Packages:            []domain.Package{{Weight: 1.5, Height: 10, Width: 20, Length: 30}},
		Origin: domain.Address{
			City:    "Santiago",
			Address: "Alameda 1000",
		},
		Destination: domain.Address{
			FirstName: "Maria",
			LastName:  "Gonzalez",
			Email:     "maria@test.cl",
			Phone:     "56911111111",
			Address:   "Apoquindo 5555",
			City:      "Las Condes",
		},
	}
}

func TestQuoteAgainstStub(t *testing.T) {
	server := newShipitAPIStub(t)
	defer server.Close()

	uc := New(clientinfra.New(newTestLogger()), newTestLogger())
	creds := domain.Credentials{Email: "test@test.cl", AccessToken: "token"}

	resp, err := uc.Quote(context.Background(), server.URL, creds, testGuideRequest())
	if err != nil {
		t.Fatalf("Quote failed: %v", err)
	}
	if len(resp.Data.Rates) != 1 {
		t.Fatalf("expected 1 available rate, got %d", len(resp.Data.Rates))
	}
	rate := resp.Data.Rates[0]
	if rate.Carrier != "starken" || rate.Flete != 3803 || rate.DeliveryDays != 1 {
		t.Fatalf("unexpected rate: %+v", rate)
	}
}

func TestGenerateAgainstStub(t *testing.T) {
	server := newShipitAPIStub(t)
	defer server.Close()

	uc := New(clientinfra.New(newTestLogger()), newTestLogger())
	creds := domain.Credentials{Email: "test@test.cl", AccessToken: "token"}

	resp, err := uc.Generate(context.Background(), server.URL, creds, testGuideRequest(), true)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}
	if resp.Data.TrackingNumber != "SHPT000123" {
		t.Fatalf("expected tracker SHPT000123, got %s", resp.Data.TrackingNumber)
	}
	if resp.Data.LabelURL != "https://shipit.example/labels/123.pdf" {
		t.Fatalf("unexpected label url: %s", resp.Data.LabelURL)
	}
	if resp.Data.Carrier != "starken" {
		t.Fatalf("unexpected carrier: %s", resp.Data.Carrier)
	}
	if resp.Data.IDOrder != 1544954 {
		t.Fatalf("unexpected idOrder: %d", resp.Data.IDOrder)
	}
}

func TestTrackAgainstStub(t *testing.T) {
	server := newShipitAPIStub(t)
	defer server.Close()

	uc := New(clientinfra.New(newTestLogger()), newTestLogger())
	creds := domain.Credentials{Email: "test@test.cl", AccessToken: "token"}

	out, err := uc.Track(context.Background(), server.URL, creds, "SHPT000123")
	if err != nil {
		t.Fatalf("Track failed: %v", err)
	}
	if out.Status != "in_transit" {
		t.Fatalf("expected in_transit, got %s", out.Status)
	}
	if len(out.History) != 2 {
		t.Fatalf("expected 2 history events, got %d", len(out.History))
	}
}

func TestGenerateFailsWithoutCommune(t *testing.T) {
	server := newShipitAPIStub(t)
	defer server.Close()

	uc := New(clientinfra.New(newTestLogger()), newTestLogger())
	creds := domain.Credentials{Email: "test@test.cl", AccessToken: "token"}

	req := testGuideRequest()
	req.Destination.City = "Ciudad Inexistente"
	req.Destination.Suburb = ""
	req.Destination.State = ""

	_, err := uc.Generate(context.Background(), server.URL, creds, req, true)
	if err == nil {
		t.Fatal("expected commune resolution error, got nil")
	}
}

func TestReferenceTruncation(t *testing.T) {
	req := testGuideRequest()
	body := buildShipmentBody(req, 308, "LAS CONDES", true)
	if len(body.Reference) > 15 {
		t.Fatalf("reference exceeds 15 chars: %q", body.Reference)
	}
	if body.Destiny.CommuneID != 308 {
		t.Fatalf("unexpected commune id: %d", body.Destiny.CommuneID)
	}
	if !body.Sandbox {
		t.Fatal("expected sandbox true")
	}
}

func TestWebhookPayloadMapping(t *testing.T) {
	payload := domain.WebhookPayload{
		ID:               1,
		Reference:        "ORD-1",
		Status:           "in_transit",
		CourierStatus:    "En camino a destino",
		TrackingNumber:   "SHPT000123",
		CourierForClient: "starken",
		UpdatedAt:        "2026-07-29T12:00:00Z",
	}

	update := payload.ToNormalizedUpdate()
	if update == nil {
		t.Fatal("expected normalized update, got nil")
	}
	if update.ProbabilityStatus != domain.StatusInTransit {
		t.Fatalf("expected in_transit, got %s", update.ProbabilityStatus)
	}
	if update.ShippedAt == nil {
		t.Fatal("expected shipped_at to be set")
	}

	delivered := payload
	delivered.Status = "delivered"
	update = delivered.ToNormalizedUpdate()
	if update.ProbabilityStatus != domain.StatusDelivered {
		t.Fatalf("expected delivered, got %s", update.ProbabilityStatus)
	}
	if update.DeliveredAt == nil {
		t.Fatal("expected delivered_at to be set")
	}

	empty := domain.WebhookPayload{Status: "delivered"}
	if empty.ToNormalizedUpdate() != nil {
		t.Fatal("expected nil update without tracking number")
	}
}

func TestMapShipitStatusFallback(t *testing.T) {
	status, known := domain.MapShipitStatus("estado raro nunca visto")
	if known {
		t.Fatal("expected unknown status")
	}
	if status != domain.StatusInTransit {
		t.Fatalf("expected in_transit fallback, got %s", status)
	}

	status, known = domain.MapShipitStatus("Entregado al destinatario")
	if !known || status != domain.StatusDelivered {
		t.Fatalf("expected delivered, got %s known=%v", status, known)
	}
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
