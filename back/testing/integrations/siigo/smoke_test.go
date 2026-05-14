package siigo_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/secamc93/probability/back/testing/integrations/siigo"
	"github.com/secamc93/probability/back/testing/shared/log"
)

func TestSiigoMockSmoke(t *testing.T) {
	logger := log.New()
	port := "19095"
	server := siigo.New(logger, port)

	go func() {
		if err := server.Start(); err != nil {
			t.Logf("server stopped: %v", err)
		}
	}()

	baseURL := "http://localhost:" + port
	if !waitForHealth(baseURL+"/health", 5*time.Second) {
		t.Fatalf("siigo mock did not become healthy on %s", baseURL)
	}

	authBody := map[string]string{
		"username":   "demo@probability.com",
		"access_key": "test-key",
	}
	authResp := postJSON(t, baseURL+"/v1/auth", "", authBody)
	token, _ := authResp["access_token"].(string)
	if token == "" {
		t.Fatalf("expected access_token, got %v", authResp)
	}

	custBody := map[string]interface{}{
		"identification": "900123456",
		"name":           []string{"Probability", "Demo"},
		"email":          "demo@probability.com",
	}
	custResp := postJSON(t, baseURL+"/v1/customers", token, custBody)
	if id, _ := custResp["id"].(string); id == "" {
		t.Fatalf("expected customer id, got %v", custResp)
	}

	invBody := map[string]interface{}{
		"document": map[string]interface{}{"id": float64(24446)},
		"date":     time.Now().Format("2006-01-02"),
		"customer": map[string]interface{}{
			"identification": "900123456",
		},
		"items": []interface{}{
			map[string]interface{}{
				"code":        "ITEM-1",
				"description": "Producto demo",
				"quantity":    float64(2),
				"price":       float64(50000),
			},
		},
	}
	invResp := postJSON(t, baseURL+"/v1/invoices", token, invBody)
	cufe, _ := invResp["cufe"].(string)
	if len(cufe) != 64 {
		t.Fatalf("expected 64-char cufe, got %q (%v)", cufe, invResp)
	}
	if _, ok := invResp["document"].(map[string]interface{}); !ok {
		t.Fatalf("missing document.prefix/number: %v", invResp)
	}
}

func waitForHealth(url string, timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		resp, err := http.Get(url)
		if err == nil {
			_ = resp.Body.Close()
			if resp.StatusCode == 200 {
				return true
			}
		}
		time.Sleep(100 * time.Millisecond)
	}
	return false
}

func postJSON(t *testing.T, url, token string, body interface{}) map[string]interface{} {
	t.Helper()
	raw, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(raw))
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Partner-Id", "test-partner")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("post %s: %v", url, err)
	}
	defer resp.Body.Close()
	data, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		t.Fatalf("post %s -> %d: %s", url, resp.StatusCode, string(data))
	}
	var out map[string]interface{}
	if err := json.Unmarshal(data, &out); err != nil {
		t.Fatalf("decode %s: %v body=%s", url, err, string(data))
	}
	fmt.Printf("[smoke] POST %s -> %d\n", url, resp.StatusCode)
	return out
}
