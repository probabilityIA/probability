package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/migration/shared/models"
	"github.com/secamc93/probability/back/testing/integrations/envioclick"
	"github.com/secamc93/probability/back/testing/integrations/shopify"
	"github.com/secamc93/probability/back/testing/integrations/softpymes"
	"github.com/secamc93/probability/back/testing/integrations/whatsapp"
	"github.com/secamc93/probability/back/testing/modules/orders"
	"github.com/secamc93/probability/back/testing/shared/db"
	"github.com/secamc93/probability/back/testing/shared/env"
	"github.com/secamc93/probability/back/testing/shared/log"
	"github.com/secamc93/probability/back/testing/shared/middleware"
	"github.com/secamc93/probability/back/testing/shared/storage"
)

func main() {
	logger := log.New()
	config := env.New()

	// 1. Connect to DB (read-only)
	database := db.New(logger, config)
	defer database.Close()

	// 2. Start Softpymes HTTP mock (background)
	softpymesPort := getEnv("SOFTPYMES_MOCK_PORT", "9090")
	softpymesServer := softpymes.New(logger, softpymesPort)

	go func() {
		if err := softpymesServer.Start(); err != nil {
			logger.Error().Msgf("Error starting Softpymes: %s", err.Error())
			os.Exit(1)
		}
	}()

	// 3. Start EnvioClick HTTP mock (background)
	envioclickPort := getEnv("ENVIOCLICK_MOCK_PORT", "9091")
	var s3Service storage.IS3Service
	urlBase := config.Get("URL_BASE_DOMAIN_S3")
	if config.Get("S3_KEY") != "" && config.Get("S3_SECRET") != "" {
		s3Service = storage.New(config, logger)
	} else {
		logger.Warn().Msg("S3 not configured — EnvioClick PDFs will use mock URL")
	}
	envioclickServer := envioclick.New(logger, envioclickPort, s3Service, urlBase)

	go func() {
		if err := envioclickServer.Start(); err != nil {
			logger.Error().Msgf("Error starting EnvioClick: %s", err.Error())
			os.Exit(1)
		}
	}()

	// 4. Initialize Shopify integration (shared between API and CLI)
	shopifyMockPort := getEnv("SHOPIFY_MOCK_PORT", "9093")
	shopifyIntegration := shopify.New(config, logger, shopifyMockPort)

	// 4b. Start Shopify Mock API (simula GET /admin/api/2024-10/orders.json)
	go func() {
		initialOrders := 500 // Pre-generar 500 órdenes distribuidas en 6 meses
		if err := shopifyIntegration.Start(initialOrders); err != nil {
			logger.Error().Msgf("Error starting Shopify Mock API: %s", err.Error())
			os.Exit(1)
		}
	}()

	// 5. Start Testing Platform API (background)
	apiPort := config.GetWithDefault("TESTING_API_PORT", "9092")
	jwtSecret := config.Get("JWT_SECRET")
	if jwtSecret == "" {
		logger.Fatal().Msg("JWT_SECRET is required")
	}

	go func() {
		if err := startAPIServer(logger, config, database, jwtSecret, apiPort, shopifyIntegration); err != nil {
			logger.Error().Msgf("Error starting Testing API: %s", err.Error())
			os.Exit(1)
		}
	}()

	fmt.Println("========================================")
	fmt.Printf("Testing Server - Simuladores\n")
	fmt.Printf("Softpymes HTTP:    http://localhost:%s\n", softpymesPort)
	fmt.Printf("EnvioClick HTTP:   http://localhost:%s\n", envioclickPort)
	fmt.Printf("Shopify Mock API:  http://localhost:%s\n", shopifyMockPort)
	fmt.Printf("Testing API:       http://localhost:%s\n", apiPort)
	fmt.Println("========================================")

	// 6. In server mode (Docker), block forever without interactive CLI
	if os.Getenv("RUN_MODE") == "server" {
		logger.Info().Msg("Running in server mode (no CLI)")
		select {}
	}

	// 7. Start interactive CLI (local development only)
	runCLIMode(logger, config, shopifyIntegration, softpymesServer, envioclickServer)
}

