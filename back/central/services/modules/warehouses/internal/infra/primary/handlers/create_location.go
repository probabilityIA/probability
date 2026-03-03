package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/dtos"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/infra/primary/handlers/request"
	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/infra/primary/handlers/response"
)

func (h *Handlers) CreateLocation(c *gin.Context) {
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

	var req request.CreateLocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	locType := "storage"
	if req.Type != "" {
		locType = req.Type
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	dto := dtos.CreateLocationDTO{
		WarehouseID:   uint(warehouseID),
		BusinessID:    businessID,
		Name:          req.Name,
		Code:          req.Code,
		Type:          locType,
		IsActive:      isActive,
		IsFulfillment: req.IsFulfillment,
		Capacity:      req.Capacity,
	}

	location, err := h.uc.CreateLocation(c.Request.Context(), dto)
	if err != nil {
		if errors.Is(err, domainerrors.ErrWarehouseNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, domainerrors.ErrDuplicateLocCode) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, response.LocationFromEntity(location))
}
