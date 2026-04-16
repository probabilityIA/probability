package handlers

import (
	"net/http"
	"strconv"

	"github.com/secamc93/probability/back/central/services/auth/resources/internal/infra/primary/handlers/response"

	"github.com/gin-gonic/gin"
)

// GetResourceByIDHandler obtiene un recurso por su ID
//
//	@Summary		Obtener recurso por ID
//	@Description	Obtiene un recurso específico del sistema por su ID único
//	@Tags			Resources
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int						true	"ID del recurso"	minimum(1)
//	@Success		200	{object}	map[string]interface{}	"Recurso obtenido exitosamente"
//	@Failure		400	{object}	map[string]interface{}	"ID de recurso inválido"
//	@Failure		401	{object}	map[string]interface{}	"No autorizado"
//	@Failure		404	{object}	map[string]interface{}	"Recurso no encontrado"
//	@Failure		500	{object}	map[string]interface{}	"Error interno del servidor"
//	@Router			/resources/{id} [get]
//	@Security		BearerAuth
func (h *ResourceHandler) GetResourceByIDHandler(c *gin.Context) {
	// Obtener el ID del recurso de los parámetros de la URL
	resourceIDStr := c.Param("id")
	resourceID, err := strconv.ParseUint(resourceIDStr, 10, 32)
	if err != nil {
		h.logger.Error().Err(err).Str("resource_id", resourceIDStr).Msg("Error al parsear ID del recurso")
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Success: false,
			Message: "ID de recurso inválido",
			Error:   "El ID del recurso debe ser un número válido",
		})
		return
	}

	h.logger.Info().Uint64("resource_id", resourceID).Msg("Iniciando obtención de recurso por ID")

	// Llamar al caso de uso
	result, err := h.usecase.GetResourceByID(c.Request.Context(), uint(resourceID))
	if err != nil {
		h.logger.Error().Err(err).Uint64("resource_id", resourceID).Msg("Error al obtener recurso por ID")

		// Determinar el tipo de error y el código de estado HTTP
		statusCode := http.StatusInternalServerError
		errorMessage := "Error interno del servidor"

		if err.Error() == "recurso con ID "+resourceIDStr+" no encontrado" {
			statusCode = http.StatusNotFound
			errorMessage = "Recurso no encontrado"
		}

		c.JSON(statusCode, response.ErrorResponse{
			Success: false,
			Message: errorMessage,
			Error:   err.Error(),
		})
		return
	}

	// Convertir a respuesta HTTP
	resourceResponse := response.ResourceResponse{
		ID:               result.ID,
		Name:             result.Name,
		Description:      result.Description,
		BusinessTypeID:   result.BusinessTypeID,
		BusinessTypeName: result.BusinessTypeName,
		CreatedAt:        result.CreatedAt,
		UpdatedAt:        result.UpdatedAt,
	}

	h.logger.Info().
		Uint64("resource_id", resourceID).
		Str("name", result.Name).
		Msg("Recurso obtenido exitosamente por ID")

	c.JSON(http.StatusOK, response.GetResourceByIDResponse{
		Success: true,
		Message: "Recurso obtenido exitosamente",
		Data:    resourceResponse,
	})
}
