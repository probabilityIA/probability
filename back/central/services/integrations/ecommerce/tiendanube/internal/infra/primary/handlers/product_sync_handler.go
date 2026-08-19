package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/tiendanube/internal/app/usecases"
	"github.com/secamc93/probability/back/central/shared/inventorycompare"
)

type syncRequest struct {
	IntegrationID uint  `json:"integration_id" binding:"required"`
	BusinessID    *uint `json:"business_id"`
}

type applyRequest struct {
	IntegrationID uint     `json:"integration_id" binding:"required"`
	BusinessID    *uint    `json:"business_id"`
	Direction     string   `json:"direction" binding:"required"`
	Mode          string   `json:"mode"`
	SKUs          []string `json:"skus"`
}

type associateRequest struct {
	IntegrationID uint     `json:"integration_id" binding:"required"`
	BusinessID    *uint    `json:"business_id"`
	SKUs          []string `json:"skus"`
}

type inventoryCompareRequest struct {
	IntegrationID uint     `json:"integration_id" binding:"required"`
	BusinessID    *uint    `json:"business_id"`
	Page          int      `json:"page"`
	PageSize      int      `json:"page_size"`
	SKUs          []string `json:"skus"`
	Source        string   `json:"source"`
	OnlyDiff      bool     `json:"only_diff"`
	Search        string   `json:"q"`
}

func (h *tiendanubeHandler) resolveBusinessID(c *gin.Context, bodyBusinessID *uint) (uint, bool) {
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

func (h *tiendanubeHandler) ReconcileProducts(c *gin.Context) {
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
		h.useCase.ReconcileProductsAsync(context.Background(), integrationID, businessID, req.IntegrationID, correlationID)
	}()

	c.JSON(http.StatusAccepted, gin.H{
		"success":        true,
		"correlation_id": correlationID,
		"message":        "Comparacion de catalogo iniciada",
	})
}

func (h *tiendanubeHandler) SyncProducts(c *gin.Context) {
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
			h.logger.Error(ctx).Err(err).Msg("Error sincronizando productos con Tiendanube")
		}
	}()

	c.JSON(http.StatusAccepted, gin.H{
		"success":        true,
		"correlation_id": correlationID,
		"message":        "Sincronizacion de productos iniciada",
	})
}

func (h *tiendanubeHandler) ApplyProducts(c *gin.Context) {
	var req applyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "integration_id y direction son requeridos"})
		return
	}

	if req.Direction != usecases.DirectionToTiendanube && req.Direction != usecases.DirectionToProbability {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "direction debe ser to_tiendanube o to_probability",
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
	skus := req.SKUs

	go func() {
		ctx := context.Background()
		var err error
		switch {
		case direction == usecases.DirectionToTiendanube && mode == usecases.ModeCreate:
			err = h.useCase.ApplyProductsToTiendanube(ctx, integrationID, businessID, correlationID, skus...)
		case direction == usecases.DirectionToTiendanube && mode == usecases.ModeUpdate:
			err = h.useCase.UpdateProductsToTiendanube(ctx, integrationID, businessID, correlationID, skus...)
		case direction == usecases.DirectionToProbability && mode == usecases.ModeCreate:
			err = h.useCase.ApplyProductsToProbability(ctx, integrationID, businessID, correlationID, skus...)
		default:
			err = h.useCase.UpdateProductsToProbability(ctx, integrationID, businessID, correlationID, skus...)
		}
		if err != nil {
			h.logger.Error(ctx).Err(err).
				Str("direction", direction).
				Str("mode", mode).
				Msg("Error aplicando productos con Tiendanube")
		}
	}()

	c.JSON(http.StatusAccepted, gin.H{
		"success":        true,
		"correlation_id": correlationID,
		"message":        "Sincronizacion de productos iniciada",
	})
}

func (h *tiendanubeHandler) AssociateProducts(c *gin.Context) {
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
			h.logger.Error(ctx).Err(err).Msg("Error asociando productos con Tiendanube")
		}
	}()

	c.JSON(http.StatusAccepted, gin.H{
		"success":        true,
		"correlation_id": correlationID,
		"message":        "Asociacion de productos iniciada",
	})
}

func (h *tiendanubeHandler) SyncInventory(c *gin.Context) {
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
			h.logger.Error(ctx).Err(err).Msg("Error sincronizando inventario a Tiendanube")
		}
	}()

	c.JSON(http.StatusAccepted, gin.H{
		"success":        true,
		"correlation_id": correlationID,
		"message":        "Sincronizacion de inventario iniciada",
	})
}

func (h *tiendanubeHandler) CompareInventory(c *gin.Context) {
	var req inventoryCompareRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "integration_id es requerido"})
		return
	}

	businessID, ok := h.resolveBusinessID(c, req.BusinessID)
	if !ok {
		return
	}

	integrationID := strconv.FormatUint(uint64(req.IntegrationID), 10)

	if req.Source == "snapshot" {
		guardado, err := h.useCase.LoadInventoryCompare(c.Request.Context(), integrationID, businessID, inventorycompare.LoadOptions{
			Page:     req.Page,
			PageSize: req.PageSize,
			OnlyDiff: req.OnlyDiff,
			Search:   req.Search,
			SKUs:     req.SKUs,
		})
		if err != nil {
			h.logger.Error(c.Request.Context()).Err(err).
				Uint("integration_id", req.IntegrationID).
				Msg("Error leyendo la foto guardada del comparativo")
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   err.Error(),
				"message": "No se pudo leer la ultima comparacion guardada",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{"success": true, "data": guardado})
		return
	}

	resultado, err := h.useCase.CompareInventory(c.Request.Context(), integrationID, businessID, req.Page, req.PageSize, req.SKUs...)
	if err != nil {
		h.logger.Error(c.Request.Context()).Err(err).
			Uint("integration_id", req.IntegrationID).
			Msg("Error comparando inventario contra Tiendanube")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
			"message": "No se pudo leer el stock de Tiendanube",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": resultado})
}
