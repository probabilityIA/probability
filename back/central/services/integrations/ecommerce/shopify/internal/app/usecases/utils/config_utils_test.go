package utils

import (
	"context"
	"errors"
	"testing"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/shopify/internal/domain"
)

// ──────────────────────────────────────────────────────────────────────────────
// NormalizeConfig
// ──────────────────────────────────────────────────────────────────────────────

func TestNormalizeConfig_ValidMap(t *testing.T) {
	// Arrange
	config := map[string]interface{}{
		"store_name": "mi-tienda.myshopify.com",
	}

	// Act
	result, err := NormalizeConfig(config, "nombre-integracion")

	// Assert
	if err != nil {
		t.Fatalf("se esperaba nil, se obtuvo: %v", err)
	}
	if result["store_name"] != "mi-tienda.myshopify.com" {
		t.Errorf("store_name incorrecto: got %v", result["store_name"])
	}
}

func TestNormalizeConfig_NilConfig_FallbackToName(t *testing.T) {
	// Arrange: config nil, pero integration name disponible

	// Act
	result, err := NormalizeConfig(nil, "fallback-store")

	// Assert
	if err != nil {
		t.Fatalf("se esperaba nil con fallback, se obtuvo: %v", err)
	}
	if result["store_name"] != "fallback-store" {
		t.Errorf("store_name fallback incorrecto: got %v, want %q", result["store_name"], "fallback-store")
	}
}

func TestNormalizeConfig_NilConfig_EmptyName_ReturnsError(t *testing.T) {
	// Arrange: config nil y name vacio

	// Act
	result, err := NormalizeConfig(nil, "")

	// Assert
	if err == nil {
		t.Fatal("se esperaba error cuando config es nil y Name esta vacio, se obtuvo nil")
	}
	if result != nil {
		t.Errorf("se esperaba nil como resultado, se obtuvo: %v", result)
	}
}

func TestNormalizeConfig_InvalidType_FallbackToName(t *testing.T) {
	// Arrange: config de tipo incorrecto (no es map[string]interface{})
	invalidConfig := "not a map"

	// Act
	result, err := NormalizeConfig(invalidConfig, "fallback-store")

	// Assert
	if err != nil {
		t.Fatalf("se esperaba nil con fallback, se obtuvo: %v", err)
	}
	if result["store_name"] != "fallback-store" {
		t.Errorf("store_name fallback incorrecto: got %v", result["store_name"])
	}
}

// ──────────────────────────────────────────────────────────────────────────────
// ExtractStoreName
// ──────────────────────────────────────────────────────────────────────────────

func TestExtractStoreName_FromStoreURL_HTTPS(t *testing.T) {
	// Arrange
	config := map[string]interface{}{
		"store_url": "https://mi-tienda.myshopify.com/",
	}

	// Act
	name, err := ExtractStoreName(config, "")

	// Assert
	if err != nil {
		t.Fatalf("se esperaba nil, se obtuvo: %v", err)
	}
	if name != "mi-tienda.myshopify.com" {
		t.Errorf("store name incorrecto: got %q, want %q", name, "mi-tienda.myshopify.com")
	}
}

func TestExtractStoreName_FromStoreURL_HTTP(t *testing.T) {
	// Arrange
	config := map[string]interface{}{
		"store_url": "http://mi-tienda.myshopify.com",
	}

	// Act
	name, err := ExtractStoreName(config, "")

	// Assert
	if err != nil {
		t.Fatalf("se esperaba nil, se obtuvo: %v", err)
	}
	if name != "mi-tienda.myshopify.com" {
		t.Errorf("store name incorrecto: got %q, want %q", name, "mi-tienda.myshopify.com")
	}
}

func TestExtractStoreName_FromStoreName(t *testing.T) {
	// Arrange: no hay store_url, solo store_name
	config := map[string]interface{}{
		"store_name": "mi-tienda.myshopify.com",
	}

	// Act
	name, err := ExtractStoreName(config, "")

	// Assert
	if err != nil {
		t.Fatalf("se esperaba nil, se obtuvo: %v", err)
	}
	if name != "mi-tienda.myshopify.com" {
		t.Errorf("store name incorrecto: got %q", name)
	}
}

