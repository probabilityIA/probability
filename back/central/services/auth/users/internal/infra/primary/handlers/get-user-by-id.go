package handlers

import (
	"net/http"

	"github.com/secamc93/probability/back/central/services/auth/middleware"
	"github.com/secamc93/probability/back/central/services/auth/users/internal/infra/primary/handlers/mapper"
	"github.com/secamc93/probability/back/central/services/auth/users/internal/infra/primary/handlers/request"
	"github.com/secamc93/probability/back/central/services/auth/users/internal/infra/primary/handlers/response"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/shared/log"
)

// GetUserByIDHandler maneja la solicitud de obtener un usuario por ID
//
//	@Summary		Obtener usuario por ID
//	@Description	Obtiene un usuario específico por su ID con sus roles y businesses
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int								true	"ID del usuario"	minimum(1)
//	@Success		200	{object}	response.UserSuccessResponse	"Usuario obtenido exitosamente"
//	@Failure		400	{object}	response.UserErrorResponse		"ID inválido"
//	@Failure		401	{object}	response.UserErrorResponse		"Token de acceso requerido"
//	@Failure		404	{object}	response.UserErrorResponse		"Usuario no encontrado"
//	@Failure		500	{object}	response.UserErrorResponse		"Error interno del servidor"
//	@Router			/users/{id} [get]
func (h *handlers) GetUserByIDHandler(c *gin.Context) {
	ctx := log.WithFunctionCtx(c.Request.Context(), "GetUserByIDHandler")
	var req request.GetUserByIDRequest

	// Binding automático con validaciones para parámetros de URL
	if err := c.ShouldBindUri(&req); err != nil {
		h.logger.Error(ctx).Err(err).Msg("Error al validar ID del usuario")
		c.JSON(http.StatusBadRequest, response.UserErrorResponse{
			Error: "ID inválido: " + err.Error(),
		})
		return
	}

	// Verificar acceso: super admin puede ver cualquier usuario, usuario normal solo los de su business
	isSuperAdmin := middleware.IsSuperAdmin(c)
	if !isSuperAdmin {
		tokenBusinessID, ok := middleware.GetBusinessID(c)
		if !ok {
			h.logger.Error(ctx).Msg("Business ID no disponible en token")
			c.JSON(http.StatusUnauthorized, response.UserErrorResponse{
				Error: "Token inválido: business_id no disponible",
			})
			return
		}
		h.logger.Info(ctx).Uint("requested_user_id", req.ID).Uint("token_business_id", tokenBusinessID).Msg("Verificando acceso a usuario")
		// TODO: Aquí deberíamos validar que el usuario pertenezca al business del token
		// Por ahora, el caso de uso lo manejará internamente
	}

	h.logger.Info(ctx).Uint("id", req.ID).Bool("is_super_admin", isSuperAdmin).Msg("Iniciando solicitud para obtener usuario por ID")

	user, err := h.usecase.GetUserByID(ctx, req.ID)
	if err != nil {
		h.logger.Error(ctx).Uint("id", req.ID).Err(err).Msg("Error al obtener usuario por ID desde el caso de uso")

		statusCode := http.StatusInternalServerError
		errorMessage := "Error interno del servidor"

		if err.Error() == "usuario no encontrado" {
			statusCode = http.StatusNotFound
			errorMessage = "Usuario no encontrado"
		}

		c.JSON(statusCode, response.UserErrorResponse{
			Error: errorMessage,
		})
		return
	}

	userResponse := mapper.ToUserResponse(*user)

	h.logger.Info(ctx).Uint("id", req.ID).Msg("Usuario obtenido exitosamente")
	c.JSON(http.StatusOK, response.UserSuccessResponse{
		Success: true,
		Data:    userResponse,
	})
}
