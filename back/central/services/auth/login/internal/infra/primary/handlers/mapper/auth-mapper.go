package mapper

import (
	"github.com/secamc93/probability/back/central/services/auth/login/internal/domain"
	"github.com/secamc93/probability/back/central/services/auth/login/internal/infra/primary/handlers/response"
)

// ToLoginResponse convierte el dominio LoginResponse a response.LoginResponse
func ToLoginResponse(domainResponse *domain.LoginResponse) *response.LoginResponse {
	if domainResponse == nil {
		return nil
	}

	// Convertir businesses
	businesses := make([]response.BusinessInfo, len(domainResponse.Businesses))
	for i, business := range domainResponse.Businesses {
		businesses[i] = response.BusinessInfo{
			ID:             business.ID,
			Name:           business.Name,
			Code:           business.Code,
			BusinessTypeID: business.BusinessTypeID,
			BusinessType: response.BusinessTypeInfo{
				ID:          business.BusinessType.ID,
				Name:        business.BusinessType.Name,
				Code:        business.BusinessType.Code,
				Description: business.BusinessType.Description,
				Icon:        business.BusinessType.Icon,
			},
			Timezone:           business.Timezone,
			Address:            business.Address,
			Description:        business.Description,
			LogoURL:            business.LogoURL,
			PrimaryColor:       business.PrimaryColor,
			SecondaryColor:     business.SecondaryColor,
			TertiaryColor:      business.TertiaryColor,
			QuaternaryColor:    business.QuaternaryColor,
			NavbarImageURL:     business.NavbarImageURL,
			CustomDomain:       business.CustomDomain,
			IsActive:           business.IsActive,
			EnableDelivery:     business.EnableDelivery,
			EnablePickup:       business.EnablePickup,
			EnableReservations: business.EnableReservations,
		}
	}

	return &response.LoginResponse{
		User: response.UserInfo{
			ID:          domainResponse.User.ID,
			Name:        domainResponse.User.Name,
			Email:       domainResponse.User.Email,
			Phone:       domainResponse.User.Phone,
			AvatarURL:   domainResponse.User.AvatarURL,
			IsActive:    domainResponse.User.IsActive,
			LastLoginAt: domainResponse.User.LastLoginAt,
		},
		Token:                 domainResponse.Token,
		RequirePasswordChange: domainResponse.RequirePasswordChange,
		Businesses:            businesses,
		Scope:                 domainResponse.Scope,
		IsSuperAdmin:          domainResponse.IsSuperAdmin,
	}
}

// ToUserRolesPermissionsResponse convierte el dominio UserRolesPermissionsResponse a response.UserRolesPermissionsResponse
func ToUserRolesPermissionsResponse(domainResponse *domain.UserRolesPermissionsResponse) response.UserRolesPermissionsResponse {
	if domainResponse == nil {
		return response.UserRolesPermissionsResponse{}
	}

	resources := groupPermissionsByResource(domainResponse.Permissions)

	return response.UserRolesPermissionsResponse{
		IsSuper:          domainResponse.IsSuper,
		BusinessID:       domainResponse.BusinessID,
		BusinessName:     domainResponse.BusinessName,
		BusinessTypeID:   domainResponse.BusinessTypeID,
		BusinessTypeName: domainResponse.BusinessTypeName,
		Role:             toRoleInfo(domainResponse.Role),
		Resources:        resources,
	}
}

// toRoleInfo convierte un dto.RoleInfo a response.RoleInfo
func toRoleInfo(domainRole domain.RoleInfo) response.RoleInfo {
	return response.RoleInfo{
		ID:          domainRole.ID,
		Name:        domainRole.Name,
		Description: domainRole.Description,
	}
}

// toRoleInfoSlice convierte un slice de dtos.RoleInfo a response.RoleInfo
func toRoleInfoSlice(domainRoles []domain.RoleInfo) []response.RoleInfo {
	if domainRoles == nil {
		return nil
	}

	roles := make([]response.RoleInfo, len(domainRoles))
	for i, role := range domainRoles {
		roles[i] = response.RoleInfo{
			ID:          role.ID,
			Name:        role.Name,
			Description: role.Description,
		}
	}
	return roles
}

// groupPermissionsByResource agrupa los permisos por recurso
func groupPermissionsByResource(permissions []domain.PermissionInfo) []response.ResourcePermissions {
	if permissions == nil {
		return nil
	}

	// Mapa para agrupar permisos por recurso con estado activo
	resourceMap := make(map[string]struct {
		actions []string
		active  bool
	})

	// Agrupar permisos por recurso y determinar si está activo
	for _, permission := range permissions {
		if existing, exists := resourceMap[permission.Resource]; exists {
			// Si ya existe el recurso, agregar acción y mantener activo si al menos uno está activo
			existing.actions = append(existing.actions, permission.Action)
			existing.active = existing.active || permission.Active
			resourceMap[permission.Resource] = existing
		} else {
			// Nuevo recurso
			resourceMap[permission.Resource] = struct {
				actions []string
				active  bool
			}{
				actions: []string{permission.Action},
				active:  permission.Active,
			}
		}
	}

	// Convertir el mapa a slice de ResourcePermissions
	var resources []response.ResourcePermissions
	for resource, data := range resourceMap {
		resources = append(resources, response.ResourcePermissions{
			Resource: resource,
			Actions:  data.actions,
			Active:   data.active,
		})
	}

	return resources
}
