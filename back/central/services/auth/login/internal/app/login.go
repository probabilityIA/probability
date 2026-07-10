package app

import (
	"context"
	"fmt"
	"strings"

	"github.com/secamc93/probability/back/central/services/auth/login/internal/domain"

	"golang.org/x/crypto/bcrypt"
)

func (uc *AuthUseCase) Login(ctx context.Context, request domain.LoginRequest) (*domain.LoginResponse, error) {
	normalizedEmail := strings.ToLower(strings.TrimSpace(request.Email))
	uc.log.Info().Str("email", normalizedEmail).Msg("Iniciando proceso de login")

	if normalizedEmail == "" || request.Password == "" {
		uc.log.Error().Msg("Email o contraseña vacíos")
		return nil, domain.ErrEmailPasswordRequired
	}

	userAuth, err := uc.repository.GetUserByEmail(ctx, normalizedEmail)
	if err != nil {
		uc.log.Error().Err(err).Str("email", request.Email).Msg("Error al obtener usuario por email")
		return nil, fmt.Errorf("error al obtener usuario por email: %w", err)
	}

	if userAuth == nil {
		uc.log.Error().Str("email", normalizedEmail).Msg("Usuario no encontrado")
		return nil, domain.ErrUserNotFound
	}

	uc.log.Debug().
		Str("email", request.Email).
		Msg("Validando contraseña con bcrypt")

	if err := bcrypt.CompareHashAndPassword([]byte(userAuth.Password), []byte(request.Password)); err != nil {
		uc.log.Error().
			Err(err).
			Str("email", request.Email).
			Msg("Contraseña inválida")
		return nil, domain.ErrInvalidCredentials
	}

	if !userAuth.IsActive {
		pending, err := uc.repository.HasPendingEmailVerification(ctx, userAuth.ID)
		if err != nil {
			uc.log.Error().Err(err).Uint("user_id", userAuth.ID).Msg("Error consultando verificación pendiente")
			return nil, domain.ErrUserInactive
		}
		if pending {
			uc.log.Warn().Str("email", request.Email).Msg("Usuario inactivo pendiente de verificación")
			return nil, domain.ErrUserPendingVerification
		}
		uc.log.Error().Str("email", request.Email).Msg("Usuario inactivo")
		return nil, domain.ErrUserInactive
	}

	roles, err := uc.repository.GetUserRoles(ctx, userAuth.ID)
	if err != nil {
		uc.log.Error().Err(err).Uint("user_id", userAuth.ID).Msg("Error al obtener roles del usuario")
		return nil, fmt.Errorf("error interno del servidor")
	}

	businesses, err := uc.repository.GetUserBusinesses(ctx, userAuth.ID)
	if err != nil {
		uc.log.Error().Err(err).Uint("user_id", userAuth.ID).Msg("Error al obtener businesses del usuario")
	}

	avatarURL := userAuth.AvatarURL
	if avatarURL != "" && !strings.HasPrefix(avatarURL, "http") {
		base := strings.TrimRight(uc.env.Get("URL_BASE_DOMAIN_S3"), "/")
		if base != "" {
			avatarURL = fmt.Sprintf("%s/%s", base, strings.TrimLeft(avatarURL, "/"))
		}
	}

	uc.log.Info().
		Uint("user_id", userAuth.ID).
		Int("businesses_count", len(businesses)).
		Msg("Businesses obtenidos del usuario")

	if len(businesses) > 0 {
		for i, business := range businesses {
			businessLogoURL := business.LogoURL
			if businessLogoURL != "" && !strings.HasPrefix(businessLogoURL, "http") {
				base := strings.TrimRight(uc.env.Get("URL_BASE_DOMAIN_S3"), "/")
				if base != "" {
					businessLogoURL = fmt.Sprintf("%s/%s", base, strings.TrimLeft(businessLogoURL, "/"))
				}
			}
			businessNavbarURL := business.NavbarImageURL
			if businessNavbarURL != "" && !strings.HasPrefix(businessNavbarURL, "http") {
				base := strings.TrimRight(uc.env.Get("URL_BASE_DOMAIN_S3"), "/")
				if base != "" {
					businessNavbarURL = fmt.Sprintf("%s/%s", base, strings.TrimLeft(businessNavbarURL, "/"))
				}
			}

			uc.log.Info().
				Uint("user_id", userAuth.ID).
				Int("business_index", i).
				Uint("business_id", business.ID).
				Str("business_name", business.Name).
				Str("business_code", business.Code).
				Msg("Business encontrado")
		}
	} else {
		uc.log.Warn().
			Uint("user_id", userAuth.ID).
			Msg("Usuario sin businesses asociados")
	}

	uc.log.Info().
		Uint("user_id", userAuth.ID).
		Int("roles_count", len(roles)).
		Msg("Roles obtenidos del usuario")

	for i, role := range roles {
		uc.log.Info().
			Uint("user_id", userAuth.ID).
			Int("role_index", i).
			Uint("role_id", role.ID).
			Str("role_name", role.Name).
			Msg("Rol encontrado")
	}

	if isSuperAdmin(roles) {
		uc.log.Info().
			Uint("user_id", userAuth.ID).
			Msg("Usuario identificado como SUPER ADMIN")
	} else {
		uc.log.Info().
			Uint("user_id", userAuth.ID).
			Msg("Usuario NO es super admin")
	}

	for i, role := range roles {
		isSuper := role.ScopeCode == "platform"
		uc.log.Info().
			Uint("user_id", userAuth.ID).
			Int("role_index", i).
			Uint("role_id", role.ID).
			Str("role_name", role.Name).
			Str("role_scope_code", role.ScopeCode).
			Str("role_scope_name", role.ScopeName).
			Bool("is_super_admin", isSuper).
			Msg("Verificación de rol super admin por scope")
	}

	roleNames := make([]string, len(roles))
	for i, role := range roles {
		roleNames[i] = role.Name
	}

	var businessID, businessTypeID, roleID uint
	isSuperAdminUser := isSuperAdmin(roles)

	if isSuperAdminUser {
		businessID = 0
		businessTypeID = 0
		if len(roles) > 0 {
			roleID = roles[0].ID
		}
		uc.log.Info().
			Uint("user_id", userAuth.ID).
			Uint("business_id", businessID).
			Uint("role_id", roleID).
			Msg("Usuario super admin - usando business_id = 0")
	} else if len(businesses) > 0 {
		businessID = businesses[0].ID
		businessTypeID = businesses[0].BusinessTypeID

		userRole, err := uc.repository.GetUserRoleByBusiness(ctx, userAuth.ID, businessID)
		if err != nil || userRole == nil {
			uc.log.Warn().
				Err(err).
				Uint("user_id", userAuth.ID).
				Uint("business_id", businessID).
				Msg("No se pudo obtener rol del usuario en el business, usando primer rol disponible")
			if len(roles) > 0 {
				roleID = roles[0].ID
			}
		} else {
			roleID = userRole.ID
		}

		uc.log.Info().
			Uint("user_id", userAuth.ID).
			Uint("business_id", businessID).
			Uint("business_type_id", businessTypeID).
			Uint("role_id", roleID).
			Msg("Usando primer business para token JWT unificado")
	} else {
		businessID = 0
		businessTypeID = 0
		if len(roles) > 0 {
			roleID = roles[0].ID
		}
		uc.log.Warn().
			Uint("user_id", userAuth.ID).
			Msg("Usuario sin businesses - usando business_id = 0")
	}

	var subscriptionStatus string
	if isSuperAdminUser {
		subscriptionStatus = "active"
	} else if len(businesses) > 0 {
		subscriptionStatus = businesses[0].SubscriptionStatus
		if subscriptionStatus == "" {
			subscriptionStatus = "active"
		}
	} else {
		subscriptionStatus = "active"
	}

	token, err := uc.jwtService.GenerateToken(userAuth.ID, businessID, businessTypeID, roleID, subscriptionStatus)
	if err != nil {
		uc.log.Error().Err(err).Uint("user_id", userAuth.ID).Msg("Error al generar token JWT")
		return nil, fmt.Errorf("error interno del servidor")
	}

	uc.log.Info().
		Uint("user_id", userAuth.ID).
		Uint("token_business_id", businessID).
		Str("user_email", userAuth.Email).
		Strs("user_roles", roleNames).
		Msg("Token JWT generado exitosamente")

	fmt.Printf("[LOGIN DEBUG] User: %s (ID: %d) | Assigned BusinessID for Session: %d | IsSuperAdmin: %v\n",
		userAuth.Email, userAuth.ID, businessID, isSuperAdminUser)

	isFirstLogin := userAuth.LastLoginAt == nil

	if isFirstLogin {
		uc.log.Info().
			Str("email", userAuth.Email).
			Uint("user_id", userAuth.ID).
			Msg("Primer login detectado - se requiere cambio de contraseña")
	}

	if err := uc.repository.UpdateLastLogin(ctx, userAuth.ID); err != nil {
		uc.log.Warn().Err(err).Uint("user_id", userAuth.ID).Msg("Error al actualizar último login")
	} else {
		uc.log.Info().Uint("user_id", userAuth.ID).Msg("Último login actualizado")
	}

	var businessesList []domain.BusinessInfo

	if len(businesses) > 0 {
		businessesList = make([]domain.BusinessInfo, len(businesses))
		for i, business := range businesses {
			businessLogoURL := business.LogoURL
			if businessLogoURL != "" && !strings.HasPrefix(businessLogoURL, "http") {
				base := strings.TrimRight(uc.env.Get("URL_BASE_DOMAIN_S3"), "/")
				if base != "" {
					businessLogoURL = fmt.Sprintf("%s/%s", base, strings.TrimLeft(businessLogoURL, "/"))
				}
			}
			businessNavbarURL := business.NavbarImageURL
			if businessNavbarURL != "" && !strings.HasPrefix(businessNavbarURL, "http") {
				base := strings.TrimRight(uc.env.Get("URL_BASE_DOMAIN_S3"), "/")
				if base != "" {
					businessNavbarURL = fmt.Sprintf("%s/%s", base, strings.TrimLeft(businessNavbarURL, "/"))
				}
			}

			businessesList[i] = domain.BusinessInfo{
				ID:             business.ID,
				Name:           business.Name,
				Code:           business.Code,
				BusinessTypeID: business.BusinessTypeID,
				BusinessType: domain.BusinessTypeInfo{
					ID:          business.BusinessTypeID,
					Name:        business.BusinessTypeName,
					Code:        business.BusinessTypeCode,
					Description: "",
					Icon:        "",
				},
				Timezone:           business.Timezone,
				Address:            business.Address,
				Description:        business.Description,
				LogoURL:            businessLogoURL,
				PrimaryColor:       business.PrimaryColor,
				SecondaryColor:     business.SecondaryColor,
				TertiaryColor:      business.TertiaryColor,
				QuaternaryColor:    business.QuaternaryColor,
				NavbarImageURL:     businessNavbarURL,
				CustomDomain:       business.CustomDomain,
				IsActive:           business.IsActive,
				EnableDelivery:     business.EnableDelivery,
				EnablePickup:       business.EnablePickup,
				EnableReservations: business.EnableReservations,
			}
		}

		uc.log.Info().
			Uint("user_id", userAuth.ID).
			Int("businesses_count", len(businesses)).
			Msg("Businesses asignados al usuario")
	} else {
		uc.log.Info().
			Uint("user_id", userAuth.ID).
			Msg("Usuario sin businesses asignados")
	}

	userScope := "business"
	isSuperAdmin := isSuperAdmin(roles)
	if isSuperAdmin {
		userScope = "platform"
	}

	uc.log.Info().
		Uint("user_id", userAuth.ID).
		Str("user_scope", userScope).
		Bool("is_super_admin", isSuperAdmin).
		Msg("Información de scope del usuario")

	response := &domain.LoginResponse{
		User: domain.UserInfo{
			ID:          userAuth.ID,
			Name:        userAuth.Name,
			Email:       userAuth.Email,
			Phone:       userAuth.Phone,
			AvatarURL:   avatarURL,
			IsActive:    userAuth.IsActive,
			LastLoginAt: userAuth.LastLoginAt,
		},
		Token:                 token,
		RequirePasswordChange: isFirstLogin,
		Businesses:            businessesList,
		Scope:                 userScope,
		IsSuperAdmin:          isSuperAdmin,
	}

	uc.log.Info().
		Str("email", userAuth.Email).
		Uint("user_id", userAuth.ID).
		Bool("require_password_change", isFirstLogin).
		Msg("Login exitoso")

	return response, nil
}

func isSuperAdmin(roles []domain.Role) bool {
	for _, role := range roles {
		if role.ScopeCode == "platform" {
			return true
		}
	}
	return false
}