func startAPIServer(logger log.ILogger, config env.IConfig, database db.IDatabase, jwtSecret, port string, shopifyIntegration *shopify.ShopifyIntegration) error {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.CORS())

	// Health check (no auth)
	router.GET("/api/v1/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "testing-platform",
		})
	})

	// Auth proxy - forward login to central backend (no auth required)
	centralAPIURL := config.GetWithDefault("CENTRAL_API_URL", "http://localhost:3050")
	router.POST("/api/v1/auth/login", func(c *gin.Context) {
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}

		resp, err := http.Post(centralAPIURL+"/api/v1/auth/login", "application/json", strings.NewReader(string(body)))
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "central backend not reachable"})
			return
		}
		defer resp.Body.Close()

		respBody, _ := io.ReadAll(resp.Body)

		// Forward Set-Cookie headers
		for _, cookie := range resp.Cookies() {
			c.SetCookie(cookie.Name, cookie.Value, cookie.MaxAge, cookie.Path, cookie.Domain, cookie.Secure, cookie.HttpOnly)
		}

		c.Data(resp.StatusCode, "application/json", respBody)
	})

	// Protected routes
	api := router.Group("/api/v1")
	api.Use(middleware.JWTAuth(jwtSecret))
	api.Use(middleware.SuperAdminGuard())

	// Businesses endpoint (for business selector dropdown)
	api.GET("/businesses", func(c *gin.Context) {
		var businesses []models.Business
		allowedIDs := middleware.GetAllowedBusinessIDs()
		database.Conn(c.Request.Context()).
			Select("id, name").
			Where("id IN ?", allowedIDs).
			Find(&businesses)

		type businessInfo struct {
			ID   uint   `json:"id"`
			Name string `json:"name"`
		}
		result := make([]businessInfo, len(businesses))
		for i, b := range businesses {
			result[i] = businessInfo{ID: b.ID, Name: b.Name}
		}
		c.JSON(200, gin.H{"data": result})
	})

	// Orders routes (require business whitelist)
	ordersGroup := api.Group("/orders")
	ordersGroup.Use(middleware.BusinessWhitelist())

	// Register webhook simulators by integration type code
	webhookSimulators := map[string]orders.IWebhookSimulator{
		"Shopify": shopifyIntegration,
	}

	orders.New(ordersGroup, database, centralAPIURL, logger, webhookSimulators)

	logger.Info().Str("port", port).Msg("Testing Platform API started")
	return router.Run(":" + port)
}

// runCLIMode starts the interactive CLI for webhook simulation
func runCLIMode(logger log.ILogger, config env.IConfig, shopifyIntegration *shopify.ShopifyIntegration, softpymesIntegration *softpymes.SoftPymesIntegration, envioclickIntegration *envioclick.EnvioClickIntegration) {
	whatsappIntegration := whatsapp.New(config, logger)

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("\n=== Simulador de APIs - Menu Principal ===")
		fmt.Println("\n1. Shopify")
		fmt.Println("2. WhatsApp")
		fmt.Println("3. Softpymes")
		fmt.Println("4. EnvioClick")
		fmt.Println("\n0. Salir")
		fmt.Print("\nSelecciona un modulo: ")

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input {
		case "1":
			showShopifyMenu(reader, shopifyIntegration, logger)
		case "2":
			showWhatsAppMenu(reader, whatsappIntegration, logger)
		case "3":
			showSoftpymesMenu(reader, softpymesIntegration, logger)
		case "4":
			showEnvioClickMenu(reader, envioclickIntegration, logger)
		case "0":
			fmt.Println("Saliendo...")
			os.Exit(0)
		default:
			fmt.Println("Opcion invalida")
		}
	}
}

// showShopifyMenu shows the Shopify menu
func showShopifyMenu(reader *bufio.Reader, integration *shopify.ShopifyIntegration, logger log.ILogger) {
	for {
		fmt.Println("\n=== Shopify - Simulador de Webhooks ===")
		fmt.Println("\n1. orders/create (crear nueva orden)")
		fmt.Println("2. orders/paid (marcar como pagada)")
		fmt.Println("3. orders/updated (actualizar orden)")
		fmt.Println("4. orders/cancelled (cancelar orden)")
		fmt.Println("5. orders/fulfilled (marcar como cumplida)")
		fmt.Println("6. orders/partially_fulfilled (parcialmente cumplida)")
		fmt.Println("7. Listar ordenes almacenadas")
		fmt.Println("\n0. Volver al menu principal")
		fmt.Print("\nOpcion: ")

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input {
		case "1", "2", "3", "4", "5", "6":
			handleShopifyOrder(input, integration, logger)
		case "7":
			listShopifyOrders(integration)
		case "0":
			return
		default:
			fmt.Println("Opcion invalida")
		}
	}
}

