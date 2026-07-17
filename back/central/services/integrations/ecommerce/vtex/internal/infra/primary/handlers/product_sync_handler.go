package handlers

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/vtex/internal/domain"
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

func (h *vtexHandler) resolveBusinessID(c *gin.Context, bodyBusinessID *uint) (uint, bool) {
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

func (h *vtexHandler) respondUseCaseError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrIntegrationNotFound), errors.Is(err, domain.ErrIntegrationNotOwned):
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": domain.ErrIntegrationNotFound.Error()})
	case errors.Is(err, domain.ErrInvalidCredentials):
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": err.Error()})
	case errors.Is(err, domain.ErrMissingAccountName), errors.Is(err, domain.ErrMissingAppKey), errors.Is(err, domain.ErrMissingAppToken),
		errors.Is(err, domain.ErrInventorySyncDisabled), errors.Is(err, domain.ErrNoWarehousesMapped):
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
	}
}

func (h *vtexHandler) ReconcileProducts(c *gin.Context) {
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
		h.respondUseCaseError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":                true,
		"matched":                result.Matched,
		"matched_not_associated": briefsToResponse(result.MatchedNotAssociated),
		"only_in_probability":    briefsToResponse(result.OnlyInProbability),
		"only_in_vtex":           briefsToResponse(result.OnlyInVTEX),
		"probability_no_sku":     result.ProbabilityNoSKU,
		"vtex_no_sku":            result.VTEXNoSKU,
	})
}

func (h *vtexHandler) SyncProducts(c *gin.Context) {
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
			h.logger.Error(ctx).Err(err).Msg("Error sincronizando productos con VTEX")
		}
	}()

	c.JSON(http.StatusAccepted, gin.H{
		"success":        true,
		"correlation_id": correlationID,
		"message":        "Sincronizacion de productos iniciada",
	})
}

func (h *vtexHandler) ApplyProducts(c *gin.Context) {
	var req applyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "integration_id y direction son requeridos"})
		return
	}

	if req.Direction != domain.DirectionToProbability {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "direction solo admite 'to_probability': la creacion de productos en VTEX requiere categoria y marca y no esta soportada",
		})
		return
	}

	businessID, ok := h.resolveBusinessID(c, req.BusinessID)
	if !ok {
		return
	}

	mode := req.Mode
	if mode == "" {
		mode = domain.ModeCreate
	}
	if mode != domain.ModeCreate && mode != domain.ModeUpdate {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "mode admite 'create' o 'update'"})
		return
	}

	integrationID := strconv.FormatUint(uint64(req.IntegrationID), 10)
	correlationID := uuid.New().String()

	go func() {
		ctx := context.Background()
		var err error
		if mode == domain.ModeCreate {
			err = h.useCase.ApplyProductsToProbability(ctx, integrationID, businessID, correlationID)
		} else {
			err = h.useCase.UpdateProductsToProbability(ctx, integrationID, businessID, correlationID)
		}
		if err != nil {
			h.logger.Error(ctx).Err(err).Msg("Error aplicando productos de VTEX hacia Probability")
		}
	}()

	c.JSON(http.StatusAccepted, gin.H{
		"success":        true,
		"correlation_id": correlationID,
		"message":        "Aplicacion de productos iniciada",
	})
}

func (h *vtexHandler) AssociateProducts(c *gin.Context) {
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

	go func() {
		ctx := context.Background()
		if err := h.useCase.AssociateProducts(ctx, integrationID, businessID, correlationID, req.SKUs); err != nil {
			h.logger.Error(ctx).Err(err).Msg("Error asociando productos con VTEX")
		}
	}()

	c.JSON(http.StatusAccepted, gin.H{
		"success":        true,
		"correlation_id": correlationID,
		"message":        "Asociacion de productos iniciada",
	})
}

func (h *vtexHandler) SyncInventory(c *gin.Context) {
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
			h.logger.Error(ctx).Err(err).Msg("Error sincronizando inventario con VTEX")
		}
	}()

	c.JSON(http.StatusAccepted, gin.H{
		"success":        true,
		"correlation_id": correlationID,
		"message":        "Sincronizacion de inventario iniciada",
	})
}

func (h *vtexHandler) GetWarehouses(c *gin.Context) {
	integrationIDParam := c.Query("integration_id")
	if integrationIDParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "integration_id es requerido"})
		return
	}

	var bodyBusinessID *uint
	if v := c.Query("business_id"); v != "" {
		parsed, err := strconv.ParseUint(v, 10, 64)
		if err != nil || parsed == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "business_id invalido"})
			return
		}
		id := uint(parsed)
		bodyBusinessID = &id
	}

	businessID, ok := h.resolveBusinessID(c, bodyBusinessID)
	if !ok {
		return
	}

	info, err := h.useCase.GetWarehouses(c.Request.Context(), integrationIDParam, businessID)
	if err != nil {
		h.respondUseCaseError(c, err)
		return
	}

	warehouses := make([]gin.H, 0, len(info.Warehouses))
	for _, w := range info.Warehouses {
		warehouses = append(warehouses, gin.H{
			"id":        w.ID,
			"name":      w.Name,
			"is_active": w.IsActive,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"warehouses": warehouses,
	})
}
