package client

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/woocommerce/internal/domain"
)

func newTestServer(status int, contentType, body string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if contentType != "" {
			w.Header().Set("Content-Type", contentType)
		}
		w.WriteHeader(status)
		_, _ = w.Write([]byte(body))
	}))
}

func TestTestConnection_HTMLComingSoon(t *testing.T) {
	srv := newTestServer(http.StatusOK, "text/html; charset=UTF-8", "<!DOCTYPE html><html><body>Coming soon</body></html>")
	defer srv.Close()

	err := New().TestConnection(context.Background(), srv.URL, "ck_x", "cs_x")
	if err == nil {
		t.Fatal("se esperaba error ante respuesta HTML (falso positivo), se obtuvo nil")
	}
}

func TestTestConnection_AuthErrorJSON(t *testing.T) {
	srv := newTestServer(http.StatusOK, "application/json", `{"code":"woocommerce_rest_authentication_error","message":"credenciales invalidas"}`)
	defer srv.Close()

	err := New().TestConnection(context.Background(), srv.URL, "ck_x", "cs_x")
	if !errors.Is(err, domain.ErrInvalidCredentials) {
		t.Fatalf("se esperaba ErrInvalidCredentials, se obtuvo: %v", err)
	}
}

func TestTestConnection_Unauthorized(t *testing.T) {
	srv := newTestServer(http.StatusUnauthorized, "application/json", `{"code":"woocommerce_rest_authentication_error"}`)
	defer srv.Close()

	err := New().TestConnection(context.Background(), srv.URL, "ck_x", "cs_x")
	if !errors.Is(err, domain.ErrInvalidCredentials) {
		t.Fatalf("se esperaba ErrInvalidCredentials ante 401, se obtuvo: %v", err)
	}
}

func TestTestConnection_ValidSystemStatus(t *testing.T) {
	srv := newTestServer(http.StatusOK, "application/json", `{"environment":{"home_url":"https://x"},"database":{}}`)
	defer srv.Close()

	err := New().TestConnection(context.Background(), srv.URL, "ck_x", "cs_x")
	if err != nil {
		t.Fatalf("se esperaba nil ante system_status valido, se obtuvo: %v", err)
	}
}
