package handlers

import (
	"net/http"
	"strconv"

	"github.com/secamc93/probability/back/central/shared/log"

	"github.com/secamc93/probability/back/central/services/auth/actions/internal/domain"
	"github.com/secamc93/probability/back/central/services/auth/actions/internal/infra/primary/handlers/request"
	"github.com/secamc93/probability/back/central/services/auth/actions/internal/infra/primary/handlers/response"
	"github.com/secamc93/probability/back/central/services/auth/middleware"

	"github.com/gin-gonic/gin"
)

// UpdateActionHandler actualiza un action existente
//
//	@Summary		Actualizar action
//	@Description	Actualiza un action existente en el sistema
//	@Tags			Actions
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int							true	"ID del action"	minimum(1)
//	@Param			request	body		request.UpdateActionRequest	true	"Datos del action a actualizar"
//	@Success		200		{object}	map[string]interface{}		"Action actualizado exitosamente"
//	@Failure		400		{object}	map[string]interface{}		"Datos de entrada inválidos"
//	@Failure		401		{object}	map[string]interface{}		"No autorizado"
//	@Failure		403		{object}	map[string]interface{}		"Solo super usuarios pueden actualizar actions"
//	@Failure		404		{object}	map[string]interface{}		"Action no encontrado"
//	@Failure		409		{object}	map[string]interface{}		"Conflicto con action existente"
//	@Failure		500		{object}	map[string]interface{}		"Error interno del servidor"
//	@Router			/actions/{id} [put]
//	@Security		BearerAuth
func (h *ActionHandler) UpdateActionHandler(c *gin.Context) {
	ctx := log.WithFunctionCtx(c.Request.Context(), "UpdateActionHandler")

	// Validar que el usuario sea super admin
	if !middleware.IsSuperAdmin(c) {
		h.logger.Warn(ctx).Msg("Intento de actualización de action por usuario no super admin")
		c.JSON(http.StatusForbidden, response.ErrorResponse{
			Success: false,
			Message: "Solo los super usuarios pueden actualizar actions",
			Error:   "permisos insuficientes",
		})
		return
	}

	// Obtener el ID del action de los parámetros de la URL
	actionIDStr := c.Param("id")
	actionID, err := strconv.ParseUint(actionIDStr, 10, 32)
	if err != nil {
		h.logger.Error(ctx).Err(err).Str("action_id", actionIDStr).Msg("Error al parsear ID del action")
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Success: false,
			Message: "ID de action inválido",
			Error:   "El ID del action debe ser un número válido",
		})
		return
	}

	h.logger.Info(ctx).Uint64("action_id", actionID).Msg("Iniciando actualización de action")

	// Parsear el cuerpo de la petición
	var req request.UpdateActionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error(ctx).Err(err).Uint64("action_id", actionID).Msg("Error al parsear el cuerpo de la petición")
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Success: false,
			Message: "Datos de entrada inválidos",
			Error:   err.Error(),
		})
		return
	}

	h.logger.Info(ctx).
		Uint64("action_id", actionID).
		Str("name", req.Name).
		Msg("Datos de actualización de action recibidos")

	// Convertir a DTO de dominio
	updateDTO := domain.UpdateActionDTO{
		Name:        req.Name,
		Description: req.Description,
	}

	// Llamar al caso de uso
	result, err := h.usecase.UpdateAction(ctx, uint(actionID), updateDTO)
	if err != nil {
		h.logger.Error(ctx).Err(err).Uint64("action_id", actionID).Msg("Error al actualizar action")

		// Determinar el tipo de error y el código de estado HTTP
		statusCode := http.StatusInternalServerError
		errorMessage := "Error interno del servidor"

		if err.Error() == "action con ID "+actionIDStr+" no encontrado" {
			statusCode = http.StatusNotFound
			errorMessage = "Action no encontrado"
		} else if err.Error() == "ya existe otro action con el nombre '"+req.Name+"'" {
			statusCode = http.StatusConflict
			errorMessage = "Conflicto con action existente"
		} else if err.Error() == "el nombre del action es obligatorio" ||
			err.Error() == "el nombre del action no puede exceder 20 caracteres" ||
			err.Error() == "la descripción del action no puede exceder 255 caracteres" {
			statusCode = http.StatusBadRequest
			errorMessage = "Datos de entrada inválidos"
		}

		c.JSON(statusCode, response.ErrorResponse{
			Success: false,
			Message: errorMessage,
			Error:   err.Error(),
		})
		return
	}

	// Convertir a respuesta HTTP
	actionResponse := response.ActionResponse{
		ID:          result.ID,
		Name:        result.Name,
		Description: result.Description,
		CreatedAt:   result.CreatedAt,
		UpdatedAt:   result.UpdatedAt,
	}

	h.logger.Info(ctx).
		Uint64("action_id", actionID).
		Str("name", result.Name).
		Msg("Action actualizado exitosamente")

	c.JSON(http.StatusOK, response.UpdateActionResponse{
		Success: true,
		Message: "Action actualizado exitosamente",
		Data:    actionResponse,
	})
}
