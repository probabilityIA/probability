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

	// Extraer datos
	customer, _ := invoiceData["customer"].(map[string]interface{})
	items, _ := invoiceData["items"].([]interface{})
	total, _ := invoiceData["total"].(float64)
	orderID, _ := invoiceData["order_id"].(string)

	// Generar factura
	invoiceNumber := s.Repository.GenerateInvoiceNumber()
	externalID := uuid.New().String()
	cufe := fmt.Sprintf("CUFE-%s", uuid.New().String()[:16])

	// Convertir items
	invoiceItems := make([]domain.InvoiceItem, 0)
	for _, item := range items {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		description := ""
		if d, ok := itemMap["description"].(string); ok {
			description = d
		}

		quantity := 1
		if q, ok := itemMap["quantity"].(float64); ok {
			quantity = int(q)
		}

		unitPrice := 0.0
		if p, ok := itemMap["unit_price"].(float64); ok {
			unitPrice = p
		}

		tax := 0.0
		if t, ok := itemMap["tax"].(float64); ok {
			tax = t
		}

		itemTotal := 0.0
		if t, ok := itemMap["total"].(float64); ok {
			itemTotal = t
		}

		invoiceItems = append(invoiceItems, domain.InvoiceItem{
			ItemCode:    fmt.Sprintf("ITEM-%d", len(invoiceItems)+1),
			Description: description,
			Quantity:    quantity,
			UnitPrice:   unitPrice,
			Tax:         tax,
			Total:       itemTotal,
		})
	}

	customerName := ""
	if name, ok := customer["name"].(string); ok {
		customerName = name
	}

	customerEmail := ""
	if email, ok := customer["email"].(string); ok {
		customerEmail = email
	}

	customerNIT := "999999999"
	if nit, ok := customer["nit"].(string); ok {
		customerNIT = nit
	} else if dni, ok := customer["dni"].(string); ok {
		customerNIT = dni
	}

	invoice := &domain.Invoice{
		ID:            externalID,
		InvoiceNumber: invoiceNumber,
		ExternalID:    externalID,
		OrderID:       orderID,
		CustomerName:  customerName,
		CustomerEmail: customerEmail,
		CustomerNIT:   customerNIT,
		Total:         total,
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
		Str("cufe", cufe).
		Msg("Factura simulada creada exitosamente")

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
