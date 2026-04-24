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

func (h *Handlers) CreateAisle(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id is required"})
		return
	}

	var req request.CreateAisleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	aisle, err := h.uc.CreateAisle(c.Request.Context(), apprequest.CreateAisleDTO{
		ZoneID:     req.ZoneID,
		BusinessID: businessID,
		Code:       req.Code,
		Name:       req.Name,
		IsActive:   isActive,
	})
	if err != nil {
		switch {
		case errors.Is(err, domainerrors.ErrZoneNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, domainerrors.ErrDuplicateAisleCode):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusCreated, response.AisleFromEntity(aisle))
}

func (h *Handlers) GetAisle(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id is required"})
		return
	}
	aisleID, err := strconv.ParseUint(c.Param("aisleId"), 10, 64)
	if err != nil || aisleID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid aisle id"})
		return
	}
	aisle, err := h.uc.GetAisle(c.Request.Context(), businessID, uint(aisleID))
	if err != nil {
		if errors.Is(err, domainerrors.ErrAisleNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response.AisleFromEntity(aisle))
}

func (h *Handlers) ListAisles(c *gin.Context) {
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
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	aisles, total, err := h.uc.ListAisles(c.Request.Context(), dtos.ListAislesParams{
		BusinessID: businessID,
		ZoneID:     uint(zoneID),
		Page:       page,
		PageSize:   pageSize,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	data := make([]response.AisleResponse, len(aisles))
	for i := range aisles {
		data[i] = response.AisleFromEntity(&aisles[i])
	}
	if pageSize < 1 {
		pageSize = 10
	}
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
	c.JSON(http.StatusOK, response.AisleListResponse{
		Data:       data,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	})
}

func (h *Handlers) UpdateAisle(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id is required"})
		return
	}
	aisleID, err := strconv.ParseUint(c.Param("aisleId"), 10, 64)
	if err != nil || aisleID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid aisle id"})
		return
	}

	var req request.UpdateAisleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	aisle, err := h.uc.UpdateAisle(c.Request.Context(), apprequest.UpdateAisleDTO{
		ID:         uint(aisleID),
		BusinessID: businessID,
		Code:       req.Code,
		Name:       req.Name,
		IsActive:   req.IsActive,
	})
	if err != nil {
		switch {
		case errors.Is(err, domainerrors.ErrAisleNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, domainerrors.ErrDuplicateAisleCode):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, response.AisleFromEntity(aisle))
}

func (h *Handlers) DeleteAisle(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id is required"})
		return
	}
	aisleID, err := strconv.ParseUint(c.Param("aisleId"), 10, 64)
	if err != nil || aisleID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid aisle id"})
		return
	}
	if err := h.uc.DeleteAisle(c.Request.Context(), businessID, uint(aisleID)); err != nil {
		if errors.Is(err, domainerrors.ErrAisleNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "aisle deleted"})
}
