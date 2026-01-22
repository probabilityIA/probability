package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetPendingTransactions returns all pending recharge transactions
func (h *WalletHandlers) GetPendingTransactions(c *gin.Context) {
	transactions, err := h.uc.GetPendingTransactions(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, transactions)
}

// ApproveTransaction approves a pending recharge transaction
func (h *WalletHandlers) ApproveTransaction(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Transaction ID is required"})
		return
	}

	if err := h.uc.ApproveTransaction(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transaction approved"})
}

// RejectTransaction rejects a pending recharge transaction
func (h *WalletHandlers) RejectTransaction(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Transaction ID is required"})
		return
	}

	if err := h.uc.RejectTransaction(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transaction rejected"})
}

// GetProcessedTransactions returns all processed (completed/failed) transactions
func (h *WalletHandlers) GetProcessedTransactions(c *gin.Context) {
	transactions, err := h.uc.GetProcessedTransactions(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, transactions)
}

// ManualDebit subtracts balance from a business wallet
func (h *WalletHandlers) ManualDebit(c *gin.Context) {
	var req struct {
		BusinessID uint    `json:"business_id" binding:"required"`
		Amount     float64 `json:"amount" binding:"required"`
		Reference  string  `json:"reference"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.uc.ManualDebit(c.Request.Context(), req.BusinessID, req.Amount, req.Reference); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Balance subtracted successfully"})
}
