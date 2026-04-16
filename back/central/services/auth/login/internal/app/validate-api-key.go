package app

// // ValidateAPIKey valida una API Key contra la base de datos y retorna la información del usuario
// func (uc *AuthUseCase) ValidateAPIKey(ctx context.Context, request domain.ValidateAPIKeyRequest) (*domain.ValidateAPIKeyResponse, error) {
// 	uc.log.Debug().Str("api_key", request.APIKey[:10]+"...").Msg("Validando API Key")

// 	// Validar formato de API Key
// 	apiKeyService := apikey.NewService()
// 	if !apiKeyService.IsValidAPIKeyFormat(request.APIKey) {
// 		uc.log.Warn().Str("api_key", request.APIKey[:10]+"...").Msg("Formato de API Key inválido")
// 		return nil, errs.New("Formato de API Key inválido")
// 	}

// 	// Buscar la API Key en la base de datos
// 	apiKeyEntity, err := uc.repository.ValidateAPIKey(ctx, request.APIKey)
// 	if err != nil {
// 		uc.log.Error().Err(err).Str("api_key", request.APIKey[:10]+"...").Msg("Error al validar API Key en base de datos")
// 		return nil, errs.New("API Key inválida o no encontrada")
// 	}

// 	// Verificar que la API Key no esté revocada
// 	if apiKeyEntity.Revoked {
// 		uc.log.Warn().Uint("api_key_id", apiKeyEntity.ID).Msg("API Key revocada")
// 		return nil, errs.New("API Key revocada")
// 	}

// 	// Obtener información del usuario asociado
// 	userInfo, err := uc.repository.GetUserByID(ctx, apiKeyEntity.UserID)
// 	if err != nil {
// 		uc.log.Error().Err(err).Uint("user_id", apiKeyEntity.UserID).Msg("Error al obtener información del usuario")
// 		return nil, errs.New("Error al obtener información del usuario")
// 	}

// 	// Verificar que el usuario esté activo
// 	if !userInfo.IsActive {
// 		uc.log.Warn().Uint("user_id", apiKeyEntity.UserID).Msg("Usuario inactivo")
// 		return nil, errs.New("Usuario inactivo")
// 	}

// 	// Obtener roles del usuario
// 	userRoles, err := uc.repository.GetUserRoles(ctx, apiKeyEntity.UserID)
// 	if err != nil {
// 		uc.log.Error().Err(err).Uint("user_id", apiKeyEntity.UserID).Msg("Error al obtener roles del usuario")
// 		return nil, errs.New("Error al obtener roles del usuario")
// 	}

// 	// Convertir roles a strings
// 	var roleNames []string
// 	for _, role := range userRoles {
// 		roleNames = append(roleNames, role.Name)
// 	}

// 	// Actualizar último uso de la API Key
// 	err = uc.repository.UpdateAPIKeyLastUsed(ctx, apiKeyEntity.ID)
// 	if err != nil {
// 		uc.log.Error().Err(err).Uint("api_key_id", apiKeyEntity.ID).Msg("Error al actualizar último uso de API Key")
// 		// No fallamos la validación por este error, solo lo registramos
// 	}

// 	uc.log.Info().
// 		Uint("api_key_id", apiKeyEntity.ID).
// 		Uint("user_id", apiKeyEntity.UserID).
// 		Uint("business_id", apiKeyEntity.BusinessID).
// 		Msg("API Key validada exitosamente")

// 	return &dtos.ValidateAPIKeyResponse{
// 		Success:    true,
// 		Message:    "API Key válida",
// 		UserID:     apiKeyEntity.UserID,
// 		Email:      userInfo.Email,
// 		BusinessID: apiKeyEntity.BusinessID,
// 		Roles:      roleNames,
// 		APIKeyID:   apiKeyEntity.ID,
// 	}, nil
// }