// showWhatsAppMenu shows the WhatsApp menu
func showWhatsAppMenu(reader *bufio.Reader, integration *whatsapp.WhatsAppIntegration, logger log.ILogger) {
	for {
		fmt.Println("\n=== WhatsApp - Simulador ===")
		fmt.Println("\n1. Simular respuesta de usuario (manual)")
		fmt.Println("2. Simular respuesta automatica (por template)")
		fmt.Println("3. Listar conversaciones almacenadas")
		fmt.Println("\n0. Volver al menu principal")
		fmt.Print("\nOpcion: ")

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input {
		case "1":
			handleWhatsAppUserResponse(reader, integration, logger)
		case "2":
			handleWhatsAppAutoResponse(reader, integration, logger)
		case "3":
			listWhatsAppConversations(integration)
		case "0":
			return
		default:
			fmt.Println("Opcion invalida")
		}
	}
}

// showSoftpymesMenu shows the Softpymes menu
func showSoftpymesMenu(reader *bufio.Reader, integration *softpymes.SoftPymesIntegration, logger log.ILogger) {
	for {
		fmt.Println("\n=== Softpymes - Facturacion ===")
		fmt.Println("\n1. Simular autenticacion")
		fmt.Println("2. Simular creacion de factura")
		fmt.Println("3. Simular nota de credito")
		fmt.Println("4. Listar facturas almacenadas")
		fmt.Println("\n0. Volver al menu principal")
		fmt.Print("\nOpcion: ")

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input {
		case "1":
			handleSoftpymesAuth(reader, integration, logger)
		case "2":
			handleSoftpymesInvoice(reader, integration, logger)
		case "3":
			handleSoftpymesCreditNote(reader, integration, logger)
		case "4":
			listSoftpymesDocuments(integration)
		case "0":
			return
		default:
			fmt.Println("Opcion invalida")
		}
	}
}

// handleShopifyOrder handles Shopify order simulation
func handleShopifyOrder(input string, integration *shopify.ShopifyIntegration, logger log.ILogger) {
	var topic string
	switch input {
	case "1":
		topic = "orders/create"
	case "2":
		topic = "orders/paid"
	case "3":
		topic = "orders/updated"
	case "4":
		topic = "orders/cancelled"
	case "5":
		topic = "orders/fulfilled"
	case "6":
		topic = "orders/partially_fulfilled"
	}

	logger.Info().Str("topic", topic).Msg("Simulating webhook")

	if err := integration.SimulateOrder(topic); err != nil {
		logger.Error().Err(err).Str("topic", topic).Msg("Error simulating webhook")
		fmt.Printf("Error: %v\n", err)
	} else {
		logger.Info().Str("topic", topic).Msg("Webhook simulated successfully")
		fmt.Printf("Webhook '%s' sent successfully\n", topic)
	}
}

// listShopifyOrders lists all Shopify orders
func listShopifyOrders(integration *shopify.ShopifyIntegration) {
	orders := integration.GetAllOrders()
	if len(orders) == 0 {
		fmt.Println("No orders stored")
	} else {
		fmt.Printf("\nStored orders (%d):\n", len(orders))
		for i, order := range orders {
			status := order.FinancialStatus
			if order.FulfillmentStatus != nil {
				status += " / " + *order.FulfillmentStatus
			}
			fmt.Printf("  %d. %s - %s - %s %s - Status: %s\n",
				i+1, order.Name, order.Email, order.Currency, order.TotalPrice, status)
		}
	}
}

