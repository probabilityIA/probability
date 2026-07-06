package client

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

type CreateWebhookResponse struct {
	Webhook struct {
		ID        int64  `json:"id"`
		Address   string `json:"address"`
		Topic     string `json:"topic"`
		Format    string `json:"format"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	} `json:"webhook"`
}

func (c *shopifyClient) CreateWebhook(ctx context.Context, storeName, accessToken, webhookURL, event string) (string, error) {
	url := buildURL(storeName, "/admin/api/2024-10/webhooks.json")

	payload := map[string]interface{}{
		"webhook": map[string]interface{}{
			"topic":   event,
			"address": webhookURL,
			"format":  "json",
		},
	}

	var result CreateWebhookResponse

	resp, err := c.client.R().
		SetContext(ctx).
		SetHeader("X-Shopify-Access-Token", accessToken).
		SetHeader("Content-Type", "application/json").
		SetBody(payload).
		SetResult(&result).
		Post(url)

	if err != nil {
		return "", fmt.Errorf("error al realizar request a Shopify: %w", err)
	}

	if resp.StatusCode() == http.StatusCreated || resp.StatusCode() == http.StatusOK {
		webhookID := fmt.Sprintf("%d", result.Webhook.ID)
		return webhookID, nil
	}

	body := strings.TrimSpace(resp.String())

	switch resp.StatusCode() {
	case http.StatusUnauthorized:
		return "", fmt.Errorf("token de acceso inválido o expirado: %s", body)
	case http.StatusForbidden:
		return "", fmt.Errorf("acceso denegado al crear el webhook: %s", body)
	case http.StatusNotFound:
		return "", fmt.Errorf("tienda no encontrada: %s", storeName)
	case http.StatusConflict:
		return "", fmt.Errorf("el webhook ya existe para este evento y URL: %s", body)
	case http.StatusUnprocessableEntity:
		return "", fmt.Errorf("datos inválidos para crear el webhook: %s", body)
	case http.StatusTooManyRequests:
		return "", fmt.Errorf("demasiadas solicitudes a Shopify. Intenta nuevamente en unos minutos")
	default:
		return "", fmt.Errorf("error al crear webhook en Shopify (código %d): %s", resp.StatusCode(), body)
	}
}
