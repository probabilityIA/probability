package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/vehicles/internal/domain/dtos"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/vehicles/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/vehicles/internal/infra/primary/handlers/request"
	"github.com/secamc93/probability/back/central/services/modules/vehicles/internal/infra/primary/handlers/response"
)

func (h *Handlers) UpdateVehicle(c *gin.Context) {
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

	var req request.UpdateVehicleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dto := dtos.UpdateVehicleDTO{
		ID:                 uint(vehicleID),
		BusinessID:         businessID,
		Type:               req.Type,
		LicensePlate:       req.LicensePlate,
		Brand:              req.Brand,
		VehicleModel:       req.VehicleModel,
		Year:               req.Year,
		Color:              req.Color,
		Status:             req.Status,
		WeightCapacityKg:   req.WeightCapacityKg,
		VolumeCapacityM3:   req.VolumeCapacityM3,
		PhotoURL:           req.PhotoURL,
		InsuranceExpiry:    req.InsuranceExpiry,
		RegistrationExpiry: req.RegistrationExpiry,
	}

	vehicle, err := h.uc.UpdateVehicle(c.Request.Context(), dto)
	if err != nil {
		if errors.Is(err, domainerrors.ErrVehicleNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, domainerrors.ErrDuplicateLicensePlate) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.FromEntity(vehicle))
}
