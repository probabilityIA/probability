package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ComplianceWebhookHandler maneja TODOS los webhooks de compliance en un solo endpoint
//
//	@Summary		Webhook unificado de compliance
//	@Description	Maneja todos los webhooks de GDPR/CCPA (data_request, redact, shop_redact)
//	@Tags			Shopify Compliance
//	@Accept			json
//	@Produce		json
//	@Param			X-Shopify-Topic		header	string	true	"Shopify webhook topic"
//	@Param			X-Shopify-Hmac-Sha256	header	string	true	"HMAC signature"
//	@Param			X-Shopify-Shop-Domain	header	string	true	"Shop domain"
//	@Success		200					{object}	map[string]interface{}
//	@Failure		401					{object}	map[string]interface{}
//	@Router			/integrations/shopify/webhooks/compliance [post]
func (h *ShopifyHandler) ComplianceWebhookHandler(c *gin.Context) {
	topic := c.GetHeader("X-Shopify-Topic")

	// Leer el body
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		h.logger.Error().Err(err).Msg("Error al leer body del webhook de compliance")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Validar HMAC (OBLIGATORIO para compliance)
	hmacHeader := c.GetHeader("X-Shopify-Hmac-Sha256")
	shopDomain := c.GetHeader("X-Shopify-Shop-Domain")

	// Recuperar el secreto específico de esta tienda
	shopifySecret, err := h.useCase.GetClientSecretByShopDomain(c.Request.Context(), shopDomain)
	if err != nil {
		h.logger.Warn().
			Err(err).
			Str("shop_domain", shopDomain).
			Msg("⚠️ No se pudo recuperar el secreto específico para webhook de compliance, intentando global")

		shopifySecret = h.config.Get("SHOPIFY_CLIENT_SECRET")
	}

	if shopifySecret == "" {
		h.logger.Error().Str("topic", topic).Msg("No hay secreto configurado para validar compliance webhook")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Webhook signature validation not configured"})
		return
	}

	if !VerifyWebhookHMAC(bodyBytes, hmacHeader, shopifySecret) {
		h.logger.Error().
			Str("topic", topic).
			Str("shop_domain", shopDomain).
			Msg("HMAC inválido en compliance webhook")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid HMAC signature"})
		return
	}

	h.logger.Info().
		Str("topic", topic).
		Str("shop_domain", c.GetHeader("X-Shopify-Shop-Domain")).
		Msg("Compliance webhook recibido")

	// Responder inmediatamente con 200 OK (requisito de Shopify)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Compliance webhook received",
	})

	// Procesar según el topic
	go func() {
		switch topic {
		case "customers/data_request":
			var payload CustomerDataRequestPayload
			if err := json.Unmarshal(bodyBytes, &payload); err == nil {
				h.logger.Info().
					Int64("customer_id", payload.Customer.ID).
					Msg("Procesando solicitud de datos de cliente (GDPR)")
				// TODO: Implementar lógica de procesamiento
			}

		case "customers/redact":
			var payload CustomerRedactPayload
			if err := json.Unmarshal(bodyBytes, &payload); err == nil {
				h.logger.Info().
					Int64("customer_id", payload.Customer.ID).
					Msg("Procesando eliminación de datos de cliente (GDPR)")
				// TODO: Implementar lógica de procesamiento
			}

		case "shop/redact":
			var payload ShopRedactPayload
			if err := json.Unmarshal(bodyBytes, &payload); err == nil {
				h.logger.Info().
					Int64("shop_id", payload.ShopID).
					Msg("Procesando eliminación de datos de tienda")
				// TODO: Implementar lógica de procesamiento
			}
		}
	}()
}

// CustomerDataRequestPayload representa la estructura del webhook customers/data_request
type CustomerDataRequestPayload struct {
	ShopID          int64   `json:"shop_id"`
	ShopDomain      string  `json:"shop_domain"`
	OrdersRequested []int64 `json:"orders_requested"`
	Customer        struct {
		ID    int64  `json:"id"`
		Email string `json:"email"`
		Phone string `json:"phone"`
	} `json:"customer"`
}

// CustomerRedactPayload representa la estructura del webhook customers/redact
type CustomerRedactPayload struct {
	ShopID     int64  `json:"shop_id"`
	ShopDomain string `json:"shop_domain"`
	Customer   struct {
		ID    int64  `json:"id"`
		Email string `json:"email"`
		Phone string `json:"phone"`
	} `json:"customer"`
	OrdersToRedact []int64 `json:"orders_to_redact"`
}

