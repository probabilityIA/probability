package mappers

import (
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/infra/primary/handlers/response"
)

// InvoiceToResponse convierte entidad de dominio a response
func InvoiceToResponse(invoice *entities.Invoice, includeItems bool) *response.Invoice {
	// Convert *uint to uint (dereference with default 0 if nil)
	var invoicingProviderID uint
	if invoice.InvoicingProviderID != nil {
		invoicingProviderID = *invoice.InvoicingProviderID
	}

	resp := &response.Invoice{
		ID:                  invoice.ID,
		CreatedAt:           invoice.CreatedAt,
		UpdatedAt:           invoice.UpdatedAt,
		OrderID:             invoice.OrderID,
		BusinessID:          invoice.BusinessID,
		InvoicingProviderID: invoicingProviderID,
		InvoiceNumber:       invoice.InvoiceNumber,
		InternalNumber:      invoice.InternalNumber,
		ExternalID:          invoice.ExternalID,
		Status:              invoice.Status,
		TotalAmount:         invoice.TotalAmount,
		Subtotal:            invoice.Subtotal,
		Tax:                 invoice.Tax,
		Discount:            invoice.Discount,
		Currency:            invoice.Currency,
		IssuedAt:            invoice.IssuedAt,
		CancelledAt:         invoice.CancelledAt,
		CUFE:                invoice.CUFE,
		PDFURL:              invoice.PDFURL,
		XMLURL:              invoice.XMLURL,
		Metadata:            invoice.Metadata,
	}

	// Incluir items si se solicita
	if includeItems && invoice.Items != nil {
		resp.Items = make([]response.InvoiceItem, 0, len(invoice.Items))
		for _, item := range invoice.Items {
			resp.Items = append(resp.Items, InvoiceItemToResponse(&item))
		}
	}

	return resp
}

// InvoiceItemToResponse convierte item de dominio a response
func InvoiceItemToResponse(item *entities.InvoiceItem) response.InvoiceItem {
	sku := item.SKU
	return response.InvoiceItem{
		ID:          item.ID,
		ProductSKU:  &sku,
		ProductName: item.Name,
		Description: item.Description,
		Quantity:    item.Quantity,
		UnitPrice:   item.UnitPrice,
		TotalPrice:  item.TotalPrice,
		Tax:         item.Tax,
		TaxRate:     item.TaxRate,
		Discount:    item.Discount,
	}
}

// InvoicesToResponse convierte lista de entidades a response
func InvoicesToResponse(invoices []*entities.Invoice, totalCount int64, page, pageSize int) *response.InvoiceList {
	items := make([]response.Invoice, 0, len(invoices))
	for _, invoice := range invoices {
		items = append(items, *InvoiceToResponse(invoice, false)) // No incluir items en listado
	}

	return &response.InvoiceList{
		Items:      items,
		TotalCount: totalCount,
		Page:       page,
		PageSize:   pageSize,
	}
}

// ProviderToResponse convierte entidad de dominio a response
func ProviderToResponse(provider *entities.InvoicingProvider) *response.Provider {
	// Ofuscar credenciales (mostrar solo primeros 4 caracteres de keys)
	credentials := make(map[string]interface{})
	if provider.Credentials != nil {
		for k, v := range provider.Credentials {
			if str, ok := v.(string); ok && len(str) > 4 {
				credentials[k] = str[:4] + "****"
			} else {
				credentials[k] = "****"
			}
		}
	}

	return &response.Provider{
		ID:               provider.ID,
		CreatedAt:        provider.CreatedAt,
		UpdatedAt:        provider.UpdatedAt,
		Name:             provider.Name,
		ProviderTypeCode: "", // TODO: Obtener del join con provider_type
		BusinessID:       provider.BusinessID,
		Config:           provider.Config,
		Credentials:      credentials,
		IsActive:         provider.IsActive,
	}
}

// ProvidersToResponse convierte lista de entidades a response
func ProvidersToResponse(providers []*entities.InvoicingProvider, totalCount int64, page, pageSize int) *response.ProviderList {
	items := make([]response.Provider, 0, len(providers))
	for _, provider := range providers {
		items = append(items, *ProviderToResponse(provider))
	}

	return &response.ProviderList{
		Items:      items,
		TotalCount: totalCount,
		Page:       page,
		PageSize:   pageSize,
	}
}

// ConfigToResponse convierte entidad de dominio a response
func ConfigToResponse(config *entities.InvoicingConfig) *response.Config {
	// Convert *uint to uint (dereference with default 0 if nil)
	var invoicingProviderID uint
	if config.InvoicingProviderID != nil {
		invoicingProviderID = *config.InvoicingProviderID
	}

	return &response.Config{
		ID:                  config.ID,
		CreatedAt:           config.CreatedAt,
		UpdatedAt:           config.UpdatedAt,
		BusinessID:          config.BusinessID,
		IntegrationID:       config.IntegrationID,
		InvoicingProviderID: invoicingProviderID,
		Enabled:             config.Enabled,
		AutoInvoice:         config.AutoInvoice,
		Filters:             config.Filters,
	}
}

// ConfigsToResponse convierte lista de entidades a response
func ConfigsToResponse(configs []*entities.InvoicingConfig, totalCount int64, page, pageSize int) *response.ConfigList {
	items := make([]response.Config, 0, len(configs))
	for _, config := range configs {
		items = append(items, *ConfigToResponse(config))
	}

	return &response.ConfigList{
		Items:      items,
		TotalCount: totalCount,
		Page:       page,
		PageSize:   pageSize,
	}
}

// CreditNoteToResponse convierte entidad de dominio a response
func CreditNoteToResponse(note *entities.CreditNote) *response.CreditNote {
	return &response.CreditNote{
		ID:               note.ID,
		CreatedAt:        note.CreatedAt,
		UpdatedAt:        note.UpdatedAt,
		InvoiceID:        note.InvoiceID,
		BusinessID:       note.BusinessID,
		CreditNoteNumber: note.CreditNoteNumber,
		ExternalID:       note.ExternalID,
		Amount:           note.Amount,
		NoteType:         note.NoteType,
		Reason:           note.Reason,
		Status:           note.Status,
		IssuedAt:         note.IssuedAt,
		CUFE:             note.CUFE,
		PDFURL:           note.PDFURL,
		XMLURL:           note.XMLURL,
		Metadata:         note.Metadata,
	}
}
