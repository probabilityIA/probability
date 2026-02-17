package usecases

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/secamc93/probability/back/testing/integrations/softpymes/internal/domain"
	"github.com/secamc93/probability/back/testing/shared/log"
)

// APISimulator simula el API de SoftPymes
type APISimulator struct {
	logger     log.ILogger
	Repository *domain.InvoiceRepository // Exportado para acceso desde bundle
}

// InvoiceWithDetails extiende Invoice con campos adicionales para el endpoint de búsqueda
type InvoiceWithDetails struct {
	domain.Invoice
	BranchCode string
	BranchName string
	Prefix     string
	SellerName string
	SellerNIT  string
}

// HandleAuth simula la autenticación de SoftPymes
func (s *APISimulator) HandleAuth(apiKey, apiSecret, referer string) (string, error) {
	s.logger.Info().
		Str("api_key", apiKey).
		Str("referer", referer).
		Msg("Simulando autenticación de SoftPymes")

	// Validar credenciales (simplificado - acepta cualquier valor para testing)
	if apiKey == "" || apiSecret == "" {
		return "", fmt.Errorf("invalid credentials")
	}

	// Generar token ficticio
	token := fmt.Sprintf("spy_token_%s", uuid.New().String()[:8])
	expiresAt := time.Now().Add(1 * time.Hour)

	// Guardar token
	authToken := &domain.AuthToken{
		Token:     token,
		ExpiresAt: expiresAt,
		APIKey:    apiKey,
		Referer:   referer,
	}
	s.Repository.SaveToken(authToken)

	s.logger.Info().
		Str("token", token).
		Msg("Token generado exitosamente")

	return token, nil
}

// HandleCreateInvoice simula la creación de una factura
// Parsea el formato real de la API de Softpymes (customerNit, items[].itemCode, items[].unitValue string, etc.)
func (s *APISimulator) HandleCreateInvoice(token string, invoiceData map[string]interface{}) (*InvoiceWithDetails, error) {
	s.logger.Info().
		Str("token", token).
		Interface("data", invoiceData).
		Msg("Simulando creación de factura")

	// Validar token
	authToken, exists := s.Repository.GetToken(token)
	if !exists {
		return nil, fmt.Errorf("invalid token")
	}

	// Verificar si expiró
	if time.Now().After(authToken.ExpiresAt) {
		return nil, fmt.Errorf("token expired")
	}

	// Extraer customerNit del top level (formato Softpymes API)
	customerNIT := "999999999"
	if nit, ok := invoiceData["customerNit"].(string); ok && nit != "" {
		customerNIT = nit
	}

	// Extraer items en formato Softpymes (itemCode, unitValue string, quantity, discount)
	items, _ := invoiceData["items"].([]interface{})

	// Generar factura
	invoiceNumber := s.Repository.GenerateInvoiceNumber()
	externalID := uuid.New().String()
	cufe := fmt.Sprintf("CUFE-%s", uuid.New().String()[:16])

	// Convertir items del formato Softpymes al formato interno
	invoiceItems := make([]domain.InvoiceItem, 0, len(items))
	subtotal := 0.0
	for _, item := range items {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		itemCode := ""
		if code, ok := itemMap["itemCode"].(string); ok {
			itemCode = code
		}

		quantity := 1
		if q, ok := itemMap["quantity"].(float64); ok {
			quantity = int(q)
		}

		// unitValue viene como string en el formato Softpymes
		unitPrice := 0.0
		if uv, ok := itemMap["unitValue"].(string); ok {
			fmt.Sscanf(uv, "%f", &unitPrice)
		} else if uv, ok := itemMap["unitValue"].(float64); ok {
			unitPrice = uv
		}

		discount := 0.0
		if d, ok := itemMap["discount"].(float64); ok {
			discount = d
		}

		itemTotal := unitPrice*float64(quantity) - discount
		itemTax := itemTotal * 0.19 // IVA 19%

		invoiceItems = append(invoiceItems, domain.InvoiceItem{
			ItemCode:  itemCode,
			ItemName:  itemCode, // Usar el código como nombre (el mock no tiene catálogo)
			Quantity:  quantity,
			UnitPrice: unitPrice,
			Tax:       itemTax,
			Discount:  discount,
			Total:     itemTotal,
		})

		subtotal += itemTotal
	}

	iva := subtotal * 0.19
	total := subtotal + iva

	// Buscar datos del cliente en el repo de customers (si existe)
	customerName := ""
	customerEmail := ""
	customerPhone := ""
	if cust, exists := s.Repository.GetCustomer(customerNIT); exists {
		customerName = cust.Name
		customerEmail = cust.Email
		customerPhone = cust.Phone
	}

	invoice := &domain.Invoice{
		ID:            externalID,
		InvoiceNumber: invoiceNumber,
		ExternalID:    externalID,
		OrderID:       "",
		CustomerName:  customerName,
		CustomerEmail: customerEmail,
		CustomerNIT:   customerNIT,
		CustomerPhone: customerPhone,
		Total:         total,
		Subtotal:      subtotal,
		IVA:           iva,
		Currency:      "COP",
		Items:         invoiceItems,
		InvoiceURL:    fmt.Sprintf("https://softpymes-mock.local/invoices/%s", externalID),
		PDFURL:        fmt.Sprintf("https://softpymes-mock.local/invoices/%s.pdf", externalID),
		XMLURL:        fmt.Sprintf("https://softpymes-mock.local/invoices/%s.xml", externalID),
		CUFE:          cufe,
		IssuedAt:      time.Now(),
		CreatedAt:     time.Now(),
	}

	s.Repository.SaveInvoice(invoice)

	s.logger.Info().
		Str("invoice_number", invoiceNumber).
		Str("customer_nit", customerNIT).
		Float64("subtotal", subtotal).
		Float64("iva", iva).
		Float64("total", total).
		Int("items_count", len(invoiceItems)).
		Msg("Factura simulada creada exitosamente")

	return &InvoiceWithDetails{
		Invoice:    *invoice,
		BranchCode: "001",
		BranchName: "Sucursal Principal",
		Prefix:     "FV",
		SellerName: "Empresa Demo S.A.S.",
		SellerNIT:  "900123456-7",
	}, nil
}

