package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/shopify/internal/infra/primary/handlers/request"
	"github.com/secamc93/probability/back/central/services/integrations/shopify/internal/infra/primary/handlers/response"
	"github.com/secamc93/probability/back/central/services/integrations/shopify/internal/infra/secondary/client/mappers"
	clientresponse "github.com/secamc93/probability/back/central/services/integrations/shopify/internal/infra/secondary/client/response"
)

// WebhookHandler maneja las peticiones de webhook de Shopify
func (h *ShopifyHandler) WebhookHandler(c *gin.Context) {
	var headers request.WebhookHeaders

	if err := c.ShouldBindHeader(&headers); err != nil {
		h.logger.Error().Err(err).Msg("Error al validar headers del webhook")
		c.JSON(http.StatusBadRequest, response.WebhookResponse{
			Success: false,
			Message: "Headers requeridos faltantes o inv?lidos: " + err.Error(),
		})
		return
	}

	h.logger.Info().
		Str("topic", headers.Topic).
		Str("shop_domain", headers.ShopDomain).
		Str("webhook_id", headers.WebhookID).
		Msg("Webhook recibido de Shopify")

	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		h.logger.Error().Err(err).Msg("Error al leer el body del webhook")
		c.JSON(http.StatusBadRequest, response.WebhookResponse{
			Success: false,
			Message: "Error al leer el body de la petici?n",
		})
		return
	}

	// HMAC validation will be done in the use case after finding the integration by shopDomain

	c.Request.Body = io.NopCloser(strings.NewReader(string(bodyBytes)))

	var payload map[string]interface{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		h.logger.Error().Err(err).Msg("Error al parsear el payload del webhook")
		c.JSON(http.StatusBadRequest, response.WebhookResponse{
			Success: false,
			Message: "Payload JSON inv?lido: " + err.Error(),
		})
		return
	}

	switch headers.Topic {
	case "orders/create":
		var orderResp clientresponse.Order
		if err := json.Unmarshal(bodyBytes, &orderResp); err != nil {
			h.logger.Error().Err(err).Msg("Error al mapear payload a Order")
			c.JSON(http.StatusBadRequest, response.WebhookResponse{
				Success: false,
				Message: "Payload JSON inv?lido: " + err.Error(),
			})
			return
		}
		mapped := mappers.MapOrderResponseToShopifyOrder(orderResp, bodyBytes, nil, 0, "shopify")
		shopifyOrder := &mapped
		if err := h.useCase.CreateOrder(c.Request.Context(), headers.ShopDomain, shopifyOrder, bodyBytes); err != nil {
			h.logger.Error().
				Err(err).
				Str("topic", headers.Topic).
				Str("shop_domain", headers.ShopDomain).
				Msg("Error al procesar webhook")
			c.JSON(http.StatusInternalServerError, response.WebhookResponse{
				Success: false,
				Message: "Error al procesar webhook",
			})
			return
		}
	case "orders/paid":
		var orderResp clientresponse.Order
		if err := json.Unmarshal(bodyBytes, &orderResp); err != nil {
			h.logger.Error().Err(err).Msg("Error al mapear payload a Order")
			c.JSON(http.StatusBadRequest, response.WebhookResponse{
				Success: false,
				Message: "Payload JSON inv?lido: " + err.Error(),
			})
			return
		}
		mapped := mappers.MapOrderResponseToShopifyOrder(orderResp, bodyBytes, nil, 0, "shopify")
		shopifyOrder := &mapped
		if err := h.useCase.ProcessOrderPaid(c.Request.Context(), headers.ShopDomain, shopifyOrder); err != nil {
			h.logger.Error().
				Err(err).
				Str("topic", headers.Topic).
				Str("shop_domain", headers.ShopDomain).
				Msg("Error al procesar webhook")
			c.JSON(http.StatusInternalServerError, response.WebhookResponse{
				Success: false,
				Message: "Error al procesar webhook",
			})
			return
		}
	case "orders/updated":
		var orderResp clientresponse.Order
		if err := json.Unmarshal(bodyBytes, &orderResp); err != nil {
			h.logger.Error().Err(err).Msg("Error al mapear payload a Order")
			c.JSON(http.StatusBadRequest, response.WebhookResponse{
				Success: false,
				Message: "Payload JSON inv?lido: " + err.Error(),
			})
			return
		}
		mapped := mappers.MapOrderResponseToShopifyOrder(orderResp, bodyBytes, nil, 0, "shopify")
		shopifyOrder := &mapped
		if err := h.useCase.ProcessOrderUpdated(c.Request.Context(), headers.ShopDomain, shopifyOrder); err != nil {
			h.logger.Error().
				Err(err).
				Str("topic", headers.Topic).
				Str("shop_domain", headers.ShopDomain).
				Msg("Error al procesar webhook")
			c.JSON(http.StatusInternalServerError, response.WebhookResponse{
				Success: false,
				Message: "Error al procesar webhook",
			})
			return
		}
	case "orders/cancelled":
		var orderResp clientresponse.Order
		if err := json.Unmarshal(bodyBytes, &orderResp); err != nil {
			h.logger.Error().Err(err).Msg("Error al mapear payload a Order")
			c.JSON(http.StatusBadRequest, response.WebhookResponse{
				Success: false,
				Message: "Payload JSON inv?lido: " + err.Error(),
			})
			return
		}
		mapped := mappers.MapOrderResponseToShopifyOrder(orderResp, bodyBytes, nil, 0, "shopify")
		shopifyOrder := &mapped
		if err := h.useCase.ProcessOrderCancelled(c.Request.Context(), headers.ShopDomain, shopifyOrder); err != nil {
			h.logger.Error().
				Err(err).
				Str("topic", headers.Topic).
				Str("shop_domain", headers.ShopDomain).
				Msg("Error al procesar webhook")
			c.JSON(http.StatusInternalServerError, response.WebhookResponse{
				Success: false,
				Message: "Error al procesar webhook",
			})
			return
		}
	case "orders/fulfilled":
		var orderResp clientresponse.Order
		if err := json.Unmarshal(bodyBytes, &orderResp); err != nil {
			h.logger.Error().Err(err).Msg("Error al mapear payload a Order")
			c.JSON(http.StatusBadRequest, response.WebhookResponse{
				Success: false,
				Message: "Payload JSON inv?lido: " + err.Error(),
			})
			return
		}
		mapped := mappers.MapOrderResponseToShopifyOrder(orderResp, bodyBytes, nil, 0, "shopify")
		shopifyOrder := &mapped
		if err := h.useCase.ProcessOrderFulfilled(c.Request.Context(), headers.ShopDomain, shopifyOrder); err != nil {
			h.logger.Error().
				Err(err).
				Str("topic", headers.Topic).
				Str("shop_domain", headers.ShopDomain).
				Msg("Error al procesar webhook")
			c.JSON(http.StatusInternalServerError, response.WebhookResponse{
				Success: false,
				Message: "Error al procesar webhook",
			})
			return
		}
	case "orders/partially_fulfilled":
		var orderResp clientresponse.Order
		if err := json.Unmarshal(bodyBytes, &orderResp); err != nil {
			h.logger.Error().Err(err).Msg("Error al mapear payload a Order")
			c.JSON(http.StatusBadRequest, response.WebhookResponse{
				Success: false,
				Message: "Payload JSON inv?lido: " + err.Error(),
			})
			return
		}
		mapped := mappers.MapOrderResponseToShopifyOrder(orderResp, bodyBytes, nil, 0, "shopify")
		shopifyOrder := &mapped
		if err := h.useCase.ProcessOrderPartiallyFulfilled(c.Request.Context(), headers.ShopDomain, shopifyOrder); err != nil {
			h.logger.Error().
				Err(err).
				Str("topic", headers.Topic).
				Str("shop_domain", headers.ShopDomain).
				Msg("Error al procesar webhook")
			c.JSON(http.StatusInternalServerError, response.WebhookResponse{
				Success: false,
				Message: "Error al procesar webhook",
			})
			return
		}
	default:
		h.logger.Info().Str("topic", headers.Topic).Msg("Topic no manejado, ignorando")
	}

	// Responder con ?xito
	c.JSON(http.StatusOK, response.WebhookResponse{
		Success: true,
		Message: "Webhook procesado exitosamente",
	})
}
