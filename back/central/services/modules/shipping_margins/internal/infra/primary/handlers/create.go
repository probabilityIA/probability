package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/shipping_margins/internal/domain/dtos"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/shipping_margins/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/shipping_margins/internal/infra/primary/handlers/request"
	"github.com/secamc93/probability/back/central/services/modules/shipping_margins/internal/infra/primary/handlers/response"
)

func (h *Handlers) Create(c *gin.Context) {
	if !h.requireSuperAdmin(c) {
		return
	}
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id is required"})
		return
	}
	var req request.CreateShippingMarginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}
	dto := dtos.CreateShippingMarginDTO{
		BusinessID:      businessID,
		CarrierCode:     req.CarrierCode,
		CarrierName:     req.CarrierName,
		MarginAmount:    req.MarginAmount,
		InsuranceMargin: req.InsuranceMargin,
		IsActive:        isActive,
	}
	m, err := h.uc.Create(c.Request.Context(), dto)
	if err != nil {
		switch {
		case errors.Is(err, domainerrors.ErrDuplicateCarrier):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		case errors.Is(err, domainerrors.ErrInvalidCarrierCode), errors.Is(err, domainerrors.ErrInvalidMargin):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusCreated, response.FromEntity(m))
}
