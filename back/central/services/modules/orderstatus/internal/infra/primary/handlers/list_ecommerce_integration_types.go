package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
	"github.com/secamc93/probability/back/central/services/modules/orderstatus/internal/infra/primary/handlers/response"
)

// ListEcommerceIntegrationTypes godoc
// @Summary      Listar tipos de integración ecommerce disponibles
// @Description  Obtiene los tipos de integración de categoría ecommerce. Super admin ve todos; scope business solo ve los que tiene configurados.
// @Tags         Channel Statuses
// @Produce      json
// @Success      200  {object}  response.IntegrationTypesResponse
// @Failure      500  {object}  map[string]string
// @Router       /ecommerce-integration-types [get]
func (h *handler) ListEcommerceIntegrationTypes(c *gin.Context) {
	// businessID == 0 → super admin (ve todo); > 0 → scope de negocio
	businessID, _ := middleware.GetBusinessID(c)

	result, err := h.uc.ListEcommerceIntegrationTypes(c.Request.Context(), businessID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error al obtener tipos de integración ecommerce",
			"error":   err.Error(),
		})
		return
	}

	data := make([]response.IntegrationTypeResponse, len(result))
	for i, it := range result {
		data[i] = response.IntegrationTypeResponse{
			ID:       it.ID,
			Code:     it.Code,
			Name:     it.Name,
			ImageURL: it.ImageURL,
		}
	}

	c.JSON(http.StatusOK, response.IntegrationTypesResponse{
		Success: true,
		Message: "Tipos de integración ecommerce obtenidos exitosamente",
		Data:    data,
	})
}
