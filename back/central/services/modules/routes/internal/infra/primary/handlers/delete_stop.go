package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/routes/internal/domain/errors"
)

func (h *Handlers) DeleteStop(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id is required"})
		return
	}

	routeID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || routeID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid route id"})
		return
	}

	stopID, err := strconv.ParseUint(c.Param("stopId"), 10, 64)
	if err != nil || stopID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid stop id"})
		return
	}

	if err := h.uc.DeleteStop(c.Request.Context(), businessID, uint(routeID), uint(stopID)); err != nil {
		if errors.Is(err, domainerrors.ErrRouteNotFound) || errors.Is(err, domainerrors.ErrStopNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "stop deleted successfully"})
}
