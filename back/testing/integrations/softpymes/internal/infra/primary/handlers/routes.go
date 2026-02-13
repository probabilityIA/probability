package handlers

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/testing/integrations/softpymes/internal/app/usecases"
	"github.com/secamc93/probability/back/testing/integrations/softpymes/internal/domain"
)

// RegisterRoutes registra todas las rutas del simulador
func (h *Handler) RegisterRoutes(router *gin.Engine) {
	// Autenticación
	router.POST("/oauth/integration/login/", h.handleAuth)

	// Crear factura
	router.POST("/app/integration/sales_invoice/", h.handleCreateInvoice)

	// Buscar documentos
	router.POST("/app/integration/search/documents/", h.handleSearchDocuments)

	// Health check
	router.GET("/health", h.handleHealth)
}

// handleHealth maneja el health check
func (h *Handler) handleHealth(c *gin.Context) {
	c.JSON(200, gin.H{
		"status":  "ok",
		"service": "softpymes-mock",
	})
}

// handleAuth maneja la autenticación
func (h *Handler) handleAuth(c *gin.Context) {
	var req struct {
		APIKey    string `json:"apiKey" binding:"required"`
		APISecret string `json:"apiSecret" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"message": "Invalid request format",
			"type":    "INVALID_DATA",
		})
		return
	}

	referer := c.GetHeader("Referer")
	if referer == "" {
		referer = "mock-referer"
	}

	token, err := h.apiSimulator.HandleAuth(req.APIKey, req.APISecret, referer)
	if err != nil {
		c.JSON(401, gin.H{
			"message": "Authentication failed",
			"type":    "UNAUTHORIZED",
		})
		return
	}

	c.JSON(200, gin.H{
		"accessToken":  token,
		"expiresInMin": 60,
		"tokenType":    "Bearer",
	})
}

// handleCreateInvoice maneja la creación de factura
func (h *Handler) handleCreateInvoice(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(401, gin.H{
			"message": "Missing authorization token",
			"type":    "UNAUTHORIZED",
		})
		return
	}

	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	var invoiceData map[string]interface{}
	if err := c.ShouldBindJSON(&invoiceData); err != nil {
		c.JSON(400, gin.H{
			"message": "Invalid request format",
			"type":    "INVALID_DATA",
		})
		return
	}

	invoice, err := h.apiSimulator.HandleCreateInvoice(token, invoiceData)
	if err != nil {
		statusCode := 500
		if err.Error() == "invalid token" || err.Error() == "token expired" {
			statusCode = 401
		}

		c.JSON(statusCode, gin.H{
			"message": err.Error(),
			"type":    "ERROR",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Se ha creado la factura de venta en Pymes+ correctamente!",
		"info": gin.H{
			"date":           invoice.IssuedAt.Format(time.RFC3339),
			"documentNumber": invoice.InvoiceNumber,
			"subtotal":       invoice.Total * 0.84,
			"discount":       0.0,
			"iva":            invoice.Total * 0.16,
			"withholding":    0.0,
			"total":          invoice.Total,
			"docsFe": gin.H{
				"status":  true,
				"message": "Documento válido enviado al proveedor tecnológico",
			},
		},
	})
}

// handleSearchDocuments maneja la búsqueda de documentos
func (h *Handler) handleSearchDocuments(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(401, gin.H{
			"message": "Missing authorization token",
			"type":    "UNAUTHORIZED",
		})
		return
	}

	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	var searchParams struct {
		DateFrom       string  `json:"dateFrom"`
		DateTo         string  `json:"dateTo"`
		DocumentNumber *string `json:"documentNumber"`
	}

	if err := c.ShouldBindJSON(&searchParams); err != nil {
		c.JSON(400, gin.H{
			"message": "Invalid request format",
			"type":    "INVALID_DATA",
		})
		return
	}

	var documents []map[string]interface{}

	if searchParams.DocumentNumber != nil && *searchParams.DocumentNumber != "" {
		invoice, err := h.apiSimulator.GetInvoiceByNumber(*searchParams.DocumentNumber)
		if err != nil {
			c.JSON(200, []interface{}{})
			return
		}

		doc := h.buildDocumentResponse(invoice)
		documents = append(documents, doc)
	} else {
		invoices, err := h.apiSimulator.HandleListDocuments(token, nil)
		if err != nil {
			statusCode := 500
			if err.Error() == "invalid token" || err.Error() == "token expired" {
				statusCode = 401
			}
			c.JSON(statusCode, gin.H{
				"message": err.Error(),
				"type":    "ERROR",
			})
			return
		}

		for _, inv := range invoices {
			invWithDetails := h.invoiceToInvoiceWithDetails(&inv)
			doc := h.buildDocumentResponse(invWithDetails)
			documents = append(documents, doc)
		}
	}

	c.JSON(200, documents)
}

// invoiceToInvoiceWithDetails convierte Invoice a InvoiceWithDetails
func (h *Handler) invoiceToInvoiceWithDetails(inv *domain.Invoice) *usecases.InvoiceWithDetails {
	return &usecases.InvoiceWithDetails{
		Invoice: domain.Invoice{
			ID:            inv.ID,
			InvoiceNumber: inv.InvoiceNumber,
			ExternalID:    inv.ExternalID,
			OrderID:       inv.OrderID,
			CustomerName:  inv.CustomerName,
			CustomerEmail: inv.CustomerEmail,
			CustomerNIT:   inv.CustomerNIT,
			Total:         inv.Total,
			Currency:      inv.Currency,
			Items:         inv.Items,
			InvoiceURL:    inv.InvoiceURL,
			PDFURL:        inv.PDFURL,
			XMLURL:        inv.XMLURL,
			CUFE:          inv.CUFE,
			IssuedAt:      inv.IssuedAt,
			CreatedAt:     inv.CreatedAt,
		},
		BranchCode:  "001",
		BranchName:  "Sucursal Principal",
		Prefix:      "SPY",
		SellerName:  "Empresa Demo",
		SellerNIT:   "900123456",
	}
}

// buildDocumentResponse construye la respuesta del documento
func (h *Handler) buildDocumentResponse(invoice *usecases.InvoiceWithDetails) map[string]interface{} {
	items := make([]map[string]interface{}, 0)
	for _, item := range invoice.Items {
		items = append(items, map[string]interface{}{
			"discount":       "0.00",
			"itemCode":       item.ItemCode,
			"itemName":       item.Description,
			"code":           item.ItemCode,
			"service":        "false",
			"iva":            fmt.Sprintf("%.2f", item.Tax),
			"ica":            "0.00",
			"quantity":       fmt.Sprintf("%d", item.Quantity),
			"sizeColor":      map[string]string{},
			"value":          fmt.Sprintf("%.2f", item.UnitPrice),
			"withholdingTax": "0.00",
			"warehouse": map[string]string{
				"code": "001",
				"name": "Principal",
			},
		})
	}

	return map[string]interface{}{
		"branchCode":             invoice.BranchCode,
		"branchName":             invoice.BranchName,
		"comment":                "",
		"customerIdentification": invoice.CustomerNIT,
		"customerName":           invoice.CustomerName,
		"details":                items,
		"documentDate":           invoice.IssuedAt.Format("2006-01-02"),
		"documentName":           "Factura de Venta",
		"documentNumber":         invoice.InvoiceNumber,
		"dueDate":                invoice.IssuedAt.AddDate(0, 0, 30).Format("2006-01-02"),
		"paymentTerm":            "Contado",
		"prefix":                 invoice.Prefix,
		"seller": map[string]string{
			"name": invoice.SellerName,
			"nit":  invoice.SellerNIT,
		},
		"shipInformation": map[string]string{
			"shipAddress":    "",
			"shipCity":       "",
			"shipCountry":    "Colombia",
			"shipDepartment": "",
			"shipPhone":      "",
			"shipTo":         "",
			"shipZipCode":    "",
		},
		"termDays":           0,
		"total":              fmt.Sprintf("%.2f", invoice.Total),
		"totalDiscount":      "0.00",
		"totalIva":           fmt.Sprintf("%.2f", invoice.Total*0.16),
		"totalWithholdingTax": "0.00",
	}
}