// handleWhatsAppUserResponse handles manual WhatsApp response
func handleWhatsAppUserResponse(reader *bufio.Reader, integration *whatsapp.WhatsAppIntegration, logger log.ILogger) {
	fmt.Print("Phone number (eg: +573001234567): ")
	phoneInput, _ := reader.ReadString('\n')
	phoneNumber := strings.TrimSpace(phoneInput)

	fmt.Print("User response (eg: Confirmar pedido): ")
	responseInput, _ := reader.ReadString('\n')
	response := strings.TrimSpace(responseInput)

	logger.Info().
		Str("phone_number", phoneNumber).
		Str("response", response).
		Msg("Simulating user response")

	if err := integration.SimulateUserResponse(phoneNumber, response); err != nil {
		logger.Error().Err(err).Msg("Error simulating user response")
		fmt.Printf("Error: %v\n", err)
	} else {
		logger.Info().Msg("User response simulated successfully")
		fmt.Printf("Response '%s' sent for %s\n", response, phoneNumber)
	}
}

// handleWhatsAppAutoResponse handles automatic WhatsApp response
func handleWhatsAppAutoResponse(reader *bufio.Reader, integration *whatsapp.WhatsAppIntegration, logger log.ILogger) {
	fmt.Print("Phone number (eg: +573001234567): ")
	phoneInput, _ := reader.ReadString('\n')
	phoneNumber := strings.TrimSpace(phoneInput)

	fmt.Print("Template name (eg: confirmacion_pedido_contraentrega): ")
	templateInput, _ := reader.ReadString('\n')
	templateName := strings.TrimSpace(templateInput)

	logger.Info().
		Str("phone_number", phoneNumber).
		Str("template", templateName).
		Msg("Simulating auto response")

	if err := integration.SimulateAutoResponse(phoneNumber, templateName); err != nil {
		logger.Error().Err(err).Msg("Error simulating auto response")
		fmt.Printf("Error: %v\n", err)
	} else {
		logger.Info().Msg("Auto response simulated successfully")
		fmt.Printf("Auto response sent for %s\n", phoneNumber)
	}
}

// listWhatsAppConversations lists all WhatsApp conversations
func listWhatsAppConversations(integration *whatsapp.WhatsAppIntegration) {
	conversations := integration.GetAllConversations()
	if len(conversations) == 0 {
		fmt.Println("No conversations stored")
	} else {
		fmt.Printf("\nStored conversations (%d):\n", len(conversations))
		for i, conv := range conversations {
			messages := integration.GetMessages(conv.ID)
			fmt.Printf("  %d. %s - State: %s - Order: %s - Messages: %d\n",
				i+1, conv.PhoneNumber, conv.CurrentState, conv.OrderNumber, len(messages))
		}
	}
}

// handleSoftpymesAuth handles Softpymes authentication
func handleSoftpymesAuth(reader *bufio.Reader, integration *softpymes.SoftPymesIntegration, logger log.ILogger) {
	fmt.Print("API Key: ")
	apiKeyInput, _ := reader.ReadString('\n')
	apiKey := strings.TrimSpace(apiKeyInput)

	fmt.Print("API Secret: ")
	apiSecretInput, _ := reader.ReadString('\n')
	apiSecret := strings.TrimSpace(apiSecretInput)

	fmt.Print("Referer (eg: https://tutienda.com): ")
	refererInput, _ := reader.ReadString('\n')
	referer := strings.TrimSpace(refererInput)

	logger.Info().Msg("Simulating SoftPymes authentication")

	token, err := integration.SimulateAuth(apiKey, apiSecret, referer)
	if err != nil {
		logger.Error().Err(err).Msg("Error authenticating")
		fmt.Printf("Error: %v\n", err)
	} else {
		logger.Info().Str("token", token).Msg("Authentication successful")
		fmt.Printf("Token generated: %s\n", token)
		fmt.Println("Save this token to create invoices")
	}
}

