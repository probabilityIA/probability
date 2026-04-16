package handlerintegrationtype

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/infra/primary/handlers/handlerintegrationtype/mapper"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/infra/primary/handlers/handlerintegrationtype/response"
)

// ListIntegrationCategoriesHandler godoc
//
//	@Summary		Listar categorías de integración
//	@Description	Obtiene todas las categorías de integración activas y visibles
//	@Tags			IntegrationCategories
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	response.IntegrationCategoryListResponse
//	@Failure		500	{object}	response.IntegrationErrorResponse
//	@Router			/integration-categories [get]
//	@Security		BearerAuth
func (h *IntegrationTypeHandler) ListIntegrationCategoriesHandler(c *gin.Context) {
	ctx := c.Request.Context()

	// Obtener categorías del use case
	categories, err := h.usecase.ListIntegrationCategories(ctx)
	if err != nil {
		h.logger.Error(ctx).Err(err).Msg("Error listing integration categories")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error al obtener las categorías de integración",
			"error":   err.Error(),
		})
		return
	}

	// Mapear a response
	categoriesResponse := make([]response.IntegrationCategoryResponse, 0, len(categories))
	for _, category := range categories {
		categoriesResponse = append(categoriesResponse, mapper.ToIntegrationCategoryResponse(category))
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Categorías de integración obtenidas exitosamente",
		"data":    categoriesResponse,
	})
}
