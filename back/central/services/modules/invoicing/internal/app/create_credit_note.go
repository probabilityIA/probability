package app

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/entities"
)

// CreateCreditNote crea una nota de crédito para una factura
// DEPRECATED: Esta funcionalidad fue migrada a integrations/invoicing/softpymes/
// TODO: Re-implementar usando integrationCore si es necesario
func (uc *useCase) CreateCreditNote(ctx context.Context, dto *dtos.CreateCreditNoteDTO) (*entities.CreditNote, error) {
	return nil, fmt.Errorf("CreateCreditNote is deprecated and was moved to softpymes integration")
}

/* IMPLEMENTACIÓN ORIGINAL (requiere providerRepo y providerClient que ya no existen):

func (uc *useCase) CreateCreditNote(ctx context.Context, dto *dtos.CreateCreditNoteDTO) (*entities.CreditNote, error) {
	uc.log.Info(ctx).Uint("invoice_id", dto.InvoiceID).Msg("Creating credit note")

	// 1. Obtener factura
	invoice, err := uc.invoiceRepo.GetByID(ctx, dto.InvoiceID)
	if err != nil {
		return nil, errors.ErrInvoiceNotFound
	}

	// 2. Validar que la factura esté emitida
	if invoice.Status != constants.InvoiceStatusIssued {
		return nil, errors.ErrInvoiceNotIssued
	}

	// 3. Validar que el monto no exceda el total
	if dto.Amount > invoice.TotalAmount {
		return nil, errors.ErrCreditNoteAmountExceeds
	}

	// 4. Obtener proveedor
	provider, err := uc.providerRepo.GetByID(ctx, invoice.InvoicingProviderID)
	if err != nil {
		return nil, errors.ErrProviderNotFound
	}

	// 5. Desencriptar credenciales
	credentials, err := uc.encryption.Decrypt(provider.Credentials)
	if err != nil {
		return nil, errors.ErrDecryptionFailed
	}

	// 6. Autenticar
	token, err := uc.providerClient.Authenticate(ctx, credentials)
	if err != nil {
		return nil, errors.ErrAuthenticationFailed
	}

	// 7. Crear entidad de nota de crédito
	creditNote := &entities.CreditNote{
		InvoiceID:      invoice.ID,
		BusinessID:     invoice.BusinessID,
		NoteType:       dto.NoteType,
		Amount:         dto.Amount,
		Currency:       invoice.Currency,
		Reason:         dto.Reason,
		Description:    dto.Description,
		Status:         constants.CreditNoteStatusPending,
		Metadata:       make(map[string]interface{}),
		CreatedByID:    dto.CreatedByUserID,
	}

	// 8. Guardar en BD
	if err := uc.creditNoteRepo.Create(ctx, creditNote); err != nil {
		uc.log.Error(ctx).Err(err).Msg("Failed to create credit note")
		return nil, err
	}

	// 9. Crear nota en el proveedor
	request := &ports.CreditNoteRequest{
		Invoice:    invoice,
		CreditNote: creditNote,
		Provider:   provider,
	}

	response, err := uc.providerClient.CreateCreditNote(ctx, token, request)
	if err != nil {
		uc.log.Error(ctx).Err(err).Msg("Failed to create credit note with provider")
		creditNote.Status = constants.CreditNoteStatusFailed
		uc.creditNoteRepo.Update(ctx, creditNote)
		return nil, errors.ErrProviderAPIError
	}

	// 10. Actualizar nota con datos del proveedor
	creditNote.CreditNoteNumber = response.CreditNoteNumber
	creditNote.ExternalID = &response.ExternalID
	creditNote.NoteURL = response.NoteURL
	creditNote.PDFURL = response.PDFURL
	creditNote.XMLURL = response.XMLURL
	creditNote.CUFE = response.CUFE
	creditNote.Status = constants.CreditNoteStatusIssued
	creditNote.ProviderResponse = response.RawResponse

	if response.IssuedAt != "" {
		issuedAt, err := time.Parse(time.RFC3339, response.IssuedAt)
		if err == nil {
			creditNote.IssuedAt = &issuedAt
		}
	}

	// 11. Actualizar en BD
	if err := uc.creditNoteRepo.Update(ctx, creditNote); err != nil {
		uc.log.Error(ctx).Err(err).Msg("Failed to update credit note")
	}

	// 12. Publicar evento
	if err := uc.eventPublisher.PublishCreditNoteCreated(ctx, creditNote); err != nil {
		uc.log.Error(ctx).Err(err).Msg("Failed to publish event")
	}

	uc.log.Info(ctx).Uint("note_id", creditNote.ID).Msg("Credit note created successfully")
	return creditNote, nil
}

// GetCreditNote obtiene una nota de crédito por ID
func (uc *useCase) GetCreditNote(ctx context.Context, id uint) (*entities.CreditNote, error) {
	return uc.creditNoteRepo.GetByID(ctx, id)
}

func (uc *useCase) ListCreditNotes(ctx context.Context, filters map[string]interface{}) ([]*entities.CreditNote, error) {
	return uc.creditNoteRepo.List(ctx, filters)
}
*/

// GetCreditNote obtiene una nota de crédito por ID (kept - still uses creditNoteRepo)
func (uc *useCase) GetCreditNote(ctx context.Context, id uint) (*entities.CreditNote, error) {
	return uc.creditNoteRepo.GetByID(ctx, id)
}

// ListCreditNotes lista notas de crédito con filtros (kept - still uses creditNoteRepo)
func (uc *useCase) ListCreditNotes(ctx context.Context, filters map[string]interface{}) ([]*entities.CreditNote, error) {
	return uc.creditNoteRepo.List(ctx, filters)
}
