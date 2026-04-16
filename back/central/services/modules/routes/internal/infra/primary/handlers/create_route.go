package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/routes/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/routes/internal/infra/primary/handlers/request"
	"github.com/secamc93/probability/back/central/services/modules/routes/internal/infra/primary/handlers/response"
)

func (h *Handlers) CreateRoute(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id is required"})
		return
	}

	var req request.CreateRouteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	stops := make([]dtos.CreateRouteStopDTO, len(req.Stops))
	for i, s := range req.Stops {
		stops[i] = dtos.CreateRouteStopDTO{
			OrderID:       s.OrderID,
			Address:       s.Address,
			City:          s.City,
			Lat:           s.Lat,
			Lng:           s.Lng,
			CustomerName:  s.CustomerName,
			CustomerPhone: s.CustomerPhone,
			DeliveryNotes: s.DeliveryNotes,
		}
	}

	dto := dtos.CreateRouteDTO{
		BusinessID:        businessID,
		DriverID:          req.DriverID,
		VehicleID:         req.VehicleID,
		Date:              req.Date,
		StartTime:         req.StartTime,
		EndTime:           req.EndTime,
		OriginWarehouseID: req.OriginWarehouseID,
		OriginAddress:     req.OriginAddress,
		OriginLat:         req.OriginLat,
		OriginLng:         req.OriginLng,
		Notes:             req.Notes,
		Stops:             stops,
	}

	route, err := h.uc.CreateRoute(c.Request.Context(), dto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, response.DetailFromEntity(route))
}
