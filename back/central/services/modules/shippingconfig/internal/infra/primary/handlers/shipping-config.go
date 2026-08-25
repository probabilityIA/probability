package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/shippingconfig/internal/domain"
)

func (h *Handlers) GetOverview(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "business_id es requerido"})
		return
	}

	overview, err := h.uc.GetOverview(c.Request.Context(), businessID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": overview})
}

func (h *Handlers) SaveBusinessConfig(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "business_id es requerido"})
		return
	}

	var req domain.SaveConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}
	req.UpdatedBy, req.UpdatedByName = h.resolveActor(c)

	cfg, err := h.uc.SaveBusinessConfig(c.Request.Context(), businessID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": cfg})
}

func (h *Handlers) SaveWarehouseConfig(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "business_id es requerido"})
		return
	}
	warehouseID, ok := parseUintParam(c, "warehouse_id")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "warehouse_id invalido"})
		return
	}

	var req domain.SaveConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}
	req.UpdatedBy, req.UpdatedByName = h.resolveActor(c)

	cfg, err := h.uc.SaveWarehouseConfig(c.Request.Context(), businessID, warehouseID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": cfg})
}

func (h *Handlers) DeleteWarehouseConfig(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "business_id es requerido"})
		return
	}
	warehouseID, ok := parseUintParam(c, "warehouse_id")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "warehouse_id invalido"})
		return
	}

	if err := h.uc.RemoveWarehouseConfig(c.Request.Context(), businessID, warehouseID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Configuracion de la bodega eliminada"})
}

func (h *Handlers) SetDefaultWarehouse(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "business_id es requerido"})
		return
	}
	warehouseID, ok := parseUintParam(c, "warehouse_id")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "warehouse_id invalido"})
		return
	}

	if err := h.uc.SetDefaultWarehouse(c.Request.Context(), businessID, warehouseID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Bodega predeterminada actualizada"})
}
