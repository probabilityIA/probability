package app

import (
	"context"
	"errors"
	"fmt"

	"github.com/secamc93/probability/back/central/services/auth/permissions/internal/domain"
	"gorm.io/gorm"
)

// UpdatePermission actualiza un permiso existente
func (uc *PermissionUseCase) UpdatePermission(ctx context.Context, id uint, permissionDTO domain.UpdatePermissionDTO) (string, error) {
	uc.logger.Info().
		Uint("id", id).
		Str("name", permissionDTO.Name).
		Str("code", permissionDTO.Code).
		Msg("Actualizando permiso")

	// Validar datos de entrada
	if err := validateUpdatePermission(permissionDTO); err != nil {
		uc.logger.Error().Err(err).Msg("Error de validación al actualizar permiso")
		return "", err
	}

	// Verificar que el permiso existe
	existingPermission, err := uc.repository.GetPermissionByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			uc.logger.Warn().Uint("id", id).Msg("Permiso no encontrado para actualizar")
			return "", fmt.Errorf("permiso no encontrado")
		}
		uc.logger.Error().Uint("id", id).Err(err).Msg("Error al obtener permiso para actualizar")
		return "", err
	}

	// Actualizar los campos del permiso existente
	updatedPermission := updatePermissionFields(*existingPermission, permissionDTO)

	// Actualizar el permiso
	result, err := uc.repository.UpdatePermission(ctx, id, updatedPermission)
	if err != nil {
		uc.logger.Error().Uint("id", id).Err(err).Msg("Error al actualizar permiso en el repositorio")
		return "", err
	}

	uc.logger.Info().Uint("id", id).Str("result", result).Msg("Permiso actualizado exitosamente")
	return result, nil
}

// validateUpdatePermission valida los datos para actualizar un permiso
func validateUpdatePermission(permission domain.UpdatePermissionDTO) error {
	if permission.Name == "" {
		return fmt.Errorf("el nombre del permiso es requerido")
	}
	if permission.ResourceID == 0 {
		return fmt.Errorf("el resource ID del permiso es requerido")
	}
	if permission.ActionID == 0 {
		return fmt.Errorf("la action ID del permiso es requerida")
	}
	if permission.ScopeID == 0 {
		return fmt.Errorf("el scope ID del permiso es requerido")
	}
	return nil
}

// updatePermissionFields actualiza los campos de un permiso existente
func updatePermissionFields(existing domain.Permission, updateDTO domain.UpdatePermissionDTO) domain.Permission {
	businessTypeID := existing.BusinessTypeID
	if updateDTO.BusinessTypeID != nil {
		businessTypeID = *updateDTO.BusinessTypeID
	}

	return domain.Permission{
		ID:               existing.ID,
		Name:             updateDTO.Name,
		Description:      updateDTO.Description,
		ResourceID:       updateDTO.ResourceID,
		ActionID:         updateDTO.ActionID,
		ScopeID:          updateDTO.ScopeID,
		BusinessTypeID:   businessTypeID,
		BusinessTypeName: existing.BusinessTypeName,
	}
}
