package app

import (
	"context"
	"fmt"
	"strings"

	"github.com/secamc93/probability/back/central/services/auth/users/internal/domain"
)

// GetUserByID obtiene un usuario por su ID
func (uc *UserUseCase) GetUserByID(ctx context.Context, id uint) (*domain.UserDTO, error) {
	uc.log.Info().Uint("id", id).Msg("Iniciando caso de uso: obtener usuario por ID")

	user, err := uc.repository.GetUserByID(ctx, id)
	if err != nil {
		uc.log.Error().Uint("id", id).Err(err).Msg("Error al obtener usuario por ID desde el repositorio")
		return nil, fmt.Errorf("usuario no encontrado")
	}

	if user == nil {
		uc.log.Error().Uint("id", id).Msg("Usuario no encontrado")
		return nil, fmt.Errorf("usuario no encontrado")
	}

	// Completar URL del avatar si es path relativo
	avatarURL := user.AvatarURL
	if avatarURL != "" && !strings.HasPrefix(avatarURL, "http") {
		base := strings.TrimRight(uc.env.Get("URL_BASE_DOMAIN_S3"), "/")
		if base != "" {
			avatarURL = fmt.Sprintf("%s/%s", base, strings.TrimLeft(avatarURL, "/"))
		}
	}

	// Convertir entidad a DTO
	userDTO := domain.UserDTO{
		ID:          user.ID,
		Name:        user.Name,
		Email:       user.Email,
		Phone:       user.Phone,
		AvatarURL:   avatarURL, // URL completa o vacía
		IsActive:    user.IsActive,
		LastLoginAt: user.LastLoginAt,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		DeletedAt:   user.DeletedAt,
	}

	// Obtener roles del usuario
	roles, err := uc.repository.GetUserRoles(ctx, user.ID)
	if err != nil {
		uc.log.Error().Uint("user_id", user.ID).Err(err).Msg("Error al obtener roles del usuario")
	} else {
		// Convertir roles a DTOs
		userDTO.Roles = make([]domain.RoleDTO, len(roles))
		for i, role := range roles {
			userDTO.Roles[i] = domain.RoleDTO{
				ID:          role.ID,
				Name:        role.Name,
				Description: role.Description,
				Level:       role.Level,
				IsSystem:    role.IsSystem,
				ScopeID:     role.ScopeID,
				ScopeName:   role.ScopeName,
				ScopeCode:   role.ScopeCode,
			}
		}
	}

	// Obtener businesses del usuario
	businesses, err := uc.repository.GetUserBusinesses(ctx, user.ID)
	if err != nil {
		uc.log.Error().Uint("user_id", user.ID).Err(err).Msg("Error al obtener businesses del usuario")
	} else {
		// Convertir businesses a DTOs
		userDTO.Businesses = make([]domain.BusinessDTO, len(businesses))
		for i, business := range businesses {
			navbarURL := business.NavbarImageURL
			if navbarURL != "" && !strings.HasPrefix(navbarURL, "http") {
				base := strings.TrimRight(uc.env.Get("URL_BASE_DOMAIN_S3"), "/")
				if base != "" {
					navbarURL = fmt.Sprintf("%s/%s", base, strings.TrimLeft(navbarURL, "/"))
				}
			}

			// Obtener el rol del usuario en este business desde business_staff
			var role *domain.RoleDTO
			uc.log.Info().
				Uint("user_id", user.ID).
				Uint("business_id", business.ID).
				Msg("Obteniendo rol por business")

			businessRole, err := uc.repository.GetUserRoleByBusiness(ctx, user.ID, business.ID)
			if err == nil && businessRole != nil {
				uc.log.Info().
					Uint("user_id", user.ID).
					Uint("business_id", business.ID).
					Uint("role_id", businessRole.ID).
					Str("role_name", businessRole.Name).
					Msg("Rol encontrado para business")
				role = &domain.RoleDTO{
					ID:               businessRole.ID,
					Name:             businessRole.Name,
					Description:      businessRole.Description,
					Level:            businessRole.Level,
					IsSystem:         businessRole.IsSystem,
					ScopeID:          businessRole.ScopeID,
					ScopeName:        businessRole.ScopeName,
					ScopeCode:        businessRole.ScopeCode,
					BusinessTypeID:   businessRole.BusinessTypeID,
					BusinessTypeName: businessRole.BusinessTypeName,
				}
			} else {
				uc.log.Warn().
					Uint("user_id", user.ID).
					Uint("business_id", business.ID).
					Err(err).
					Msg("No se encontró rol para business")
			}

			userDTO.Businesses[i] = domain.BusinessDTO{
				ID:                 business.ID,
				Name:               business.Name,
				Code:               business.Code,
				BusinessTypeID:     business.BusinessTypeID,
				Timezone:           business.Timezone,
				Address:            business.Address,
				Description:        business.Description,
				LogoURL:            business.LogoURL,
				PrimaryColor:       business.PrimaryColor,
				SecondaryColor:     business.SecondaryColor,
				TertiaryColor:      business.TertiaryColor,
				QuaternaryColor:    business.QuaternaryColor,
				NavbarImageURL:     navbarURL,
				CustomDomain:       business.CustomDomain,
				IsActive:           business.IsActive,
				EnableDelivery:     business.EnableDelivery,
				EnablePickup:       business.EnablePickup,
				EnableReservations: business.EnableReservations,
				BusinessTypeName:   business.BusinessTypeName,
				BusinessTypeCode:   business.BusinessTypeCode,
				Role:               role,
			}
		}
	}

	// Construir assignments SIEMPRE desde user_businesses; si no hay rol asignado, dejarlo en blanco
	if len(userDTO.Businesses) > 0 {
		assignments := make([]domain.BusinessRoleAssignmentDetailed, 0, len(userDTO.Businesses))
		for _, b := range userDTO.Businesses {
			assignment := domain.BusinessRoleAssignmentDetailed{
				BusinessID:   b.ID,
				BusinessName: b.Name,
			}
			if b.Role != nil {
				assignment.RoleID = b.Role.ID
				assignment.RoleName = b.Role.Name
			}
			assignments = append(assignments, assignment)
		}
		userDTO.BusinessRoleAssignments = assignments
	}

	// Determinar si es super usuario: tiene un rol con scope_id = 1 o scope code = "platform"
	isSuperUser := false
	var superUserRoleID uint
	for _, role := range userDTO.Roles {
		if role.ScopeID == 1 || role.ScopeCode == "platform" {
			isSuperUser = true
			superUserRoleID = role.ID
			break
		}
	}

	userDTO.IsSuperUser = isSuperUser

	// Si es super usuario, agregar assignment con business_id = 0
	if isSuperUser {
		// Buscar el nombre del rol
		roleName := ""
		for _, role := range userDTO.Roles {
			if role.ID == superUserRoleID {
				roleName = role.Name
				break
			}
		}

		superUserAssignment := domain.BusinessRoleAssignmentDetailed{
			BusinessID:   0,
			BusinessName: "", // Super usuarios no tienen business específico
			RoleID:       superUserRoleID,
			RoleName:     roleName,
		}
		// Agregar al inicio del array
		userDTO.BusinessRoleAssignments = append([]domain.BusinessRoleAssignmentDetailed{superUserAssignment}, userDTO.BusinessRoleAssignments...)
	}

	uc.log.Info().Uint("id", id).Bool("is_super_user", isSuperUser).Msg("Usuario obtenido exitosamente")
	return &userDTO, nil
}
