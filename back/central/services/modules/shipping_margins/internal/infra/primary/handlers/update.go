package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/shipping_margins/internal/domain/dtos"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/shipping_margins/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/shipping_margins/internal/infra/primary/handlers/request"
	"github.com/secamc93/probability/back/central/services/modules/shipping_margins/internal/infra/primary/handlers/response"
)

func (h *Handlers) Update(c *gin.Context) {
	if !h.requireSuperAdmin(c) {
		return
	}
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id is required"})
		return
	}
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req request.UpdateShippingMarginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}
	dto := dtos.UpdateShippingMarginDTO{
		ID:              uint(id),
		BusinessID:      businessID,
		CarrierName:     req.CarrierName,
		MarginAmount:    req.MarginAmount,
		InsuranceMargin: req.InsuranceMargin,
		IsActive:        isActive,
	}
	m, err := h.uc.Update(c.Request.Context(), dto)
	if err != nil {
		switch {
		case errors.Is(err, domainerrors.ErrShippingMarginNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, domainerrors.ErrInvalidMargin):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, response.FromEntity(m))
}
