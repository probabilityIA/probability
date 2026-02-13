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
// Valida el formato correcto según documentación de Softpymes:
// https://api-integracion.softpymes.com.co/doc/#api-Documentos-PostSaleInvoice
func (s *APISimulator) HandleCreateInvoice(token string, invoiceData map[string]interface{}) (*domain.Invoice, error) {
	s.logger.Info().
		Str("token", token).
		Interface("data", invoiceData).
		Msg("Simulando creación de factura con formato Softpymes")

	// Validar token
	authToken, exists := s.Repository.GetToken(token)
	if !exists {
		return nil, fmt.Errorf("invalid token")
	}

	// Verificar si expiró
	if time.Now().After(authToken.ExpiresAt) {
		return nil, fmt.Errorf("token expired")
	}

	// ====================================================
	// VALIDAR FORMATO SEGÚN DOCUMENTACIÓN DE SOFTPYMES
	// ====================================================

	// Campos requeridos según documentación
	currencyCode, ok := invoiceData["currencyCode"].(string)
	if !ok || currencyCode == "" {
		return nil, fmt.Errorf("missing required field: currencyCode")
	}

	items, ok := invoiceData["items"].([]interface{})
	if !ok || len(items) == 0 {
		return nil, fmt.Errorf("missing required field: items")
	}

	// Campos opcionales con defaults
	customerNit := ""
	if nit, ok := invoiceData["customerNit"].(string); ok {
		customerNit = nit
	}

	branchCode := "001"
	if branch, ok := invoiceData["branchCode"].(string); ok {
		branchCode = branch
	}

	sellerNit := ""
	if seller, ok := invoiceData["sellerNit"].(string); ok {
		sellerNit = seller
	}

	s.logger.Info().
		Str("currency", currencyCode).
		Str("customer_nit", customerNit).
		Str("seller_nit", sellerNit).
		Str("branch", branchCode).
		Int("items_count", len(items)).
		Msg("Request validado correctamente según formato Softpymes")

	// Calcular total desde items (Softpymes lo calcula en backend)
	total := 0.0
	invoiceItems := make([]domain.InvoiceItem, 0)
	for _, item := range items {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		// Validar campos requeridos del item
		itemCode, ok := itemMap["itemCode"].(string)
		if !ok {
			return nil, fmt.Errorf("missing required field: itemCode in item")
		}

		quantity := 0.0
		if q, ok := itemMap["quantity"].(float64); ok {
			quantity = q
		}

		discount := 0.0
		if d, ok := itemMap["discount"].(float64); ok {
			discount = d
		}

		// Calcular subtotal del item (simplificado)
		itemTotal := quantity * 100.0 // Asumimos precio ficticio
		itemTotal -= discount
		total += itemTotal

		invoiceItems = append(invoiceItems, domain.InvoiceItem{
			Description: itemCode,
			Quantity:    int(quantity),
			UnitPrice:   100.0, // Precio ficticio para el mock
			Tax:         0.0,
			Total:       itemTotal,
		})
	}

	// Generar factura
	invoiceNumber := s.Repository.GenerateInvoiceNumber()
	externalID := uuid.New().String()
	cufe := fmt.Sprintf("CUFE-%s", uuid.New().String()[:16])

	invoice := &domain.Invoice{
		ID:            externalID,
		InvoiceNumber: invoiceNumber,
		ExternalID:    externalID,
		OrderID:       "", // No disponible en formato Softpymes
		CustomerName:  customerNit,
		CustomerEmail: "",
		CustomerNIT:   customerNit,
		Total:         total,
		Currency:      currencyCode,
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
		Float64("total", total).
		Msg("✅ Factura simulada creada exitosamente con formato Softpymes")

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
