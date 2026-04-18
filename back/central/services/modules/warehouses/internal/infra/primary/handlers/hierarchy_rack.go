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

func (h *Handlers) CreateRack(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id is required"})
		return
	}

	var req request.CreateRackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	rack, err := h.uc.CreateRack(c.Request.Context(), apprequest.CreateRackDTO{
		AisleID:     req.AisleID,
		BusinessID:  businessID,
		Code:        req.Code,
		Name:        req.Name,
		LevelsCount: req.LevelsCount,
		IsActive:    isActive,
	})
	if err != nil {
		switch {
		case errors.Is(err, domainerrors.ErrAisleNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, domainerrors.ErrDuplicateRackCode):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusCreated, response.RackFromEntity(rack))
}

func (h *Handlers) GetRack(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id is required"})
		return
	}
	rackID, err := strconv.ParseUint(c.Param("rackId"), 10, 64)
	if err != nil || rackID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid rack id"})
		return
	}
	rack, err := h.uc.GetRack(c.Request.Context(), businessID, uint(rackID))
	if err != nil {
		if errors.Is(err, domainerrors.ErrRackNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response.RackFromEntity(rack))
}

func (h *Handlers) ListRacks(c *gin.Context) {
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
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	racks, total, err := h.uc.ListRacks(c.Request.Context(), dtos.ListRacksParams{
		BusinessID: businessID,
		AisleID:    uint(aisleID),
		Page:       page,
		PageSize:   pageSize,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	data := make([]response.RackResponse, len(racks))
	for i := range racks {
		data[i] = response.RackFromEntity(&racks[i])
	}
	if pageSize < 1 {
		pageSize = 10
	}
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
	c.JSON(http.StatusOK, response.RackListResponse{
		Data:       data,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	})
}

func (h *Handlers) UpdateRack(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id is required"})
		return
	}
	rackID, err := strconv.ParseUint(c.Param("rackId"), 10, 64)
	if err != nil || rackID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid rack id"})
		return
	}

	var req request.UpdateRackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rack, err := h.uc.UpdateRack(c.Request.Context(), apprequest.UpdateRackDTO{
		ID:          uint(rackID),
		BusinessID:  businessID,
		Code:        req.Code,
		Name:        req.Name,
		LevelsCount: req.LevelsCount,
		IsActive:    req.IsActive,
	})
	if err != nil {
		switch {
		case errors.Is(err, domainerrors.ErrRackNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, domainerrors.ErrDuplicateRackCode):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, response.RackFromEntity(rack))
}

func (h *Handlers) DeleteRack(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id is required"})
		return
	}
	rackID, err := strconv.ParseUint(c.Param("rackId"), 10, 64)
	if err != nil || rackID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid rack id"})
		return
	}
	if err := h.uc.DeleteRack(c.Request.Context(), businessID, uint(rackID)); err != nil {
		if errors.Is(err, domainerrors.ErrRackNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "rack deleted"})
}
