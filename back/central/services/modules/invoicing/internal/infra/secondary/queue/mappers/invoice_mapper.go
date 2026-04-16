package mappers

import (
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/infra/secondary/queue/messages"
)

// InvoiceToCreatedEvent convierte una entidad Invoice a un evento de factura creada
func InvoiceToCreatedEvent(invoice *entities.Invoice) *messages.InvoiceEvent {
	return &messages.InvoiceEvent{
		EventType:     "invoice.created",
		InvoiceID:     invoice.ID,
		OrderID:       invoice.OrderID,
		BusinessID:    invoice.BusinessID,
		InvoiceNumber: invoice.InvoiceNumber,
		TotalAmount:   invoice.TotalAmount,
		Status:        invoice.Status,
		Timestamp:     invoice.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// InvoiceToCancelledEvent convierte una entidad Invoice a un evento de factura cancelada
func InvoiceToCancelledEvent(invoice *entities.Invoice) *messages.InvoiceEvent {
	return &messages.InvoiceEvent{
		EventType:     "invoice.cancelled",
		InvoiceID:     invoice.ID,
		OrderID:       invoice.OrderID,
		BusinessID:    invoice.BusinessID,
		InvoiceNumber: invoice.InvoiceNumber,
		TotalAmount:   invoice.TotalAmount,
		Status:        invoice.Status,
		Timestamp:     invoice.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// InvoiceToFailedEvent convierte una entidad Invoice a un evento de factura fallida
func InvoiceToFailedEvent(invoice *entities.Invoice, errorMsg string) *messages.InvoiceEvent {
	return &messages.InvoiceEvent{
		EventType:     "invoice.failed",
		InvoiceID:     invoice.ID,
		OrderID:       invoice.OrderID,
		BusinessID:    invoice.BusinessID,
		InvoiceNumber: invoice.InvoiceNumber,
		TotalAmount:   invoice.TotalAmount,
		Status:        invoice.Status,
		ErrorMessage:  errorMsg,
		Timestamp:     invoice.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// CreditNoteToEvent convierte una entidad CreditNote a un evento
func CreditNoteToEvent(creditNote *entities.CreditNote) *messages.CreditNoteEvent {
	return &messages.CreditNoteEvent{
		EventType:        "credit_note.created",
		CreditNoteID:     creditNote.ID,
		InvoiceID:        creditNote.InvoiceID,
		BusinessID:       creditNote.BusinessID,
		CreditNoteNumber: creditNote.CreditNoteNumber,
		Amount:           creditNote.Amount,
		NoteType:         creditNote.NoteType,
		Timestamp:        creditNote.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
