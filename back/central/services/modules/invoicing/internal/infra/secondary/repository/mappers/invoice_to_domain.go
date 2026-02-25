package mappers

import (
	"encoding/json"

	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/entities"
	"github.com/secamc93/probability/back/migration/shared/models"
)

// InvoiceToDomain convierte un modelo GORM a entidad de dominio
func InvoiceToDomain(model *models.Invoice) *entities.Invoice {
	if model == nil {
		return nil
	}

	entity := &entities.Invoice{
		ID:                     model.ID,
		CreatedAt:              model.CreatedAt,
		UpdatedAt:              model.UpdatedAt,
		OrderID:                model.OrderID,
		BusinessID:             model.BusinessID,
		InvoicingProviderID:    model.InvoicingProviderID,    // Direct assignment (both are *uint)
		InvoicingIntegrationID: model.InvoicingIntegrationID, // New field for integrationCore
		InvoiceNumber:       model.InvoiceNumber,
		ExternalID:          model.ExternalID,
		InternalNumber:      model.InternalNumber,
		Subtotal:            model.Subtotal,
		Tax:                 model.Tax,
		Discount:            model.Discount,
		ShippingCost:        model.ShippingCost,
		TotalAmount:         model.TotalAmount,
		Currency:            model.Currency,
		CustomerName:        model.CustomerName,
		CustomerEmail:       model.CustomerEmail,
		CustomerPhone:       model.CustomerPhone,
		CustomerDNI:         model.CustomerDNI,
		Status:              model.Status,
		IssuedAt:            model.IssuedAt,
		CancelledAt:         model.CancelledAt,
		ExpiresAt:           model.ExpiresAt,
		InvoiceURL:          model.InvoiceURL,
		PDFURL:              model.PDFURL,
		XMLURL:              model.XMLURL,
		CUFE:                model.CUFE,
		Notes:               model.Notes,
	}

	if model.DeletedAt.Valid {
		entity.DeletedAt = &model.DeletedAt.Time
	}

	// Convertir JSONB a map
	if model.Metadata != nil {
		var metadata map[string]interface{}
		if err := json.Unmarshal(model.Metadata, &metadata); err == nil {
			entity.Metadata = metadata
		}
	}

	if model.ProviderResponse != nil {
		var response map[string]interface{}
		if err := json.Unmarshal(model.ProviderResponse, &response); err == nil {
			entity.ProviderResponse = response
		}
	}

	// Extraer logo y nombre del proveedor desde la relación preloaded
	if model.InvoicingIntegration.IntegrationType != nil {
		if model.InvoicingIntegration.IntegrationType.ImageURL != "" {
			logoURL := model.InvoicingIntegration.IntegrationType.ImageURL
			entity.ProviderLogoURL = &logoURL
		}
		if model.InvoicingIntegration.IntegrationType.Name != "" {
			providerName := model.InvoicingIntegration.IntegrationType.Name
			entity.ProviderName = &providerName
		}
	}

	// Mapear items de factura si fueron preloaded
	if len(model.InvoiceItems) > 0 {
		items := make([]entities.InvoiceItem, 0, len(model.InvoiceItems))
		for _, it := range model.InvoiceItems {
			item := entities.InvoiceItem{
				ID:             it.ID,
				CreatedAt:      it.CreatedAt,
				UpdatedAt:      it.UpdatedAt,
				InvoiceID:      it.InvoiceID,
				ProductID:      it.ProductID,
				SKU:            it.SKU,
				Name:           it.Name,
				Description:    it.Description,
				Quantity:       it.Quantity,
				UnitPrice:      it.UnitPrice,
				TotalPrice:     it.TotalPrice,
				Currency:       it.Currency,
				Tax:            it.Tax,
				TaxRate:        it.TaxRate,
				Discount:       it.Discount,
				ProviderItemID: it.ProviderItemID,
			}
			if it.DeletedAt.Valid {
				item.DeletedAt = &it.DeletedAt.Time
			}
			items = append(items, item)
		}
		entity.Items = items
	}

	return entity
}

// InvoiceListToDomain convierte una lista de modelos a entidades
func InvoiceListToDomain(models []*models.Invoice) []*entities.Invoice {
	entities := make([]*entities.Invoice, 0, len(models))
	for _, model := range models {
		entities = append(entities, InvoiceToDomain(model))
	}
	return entities
}
