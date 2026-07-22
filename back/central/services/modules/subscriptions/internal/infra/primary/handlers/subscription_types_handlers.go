package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/infra/primary/handlers/request"
	"github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/infra/primary/handlers/response"
)

func (h *Handlers) ListSubscriptionTypes(c *gin.Context) {
	activeOnly := c.Query("active_only") == "true"

	types, err := h.uc.ListSubscriptionTypes(c.Request.Context(), activeOnly)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list subscription types"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": response.FromSubscriptionTypes(types)})
}

func (h *Handlers) GetSubscriptionType(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid subscription type id"})
		return
	}

	subType, err := h.uc.GetSubscriptionType(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get subscription type"})
		return
	}
	if subType == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "subscription type not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": response.FromSubscriptionType(subType)})
}

func (h *Handlers) CreateSubscriptionType(c *gin.Context) {
	if !h.requireSuperAdmin(c) {
		return
	}

	var req request.CreateSubscriptionTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	subType, err := h.uc.CreateSubscriptionType(c.Request.Context(), dtos.CreateSubscriptionTypeDTO{
		Name:                 req.Name,
		Code:                 req.Code,
		Description:          req.Description,
		Price:                req.Price,
		BillingPeriod:        req.BillingPeriod,
		ModuleCodes:          req.ModuleCodes,
		MaxEcommerceChannels: req.MaxEcommerceChannels,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": response.FromSubscriptionType(subType)})
}

func (h *Handlers) UpdateSubscriptionType(c *gin.Context) {
	if !h.requireSuperAdmin(c) {
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid subscription type id"})
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

func (h *Handlers) DeleteSubscriptionType(c *gin.Context) {
	if !h.requireSuperAdmin(c) {
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid subscription type id"})
		return
	}

	if err := h.uc.DeleteSubscriptionType(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete subscription type"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "subscription type deleted"})
}
