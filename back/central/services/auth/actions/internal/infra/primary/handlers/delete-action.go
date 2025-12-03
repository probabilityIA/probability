package handlers

import (
	"net/http"
	"strconv"

	"github.com/secamc93/probability/back/central/services/auth/actions/internal/infra/primary/handlers/response"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
	"github.com/secamc93/probability/back/central/shared/log"

	"github.com/gin-gonic/gin"
)

// DeleteActionHandler elimina un action por su ID
//
//	@Summary		Eliminar action
//	@Description	Elimina un action del sistema por su ID único. No se puede eliminar si tiene permisos asociados.
//	@Tags			Actions
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int						true	"ID del action"	minimum(1)
//	@Success		200	{object}	map[string]interface{}	"Action eliminado exitosamente"
//	@Failure		400	{object}	map[string]interface{}	"ID de action inválido"
//	@Failure		401	{object}	map[string]interface{}	"No autorizado"
//	@Failure		403	{object}	map[string]interface{}	"Solo super usuarios pueden eliminar actions"
//	@Failure		404	{object}	map[string]interface{}	"Action no encontrado"
//	@Failure		409	{object}	map[string]interface{}	"Action tiene permisos asociados"
//	@Failure		500	{object}	map[string]interface{}	"Error interno del servidor"
//	@Router			/actions/{id} [delete]
//	@Security		BearerAuth
func (h *ActionHandler) DeleteActionHandler(c *gin.Context) {
	ctx := log.WithFunctionCtx(c.Request.Context(), "DeleteActionHandler")

	// Validar que el usuario sea super admin
	if !middleware.IsSuperAdmin(c) {
		h.logger.Warn(ctx).Msg("Intento de eliminación de action por usuario no super admin")
		c.JSON(http.StatusForbidden, response.ErrorResponse{
			Success: false,
			Message: "Solo los super usuarios pueden eliminar actions",
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

	h.logger.Info(ctx).Uint64("action_id", actionID).Msg("Iniciando eliminación de action")

	// Llamar al caso de uso
	message, err := h.usecase.DeleteAction(ctx, uint(actionID))
	if err != nil {
		h.logger.Error(ctx).Err(err).Uint64("action_id", actionID).Msg("Error al eliminar action")

		// Determinar el tipo de error y el código de estado HTTP
		statusCode := http.StatusInternalServerError
		errorMessage := "Error interno del servidor"

		errMsg := err.Error()
		if errMsg == "action con ID "+actionIDStr+" no encontrado" {
			statusCode = http.StatusNotFound
			errorMessage = "Action no encontrado"
		} else if len(errMsg) > 20 && errMsg[:20] == "no se puede eliminar" {
			statusCode = http.StatusConflict
			errorMessage = errMsg
		}

		c.JSON(statusCode, response.ErrorResponse{
			Success: false,
			Message: errorMessage,
			Error:   err.Error(),
		})
		return
	}

	h.logger.Info(ctx).
		Uint64("action_id", actionID).
		Str("message", message).
		Msg("Action eliminado exitosamente")

	c.JSON(http.StatusOK, response.DeleteActionResponse{
		Success: true,
		Message: message,
	})
}
