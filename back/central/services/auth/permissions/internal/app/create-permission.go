package app

import (
	"context"
	"fmt"
	"strings"

	"github.com/secamc93/probability/back/central/services/auth/permissions/internal/domain"
)

// CreatePermission crea un nuevo permiso
func (uc *PermissionUseCase) CreatePermission(ctx context.Context, permissionDTO domain.CreatePermissionDTO) (string, error) {
	uc.logger.Info().
		Str("name", permissionDTO.Name).
		Str("code", permissionDTO.Code).
		Uint("resource_id", permissionDTO.ResourceID).
		Uint("action_id", permissionDTO.ActionID).
		Msg("Creando nuevo permiso")

	// Validar datos de entrada
	if err := validateCreatePermission(permissionDTO); err != nil {
		uc.logger.Error().Err(err).Msg("Error de validación al crear permiso")
		return "", err
	}

	// Validar que no existe un permiso con el mismo nombre
	exists, err := uc.repository.PermissionExistsByName(ctx, permissionDTO.Name)
	if err != nil {
		uc.logger.Error().Err(err).Str("name", permissionDTO.Name).Msg("Error verificando existencia de permiso por nombre")
		return "", fmt.Errorf("error verificando existencia de permiso: %w", err)
	}
	if exists {
		uc.logger.Warn().Str("name", permissionDTO.Name).Msg("Ya existe un permiso con este nombre")
		return "", fmt.Errorf("ya existe un permiso con el nombre '%s'", permissionDTO.Name)
	}

	// Generar código automáticamente si no se proporciona
	if permissionDTO.Code == "" {
		generatedCode, err := uc.generatePermissionCode(ctx, permissionDTO)
		if err != nil {
			uc.logger.Error().Err(err).Msg("Error al generar código de permiso")
			return "", fmt.Errorf("error al generar código de permiso: %w", err)
		}
		permissionDTO.Code = generatedCode
		uc.logger.Info().Str("generated_code", generatedCode).Msg("Código generado automáticamente")
	}

	// Convertir DTO a entidad
	permission := dtosToPermissionEntity(permissionDTO)

	// Crear el permiso
	result, err := uc.repository.CreatePermission(ctx, permission)
	if err != nil {
		uc.logger.Error().Err(err).Msg("Error al crear permiso en el repositorio")
		return "", err
	}

	uc.logger.Info().Str("result", result).Msg("Permiso creado exitosamente")
	return result, nil
}

// validateCreatePermission valida los datos para crear un permiso
func validateCreatePermission(permission domain.CreatePermissionDTO) error {
	if permission.Name == "" {
		return fmt.Errorf("el nombre del permiso es requerido")
	}
	// El código ya no es obligatorio, se genera automáticamente
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

// generatePermissionCode genera un código automáticamente basado en el nombre del tipo de business y el nombre del permiso
func (uc *PermissionUseCase) generatePermissionCode(ctx context.Context, permissionDTO domain.CreatePermissionDTO) (string, error) {
	// Si tiene business_type_id, obtener el nombre del tipo de business
	businessTypePrefix := "generic" // Por defecto para permisos genéricos

	if permissionDTO.BusinessTypeID != nil && *permissionDTO.BusinessTypeID > 0 {
		// Simplificación: usar "bt_" + ID si no podemos obtener el nombre sin el repositorio
		businessTypePrefix = fmt.Sprintf("bt_%d", *permissionDTO.BusinessTypeID)
	}

	// Normalizar el nombre del permiso: lowercase y reemplazar espacios con guiones bajos
	permissionName := strings.ToLower(permissionDTO.Name)
	permissionName = strings.ReplaceAll(permissionName, " ", "_")

	// Combinar: business_type_name_permission_name
	generatedCode := fmt.Sprintf("%s_%s", businessTypePrefix, permissionName)

	return generatedCode, nil
}

// dtosToPermissionEntity convierte un CreatePermissionDTO a entidad Permission
func dtosToPermissionEntity(permissionDTO domain.CreatePermissionDTO) domain.Permission {
	return domain.Permission{
		Name:           permissionDTO.Name,
		Description:    permissionDTO.Description,
		ResourceID:     permissionDTO.ResourceID,
		ActionID:       permissionDTO.ActionID,
		ScopeID:        permissionDTO.ScopeID,
		BusinessTypeID: *permissionDTO.BusinessTypeID,
		// BusinessTypeName se establecerá desde el modelo si es necesario
	}
}
