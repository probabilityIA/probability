package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/dtos"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/infra/primary/handlers/response"
)

func (h *Handlers) ListLocations(c *gin.Context) {
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

	params := dtos.ListLocationsParams{
		WarehouseID: uint(warehouseID),
		BusinessID:  businessID,
	}

	locations, err := h.uc.ListLocations(c.Request.Context(), params)
	if err != nil {
		if errors.Is(err, domainerrors.ErrWarehouseNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	data := make([]response.LocationResponse, len(locations))
	for i, loc := range locations {
		data[i] = response.LocationFromEntity(&loc)
	}

	c.JSON(http.StatusOK, data)
}
