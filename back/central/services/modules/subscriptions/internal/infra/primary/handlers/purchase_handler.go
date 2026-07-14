package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
	"github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/domain/dtos"
	errs "github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/infra/primary/handlers/request"
	"github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/infra/primary/handlers/response"
)

func (h *Handlers) PurchaseSubscription(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id is required"})
		return
	}

	userID, _ := middleware.GetUserID(c)

	var req request.PurchaseSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sub, err := h.uc.PurchaseSubscription(c.Request.Context(), dtos.PurchaseSubscriptionDTO{
		BusinessID:         businessID,
		SubscriptionTypeID: req.SubscriptionTypeID,
		Months:             req.Months,
		UserID:             userID,
	})
	if err != nil {
		if errors.Is(err, errs.ErrInsufficientBalance) {
			c.JSON(http.StatusPaymentRequired, gin.H{"error": "insufficient wallet balance"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": response.FromSubscription(sub)})
}
