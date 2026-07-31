package client_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/tiktok/internal/domain"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/tiktok/internal/infra/secondary/client"
	"github.com/secamc93/probability/back/central/shared/log"
)

const testToken = "clt.test-token"

func newMockServer(t *testing.T) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()

	mux.HandleFunc("/v2/oauth/token/", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if r.PostFormValue("client_key") == "" || r.PostFormValue("client_secret") == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error":             "invalid_request",
				"error_description": "missing credentials",
			})
			return
		}
		if r.PostFormValue("client_key") == "bad-key" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error":             "invalid_client",
				"error_description": "client_key is invalid",
			})
			return
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"access_token": testToken,
			"expires_in":   7200,
			"token_type":   "Bearer",
		})
	})

	mux.HandleFunc("/v2/research/tts/shop/", func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.Header.Get("Authorization"), testToken) {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"data":  map[string]interface{}{},
				"error": map[string]string{"code": "access_token_invalid", "message": "invalid token"},
			})
			return
		}
		var req struct {
			ShopName string `json:"shop_name"`
		}
		json.NewDecoder(r.Body).Decode(&req)
		shops := []map[string]interface{}{}
		if req.ShopName == "mitienda" {
			shops = append(shops, map[string]interface{}{
				"shop_name":              "mitienda",
				"shop_rating":            "4.5",
				"shop_review_count":      int64(100),
				"item_sold_count":        int64(2500),
				"shop_id":                int64(123456789),
				"shop_performance_value": int64(90),
			})
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data":  map[string]interface{}{"shop_data": shops},
			"error": map[string]string{"code": "ok", "message": ""},
		})
	})

	return httptest.NewServer(mux)
}

func TestGetClientToken(t *testing.T) {
	server := newMockServer(t)
	defer server.Close()

	c := client.New(log.New())

	token, err := c.GetClientToken(context.Background(), server.URL, "good-key", "good-secret")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if token != testToken {
		t.Fatalf("expected %q, got %q", testToken, token)
	}
}

func TestGetClientTokenInvalid(t *testing.T) {
	server := newMockServer(t)
	defer server.Close()

	c := client.New(log.New())

	_, err := c.GetClientToken(context.Background(), server.URL, "bad-key", "any")
	if err == nil {
		t.Fatal("expected error for invalid client_key")
	}
	if !strings.Contains(err.Error(), "invalid_client") {
		t.Fatalf("expected invalid_client error, got: %v", err)
	}
}

func TestQueryShopInfo(t *testing.T) {
	server := newMockServer(t)
	defer server.Close()

	c := client.New(log.New())
	cred := domain.Credential{AccessToken: testToken, BaseURL: server.URL}

	shops, err := c.QueryShopInfo(context.Background(), cred, "mitienda", 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(shops) != 1 {
		t.Fatalf("expected 1 shop, got %d", len(shops))
	}
	if shops[0].ShopName != "mitienda" || shops[0].ShopRating != "4.5" || shops[0].ShopID != 123456789 {
		t.Fatalf("unexpected shop data: %+v", shops[0])
	}
}

func TestQueryShopInfoNotFound(t *testing.T) {
	server := newMockServer(t)
	defer server.Close()

	c := client.New(log.New())
	cred := domain.Credential{AccessToken: testToken, BaseURL: server.URL}

	shops, err := c.QueryShopInfo(context.Background(), cred, "desconocida", 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(shops) != 0 {
		t.Fatalf("expected 0 shops, got %d", len(shops))
	}
}

func TestQueryShopInfoUnauthorized(t *testing.T) {
	server := newMockServer(t)
	defer server.Close()

	c := client.New(log.New())
	cred := domain.Credential{AccessToken: "otro-token", BaseURL: server.URL}

	_, err := c.QueryShopInfo(context.Background(), cred, "mitienda", 1)
	if err == nil {
		t.Fatal("expected error for invalid token")
	}
}
