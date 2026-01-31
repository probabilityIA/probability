package handlers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/shopify/internal/infra/primary/handlers/request"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/shopify/internal/infra/primary/handlers/response"
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

	// HMAC Validation
	shopifySecret := os.Getenv("SHOPIFY_API_SECRET")
	if shopifySecret != "" {
		if !VerifyWebhookHMAC(bodyBytes, headers.Hmac, shopifySecret) {
			h.logger.Error().Msg("Firma HMAC inválida")
			c.JSON(http.StatusUnauthorized, response.WebhookResponse{
				Success: false,
				Message: "Firma HMAC inválida",
			})
			return
		}
	}

	// Respond 200 OK as fast as possible as required by Shopify
	c.JSON(http.StatusOK, response.WebhookResponse{
		Success: true,
		Message: "Recibido",
	})

	// TODO: En el futuro, si se requiere procesar el payload (ej. orders/create),
	// se debería hacer de forma asíncrona (Goroutine o Queue) para no bloquear la respuesta.
	h.logger.Debug().Str("topic", headers.Topic).Msg("Webhook aceptado (procesamiento asíncrono no implementado)")
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
