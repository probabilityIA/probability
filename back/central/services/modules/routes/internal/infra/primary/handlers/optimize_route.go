package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/routes/internal/domain/dtos"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/routes/internal/domain/errors"
)

func (h *Handlers) OptimizeRoute(c *gin.Context) {
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

	result, err := h.uc.OptimizeRoute(c.Request.Context(), dtos.OptimizeRouteDTO{
		RouteID:    uint(routeID),
		BusinessID: businessID,
	})
	if err != nil {
		switch {
		case errors.Is(err, domainerrors.ErrRouteNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, domainerrors.ErrRouteNotEditable),
			errors.Is(err, domainerrors.ErrNotEnoughStops),
			errors.Is(err, domainerrors.ErrOriginMissing):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		case errors.Is(err, domainerrors.ErrOptimizerNotConfigured):
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":           "ruta optimizada",
		"stop_ids":          result.StopIDs,
		"total_distance_km": result.DistanceKm,
		"total_duration_min": result.DurationMin,
		"stops_optimized":   result.StopsOptimized,
		"stops_sin_coords":  result.StopsSinCoords,
	})
}
