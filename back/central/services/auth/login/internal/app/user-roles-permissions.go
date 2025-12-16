package app

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/auth/login/internal/domain"
)

// GetUserRolesPermissions maneja la lógica para obtener roles y permisos del usuario
func (uc *AuthUseCase) GetUserRolesPermissions(ctx context.Context, userID uint, businessID uint, token string) (*domain.UserRolesPermissionsResponse, error) {
	uc.log.Info().Uint("user_id", userID).Uint("business_id", businessID).Msg("Obteniendo roles y permisos del usuario")

	// Validar token
	if token == "" {
		uc.log.Error().Msg("Token requerido")
		return nil, fmt.Errorf("token inválido")
	}

	// Verificar que el token sea válido y obtener claims
	claims, err := uc.jwtService.ValidateToken(token)
	if err != nil {
		uc.log.Error().Err(err).Msg("Token inválido")
		return nil, fmt.Errorf("token inválido")
	}

	// Verificar que el usuario del token coincida con el userID solicitado
	if claims.UserID != userID {
		uc.log.Error().
			Uint("token_user_id", claims.UserID).
			Uint("requested_user_id", userID).
			Msg("El token no corresponde al usuario solicitado")
		return nil, fmt.Errorf("acceso denegado")
	}

	// Obtener usuario para verificar que existe
	user, err := uc.repository.GetUserByID(ctx, claims.UserID)
	if err != nil {
		uc.log.Error().Err(err).Uint("user_id", userID).Msg("Error al obtener usuario")
		return nil, fmt.Errorf("usuario no encontrado")
	}

	if user == nil {
		uc.log.Error().Uint("user_id", userID).Msg("Usuario no encontrado")
		return nil, fmt.Errorf("usuario no encontrado")
	}

	// Verificar que el usuario esté activo
	if !user.IsActive {
		uc.log.Error().Uint("user_id", userID).Msg("Usuario inactivo")
		return nil, fmt.Errorf("usuario inactivo")
	}

	// Obtener relación completa desde business_staff (única fuente de verdad)
	var bidPtr *uint
	if businessID != 0 {
		bidPtr = &businessID
	}
	bsRelation, err := uc.repository.GetBusinessStaffRelation(ctx, userID, bidPtr)
	if err != nil {
		uc.log.Error().Err(err).Uint("user_id", userID).Msg("Error al obtener relación desde business_staff")
		return nil, fmt.Errorf("error interno del servidor")
	}

	// Si no existe la relación en business_staff, retornar error
	if bsRelation == nil {
		uc.log.Warn().Uint("user_id", userID).Uint("business_id", businessID).Msg("Relación usuario-business no encontrada en business_staff")
		return nil, fmt.Errorf("relación usuario-business no encontrada")
	}

	uc.log.Info().
		Uint("user_id", userID).
		Any("business_id", bsRelation.BusinessID).
		Any("role_id", bsRelation.RoleID).
		Msg("Relación business_staff obtenida")

	// Determinar si es super: business_id NULL en business_staff
	isSuper := (bsRelation.BusinessID == nil)

	// Obtener información del business desde la relación
	var business *domain.BusinessInfo
	var businessType *domain.BusinessTypeInfo
	var activeResourcesMap map[uint]bool
	var roleIDPtr *uint = bsRelation.RoleID

	if roleIDPtr == nil {
		uc.log.Warn().Uint("user_id", userID).Any("business_id", bsRelation.BusinessID).Msg("Relación encontrada pero role_id es NULL")
	} else {
		uc.log.Info().Uint("user_id", userID).Uint("role_id", *roleIDPtr).Msg("Role ID encontrado en relación")
	}

	if bsRelation.BusinessID == nil {
		// Super admin sin business específico
		activeResourcesMap = make(map[uint]bool)
	} else {
		// Business existe en la relación
		if bsRelation.Business == nil {
			uc.log.Error().Uint("business_id", *bsRelation.BusinessID).Msg("Business no encontrado en relación")
			return nil, fmt.Errorf("business no encontrado")
		}

		business = &domain.BusinessInfo{
			ID:             bsRelation.Business.ID,
			Name:           bsRelation.Business.Name,
			Code:           bsRelation.Business.Code,
			BusinessTypeID: bsRelation.Business.BusinessTypeID,
		}

		businessType = &domain.BusinessTypeInfo{
			ID:          bsRelation.Business.BusinessTypeID,
			Name:        bsRelation.Business.BusinessTypeName,
			Code:        bsRelation.Business.BusinessTypeCode,
			Description: "",
			Icon:        "",
		}

		// Obtener recursos configurados para el business
		businessResourcesIDs, err := uc.repository.GetBusinessConfiguredResourcesIDs(ctx, *bsRelation.BusinessID)
		if err != nil {
			uc.log.Error().Err(err).Uint("business_id", *bsRelation.BusinessID).Msg("Error al obtener recursos configurados del business")
			return nil, fmt.Errorf("error interno del servidor")
		}

		// Crear mapa de recursos activos para búsqueda rápida
		activeResourcesMap = make(map[uint]bool)
		for _, resourceID := range businessResourcesIDs {
			activeResourcesMap[resourceID] = true
		}
	}

	var currentRole *domain.Role
	var allPermissions []domain.Permission
	if roleIDPtr != nil {
		// Cargar info del rol y permisos
		currentRole, _ = uc.repository.GetRoleByID(ctx, *roleIDPtr)
		perms, err := uc.repository.GetRolePermissions(ctx, *roleIDPtr)
		if err == nil {
			allPermissions = perms
		}
		// Ajustar isSuper según scope del rol si está disponible
		if currentRole != nil && (currentRole.ScopeID == 1 || currentRole.ScopeCode == "platform") {
			isSuper = true
		}
	} else {
		// Si no hay role_id asignado aún, verificar si es super por business_id NULL
		if bsRelation.BusinessID == nil {
			isSuper = true
		}
	}

	// Construir respuesta
	response := &domain.UserRolesPermissionsResponse{
		Success:     true,
		Message:     "Roles y permisos obtenidos exitosamente",
		UserID:      userID,
		Email:       user.Email,
		IsSuper:     isSuper,
		Role:        domain.RoleInfo{}, // Se llenará después
		Permissions: make([]domain.PermissionInfo, 0),
	}

	// Setear campos de business solo si existe
	if business != nil && businessType != nil {
		response.BusinessID = business.ID
		response.BusinessName = business.Name
		response.BusinessTypeID = businessType.ID
		response.BusinessTypeName = businessType.Name
	}

	// Mapear el primer rol (ya que ahora solo hay uno por business)
	if currentRole != nil {
		response.Role = domain.RoleInfo{
			ID:          currentRole.ID,
			Name:        currentRole.Name,
			Description: currentRole.Description,
			Level:       currentRole.Level,
			IsSystem:    currentRole.IsSystem,
			Scope:       currentRole.ScopeName,
		}
	}

	// Mapear permisos (eliminar duplicados)
	// Si hay recursos configurados para el business, solo incluir permisos de esos recursos
	// Si NO hay recursos configurados (mapa vacío), permitir todos los permisos del rol
	permissionMap := make(map[string]domain.PermissionInfo)
	hasConfiguredResources := len(activeResourcesMap) > 0

	uc.log.Info().
		Uint("user_id", userID).
		Int("active_resources_count", len(activeResourcesMap)).
		Int("all_permissions_count", len(allPermissions)).
		Bool("has_configured_resources", hasConfiguredResources).
		Msg("Filtrando permisos por recursos activos")

	for _, permission := range allPermissions {
		var isActive bool

		if bsRelation.BusinessID != nil {
			// Usuario con business asignado
			if hasConfiguredResources {
				// Filtrar solo si hay recursos configurados
				isActive = activeResourcesMap[permission.ResourceID]
				if !isActive {
					// Omitir permiso si el recurso no está activo para el business
					continue
				}
			} else {
				// Si no hay recursos configurados, permitir todos los permisos del rol
				isActive = true
			}
		} else {
			// Super admin sin business: mostrar todos los permisos como activos
			isActive = true
		}

		key := permission.Resource + ":" + permission.Action
		if _, exists := permissionMap[key]; !exists {
			permissionMap[key] = domain.PermissionInfo{
				ID:          permission.ID,
				Name:        permission.Name,
				Description: permission.Description,
				Resource:    permission.Resource,
				Action:      permission.Action,
				Active:      isActive,
			}
		}
	}

	// Convertir map a slice (solo incluye permisos de recursos activos)
	for _, permission := range permissionMap {
		response.Permissions = append(response.Permissions, permission)
	}

	uc.log.Info().
		Uint("user_id", userID).
		Int("permissions_count", len(response.Permissions)).
		Bool("is_super", isSuper).
		Msg("Roles y permisos obtenidos exitosamente")

	return response, nil
}
