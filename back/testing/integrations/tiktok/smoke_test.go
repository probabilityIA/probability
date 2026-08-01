package tiktok_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/secamc93/probability/back/testing/integrations/tiktok"
	"github.com/secamc93/probability/back/testing/shared/log"
)

func TestTiktokMockSmoke(t *testing.T) {
	logger := log.New()
	port := "19101"
	server := tiktok.New(logger, port)

	go func() {
		if err := server.Start(); err != nil {
			t.Logf("server stopped: %v", err)
		}
	}()

	baseURL := "http://localhost:" + port
	if !waitForHealth(baseURL+"/health", 5*time.Second) {
		t.Fatalf("tiktok mock did not become healthy on %s", baseURL)
	}

	form := url.Values{}
	form.Set("client_key", "mock-key")
	form.Set("client_secret", "mock-secret")
	form.Set("grant_type", "client_credentials")

	resp, err := http.Post(baseURL+"/v2/oauth/token/", "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatalf("token request failed: %v", err)
	}
	defer resp.Body.Close()

	var tokenResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		t.Fatalf("decoding token response: %v", err)
	}
	token, _ := tokenResp["access_token"].(string)
	if token == "" {
		t.Fatalf("expected access_token, got %v", tokenResp)
	}

	shopResp := queryShop(t, baseURL, token, "probability-mock-shop")
	data := shopResp["data"].(map[string]interface{})
	shops := data["shop_data"].([]interface{})
	if len(shops) != 1 {
		t.Fatalf("expected 1 shop, got %d", len(shops))
	}
	shop := shops[0].(map[string]interface{})
	if shop["shop_name"] != "probability-mock-shop" {
		t.Fatalf("unexpected shop_name: %v", shop["shop_name"])
	}
	if shop["shop_rating"] != "4.7" {
		t.Fatalf("unexpected shop_rating: %v", shop["shop_rating"])
	}

	emptyResp := queryShop(t, baseURL, token, "no-existe")
	emptyData := emptyResp["data"].(map[string]interface{})
	emptyShops := emptyData["shop_data"].([]interface{})
	if len(emptyShops) != 0 {
		t.Fatalf("expected 0 shops for unknown name, got %d", len(emptyShops))
	}

	badReq, _ := http.NewRequest(http.MethodPost, baseURL+"/v2/research/tts/shop/", bytes.NewReader([]byte(`{"shop_name":"x"}`)))
	badReq.Header.Set("Content-Type", "application/json")
	badReq.Header.Set("Authorization", "Bearer token-invalido")
	badResp, err := http.DefaultClient.Do(badReq)
	if err != nil {
		t.Fatalf("unauthorized request failed: %v", err)
	}
	defer badResp.Body.Close()
	if badResp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401 for invalid token, got %d", badResp.StatusCode)
	}
}

func queryShop(t *testing.T, baseURL, token, shopName string) map[string]interface{} {
	t.Helper()
	body, _ := json.Marshal(map[string]interface{}{
		"shop_name": shopName,
		"fields":    "shop_name,shop_rating,shop_review_count,item_sold_count,shop_id,shop_performance_value",
		"limit":     1,
	})
	req, _ := http.NewRequest(http.MethodPost, baseURL+"/v2/research/tts/shop/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("shop query failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var parsed map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		t.Fatalf("decoding shop response: %v", err)
	}
	return parsed
}

func waitForHealth(healthURL string, timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		resp, err := http.Get(healthURL)
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				return true
			}
		}
		time.Sleep(100 * time.Millisecond)
	}
	return false
}
