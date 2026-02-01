package mappers

import (
	"github.com/secamc93/probability/back/central/services/modules/payments/internal/domain/entities"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
)

// PaymentMethodToModel convierte entidad de dominio a modelo GORM
func PaymentMethodToModel(e *entities.PaymentMethod) *models.PaymentMethod {
	return &models.PaymentMethod{
		Model: gorm.Model{
			ID:        e.ID,
			CreatedAt: e.CreatedAt,
			UpdatedAt: e.UpdatedAt,
		},
		Code:        e.Code,
		Name:        e.Name,
		Description: e.Description,
		Category:    e.Category,
		Provider:    e.Provider,
		IsActive:    e.IsActive,
		Icon:        e.Icon,
		Color:       e.Color,
	}
}

// PaymentMappingToModel convierte entidad de dominio a modelo GORM
func PaymentMappingToModel(e *entities.PaymentMethodMapping) *models.PaymentMethodMapping {
	return &models.PaymentMethodMapping{
		Model: gorm.Model{
			ID:        e.ID,
			CreatedAt: e.CreatedAt,
			UpdatedAt: e.UpdatedAt,
		},
		IntegrationType: e.IntegrationType,
		OriginalMethod:  e.OriginalMethod,
		PaymentMethodID: e.PaymentMethodID,
		IsActive:        e.IsActive,
		Priority:        e.Priority,
	}
}
