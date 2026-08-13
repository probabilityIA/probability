package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type inventoryCompareRequest struct {
	IntegrationID uint     `json:"integration_id" binding:"required"`
	BusinessID    *uint    `json:"business_id"`
	Page          int      `json:"page"`
	PageSize      int      `json:"page_size"`
	SKUs          []string `json:"skus"`
}

func (h *meliHandler) CompareInventory(c *gin.Context) {
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

	resultado, err := h.useCase.CompareInventory(c.Request.Context(), integrationID, businessID, req.Page, req.PageSize, req.SKUs...)
	if err != nil {
		h.logger.Error(c.Request.Context()).Err(err).
			Uint("integration_id", req.IntegrationID).
			Msg("Error comparando inventario contra MercadoLibre")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
			"message": "No se pudo leer el stock de MercadoLibre",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": resultado})
}
