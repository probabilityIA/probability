package mappers

import (
	"github.com/secamc93/probability/back/central/services/auth/bussines/internal/domain"
	"github.com/secamc93/probability/back/migration/shared/models"
)

// ToBusinessEntity convierte models.Business a entities.Business
func ToBusinessEntity(model models.Business) domain.Business {
	business := domain.Business{
		ID:                 model.Model.ID,
		Name:               model.Name,
		Code:               model.Code,
		BusinessTypeID:     model.BusinessTypeID,
		Timezone:           model.Timezone,
		Address:            model.Address,
		Description:        model.Description,
		LogoURL:            model.LogoURL,
		PrimaryColor:       model.PrimaryColor,
		SecondaryColor:     model.SecondaryColor,
		TertiaryColor:      model.TertiaryColor,
		QuaternaryColor:    model.QuaternaryColor,
		NavbarImageURL:     model.NavbarImageURL,
		CustomDomain:       model.CustomDomain,
		IsActive:           model.IsActive,
		EnableDelivery:     model.EnableDelivery,
		EnablePickup:       model.EnablePickup,
		EnableReservations: model.EnableReservations,
		CreatedAt:          model.Model.CreatedAt,
		UpdatedAt:          model.Model.UpdatedAt,
		DeletedAt:          &model.Model.DeletedAt.Time,
	}

	// Mapear BusinessType si existe
	if model.BusinessType.ID != 0 {
		business.BusinessType = &domain.BusinessType{
			ID:          model.BusinessType.ID,
			Name:        model.BusinessType.Name,
			Code:        model.BusinessType.Code,
			Description: model.BusinessType.Description,
			Icon:        model.BusinessType.Icon,
			IsActive:    model.BusinessType.IsActive,
			CreatedAt:   model.BusinessType.CreatedAt,
			UpdatedAt:   model.BusinessType.UpdatedAt,
			DeletedAt:   &model.BusinessType.DeletedAt.Time,
		}
	}

	return business
}

// ToBusinessEntitySlice convierte un slice de models.Business a entities.Business
func ToBusinessEntitySlice(models []models.Business) []domain.Business {
	if models == nil {
		return nil
	}

	entities := make([]domain.Business, len(models))
	for i, model := range models {
		entities[i] = ToBusinessEntity(model)
	}
	return entities
}

// ToBusinessModel convierte entities.Business a models.Business
func ToBusinessModel(entity domain.Business) models.Business {
	return models.Business{
		Name:               entity.Name,
		Code:               entity.Code,
		BusinessTypeID:     entity.BusinessTypeID,
		Timezone:           entity.Timezone,
		Address:            entity.Address,
		Description:        entity.Description,
		LogoURL:            entity.LogoURL,
		PrimaryColor:       entity.PrimaryColor,
		SecondaryColor:     entity.SecondaryColor,
		TertiaryColor:      entity.TertiaryColor,
		QuaternaryColor:    entity.QuaternaryColor,
		NavbarImageURL:     entity.NavbarImageURL,
		CustomDomain:       entity.CustomDomain,
		IsActive:           entity.IsActive,
		EnableDelivery:     entity.EnableDelivery,
		EnablePickup:       entity.EnablePickup,
		EnableReservations: entity.EnableReservations,
	}
}
