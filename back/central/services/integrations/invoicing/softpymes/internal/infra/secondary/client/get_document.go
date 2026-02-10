package client

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// GetDocumentByNumber obtiene un documento específico por su número
// Reutiliza el endpoint de búsqueda (/app/integration/search/documents/)
// filtrando por DocumentNumber
//
// Parámetros:
// - documentNumber: Número del documento a buscar (ej: "ABC0000000000")
//
// Comportamiento:
// - Busca en los últimos 30 días (límite de la API)
// - Retorna el primer documento que coincida como map[string]interface{}
// - Si no encuentra, retorna error
// - Si encuentra múltiples, retorna el primero y logea warning
//
// Uso típico: Después de crear una factura, esperar 3 segundos y consultar
// el documento completo para obtener URLs de PDF/XML y CUFE
//
// Implementa: ports.ISoftpymesClient.GetDocumentByNumber
func (c *Client) GetDocumentByNumber(ctx context.Context, apiKey, apiSecret, referer, documentNumber string) (map[string]interface{}, error) {
	c.log.Info(ctx).
		Str("document_number", documentNumber).
		Msg("📄 Getting document by number from Softpymes")

	// Validar parámetro
	if documentNumber == "" {
		c.log.Error(ctx).Msg("❌ documentNumber is required")
		return nil, fmt.Errorf("documentNumber is required")
	}

	// Preparar rango de fechas: últimos 30 días (máximo permitido por API)
	now := time.Now()
	dateFrom := now.AddDate(0, 0, -30).Format("2006-01-02") // 30 días atrás
	dateTo := now.Format("2006-01-02")                       // Hoy

	// Preparar parámetros de búsqueda
	params := ListDocumentsParams{
		DateFrom:       dateFrom,
		DateTo:         dateTo,
		DocumentNumber: &documentNumber,
	}

	c.log.Info(ctx).
		Str("date_from", dateFrom).
		Str("date_to", dateTo).
		Str("document_number", documentNumber).
		Msg("📤 Searching for document in last 30 days")

	// Llamar al endpoint de lista con filtro de número
	documents, err := c.ListDocuments(ctx, apiKey, apiSecret, referer, params)
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
