package handlerintegrations

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/infra/primary/handlers/handlerintegrations/response"
)

// CheckIntegrationExistsHandler verifica si existe una integración activa por tipo y business_id
//
//	@Summary		Verificar existencia de integración
//	@Description	Verifica si existe una integración activa. Acepta integration_type_id o integration_type_code.
//	@Tags			Integrations
//	@Produce		json
//	@Security		BearerAuth
//	@Param			integration_type_id		query	int		false	"ID del tipo de integración"
//	@Param			integration_type_code	query	string	false	"Código del tipo de integración (ej: whatsapp)"
//	@Param			business_id				query	int		false	"ID del business (opcional)"
//	@Success		200	{object}	map[string]interface{}
//	@Failure		400	{object}	map[string]interface{}
//	@Failure		500	{object}	map[string]interface{}
func (h *IntegrationHandler) CheckIntegrationExistsHandler(c *gin.Context) {
	var businessID *uint
	if businessIDStr := c.Query("business_id"); businessIDStr != "" {
		if id, err := strconv.ParseUint(businessIDStr, 10, 32); err == nil {
			bid := uint(id)
			businessID = &bid
		}
	}

	var exists bool
	var err error

	// Opción 1: por integration_type_id directo
	if idStr := c.Query("integration_type_id"); idStr != "" {
		id, parseErr := strconv.ParseUint(idStr, 10, 32)
		if parseErr != nil {
			c.JSON(http.StatusBadRequest, response.IntegrationErrorResponse{
				Success: false,
				Message: "integration_type_id inválido",
				Error:   parseErr.Error(),
			})
			return
		}
		exists, err = h.usecase.HasActiveIntegration(c.Request.Context(), uint(id), businessID)
	} else if code := c.Query("integration_type_code"); code != "" {
		// Opción 2: resolver por código
		exists, err = h.usecase.HasActiveIntegrationByCode(c.Request.Context(), code, businessID)
	} else {
		c.JSON(http.StatusBadRequest, response.IntegrationErrorResponse{
			Success: false,
			Message: "Se requiere integration_type_id o integration_type_code",
		})
		return
	}

	if err != nil {
		h.logger.Error().Err(err).Msg("Error verificando existencia de integración")
		c.JSON(http.StatusInternalServerError, response.IntegrationErrorResponse{
			Success: false,
			Message: "Error al verificar integración",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"exists":  exists,
	})
}