// handleSoftpymesInvoice handles Softpymes invoice creation
func handleSoftpymesInvoice(reader *bufio.Reader, integration *softpymes.SoftPymesIntegration, logger log.ILogger) {
	fmt.Print("Token: ")
	tokenInput, _ := reader.ReadString('\n')
	token := strings.TrimSpace(tokenInput)

	fmt.Print("Order ID (eg: ORD-001): ")
	orderIDInput, _ := reader.ReadString('\n')
	orderID := strings.TrimSpace(orderIDInput)

	fmt.Print("Customer name: ")
	customerNameInput, _ := reader.ReadString('\n')
	customerName := strings.TrimSpace(customerNameInput)

	fmt.Print("Customer email: ")
	customerEmailInput, _ := reader.ReadString('\n')
	customerEmail := strings.TrimSpace(customerEmailInput)

	fmt.Print("Customer NIT: ")
	customerNITInput, _ := reader.ReadString('\n')
	customerNIT := strings.TrimSpace(customerNITInput)

	fmt.Print("Total (eg: 100000): ")
	totalInput, _ := reader.ReadString('\n')
	totalStr := strings.TrimSpace(totalInput)
	var total float64
	fmt.Sscanf(totalStr, "%f", &total)

	invoiceData := map[string]interface{}{
		"order_id": orderID,
		"customer": map[string]interface{}{
			"name":  customerName,
			"email": customerEmail,
			"nit":   customerNIT,
		},
		"items": []interface{}{
			map[string]interface{}{
				"description": "Producto Test",
				"quantity":    1.0,
				"unit_price":  total,
				"tax":         total * 0.19,
				"total":       total * 1.19,
			},
		},
		"total": total,
	}

	logger.Info().Msg("Simulating invoice creation")

	invoice, err := integration.SimulateInvoice(token, invoiceData)
	if err != nil {
		logger.Error().Err(err).Msg("Error creating invoice")
		fmt.Printf("Error: %v\n", err)
	} else {
		logger.Info().
			Str("invoice_number", invoice.InvoiceNumber).
			Str("cufe", invoice.CUFE).
			Msg("Invoice created successfully")
		fmt.Printf("Invoice created:\n")
		fmt.Printf("  Number: %s\n", invoice.InvoiceNumber)
		fmt.Printf("  CUFE: %s\n", invoice.CUFE)
		fmt.Printf("  Total: $%.2f %s\n", invoice.Total, invoice.Currency)
		fmt.Printf("  PDF: %s\n", invoice.PDFURL)
	}
}

// handleSoftpymesCreditNote handles Softpymes credit note creation
func handleSoftpymesCreditNote(reader *bufio.Reader, integration *softpymes.SoftPymesIntegration, logger log.ILogger) {
	fmt.Print("Token: ")
	tokenInput, _ := reader.ReadString('\n')
	token := strings.TrimSpace(tokenInput)

	fmt.Print("Invoice ID (external_id): ")
	invoiceIDInput, _ := reader.ReadString('\n')
	invoiceID := strings.TrimSpace(invoiceIDInput)

	fmt.Print("Amount to credit: ")
	amountInput, _ := reader.ReadString('\n')
	amountStr := strings.TrimSpace(amountInput)
	var amount float64
	fmt.Sscanf(amountStr, "%f", &amount)

	fmt.Print("Reason (eg: Product return): ")
	reasonInput, _ := reader.ReadString('\n')
	reason := strings.TrimSpace(reasonInput)

	fmt.Print("Type (total/partial): ")
	noteTypeInput, _ := reader.ReadString('\n')
	noteType := strings.TrimSpace(noteTypeInput)

	creditNoteData := map[string]interface{}{
		"invoice_id": invoiceID,
		"amount":     amount,
		"reason":     reason,
		"note_type":  noteType,
	}

	logger.Info().Msg("Simulating credit note creation")

	creditNote, err := integration.SimulateCreditNote(token, creditNoteData)
	if err != nil {
		logger.Error().Err(err).Msg("Error creating credit note")
		fmt.Printf("Error: %v\n", err)
	} else {
		logger.Info().
			Str("note_number", creditNote.CreditNoteNumber).
			Str("cufe", creditNote.CUFE).
			Msg("Credit note created successfully")
		fmt.Printf("Credit note created:\n")
		fmt.Printf("  Number: %s\n", creditNote.CreditNoteNumber)
		fmt.Printf("  CUFE: %s\n", creditNote.CUFE)
		fmt.Printf("  Amount: $%.2f\n", creditNote.Amount)
		fmt.Printf("  Type: %s\n", creditNote.NoteType)
		fmt.Printf("  PDF: %s\n", creditNote.PDFURL)
	}
}

