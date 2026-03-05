package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/routes/internal/domain/dtos"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/routes/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/routes/internal/infra/primary/handlers/request"
)

func (h *Handlers) UpdateStopStatus(c *gin.Context) {
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

	var req request.UpdateStopStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dto := dtos.UpdateStopStatusDTO{
		ID:            uint(stopID),
		RouteID:       uint(routeID),
		BusinessID:    businessID,
		Status:        req.Status,
		FailureReason: req.FailureReason,
		SignatureURL:  req.SignatureURL,
		PhotoURL:      req.PhotoURL,
	}

	if err := h.uc.UpdateStopStatus(c.Request.Context(), dto); err != nil {
		if errors.Is(err, domainerrors.ErrRouteNotFound) || errors.Is(err, domainerrors.ErrStopNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, domainerrors.ErrRouteNotInProgress) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "stop status updated"})
}
