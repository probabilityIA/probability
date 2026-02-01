package client

import (
	"context"
	"fmt"
)

// CreditNoteResponse representa la respuesta de creación de nota de crédito de Softpymes
type CreditNoteResponse struct {
	Success          bool   `json:"success"`
	Message          string `json:"message"`
	Error            string `json:"error"`
	CreditNoteNumber string `json:"credit_note_number"`
	ExternalID       string `json:"external_id"`
	NoteURL          string `json:"note_url"`
	PDFURL           string `json:"pdf_url"`
	XMLURL           string `json:"xml_url"`
	CUFE             string `json:"cufe"`
	IssuedAt         string `json:"issued_at"`
}

// CreateCreditNote crea una nota de crédito en Softpymes
func (c *Client) CreateCreditNote(ctx context.Context, creditNoteData map[string]interface{}) error {
	c.log.Info(ctx).Interface("data", creditNoteData).Msg("Creating credit note in Softpymes")

	// Extraer credenciales del map
	credentials, ok := creditNoteData["credentials"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("credentials not found in credit note data")
	}

	apiKey, ok := credentials["api_key"].(string)
	if !ok || apiKey == "" {
		return fmt.Errorf("api_key not found in credentials")
	}

	apiSecret, ok := credentials["api_secret"].(string)
	if !ok || apiSecret == "" {
		return fmt.Errorf("api_secret not found in credentials")
	}

	// Autenticar
	token, err := c.authenticate(ctx, apiKey, apiSecret)
	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	// Preparar request de nota de crédito (simplificado)
	noteReq := map[string]interface{}{
		"invoice_id": creditNoteData["invoice_id"],
		"amount":     creditNoteData["amount"],
		"reason":     creditNoteData["reason"],
		"note_type":  creditNoteData["note_type"],
	}

	var noteResp CreditNoteResponse

	// Hacer llamado a la API
	resp, err := c.httpClient.R().
		SetContext(ctx).
		SetAuthToken(token).
		SetBody(noteReq).
		SetResult(&noteResp).
		Post("/search/documents/notes/")

	if err != nil {
		c.log.Error(ctx).Err(err).Msg("Failed to create credit note")
		return fmt.Errorf("credit note creation request failed: %w", err)
	}

	// Manejar errores HTTP
	if resp.IsError() {
		c.log.Error(ctx).
			Int("status", resp.StatusCode()).
			Str("error", noteResp.Error).
			Msg("Credit note creation failed")

		// Si es 401, el token expiró
		if resp.StatusCode() == 401 {
			c.tokenCache.Clear()
			return fmt.Errorf("authentication token expired")
		}

		return fmt.Errorf("credit note creation failed: %s", noteResp.Error)
	}

	// Verificar respuesta
	if !noteResp.Success {
		c.log.Error(ctx).
			Str("message", noteResp.Message).
			Msg("Credit note creation unsuccessful")
		return fmt.Errorf("credit note creation unsuccessful: %s", noteResp.Message)
	}

	c.log.Info(ctx).
		Str("note_number", noteResp.CreditNoteNumber).
		Str("cufe", noteResp.CUFE).
		Msg("Credit note created successfully in Softpymes")

	// Actualizar creditNoteData con los datos de respuesta
	creditNoteData["external_id"] = noteResp.ExternalID
	creditNoteData["note_number"] = noteResp.CreditNoteNumber
	creditNoteData["cufe"] = noteResp.CUFE
	creditNoteData["note_url"] = noteResp.NoteURL
	creditNoteData["pdf_url"] = noteResp.PDFURL
	creditNoteData["xml_url"] = noteResp.XMLURL

	return nil
}
