package app

import (
	"context"
	"fmt"
	"strings"

	"github.com/secamc93/probability/back/central/services/auth/users/internal/domain"
)

// GetUsers obtiene usuarios filtrados y paginados
func (uc *UserUseCase) GetUsers(ctx context.Context, filters domain.UserFilters) (*domain.UserListDTO, error) {
	uc.log.Info().
		Int("page", filters.Page).
		Int("page_size", filters.PageSize).
		Str("name", filters.Name).
		Str("email", filters.Email).
		Str("phone", filters.Phone).
		Str("sort_by", filters.SortBy).
		Str("sort_order", filters.SortOrder).
		Msg("Iniciando caso de uso: obtener usuarios filtrados y paginados")

	// Validar y normalizar parámetros de paginación
	if filters.Page < 1 {
		filters.Page = 1
	}
	if filters.PageSize < 1 {
		filters.PageSize = 10
	}
	if filters.PageSize > 100 {
		filters.PageSize = 100
	}

	// Validar ordenamiento
	if filters.SortBy != "" {
		allowedSortFields := map[string]bool{
			"id": true, "name": true, "email": true, "phone": true, "is_active": true,
			"created_at": true, "updated_at": true,
		}
		if !allowedSortFields[filters.SortBy] {
			filters.SortBy = "created_at"
		}
	}

	if filters.SortOrder != "" {
		if filters.SortOrder != "asc" && filters.SortOrder != "desc" {
			filters.SortOrder = "desc"
		}
	}

	users, total, err := uc.repository.GetUsers(ctx, filters)
	if err != nil {
		uc.log.Error().Err(err).Msg("Error al obtener usuarios desde el repositorio")
		return nil, err
	}

	// Convertir UserQueryDTO a UserDTO y obtener relaciones
	userDTOs := make([]domain.UserDTO, len(users))
	for i, user := range users {
		// Procesar URL del avatar
		avatarURL := user.AvatarURL
		if avatarURL != "" {
			// Verificar si es un path relativo (no empieza con http)
			if !strings.HasPrefix(avatarURL, "http") {
				// Verificar si la imagen existe en S3
				exists, err := uc.s3.ImageExists(ctx, avatarURL)
				if err != nil {
					uc.log.Error().Err(err).Str("avatar_path", avatarURL).Msg("Error al verificar existencia de imagen en S3")
					// No fallar si hay error al verificar, continuar con URL por defecto
				} else if exists {
					// Generar URL completa usando el dominio de media
					mediaBaseURL := uc.getMediaBaseURL()
					if mediaBaseURL != "" {
						// Construir URL completa: MEDIA_BASE_URL + / + avatar_path
						avatarURL = fmt.Sprintf("%s/%s", mediaBaseURL, strings.TrimLeft(avatarURL, "/"))
						uc.log.Debug().Str("avatar_path", user.AvatarURL).Str("full_url", avatarURL).Msg("URL de avatar generada (media)")
					} else {
						uc.log.Warn().Str("avatar_path", avatarURL).Msg("URL_BASE_DOMAIN_S3 no configurada, usando path relativo")
					}
				} else {
					uc.log.Warn().Str("avatar_path", avatarURL).Msg("Imagen de avatar no encontrada en S3")
					// Si la imagen no existe, limpiar la URL
					avatarURL = ""
				}
			}
			// Si ya es una URL completa (empieza con http), mantenerla tal como está
		}

		// Convertir datos básicos
		userDTOs[i] = domain.UserDTO{
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
			userDTOs[i].Roles = make([]domain.RoleDTO, len(roles))
			for j, role := range roles {
				userDTOs[i].Roles[j] = domain.RoleDTO{
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

		// Obtener relaciones business_staff directamente (incluye roles)
		staffRelationships, err := uc.repository.GetBusinessStaffRelationships(ctx, user.ID)
		if err != nil {
			uc.log.Error().Uint("user_id", user.ID).Err(err).Msg("Error al obtener relaciones business_staff del usuario")
		} else {
			// Construir assignments desde business_staff (incluye roles)
			assignments := make([]domain.BusinessRoleAssignmentDetailed, 0)
			businessMap := make(map[uint]domain.BusinessDTO)

			// Primero obtener businesses para construir BusinessDTO
			businesses, err := uc.repository.GetUserBusinesses(ctx, user.ID)
			if err == nil {
				for _, business := range businesses {
					navbarURL := business.NavbarImageURL
					if navbarURL != "" && !strings.HasPrefix(navbarURL, "http") {
						base := strings.TrimRight(uc.getMediaBaseURL(), "/")
						if base != "" {
							navbarURL = fmt.Sprintf("%s/%s", base, strings.TrimLeft(navbarURL, "/"))
						}
					}

					businessMap[business.ID] = domain.BusinessDTO{
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
						Role:               nil, // Se completará desde staffRelationships
					}
				}
			}

			// Construir assignments desde staffRelationships (solo los que tienen business_id)
			for _, rel := range staffRelationships {
				if rel.BusinessID == 0 {
					continue // Saltar super usuarios (business_id = 0)
				}

				// Buscar el rol completo si existe
				var role *domain.RoleDTO
				if rel.RoleID > 0 {
					// Obtener información completa del rol
					roleInfo, err := uc.repository.GetRoleByID(ctx, rel.RoleID)
					if err == nil && roleInfo != nil {
						role = &domain.RoleDTO{
							ID:               roleInfo.ID,
							Name:             roleInfo.Name,
							Description:      roleInfo.Description,
							Level:            roleInfo.Level,
							IsSystem:         roleInfo.IsSystem,
							ScopeID:          roleInfo.ScopeID,
							ScopeName:        roleInfo.ScopeName,
							ScopeCode:        roleInfo.ScopeCode,
							BusinessTypeID:   roleInfo.BusinessTypeID,
							BusinessTypeName: roleInfo.BusinessTypeName,
						}
					} else {
						uc.log.Warn().
							Uint("user_id", user.ID).
							Uint("role_id", rel.RoleID).
							Err(err).
							Msg("No se pudo obtener información completa del rol")
					}
				}

				// Actualizar el Role en el BusinessDTO si existe
				if businessDTO, exists := businessMap[rel.BusinessID]; exists {
					businessDTO.Role = role
					businessMap[rel.BusinessID] = businessDTO
				}

				// Agregar al assignment (usar directamente rel que ya tiene RoleID y RoleName)
				assignments = append(assignments, rel)
			}

			// Convertir businessMap a slice
			businessDTOs := make([]domain.BusinessDTO, 0, len(businessMap))
			for _, b := range businessMap {
				businessDTOs = append(businessDTOs, b)
			}
			userDTOs[i].Businesses = businessDTOs
			userDTOs[i].BusinessRoleAssignments = assignments
		}

		// Determinar si es super usuario: tiene un rol con scope_id = 1 o scope code = "platform"
		isSuperUser := false
		var superUserRoleID uint
		for _, role := range userDTOs[i].Roles {
			if role.ScopeID == 1 || role.ScopeCode == "platform" {
				isSuperUser = true
				superUserRoleID = role.ID
				break
			}
		}

		userDTOs[i].IsSuperUser = isSuperUser

		// Si es super usuario, agregar assignment con business_id = 0
		if isSuperUser {
			// Buscar el nombre del rol
			roleName := ""
			for _, role := range userDTOs[i].Roles {
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
			userDTOs[i].BusinessRoleAssignments = append([]domain.BusinessRoleAssignmentDetailed{superUserAssignment}, userDTOs[i].BusinessRoleAssignments...)
		}
	}

	// Calcular total de páginas
	totalPages := int((total + int64(filters.PageSize) - 1) / int64(filters.PageSize))

	userListDTO := &domain.UserListDTO{
		Users:      userDTOs,
		Total:      total,
		Page:       filters.Page,
		PageSize:   filters.PageSize,
		TotalPages: totalPages,
	}

	uc.log.Info().
		Int("count", len(userDTOs)).
		Int64("total", total).
		Int("total_pages", totalPages).
		Msg("Usuarios obtenidos exitosamente")
	return userListDTO, nil
}
