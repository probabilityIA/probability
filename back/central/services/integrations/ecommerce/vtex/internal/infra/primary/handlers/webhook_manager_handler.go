package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/vtex/internal/app/usecases"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/vtex/internal/domain"
)

type webhookRequest struct {
	IntegrationID uint  `json:"integration_id" binding:"required"`
	BusinessID    *uint `json:"business_id"`
	Force         bool  `json:"force"`
}

func webhookItemToResponse(item *domain.WebhookItem) gin.H {
	return gin.H{
		"id":       item.ID,
		"address":  item.Address,
		"statuses": item.Statuses,
		"is_ours":  item.IsOurs,
	}
}

func (h *vtexHandler) deliveryURL(integrationID uint) string {
	return usecases.WebhookDeliveryURL(h.baseURL, integrationID)
}

func (h *vtexHandler) ownedIntegrationID(c *gin.Context, rawID uint, bodyBusinessID *uint) (string, bool) {
	businessID, ok := h.resolveBusinessID(c, bodyBusinessID)
	if !ok {
		return "", false
	}

	integrationID := strconv.FormatUint(uint64(rawID), 10)
	if err := h.useCase.AssertIntegrationOwned(c.Request.Context(), integrationID, businessID); err != nil {
		h.respondUseCaseError(c, err)
		return "", false
	}

	return integrationID, true
}

func (h *vtexHandler) GetWebhookStatus(c *gin.Context) {
	rawID, err := strconv.ParseUint(c.Query("integration_id"), 10, 64)
	if err != nil || rawID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "integration_id es requerido"})
		return
	}

	var queryBusinessID *uint
	if v := c.Query("business_id"); v != "" {
		parsed, perr := strconv.ParseUint(v, 10, 64)
		if perr != nil || parsed == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "business_id invalido"})
			return
		}
		id := uint(parsed)
		queryBusinessID = &id
	}

	integrationID, ok := h.ownedIntegrationID(c, uint(rawID), queryBusinessID)
	if !ok {
		return
	}

	item, err := h.useCase.InspectWebhook(c.Request.Context(), integrationID, h.baseURL)
	if err != nil {
		h.respondUseCaseError(c, err)
		return
	}

	response := gin.H{
		"success":     true,
		"registered":  item != nil,
		"webhook_url": h.deliveryURL(uint(rawID)),
		"webhook":     nil,
	}
	if item != nil {
		response["webhook"] = webhookItemToResponse(item)
	}

	c.JSON(http.StatusOK, response)
}

func (h *vtexHandler) RegisterWebhook(c *gin.Context) {
	var req webhookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "integration_id es requerido"})
		return
	}

	integrationID, ok := h.ownedIntegrationID(c, req.IntegrationID, req.BusinessID)
	if !ok {
		return
	}

	if err := h.useCase.CreateWebhooks(c.Request.Context(), integrationID, h.baseURL, req.Force); err != nil {
		if errors.Is(err, domain.ErrForeignHookExists) {
			c.JSON(http.StatusConflict, gin.H{
				"success":       false,
				"error":         err.Error(),
				"foreign_hook":  true,
				"needs_confirm": true,
			})
			return
		}
		h.respondUseCaseError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Webhook registrado en VTEX",
	})
}

func (h *vtexHandler) UnregisterWebhook(c *gin.Context) {
	var req webhookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "integration_id es requerido"})
		return
	}

	integrationID, ok := h.ownedIntegrationID(c, req.IntegrationID, req.BusinessID)
	if !ok {
		return
	}

	if err := h.useCase.DeleteWebhook(c.Request.Context(), integrationID, ""); err != nil {
		h.respondUseCaseError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Webhook eliminado en VTEX",
	})
}