func TestExtractStoreName_FallbackToIntegrationName(t *testing.T) {
	// Arrange: config sin store_url ni store_name
	config := map[string]interface{}{}

	// Act
	name, err := ExtractStoreName(config, "integracion-nombre")

	// Assert
	if err != nil {
		t.Fatalf("se esperaba nil con fallback, se obtuvo: %v", err)
	}
	if name != "integracion-nombre" {
		t.Errorf("store name fallback incorrecto: got %q, want %q", name, "integracion-nombre")
	}
}

func TestExtractStoreName_AllEmpty_ReturnsError(t *testing.T) {
	// Arrange: sin store_url, sin store_name, sin integration name
	config := map[string]interface{}{}

	// Act
	name, err := ExtractStoreName(config, "")

	// Assert
	if err == nil {
		t.Fatal("se esperaba error cuando no hay ninguna fuente de store_name, se obtuvo nil")
	}
	if name != "" {
		t.Errorf("se esperaba string vacio, se obtuvo %q", name)
	}
}

func TestExtractStoreName_StoreURL_RemovesTrailingSlash(t *testing.T) {
	// Arrange
	config := map[string]interface{}{
		"store_url": "https://mi-tienda.myshopify.com/",
	}

	// Act
	name, err := ExtractStoreName(config, "")

	// Assert
	if err != nil {
		t.Fatalf("se esperaba nil, se obtuvo: %v", err)
	}
	// No debe tener trailing slash
	if len(name) > 0 && name[len(name)-1] == '/' {
		t.Errorf("el store name no debe tener trailing slash: %q", name)
	}
}

// ──────────────────────────────────────────────────────────────────────────────
// GetAccessToken
// ──────────────────────────────────────────────────────────────────────────────

func TestGetAccessToken_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	expectedToken := "shpat_abc123xyz"
	integrationID := "integration-5"

	mockSvc := &mockIntegrationSvc{
		DecryptCredentialFn: func(ctx context.Context, id string, fieldName string) (string, error) {
			if id != integrationID {
				t.Errorf("integrationID incorrecto: got %q, want %q", id, integrationID)
			}
			if fieldName != "access_token" {
				t.Errorf("fieldName incorrecto: got %q, want %q", fieldName, "access_token")
			}
			return expectedToken, nil
		},
	}

	// Act
	token, err := GetAccessToken(ctx, mockSvc, integrationID)

	// Assert
	if err != nil {
		t.Fatalf("se esperaba nil, se obtuvo: %v", err)
	}
	if token != expectedToken {
		t.Errorf("token incorrecto: got %q, want %q", token, expectedToken)
	}
}

func TestGetAccessToken_DecryptError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	decryptErr := errors.New("key not found")

	mockSvc := &mockIntegrationSvc{
		DecryptCredentialFn: func(ctx context.Context, id string, fieldName string) (string, error) {
			return "", decryptErr
		},
	}

	// Act
	token, err := GetAccessToken(ctx, mockSvc, "integration-1")

	// Assert
	if err == nil {
		t.Fatal("se esperaba error, se obtuvo nil")
	}
	if !errors.Is(err, decryptErr) {
		t.Errorf("error incorrecto: got %v, want %v", err, decryptErr)
	}
	if token != "" {
		t.Errorf("se esperaba token vacio, se obtuvo %q", token)
	}
}

// ──────────────────────────────────────────────────────────────────────────────
// Mock local para domain.IIntegrationService (solo lo que usa utils)
// ──────────────────────────────────────────────────────────────────────────────

type mockIntegrationSvc struct {
	DecryptCredentialFn func(ctx context.Context, integrationID string, fieldName string) (string, error)
}

func (m *mockIntegrationSvc) GetIntegrationByID(ctx context.Context, integrationID string) (*domain.Integration, error) {
	return nil, nil
}

func (m *mockIntegrationSvc) GetIntegrationByExternalID(ctx context.Context, externalID string, integrationType int) (*domain.Integration, error) {
	return nil, nil
}

func (m *mockIntegrationSvc) DecryptCredential(ctx context.Context, integrationID string, fieldName string) (string, error) {
	if m.DecryptCredentialFn != nil {
		return m.DecryptCredentialFn(ctx, integrationID, fieldName)
	}
	return "", nil
}

func (m *mockIntegrationSvc) UpdateIntegrationConfig(ctx context.Context, integrationID string, config map[string]interface{}) error {
	return nil
}
