package mappers

import (
	"github.com/secamc93/probability/back/central/services/auth/permissions/internal/domain"
	"github.com/secamc93/probability/back/migration/shared/models"
)

// ToPermissionEntitySlice convierte un slice de models.Permission a entities.Permission
func ToPermissionEntitySlice(models []models.Permission) []domain.Permission {
	if models == nil {
		return nil
	}

	entities := make([]domain.Permission, len(models))
	for i, model := range models {
		entities[i] = ToPermissionEntity(model)
	}
	return entities
}

// ToPermissionModel convierte entities.Permission a models.Permission
func ToPermissionModel(entity domain.Permission) models.Permission {
	return models.Permission{
		ResourceID: 0, // Se debe establecer según el nombre del recurso
		ActionID:   0, // Se debe establecer según el nombre de la acción

	}
}

// ToPermissionEntity convierte models.Permission a entities.Permission
func ToPermissionEntity(model models.Permission) domain.Permission {
	businessTypeID := uint(0)
	businessTypeName := ""
	if model.BusinessTypeID != nil {
		businessTypeID = *model.BusinessTypeID
	}
	if model.BusinessType != nil {
		businessTypeName = model.BusinessType.Name
	}

	scopeName := ""
	scopeCode := ""
	if model.Scope.ID != 0 {
		scopeName = model.Scope.Name
		scopeCode = model.Scope.Code
	}

	return domain.Permission{
		ID:               model.Model.ID,
		Name:             model.Name,
		Code:             "", // Code no existe en el modelo, debe generarse o venir de otra fuente
		Description:      model.Description,
		Resource:         model.Resource.Name,
		Action:           model.Action.Name,
		ResourceID:       model.ResourceID,
		ActionID:         model.ActionID,
		ScopeID:          model.ScopeID,
		ScopeName:        scopeName,
		ScopeCode:        scopeCode,
		BusinessTypeID:   businessTypeID,
		BusinessTypeName: businessTypeName,
	}
}
