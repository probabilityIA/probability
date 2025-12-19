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

	"github.com/secamc93/probability/back/integrationTest/integrations/shopify/internal/domain"
	"github.com/secamc93/probability/back/integrationTest/shared/env"
	"github.com/secamc93/probability/back/integrationTest/shared/log"
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
	req.Header.Set("X-Shopify-API-Version", c.config.GetWithDefault("SHOPIFY_API_VERSION", "2024-10"))
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

// generateHMAC genera un HMAC simulado para pruebas
func (c *WebhookClient) generateHMAC(payload []byte) string {
	// En pruebas, generamos un HMAC simulado
	// En producción, esto debería usar la clave secreta real de Shopify
	secret := "test_secret_key_for_integration_tests"
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}



