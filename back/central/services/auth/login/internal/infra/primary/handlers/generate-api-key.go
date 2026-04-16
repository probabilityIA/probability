package authhandler

// // GenerateAPIKeyHandler maneja la solicitud de generación de API Key
// // @Summary Generar API Key
// // @Description Genera una API Key para un usuario específico (solo super administradores)
// // @Tags Auth
// // @Accept json
// // @Produce json
// // @Security     BearerAuth
// // @Param request body request.GenerateAPIKeyRequest true "Datos para generar API Key"
// // @Success 200 {object} response.GenerateAPIKeySuccessResponse "API Key generada exitosamente"
// // @Failure 400 {object} response.GenerateAPIKeyErrorResponse "Datos de entrada inválidos"
// // @Failure 401 {object} response.GenerateAPIKeyErrorResponse "No autorizado"
// // @Failure 403 {object} response.GenerateAPIKeyErrorResponse "Acceso denegado - solo super administradores"
// // @Failure 404 {object} response.GenerateAPIKeyErrorResponse "Usuario o business no encontrado"
// // @Failure 500 {object} response.GenerateAPIKeyErrorResponse "Error interno del servidor"
// // @Router /auth/generate-api-key [post]
// func (h *AuthHandler) GenerateAPIKeyHandler(c *gin.Context) {
// 	var apiKeyRequest request.GenerateAPIKeyRequest

// 	// Validar y bindear el request
// 	if err := c.ShouldBindJSON(&apiKeyRequest); err != nil {
// 		h.logger.Error().Err(err).Msg("Error al validar request de generación de API Key")
// 		c.JSON(http.StatusBadRequest, response.GenerateAPIKeyErrorResponse{
// 			Error:   "Datos de entrada inválidos",
// 			Details: err.Error(),
// 		})
// 		return
// 	}

// 	// Obtener el ID del usuario que hace la solicitud desde el contexto
// 	requesterID, exists := middleware.GetUserID(c)
// 	if !exists {
// 		h.logger.Error().Msg("Usuario no autenticado")
// 		c.JSON(http.StatusUnauthorized, response.GenerateAPIKeyErrorResponse{
// 			Error: "Usuario no autenticado",
// 		})
// 		return
// 	}

// 	// Convertir request a dominio
// 	domainRequest := mapper.ToGenerateAPIKeyRequest(apiKeyRequest, requesterID)

// 	// Ejecutar caso de uso
// 	domainResponse, err := h.usecase.GenerateAPIKey(c.Request.Context(), domainRequest)
// 	if err != nil {
// 		h.logger.Error().Err(err).
// 			Uint("requester_id", requesterID).
// 			Uint("target_user_id", apiKeyRequest.UserID).
// 			Uint("business_id", apiKeyRequest.BusinessID).
// 			Msg("Error en proceso de generación de API Key")

// 		// Determinar el código de estado HTTP apropiado
// 		statusCode := http.StatusInternalServerError
// 		errorMessage := "Error interno del servidor"

// 		if err.Error() == "solo super administradores pueden generar API Keys" {
// 			statusCode = http.StatusForbidden
// 			errorMessage = "Acceso denegado: solo super administradores pueden generar API Keys"
// 		} else if err.Error() == "usuario objetivo no encontrado" || err.Error() == "usuario objetivo está inactivo" {
// 			statusCode = http.StatusNotFound
// 			errorMessage = err.Error()
// 		} else if err.Error() == "error al verificar permisos del solicitante" {
// 			statusCode = http.StatusUnauthorized
// 			errorMessage = "Error al verificar permisos"
// 		}

// 		c.JSON(statusCode, response.GenerateAPIKeyErrorResponse{
// 			Error: errorMessage,
// 		})
// 		return
// 	}

// 	// Convertir respuesta de dominio a response
// 	apiKeyResponse := mapper.ToGenerateAPIKeyResponse(domainResponse)

// 	h.logger.Info().
// 		Uint("requester_id", requesterID).
// 		Uint("target_user_id", apiKeyRequest.UserID).
// 		Uint("business_id", apiKeyRequest.BusinessID).
// 		Msg("API Key generada exitosamente")

// 	// Retornar respuesta exitosa
// 	c.JSON(http.StatusOK, response.GenerateAPIKeySuccessResponse{
// 		Success: true,
// 		Data:    apiKeyResponse,
// 	})
// }
