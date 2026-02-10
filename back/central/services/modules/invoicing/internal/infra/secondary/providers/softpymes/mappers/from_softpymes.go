package mappers

import (
<<<<<<< HEAD
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/ports"
=======
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/dtos"
>>>>>>> 7b7c2054fa8e6cf0840b58d299ba6b7ca4e6b49e
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/infra/secondary/providers/softpymes/response"
)

// FromInvoiceResponse convierte una respuesta de Softpymes a InvoiceResponse de dominio
<<<<<<< HEAD
func FromInvoiceResponse(resp *response.InvoiceResponse) *ports.InvoiceResponse {
=======
func FromInvoiceResponse(resp *response.InvoiceResponse) *dtos.InvoiceResponse {
>>>>>>> 7b7c2054fa8e6cf0840b58d299ba6b7ca4e6b49e
	if resp == nil || resp.Data == nil {
		return nil
	}

	issuedAt := resp.Data.IssuedAt.Format("2006-01-02T15:04:05Z07:00")

<<<<<<< HEAD
	return &ports.InvoiceResponse{
=======
	return &dtos.InvoiceResponse{
>>>>>>> 7b7c2054fa8e6cf0840b58d299ba6b7ca4e6b49e
		InvoiceNumber: resp.Data.InvoiceNumber,
		ExternalID:    resp.Data.InvoiceID,
		InvoiceURL:    &resp.Data.InvoiceURL,
		PDFURL:        &resp.Data.PDFURL,
		XMLURL:        &resp.Data.XMLURL,
		CUFE:          &resp.Data.CUFE,
		IssuedAt:      issuedAt,
		RawResponse: map[string]interface{}{
			"invoice_id":     resp.Data.InvoiceID,
			"invoice_number": resp.Data.InvoiceNumber,
			"cufe":           resp.Data.CUFE,
			"status":         resp.Data.Status,
			"qr_code":        resp.Data.QRCode,
		},
	}
}

// FromCreditNoteResponse convierte una respuesta de Softpymes a CreditNoteResponse de dominio
<<<<<<< HEAD
func FromCreditNoteResponse(resp *response.CreditNoteResponse) *ports.CreditNoteResponse {
=======
func FromCreditNoteResponse(resp *response.CreditNoteResponse) *dtos.CreditNoteResponse {
>>>>>>> 7b7c2054fa8e6cf0840b58d299ba6b7ca4e6b49e
	if resp == nil || resp.Data == nil {
		return nil
	}

	issuedAt := resp.Data.IssuedAt.Format("2006-01-02T15:04:05Z07:00")

<<<<<<< HEAD
	return &ports.CreditNoteResponse{
=======
	return &dtos.CreditNoteResponse{
>>>>>>> 7b7c2054fa8e6cf0840b58d299ba6b7ca4e6b49e
		CreditNoteNumber: resp.Data.NoteNumber,
		ExternalID:       resp.Data.NoteID,
		NoteURL:          &resp.Data.NoteURL,
		PDFURL:           &resp.Data.PDFURL,
		XMLURL:           &resp.Data.XMLURL,
		CUFE:             &resp.Data.CUFE,
		IssuedAt:         issuedAt,
		RawResponse: map[string]interface{}{
			"note_id":     resp.Data.NoteID,
			"note_number": resp.Data.NoteNumber,
			"cufe":        resp.Data.CUFE,
			"status":      resp.Data.Status,
		},
	}
}
