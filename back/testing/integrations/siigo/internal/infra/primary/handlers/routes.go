package handlers

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/testing/integrations/siigo/internal/domain"
)

func (h *Handler) RegisterRoutes(router *gin.Engine) {
	router.GET("/health", h.handleHealth)
	router.POST("/auth", h.handleAuth)
	router.POST("/v1/auth", h.handleAuth)
	router.POST("/v1/customers", h.handleCreateCustomer)
	router.GET("/v1/customers", h.handleGetCustomer)
	router.POST("/v1/invoices", h.handleCreateInvoice)
	router.GET("/v1/invoices", h.handleListInvoices)
	router.GET("/v1/invoices/:id", h.handleGetInvoice)
	router.POST("/v1/invoices/:id/annul", h.handleAnnulInvoice)
	router.GET("/v1/invoices/:id/stamp/errors", h.handleStampErrors)
	router.GET("/v1/products", h.handleListProducts)
	router.GET("/v1/warehouses", h.handleListWarehouses)
	router.GET("/v1/payment-types", h.handleListPaymentTypes)
	router.POST("/v1/vouchers", h.handleCreateVoucher)
	router.POST("/v1/credit-notes", h.handleCreateCreditNote)
	router.POST("/v1/journals", h.handleCreateJournal)
	router.GET("/v1/webhooks", h.handleListWebhooks)
	router.POST("/v1/webhooks", h.handleCreateWebhook)
	router.DELETE("/v1/webhooks/:id", h.handleDeleteWebhook)
}

var validWebhookTopics = map[string]bool{
	"public.siigoapi.products.stock.update": true,
	"public.siigoapi.products.create":       true,
	"public.siigoapi.products.update":       true,
}

func webhookPayload(w *domain.Webhook) gin.H {
	return gin.H{
		"id":             w.ID,
		"application_id": w.ApplicationID,
		"url":            w.URL,
		"topic":          w.Topic,
		"company_key":    w.CompanyKey,
		"active":         w.Active,
		"created_at":     w.CreatedAt,
	}
}

func (h *Handler) handleListWebhooks(c *gin.Context) {
	if !h.requireAuth(c) {
		return
	}
	webhooks := h.apiSimulator.HandleListWebhooks()
	out := make([]gin.H, 0, len(webhooks))
	for _, w := range webhooks {
		out = append(out, webhookPayload(w))
	}
	c.JSON(200, out)
}

