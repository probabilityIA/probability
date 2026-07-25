package handlers

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/shopify/internal/infra/primary/handlers/request"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/shopify/internal/infra/primary/handlers/response"
	shopifyqueue "github.com/secamc93/probability/back/central/services/integrations/ecommerce/shopify/internal/infra/primary/queue"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

func (h *ShopifyHandler) WebhookHandler(c *gin.Context) {
	var headers request.WebhookHeaders

	if err := c.ShouldBindHeader(&headers); err != nil {
		h.logger.Error().Err(err).Msg("Error al validar headers del webhook")
		c.JSON(http.StatusBadRequest, response.WebhookResponse{
			Success: false,
			Message: "Headers requeridos faltantes o inválidos",
		})
		return
	}

	h.logger.Info().
		Str("topic", headers.Topic).
		Str("shop_domain", headers.ShopDomain).
		Str("hmac", headers.Hmac).
		Msg("Webhook recibido de Shopify")

	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		h.logger.Error().Err(err).Msg("Error al leer el body del webhook")
		c.JSON(http.StatusBadRequest, response.WebhookResponse{
			Success: false,
			Message: "Error al leer el body",
		})
		return
	}

	payloadPreview := string(bodyBytes)
	if len(payloadPreview) > 500 {
		payloadPreview = payloadPreview[:500] + "..."
	}
	h.logger.Info().
		Str("topic", headers.Topic).
		Str("shop_domain", headers.ShopDomain).
		Int("payload_size", len(bodyBytes)).
		Str("payload_preview", payloadPreview).
		Msg("📦 Payload del webhook")

	shopifySecret, err := h.useCase.GetClientSecretByShopDomain(c.Request.Context(), headers.ShopDomain)
	if err != nil {
		h.logger.Warn().
			Err(err).
			Str("shop_domain", headers.ShopDomain).
			Msg("⚠️ No se pudo recuperar el secreto específico de la tienda, intentando fallback global")

		shopifySecret = h.config.Get("SHOPIFY_CLIENT_SECRET")
		if shopifySecret == "" {
			shopifySecret = h.config.Get("SHOPIFY_API_SECRET")
		}
	}

	h.logger.Debug().
		Bool("has_secret", shopifySecret != "").
		Str("shop_domain", headers.ShopDomain).
		Msg("🔐 Verificando HMAC dinámico")

	if shopifySecret != "" {
		if !VerifyWebhookHMAC(bodyBytes, headers.Hmac, shopifySecret) {
			h.logger.Error().
				Str("shop_domain", headers.ShopDomain).
				Str("hmac_header", headers.Hmac).
				Msg("❌ Firma HMAC inválida para esta tienda")
			c.JSON(http.StatusUnauthorized, response.WebhookResponse{
				Success: false,
				Message: "Firma HMAC inválida",
			})
			return
		}
		h.logger.Info().Str("shop_domain", headers.ShopDomain).Msg("✅ HMAC válido")
	} else {
		h.logger.Error().Str("shop_domain", headers.ShopDomain).Msg("❌ No hay secreto configurado para validar HMAC de esta tienda")
		c.JSON(http.StatusUnauthorized, response.WebhookResponse{
			Success: false,
			Message: "Configuración de seguridad faltante para esta tienda",
		})
		return
	}

	isTest := headers.ProbabilityTesting == "true"

	c.JSON(http.StatusOK, response.WebhookResponse{
		Success: true,
		Message: "Recibido",
	})

	h.enqueueWebhook(headers.Topic, headers.ShopDomain, bodyBytes, isTest)
}

func (h *ShopifyHandler) enqueueWebhook(topic string, shopDomain string, bodyBytes []byte, isTest bool) {
	ctx := context.Background()

	if h.rabbit != nil {
		msg := shopifyqueue.WebhookMessage{
			Topic:      topic,
			ShopDomain: shopDomain,
			IsTest:     isTest,
			Body:       bodyBytes,
		}
		payload, err := json.Marshal(msg)
		if err == nil {
			if err := h.rabbit.Publish(ctx, rabbitmq.QueueWebhooksShopifyReceived, payload); err == nil {
				return
			}
			h.logger.Warn(ctx).
				Str("topic", topic).
				Str("shop_domain", shopDomain).
				Msg("No se pudo encolar el webhook, procesando inline como fallback")
		}
	}

	go shopifyqueue.DispatchWebhook(ctx, h.useCase, h.logger, topic, shopDomain, bodyBytes, isTest)
}

func VerifyWebhookHMAC(message []byte, hmacHeader string, secret string) bool {
	if secret == "" {
		return false
	}

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(message)
	expectedMAC := mac.Sum(nil)
	expectedMACB64 := base64.StdEncoding.EncodeToString(expectedMAC)

	return hmac.Equal([]byte(hmacHeader), []byte(expectedMACB64))
}
