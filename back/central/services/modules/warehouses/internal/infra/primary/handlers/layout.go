package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/infra/primary/handlers/request"
	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/infra/primary/handlers/response"
)

func (h *Handlers) GetLayout(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id is required"})
		return
	}
	warehouseID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || warehouseID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid warehouse id"})
		return
	}

	layout, err := h.uc.GetLayout(c.Request.Context(), businessID, uint(warehouseID))
	if err != nil {
		if errors.Is(err, domainerrors.ErrWarehouseNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response.LayoutFromEntity(layout))
}

func (h *Handlers) SaveLayout(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id is required"})
		return
	}
	warehouseID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || warehouseID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid warehouse id"})
		return
	}

	var req request.SaveLayoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	nodes := make([]entities.LayoutNode, len(req.Nodes))
	for i, n := range req.Nodes {
		nodes[i] = entities.LayoutNode{
			NodeID:   n.NodeID,
			RefType:  n.RefType,
			RefID:    n.RefID,
			X:        n.X,
			Y:        n.Y,
			Width:    n.Width,
			Height:   n.Height,
			Rotation: n.Rotation,
			Color:    n.Color,
			Label:    n.Label,
		}
	}

	layout, err := h.uc.SaveLayout(c.Request.Context(), dtos.SaveLayoutDTO{
		WarehouseID:  uint(warehouseID),
		BusinessID:   businessID,
		CanvasWidth:  req.CanvasWidth,
		CanvasHeight: req.CanvasHeight,
		GridSize:     req.GridSize,
		Nodes:        nodes,
	})
	if err != nil {
		if errors.Is(err, domainerrors.ErrWarehouseNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response.LayoutFromEntity(layout))
}
