package app

import (
	"context"
	"fmt"
)

// GenerateBusinessToken genera un token específico para un business
func (uc *AuthUseCase) GenerateBusinessToken(ctx context.Context, userID uint, businessID uint) (string, error) {
	uc.log.Info().
		Uint("user_id", userID).
		Uint("business_id", businessID).
		Msg("Generando token de business")

	// Verificar que el usuario existe y está activo
	user, err := uc.repository.GetUserByID(ctx, userID)
	if err != nil || user == nil {
		uc.log.Error().Err(err).Uint("user_id", userID).Msg("Usuario no encontrado")
		return "", fmt.Errorf("usuario no encontrado")
	}

	if !user.IsActive {
		uc.log.Error().Uint("user_id", userID).Msg("Usuario inactivo")
		return "", fmt.Errorf("usuario inactivo")
	}

	// Obtener roles del usuario para verificar si es super admin
	roles, err := uc.repository.GetUserRoles(ctx, userID)
	if err != nil {
		uc.log.Error().Err(err).Uint("user_id", userID).Msg("Error al obtener roles del usuario")
		return "", fmt.Errorf("error al obtener roles del usuario")
	}

	// Verificar si es super admin (scope platform o scope_id = 1)
	isSuperAdmin := false
	for _, role := range roles {
		if role.ScopeCode == "platform" || role.ScopeID == 1 {
			isSuperAdmin = true
			uc.log.Info().
				Uint("user_id", userID).
				Uint("role_id", role.ID).
				Str("role_name", role.Name).
				Str("scope_code", role.ScopeCode).
				Uint("scope_id", role.ScopeID).
				Msg("Usuario identificado como SUPER ADMIN")
			break
		}
	}

	// Si es super admin, generar token con business_id = 0
	if isSuperAdmin {
		// Usar el primer rol para el token
		var roleID uint
		if len(roles) > 0 {
			roleID = roles[0].ID
		}

		uc.log.Info().
			Uint("user_id", userID).
			Uint("role_id", roleID).
			Msg("Generando token de super admin con business_id = 0")

		businessToken, err := uc.jwtService.GenerateBusinessToken(
			userID,
			0, // business_id = 0 para super admin
			0, // business_type_id = 0 para super admin
			roleID,
		)
		if err != nil {
			uc.log.Error().Err(err).
				Uint("user_id", userID).
				Msg("Error al generar token de super admin")
			return "", fmt.Errorf("error al generar token de super admin")
		}

		uc.log.Info().
			Uint("user_id", userID).
			Msg("Token de super admin generado exitosamente")

		return businessToken, nil
	}

	// Para usuarios normales, continuar con la validación de business
	// Obtener información del business
	business, err := uc.repository.GetBusinessByID(ctx, businessID)
	if err != nil || business == nil {
		uc.log.Error().
			Err(err).
			Uint("business_id", businessID).
			Msg("Business no encontrado")
		return "", fmt.Errorf("business no encontrado")
	}

	if !business.IsActive {
		uc.log.Error().Uint("business_id", businessID).Msg("Business inactivo")
		return "", fmt.Errorf("business inactivo")
	}

	// Verificar que el usuario esté relacionado al business
	userBusinesses, err := uc.repository.GetUserBusinesses(ctx, userID)
	if err != nil {
		uc.log.Error().Err(err).Uint("user_id", userID).Msg("Error al obtener businesses del usuario")
		return "", fmt.Errorf("error al validar relación usuario-business")
	}

	// Verificar que el usuario tenga acceso al business
	hasAccess := false
	for _, ub := range userBusinesses {
		if ub.ID == businessID {
			hasAccess = true
			break
		}
	}

	if !hasAccess {
		uc.log.Error().
			Uint("user_id", userID).
			Uint("business_id", businessID).
			Msg("El usuario no tiene acceso a este business")
		return "", fmt.Errorf("el usuario no tiene acceso a este business")
	}

	// Obtener el rol del usuario en este business
	userRole, err := uc.repository.GetUserRoleByBusiness(ctx, userID, businessID)
	if err != nil || userRole == nil {
		uc.log.Error().
			Err(err).
			Uint("user_id", userID).
			Uint("business_id", businessID).
			Msg("Error al obtener rol del usuario en el business")
		return "", fmt.Errorf("error al obtener rol del usuario")
	}

	uc.log.Info().
		Uint("user_id", userID).
		Uint("business_id", businessID).
		Uint("business_type_id", business.BusinessTypeID).
		Uint("role_id", userRole.ID).
		Str("role_name", userRole.Name).
		Msg("Datos obtenidos para generar business token")

	// Generar el business token
	businessToken, err := uc.jwtService.GenerateBusinessToken(
		userID,
		businessID,
		business.BusinessTypeID,
		userRole.ID,
	)
	if err != nil {
		uc.log.Error().Err(err).
			Uint("user_id", userID).
			Uint("business_id", businessID).
			Msg("Error al generar business token")
		return "", fmt.Errorf("error al generar token de business")
	}

	uc.log.Info().
		Uint("user_id", userID).
		Uint("business_id", businessID).
		Msg("Business token generado exitosamente")

	return businessToken, nil
}
