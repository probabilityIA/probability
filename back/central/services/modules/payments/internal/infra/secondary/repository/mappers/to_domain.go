package mappers

import (
	"github.com/secamc93/probability/back/central/services/modules/payments/internal/domain/entities"
	"github.com/secamc93/probability/back/migration/shared/models"
)

// PaymentMethodToDomain convierte modelo GORM a entidad de dominio
func PaymentMethodToDomain(m *models.PaymentMethod) entities.PaymentMethod {
	return entities.PaymentMethod{
		ID:          m.ID,
		Code:        m.Code,
		Name:        m.Name,
		Description: m.Description,
		Category:    m.Category,
		Provider:    m.Provider,
		IsActive:    m.IsActive,
		Icon:        m.Icon,
		Color:       m.Color,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

// PaymentMethodsToDomain convierte slice de modelos GORM a slice de entidades de dominio
func PaymentMethodsToDomain(models []models.PaymentMethod) []entities.PaymentMethod {
	result := make([]entities.PaymentMethod, len(models))
	for i, m := range models {
		result[i] = PaymentMethodToDomain(&m)
	}
	return result
}

// PaymentMappingToDomain convierte modelo GORM a entidad de dominio
func PaymentMappingToDomain(m *models.PaymentMethodMapping) entities.PaymentMethodMapping {
	return entities.PaymentMethodMapping{
		ID:              m.ID,
		IntegrationType: m.IntegrationType,
		OriginalMethod:  m.OriginalMethod,
		PaymentMethodID: m.PaymentMethodID,
		PaymentMethod:   PaymentMethodToDomain(&m.PaymentMethod),
		IsActive:        m.IsActive,
		Priority:        m.Priority,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
	}
}

// PaymentMappingsToDomain convierte slice de modelos GORM a slice de entidades de dominio
func PaymentMappingsToDomain(models []models.PaymentMethodMapping) []entities.PaymentMethodMapping {
	result := make([]entities.PaymentMethodMapping, len(models))
	for i, m := range models {
		result[i] = PaymentMappingToDomain(&m)
	}
	return result
}
