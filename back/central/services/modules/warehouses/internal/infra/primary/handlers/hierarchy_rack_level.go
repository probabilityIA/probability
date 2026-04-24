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

func (h *Handlers) CreateRackLevel(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id is required"})
		return
	}

	var req request.CreateRackLevelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	level, err := h.uc.CreateRackLevel(c.Request.Context(), apprequest.CreateRackLevelDTO{
		RackID:     req.RackID,
		BusinessID: businessID,
		Code:       req.Code,
		Ordinal:    req.Ordinal,
		IsActive:   isActive,
	})
	if err != nil {
		switch {
		case errors.Is(err, domainerrors.ErrRackNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, domainerrors.ErrDuplicateLevelCode):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusCreated, response.RackLevelFromEntity(level))
}

func (h *Handlers) GetRackLevel(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id is required"})
		return
	}
	levelID, err := strconv.ParseUint(c.Param("levelId"), 10, 64)
	if err != nil || levelID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid level id"})
		return
	}
	level, err := h.uc.GetRackLevel(c.Request.Context(), businessID, uint(levelID))
	if err != nil {
		if errors.Is(err, domainerrors.ErrRackLevelNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response.RackLevelFromEntity(level))
}

func (h *Handlers) ListRackLevels(c *gin.Context) {
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
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	levels, total, err := h.uc.ListRackLevels(c.Request.Context(), dtos.ListRackLevelsParams{
		BusinessID: businessID,
		RackID:     uint(rackID),
		Page:       page,
		PageSize:   pageSize,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	data := make([]response.RackLevelResponse, len(levels))
	for i := range levels {
		data[i] = response.RackLevelFromEntity(&levels[i])
	}
	if pageSize < 1 {
		pageSize = 10
	}
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
	c.JSON(http.StatusOK, response.RackLevelListResponse{
		Data:       data,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	})
}

func (h *Handlers) UpdateRackLevel(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id is required"})
		return
	}
	levelID, err := strconv.ParseUint(c.Param("levelId"), 10, 64)
	if err != nil || levelID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid level id"})
		return
	}

	var req request.UpdateRackLevelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	level, err := h.uc.UpdateRackLevel(c.Request.Context(), apprequest.UpdateRackLevelDTO{
		ID:         uint(levelID),
		BusinessID: businessID,
		Code:       req.Code,
		Ordinal:    req.Ordinal,
		IsActive:   req.IsActive,
	})
	if err != nil {
		switch {
		case errors.Is(err, domainerrors.ErrRackLevelNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, domainerrors.ErrDuplicateLevelCode):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, response.RackLevelFromEntity(level))
}

func (h *Handlers) DeleteRackLevel(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id is required"})
		return
	}
	levelID, err := strconv.ParseUint(c.Param("levelId"), 10, 64)
	if err != nil || levelID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid level id"})
		return
	}
	if err := h.uc.DeleteRackLevel(c.Request.Context(), businessID, uint(levelID)); err != nil {
		if errors.Is(err, domainerrors.ErrRackLevelNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "rack level deleted"})
}
