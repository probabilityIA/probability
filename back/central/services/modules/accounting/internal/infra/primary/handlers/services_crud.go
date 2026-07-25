package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/accounting/internal/domain/dtos"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/accounting/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/accounting/internal/infra/primary/handlers/request"
	"github.com/secamc93/probability/back/central/services/modules/accounting/internal/infra/primary/handlers/response"
)

func (h *Handlers) serviceError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domainerrors.ErrServiceNotFound):
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": err.Error()})
	case errors.Is(err, domainerrors.ErrDuplicateCode):
		c.JSON(http.StatusConflict, gin.H{"success": false, "error": err.Error()})
	case errors.Is(err, domainerrors.ErrCodeRequired), errors.Is(err, domainerrors.ErrNameRequired), errors.Is(err, domainerrors.ErrInvalidAmount):
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
	}
}

func (h *Handlers) ListServices(c *gin.Context) {
	if !h.requireSuperAdmin(c) {
		return
	}
	items, err := h.uc.ListServices(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	data := make([]response.ServiceResponse, len(items))
	for i := range items {
		data[i] = response.FromService(&items[i])
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": data})
}

func (h *Handlers) CreateService(c *gin.Context) {
	if !h.requireSuperAdmin(c) {
		return
	}
	var req request.SaveServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}
	service, err := h.uc.CreateService(c.Request.Context(), dtos.SaveServiceDTO{
		Code:       req.Code,
		Name:       req.Name,
		UnspscCode: req.UnspscCode,
		UnitPrice:  req.UnitPrice,
	})
	if err != nil {
		h.serviceError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "data": response.FromService(service)})
}

func (h *Handlers) UpdateService(c *gin.Context) {
	if !h.requireSuperAdmin(c) {
		return
	}
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "id invalido"})
		return
	}
	var req request.SaveServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}
	service, err := h.uc.UpdateService(c.Request.Context(), dtos.SaveServiceDTO{
		ID:         uint(id),
		Code:       req.Code,
		Name:       req.Name,
		UnspscCode: req.UnspscCode,
		UnitPrice:  req.UnitPrice,
		IsActive:   isActive,
	})
	if err != nil {
		h.serviceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": response.FromService(service)})
}

func (h *Handlers) DeleteService(c *gin.Context) {
	if !h.requireSuperAdmin(c) {
		return
	}
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "id invalido"})
		return
	}
	if err := h.uc.DeleteService(c.Request.Context(), uint(id)); err != nil {
		h.serviceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}
