package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/dtos"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/infra/primary/handlers/request"
	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/infra/primary/handlers/response"
)

func (h *Handlers) CreateWarehouse(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id is required"})
		return
	}

	var req request.CreateWarehouseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	dto := dtos.CreateWarehouseDTO{
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
		StructureType: req.StructureType,
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

		ShippingPackageStrategy: req.ShippingPackageStrategy,
		StandardBoxes:           toStandardBoxDTOs(req.StandardBoxes),
	}

	warehouse, err := h.uc.CreateWarehouse(c.Request.Context(), dto)
	if err != nil {
		if errors.Is(err, domainerrors.ErrDuplicateCode) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, response.FromEntity(warehouse))
}

func toStandardBoxDTOs(boxes []request.StandardBoxRequest) []dtos.StandardBoxDTO {
	result := make([]dtos.StandardBoxDTO, len(boxes))
	for i, b := range boxes {
		result[i] = dtos.StandardBoxDTO{
			Name:     b.Name,
			Weight:   b.Weight,
			Length:   b.Length,
			Width:    b.Width,
			Height:   b.Height,
			MaxItems: b.MaxItems,
		}
	}
	return result
}
