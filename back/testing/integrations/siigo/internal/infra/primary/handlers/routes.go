package handlers

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/testing/integrations/siigo/internal/domain"
)

func (h *Handler) RegisterRoutes(router *gin.Engine) {
	router.GET("/health", h.handleHealth)
	router.POST("/v1/auth", h.handleAuth)
	router.POST("/v1/customers", h.handleCreateCustomer)
	router.GET("/v1/customers", h.handleGetCustomer)
	router.POST("/v1/invoices", h.handleCreateInvoice)
	router.GET("/v1/invoices/:id", h.handleGetInvoice)
	router.POST("/v1/journals", h.handleCreateJournal)
}

func (h *Handler) handleHealth(c *gin.Context) {
	c.JSON(200, gin.H{"status": "ok", "service": "siigo-mock"})
}

func (h *Handler) handleAuth(c *gin.Context) {
	var req struct {
		Username  string `json:"username"`
		AccessKey string `json:"access_key"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"Status": 400, "Errors": []gin.H{{"Code": "invalid_body", "Message": "invalid request body"}}})
		return
	}

	partnerID := c.GetHeader("Partner-Id")
	token, err := h.apiSimulator.HandleAuth(req.Username, req.AccessKey, partnerID)
	if err != nil {
		c.JSON(401, gin.H{"Status": 401, "Errors": []gin.H{{"Code": "unauthorized", "Message": err.Error()}}})
		return
	}

	c.JSON(200, gin.H{
		"access_token": token,
		"expires_in":   86400,
		"token_type":   "Bearer",
		"scope":        "siigoAPI",
	})
}

func (h *Handler) handleCreateCustomer(c *gin.Context) {
	if !h.requireAuth(c) {
		return
	}

	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"Status": 400, "Errors": []gin.H{{"Code": "invalid_body", "Message": "invalid request body"}}})
		return
	}

	customer, err := h.apiSimulator.HandleCreateCustomer(body)
	if err != nil {
		c.JSON(400, gin.H{"Status": 400, "Errors": []gin.H{{"Code": "invalid_data", "Message": err.Error()}}})
		return
	}

	c.JSON(201, gin.H{
		"id":             customer.ID,
		"identification": customer.Identification,
		"name":           []string{customer.Name},
		"email":          customer.Email,
		"phone":          customer.Phone,
	})
}

func (h *Handler) handleGetCustomer(c *gin.Context) {
	if !h.requireAuth(c) {
		return
	}

	identification := c.Query("identification")
	if identification == "" {
		c.JSON(200, gin.H{"results": []interface{}{}, "pagination": gin.H{"total_results": 0}})
		return
	}

	customer, ok := h.apiSimulator.HandleGetCustomer(identification)
	if !ok {
		c.JSON(200, gin.H{"results": []interface{}{}, "pagination": gin.H{"total_results": 0}})
		return
	}

	c.JSON(200, gin.H{
		"results": []gin.H{{
			"id":             customer.ID,
			"identification": customer.Identification,
			"name":           []string{customer.Name},
			"email":          customer.Email,
		}},
		"pagination": gin.H{"total_results": 1},
	})
}

func (h *Handler) handleCreateInvoice(c *gin.Context) {
	if !h.requireAuth(c) {
		return
	}

	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"Status": 400, "Errors": []gin.H{{"Code": "invalid_body", "Message": "invalid request body"}}})
		return
	}

	invoice, err := h.apiSimulator.HandleCreateInvoice(body)
	if err != nil {
		c.JSON(400, gin.H{"Status": 400, "Errors": []gin.H{{"Code": "invalid_data", "Message": err.Error()}}})
		return
	}

	c.JSON(201, invoicePayload(invoice))
}

func (h *Handler) handleGetInvoice(c *gin.Context) {
	if !h.requireAuth(c) {
		return
	}

	id := c.Param("id")
	invoice, ok := h.apiSimulator.HandleGetInvoice(id)
	if !ok {
		c.JSON(404, gin.H{"Status": 404, "Errors": []gin.H{{"Code": "not_found", "Message": "invoice not found"}}})
		return
	}

	c.JSON(200, invoicePayload(invoice))
}

func (h *Handler) handleCreateJournal(c *gin.Context) {
	if !h.requireAuth(c) {
		return
	}

	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"Status": 400, "Errors": []gin.H{{"Code": "invalid_body", "Message": "invalid request body"}}})
		return
	}

	journal, err := h.apiSimulator.HandleCreateJournal(body)
	if err != nil {
		c.JSON(400, gin.H{"Status": 400, "Errors": []gin.H{{"Code": "invalid_data", "Message": err.Error()}}})
		return
	}

	c.JSON(201, gin.H{
		"id":       journal.ID,
		"date":     journal.Date,
		"document": gin.H{"id": journal.DocumentID},
		"items":    journal.Items,
		"total":    journal.Total,
	})
}

func (h *Handler) requireAuth(c *gin.Context) bool {
	auth := c.GetHeader("Authorization")
	if auth == "" || !strings.HasPrefix(strings.ToLower(auth), "bearer ") {
		c.JSON(401, gin.H{"Status": 401, "Errors": []gin.H{{"Code": "unauthorized", "Message": "missing bearer token"}}})
		return false
	}
	return true
}

func invoicePayload(invoice *domain.Invoice) gin.H {
	items := make([]gin.H, 0, len(invoice.Items))
	for _, it := range invoice.Items {
		items = append(items, gin.H{
			"code":        it.Code,
			"description": it.Description,
			"quantity":    it.Quantity,
			"price":       it.Price,
			"total":       it.Total,
		})
	}

	return gin.H{
		"id":     invoice.ID,
		"name":   invoice.Name,
		"date":   invoice.Date,
		"number": invoice.Number,
		"document": gin.H{
			"prefix": invoice.Prefix,
			"number": invoice.Number,
		},
		"customer": gin.H{
			"id":             invoice.CustomerID,
			"identification": invoice.CustomerNIT,
		},
		"items":      items,
		"total":      invoice.Total,
		"cufe":       invoice.CUFE,
		"public_url": invoice.PublicURL,
	}
}
