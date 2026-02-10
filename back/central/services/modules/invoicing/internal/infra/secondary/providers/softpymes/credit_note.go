package softpymes

import (
	"context"
	"fmt"

<<<<<<< HEAD
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/ports"
=======
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/dtos"
>>>>>>> 7b7c2054fa8e6cf0840b58d299ba6b7ca4e6b49e
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/infra/secondary/providers/softpymes/mappers"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/infra/secondary/providers/softpymes/response"
)

// CreateCreditNote crea una nota de crédito en Softpymes
<<<<<<< HEAD
func (c *Client) CreateCreditNote(ctx context.Context, token string, request *ports.CreditNoteRequest) (*ports.CreditNoteResponse, error) {
=======
func (c *Client) CreateCreditNote(ctx context.Context, token string, request *dtos.CreditNoteRequest) (*dtos.CreditNoteResponse, error) {
>>>>>>> 7b7c2054fa8e6cf0840b58d299ba6b7ca4e6b49e
	c.log.Info(ctx).
		Uint("invoice_id", request.Invoice.ID).
		Str("note_type", request.CreditNote.NoteType).
		Msg("Creating credit note in Softpymes")

	// Convertir request a formato Softpymes
	softpymesReq := mappers.ToCreditNoteRequest(request)

	var noteResp response.CreditNoteResponse

	// Hacer llamado a la API
	resp, err := c.httpClient.R().
		SetContext(ctx).
		SetAuthToken(token).
		SetBody(softpymesReq).
		SetResult(&noteResp).
		Post("/search/documents/notes/")

	if err != nil {
		c.log.Error(ctx).Err(err).Msg("Failed to create credit note")
		return nil, fmt.Errorf("credit note creation request failed: %w", err)
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
			return nil, fmt.Errorf("authentication token expired")
		}

		return nil, fmt.Errorf("credit note creation failed: %s", noteResp.Error)
	}

	// Verificar respuesta
	if !noteResp.Success {
		c.log.Error(ctx).
			Str("message", noteResp.Message).
			Msg("Credit note creation unsuccessful")
		return nil, fmt.Errorf("credit note creation unsuccessful: %s", noteResp.Message)
	}

	// Convertir respuesta a formato de dominio
	result := mappers.FromCreditNoteResponse(&noteResp)
	if result == nil {
		return nil, fmt.Errorf("failed to parse credit note response")
	}

	c.log.Info(ctx).
		Str("note_number", result.CreditNoteNumber).
		Str("cufe", *result.CUFE).
		Msg("Credit note created successfully in Softpymes")

	return result, nil
}
