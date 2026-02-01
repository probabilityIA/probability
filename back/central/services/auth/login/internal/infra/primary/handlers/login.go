package authhandler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/login/internal/domain"
	"github.com/secamc93/probability/back/central/services/auth/login/internal/infra/primary/handlers/mapper"
	"github.com/secamc93/probability/back/central/services/auth/login/internal/infra/primary/handlers/request"
	"github.com/secamc93/probability/back/central/services/auth/login/internal/infra/primary/handlers/response"
	"github.com/secamc93/probability/back/central/shared/log"
)

// LoginHandler maneja la solicitud de login
//
//	@Summary		Autenticar usuario
//	@Description	Autentica un usuario con email y contraseña, retornando información del usuario y token de acceso
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		request.LoginRequest				true	"Credenciales de login"
//	@Success		200		{object}	response.LoginSuccessResponse		"Login exitoso"
//	@Failure		400		{object}	response.LoginBadRequestResponse	"Datos de entrada inválidos"
//	@Failure		401		{object}	response.LoginErrorResponse			"Credenciales inválidas"
//	@Failure		403		{object}	response.LoginErrorResponse			"Usuario inactivo"
//	@Failure		500		{object}	response.LoginErrorResponse			"Error interno del servidor"
//	@Router			/auth/login [post]
func (h *AuthHandler) LoginHandler(c *gin.Context) {
	ctx := log.WithFunctionCtx(c.Request.Context(), "LoginHandler")

	var loginRequest request.LoginRequest

	// Validar y bindear el request
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		h.logger.Error(ctx).Err(err).Msg("Error al validar request de login")
		c.JSON(http.StatusBadRequest, response.LoginBadRequestResponse{
			Error:   "Datos de entrada inválidos",
			Details: err.Error(),
		})
		return
	}

	// Convertir request a dominio
	domainRequest := domain.LoginRequest{
		Email:    loginRequest.Email,
		Password: loginRequest.Password,
	}

	// Ejecutar caso de uso
	domainResponse, err := h.usecase.Login(ctx, domainRequest)
	if err != nil {
		h.logger.Error(ctx).Err(err).Str("email", loginRequest.Email).Msg("Error en proceso de login")

		// Determinar el código de estado HTTP apropiado
		statusCode := http.StatusInternalServerError
		errorMessage := "Error interno del servidor"

		switch {
		case errors.Is(err, domain.ErrInvalidCredentials):
			statusCode = http.StatusUnauthorized
			errorMessage = domain.ErrInvalidCredentials.Error()
		case errors.Is(err, domain.ErrUserNotFound):
			statusCode = http.StatusUnauthorized
			errorMessage = domain.ErrInvalidCredentials.Error() // Generic message for security
		case errors.Is(err, domain.ErrUserInactive):
			statusCode = http.StatusForbidden
			errorMessage = domain.ErrUserInactive.Error()
		case errors.Is(err, domain.ErrEmailPasswordRequired):
			statusCode = http.StatusBadRequest
			errorMessage = domain.ErrEmailPasswordRequired.Error()
		}

		c.JSON(statusCode, response.LoginErrorResponse{
			Error: errorMessage,
		})
		return
	}

	// Convertir respuesta de dominio a response
	loginResponse := mapper.ToLoginResponse(domainResponse)

	// IMPORTANTE: Setear token como cookie HttpOnly para seguridad
	// SameSite=None permite que funcione en iframes (Shopify)
	c.SetCookie(
		"session_token",          // name
		domainResponse.Token,     // value
		7*24*60*60,               // maxAge: 7 días en segundos
		"/",                      // path
		"",                       // domain: vacío = current domain
		true,                     // secure: solo HTTPS
		true,                     // httpOnly: JavaScript no puede leer (seguridad)
	)
	c.SetSameSite(http.SameSiteNoneMode) // Para iframes de terceros (Shopify)

	// NO retornar el token en el JSON por seguridad
	// El token solo estará en la cookie HttpOnly
	loginResponse.Token = ""

	h.logger.Info(ctx).
		Str("email", loginRequest.Email).
		Uint("user_id", domainResponse.User.ID).
		Str("scope", domainResponse.Scope).
		Bool("is_super_admin", domainResponse.IsSuperAdmin).
		Msg("Login exitoso - Cookie HttpOnly seteada")

	// Retornar respuesta exitosa (sin token en JSON)
	c.JSON(http.StatusOK, response.LoginSuccessResponse{
		Success: true,
		Data:    *loginResponse,
	})
}