// GetInvoiceByNumber busca una factura por su número
func (s *APISimulator) GetInvoiceByNumber(invoiceNumber string) (*InvoiceWithDetails, error) {
	s.logger.Info().
		Str("invoice_number", invoiceNumber).
		Msg("Buscando factura por número")

	invoice, exists := s.Repository.GetInvoiceByNumber(invoiceNumber)
	if !exists {
		return nil, fmt.Errorf("invoice not found: %s", invoiceNumber)
	}

	// Retornar con detalles adicionales
	return &InvoiceWithDetails{
		Invoice:    *invoice,
		BranchCode: "001",
		BranchName: "Sucursal Principal",
		Prefix:     "FV",
		SellerName: "Empresa Demo S.A.S.",
		SellerNIT:  "900123456-7",
	}, nil
}

// HandleGetCustomer busca un cliente por identificación
func (s *APISimulator) HandleGetCustomer(token, identification string) (*domain.Customer, error) {
	// Validar token
	authToken, exists := s.Repository.GetToken(token)
	if !exists {
		return nil, fmt.Errorf("invalid token")
	}
	if time.Now().After(authToken.ExpiresAt) {
		return nil, fmt.Errorf("token expired")
	}

	customer, exists := s.Repository.GetCustomer(identification)
	if !exists {
		return nil, fmt.Errorf("customer not found")
	}
	return customer, nil
}

