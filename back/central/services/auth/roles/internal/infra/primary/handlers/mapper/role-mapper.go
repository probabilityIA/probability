package mapper

import (
	"github.com/secamc93/probability/back/central/services/auth/roles/internal/domain"
	"github.com/secamc93/probability/back/central/services/auth/roles/internal/infra/primary/handlers/request"
	"github.com/secamc93/probability/back/central/services/auth/roles/internal/infra/primary/handlers/response"
)

// ToRoleFilters convierte GetRolesByLevelRequest a filtros del dominio
func ToRoleFilters(req request.GetRolesByLevelRequest) domain.RoleFilters {
	level := req.Level
	return domain.RoleFilters{
		Level: &level,
	}
}

// ToRoleResponse convierte RoleDTO a RoleResponse
func ToRoleResponse(dto domain.RoleDTO) response.RoleResponse {
	return response.RoleResponse{
		ID:               dto.ID,
		Name:             dto.Name,
		Code:             dto.Code,
		Description:      dto.Description,
		Level:            dto.Level,
		IsSystem:         dto.IsSystem,
		ScopeID:          dto.ScopeID,
		ScopeName:        dto.ScopeName,
		ScopeCode:        dto.ScopeCode,
		BusinessTypeID:   dto.BusinessTypeID,
		BusinessTypeName: dto.BusinessTypeName,
	}
}

// ToRoleListResponse convierte un slice de RoleDTO a RoleListResponse
func ToRoleListResponse(roles []domain.RoleDTO) response.RoleListResponse {
	roleResponses := make([]response.RoleResponse, len(roles))
	for i, role := range roles {
		roleResponses[i] = ToRoleResponse(role)
	}

	return response.RoleListResponse{
		Success: true,
		Data:    roleResponses,
		Count:   len(roleResponses),
	}
}

// ToAssignPermissionsToRoleResponse construye la respuesta para asignar permisos a un rol
func ToAssignPermissionsToRoleResponse(roleID uint, permissionIDs []uint) response.AssignPermissionsToRoleResponse {
	return response.AssignPermissionsToRoleResponse{
		Success:       true,
		Message:       "Permisos asignados exitosamente al rol",
		RoleID:        roleID,
		PermissionIDs: permissionIDs,
	}
}

// ToGetRolePermissionsResponse construye la respuesta con los permisos de un rol
func ToGetRolePermissionsResponse(roleID uint, permissions []PermissionDTO) response.GetRolePermissionsResponse {
	permissionResponses := make([]response.PermissionResponse, len(permissions))
	for i, perm := range permissions {
		permissionResponses[i] = response.PermissionResponse{
			ID:          perm.ID,
			Resource:    perm.Resource,
			Action:      perm.Action,
			Description: perm.Description,
			ScopeID:     perm.ScopeID,
			ScopeName:   perm.ScopeName,
			ScopeCode:   perm.ScopeCode,
		}
	}

	return response.GetRolePermissionsResponse{
		Success:     true,
		Message:     "Permisos del rol obtenidos exitosamente",
		RoleID:      roleID,
		Permissions: permissionResponses,
		Count:       len(permissionResponses),
	}
}

// PermissionDTO representa un permiso para mapeo
type PermissionDTO struct {
	ID          uint
	Resource    string
	Action      string
	Description string
	ScopeID     uint
	ScopeName   string
	ScopeCode   string
}

// PermissionToDTO convierte un permiso del dominio a PermissionDTO
func PermissionToDTO(perm domain.Permission) PermissionDTO {
	return PermissionDTO{
		ID:          perm.ID,
		Resource:    perm.Resource,
		Action:      perm.Action,
		Description: perm.Description,
		ScopeID:     perm.ResourceID,
		ScopeName:   "", // Se puede obtener de la relación Scope si está disponible
		ScopeCode:   "", // Se puede obtener de la relación Scope si está disponible
	}
}
