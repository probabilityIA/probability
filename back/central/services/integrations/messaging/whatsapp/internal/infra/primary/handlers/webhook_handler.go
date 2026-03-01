package handlers

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/infra/primary/handlers/mappers"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/infra/primary/handlers/request"
)

// VerifyWebhook maneja la verificación inicial del webhook (GET request)
// @Summary Verifica el webhook de WhatsApp
// @Description Endpoint para verificación del webhook por Meta. Retorna el challenge si el token es válido.
// @Tags WhatsApp Webhooks
// @Accept json
// @Produce plain
// @Param hub.mode query string true "Modo de suscripción (debe ser 'subscribe')"
// @Param hub.verify_token query string true "Token de verificación"
// @Param hub.challenge query string true "Challenge a retornar"
// @Success 200 {string} string "Challenge token"
// @Failure 403 {string} string "Forbidden"
// @Router /integrations/whatsapp/webhook [get]
func (h *handler) VerifyWebhook(c *gin.Context) {
	ctx := c.Request.Context()

	mode := c.Query("hub.mode")
	token := c.Query("hub.verify_token")
	challenge := c.Query("hub.challenge")

	h.log.Info(ctx).
		Str("mode", mode).
		Str("token", token).
		Str("challenge", challenge).
		Msg("[Webhook Handler] - solicitud de verificación de webhook")

	// Verificar que el modo sea "subscribe"
	if mode != "subscribe" {
		h.log.Warn(ctx).
			Str("mode", mode).
			Msg("[Webhook Handler] - modo de suscripción inválido")
		c.String(http.StatusForbidden, "Modo inválido")
		return
	}

	// Obtener token de verificación de la configuración
	expectedToken := h.config.Get("WHATSAPP_VERIFY_TOKEN")
	if expectedToken == "" {
		h.log.Error(ctx).Msg("[Webhook Handler] - WHATSAPP_VERIFY_TOKEN no configurado")
		c.String(http.StatusForbidden, "Token de verificación no configurado")
		return
	}

	// Verificar que el token coincida
	if token != expectedToken {
		h.log.Warn(ctx).
			Str("received_token", token).
			Msg("[Webhook Handler] - token de verificación inválido")
		c.String(http.StatusForbidden, "Token inválido")
		return
	}

	// Retornar el challenge para completar la verificación
	h.log.Info(ctx).
		Str("challenge", challenge).
		Msg("[Webhook Handler] - webhook verificado exitosamente")

	c.String(http.StatusOK, challenge)
}

// ReceiveWebhook maneja los eventos entrantes del webhook (POST request)
// @Summary Recibe eventos de WhatsApp
// @Description Endpoint para recibir eventos de mensajes y estados desde Meta WhatsApp Business API
// @Tags WhatsApp Webhooks
// @Accept json
// @Produce json
// @Param X-Hub-Signature-256 header string true "Firma HMAC-SHA256 del payload"
// @Param payload body request.WebhookPayload true "Payload del webhook"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /integrations/whatsapp/webhook [post]
func (h *handler) ReceiveWebhook(c *gin.Context) {
	ctx := c.Request.Context()

	h.log.Info(ctx).Msg("[Webhook Handler] - recibiendo webhook de WhatsApp")

	// 1. Leer el body completo para validar firma
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		h.log.Error(ctx).Err(err).Msg("[Webhook Handler] - error leyendo body")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request",
			"message": "Error leyendo el body del request",
		})
		return
	}

	// 2. Validar firma HMAC-SHA256
	signature := c.GetHeader("X-Hub-Signature-256")
	if signature == "" {
		h.log.Warn(ctx).Msg("[Webhook Handler] - falta header X-Hub-Signature-256")
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "missing_signature",
			"message": "Falta la firma del webhook",
		})
		return
	}

	if !h.verifySignature(bodyBytes, signature) {
		h.log.Error(ctx).
			Str("signature", signature).
			Msg("[Webhook Handler] - firma inválida")
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "invalid_signature",
			"message": "La firma del webhook es inválida",
		})
		return
	}

	// 3. Parsear payload
	var webhook request.WebhookPayload
	if err := json.Unmarshal(bodyBytes, &webhook); err != nil {
		h.log.Error(ctx).Err(err).Msg("[Webhook Handler] - error parseando webhook")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_payload",
			"message": "El payload del webhook es inválido",
			"details": err.Error(),
		})
		return
	}

	h.log.Info(ctx).
		Str("object", webhook.Object).
		Int("entries", len(webhook.Entry)).
		Msg("[Webhook Handler] - webhook parseado correctamente")

	// 4. Retornar 200 inmediatamente (requisito de Meta: responder en <5s)
	c.JSON(http.StatusOK, gin.H{
		"status": "received",
	})

	// 5. Procesar en background (goroutine)
	go h.processWebhookAsync(webhook)
}

