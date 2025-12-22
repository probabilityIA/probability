package client

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/secamc93/probability/back/central/services/integrations/shopify/internal/domain"
)

// ListWebhooksResponse representa la respuesta de Shopify al listar webhooks
type ListWebhooksResponse struct {
	Webhooks []struct {
		ID        int64  `json:"id"`
		Address   string `json:"address"`
		Topic     string `json:"topic"`
		Format    string `json:"format"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	} `json:"webhooks"`
}

// ListWebhooks lista todos los webhooks de la tienda de Shopify
func (c *shopifyClient) ListWebhooks(ctx context.Context, storeName, accessToken string) ([]domain.WebhookInfo, error) {
	// Normalizar el nombre de la tienda (remover .myshopify.com si está presente)
	shop := strings.TrimSuffix(storeName, ".myshopify.com")

	// Construir la URL según el formato de Shopify API
	url := fmt.Sprintf("https://%s.myshopify.com/admin/api/2024-10/webhooks.json", shop)

	var result ListWebhooksResponse
	var errorResponse struct {
		Errors map[string]interface{} `json:"errors"`
	}

	resp, err := c.client.R().
		SetContext(ctx).
		SetHeader("X-Shopify-Access-Token", accessToken).
		SetHeader("Content-Type", "application/json").
		SetResult(&result).
		SetError(&errorResponse).
		Get(url)

	if err != nil {
		return nil, fmt.Errorf("error al realizar request a Shopify: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		// Si hay errores en la respuesta
		if len(errorResponse.Errors) > 0 {
			return nil, fmt.Errorf("error al listar webhooks en Shopify: %v", errorResponse.Errors)
		}

		// Mensajes de error descriptivos según el código de estado
		switch resp.StatusCode() {
		case http.StatusUnauthorized: // 401
			return nil, fmt.Errorf("token de acceso inválido o expirado")
		case http.StatusForbidden: // 403
			return nil, fmt.Errorf("acceso denegado. Verifica que la app tenga permisos para listar webhooks")
		case http.StatusNotFound: // 404
			return nil, fmt.Errorf("tienda no encontrada: %s", storeName)
		case http.StatusTooManyRequests: // 429
			return nil, fmt.Errorf("demasiadas solicitudes a Shopify. Intenta nuevamente en unos minutos")
		default:
			return nil, fmt.Errorf("error al listar webhooks en Shopify (código %d): %s", resp.StatusCode(), resp.String())
		}
	}

	// Convertir a domain.WebhookInfo
	webhooks := make([]domain.WebhookInfo, len(result.Webhooks))
	for i, wh := range result.Webhooks {
		webhooks[i] = domain.WebhookInfo{
			ID:        strconv.FormatInt(wh.ID, 10),
			Address:   wh.Address,
			Topic:     wh.Topic,
			Format:    wh.Format,
			CreatedAt: wh.CreatedAt,
			UpdatedAt: wh.UpdatedAt,
		}
	}

	return webhooks, nil
}








