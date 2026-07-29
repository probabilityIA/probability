package shipit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/testing/integrations/shipit/internal/app/usecases"
	"github.com/secamc93/probability/back/testing/integrations/shipit/internal/infra/primary/handlers"
	"github.com/secamc93/probability/back/testing/shared/log"
)

type webhookCapture struct {
	mu       sync.Mutex
	payloads []map[string]any
}

func (w *webhookCapture) handler(rw http.ResponseWriter, r *http.Request) {
	var payload map[string]any
	json.NewDecoder(r.Body).Decode(&payload)
	w.mu.Lock()
	w.payloads = append(w.payloads, payload)
	w.mu.Unlock()
	rw.WriteHeader(200)
}

func (w *webhookCapture) count() int {
	w.mu.Lock()
	defer w.mu.Unlock()
	return len(w.payloads)
}

func newMockServer(t *testing.T, webhookTarget string) (*httptest.Server, *usecases.APISimulator) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	logger := log.New()

	simulator := usecases.NewAPISimulator(logger, webhookTarget, "http://mock.local")
	handler := handlers.New(simulator, logger)

	router := gin.New()
	handler.RegisterRoutes(router)

	return httptest.NewServer(router), simulator
}

func authedRequest(t *testing.T, method, url string, body any) *http.Response {
	t.Helper()
	var reader *bytes.Reader
	if body != nil {
		data, _ := json.Marshal(body)
		reader = bytes.NewReader(data)
	} else {
		reader = bytes.NewReader(nil)
	}
	req, _ := http.NewRequest(method, url, reader)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/vnd.shipit.v4")
	req.Header.Set("X-Shipit-Email", "mock@test.cl")
	req.Header.Set("X-Shipit-Access-Token", "mock-token")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	return resp
}

func TestMockFullFlow(t *testing.T) {
	capture := &webhookCapture{}
	webhookServer := httptest.NewServer(http.HandlerFunc(capture.handler))
	defer webhookServer.Close()

	server, _ := newMockServer(t, webhookServer.URL)
	defer server.Close()

	resp := authedRequest(t, "GET", server.URL+"/v/communes", nil)
	if resp.StatusCode != 200 {
		t.Fatalf("communes status: %d", resp.StatusCode)
	}
	var communes []map[string]any
	json.NewDecoder(resp.Body).Decode(&communes)
	resp.Body.Close()
	if len(communes) == 0 {
		t.Fatal("expected communes")
	}

	resp = authedRequest(t, "POST", server.URL+"/v/rates", map[string]any{
		"parcel": map[string]any{
			"length": 30, "width": 20, "height": 10, "weight": 1.5,
			"origin_id": 309, "destiny_id": 308, "type_of_destiny": "domicilio",
		},
	})
	if resp.StatusCode != 201 {
		t.Fatalf("rates status: %d", resp.StatusCode)
	}
	var rates map[string]any
	json.NewDecoder(resp.Body).Decode(&rates)
	resp.Body.Close()
	if prices, ok := rates["prices"].([]any); !ok || len(prices) == 0 {
		t.Fatalf("expected prices, got %v", rates)
	}

	resp = authedRequest(t, "POST", server.URL+"/v/shipments", map[string]any{
		"shipment": map[string]any{
			"kind": 0, "platform": 2, "reference": "ORD-SMOKE-1", "items": 1, "sandbox": true,
			"sizes": map[string]any{"width": 20, "height": 10, "length": 30, "weight": 1.5},
			"destiny": map[string]any{
				"street": "Apoquindo", "number": "5555",
				"commune_id": 308, "commune_name": "LAS CONDES",
				"full_name": "Maria Gonzalez", "email": "maria@test.cl",
				"phone": "56911111111", "kind": "home_delivery",
			},
			"courier":   map[string]any{"id": 0, "client": "starken", "selected": true, "payable": false, "algorithm": 1},
			"insurance": map[string]any{"ticket_number": "ORD-SMOKE-1", "ticket_amount": 25000, "detail": "Zapatos", "extra": false},
		},
	})
	if resp.StatusCode != 200 {
		t.Fatalf("create shipment status: %d", resp.StatusCode)
	}
	var shipment map[string]any
	json.NewDecoder(resp.Body).Decode(&shipment)
	resp.Body.Close()

	tracking, _ := shipment["tracking_number"].(string)
	if !strings.HasPrefix(tracking, "SHPT") {
		t.Fatalf("unexpected tracking number: %v", shipment["tracking_number"])
	}
	if shipment["ticket_shipit_pdf_url"] == nil {
		t.Fatal("expected label url")
	}

	resp = authedRequest(t, "GET", server.URL+"/v/trackings/number/"+tracking, nil)
	if resp.StatusCode != 200 {
		t.Fatalf("tracking status: %d", resp.StatusCode)
	}
	var trackingResp map[string]any
	json.NewDecoder(resp.Body).Decode(&trackingResp)
	resp.Body.Close()
	data, _ := trackingResp["data"].(map[string]any)
	if data == nil || data["number"] != tracking {
		t.Fatalf("unexpected tracking response: %v", trackingResp)
	}

	simBody, _ := json.Marshal(map[string]any{"tracking_number": tracking, "status": "in_transit"})
	simResp, err := http.Post(server.URL+"/simulate/status", "application/json", bytes.NewReader(simBody))
	if err != nil || simResp.StatusCode != 200 {
		t.Fatalf("simulate status failed: %v %d", err, simResp.StatusCode)
	}
	simResp.Body.Close()

	deadline := time.Now().Add(3 * time.Second)
	for capture.count() == 0 && time.Now().Before(deadline) {
		time.Sleep(50 * time.Millisecond)
	}
	if capture.count() == 0 {
		t.Fatal("expected webhook delivery to target")
	}

	capture.mu.Lock()
	webhook := capture.payloads[0]
	capture.mu.Unlock()
	if webhook["tracking_number"] != tracking || webhook["status"] != "in_transit" {
		t.Fatalf("unexpected webhook payload: %v", webhook)
	}

	resp = authedRequest(t, "GET", server.URL+"/v/trackings/number/"+tracking, nil)
	json.NewDecoder(resp.Body).Decode(&trackingResp)
	resp.Body.Close()
	data, _ = trackingResp["data"].(map[string]any)
	statuses, _ := data["statuses"].([]any)
	if len(statuses) != 2 {
		t.Fatalf("expected 2 statuses after simulate, got %d", len(statuses))
	}
}

func TestMockRejectsMissingAuth(t *testing.T) {
	server, _ := newMockServer(t, "")
	defer server.Close()

	req, _ := http.NewRequest("GET", server.URL+"/v/communes", nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 401 {
		t.Fatalf("expected 401 without auth headers, got %d", resp.StatusCode)
	}
}

func TestMockUnknownTracking(t *testing.T) {
	server, _ := newMockServer(t, "")
	defer server.Close()

	resp := authedRequest(t, "GET", fmt.Sprintf("%s/v/trackings/number/NOPE", server.URL), nil)
	defer resp.Body.Close()
	if resp.StatusCode != 400 {
		t.Fatalf("expected 400 for unknown tracking, got %d", resp.StatusCode)
	}
}
