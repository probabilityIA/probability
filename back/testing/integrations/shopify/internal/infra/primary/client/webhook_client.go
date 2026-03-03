package client

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/secamc93/probability/back/testing/integrations/shopify/internal/domain"
	"github.com/secamc93/probability/back/testing/shared/env"
	"github.com/secamc93/probability/back/testing/shared/log"
	sharedtypes "github.com/secamc93/probability/back/testing/shared/types"
)

// Verificar que WebhookClient implementa domain.IWebhookClient
var _ domain.IWebhookClient = (*WebhookClient)(nil)

// WebhookClient maneja el envío de webhooks a la URL configurada
type WebhookClient struct {
	baseURL    string
	config     env.IConfig
	logger     log.ILogger
	httpClient *http.Client
}

// New crea una nueva instancia del cliente de webhook
func New(config env.IConfig, logger log.ILogger) *WebhookClient {
	return &WebhookClient{
		baseURL: config.Get("WEBHOOK_BASE_URL"),
		config:  config,
		logger:  logger,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SendWebhook envía un webhook a la URL configurada
func (c *WebhookClient) SendWebhook(topic string, shopDomain string, payload interface{}) error {
	url := fmt.Sprintf("%s/api/v1/integrations/shopify/webhook", c.baseURL)

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error al serializar payload: %w", err)
	}

	// Generar HMAC (simulado para pruebas - en producción usar la clave real)
	hmacValue := c.generateHMAC(payloadBytes)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return fmt.Errorf("error al crear request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Shopify-Topic", topic)
	req.Header.Set("X-Shopify-Shop-Domain", shopDomain)
	req.Header.Set("X-Shopify-Hmac-Sha256", hmacValue)
	// API version from config (matches integration config in database)
	req.Header.Set("X-Shopify-API-Version", c.config.GetWithDefault("SHOPIFY_API_VERSION", "2024-01"))
	req.Header.Set("X-Shopify-Webhook-Id", fmt.Sprintf("test-%d", time.Now().Unix()))

	c.logger.Info().
		Str("url", url).
		Str("topic", topic).
		Str("shop_domain", shopDomain).
		Msg("Enviando webhook")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error al enviar webhook: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("webhook recibió status code %d", resp.StatusCode)
	}

	c.logger.Info().
		Int("status_code", resp.StatusCode).
		Str("topic", topic).
		Msg("Webhook enviado exitosamente")

	return nil
}

// BuildWebhook builds the webhook payload without sending it.
// baseURL is the central API URL (e.g. http://localhost:3050) — NOT the WEBHOOK_BASE_URL env var.
func (c *WebhookClient) BuildWebhook(topic string, shopDomain string, payload interface{}, baseURL string) (*sharedtypes.WebhookPayload, error) {
	url := fmt.Sprintf("%s/api/v1/integrations/shopify/webhook", baseURL)

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error al serializar payload: %w", err)
	}

	hmacValue := c.generateHMAC(payloadBytes)

	headers := map[string]string{
		"Content-Type":            "application/json",
		"X-Shopify-Topic":         topic,
		"X-Shopify-Shop-Domain":   shopDomain,
		"X-Shopify-Hmac-Sha256":   hmacValue,
		"X-Shopify-API-Version":   c.config.GetWithDefault("SHOPIFY_API_VERSION", "2024-01"),
		"X-Shopify-Webhook-Id":    fmt.Sprintf("test-%d", time.Now().Unix()),
	}

	// Unmarshal back to map so the body is a JSON object in the response
	var bodyMap map[string]interface{}
	if err := json.Unmarshal(payloadBytes, &bodyMap); err != nil {
		return nil, fmt.Errorf("error al parsear payload: %w", err)
	}

	return &sharedtypes.WebhookPayload{
		URL:     url,
		Method:  "POST",
		Headers: headers,
		Body:    bodyMap,
		RawBody: string(payloadBytes), // exact bytes used for HMAC — frontend must send this, not re-serialize Body
	}, nil
}

// generateHMAC genera un HMAC usando el client_secret real de Shopify
func (c *WebhookClient) generateHMAC(payload []byte) string {
	// Leer el client_secret desde variable de entorno
	// Este debe coincidir con el client_secret guardado en la BD (integración ID=1)
	secret := c.config.GetWithDefault("SHOPIFY_CLIENT_SECRET", "test-secret-for-development")

	if secret == "test-secret-for-development" {
		c.logger.Warn().Msg("⚠️ Usando SHOPIFY_CLIENT_SECRET por defecto - configura la variable de entorno para testing real")
	}

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}













