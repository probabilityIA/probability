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

func (h *Handlers) UpdateWarehouse(c *gin.Context) {
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

	var req request.UpdateWarehouseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	dto := dtos.UpdateWarehouseDTO{
		ID:            uint(warehouseID),
		BusinessID:    businessID,
		Name:          req.Name,
		Code:          req.Code,
		Address:       req.Address,
		City:          req.City,
		State:         req.State,
		Country:       req.Country,
		ZipCode:       req.ZipCode,
		Phone:         req.Phone,
		ContactName:   req.ContactName,
		ContactEmail:  req.ContactEmail,
		IsActive:      isActive,
		IsDefault:     req.IsDefault,
		IsFulfillment: req.IsFulfillment,
		Company:       req.Company,
		FirstName:     req.FirstName,
		LastName:      req.LastName,
		Email:         req.Email,
		Suburb:        req.Suburb,
		CityDaneCode:  req.CityDaneCode,
		PostalCode:    req.PostalCode,
		Street:        req.Street,
		Latitude:      req.Latitude,
		Longitude:     req.Longitude,
	}

	warehouse, err := h.uc.UpdateWarehouse(c.Request.Context(), dto)
	if err != nil {
		if errors.Is(err, domainerrors.ErrWarehouseNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, domainerrors.ErrDuplicateCode) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.FromEntity(warehouse))
}
