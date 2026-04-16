package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/routes/internal/domain/dtos"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/routes/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/routes/internal/infra/primary/handlers/request"
	"github.com/secamc93/probability/back/central/services/modules/routes/internal/infra/primary/handlers/response"
)

func (h *Handlers) AddStop(c *gin.Context) {
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

	var req request.AddStopRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dto := dtos.AddStopDTO{
		RouteID:       uint(routeID),
		BusinessID:    businessID,
		OrderID:       req.OrderID,
		Address:       req.Address,
		City:          req.City,
		Lat:           req.Lat,
		Lng:           req.Lng,
		CustomerName:  req.CustomerName,
		CustomerPhone: req.CustomerPhone,
		DeliveryNotes: req.DeliveryNotes,
	}

	stop, err := h.uc.AddStop(c.Request.Context(), dto)
	if err != nil {
		if errors.Is(err, domainerrors.ErrRouteNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, response.StopFromEntity(stop))
}
