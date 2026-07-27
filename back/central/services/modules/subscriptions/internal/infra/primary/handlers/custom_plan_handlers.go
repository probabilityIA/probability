package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/infra/primary/handlers/request"
	"github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/infra/primary/handlers/response"
)

func (h *Handlers) ListCustomPlans(c *gin.Context) {
	if !h.requireSuperAdmin(c) {
		return
	}

	var businessID *uint
	if param := c.Query("business_id"); param != "" {
		if id, err := strconv.ParseUint(param, 10, 64); err == nil && id > 0 {
			v := uint(id)
			businessID = &v
		}
	}

	types, err := h.uc.ListCustomPlans(c.Request.Context(), businessID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list custom plans"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": response.FromSubscriptionTypes(types)})
}

func (h *Handlers) CreateCustomPlan(c *gin.Context) {
	if !h.requireSuperAdmin(c) {
		return
	}

	var req request.CreateCustomPlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	subType, err := h.uc.CreateCustomPlan(c.Request.Context(), dtos.CreateCustomPlanDTO{
		Name:                 req.Name,
		Code:                 req.Code,
		Description:          req.Description,
		Price:                req.Price,
		BillingPeriod:        req.BillingPeriod,
		ModuleCodes:          req.ModuleCodes,
		MaxEcommerceChannels: req.MaxEcommerceChannels,
		BusinessID:           req.BusinessID,
		Months:               req.Months,
		PaymentReference:     req.PaymentReference,
		Notes:                req.Notes,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": response.FromSubscriptionType(subType)})
}

func (h *Handlers) UpdateCustomPlan(c *gin.Context) {
	if !h.requireSuperAdmin(c) {
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid custom plan id"})
		return
	}

	var req request.UpdateSubscriptionTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	subType, err := h.uc.UpdateSubscriptionType(c.Request.Context(), dtos.UpdateSubscriptionTypeDTO{
		ID:                   uint(id),
		Name:                 req.Name,
		Description:          req.Description,
		Price:                req.Price,
		BillingPeriod:        req.BillingPeriod,
		Active:               req.Active,
		ModuleCodes:          req.ModuleCodes,
		MaxEcommerceChannels: req.MaxEcommerceChannels,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": response.FromSubscriptionType(subType)})
}

func (h *Handlers) DeleteCustomPlan(c *gin.Context) {
	if !h.requireSuperAdmin(c) {
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid custom plan id"})
		return
	}

	if err := h.uc.DeleteSubscriptionType(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete custom plan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "custom plan deleted"})
}
