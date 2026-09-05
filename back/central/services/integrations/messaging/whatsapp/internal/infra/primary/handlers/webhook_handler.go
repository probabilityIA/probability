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
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/infra/primary/handlers/request"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/infra/primary/queue/consumerwebhook"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

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

	if mode != "subscribe" {
		h.log.Warn(ctx).
			Str("mode", mode).
			Msg("[Webhook Handler] - modo de suscripción inválido")
		c.String(http.StatusForbidden, "Modo inválido")
		return
	}

	expectedToken := h.getVerifyToken(ctx)
	if expectedToken == "" {
		h.log.Error(ctx).Msg("[Webhook Handler] - verify_token no encontrado en cache ni en env")
		c.String(http.StatusForbidden, "Token de verificación no configurado")
		return
	}

	if token != expectedToken {
		h.log.Warn(ctx).
			Str("received_token", token).
			Msg("[Webhook Handler] - token de verificación inválido")
		c.String(http.StatusForbidden, "Token inválido")
		return
	}

	h.log.Info(ctx).
		Str("challenge", challenge).
		Msg("[Webhook Handler] - webhook verificado exitosamente")

	c.String(http.StatusOK, challenge)
}

func (h *handler) ReceiveWebhook(c *gin.Context) {
	ctx := c.Request.Context()

	h.log.Info(ctx).Msg("[Webhook Handler] - recibiendo webhook de WhatsApp")

	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		h.log.Error(ctx).Err(err).Msg("[Webhook Handler] - error leyendo body")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request",
			"message": "Error leyendo el body del request",
		})
		return
	}

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

	c.JSON(http.StatusOK, gin.H{
		"status": "received",
	})

	h.enqueueWebhook(bodyBytes, webhook)
}

func (h *handler) enqueueWebhook(bodyBytes []byte, webhook request.WebhookPayload) {
	ctx := context.Background()

	if h.rabbit != nil {
		if err := h.rabbit.Publish(ctx, rabbitmq.QueueWebhooksWhatsappReceived, bodyBytes); err == nil {
			return
		}
		h.log.Warn(ctx).Msg("[Webhook Handler] - no se pudo encolar el webhook, procesando inline como fallback")
	}

	go consumerwebhook.Dispatch(ctx, h.useCase, h.templatesUseCase, h.log, webhook)
}

const whatsAppTypeID = uint(2)

func (h *handler) getPlatformCredField(ctx context.Context, field string) string {
	if h.platformCredsGetter == nil {
		return ""
	}
	creds, err := h.platformCredsGetter.GetCachedPlatformCredentials(ctx, whatsAppTypeID)
	if err != nil {
		return ""
	}
	if val, ok := creds[field].(string); ok {
		return val
	}
	return ""
}

func (h *handler) getVerifyToken(ctx context.Context) string {
	return h.getPlatformCredField(ctx, "verify_token")
}

func (h *handler) getWebhookSecret(ctx context.Context) string {
	return h.getPlatformCredField(ctx, "webhook_secret")
}

func (h *handler) verifySignature(payload []byte, signatureHeader string) bool {
	secret := h.getWebhookSecret(context.Background())
	if secret == "" {
		h.log.Error().Msg("[Webhook Handler] - webhook_secret no encontrado en cache ni en env")
		return false
	}

	signatureParts := strings.Split(signatureHeader, "=")
	if len(signatureParts) != 2 || signatureParts[0] != "sha256" {
		h.log.Warn().
			Str("signature_header", signatureHeader).
			Msg("[Webhook Handler] - formato de firma inválido")
		return false
	}

	expectedSignature := signatureParts[1]

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	calculatedSignature := hex.EncodeToString(mac.Sum(nil))

	valid := hmac.Equal([]byte(calculatedSignature), []byte(expectedSignature))

	if !valid {
		h.log.Warn().
			Str("expected", expectedSignature).
			Str("calculated", calculatedSignature).
			Msg("[Webhook Handler] - firma no coincide")
	}

	return valid
}
