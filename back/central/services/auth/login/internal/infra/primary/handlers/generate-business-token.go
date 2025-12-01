package authhandler

import (
	"net/http"

	"github.com/secamc93/probability/back/central/shared/log"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/login/internal/infra/primary/handlers/request"
	"github.com/secamc93/probability/back/central/services/auth/login/internal/infra/primary/handlers/response"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
)

// GenerateBusinessTokenHandler maneja la solicitud para generar un token de business
//
//	@Summary		Generar token de business
//	@Description	Genera un token específico para un business basado en el token principal del usuario. Para super admins (scope platform), usar business_id = 0
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Security		BusinessTokenAuth
//	@Param			request	body		request.GenerateBusinessTokenRequest						true	"Datos del business (business_id = 0 para super admin)"
//	@Success		200		{object}	response.GenerateBusinessTokenSuccessResponse				"Token generado exitosamente"
//	@Failure		400		{object}	response.GenerateBusinessTokenErrorResponse				"Datos de entrada inválidos"
//	@Failure		401		{object}	response.GenerateBusinessTokenErrorResponse				"No autorizado"
//	@Failure		404		{object}	response.GenerateBusinessTokenErrorResponse				"Business no encontrado o sin acceso"
//	@Failure		500		{object}	response.GenerateBusinessTokenErrorResponse				"Error interno del servidor"
//	@Router			/auth/business-token [post]
func (h *AuthHandler) GenerateBusinessTokenHandler(c *gin.Context) {
	ctx := log.WithFunctionCtx(c.Request.Context(), "GenerateBusinessTokenHandler")

	// Obtener el userID del contexto (ya validado por el middleware JWT)
	userID, exists := middleware.GetUserID(c)
	if !exists || userID == 0 {
		h.logger.Error(ctx).Msg("No se pudo obtener el user_id del contexto")
		c.JSON(http.StatusUnauthorized, response.GenerateBusinessTokenErrorResponse{
			Error: "Token inválido o no autorizado",
		})
		return
	}

	// Parsear el body
	var businessTokenRequest request.GenerateBusinessTokenRequest
	if err := c.ShouldBindJSON(&businessTokenRequest); err != nil {
		h.logger.Error(ctx).Err(err).Msg("Error al validar request")
		c.JSON(http.StatusBadRequest, response.GenerateBusinessTokenErrorResponse{
			Error: "Datos de entrada inválidos: " + err.Error(),
		})
		return
	}

	// Para super admins, permitir business_id = 0
	// Para usuarios normales, validar que business_id sea válido
	if businessTokenRequest.BusinessID == 0 {
		h.logger.Info(ctx).
			Uint("user_id", userID).
			Msg("Solicitud de token con business_id = 0 (posible super admin)")
	} else {
		h.logger.Info(ctx).
			Uint("user_id", userID).
			Uint("business_id", businessTokenRequest.BusinessID).
			Msg("Solicitud de token para business específico")
	}

	// Ejecutar el caso de uso
	businessToken, err := h.usecase.GenerateBusinessToken(
		ctx,
		userID,
		businessTokenRequest.BusinessID,
	)

	if err != nil {
		h.logger.Error(ctx).Err(err).
			Uint("user_id", userID).
			Uint("business_id", businessTokenRequest.BusinessID).
			Msg("Error al generar business token")

		// Determinar el código de estado apropiado
		statusCode := http.StatusInternalServerError
		if err.Error() == "el usuario no tiene acceso a este business" {
			statusCode = http.StatusForbidden
		} else if err.Error() == "business no encontrado" {
			statusCode = http.StatusNotFound
		} else if err.Error() == "usuario no encontrado" || err.Error() == "usuario inactivo" {
			statusCode = http.StatusUnauthorized
		}

		c.JSON(statusCode, response.GenerateBusinessTokenErrorResponse{
			Error: err.Error(),
		})
		return
	}

	h.logger.Info(ctx).
		Uint("user_id", userID).
		Uint("business_id", businessTokenRequest.BusinessID).
		Msg("Business token generado exitosamente")

	// Retornar respuesta exitosa
	c.JSON(http.StatusOK, response.GenerateBusinessTokenSuccessResponse{
		Success: true,
		Message: "Business token generado exitosamente",
		Data: response.BusinessTokenResponse{
			Token: businessToken,
		},
	})
}
