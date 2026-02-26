package handlers

import (
	"fmt"
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

// ClearRechargeHistory deletes all recharge transactions for a business
func (h *WalletHandlers) ClearRechargeHistory(c *gin.Context) {
	businessIDStr := c.Param("business_id")
	if businessIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Business ID is required"})
		return
	}

	var businessID uint
	if _, err := fmt.Sscanf(businessIDStr, "%d", &businessID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid business ID"})
		return
	}

	if err := h.uc.ClearRechargeHistory(c.Request.Context(), businessID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Recharge history cleared successfully"})
}

// DebitForGuide debits from the current user's wallet for guide generation
func (h *WalletHandlers) DebitForGuide(c *gin.Context) {
	var req struct {
		Amount         float64 `json:"amount" binding:"required"`
		TrackingNumber string  `json:"tracking_number" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Amount must be greater than 0"})
		return
	}

	// Get business ID from JWT claims (assuming middleware sets it)
	businessID, exists := c.Get("business_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Business ID not found in token"})
		return
	}

	ref := fmt.Sprintf("Guide generation: %s", req.TrackingNumber)
	if err := h.uc.ManualDebit(c.Request.Context(), businessID.(uint), req.Amount, ref); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Wallet debited successfully", "tracking_number": req.TrackingNumber})
}
