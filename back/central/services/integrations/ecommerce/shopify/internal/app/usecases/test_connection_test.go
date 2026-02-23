package usecases

import (
	"context"
	"errors"
	"testing"
)

func TestTestConnection_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()

	shopifyClient := &mockShopifyClient{
		ValidateTokenFn: func(ctx context.Context, storeName, accessToken string) (bool, map[string]interface{}, error) {
			if storeName != "mi-tienda.myshopify.com" {
				t.Errorf("storeName inesperado: got %q", storeName)
			}
			if accessToken != "shpat_abc123" {
				t.Errorf("accessToken inesperado: got %q", accessToken)
			}
			return true, map[string]interface{}{"shop": storeName}, nil
		},
	}

	uc := newTestUseCase(&mockIntegrationService{}, shopifyClient, &mockOrderPublisher{})

	config := map[string]interface{}{
		"store_name": "mi-tienda.myshopify.com",
	}
	credentials := map[string]interface{}{
		"access_token": "shpat_abc123",
	}

	// Act
	err := uc.TestConnection(ctx, config, credentials)

	// Assert
	if err != nil {
		t.Fatalf("se esperaba nil, se obtuvo: %v", err)
	}
}

func TestTestConnection_MissingStoreName(t *testing.T) {
	// Arrange
	ctx := context.Background()
	uc := newTestUseCase(&mockIntegrationService{}, &mockShopifyClient{}, &mockOrderPublisher{})

	config := map[string]interface{}{} // sin store_name
	credentials := map[string]interface{}{
		"access_token": "shpat_abc123",
	}

	// Act
	err := uc.TestConnection(ctx, config, credentials)

	// Assert
	if err == nil {
		t.Fatal("se esperaba error por store_name faltante, se obtuvo nil")
	}
}

func TestTestConnection_EmptyStoreName(t *testing.T) {
	// Arrange
	ctx := context.Background()
	uc := newTestUseCase(&mockIntegrationService{}, &mockShopifyClient{}, &mockOrderPublisher{})

	config := map[string]interface{}{
		"store_name": "",
	}
	credentials := map[string]interface{}{
		"access_token": "shpat_abc123",
	}

	// Act
	err := uc.TestConnection(ctx, config, credentials)

	// Assert
	if err == nil {
		t.Fatal("se esperaba error por store_name vacio, se obtuvo nil")
	}
}

func TestTestConnection_MissingAccessToken(t *testing.T) {
	// Arrange
	ctx := context.Background()
	uc := newTestUseCase(&mockIntegrationService{}, &mockShopifyClient{}, &mockOrderPublisher{})

	config := map[string]interface{}{
		"store_name": "mi-tienda.myshopify.com",
	}
	credentials := map[string]interface{}{} // sin access_token

	// Act
	err := uc.TestConnection(ctx, config, credentials)

	// Assert
	if err == nil {
		t.Fatal("se esperaba error por access_token faltante, se obtuvo nil")
	}
}

func TestTestConnection_InvalidCredentials(t *testing.T) {
	// Arrange
	ctx := context.Background()

	shopifyClient := &mockShopifyClient{
		ValidateTokenFn: func(ctx context.Context, storeName, accessToken string) (bool, map[string]interface{}, error) {
			// Shopify responde con valid=false (credenciales incorrectas)
			return false, nil, nil
		},
	}

	uc := newTestUseCase(&mockIntegrationService{}, shopifyClient, &mockOrderPublisher{})

	config := map[string]interface{}{
		"store_name": "tienda-invalida.myshopify.com",
	}
	credentials := map[string]interface{}{
		"access_token": "token-incorrecto",
	}

	// Act
	err := uc.TestConnection(ctx, config, credentials)

	// Assert
	if err == nil {
		t.Fatal("se esperaba error por credenciales invalidas, se obtuvo nil")
	}
}

func TestTestConnection_ClientError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	clientErr := errors.New("connection refused")

	shopifyClient := &mockShopifyClient{
		ValidateTokenFn: func(ctx context.Context, storeName, accessToken string) (bool, map[string]interface{}, error) {
			return false, nil, clientErr
		},
	}

	uc := newTestUseCase(&mockIntegrationService{}, shopifyClient, &mockOrderPublisher{})

	config := map[string]interface{}{
		"store_name": "mi-tienda.myshopify.com",
	}
	credentials := map[string]interface{}{
		"access_token": "shpat_abc123",
	}

	// Act
	err := uc.TestConnection(ctx, config, credentials)

	// Assert
	if err == nil {
		t.Fatal("se esperaba error del cliente Shopify, se obtuvo nil")
	}
	if !errors.Is(err, clientErr) {
		t.Errorf("error incorrecto: got %v, want %v", err, clientErr)
	}
}
