package handlers

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/jumpseller/internal/app/usecases"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/jumpseller/internal/domain"
)

type syncRequest struct {
	IntegrationID uint  `json:"integration_id" binding:"required"`
	BusinessID    *uint `json:"business_id"`
}

type applyRequest struct {
	IntegrationID uint   `json:"integration_id" binding:"required"`
	BusinessID    *uint  `json:"business_id"`
	Direction     string `json:"direction" binding:"required"`
	Mode          string `json:"mode"`
}

type associateRequest struct {
	IntegrationID uint     `json:"integration_id" binding:"required"`
	BusinessID    *uint    `json:"business_id"`
	SKUs          []string `json:"skus"`
}

func (h *jumpsellerHandler) resolveBusinessID(c *gin.Context, bodyBusinessID *uint) (uint, bool) {
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

func briefsToResponse(items []domain.ProductBrief) []gin.H {
	out := make([]gin.H, 0, len(items))
	for _, item := range items {
		out = append(out, gin.H{"sku": item.SKU, "name": item.Name})
	}
	return out
}

func (h *jumpsellerHandler) ReconcileProducts(c *gin.Context) {
	var req syncRequest
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
		if errors.Is(err, domain.ErrIntegrationNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":                true,
		"matched":                result.Matched,
		"matched_not_associated": briefsToResponse(result.MatchedNotAssociated),
		"only_in_probability":    briefsToResponse(result.OnlyInProbability),
		"only_in_jumpseller":     briefsToResponse(result.OnlyInJumpseller),
		"probability_no_sku":     result.ProbabilityNoSKU,
		"jumpseller_no_sku":      result.JumpsellerNoSKU,
	})
}

func (h *jumpsellerHandler) SyncProducts(c *gin.Context) {
	var req syncRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "integration_id es requerido"})
		return
	}

	businessID, ok := h.resolveBusinessID(c, req.BusinessID)
	if !ok {
		return
	}

	integrationID := strconv.FormatUint(uint64(req.IntegrationID), 10)
	correlationID := uuid.New().String()

	go func() {
		ctx := context.Background()
		if err := h.useCase.SyncProducts(ctx, integrationID, businessID, correlationID); err != nil {
			h.logger.Error(ctx).Err(err).Msg("Error sincronizando productos con Jumpseller")
		}
	}()

	c.JSON(http.StatusAccepted, gin.H{
		"success":        true,
		"correlation_id": correlationID,
		"message":        "Sincronizacion de productos iniciada",
	})
}

func (h *jumpsellerHandler) ApplyProducts(c *gin.Context) {
	var req applyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "integration_id y direction son requeridos"})
		return
	}

	if req.Direction != usecases.DirectionToJumpseller && req.Direction != usecases.DirectionToProbability {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "direction debe ser to_jumpseller o to_probability",
		})
		return
	}

	mode := req.Mode
	if mode == "" {
		mode = usecases.ModeCreate
	}
	if mode != usecases.ModeCreate && mode != usecases.ModeUpdate {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "mode debe ser create o update",
		})
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
		switch {
		case direction == usecases.DirectionToJumpseller && mode == usecases.ModeCreate:
			err = h.useCase.ApplyProductsToJumpseller(ctx, integrationID, businessID, correlationID)
		case direction == usecases.DirectionToJumpseller && mode == usecases.ModeUpdate:
			err = h.useCase.UpdateProductsToJumpseller(ctx, integrationID, businessID, correlationID)
		case direction == usecases.DirectionToProbability && mode == usecases.ModeCreate:
			err = h.useCase.ApplyProductsToProbability(ctx, integrationID, businessID, correlationID)
		default:
			err = h.useCase.UpdateProductsToProbability(ctx, integrationID, businessID, correlationID)
		}
		if err != nil {
			h.logger.Error(ctx).Err(err).
				Str("direction", direction).
				Str("mode", mode).
				Msg("Error aplicando productos con Jumpseller")
		}
	}()

	c.JSON(http.StatusAccepted, gin.H{
		"success":        true,
		"correlation_id": correlationID,
		"message":        "Sincronizacion de productos iniciada",
	})
}

func (h *jumpsellerHandler) AssociateProducts(c *gin.Context) {
	var req associateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "integration_id es requerido"})
		return
	}

	businessID, ok := h.resolveBusinessID(c, req.BusinessID)
	if !ok {
		return
	}

	integrationID := strconv.FormatUint(uint64(req.IntegrationID), 10)
	correlationID := uuid.New().String()
	skus := req.SKUs

	go func() {
		ctx := context.Background()
		if err := h.useCase.AssociateProducts(ctx, integrationID, businessID, correlationID, skus); err != nil {
			h.logger.Error(ctx).Err(err).Msg("Error asociando productos con Jumpseller")
		}
	}()

	c.JSON(http.StatusAccepted, gin.H{
		"success":        true,
		"correlation_id": correlationID,
		"message":        "Asociacion de productos iniciada",
	})
}

func (h *jumpsellerHandler) SyncInventory(c *gin.Context) {
	var req syncRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "integration_id es requerido"})
		return
	}

	businessID, ok := h.resolveBusinessID(c, req.BusinessID)
	if !ok {
		return
	}

	integrationID := strconv.FormatUint(uint64(req.IntegrationID), 10)
	correlationID := uuid.New().String()

	go func() {
		ctx := context.Background()
		if err := h.useCase.SyncInventory(ctx, integrationID, businessID, correlationID); err != nil {
			h.logger.Error(ctx).Err(err).Msg("Error sincronizando inventario a Jumpseller")
		}
	}()

	c.JSON(http.StatusAccepted, gin.H{
		"success":        true,
		"correlation_id": correlationID,
		"message":        "Sincronizacion de inventario iniciada",
	})
}