// ShopRedactPayload representa la estructura del webhook shop/redact
type ShopRedactPayload struct {
	ShopID     int64  `json:"shop_id"`
	ShopDomain string `json:"shop_domain"`
}

// CustomerDataRequestHandler maneja solicitudes de datos de clientes (GDPR)
//
//	@Summary		Webhook de solicitud de datos de cliente
//	@Description	Maneja solicitudes de acceso a datos de clientes según GDPR/CCPA
//	@Tags			Shopify Compliance
//	@Accept			json
//	@Produce		json
//	@Param			X-Shopify-Topic		header	string	true	"Shopify webhook topic"
//	@Param			X-Shopify-Hmac-Sha256	header	string	true	"HMAC signature"
//	@Param			X-Shopify-Shop-Domain	header	string	true	"Shop domain"
//	@Success		200					{object}	map[string]interface{}
//	@Failure		401					{object}	map[string]interface{}
//	@Router			/integrations/shopify/webhooks/customers/data_request [post]
func (h *ShopifyHandler) CustomerDataRequestHandler(c *gin.Context) {
	// Leer el body
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		h.logger.Error().Err(err).Msg("Error al leer body del webhook customers/data_request")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Validar HMAC
	hmacHeader := c.GetHeader("X-Shopify-Hmac-Sha256")
	shopDomain := c.GetHeader("X-Shopify-Shop-Domain")

	shopifySecret, err := h.useCase.GetClientSecretByShopDomain(c.Request.Context(), shopDomain)
	if err != nil {
		shopifySecret = h.config.Get("SHOPIFY_CLIENT_SECRET")
	}

	if shopifySecret != "" && !VerifyWebhookHMAC(bodyBytes, hmacHeader, shopifySecret) {
		h.logger.Error().Str("shop_domain", shopDomain).Msg("HMAC inválido en customers/data_request")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid HMAC signature"})
		return
	}

	// Parsear payload
	var payload CustomerDataRequestPayload
	if err := json.Unmarshal(bodyBytes, &payload); err != nil {
		h.logger.Error().Err(err).Msg("Error al parsear payload customers/data_request")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON payload"})
		return
	}

	h.logger.Info().
		Str("shop_domain", payload.ShopDomain).
		Int64("customer_id", payload.Customer.ID).
		Str("customer_email", payload.Customer.Email).
		Msg("Solicitud de datos de cliente recibida (GDPR)")

	// Responder inmediatamente con 200 OK (requisito de Shopify)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Customer data request received",
	})

	// TODO: Procesar de forma asíncrona
	// - Buscar todos los datos del cliente en la base de datos
	// - Generar archivo con los datos del cliente
	// - Enviar datos al cliente o a Shopify según el proceso
	// - Completar dentro de 30 días
	go func() {
		h.logger.Info().
			Int64("customer_id", payload.Customer.ID).
			Msg("Procesando solicitud de datos de cliente de forma asíncrona")
		// Aquí iría la lógica de procesamiento asíncrono
	}()
}