func (h *Handler) handleCreateWebhook(c *gin.Context) {
	if !h.requireAuth(c) {
		return
	}
	var req struct {
		ApplicationID string `json:"application_id"`
		URL           string `json:"url"`
		Topic         string `json:"topic"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"Status": 400, "Errors": []gin.H{{"Code": "invalid_body", "Message": "invalid request body"}}})
		return
	}
	if req.URL == "" || req.ApplicationID == "" {
		c.JSON(400, gin.H{"Status": 400, "Errors": []gin.H{{"Code": "invalid_data", "Message": "application_id and url are required"}}})
		return
	}
	if !validWebhookTopics[req.Topic] {
		c.JSON(400, gin.H{"Status": 400, "Errors": []gin.H{{"Code": "invalid_topic", "Message": "The topic doesn't exist"}}})
		return
	}
	w := h.apiSimulator.HandleCreateWebhook(req.ApplicationID, req.URL, req.Topic)
	c.JSON(201, webhookPayload(w))
}

func (h *Handler) handleDeleteWebhook(c *gin.Context) {
	if !h.requireAuth(c) {
		return
	}
	id := c.Param("id")
	if !h.apiSimulator.HandleDeleteWebhook(id) {
		c.JSON(404, gin.H{"Status": 404, "Errors": []gin.H{{"Code": "not_found", "Message": "webhook not found"}}})
		return
	}
	c.JSON(200, gin.H{"deleted": true})
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

func (h *Handler) handleListInvoices(c *gin.Context) {
	if !h.requireAuth(c) {
		return
	}

	name := c.Query("name")
	invoices := h.apiSimulator.HandleListInvoices()
	results := make([]gin.H, 0, len(invoices))
	for _, inv := range invoices {
		if name != "" && inv.Name != name {
			continue
		}
		results = append(results, invoicePayload(inv))
	}

	c.JSON(200, gin.H{
		"results": results,
		"pagination": gin.H{
			"page":          1,
			"page_size":     len(results),
			"total_results": len(results),
		},
	})
}

func (h *Handler) handleAnnulInvoice(c *gin.Context) {
	if !h.requireAuth(c) {
		return
	}

	id := c.Param("id")
	invoice, err := h.apiSimulator.HandleAnnulInvoice(id)
	if err != nil {
		switch err.Error() {
		case "not_found":
			c.JSON(404, gin.H{"Status": 404, "Errors": []gin.H{{"Code": "not_found", "Message": "invoice not found"}}})
		case "annul_not_allowed":
			c.JSON(409, gin.H{"Status": 409, "Errors": []gin.H{{"Code": "annul_not_allowed", "Message": "invoice already annulled or has associated documents"}}})
		default:
			c.JSON(400, gin.H{"Status": 400, "Errors": []gin.H{{"Code": "invalid_data", "Message": err.Error()}}})
		}
		return
	}

	c.JSON(200, invoicePayload(invoice))
}

func (h *Handler) handleStampErrors(c *gin.Context) {
	if !h.requireAuth(c) {
		return
	}

	id := c.Param("id")
	stampErrors, ok := h.apiSimulator.HandleGetStampErrors(id)
	if !ok {
		c.JSON(404, gin.H{"Status": 404, "Errors": []gin.H{{"Code": "not_found", "Message": "invoice not found"}}})
		return
	}

	errs := make([]gin.H, 0, len(stampErrors))
	for _, e := range stampErrors {
		errs = append(errs, gin.H{"Code": e.Code, "Message": e.Message})
	}
	c.JSON(200, gin.H{"Errors": errs})
}

func (h *Handler) handleListProducts(c *gin.Context) {
	if !h.requireAuth(c) {
		return
	}

	products := h.apiSimulator.HandleListProducts()
	results := make([]gin.H, 0, len(products))
	for _, p := range products {
		warehouses := make([]gin.H, 0, len(p.Warehouses))
		for _, w := range p.Warehouses {
			warehouses = append(warehouses, gin.H{
				"id":       w.ID,
				"name":     w.Name,
				"quantity": w.Quantity,
			})
		}
		results = append(results, gin.H{
			"id":                 p.ID,
			"code":               p.Code,
			"name":               p.Name,
			"description":        p.Description,
			"stock_control":      p.StockControl,
			"available_quantity": p.AvailableQuantity,
			"warehouses":         warehouses,
			"prices": []gin.H{{
				"price_list": []gin.H{{"position": 1, "value": p.Price}},
			}},
		})
	}

	c.JSON(200, gin.H{
		"results": results,
		"pagination": gin.H{
			"page":          1,
			"page_size":     len(results),
			"total_results": len(results),
		},
	})
}

func (h *Handler) handleListWarehouses(c *gin.Context) {
	if !h.requireAuth(c) {
		return
	}

	warehouses := h.apiSimulator.HandleListWarehouses()
	results := make([]gin.H, 0, len(warehouses))
	for _, w := range warehouses {
		results = append(results, gin.H{
			"id":     w.ID,
			"name":   w.Name,
			"active": true,
		})
	}

	c.JSON(200, gin.H{"results": results})
}

func (h *Handler) handleListPaymentTypes(c *gin.Context) {
	if !h.requireAuth(c) {
		return
	}

	paymentTypes := h.apiSimulator.HandleListPaymentTypes()
	results := make([]gin.H, 0, len(paymentTypes))
	for _, pt := range paymentTypes {
		results = append(results, gin.H{
			"id":     pt.ID,
			"name":   pt.Name,
			"type":   pt.Type,
			"active": true,
		})
	}

	c.JSON(200, results)
}

func (h *Handler) handleCreateCreditNote(c *gin.Context) {
	if !h.requireAuth(c) {
		return
	}

	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"Status": 400, "Errors": []gin.H{{"Code": "invalid_body", "Message": "invalid request body"}}})
		return
	}

	note, err := h.apiSimulator.HandleCreateCreditNote(body)
	if err != nil {
		c.JSON(400, gin.H{"Status": 400, "Errors": []gin.H{{"Code": "invalid_data", "Message": err.Error()}}})
		return
	}

	c.JSON(201, gin.H{
		"id":       note.ID,
		"name":     note.Name,
		"number":   note.Number,
		"date":     note.Date,
		"total":    note.Amount,
		"metadata": gin.H{"cufe": note.CUFE},
		"stamp":    gin.H{"cufe": note.CUFE, "status": "Stamped"},
	})
}

func (h *Handler) handleCreateVoucher(c *gin.Context) {
	if !h.requireAuth(c) {
		return
	}

	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"Status": 400, "Errors": []gin.H{{"Code": "invalid_body", "Message": "invalid request body"}}})
		return
	}

	voucher, err := h.apiSimulator.HandleCreateVoucher(body)
	if err != nil {
		c.JSON(400, gin.H{"Status": 400, "Errors": []gin.H{{"Code": "invalid_data", "Message": err.Error()}}})
		return
	}

	c.JSON(201, gin.H{
		"id":     voucher.ID,
		"name":   voucher.Name,
		"number": voucher.Number,
		"date":   voucher.Date,
		"total":  voucher.Value,
	})
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
		"prefix": invoice.Prefix,
		"date":   invoice.Date,
		"number": invoice.Number,
		"status": invoice.Status,
		"document": gin.H{
			"id":     invoice.Number,
			"prefix": invoice.Prefix,
			"number": invoice.Number,
		},
		"customer": gin.H{
			"id":             invoice.CustomerID,
			"identification": invoice.CustomerNIT,
			"branch_office":  0,
		},
		"items":      items,
		"total":      invoice.Total,
		"balance":    invoice.Balance,
		"cufe":       invoice.CUFE,
		"public_url": invoice.PublicURL,
		"stamp": gin.H{
			"status": invoice.StampStatus,
			"cufe":   invoice.CUFE,
		},
		"metadata": gin.H{
			"cufe": invoice.CUFE,
		},
	}
}
