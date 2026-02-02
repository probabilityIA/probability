package usecases

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/secamc93/probability/back/integrationTest/integrations/softpymes/internal/domain"
	"github.com/secamc93/probability/back/integrationTest/shared/log"
)

// APISimulator simula el API de SoftPymes
type APISimulator struct {
	logger     log.ILogger
	Repository *domain.InvoiceRepository // Exportado para acceso desde bundle
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
func (s *APISimulator) HandleCreateInvoice(token string, invoiceData map[string]interface{}) (*domain.Invoice, error) {
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
		invoiceItems = append(invoiceItems, domain.InvoiceItem{
			Description: itemMap["description"].(string),
			Quantity:    int(itemMap["quantity"].(float64)),
			UnitPrice:   itemMap["unit_price"].(float64),
			Tax:         itemMap["tax"].(float64),
			Total:       itemMap["total"].(float64),
		})
	}

	customerName, _ := customer["name"].(string)
	customerEmail, _ := customer["email"].(string)
	customerNIT, _ := customer["nit"].(string)

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

	return invoice, nil
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
