package client

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// DeleteWebhook elimina un webhook de Shopify por su ID
func (c *shopifyClient) DeleteWebhook(ctx context.Context, storeName, accessToken, webhookID string) error {
	// Normalizar el nombre de la tienda (remover .myshopify.com si está presente)
	shop := strings.TrimSuffix(storeName, ".myshopify.com")

	// Validar que webhookID sea un número válido
	_, err := strconv.ParseInt(webhookID, 10, 64)
	if err != nil {
		return fmt.Errorf("webhook ID inválido: %s", webhookID)
	}

	// Construir la URL según el formato de Shopify API
	url := fmt.Sprintf("https://%s.myshopify.com/admin/api/2024-10/webhooks/%s.json", shop, webhookID)

	var errorResponse struct {
		Errors map[string]interface{} `json:"errors"`
	}

	resp, err := c.client.R().
		SetContext(ctx).
		SetHeader("X-Shopify-Access-Token", accessToken).
		SetHeader("Content-Type", "application/json").
		SetError(&errorResponse).
		Delete(url)

	if err != nil {
		return fmt.Errorf("error al realizar request a Shopify: %w", err)
	}

	// Shopify retorna 204 No Content cuando se elimina exitosamente
	if resp.StatusCode() == http.StatusNoContent || resp.StatusCode() == http.StatusOK {
		return nil
	}

	// Si hay errores en la respuesta
	if len(errorResponse.Errors) > 0 {
		return fmt.Errorf("error al eliminar webhook en Shopify: %v", errorResponse.Errors)
	}

	// Mensajes de error descriptivos según el código de estado
	switch resp.StatusCode() {
	case http.StatusUnauthorized: // 401
		return fmt.Errorf("token de acceso inválido o expirado")
	case http.StatusForbidden: // 403
		return fmt.Errorf("acceso denegado. Verifica que la app tenga permisos para eliminar webhooks")
	case http.StatusNotFound: // 404
		return fmt.Errorf("webhook no encontrado con ID: %s", webhookID)
	case http.StatusTooManyRequests: // 429
		return fmt.Errorf("demasiadas solicitudes a Shopify. Intenta nuevamente en unos minutos")
	default:
		return fmt.Errorf("error al eliminar webhook en Shopify (código %d): %s", resp.StatusCode(), resp.String())
	}
}













