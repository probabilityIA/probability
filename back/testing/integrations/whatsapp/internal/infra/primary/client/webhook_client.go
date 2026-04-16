package client

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/secamc93/probability/back/testing/integrations/whatsapp/internal/domain"
	"github.com/secamc93/probability/back/testing/shared/env"
	"github.com/secamc93/probability/back/testing/shared/log"
)

// Verificar que WebhookClient implementa domain.IWebhookClient
var _ domain.IWebhookClient = (*WebhookClient)(nil)

// WebhookClient maneja el envío de webhooks de WhatsApp
type WebhookClient struct {
	baseURL       string
	webhookSecret string
	config        env.IConfig
	logger        log.ILogger
	httpClient    *http.Client
}

// New crea una nueva instancia del cliente de webhook de WhatsApp
func New(config env.IConfig, logger log.ILogger) *WebhookClient {
	return &WebhookClient{
		baseURL:       config.Get("WEBHOOK_BASE_URL"),
		webhookSecret: config.GetWithDefault("WHATSAPP_WEBHOOK_SECRET", "test_webhook_secret"),
		config:        config,
		logger:        logger,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SendWebhook envía un webhook de WhatsApp al sistema real
func (c *WebhookClient) SendWebhook(payload domain.WebhookPayload) error {
	url := fmt.Sprintf("%s/api/integrations/whatsapp/webhook", c.baseURL)

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error al serializar payload: %w", err)
	}

	// Generar HMAC-SHA256 signature (igual que Meta)
	signature := c.generateSignature(payloadBytes)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return fmt.Errorf("error al crear request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Hub-Signature-256", "sha256="+signature)

	c.logger.Info().
		Str("url", url).
		Msg("Enviando webhook de WhatsApp")

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
		Msg("Webhook de WhatsApp enviado exitosamente")

	return nil
}

// generateSignature genera la firma HMAC-SHA256 (igual que Meta)
func (c *WebhookClient) generateSignature(payload []byte) string {
	h := hmac.New(sha256.New, []byte(c.webhookSecret))
	h.Write(payload)
	return hex.EncodeToString(h.Sum(nil))
}
