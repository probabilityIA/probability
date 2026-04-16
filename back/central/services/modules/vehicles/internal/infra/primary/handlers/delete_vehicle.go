package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/vehicles/internal/domain/errors"
)

func (h *Handlers) DeleteVehicle(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id is required"})
		return
	}

	vehicleID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || vehicleID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid vehicle id"})
		return
	}

	if err := h.uc.DeleteVehicle(c.Request.Context(), businessID, uint(vehicleID)); err != nil {
		if errors.Is(err, domainerrors.ErrVehicleNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "vehicle deleted successfully"})
}