// CustomerRedactHandler maneja solicitudes de eliminación de datos de clientes (GDPR)
//
//	@Summary		Webhook de eliminación de datos de cliente
//	@Description	Maneja solicitudes de eliminación de datos de clientes según GDPR/CCPA
//	@Tags			Shopify Compliance
//	@Accept			json
//	@Produce		json
//	@Param			X-Shopify-Topic		header	string	true	"Shopify webhook topic"
//	@Param			X-Shopify-Hmac-Sha256	header	string	true	"HMAC signature"
//	@Param			X-Shopify-Shop-Domain	header	string	true	"Shop domain"
//	@Success		200					{object}	map[string]interface{}
//	@Failure		401					{object}	map[string]interface{}
//	@Router			/integrations/shopify/webhooks/customers/redact [post]
func (h *ShopifyHandler) CustomerRedactHandler(c *gin.Context) {
	// Leer el body
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		h.logger.Error().Err(err).Msg("Error al leer body del webhook customers/redact")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Validar HMAC
	hmacHeader := c.GetHeader("X-Shopify-Hmac-Sha256")
	shopDomain := c.GetHeader("X-Shopify-Shop-Domain")

	shopifySecret, err := h.useCase.GetClientSecretByShopDomain(c.Request.Context(), shopDomain)
	if err != nil {
		shopifySecret = h.config.Get("SHOPIFY_CLIENT_SECRET")
	}

	if shopifySecret != "" && !VerifyWebhookHMAC(bodyBytes, hmacHeader, shopifySecret) {
		h.logger.Error().Str("shop_domain", shopDomain).Msg("HMAC inválido en customers/redact")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid HMAC signature"})
		return
	}

	// Parsear payload
	var payload CustomerRedactPayload
	if err := json.Unmarshal(bodyBytes, &payload); err != nil {
		h.logger.Error().Err(err).Msg("Error al parsear payload customers/redact")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON payload"})
		return
	}

	h.logger.Info().
		Str("shop_domain", payload.ShopDomain).
		Int64("customer_id", payload.Customer.ID).
		Str("customer_email", payload.Customer.Email).
		Msg("Solicitud de eliminación de datos de cliente recibida (GDPR)")

	// Responder inmediatamente con 200 OK (requisito de Shopify)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Customer redaction request received",
	})

	// TODO: Procesar de forma asíncrona
	// - Anonimizar o eliminar datos del cliente de la base de datos
	// - Eliminar datos de pedidos relacionados (si aplica)
	// - Actualizar registros de auditoría
	// - Completar dentro de 30 días
	go func() {
		h.logger.Info().
			Int64("customer_id", payload.Customer.ID).
			Msg("Procesando eliminación de datos de cliente de forma asíncrona")
		// Aquí iría la lógica de procesamiento asíncrono
	}()
}

// ShopRedactHandler maneja solicitudes de eliminación de datos de tienda (cuando se desinstala la app)
//
//	@Summary		Webhook de eliminación de datos de tienda
//	@Description	Maneja la eliminación de datos cuando la app es desinstalada de una tienda
//	@Tags			Shopify Compliance
//	@Accept			json
//	@Produce		json
//	@Param			X-Shopify-Topic		header	string	true	"Shopify webhook topic"
//	@Param			X-Shopify-Hmac-Sha256	header	string	true	"HMAC signature"
//	@Param			X-Shopify-Shop-Domain	header	string	true	"Shop domain"
//	@Success		200					{object}	map[string]interface{}
//	@Failure		401					{object}	map[string]interface{}
//	@Router			/integrations/shopify/webhooks/shop/redact [post]
func (h *ShopifyHandler) ShopRedactHandler(c *gin.Context) {
	// Leer el body
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		h.logger.Error().Err(err).Msg("Error al leer body del webhook shop/redact")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Validar HMAC
	hmacHeader := c.GetHeader("X-Shopify-Hmac-Sha256")
	shopDomain := c.GetHeader("X-Shopify-Shop-Domain")

	shopifySecret, err := h.useCase.GetClientSecretByShopDomain(c.Request.Context(), shopDomain)
	if err != nil {
		shopifySecret = h.config.Get("SHOPIFY_CLIENT_SECRET")
	}

	if shopifySecret != "" && !VerifyWebhookHMAC(bodyBytes, hmacHeader, shopifySecret) {
		h.logger.Error().Str("shop_domain", shopDomain).Msg("HMAC inválido en shop/redact")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid HMAC signature"})
		return
	}

	// Parsear payload
	var payload ShopRedactPayload
	if err := json.Unmarshal(bodyBytes, &payload); err != nil {
		h.logger.Error().Err(err).Msg("Error al parsear payload shop/redact")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON payload"})
		return
	}

	h.logger.Info().
		Str("shop_domain", payload.ShopDomain).
		Int64("shop_id", payload.ShopID).
		Msg("Solicitud de eliminación de datos de tienda recibida (app desinstalada)")

	// Responder inmediatamente con 200 OK (requisito de Shopify)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Shop redaction request received",
	})

	// TODO: Procesar de forma asíncrona
	// - Eliminar todos los datos relacionados con la tienda
	// - Desactivar la integración en la base de datos
	// - Limpiar webhooks y configuraciones
	// - Completar dentro de 30 días (o antes si no hay requisitos legales de retención)
	go func() {
		h.logger.Info().
			Int64("shop_id", payload.ShopID).
			Str("shop_domain", payload.ShopDomain).
			Msg("Procesando eliminación de datos de tienda de forma asíncrona")
		// Aquí iría la lógica de procesamiento asíncrono
	}()
}