// processWebhookAsync procesa el webhook de forma asíncrona
func (h *handler) processWebhookAsync(webhook request.WebhookPayload) {
	// Crear nuevo contexto para operación asíncrona
	ctx := context.Background()

	h.log.Info(ctx).
		Str("object", webhook.Object).
		Int("entries", len(webhook.Entry)).
		Msg("[Webhook Handler] - procesando webhook de forma asíncrona")

	// Mapear de infra → domain ANTES de invocar use case
	webhookDTO := mappers.WebhookPayloadToDomain(webhook)

	// Determinar tipo de evento y procesar
	for _, entry := range webhookDTO.Entry {
		for _, change := range entry.Changes {
			switch change.Field {
			case "messages":
				// Procesar mensajes entrantes o cambios de estado
				if len(change.Value.Messages) > 0 {
					if err := h.useCase.HandleIncomingMessage(ctx, webhookDTO); err != nil {
						h.log.Error(ctx).Err(err).Msg("[Webhook Handler] - error procesando mensajes entrantes")
					}
				}
				if len(change.Value.Statuses) > 0 {
					if err := h.useCase.HandleMessageStatus(ctx, webhookDTO); err != nil {
						h.log.Error(ctx).Err(err).Msg("[Webhook Handler] - error procesando estados de mensajes")
					}
				}
			case "message_template_status_update":
				h.log.Info(ctx).Msg("[Webhook Handler] - actualización de estado de plantilla recibida")
				// TODO: Implementar si se necesita tracking de estado de plantillas
			default:
				h.log.Warn(ctx).
					Str("field", change.Field).
					Msg("[Webhook Handler] - campo de webhook no reconocido")
			}
		}
	}

	h.log.Info(ctx).Msg("[Webhook Handler] - procesamiento asíncrono completado")
}

// verifySignature verifica la firma HMAC-SHA256 del webhook
func (h *handler) verifySignature(payload []byte, signatureHeader string) bool {
	// Obtener secret de configuración
	secret := h.config.Get("WHATSAPP_WEBHOOK_SECRET")
	if secret == "" {
		h.log.Error().Msg("[Webhook Handler] - WHATSAPP_WEBHOOK_SECRET no configurado")
		return false
	}

	// La firma viene en formato "sha256=<hex>"
	signatureParts := strings.Split(signatureHeader, "=")
	if len(signatureParts) != 2 || signatureParts[0] != "sha256" {
		h.log.Warn().
			Str("signature_header", signatureHeader).
			Msg("[Webhook Handler] - formato de firma inválido")
		return false
	}

	expectedSignature := signatureParts[1]

	// Calcular HMAC-SHA256
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	calculatedSignature := hex.EncodeToString(mac.Sum(nil))

	// Comparar firmas
	valid := hmac.Equal([]byte(calculatedSignature), []byte(expectedSignature))

	if !valid {
		h.log.Warn().
			Str("expected", expectedSignature).
			Str("calculated", calculatedSignature).
			Msg("[Webhook Handler] - firma no coincide")
	}

	return valid
}
