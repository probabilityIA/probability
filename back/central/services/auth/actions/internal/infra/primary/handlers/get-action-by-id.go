package handlers

import (
	"net/http"
	"strconv"

	"github.com/secamc93/probability/back/central/services/auth/actions/internal/infra/primary/handlers/response"
	"github.com/secamc93/probability/back/central/shared/log"

	"github.com/gin-gonic/gin"
)

// GetActionByIDHandler obtiene un action por su ID
//
//	@Summary		Obtener action por ID
//	@Description	Obtiene un action específico del sistema por su ID único
//	@Tags			Actions
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int						true	"ID del action"	minimum(1)
//	@Success		200	{object}	map[string]interface{}	"Action obtenido exitosamente"
//	@Failure		400	{object}	map[string]interface{}	"ID de action inválido"
//	@Failure		401	{object}	map[string]interface{}	"No autorizado"
//	@Failure		404	{object}	map[string]interface{}	"Action no encontrado"
//	@Failure		500	{object}	map[string]interface{}	"Error interno del servidor"
//	@Router			/actions/{id} [get]
//	@Security		BearerAuth
func (h *ActionHandler) GetActionByIDHandler(c *gin.Context) {
	ctx := log.WithFunctionCtx(c.Request.Context(), "GetActionByIDHandler")

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

	h.logger.Info(ctx).Uint64("action_id", actionID).Msg("Iniciando obtención de action por ID")

	// Llamar al caso de uso
	result, err := h.usecase.GetActionByID(ctx, uint(actionID))
	if err != nil {
		h.logger.Error(ctx).Err(err).Uint64("action_id", actionID).Msg("Error al obtener action por ID")

		// Determinar el tipo de error y el código de estado HTTP
		statusCode := http.StatusInternalServerError
		errorMessage := "Error interno del servidor"

		if err.Error() == "action con ID "+actionIDStr+" no encontrado" {
			statusCode = http.StatusNotFound
			errorMessage = "Action no encontrado"
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
		Msg("Action obtenido exitosamente por ID")

	c.JSON(http.StatusOK, response.GetActionByIDResponse{
		Success: true,
		Message: "Action obtenido exitosamente",
		Data:    actionResponse,
	})
}
