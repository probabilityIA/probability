package authhandler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/secamc93/probability/back/central/services/auth/login/internal/infra/primary/handlers/mapper"
	"github.com/secamc93/probability/back/central/services/auth/login/internal/infra/primary/handlers/response"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
	"github.com/secamc93/probability/back/central/shared/log"
)

// GetUserRolesPermissionsHandler maneja la solicitud de obtener roles y permisos del usuario
//
//	@Summary		Obtener roles y permisos del usuario
//	@Description	Obtiene los roles y permisos del usuario autenticado desde el token
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	response.UserRolesPermissionsSuccessResponse	"Roles y permisos obtenidos exitosamente"
//	@Failure		401	{object}	response.LoginErrorResponse						"Token de acceso requerido"
//	@Failure		404	{object}	response.LoginErrorResponse						"Usuario no encontrado"
//	@Failure		500	{object}	response.LoginErrorResponse						"Error interno del servidor"
//	@Router			/auth/roles-permissions [get]
func (h *AuthHandler) GetUserRolesPermissionsHandler(c *gin.Context) {
	ctx := log.WithFunctionCtx(c.Request.Context(), "GetUserRolesPermissionsHandler")

	// Obtener el ID del usuario autenticado desde el middleware
	userID, exists := middleware.GetUserID(c)
	if !exists {
		h.logger.Error(ctx).Msg("Usuario no autenticado")
		c.JSON(http.StatusUnauthorized, response.LoginErrorResponse{
			Error: "Usuario no autenticado",
		})
		return
	}

	// Obtener el business_id desde el token (ya validado por el middleware)
	businessID, exists := middleware.GetBusinessID(c)
	if !exists {
		h.logger.Error(ctx).Msg("Business ID no encontrado en el token")
		c.JSON(http.StatusUnauthorized, response.LoginErrorResponse{
			Error: "Token inválido - Business ID no encontrado",
		})
		return
	}

	h.logger.Info(ctx).Uint("user_id", uint(userID)).Uint("business_id", businessID).Msg("Obteniendo roles y permisos desde business_staff")

	// Obtener el token del header de autorización (ya validado por el middleware)
	token := c.GetHeader("Authorization")
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	// Ejecutar caso de uso para obtener roles y permisos
	rolesPermissions, err := h.usecase.GetUserRolesPermissions(ctx, uint(userID), businessID, token)
	if err != nil {
		h.logger.Error(ctx).Err(err).Uint("user_id", uint(userID)).Msg("Error al obtener roles y permisos del usuario")

		// Determinar el código de estado HTTP apropiado
		statusCode := http.StatusInternalServerError
		errorMessage := "Error interno del servidor"

		if err.Error() == "usuario no encontrado" {
			statusCode = http.StatusNotFound
			errorMessage = "Usuario no encontrado"
		} else if err.Error() == "token inválido" {
			statusCode = http.StatusUnauthorized
			errorMessage = "Token inválido"
		} else if err.Error() == "acceso denegado" {
			statusCode = http.StatusForbidden
			errorMessage = "Acceso denegado"
		}

		c.JSON(statusCode, response.LoginErrorResponse{
			Error: errorMessage,
		})
		return
	}

	// Convertir respuesta de dominio a response
	rolesPermissionsResponse := mapper.ToUserRolesPermissionsResponse(rolesPermissions)

	h.logger.Info(ctx).
		Uint("user_id", uint(userID)).
		Uint("business_id", businessID).
		Bool("is_super", rolesPermissions.IsSuper).
		Int("permissions_count", len(rolesPermissions.Permissions)).
		Msg("Roles y permisos obtenidos exitosamente")

	// Retornar respuesta exitosa
	c.JSON(http.StatusOK, response.UserRolesPermissionsSuccessResponse{
		Success: true,
		Data:    rolesPermissionsResponse,
	})
}
