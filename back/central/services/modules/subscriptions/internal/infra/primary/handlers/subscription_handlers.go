package handlers

import (
	"net/http"
	"strconv"
	"time"

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

	var startDate *time.Time
	if req.StartDate != nil && *req.StartDate != "" {
		parsed, err := time.Parse("2006-01-02", *req.StartDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "fecha de inicio invalida, formato esperado YYYY-MM-DD"})
			return
		}
		startDate = &parsed
	}

	sub, err := h.uc.RegisterPayment(c.Request.Context(), dtos.RegisterPaymentDTO{
		BusinessID:         req.BusinessID,
		SubscriptionTypeID: req.SubscriptionTypeID,
		Months:             req.Months,
		PaymentReference:   req.PaymentReference,
		Notes:              req.Notes,
		StartDate:          startDate,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": response.FromSubscription(sub)})
}

func (h *Handlers) EditSubscriptionDates(c *gin.Context) {
	if !h.requireSuperAdmin(c) {
		return
	}

	var req request.EditSubscriptionDatesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "fecha de inicio invalida, formato esperado YYYY-MM-DD"})
		return
	}
	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "fecha de fin invalida, formato esperado YYYY-MM-DD"})
		return
	}

	sub, err := h.uc.EditSubscriptionDates(c.Request.Context(), dtos.EditSubscriptionDatesDTO{
		BusinessID: req.BusinessID,
		StartDate:  startDate,
		EndDate:    endDate,
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
