package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	apprequest "github.com/secamc93/probability/back/central/services/modules/warehouses/internal/app/request"
	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/dtos"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/infra/primary/handlers/request"
	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/infra/primary/handlers/response"
)

func (h *Handlers) CreateZone(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id is required"})
		return
	}

	var req request.CreateZoneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}
	purpose := req.Purpose
	if purpose == "" {
		purpose = "storage"
	}

	zone, err := h.uc.CreateZone(c.Request.Context(), apprequest.CreateZoneDTO{
		WarehouseID: req.WarehouseID,
		BusinessID:  businessID,
		Code:        req.Code,
		Name:        req.Name,
		Purpose:     purpose,
		ColorHex:    req.ColorHex,
		IsActive:    isActive,
	})
	if err != nil {
		switch {
		case errors.Is(err, domainerrors.ErrWarehouseNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, domainerrors.ErrDuplicateZoneCode):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusCreated, response.ZoneFromEntity(zone))
}

func (h *Handlers) GetZone(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id is required"})
		return
	}
	zoneID, err := strconv.ParseUint(c.Param("zoneId"), 10, 64)
	if err != nil || zoneID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid zone id"})
		return
	}
	zone, err := h.uc.GetZone(c.Request.Context(), businessID, uint(zoneID))
	if err != nil {
		if errors.Is(err, domainerrors.ErrZoneNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response.ZoneFromEntity(zone))
}

func (h *Handlers) ListZones(c *gin.Context) {
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

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	activeOnly := c.Query("active_only") == "true"

	zones, total, err := h.uc.ListZones(c.Request.Context(), dtos.ListZonesParams{
		BusinessID:  businessID,
		WarehouseID: uint(warehouseID),
		ActiveOnly:  activeOnly,
		Page:        page,
		PageSize:    pageSize,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	data := make([]response.ZoneResponse, len(zones))
	for i := range zones {
		data[i] = response.ZoneFromEntity(&zones[i])
	}
	if pageSize < 1 {
		pageSize = 10
	}
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
	c.JSON(http.StatusOK, response.ZoneListResponse{
		Data:       data,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	})
}

func (h *Handlers) UpdateZone(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id is required"})
		return
	}
	zoneID, err := strconv.ParseUint(c.Param("zoneId"), 10, 64)
	if err != nil || zoneID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid zone id"})
		return
	}

	var req request.UpdateZoneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	zone, err := h.uc.UpdateZone(c.Request.Context(), apprequest.UpdateZoneDTO{
		ID:         uint(zoneID),
		BusinessID: businessID,
		Code:       req.Code,
		Name:       req.Name,
		Purpose:    req.Purpose,
		ColorHex:   req.ColorHex,
		IsActive:   req.IsActive,
	})
	if err != nil {
		switch {
		case errors.Is(err, domainerrors.ErrZoneNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, domainerrors.ErrDuplicateZoneCode):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, response.ZoneFromEntity(zone))
}

func (h *Handlers) DeleteZone(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id is required"})
		return
	}
	zoneID, err := strconv.ParseUint(c.Param("zoneId"), 10, 64)
	if err != nil || zoneID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid zone id"})
		return
	}
	if err := h.uc.DeleteZone(c.Request.Context(), businessID, uint(zoneID)); err != nil {
		if errors.Is(err, domainerrors.ErrZoneNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "zone deleted"})
}
