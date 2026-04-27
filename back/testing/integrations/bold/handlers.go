package bold

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const mockSecretKey = "test-secret"

func (b *Bold) handleCreateLink(c *gin.Context) {
	apiKey := strings.TrimPrefix(c.GetHeader("Authorization"), "x-api-key ")
	if apiKey == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"errors": []string{"missing api key"}})
		return
	}

	var req struct {
		AmountType  string  `json:"amount_type"`
		Amount      any     `json:"amount"`
		Currency    string  `json:"currency"`
		Description string  `json:"description"`
		Reference   string  `json:"merchant_reference"`
		TotalAmount float64 `json:"total_amount"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": []string{"invalid body"}})
		return
	}

	linkID := "BOLD-" + randHex(6)
	state := &linkState{
		ID:        linkID,
		Reference: req.Reference,
		Amount:    req.TotalAmount,
		Currency:  req.Currency,
		CreatedAt: time.Now(),
		Status:    "PENDING",
		PaymentID: "PAY-" + randHex(8),
	}

	b.mu.Lock()
	b.links[linkID] = state
	b.mu.Unlock()

	checkoutURL := fmt.Sprintf("http://localhost:%s/checkout/%s", b.port, linkID)

	c.JSON(http.StatusOK, gin.H{
		"payload": gin.H{
			"payment_link": linkID,
			"url":          checkoutURL,
			"expires_at":   time.Now().Add(24 * time.Hour).Format(time.RFC3339),
		},
	})
}

func (b *Bold) handleGetLink(c *gin.Context) {
	linkID := c.Param("id")
	b.mu.RLock()
	st, ok := b.links[linkID]
	b.mu.RUnlock()
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"errors": []string{"link not found"}})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":             linkID,
		"payment_status": st.Status,
		"amount":         st.Amount,
		"currency":       st.Currency,
		"payment_method": "CREDIT_CARD",
	})
}

func (b *Bold) handleSimulateApproved(c *gin.Context) {
	b.simulateEvent(c, "SALE_APPROVED", "APPROVED")
}

func (b *Bold) handleSimulateRejected(c *gin.Context) {
	b.simulateEvent(c, "SALE_REJECTED", "REJECTED")
}

func (b *Bold) handleSimulateVoidApproved(c *gin.Context) {
	b.simulateEvent(c, "VOID_APPROVED", "VOIDED")
}

func (b *Bold) simulateEvent(c *gin.Context, eventType, newStatus string) {
	linkID := c.Param("id")
	b.mu.Lock()
	st, ok := b.links[linkID]
	if !ok {
		b.mu.Unlock()
		c.JSON(http.StatusNotFound, gin.H{"error": "link not found"})
		return
	}
	st.Status = newStatus
	state := *st
	b.mu.Unlock()

	if b.webhookTarget == "" {
		c.JSON(http.StatusOK, gin.H{"sent": false, "reason": "webhook target not configured"})
		return
	}

	payload := gin.H{
		"id":      uuid.New().String(),
		"type":    eventType,
		"subject": linkID,
		"source":  "/payments/links",
		"time":    time.Now().UnixNano(),
		"data": gin.H{
			"payment_id":         state.PaymentID,
			"merchant_id":        "MOCK_MERCHANT",
			"amount":             state.Amount,
			"currency":           state.Currency,
			"payment_method":     "CREDIT_CARD",
			"merchant_reference": state.Reference,
			"payer_email":        "test@probability.com",
			"created_at":         state.CreatedAt.Format(time.RFC3339),
		},
	}
	body, _ := json.Marshal(payload)

	signature := computeBoldSignature(body, mockSecretKey)

	req, _ := http.NewRequest(http.MethodPost, b.webhookTarget, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-bold-signature", signature)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		b.logger.Error().Msgf("bold mock: failed to deliver webhook: %s", err.Error())
		c.JSON(http.StatusOK, gin.H{"sent": false, "error": err.Error()})
		return
	}
	defer resp.Body.Close()

	st.WebhookSent = resp.StatusCode == http.StatusOK
	c.JSON(http.StatusOK, gin.H{
		"sent":            true,
		"webhook_status":  resp.StatusCode,
		"event_type":      eventType,
		"new_link_status": newStatus,
		"payment_id":      state.PaymentID,
	})
}

func (b *Bold) handleListLinks(c *gin.Context) {
	b.mu.RLock()
	out := make([]gin.H, 0, len(b.links))
	for _, st := range b.links {
		out = append(out, gin.H{
			"id":           st.ID,
			"reference":    st.Reference,
			"amount":       st.Amount,
			"status":       st.Status,
			"payment_id":   st.PaymentID,
			"created_at":   st.CreatedAt.Format(time.RFC3339),
			"webhook_sent": st.WebhookSent,
		})
	}
	b.mu.RUnlock()
	c.JSON(http.StatusOK, gin.H{"links": out})
}

func computeBoldSignature(body []byte, secret string) string {
	bodyB64 := base64.StdEncoding.EncodeToString(body)
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(bodyB64))
	return hex.EncodeToString(mac.Sum(nil))
}
