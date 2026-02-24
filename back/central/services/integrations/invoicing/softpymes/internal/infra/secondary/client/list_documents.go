package client

import (
	"context"
	"encoding/json"
	"fmt"
)

// Document representa un documento (factura/nota crédito) de Softpymes
// Estructura basada en la documentación oficial:
// https://api-integracion.softpymes.com.co/doc/#api-Documentos-GetSearchDocument
type Document struct {
	BranchCode             string            `json:"branchCode"`
	BranchName             string            `json:"branchName"`
	Comment                string            `json:"comment"`
	CustomerIdentification string            `json:"customerIdentification"` // NIT/CC del cliente
	CustomerName           string            `json:"customerName"`
	Details                []DocumentDetail  `json:"details"`
	DocumentDate           string            `json:"documentDate"` // Formato: string en respuesta
	DocumentName           string            `json:"documentName"` // Tipo de documento (ej: "Factura de Venta")
	DocumentNumber         string            `json:"documentNumber"`
	DueDate                string            `json:"dueDate"`
	PaymentTerm            string            `json:"paymentTerm"`
	Prefix                 string            `json:"prefix"`
	Seller                 DocumentSeller    `json:"seller"`
	ShipInformation        ShipInformation   `json:"shipInformation"`
	TermDays               int               `json:"termDays"`
	Total                  string            `json:"total"`           // Viene como string en la API
	TotalDiscount          string            `json:"totalDiscount"`   // Viene como string en la API
	TotalIva               string            `json:"totalIva"`        // Viene como string en la API
	TotalWithholdingTax    string            `json:"totalWithholdingTax"` // Viene como string en la API
}

// DocumentDetail representa el detalle de un ítem en el documento
type DocumentDetail struct {
	Discount       string            `json:"discount"`
	ItemCode       string            `json:"itemCode"`
	ItemName       string            `json:"itemName"`
	Code           string            `json:"code"`
	Service        string            `json:"service"`
	Iva            string            `json:"iva"`
	Ica            string            `json:"ica"` // Solo para Facturas de Servicios Profesionales
	Quantity       string            `json:"quantity"`
	SizeColor      map[string]string `json:"sizeColor"`
	Value          string            `json:"value"`
	WithholdingTax string            `json:"withholdingTax"`
	Warehouse      DocumentWarehouse `json:"warehouse"`
}

// DocumentWarehouse representa la bodega asociada a un ítem
type DocumentWarehouse struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

// DocumentSeller representa el vendedor del documento
type DocumentSeller struct {
	Name string `json:"name"`
	Nit  string `json:"nit"`
}

// ShipInformation representa la información de envío
type ShipInformation struct {
	ShipAddress    string `json:"shipAddress"`
	ShipCity       string `json:"shipCity"`
	ShipCountry    string `json:"shipCountry"`
	ShipDepartment string `json:"shipDepartment"`
	ShipPhone      string `json:"shipPhone"`
	ShipTo         string `json:"shipTo"`
	ShipZipCode    string `json:"shipZipCode"`
}

// ListDocumentsParams parámetros para filtrar documentos
// Documentación: https://api-integracion.softpymes.com.co/doc/#api-Documentos-GetSearchDocument
type ListDocumentsParams struct {
	DateFrom       string  `json:"dateFrom"`       // REQUERIDO - Formato: YYYY-MM-DD
	DateTo         string  `json:"dateTo"`         // REQUERIDO - Formato: YYYY-MM-DD (máx 30 días desde dateFrom)
	DocumentType   *string `json:"documentType,omitempty"`   // OPCIONAL - Tipo de documento
	DocumentNumber *string `json:"documentNumber,omitempty"` // OPCIONAL - Número documento
	Prefix         *string `json:"prefix,omitempty"`         // OPCIONAL - Prefijo documento
	BranchCode     *string `json:"branchCode,omitempty"`     // OPCIONAL - Código de sucursal
	Page           *string `json:"page,omitempty"`           // OPCIONAL - Número de página (para paginación)
	PageSize       *string `json:"pageSize,omitempty"`       // OPCIONAL - Registros por página
}

// ListDocumentsResponse respuesta de la lista de documentos
// La API retorna un array de documentos directamente (no un objeto con metadata)
type ListDocumentsResponse []Document

// ListDocuments obtiene la lista de documentos de Softpymes
// Documentación: https://api-integracion.softpymes.com.co/doc/#api-Documentos-GetSearchDocument
// Endpoint: POST /app/integration/search/documents/
// IMPORTANTE:
// - dateFrom y dateTo son REQUERIDOS (formato YYYY-MM-DD)
// - El rango máximo entre fechas es de 30 días
// - La respuesta es un array de documentos directamente (no un objeto con metadata)
// ListDocuments obtiene la lista de documentos de Softpymes
// baseURL: URL base efectiva (producción o testing); vacío usa c.baseURL.
func (c *Client) ListDocuments(ctx context.Context, apiKey, apiSecret, referer string, params ListDocumentsParams, baseURL string) (*ListDocumentsResponse, error) {
	c.log.Info(ctx).
		Interface("params", params).
		Msg("📋 Listing documents from Softpymes")

	// Validar parámetros requeridos
	if params.DateFrom == "" || params.DateTo == "" {
		c.log.Error(ctx).Msg("❌ dateFrom and dateTo are required")
		return nil, fmt.Errorf("dateFrom and dateTo are required parameters")
	}

	// Autenticar usando la URL efectiva
	token, err := c.authenticate(ctx, apiKey, apiSecret, referer, baseURL)
	if err != nil {
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	c.log.Info(ctx).
		Str("dateFrom", params.DateFrom).
		Str("dateTo", params.DateTo).
		Msg("📤 Sending list documents request")

	var listResp ListDocumentsResponse

	// Endpoint confirmado según documentación oficial
	resp, err := c.httpClient.R().
		SetContext(ctx).
		SetAuthToken(token).
		SetHeader("Referer", referer).
		SetHeader("Content-Type", "application/json").
		SetBody(params).
		SetResult(&listResp).
		SetDebug(true).
		Post(c.resolveURL(baseURL, "/app/integration/search/documents/"))

	if err != nil {
		c.log.Error(ctx).Err(err).Msg("❌ Failed to list documents")
		return nil, fmt.Errorf("list documents request failed: %w", err)
	}

	c.log.Info(ctx).
		Int("status_code", resp.StatusCode()).
		Str("status", resp.Status()).
		Msg("📥 Received list documents response")

	// Manejar errores HTTP
	if resp.IsError() {
		// Intentar parsear mensaje de error
		// La API puede retornar: {"message": "...", "type": "INVALID_DATA"}
		var errorBody map[string]interface{}
		if err := json.Unmarshal(resp.Body(), &errorBody); err == nil {
			if msg, ok := errorBody["message"].(string); ok {
				errorType := errorBody["type"]
				c.log.Error(ctx).
					Int("status", resp.StatusCode()).
					Str("error", msg).
					Interface("type", errorType).
					Msg("❌ List documents failed")
				return nil, fmt.Errorf("list documents failed (status %d): %s", resp.StatusCode(), msg)
			}
		}

		c.log.Error(ctx).
			Int("status", resp.StatusCode()).
			Str("body", resp.String()).
			Msg("❌ List documents failed - unknown error")

		return nil, fmt.Errorf("list documents failed (status %d): %s", resp.StatusCode(), resp.Status())
	}

	c.log.Info(ctx).
		Int("documents_count", len(listResp)).
		Msg("✅ Documents retrieved successfully")

	return &listResp, nil
}
