package mappers

import (
	"strings"

	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/infra/primary/handlers/response"
)

// buildImageURL construye la URL completa de una imagen desde su ruta relativa
func buildImageURL(relativePath string, baseURL string, fallbackBucket string) string {
	if relativePath == "" {
		return ""
	}

	// Si ya es una URL completa, retornarla directamente
	if strings.HasPrefix(relativePath, "http://") || strings.HasPrefix(relativePath, "https://") {
		return relativePath
	}

	// Usar URL base de S3 desde config
	if baseURL != "" {
		return strings.TrimRight(baseURL, "/") + "/" + strings.TrimLeft(relativePath, "/")
	}

	// Fallback a formato por defecto de S3
	if fallbackBucket != "" {
		return "https://" + fallbackBucket + ".s3.amazonaws.com/" + strings.TrimLeft(relativePath, "/")
	}

	return relativePath
}

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
		CustomerName:        invoice.CustomerName,
		CustomerEmail:       invoice.CustomerEmail,
		IssuedAt:            invoice.IssuedAt,
		CancelledAt:         invoice.CancelledAt,
		CUFE:                invoice.CUFE,
		PDFURL:              invoice.PDFURL,
		XMLURL:              invoice.XMLURL,
		Metadata:            invoice.Metadata,
		ProviderResponse:    invoice.ProviderResponse, // Incluir respuesta completa del proveedor
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
func ConfigToResponse(config *entities.InvoicingConfig, baseURL string, bucket string) *response.Config {
	// Convert *uint to uint (dereference with default 0 if nil)
	var invoicingProviderID uint
	if config.InvoicingProviderID != nil {
		invoicingProviderID = *config.InvoicingProviderID
	}

	resp := &response.Config{
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

	// Incluir nombres de relaciones si están disponibles
	if config.IntegrationName != nil {
		resp.IntegrationName = config.IntegrationName
	}

	if config.ProviderName != nil {
		resp.ProviderName = config.ProviderName
	}

	if config.ProviderImageURL != nil {
		fullImageURL := buildImageURL(*config.ProviderImageURL, baseURL, bucket)
		resp.ProviderImageURL = &fullImageURL
	}

	if config.Description != "" {
		resp.Description = &config.Description
	}

	return resp
}

// ConfigsToResponse convierte lista de entidades a response
func ConfigsToResponse(configs []*entities.InvoicingConfig, totalCount int64, page, pageSize int, baseURL string, bucket string) *response.ConfigList {
	items := make([]response.Config, 0, len(configs))
	for _, config := range configs {
		items = append(items, *ConfigToResponse(config, baseURL, bucket))
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

// SyncLogToResponse convierte entidad de dominio a response
func SyncLogToResponse(log *entities.InvoiceSyncLog) response.SyncLog {
	return response.SyncLog{
		ID:             log.ID,
		InvoiceID:      log.InvoiceID,
		OperationType:  log.OperationType,
		Status:         log.Status,
		ErrorMessage:   log.ErrorMessage,
		ErrorCode:      log.ErrorCode,
		RetryCount:     log.RetryCount,
		MaxRetries:     log.MaxRetries,
		NextRetryAt:    log.NextRetryAt,
		TriggeredBy:    log.TriggeredBy,
		Duration:       log.Duration,
		StartedAt:      log.StartedAt,
		CompletedAt:    log.CompletedAt,
		CreatedAt:      log.CreatedAt,
		RequestPayload: log.RequestPayload,
		RequestURL:     log.RequestURL,
		ResponseStatus: log.ResponseStatus,
		ResponseBody:   log.ResponseBody,
	}
}

// SyncLogsToResponse convierte lista de entidades a response
func SyncLogsToResponse(logs []*entities.InvoiceSyncLog) []response.SyncLog {
	items := make([]response.SyncLog, 0, len(logs))
	for _, log := range logs {
		items = append(items, SyncLogToResponse(log))
	}
	return items
}

// ToInvoiceableOrderResponse convierte OrderData a InvoiceableOrder response
func ToInvoiceableOrderResponse(order *dtos.OrderData) response.InvoiceableOrder {
	return response.InvoiceableOrder{
		ID:           order.ID,
		BusinessID:   order.BusinessID,
		OrderNumber:  order.OrderNumber,
		CustomerName: order.CustomerName,
		TotalAmount:  order.TotalAmount,
		Currency:     order.Currency,
		CreatedAt:    order.CreatedAt,
	}
}
