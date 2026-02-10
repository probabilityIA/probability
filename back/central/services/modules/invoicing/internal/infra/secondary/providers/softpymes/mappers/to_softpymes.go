package mappers

import (
	"time"

<<<<<<< HEAD
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/ports"
=======
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/dtos"
>>>>>>> 7b7c2054fa8e6cf0840b58d299ba6b7ca4e6b49e
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/infra/secondary/providers/softpymes/request"
)

// ToInvoiceRequest convierte un InvoiceRequest de dominio a request de Softpymes
<<<<<<< HEAD
func ToInvoiceRequest(req *ports.InvoiceRequest) *request.InvoiceRequest {
=======
func ToInvoiceRequest(req *dtos.InvoiceRequest) *request.InvoiceRequest {
>>>>>>> 7b7c2054fa8e6cf0840b58d299ba6b7ca4e6b49e
	// Extraer configuración del proveedor
	referer := ""
	branchCode := "001" // Por defecto
	if req.Provider.Config != nil {
		if ref, ok := req.Provider.Config["referer"].(string); ok {
			referer = ref
		}
		if branch, ok := req.Provider.Config["branch_code"].(string); ok {
			branchCode = branch
		}
	}

	// Mapear items
	items := make([]request.InvoiceItem, 0, len(req.InvoiceItems))
	for _, item := range req.InvoiceItems {
		softpymesItem := request.InvoiceItem{
			Code:        item.SKU,
			Description: item.Name,
			Quantity:    float64(item.Quantity),
			UnitPrice:   item.UnitPrice,
			TotalPrice:  item.TotalPrice,
			Tax:         item.Tax,
			Discount:    item.Discount,
		}

		if item.TaxRate != nil {
			softpymesItem.TaxRate = *item.TaxRate
		}

		items = append(items, softpymesItem)
	}

	// Determinar tipo de identificación
	identificationType := "CC" // Cédula por defecto
	if len(req.Invoice.CustomerDNI) > 10 {
		identificationType = "NIT" // Si es más largo, probablemente es NIT
	}

	return &request.InvoiceRequest{
		Referer:    referer,
		BranchCode: branchCode,
		Date:       time.Now(),
		Customer: request.CustomerData{
			IdentificationType:   identificationType,
			IdentificationNumber: req.Invoice.CustomerDNI,
			Name:                 req.Invoice.CustomerName,
			Email:                req.Invoice.CustomerEmail,
			Phone:                req.Invoice.CustomerPhone,
		},
		Items:         items,
		Currency:      req.Invoice.Currency,
		Notes:         req.Invoice.Notes,
		PaymentMethod: "Contado", // Por defecto
	}
}

// ToCreditNoteRequest convierte un CreditNoteRequest de dominio a request de Softpymes
<<<<<<< HEAD
func ToCreditNoteRequest(req *ports.CreditNoteRequest) *request.CreditNoteRequest {
=======
func ToCreditNoteRequest(req *dtos.CreditNoteRequest) *request.CreditNoteRequest {
>>>>>>> 7b7c2054fa8e6cf0840b58d299ba6b7ca4e6b49e
	// Extraer configuración del proveedor
	referer := ""
	if req.Provider.Config != nil {
		if ref, ok := req.Provider.Config["referer"].(string); ok {
			referer = ref
		}
	}

	// Mapear tipo de nota
	noteType := "03" // Devolución por defecto
	switch req.CreditNote.NoteType {
	case "cancellation":
		noteType = "01" // Anulación
	case "correction":
		noteType = "02" // Corrección
	case "full_refund", "partial_refund":
		noteType = "03" // Devolución
	}

	description := req.CreditNote.Reason
	if req.CreditNote.Description != nil {
		description = *req.CreditNote.Description
	}

	return &request.CreditNoteRequest{
		InvoiceNumber: req.Invoice.InvoiceNumber,
		Referer:       referer,
		Amount:        req.CreditNote.Amount,
		Reason:        req.CreditNote.Reason,
		Description:   description,
		NoteType:      noteType,
	}
}