// listSoftpymesDocuments lists all Softpymes documents
func listSoftpymesDocuments(integration *softpymes.SoftPymesIntegration) {
	repo := integration.GetRepository()
	invoices := repo.GetAllInvoices()
	creditNotes := repo.GetAllCreditNotes()

	if len(invoices) == 0 && len(creditNotes) == 0 {
		fmt.Println("No documents stored")
	} else {
		if len(invoices) > 0 {
			fmt.Printf("\nStored invoices (%d):\n", len(invoices))
			for i, invoice := range invoices {
				fmt.Printf("  %d. %s - %s - $%.2f %s - Customer: %s\n",
					i+1, invoice.InvoiceNumber, invoice.OrderID, invoice.Total, invoice.Currency, invoice.CustomerName)
			}
		}
		if len(creditNotes) > 0 {
			fmt.Printf("\nStored credit notes (%d):\n", len(creditNotes))
			for i, note := range creditNotes {
				fmt.Printf("  %d. %s - Invoice: %s - $%.2f - Type: %s\n",
					i+1, note.CreditNoteNumber, note.InvoiceID, note.Amount, note.NoteType)
			}
		}
	}
}

// showEnvioClickMenu shows the EnvioClick menu
func showEnvioClickMenu(reader *bufio.Reader, integration *envioclick.EnvioClickIntegration, logger log.ILogger) {
	for {
		fmt.Println("\n=== EnvioClick - Envios ===")
		fmt.Println("\n1. Cotizar envio")
		fmt.Println("2. Generar guia")
		fmt.Println("3. Rastrear envio")
		fmt.Println("4. Cancelar envio")
		fmt.Println("5. Listar envios almacenados")
		fmt.Println("\n0. Volver al menu principal")
		fmt.Print("\nOpcion: ")

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input {
		case "1":
			handleEnvioClickQuote(reader, integration, logger)
		case "2":
			handleEnvioClickGenerate(reader, integration, logger)
		case "3":
			handleEnvioClickTrack(reader, integration, logger)
		case "4":
			handleEnvioClickCancel(reader, integration, logger)
		case "5":
			listEnvioClickShipments(integration)
		case "0":
			return
		default:
			fmt.Println("Opcion invalida")
		}
	}
}

// readEnvioClickRequest reads common shipment data from user
func readEnvioClickRequest(reader *bufio.Reader) envioclick.QuoteRequest {
	fmt.Print("DANE code origin (eg: 11001 for Bogota): ")
	originInput, _ := reader.ReadString('\n')
	originDane := strings.TrimSpace(originInput)

	fmt.Print("DANE code destination (eg: 05001 for Medellin): ")
	destInput, _ := reader.ReadString('\n')
	destDane := strings.TrimSpace(destInput)

	fmt.Print("Weight in kg (eg: 2.5): ")
	weightInput, _ := reader.ReadString('\n')
	var weight float64
	fmt.Sscanf(strings.TrimSpace(weightInput), "%f", &weight)
	if weight <= 0 {
		weight = 1.0
	}

	fmt.Print("Height in cm (eg: 20): ")
	heightInput, _ := reader.ReadString('\n')
	var height float64
	fmt.Sscanf(strings.TrimSpace(heightInput), "%f", &height)
	if height <= 0 {
		height = 10.0
	}

	fmt.Print("Width in cm (eg: 15): ")
	widthInput, _ := reader.ReadString('\n')
	var width float64
	fmt.Sscanf(strings.TrimSpace(widthInput), "%f", &width)
	if width <= 0 {
		width = 10.0
	}

	fmt.Print("Length in cm (eg: 30): ")
	lengthInput, _ := reader.ReadString('\n')
	var length float64
	fmt.Sscanf(strings.TrimSpace(lengthInput), "%f", &length)
	if length <= 0 {
		length = 10.0
	}

	fmt.Print("Declared value COP (eg: 50000): ")
	valueInput, _ := reader.ReadString('\n')
	var contentValue float64
	fmt.Sscanf(strings.TrimSpace(valueInput), "%f", &contentValue)
	if contentValue <= 0 {
		contentValue = 10000
	}

	return envioclick.QuoteRequest{
		Origin:       envioclick.Address{DaneCode: originDane},
		Destination:  envioclick.Address{DaneCode: destDane},
		Packages:     []envioclick.Package{{Weight: weight, Height: height, Width: width, Length: length}},
		ContentValue: contentValue,
	}
}

