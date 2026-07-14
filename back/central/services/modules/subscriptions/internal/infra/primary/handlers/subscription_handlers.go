package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/infra/primary/handlers/request"
	"github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/infra/primary/handlers/response"
)

func (h *Handlers) GetCurrentSubscription(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id is required"})
		return
	}

	sub, err := h.uc.GetBusinessSubscription(c.Request.Context(), businessID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch subscription"})
		return
	}
	if sub == nil {
		c.JSON(http.StatusOK, gin.H{"data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": response.FromSubscription(sub)})
}

func (h *Handlers) RegisterPayment(c *gin.Context) {
	if !h.requireSuperAdmin(c) {
		return
	}

	var req request.RegisterPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sub, err := h.uc.RegisterPayment(c.Request.Context(), dtos.RegisterPaymentDTO{
		BusinessID:         req.BusinessID,
		SubscriptionTypeID: req.SubscriptionTypeID,
		Months:             req.Months,
		PaymentReference:   req.PaymentReference,
		Notes:              req.Notes,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": response.FromSubscription(sub)})
}

func (h *Handlers) DisableSubscription(c *gin.Context) {
	if !h.requireSuperAdmin(c) {
		return
	}

	businessID, err := strconv.ParseUint(c.Query("business_id"), 10, 64)
	if err != nil || businessID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id invalido"})
		return
	}

	if err := h.uc.DisableSubscription(c.Request.Context(), uint(businessID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to disable subscription"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "subscription disabled"})
}
