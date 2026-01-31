package mappers

import (
	"encoding/json"

	"gorm.io/datatypes"
	"gorm.io/gorm"

	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/entities"
	"github.com/secamc93/probability/back/migration/shared/models"
)

// InvoiceToModel convierte una entidad de dominio a modelo GORM
func InvoiceToModel(entity *entities.Invoice) *models.Invoice {
	if entity == nil {
		return nil
	}

	model := &models.Invoice{
		Model: gorm.Model{
			ID:        entity.ID,
			CreatedAt: entity.CreatedAt,
			UpdatedAt: entity.UpdatedAt,
		},
		OrderID:             entity.OrderID,
		BusinessID:          entity.BusinessID,
		InvoicingProviderID: entity.InvoicingProviderID,
		InvoiceNumber:       entity.InvoiceNumber,
		ExternalID:          entity.ExternalID,
		InternalNumber:      entity.InternalNumber,
		Subtotal:            entity.Subtotal,
		Tax:                 entity.Tax,
		Discount:            entity.Discount,
		ShippingCost:        entity.ShippingCost,
		TotalAmount:         entity.TotalAmount,
		Currency:            entity.Currency,
		CustomerName:        entity.CustomerName,
		CustomerEmail:       entity.CustomerEmail,
		CustomerPhone:       entity.CustomerPhone,
		CustomerDNI:         entity.CustomerDNI,
		Status:              entity.Status,
		IssuedAt:            entity.IssuedAt,
		CancelledAt:         entity.CancelledAt,
		ExpiresAt:           entity.ExpiresAt,
		InvoiceURL:          entity.InvoiceURL,
		PDFURL:              entity.PDFURL,
		XMLURL:              entity.XMLURL,
		CUFE:                entity.CUFE,
		Notes:               entity.Notes,
	}

	if entity.DeletedAt != nil {
		model.DeletedAt = gorm.DeletedAt{Time: *entity.DeletedAt, Valid: true}
	}

	// Convertir map a JSONB
	if entity.Metadata != nil {
		if data, err := json.Marshal(entity.Metadata); err == nil {
			model.Metadata = datatypes.JSON(data)
		}
	}

	if entity.ProviderResponse != nil {
		if data, err := json.Marshal(entity.ProviderResponse); err == nil {
			model.ProviderResponse = datatypes.JSON(data)
		}
	}

	return model
}
