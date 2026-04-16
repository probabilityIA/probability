package usecases

import (
	"context"
	"fmt"
)

// TestConnection valida que las credenciales sean correctas contra la API de Shopify.
func (uc *SyncOrdersUseCase) TestConnection(ctx context.Context, config map[string]interface{}, credentials map[string]interface{}) error {
	storeName, ok := config["store_name"].(string)
	if !ok || storeName == "" {
		return fmt.Errorf("el nombre de la tienda (store_name) es requerido")
	}

	accessToken, ok := credentials["access_token"].(string)
	if !ok || accessToken == "" {
		return fmt.Errorf("el token de acceso (access_token) es requerido")
	}

	valid, _, err := uc.shopifyClient.ValidateToken(ctx, storeName, accessToken)
	if err != nil {
		return err
	}

	if !valid {
		return fmt.Errorf("credenciales o nombre de tienda inválidos")
	}

	return nil
}
