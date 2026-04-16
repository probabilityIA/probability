package authhandler

// // ValidateAPIKeyHandler maneja la validación de API Keys
// func (h *AuthHandler) ValidateAPIKeyHandler(c *gin.Context) {
// 	// Obtener la API Key del header o query parameter
// 	apiKey := c.GetHeader("X-API-Key")
// 	if apiKey == "" {
// 		apiKey = c.Query("api_key")
// 	}

// 	if apiKey == "" {
// 		errorResponse := mapper.ToValidateAPIKeyErrorResponse("API Key requerida")
// 		c.JSON(http.StatusBadRequest, errorResponse)
// 		return
// 	}

// 	// Crear la solicitud para el caso de uso
// 	request := domain.ValidateAPIKeyRequest{
// 		APIKey: apiKey,
// 	}

// 	// Llamar al caso de uso
// 	domainResponse, err := h.usecase.ValidateAPIKey(c.Request.Context(), request)
// 	if err != nil {
// 		h.logger.Error().Err(err).Str("api_key", apiKey[:10]+"...").Msg("Error al validar API Key")
// 		errorResponse := mapper.ToValidateAPIKeyErrorResponse(err.Error())
// 		c.JSON(http.StatusUnauthorized, errorResponse)
// 		return
// 	}

// 	// Convertir la respuesta del dominio a la respuesta HTTP
// 	httpResponse := mapper.ToValidateAPIKeySuccessResponse(domainResponse)
// 	c.JSON(http.StatusOK, httpResponse)
// }
