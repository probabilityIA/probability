package client

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

func (c *shopifyClient) ValidateToken(ctx context.Context, storeName, accessToken string) (bool, map[string]interface{}, error) {
	if !strings.HasSuffix(storeName, ".myshopify.com") {
		storeName = storeName + ".myshopify.com"
	}

	url := fmt.Sprintf("https://%s/admin/api/2024-10/shop.json", storeName)

	resp, err := c.client.R().
		SetContext(ctx).
		SetHeader("X-Shopify-Access-Token", accessToken).
		SetHeader("Content-Type", "application/json").
		SetResult(map[string]interface{}{}).
		Get(url)

	if err != nil {
		return false, nil, err
	}

	if resp.StatusCode() == http.StatusOK {
		result, ok := resp.Result().(*map[string]interface{})
		if ok && result != nil {
			return true, *result, nil
		}
		return true, nil, nil
	}

	// Mensajes de error descriptivos según el código de estado
	switch resp.StatusCode() {
	case http.StatusUnauthorized: // 401
		return false, nil, fmt.Errorf("token de acceso inválido o expirado. Verifica que el Access Token sea correcto")
	case http.StatusForbidden: // 403
		return false, nil, fmt.Errorf("acceso denegado. Verifica que la app tenga los permisos necesarios (read_orders, read_products)")
	case http.StatusNotFound: // 404
		return false, nil, fmt.Errorf("tienda no encontrada. Verifica que el nombre de la tienda '%s' sea correcto (ej: mitienda.myshopify.com)", storeName)
	case http.StatusTooManyRequests: // 429
		return false, nil, fmt.Errorf("demasiadas solicitudes a Shopify. Intenta nuevamente en unos minutos")
	default:
		return false, nil, fmt.Errorf("error de conexión con Shopify (código %d). Verifica tus credenciales e intenta nuevamente", resp.StatusCode())
	}
}
