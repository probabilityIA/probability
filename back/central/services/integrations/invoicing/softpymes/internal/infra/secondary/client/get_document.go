package client

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

// GetDocumentByNumber obtiene un documento específico por su número.
// Acepta el formato combinado que retorna la creación (ej: "FEV23") o el número bare ("23").
//
// La API de Softpymes usa formatos distintos según el endpoint:
//   - Creación  → documentNumber: "FEV23"  (prefix + número, sin ceros)
//   - Búsqueda  → documentNumber: "0000000023" (solo número, con ceros), prefix: "FEV" (separado)
//
// Este método parsea el combinado en prefix + número bare antes de llamar al endpoint de búsqueda.
//
// Implementa: ports.ISoftpymesClient.GetDocumentByNumber
func (c *Client) GetDocumentByNumber(ctx context.Context, apiKey, apiSecret, referer, documentNumber, baseURL string) (map[string]interface{}, error) {
	c.log.Info(ctx).
		Str("document_number", documentNumber).
		Msg("📄 Getting document by number from Softpymes")

	// Validar parámetro
	if documentNumber == "" {
		c.log.Error(ctx).Msg("❌ documentNumber is required")
		return nil, fmt.Errorf("documentNumber is required")
	}

	// Parsear el formato combinado que retorna la creación (ej: "FEV26").
	// La API de búsqueda espera:
	//   - documentNumber: "0000000026" (10 dígitos, zero-padded)
	//   - prefix: "FEV" (campo separado)
	//
	// El creation response retorna el número SIN ceros y con el prefix pegado.
	// Ejemplo: creation devuelve "FEV26", búsqueda necesita documentNumber="0000000026" + prefix="FEV"
	prefix := ""
	bareNumber := documentNumber
	for i, ch := range documentNumber {
		if ch >= '0' && ch <= '9' {
			prefix = documentNumber[:i]
			bareNumber = documentNumber[i:]
			break
		}
	}

	// Zero-pad el número bare a 10 dígitos (formato estándar de Softpymes)
	// Ej: "26" → "0000000026"
	paddedNumber := bareNumber
	if n, err := strconv.ParseInt(bareNumber, 10, 64); err == nil {
		paddedNumber = fmt.Sprintf("%010d", n)
	}

	c.log.Info(ctx).
		Str("raw_document_number", documentNumber).
		Str("prefix", prefix).
		Str("bare_number", bareNumber).
		Str("padded_number", paddedNumber).
		Msg("📄 Parsed document number for search")

	// Preparar rango de fechas: últimos 30 días (máximo permitido por API)
	now := time.Now()
	dateFrom := now.AddDate(0, 0, -30).Format("2006-01-02")
	dateTo := now.Format("2006-01-02")

	// Buscar con número padded y prefix separados
	params := ListDocumentsParams{
		DateFrom:       dateFrom,
		DateTo:         dateTo,
		DocumentNumber: &paddedNumber,
	}
	if prefix != "" {
		params.Prefix = &prefix
	}

	c.log.Info(ctx).
		Str("date_from", dateFrom).
		Str("date_to", dateTo).
		Str("document_number", paddedNumber).
		Str("prefix", prefix).
		Msg("📤 Searching for document in last 30 days")

	// Llamar al endpoint de lista con filtro de número
	documents, err := c.listDocuments(ctx, apiKey, apiSecret, referer, params, baseURL)
	if err != nil {
		c.log.Error(ctx).
			Err(err).
			Str("document_number", documentNumber).
			Msg("❌ Failed to search document")
		return nil, fmt.Errorf("failed to search document %s: %w", documentNumber, err)
	}

	// Validar resultado
	if documents == nil || len(*documents) == 0 {
		c.log.Warn(ctx).
			Str("document_number", documentNumber).
			Msg("⚠️ Document not found - may not be processed yet")
		return nil, fmt.Errorf("document %s not found - it may not have been processed by DIAN yet", documentNumber)
	}

	// Si hay múltiples resultados, logear warning (no debería pasar)
	if len(*documents) > 1 {
		c.log.Warn(ctx).
			Str("document_number", documentNumber).
			Int("count", len(*documents)).
			Msg("⚠️ Multiple documents found with same number - using first one")
	}

	// Obtener el primer documento
	document := (*documents)[0]

	c.log.Info(ctx).
		Str("document_number", document.DocumentNumber).
		Str("document_date", document.DocumentDate).
		Str("customer_name", document.CustomerName).
		Str("total", document.Total).
		Msg("✅ Document retrieved successfully")

	// Convertir Document a map[string]interface{} para mantener consistencia
	// con otros métodos del cliente (CreateInvoice, etc.)
	var documentMap map[string]interface{}
	documentBytes, err := json.Marshal(document)
	if err != nil {
		c.log.Error(ctx).
			Err(err).
			Msg("❌ Failed to marshal document to JSON")
		return nil, fmt.Errorf("failed to marshal document: %w", err)
	}

	if err := json.Unmarshal(documentBytes, &documentMap); err != nil {
		c.log.Error(ctx).
			Err(err).
			Msg("❌ Failed to unmarshal document to map")
		return nil, fmt.Errorf("failed to unmarshal document: %w", err)
	}

	// Documento convertido exitosamente
	return documentMap, nil
}
