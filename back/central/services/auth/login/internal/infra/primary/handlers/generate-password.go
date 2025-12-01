package authhandler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/login/internal/domain"
	"github.com/secamc93/probability/back/central/services/auth/login/internal/infra/primary/handlers/request"
	"github.com/secamc93/probability/back/central/services/auth/login/internal/infra/primary/handlers/response"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
	"github.com/secamc93/probability/back/central/shared/log"
)

// GeneratePasswordHandler maneja la solicitud de generar una nueva contraseña aleatoria
//
//	@Summary		Generar nueva contraseña aleatoria
//	@Description	Genera una nueva contraseña aleatoria. Si el usuario es super admin, puede especificar user_id para generar contraseña de otro usuario. Si no se envía user_id, se genera para el usuario autenticado. La contraseña solo se muestra una vez en esta respuesta.
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		request.GeneratePasswordRequest	false	"Request body (user_id opcional, solo para super usuarios)"
//	@Success		200		{object}	response.GeneratePasswordResponse	"Nueva contraseña generada exitosamente (solo se muestra una vez)"
//	@Failure		400		{object}	response.LoginErrorResponse		"Datos inválidos"
//	@Failure		401		{object}	response.LoginErrorResponse		"Token de acceso requerido"
//	@Failure		403		{object}	response.LoginErrorResponse		"No tienes permisos para generar contraseña de otro usuario o usuario inactivo"
//	@Failure		404		{object}	response.LoginErrorResponse		"Usuario no encontrado"
//	@Failure		500		{object}	response.LoginErrorResponse		"Error interno del servidor"
//	@Router			/auth/generate-password [post]
func (h *AuthHandler) GeneratePasswordHandler(c *gin.Context) {
	ctx := log.WithFunctionCtx(c.Request.Context(), "GeneratePasswordHandler")

	// Obtener el ID del usuario autenticado desde el middleware
	authenticatedUserID, exists := middleware.GetUserID(c)
	if !exists {
		h.logger.Error(ctx).Msg("Usuario no autenticado")
		c.JSON(http.StatusUnauthorized, response.LoginErrorResponse{
			Error: "Usuario no autenticado",
		})
		return
	}

	// Verificar si es super admin
	isSuperAdmin := middleware.IsSuperAdmin(c)

	// Parsear request body (opcional)
	var req request.GeneratePasswordRequest
	// Si hay contenido en el body, intentar parsearlo
	if c.Request.ContentLength > 0 {
		if err := c.ShouldBindJSON(&req); err != nil {
			h.logger.Error(ctx).Err(err).Msg("Error al validar request")
			c.JSON(http.StatusBadRequest, response.LoginErrorResponse{
				Error: "Datos de entrada inválidos: " + err.Error(),
			})
			return
		}
	}

	// Determinar el user_id a usar
	var targetUserID uint
	if req.UserID != nil {
		// Si se envió user_id, verificar que sea super admin
		if !isSuperAdmin {
			h.logger.Warn(ctx).
				Uint("authenticated_user_id", authenticatedUserID).
				Uint("requested_user_id", *req.UserID).
				Msg("Usuario no super admin intentando generar contraseña para otro usuario")
			c.JSON(http.StatusForbidden, response.LoginErrorResponse{
				Error: "No tienes permisos para generar contraseña de otro usuario",
			})
			return
		}
		targetUserID = *req.UserID
		h.logger.Info(ctx).
			Uint("authenticated_user_id", authenticatedUserID).
			Uint("target_user_id", targetUserID).
			Msg("Super admin generando contraseña para otro usuario")
	} else {
		// Si no se envió user_id, usar el del usuario autenticado
		targetUserID = authenticatedUserID
		h.logger.Info(ctx).
			Uint("user_id", targetUserID).
			Msg("Generando contraseña para usuario autenticado")
	}

	// Convertir request a dominio
	domainRequest := domain.GeneratePasswordRequest{
		UserID: targetUserID,
	}

	// Ejecutar caso de uso
	domainResponse, err := h.usecase.GeneratePassword(ctx, domainRequest)
	if err != nil {
		h.logger.Error(ctx).Err(err).Uint("target_user_id", targetUserID).Msg("Error en proceso de generación de contraseña")

		// Determinar el código de estado HTTP apropiado
		statusCode := http.StatusInternalServerError
		errorMessage := "Error interno del servidor"

		if err.Error() == "usuario no encontrado" {
			statusCode = http.StatusNotFound
			errorMessage = "Usuario no encontrado"
		} else if err.Error() == "usuario inactivo" {
			statusCode = http.StatusForbidden
			errorMessage = "Usuario inactivo"
		}

		c.JSON(statusCode, response.LoginErrorResponse{
			Error: errorMessage,
		})
		return
	}

	h.logger.Info(ctx).
		Uint("target_user_id", targetUserID).
		Uint("authenticated_user_id", authenticatedUserID).
		Bool("is_super_admin", isSuperAdmin).
		Str("email", domainResponse.Email).
		Msg("Nueva contraseña generada exitosamente")

	// Retornar respuesta exitosa con la contraseña (solo se muestra una vez)
	c.JSON(http.StatusOK, response.GeneratePasswordResponse{
		Success:  domainResponse.Success,
		Email:    domainResponse.Email,
		Password: domainResponse.Password,
		Message:  domainResponse.Message,
	})
}
