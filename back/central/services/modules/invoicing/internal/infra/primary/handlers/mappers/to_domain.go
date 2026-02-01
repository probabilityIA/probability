package mappers

import (
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/infra/primary/handlers/request"
)

// CreateInvoiceRequestToDTO convierte request a DTO de dominio
func CreateInvoiceRequestToDTO(req *request.CreateInvoice) *dtos.CreateInvoiceDTO {
	return &dtos.CreateInvoiceDTO{
		OrderID:  req.OrderID,
		IsManual: true, // Siempre true en creación manual
	}
}

// CreateCreditNoteRequestToDTO convierte request a DTO de dominio
func CreateCreditNoteRequestToDTO(req *request.CreateCreditNote) *dtos.CreateCreditNoteDTO {
	return &dtos.CreateCreditNoteDTO{
		InvoiceID: req.InvoiceID,
		Amount:    req.Amount,
		Reason:    req.Reason,
		NoteType:  req.NoteType,
	}
}

// CreateProviderRequestToDTO convierte request a DTO de dominio
func CreateProviderRequestToDTO(req *request.CreateProvider) *dtos.CreateProviderDTO {
	return &dtos.CreateProviderDTO{
		Name:             req.Name,
		ProviderTypeCode: req.ProviderTypeCode,
		BusinessID:       req.BusinessID,
		Config:           req.Config,
		Credentials:      req.Credentials,
	}
}

// UpdateProviderRequestToDTO convierte request a DTO de dominio
func UpdateProviderRequestToDTO(req *request.UpdateProvider) *dtos.UpdateProviderDTO {
	dto := &dtos.UpdateProviderDTO{}

	if req.Name != nil {
		dto.Name = req.Name
	}

	if req.Config != nil {
		dto.Config = *req.Config
	}

	if req.Credentials != nil {
		dto.Credentials = *req.Credentials
	}

	if req.IsActive != nil {
		dto.IsActive = req.IsActive
	}

	return dto
}

// CreateConfigRequestToDTO convierte request a DTO de dominio
func CreateConfigRequestToDTO(req *request.CreateConfig, userID uint) *dtos.CreateConfigDTO {
	enabled := true
	if req.Enabled != nil {
		enabled = *req.Enabled
	}

	autoInvoice := false
	if req.AutoInvoice != nil {
		autoInvoice = *req.AutoInvoice
	}

	var invoicingProviderID uint
	if req.InvoicingProviderID != nil {
		invoicingProviderID = *req.InvoicingProviderID
	}

	return &dtos.CreateConfigDTO{
		BusinessID:             req.BusinessID,
		IntegrationID:          req.IntegrationID,
		InvoicingIntegrationID: req.InvoicingIntegrationID,
		InvoicingProviderID:    invoicingProviderID, // Deprecado pero mantener compatibilidad
		Enabled:                enabled,
		AutoInvoice:            autoInvoice,
		Filters:                req.Filters,
		CreatedByUserID:        userID,
	}
}

// UpdateConfigRequestToDTO convierte request a DTO de dominio
func UpdateConfigRequestToDTO(req *request.UpdateConfig) *dtos.UpdateConfigDTO {
	dto := &dtos.UpdateConfigDTO{}

	if req.Enabled != nil {
		dto.Enabled = req.Enabled
	}

	if req.AutoInvoice != nil {
		dto.AutoInvoice = req.AutoInvoice
	}

	if req.Filters != nil {
		dto.Filters = *req.Filters
	}

	return dto
}
