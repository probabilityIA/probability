package businesshandler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/bussines/internal/infra/primary/controllers/businesshandler/response"
)

// GetBusinessesSimple godoc
// @Summary Obtener lista simple de negocios
// @Description Retorna solo ID y nombre de negocios activos para dropdowns/selectores
// @Tags businesses
// @Produce json
// @Success 200 {object} response.GetBusinessesSimpleResponse
// @Failure 500 {object} map[string]interface{}
// @Router /businesses/simple [get]
func (h *BusinessHandler) GetBusinessesSimple(c *gin.Context) {
	// Parámetros para obtener todos los businesses activos sin paginación
	page := 1
	perPage := 1000 // Suficiente para obtener todos los businesses
	name := ""      // Sin filtro de nombre
	isActive := true
	isActivePtr := &isActive

	// Obtener todos los negocios activos
	businesses, _, err := h.usecase.GetBusinesses(c.Request.Context(), page, perPage, name, nil, isActivePtr)
	if err != nil {
		h.logger.Error().Err(err).Msg("Error getting businesses for simple list")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error al obtener negocios",
			"error":   err.Error(),
		})
		return
	}

	// Mapear solo a formato simple (id, name)
	simpleBusinesses := make([]response.BusinessSimpleResponse, 0, len(businesses))
	for _, business := range businesses {
		simpleBusinesses = append(simpleBusinesses, response.BusinessSimpleResponse{
			ID:   business.ID,
			Name: business.Name,
		})
	}

	c.JSON(http.StatusOK, response.GetBusinessesSimpleResponse{
		Success: true,
		Message: "Negocios obtenidos exitosamente",
		Data:    simpleBusinesses,
	})
}
