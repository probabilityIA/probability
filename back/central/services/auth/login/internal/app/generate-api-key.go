package app

// // GenerateAPIKey genera una API Key para un usuario específico
// func (uc *AuthUseCase) GenerateAPIKey(ctx context.Context, request domain.GenerateAPIKeyRequest) (*domain.GenerateAPIKeyResponse, error) {
// 	uc.log.Info().
// 		Uint("requester_id", request.RequesterID).
// 		Uint("target_user_id", request.UserID).
// 		Uint("business_id", request.BusinessID).
// 		Msg("Iniciando generación de API Key")

// 	// 1. Verificar que el solicitante sea super administrador (rol 1)
// 	requesterRoles, err := uc.repository.GetUserRoles(ctx, request.RequesterID)
// 	if err != nil {
// 		uc.log.Error().Err(err).Uint("requester_id", request.RequesterID).Msg("Error al obtener roles del solicitante")
// 		return nil, fmt.Errorf("error al verificar permisos del solicitante")
// 	}

// 	// Verificar si tiene el rol de super administrador
// 	hasSuperAdminRole := false
// 	for _, role := range requesterRoles {
// 		if role.ID == domain.SuperAdministrator {
// 			hasSuperAdminRole = true
// 			break
// 		}
// 	}

// 	if !hasSuperAdminRole {
// 		uc.log.Error().
// 			Uint("requester_id", request.RequesterID).
// 			Interface("roles", requesterRoles).
// 			Msg("Solicitante no tiene rol de super administrador")
// 		return nil, fmt.Errorf("solo super administradores pueden generar API Keys")
// 	}

// 	// 2. Verificar que el usuario objetivo existe y está activo
// 	targetUser, err := uc.repository.GetUserByID(ctx, request.UserID)
// 	if err != nil {
// 		uc.log.Error().Err(err).Uint("target_user_id", request.UserID).Msg("Error al obtener usuario objetivo")
// 		return nil, fmt.Errorf("usuario objetivo no encontrado")
// 	}

// 	if targetUser == nil {
// 		uc.log.Error().Uint("target_user_id", request.UserID).Msg("Usuario objetivo no encontrado")
// 		return nil, fmt.Errorf("usuario objetivo no encontrado")
// 	}

// 	if !targetUser.IsActive {
// 		uc.log.Error().Uint("target_user_id", request.UserID).Msg("Usuario objetivo está inactivo")
// 		return nil, fmt.Errorf("usuario objetivo está inactivo")
// 	}

// 	// 3. Verificar que el business existe y está activo
// 	// TODO: Implementar verificación de business cuando tengamos el repositorio
// 	// Por ahora, asumimos que el business existe

// 	// 4. Generar API Key segura
// 	apiKeyService := apikey.NewService()
// 	apiKey, err := apiKeyService.GenerateAPIKey()
// 	if err != nil {
// 		uc.log.Error().Err(err).Msg("Error al generar API Key")
// 		return nil, fmt.Errorf("error al generar API Key: %w", err)
// 	}

// 	// 5. Generar hash de la API Key
// 	keyHash, err := apiKeyService.HashAPIKey(apiKey)
// 	if err != nil {
// 		uc.log.Error().Err(err).Msg("Error al hashear API Key")
// 		return nil, fmt.Errorf("error al procesar API Key: %w", err)
// 	}

// 	// 6. Crear entidad API Key
// 	apiKeyEntity := domain.APIKey{
// 		UserID:      request.UserID,
// 		BusinessID:  request.BusinessID,
// 		CreatedByID: request.RequesterID,
// 		Name:        request.Name,
// 		Description: request.Description,
// 		RateLimit:   1000, // Límite por defecto
// 		CreatedAt:   time.Now(),
// 		UpdatedAt:   time.Now(),
// 	}

// 	// 7. Guardar en la base de datos
// 	apiKeyID, err := uc.repository.CreateAPIKey(ctx, apiKeyEntity, keyHash)
// 	if err != nil {
// 		uc.log.Error().Err(err).Msg("Error al guardar API Key en base de datos")
// 		return nil, fmt.Errorf("error al guardar API Key: %w", err)
// 	}

// 	// 8. Construir respuesta
// 	response := &domain.GenerateAPIKeyResponse{
// 		Success: true,
// 		Message: "API Key generada exitosamente",
// 		APIKey:  apiKey, // Solo se muestra una vez
// 		APIKeyInfo: domain.APIKeyInfo{
// 			ID:          apiKeyID,
// 			UserID:      request.UserID,
// 			BusinessID:  request.BusinessID,
// 			Name:        request.Name,
// 			Description: request.Description,
// 			RateLimit:   1000,
// 			CreatedAt:   time.Now(),
// 		},
// 	}

// 	uc.log.Info().
// 		Uint("requester_id", request.RequesterID).
// 		Uint("target_user_id", request.UserID).
// 		Uint("business_id", request.BusinessID).
// 		Uint("api_key_id", apiKeyID).
// 		Msg("API Key generada exitosamente")

// 	return response, nil
// }
