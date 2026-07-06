package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/secamc93/probability/back/central/services/auth/middleware"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/meli/internal/domain"
)

func (h *meliHandler) resolveBusinessID(c *gin.Context, bodyBusinessID *uint) (uint, bool) {
	businessID, ok := middleware.GetBusinessIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "contexto de negocio no encontrado"})
		return 0, false
	}
	if businessID == 0 {
		if bodyBusinessID == nil || *bodyBusinessID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "business_id es requerido para super admin"})
			return 0, false
		}
		businessID = *bodyBusinessID
	}
	return businessID, true
}

type reconcileRequest struct {
	IntegrationID uint  `json:"integration_id" binding:"required"`
	BusinessID    *uint `json:"business_id"`
}

func briefsToResponse(items []domain.ProductBrief) []gin.H {
	out := make([]gin.H, 0, len(items))
	for _, b := range items {
		out = append(out, gin.H{"sku": b.SKU, "name": b.Name})
	}
	return out
}

func (h *meliHandler) ReconcileProducts(c *gin.Context) {
	var req reconcileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "integration_id es requerido"})
		return
	}
	businessID, ok := h.resolveBusinessID(c, req.BusinessID)
	if !ok {
		return
	}

	integrationID := strconv.FormatUint(uint64(req.IntegrationID), 10)
	result, err := h.useCase.ReconcileProducts(c.Request.Context(), integrationID, businessID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":             true,
		"matched":             result.Matched,
		"only_in_probability": briefsToResponse(result.OnlyInProbability),
		"only_in_meli":        briefsToResponse(result.OnlyInMeli),
		"probability_no_sku":  result.ProbabilityNoSKU,
		"meli_no_sku":         result.MeliNoSKU,
	})
}

type applyRequest struct {
	IntegrationID uint   `json:"integration_id" binding:"required"`
	BusinessID    *uint  `json:"business_id"`
	Direction     string `json:"direction" binding:"required"`
}

func (h *meliHandler) ApplyProducts(c *gin.Context) {
	var req applyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "integration_id y direction son requeridos"})
		return
	}
	if req.Direction != "to_meli" && req.Direction != "to_probability" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "direction debe ser to_meli o to_probability"})
		return
	}
	businessID, ok := h.resolveBusinessID(c, req.BusinessID)
	if !ok {
		return
	}

	integrationID := strconv.FormatUint(uint64(req.IntegrationID), 10)
	correlationID := uuid.New().String()
	direction := req.Direction

	go func() {
		ctx := context.Background()
		var err error
		if direction == "to_meli" {
			err = h.useCase.ApplyProductsToMeli(ctx, integrationID, businessID, correlationID)
		} else {
			err = h.useCase.ApplyProductsToProbability(ctx, integrationID, businessID, correlationID)
		}
		if err != nil {
			h.logger.Error(ctx).Err(err).Str("direction", direction).Msg("Error aplicando reconciliacion de productos MeLi")
		}
	}()

	c.JSON(http.StatusAccepted, gin.H{
		"success":        true,
		"correlation_id": correlationID,
		"message":        "Reconciliacion iniciada",
	})
}
