package handlers

import (
	stderrors "errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/pay/internal/domain/dtos"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/pay/internal/domain/errors"
)

func (h *walletHandler) BoldGenerateSignature(c *gin.Context) {
	businessID, ok := resolveBusinessID(c)
	if !ok {
		return
	}

	amountStr := c.Query("amount")
	if amountStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "amount is required"})
		return
	}
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil || amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "amount must be a positive number"})
		return
	}

	currency := c.Query("currency")
	if currency == "" {
		currency = "COP"
	}

	resp, err := h.walletUC.BoldGenerateSignature(c.Request.Context(), businessID, amount, currency)
	if err != nil {
		status, body := mapBoldError(err)
		h.log.Error(c.Request.Context()).Err(err).Int("status", status).Msg("BoldGenerateSignature failed")
		c.JSON(status, body)
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

func (h *walletHandler) GetBoldStatus(c *gin.Context) {
	boldOrderID := c.Param("id")
	if boldOrderID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "order id is required"})
		return
	}

	resp, err := h.walletUC.GetBoldStatus(c.Request.Context(), boldOrderID)
	if err != nil {
		status, body := mapBoldError(err)
		h.log.Error(c.Request.Context()).Err(err).Str("bold_order_id", boldOrderID).Int("status", status).Msg("GetBoldStatus failed")
		c.JSON(status, body)
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

func (h *walletHandler) BoldSimulatePayment(c *gin.Context) {
	businessID, ok := resolveBusinessID(c)
	if !ok {
		return
	}

	var req struct {
		OrderID string  `json:"order_id" binding:"required"`
		Amount  float64 `json:"amount" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "order_id and amount are required"})
		return
	}
	if req.Amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "amount must be greater than zero"})
		return
	}

	resp, err := h.walletUC.BoldSimulatePayment(c.Request.Context(), &dtos.BoldSimulateDTO{
		BusinessID: businessID,
		OrderID:    req.OrderID,
		Amount:     req.Amount,
	})
	if err != nil {
		status, body := mapBoldError(err)
		h.log.Error(c.Request.Context()).Err(err).Int("status", status).Msg("BoldSimulatePayment failed")
		c.JSON(status, body)
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

func mapBoldError(err error) (int, gin.H) {
	switch {
	case stderrors.Is(err, domainerrors.ErrBoldOrderNotFound):
		return http.StatusNotFound, gin.H{"success": false, "code": "BOLD_ORDER_NOT_FOUND", "message": "bold order not found"}
	case stderrors.Is(err, domainerrors.ErrBoldUnauthorized):
		return http.StatusBadGateway, gin.H{"success": false, "code": "BOLD_UNAUTHORIZED", "message": "bold credentials invalid"}
	case stderrors.Is(err, domainerrors.ErrBoldConfigNotFound):
		return http.StatusServiceUnavailable, gin.H{"success": false, "code": "BOLD_CONFIG_NOT_FOUND", "message": "bold integration not configured"}
	case stderrors.Is(err, domainerrors.ErrBoldCredentialsMissing):
		return http.StatusServiceUnavailable, gin.H{"success": false, "code": "BOLD_CREDENTIALS_MISSING", "message": "bold credentials not configured"}
	case stderrors.Is(err, domainerrors.ErrBoldUpstreamUnavailable):
		return http.StatusBadGateway, gin.H{"success": false, "code": "BOLD_UPSTREAM_UNAVAILABLE", "message": "bold api unavailable"}
	default:
		return http.StatusInternalServerError, gin.H{"success": false, "code": "INTERNAL", "message": err.Error()}
	}
}
