package handlers

import (
	"net/http"

	"github.com/secamc93/probability/back/central/services/auth/actions/internal/infra/primary/handlers/request"
	"github.com/secamc93/probability/back/central/services/auth/actions/internal/infra/primary/handlers/response"
	"github.com/secamc93/probability/back/central/shared/log"

	"github.com/gin-gonic/gin"
)

// GetActionsHandler obtiene todos los actions con filtros y paginación
//
//	@Summary		Obtener actions
//	@Description	Obtiene una lista paginada de actions del sistema con opciones de filtrado
//	@Tags			Actions
//	@Accept			json
//	@Produce		json
//	@Param			page		query		int						false	"Número de página"		default(1)	minimum(1)
//	@Param			page_size	query		int						false	"Tamaño de página"		default(10)	minimum(1)	maximum(100)
//	@Param			name		query		string					false	"Filtrar por nombre (búsqueda parcial)"
//	@Success		200			{object}	map[string]interface{}	"Lista de actions obtenida exitosamente"
//	@Failure		400			{object}	map[string]interface{}	"Parámetros de consulta inválidos"
//	@Failure		401			{object}	map[string]interface{}	"No autorizado"
//	@Failure		500			{object}	map[string]interface{}	"Error interno del servidor"
//	@Router			/actions [get]
//	@Security		BearerAuth
func (h *ActionHandler) GetActionsHandler(c *gin.Context) {
	ctx := log.WithFunctionCtx(c.Request.Context(), "GetActionsHandler")

	// Parsear parámetros de consulta
	var req request.GetActionsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		h.logger.Error(ctx).Err(err).Msg("Error al parsear parámetros de consulta")
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Success: false,
			Message: "Parámetros de consulta inválidos",
			Error:   err.Error(),
		})
		return
	}

	// Llamar al caso de uso
	result, err := h.usecase.GetActions(ctx, req.Page, req.PageSize, req.Name)
	if err != nil {
		h.logger.Error(ctx).Err(err).Msg("Error al obtener actions")
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Success: false,
			Message: "Error interno del servidor",
			Error:   err.Error(),
		})
		return
	}

	// Convertir a respuesta HTTP
	var actionResponses []response.ActionResponse
	for _, action := range result.Actions {
		actionResponses = append(actionResponses, response.ActionResponse{
			ID:          action.ID,
			Name:        action.Name,
			Description: action.Description,
			CreatedAt:   action.CreatedAt,
			UpdatedAt:   action.UpdatedAt,
		})
	}

	listResponse := response.ActionListResponse{
		Actions:    actionResponses,
		Total:      result.Total,
		Page:       result.Page,
		PageSize:   result.PageSize,
		TotalPages: result.TotalPages,
	}

	h.logger.Info(ctx).
		Int64("total", result.Total).
		Int("returned", len(actionResponses)).
		Int("page", result.Page).
		Msg("Actions obtenidos exitosamente")

	c.JSON(http.StatusOK, response.GetActionsResponse{
		Success: true,
		Message: "Actions obtenidos exitosamente",
		Data:    listResponse,
	})
}