// handleEnvioClickQuote handles shipment quotation
func handleEnvioClickQuote(reader *bufio.Reader, integration *envioclick.EnvioClickIntegration, logger log.ILogger) {
	req := readEnvioClickRequest(reader)
	logger.Info().Msg("Simulating shipment quotation")

	resp, err := integration.SimulateQuote(req)
	if err != nil {
		logger.Error().Err(err).Msg("Error quoting")
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("\nQuotation successful - %d rates available:\n", len(resp.Data.Rates))
	for i, rate := range resp.Data.Rates {
		fmt.Printf("  %d. [%s] %s - $%.0f COP - %d days - ID: %d\n",
			i+1, rate.Carrier, rate.Product, rate.Flete, rate.DeliveryDays, rate.IDRate)
	}
}

// handleEnvioClickGenerate handles shipment label generation
func handleEnvioClickGenerate(reader *bufio.Reader, integration *envioclick.EnvioClickIntegration, logger log.ILogger) {
	req := readEnvioClickRequest(reader)

	fmt.Print("Rate ID (from quotation, eg: 1001): ")
	rateInput, _ := reader.ReadString('\n')
	var rateID int64
	fmt.Sscanf(strings.TrimSpace(rateInput), "%d", &rateID)
	req.IDRate = rateID

	fmt.Print("Shipment reference (eg: ORD-001): ")
	refInput, _ := reader.ReadString('\n')
	req.MyShipmentReference = strings.TrimSpace(refInput)

	logger.Info().Msg("Simulating label generation")

	resp, err := integration.SimulateGenerate(req)
	if err != nil {
		logger.Error().Err(err).Msg("Error generating label")
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("\nLabel generated successfully:\n")
	fmt.Printf("  Tracking: %s\n", resp.Data.TrackingNumber)
	fmt.Printf("  Label URL: %s\n", resp.Data.LabelURL)
	fmt.Printf("  Reference: %s\n", resp.Data.MyGuideReference)
}

// handleEnvioClickTrack handles shipment tracking
func handleEnvioClickTrack(reader *bufio.Reader, integration *envioclick.EnvioClickIntegration, logger log.ILogger) {
	fmt.Print("Tracking number: ")
	trackInput, _ := reader.ReadString('\n')
	trackingNumber := strings.TrimSpace(trackInput)

	logger.Info().Str("tracking", trackingNumber).Msg("Simulating tracking")

	resp, err := integration.SimulateTrack(trackingNumber)
	if err != nil {
		logger.Error().Err(err).Msg("Error tracking")
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("\nTracking %s:\n", resp.Data.TrackingNumber)
	fmt.Printf("  Carrier: %s\n", resp.Data.Carrier)
	fmt.Printf("  Status: %s\n", resp.Data.Status)
	fmt.Printf("  Events (%d):\n", len(resp.Data.Events))
	for i, event := range resp.Data.Events {
		fmt.Printf("    %d. [%s] %s - %s (%s)\n",
			i+1, event.Date, event.Status, event.Description, event.Location)
	}
}

// handleEnvioClickCancel handles shipment cancellation
func handleEnvioClickCancel(reader *bufio.Reader, integration *envioclick.EnvioClickIntegration, logger log.ILogger) {
	fmt.Print("Shipment ID (eg: EC-005001): ")
	idInput, _ := reader.ReadString('\n')
	shipmentID := strings.TrimSpace(idInput)

	logger.Info().Str("shipment_id", shipmentID).Msg("Simulating cancellation")

	resp, err := integration.SimulateCancel(shipmentID)
	if err != nil {
		logger.Error().Err(err).Msg("Error cancelling")
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("%s: %s\n", resp.Status, resp.Message)
}

// listEnvioClickShipments lists all stored shipments
func listEnvioClickShipments(integration *envioclick.EnvioClickIntegration) {
	shipments := integration.GetAllShipments()
	if len(shipments) == 0 {
		fmt.Println("No shipments stored")
		return
	}

	fmt.Printf("\nStored shipments (%d):\n", len(shipments))
	for i, s := range shipments {
		fmt.Printf("  %d. ID: %s - Tracking: %s - %s - %s -> %s - $%.0f COP - Status: %s\n",
			i+1, s.ID, s.TrackingNumber, s.Carrier,
			s.Origin.DaneCode, s.Destination.DaneCode,
			s.Flete, s.Status)
	}
}

// getEnv gets environment variable or returns default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
