package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/secamc93/probability/back/central/shared/inventorycompare"
)

type inventorySyncRequest struct {
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

func (h *wooCommerceHandler) CompareInventory(c *gin.Context) {
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
			Msg("Error comparando inventario contra WooCommerce")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
			"message": "No se pudo leer el stock de WooCommerce",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": resultado})
}

func (h *wooCommerceHandler) SyncInventory(c *gin.Context) {
	var req inventorySyncRequest
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
		if err := h.useCase.SyncInventory(ctx, integrationID, businessID, correlationID, req.SKUs...); err != nil {
			h.logger.Error(ctx).Err(err).Msg("Error sincronizando inventario a WooCommerce")
		}
	}()

	c.JSON(http.StatusAccepted, gin.H{
		"success":        true,
		"correlation_id": correlationID,
		"message":        "Sincronizacion de inventario iniciada",
	})
}
