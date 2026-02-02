package main

import (
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/integrationTest/integrations/softpymes/internal/app/usecases"
	customlog "github.com/secamc93/probability/back/integrationTest/shared/log"
)

// Global simulator
var apiSimulator *usecases.APISimulator

func main() {
	// Inicializar logger
	logger := customlog.New()

	// Inicializar simulador
	apiSimulator = usecases.NewAPISimulator(logger).(*usecases.APISimulator)

	port := getEnv("SOFTPYMES_MOCK_PORT", "8082")

	// Configurar Gin
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// Middleware de logging
	router.Use(func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		c.Next()

		duration := time.Since(start)
		status := c.Writer.Status()

		log.Printf("[%s] %s %s - Status: %d - Duration: %v",
			time.Now().Format("15:04:05"),
			method,
			path,
			status,
			duration,
		)
	})

	// Endpoints de SoftPymes
	router.POST("/oauth/integration/login/", handleAuth)
	router.POST("/sales_invoice/", handleCreateInvoice)
	router.POST("/search/documents/notes/", handleCreateCreditNote)
	router.GET("/search/documents/", handleListDocuments)

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "softpymes-mock",
			"port":    port,
		})
	})

	log.Printf("🚀 SoftPymes Mock Server running on port %s", port)
	log.Printf("📋 Endpoints available:")
	log.Printf("  POST /oauth/integration/login/")
	log.Printf("  POST /sales_invoice/")
	log.Printf("  POST /search/documents/notes/")
	log.Printf("  GET  /search/documents/")
	log.Printf("  GET  /health")

	router.Run(":" + port)
}

// handleAuth maneja la autenticación
func handleAuth(c *gin.Context) {
	var req struct {
		APIKey    string `json:"apiKey" binding:"required"`
		APISecret string `json:"apiSecret" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"error":   err.Error(),
			"message": "Invalid request format",
		})
		return
	}

	referer := c.GetHeader("Referer")
	if referer == "" {
		referer = "http://localhost"
	}

	// Simular autenticación
	token, err := apiSimulator.HandleAuth(req.APIKey, req.APISecret, referer)
	if err != nil {
		c.JSON(401, gin.H{
			"success": false,
			"error":   err.Error(),
			"message": "Authentication failed",
		})
		return
	}

	c.JSON(200, gin.H{
		"success":      true,
		"accessToken":  token,
		"expiresInMin": 60,
		"tokenType":    "Bearer",
	})
}

// handleCreateInvoice maneja la creación de factura
func handleCreateInvoice(c *gin.Context) {
	// Extraer token
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(401, gin.H{
			"success": false,
			"error":   "Missing authorization header",
			"message": "Unauthorized",
		})
		return
	}

	// Remover "Bearer " si existe
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	// Obtener body
	var invoiceData map[string]interface{}
	if err := c.ShouldBindJSON(&invoiceData); err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"error":   err.Error(),
			"message": "Invalid request format",
		})
		return
	}

	// Simular creación de factura
	invoice, err := apiSimulator.HandleCreateInvoice(token, invoiceData)
	if err != nil {
		statusCode := 500
		if err.Error() == "invalid token" || err.Error() == "token expired" {
			statusCode = 401
		}

		c.JSON(statusCode, gin.H{
			"success": false,
			"error":   err.Error(),
			"message": "Invoice creation failed",
		})
		return
	}

	c.JSON(200, gin.H{
		"success":        true,
		"message":        "Invoice created successfully",
		"invoice_number": invoice.InvoiceNumber,
		"external_id":    invoice.ExternalID,
		"invoice_url":    invoice.InvoiceURL,
		"pdf_url":        invoice.PDFURL,
		"xml_url":        invoice.XMLURL,
		"cufe":           invoice.CUFE,
		"issued_at":      invoice.IssuedAt.Format(time.RFC3339),
	})
}

// handleCreateCreditNote maneja la creación de nota de crédito
func handleCreateCreditNote(c *gin.Context) {
	// Extraer token
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(401, gin.H{
			"success": false,
			"error":   "Missing authorization header",
			"message": "Unauthorized",
		})
		return
	}

	// Remover "Bearer " si existe
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	// Obtener body
	var creditNoteData map[string]interface{}
	if err := c.ShouldBindJSON(&creditNoteData); err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"error":   err.Error(),
			"message": "Invalid request format",
		})
		return
	}

	// Simular creación de nota de crédito
	creditNote, err := apiSimulator.HandleCreateCreditNote(token, creditNoteData)
	if err != nil {
		statusCode := 500
		if err.Error() == "invalid token" || err.Error() == "token expired" {
			statusCode = 401
		}

		c.JSON(statusCode, gin.H{
			"success": false,
			"error":   err.Error(),
			"message": "Credit note creation failed",
		})
		return
	}

	c.JSON(200, gin.H{
		"success":            true,
		"message":            "Credit note created successfully",
		"credit_note_number": creditNote.CreditNoteNumber,
		"external_id":        creditNote.ExternalID,
		"note_url":           creditNote.NoteURL,
		"pdf_url":            creditNote.PDFURL,
		"xml_url":            creditNote.XMLURL,
		"cufe":               creditNote.CUFE,
		"issued_at":          creditNote.IssuedAt.Format(time.RFC3339),
	})
}

// handleListDocuments maneja el listado de documentos
func handleListDocuments(c *gin.Context) {
	// Extraer token
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(401, gin.H{
			"success": false,
			"error":   "Missing authorization header",
			"message": "Unauthorized",
		})
		return
	}

	// Remover "Bearer " si existe
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	// Simular listado
	invoices, err := apiSimulator.HandleListDocuments(token, nil)
	if err != nil {
		statusCode := 500
		if err.Error() == "invalid token" || err.Error() == "token expired" {
			statusCode = 401
		}

		c.JSON(statusCode, gin.H{
			"success": false,
			"error":   err.Error(),
			"message": "List documents failed",
		})
		return
	}

	c.JSON(200, gin.H{
		"success":   true,
		"count":     len(invoices),
		"documents": invoices,
	})
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
