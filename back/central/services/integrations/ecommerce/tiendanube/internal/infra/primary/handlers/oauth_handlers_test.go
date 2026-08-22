package handlers

import "testing"

func TestParseTokenResponse_UserIDNumerico(t *testing.T) {
	raw := []byte(`{"access_token":"abc123","token_type":"bearer","scope":"read_orders write_products","user_id":1234567}`)

	got, err := parseTokenResponse(raw)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if got.AccessToken != "abc123" {
		t.Errorf("access_token = %q, se esperaba abc123", got.AccessToken)
	}
	if got.UserID != "1234567" {
		t.Errorf("user_id = %q, se esperaba 1234567", got.UserID)
	}
	if got.Scope != "read_orders write_products" {
		t.Errorf("scope = %q inesperado", got.Scope)
	}
}

func TestParseTokenResponse_UserIDComoString(t *testing.T) {
	raw := []byte(`{"access_token":"abc123","token_type":"bearer","user_id":"1234567"}`)

	got, err := parseTokenResponse(raw)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if got.UserID != "1234567" {
		t.Errorf("user_id = %q, se esperaba 1234567", got.UserID)
	}
}

func TestParseTokenResponse_SinAccessToken(t *testing.T) {
	raw := []byte(`{"access_token":"","user_id":1234567}`)

	if _, err := parseTokenResponse(raw); err == nil {
		t.Fatal("se esperaba error cuando falta access_token")
	}
}

func TestParseTokenResponse_SinUserID(t *testing.T) {
	raw := []byte(`{"access_token":"abc123"}`)

	if _, err := parseTokenResponse(raw); err == nil {
		t.Fatal("se esperaba error cuando falta user_id: sin el no hay store_id y la integracion queda rota")
	}
}

func TestParseTokenResponse_JSONInvalido(t *testing.T) {
	if _, err := parseTokenResponse([]byte(`no soy json`)); err == nil {
		t.Fatal("se esperaba error de decodificacion")
	}
}