// HandleCreateCustomer crea un nuevo cliente
func (s *APISimulator) HandleCreateCustomer(token string, customerData map[string]interface{}) (*domain.Customer, error) {
	// Validar token
	authToken, exists := s.Repository.GetToken(token)
	if !exists {
		return nil, fmt.Errorf("invalid token")
	}
	if time.Now().After(authToken.ExpiresAt) {
		return nil, fmt.Errorf("token expired")
	}

	identification := ""
	if id, ok := customerData["identification"].(string); ok {
		identification = id
	}
	if identification == "" {
		return nil, fmt.Errorf("identification is required")
	}

	name := ""
	if n, ok := customerData["name"].(string); ok {
		name = n
	}
	email := ""
	if e, ok := customerData["email"].(string); ok {
		email = e
	}
	phone := ""
	if p, ok := customerData["phone"].(string); ok {
		phone = p
	}

	customer := &domain.Customer{
		Identification: identification,
		Name:           name,
		Email:          email,
		Phone:          phone,
		Branch:         "000",
	}

	s.Repository.SaveCustomer(customer)

	s.logger.Info().
		Str("identification", identification).
		Str("name", name).
		Msg("Cliente simulado creado exitosamente")

	return customer, nil
}

// HandleCreateCreditNote simula la creación de una nota de crédito
func (s *APISimulator) HandleCreateCreditNote(token string, creditNoteData map[string]interface{}) (*domain.CreditNote, error) {
	s.logger.Info().
		Str("token", token).
		Interface("data", creditNoteData).
		Msg("Simulando creación de nota de crédito")

	// Validar token
	authToken, exists := s.Repository.GetToken(token)
	if !exists {
		return nil, fmt.Errorf("invalid token")
	}

	// Verificar si expiró
	if time.Now().After(authToken.ExpiresAt) {
		return nil, fmt.Errorf("token expired")
	}

	// Extraer datos
	invoiceID, _ := creditNoteData["invoice_id"].(string)
	amount, _ := creditNoteData["amount"].(float64)
	reason, _ := creditNoteData["reason"].(string)
	noteType, _ := creditNoteData["note_type"].(string)

	// Verificar que exista la factura
	_, exists = s.Repository.GetInvoice(invoiceID)
	if !exists {
		return nil, fmt.Errorf("invoice not found: %s", invoiceID)
	}

	// Generar nota de crédito
	creditNoteNumber := s.Repository.GenerateCreditNoteNumber()
	externalID := uuid.New().String()
	cufe := fmt.Sprintf("CUFE-NC-%s", uuid.New().String()[:16])

	creditNote := &domain.CreditNote{
		ID:               externalID,
		CreditNoteNumber: creditNoteNumber,
		ExternalID:       externalID,
		InvoiceID:        invoiceID,
		Amount:           amount,
		Reason:           reason,
		NoteType:         noteType,
		NoteURL:          fmt.Sprintf("https://softpymes-mock.local/credit-notes/%s", externalID),
		PDFURL:           fmt.Sprintf("https://softpymes-mock.local/credit-notes/%s.pdf", externalID),
		XMLURL:           fmt.Sprintf("https://softpymes-mock.local/credit-notes/%s.xml", externalID),
		CUFE:             cufe,
		IssuedAt:         time.Now(),
		CreatedAt:        time.Now(),
	}

	s.Repository.SaveCreditNote(creditNote)

	s.logger.Info().
		Str("note_number", creditNoteNumber).
		Str("cufe", cufe).
		Msg("Nota de crédito simulada creada exitosamente")

	return creditNote, nil
}

// HandleListDocuments simula el listado de documentos
func (s *APISimulator) HandleListDocuments(token string, filters map[string]interface{}) ([]domain.Invoice, error) {
	s.logger.Info().
		Str("token", token).
		Interface("filters", filters).
		Msg("Simulando listado de documentos")

	// Validar token
	authToken, exists := s.Repository.GetToken(token)
	if !exists {
		return nil, fmt.Errorf("invalid token")
	}

	// Verificar si expiró
	if time.Now().After(authToken.ExpiresAt) {
		return nil, fmt.Errorf("token expired")
	}

	// Por ahora retornar todas las facturas (simplificado)
	invoices := s.Repository.GetAllInvoices()

	invoiceList := make([]domain.Invoice, len(invoices))
	for i, inv := range invoices {
		invoiceList[i] = *inv
	}

	s.logger.Info().
		Int("count", len(invoiceList)).
		Msg("Documentos listados exitosamente")

	return invoiceList, nil
}
