package handlers

import (
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/testing/integrations/softpymes/internal/app/usecases"
	"github.com/secamc93/probability/back/testing/integrations/softpymes/internal/domain"
)

// RegisterRoutes registra todas las rutas del simulador
func (h *Handler) RegisterRoutes(router *gin.Engine) {
	// Autenticación
	router.POST("/oauth/integration/login/", h.handleAuth)

	// Clientes
	router.GET("/app/integration/customer", h.handleGetCustomer)
	router.POST("/app/integration/customer", h.handleCreateCustomer)
	router.POST("/app/integration/customer_new/", h.handleCreateCustomer)

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

// handleGetCustomer busca un cliente por identificación
func (h *Handler) handleGetCustomer(c *gin.Context) {
	token := h.extractToken(c)
	if token == "" {
		c.JSON(401, gin.H{"message": "Missing authorization token"})
		return
	}

	identification := c.Query("identification")
	if identification == "" {
		c.JSON(400, gin.H{"message": "identification is required"})
		return
	}

	customer, err := h.apiSimulator.HandleGetCustomer(token, identification)
	if err != nil {
		if err.Error() == "customer not found" {
			// Softpymes retorna 200 con solo "message" cuando no existe
			c.JSON(200, gin.H{"message": fmt.Sprintf("No existe cliente con identificación: %s", identification)})
			return
		}
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"identification": customer.Identification,
		"name":           customer.Name,
		"email":          customer.Email,
		"phone":          customer.Phone,
		"branch":         customer.Branch,
	})
}

// handleCreateCustomer crea un nuevo cliente
// El cliente real envía: identificationNumber, firstName, lastName, email, phone, cellPhone, etc.
func (h *Handler) handleCreateCustomer(c *gin.Context) {
	token := h.extractToken(c)
	if token == "" {
		c.JSON(401, gin.H{"message": "Missing authorization token"})
		return
	}

	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"message": "Invalid request format"})
		return
	}

	// Normalizar al formato interno: el cliente real envía identificationNumber, firstName, lastName
	customerData := map[string]interface{}{}
	if id, ok := body["identificationNumber"].(string); ok {
		customerData["identification"] = id
	} else if id, ok := body["identification"].(string); ok {
		customerData["identification"] = id
	}

	// Combinar firstName + lastName como name
	firstName, _ := body["firstName"].(string)
	lastName, _ := body["lastName"].(string)
	if firstName != "" || lastName != "" {
		customerData["name"] = fmt.Sprintf("%s %s", firstName, lastName)
	} else if name, ok := body["name"].(string); ok {
		customerData["name"] = name
	}

	if email, ok := body["email"].(string); ok {
		customerData["email"] = email
	}
	if phone, ok := body["phone"].(string); ok {
		customerData["phone"] = phone
	} else if phone, ok := body["cellPhone"].(string); ok {
		customerData["phone"] = phone
	}

	customer, err := h.apiSimulator.HandleCreateCustomer(token, customerData)
	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"identification": customer.Identification,
		"name":           customer.Name,
		"email":          customer.Email,
		"phone":          customer.Phone,
		"branch":         customer.Branch,
	})
}

// handleCreateInvoice maneja la creación de factura
func (h *Handler) handleCreateInvoice(c *gin.Context) {
	token := h.extractToken(c)
	if token == "" {
		c.JSON(401, gin.H{"message": "Missing authorization token", "type": "UNAUTHORIZED"})
		return
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
			// Igual que la API real: prefix + número SIN ceros (ej: "FEV1001", no "FEV0000001001")
		"documentNumber": fmt.Sprintf("%s%s", invoice.Prefix, strings.TrimLeft(invoice.InvoiceNumber, "0")),
			"subtotal":       invoice.Subtotal,
			"discount":       0.0,
			"iva":            invoice.IVA,
			"withholding":    0.0,
			"total":          invoice.Total,
			"docsFe": gin.H{
				"status":          "Aceptado",
				"message":         "Documento válido enviado al proveedor tecnológico",
				"quantitySlopes":  nil,
			},
		},
	})
}

// extractToken extrae el Bearer token del header Authorization
func (h *Handler) extractToken(c *gin.Context) string {
	token := c.GetHeader("Authorization")
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}
	return token
}

// handleSearchDocuments maneja la búsqueda de documentos
func (h *Handler) handleSearchDocuments(c *gin.Context) {
	token := h.extractToken(c)
	if token == "" {
		c.JSON(401, gin.H{"message": "Missing authorization token", "type": "UNAUTHORIZED"})
		return
	}

	var searchParams struct {
		DateFrom       string  `json:"dateFrom"`
		DateTo         string  `json:"dateTo"`
		DocumentNumber *string `json:"documentNumber"`
		Prefix         *string `json:"prefix"`
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
		Invoice:    *inv,
		BranchCode: "001",
		BranchName: "Sucursal Principal",
		Prefix:     "FEV",
		SellerName: "Empresa Demo S.A.S.",
		SellerNIT:  "900123456-7",
	}
}

// buildDocumentResponse construye la respuesta del documento con datos reales
func (h *Handler) buildDocumentResponse(invoice *usecases.InvoiceWithDetails) map[string]interface{} {
	items := make([]map[string]interface{}, 0, len(invoice.Items))
	for _, item := range invoice.Items {
		items = append(items, map[string]interface{}{
			"discount":       fmt.Sprintf("%.2f", item.Discount),
			"itemCode":       item.ItemCode,
			"itemName":       item.ItemName,
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
		"comment":                invoice.Comment,
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
		"termDays":            0,
		"total":               fmt.Sprintf("%.2f", invoice.Total),
		"totalDiscount":       "0.00",
		"totalIva":            fmt.Sprintf("%.2f", invoice.IVA),
		"totalWithholdingTax": "0.00",
	}
}
