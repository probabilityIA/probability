package handlers

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/shopify/internal/infra/primary/handlers/request"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/shopify/internal/infra/primary/handlers/response"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/shopify/internal/infra/secondary/client/mappers"
	clientresponse "github.com/secamc93/probability/back/central/services/integrations/ecommerce/shopify/internal/infra/secondary/client/response"
)

// WebhookHandler maneja las peticiones de webhook de Shopify
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

	// Log del payload recibido (primeros 500 caracteres para no saturar logs)
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

	// HMAC Validation
	shopifySecret := os.Getenv("SHOPIFY_API_SECRET")
	h.logger.Info().
		Bool("has_secret", shopifySecret != "").
		Str("secret_prefix", func() string {
			if shopifySecret != "" && len(shopifySecret) > 4 {
				return shopifySecret[:4] + "..."
			}
			return "NO_SECRET"
		}()).
		Msg("🔐 Verificando HMAC")

	if shopifySecret != "" {
		if !VerifyWebhookHMAC(bodyBytes, headers.Hmac, shopifySecret) {
			h.logger.Error().
				Str("hmac_header", headers.Hmac).
				Msg("❌ Firma HMAC inválida")
			c.JSON(http.StatusUnauthorized, response.WebhookResponse{
				Success: false,
				Message: "Firma HMAC inválida",
			})
			return
		}
		h.logger.Info().Msg("✅ HMAC válido")
	} else {
		h.logger.Warn().Msg("⚠️ SHOPIFY_API_SECRET no configurado - omitiendo validación HMAC")
	}

	// Respond 200 OK as fast as possible as required by Shopify
	c.JSON(http.StatusOK, response.WebhookResponse{
		Success: true,
		Message: "Recibido",
	})

	// Procesar el webhook de forma asíncrona para no bloquear la respuesta
	go h.processWebhookAsync(headers.Topic, headers.ShopDomain, bodyBytes)
}

// processWebhookAsync procesa el webhook de forma asíncrona
func (h *ShopifyHandler) processWebhookAsync(topic string, shopDomain string, bodyBytes []byte) {
	ctx := context.Background()

	h.logger.Info().
		Str("topic", topic).
		Str("shop_domain", shopDomain).
		Msg("🔄 Iniciando procesamiento asíncrono del webhook")

	// Parsear el payload a Order de Shopify
	var orderResp clientresponse.Order
	if err := json.Unmarshal(bodyBytes, &orderResp); err != nil {
		h.logger.Error(ctx).
			Err(err).
			Str("topic", topic).
			Str("shop_domain", shopDomain).
			Msg("❌ Error al parsear payload de Shopify a Order")
		return
	}

	// Mapear la orden de Shopify a dominio
	mapped := mappers.MapOrderResponseToShopifyOrder(orderResp, bodyBytes, nil, 0, "shopify")
	shopifyOrder := &mapped

	// Procesar según el topic del webhook
	var err error
	switch topic {
	case "orders/create":
		h.logger.Info(ctx).
			Str("shop_domain", shopDomain).
			Str("order_id", orderResp.Name).
			Msg("📦 Procesando orden nueva (orders/create)")
		err = h.useCase.CreateOrder(ctx, shopDomain, shopifyOrder, bodyBytes)

	case "orders/paid":
		h.logger.Info(ctx).
			Str("shop_domain", shopDomain).
			Str("order_id", orderResp.Name).
			Msg("💰 Procesando orden pagada (orders/paid)")
		err = h.useCase.ProcessOrderPaid(ctx, shopDomain, shopifyOrder)

	case "orders/updated":
		h.logger.Info(ctx).
			Str("shop_domain", shopDomain).
			Str("order_id", orderResp.Name).
			Msg("🔄 Procesando orden actualizada (orders/updated)")
		err = h.useCase.ProcessOrderUpdated(ctx, shopDomain, shopifyOrder)

	case "orders/cancelled":
		h.logger.Info(ctx).
			Str("shop_domain", shopDomain).
			Str("order_id", orderResp.Name).
			Msg("❌ Procesando orden cancelada (orders/cancelled)")
		err = h.useCase.ProcessOrderCancelled(ctx, shopDomain, shopifyOrder)

	case "orders/fulfilled":
		h.logger.Info(ctx).
			Str("shop_domain", shopDomain).
			Str("order_id", orderResp.Name).
			Msg("✅ Procesando orden cumplida (orders/fulfilled)")
		err = h.useCase.ProcessOrderFulfilled(ctx, shopDomain, shopifyOrder)

	case "orders/partially_fulfilled":
		h.logger.Info(ctx).
			Str("shop_domain", shopDomain).
			Str("order_id", orderResp.Name).
			Msg("📦 Procesando orden parcialmente cumplida (orders/partially_fulfilled)")
		err = h.useCase.ProcessOrderPartiallyFulfilled(ctx, shopDomain, shopifyOrder)

	default:
		h.logger.Info(ctx).
			Str("topic", topic).
			Str("shop_domain", shopDomain).
			Msg("ℹ️ Topic no manejado, ignorando webhook")
		return
	}

	// Log del resultado
	if err != nil {
		h.logger.Error(ctx).
			Err(err).
			Str("topic", topic).
			Str("shop_domain", shopDomain).
			Str("order_id", orderResp.Name).
			Msg("❌ Error al procesar webhook de Shopify")
	} else {
		h.logger.Info(ctx).
			Str("topic", topic).
			Str("shop_domain", shopDomain).
			Str("order_id", orderResp.Name).
			Msg("✅ Webhook procesado exitosamente")
	}
}

// VerifyWebhookHMAC validates the Shopify HMAC signature
func VerifyWebhookHMAC(message []byte, hmacHeader string, secret string) bool {
	// If no secret provided (e.g. dev env without configuring it), we might skip validation or fail.
	// For security, we should fail or log warning. Ideally fail.
	// But to avoid blocking dev if env is missing:
	if secret == "" {
		// Try to load from env here or cleaner: pass from caller.
		// Caller passed "YOUR_SHOPIFY_API_SECRET", we need to fix that in caller.
		return false
	}

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(message)
	expectedMAC := mac.Sum(nil)
	expectedMACB64 := base64.StdEncoding.EncodeToString(expectedMAC)

	return hmac.Equal([]byte(hmacHeader), []byte(expectedMACB64))
}
