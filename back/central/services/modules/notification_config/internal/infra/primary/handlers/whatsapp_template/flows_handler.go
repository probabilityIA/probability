package whatsapp_template

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/dtos"
)

func (h *handler) ListAllFlows(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "business_id es requerido para super admin"})
		return
	}

	flows, err := h.useCase.ListBusinessFlows(c.Request.Context(), businessID)
	if err != nil {
		h.logger.Error().Err(err).Uint("business_id", businessID).Msg("Error listando flujos de plantillas")
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": flows})
}

func (h *handler) GetFlows(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "business_id es requerido para super admin"})
		return
	}

	id, ok := parseID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "id invalido"})
		return
	}

	flows, err := h.useCase.ListFlows(c.Request.Context(), id, businessID)
	if err != nil {
		h.logger.Error().Err(err).Uint("id", id).Msg("Error listando flujos de la plantilla")
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": flows})
}

func (h *handler) ReplaceFlows(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "business_id es requerido para super admin"})
		return
	}

	id, ok := parseID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "id invalido"})
		return
	}

	var body dtos.ReplaceTemplateFlowsDTO
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "cuerpo invalido"})
		return
	}

	body.BusinessID = businessID
	body.SourceTemplateID = id

	flows, err := h.useCase.ReplaceFlows(c.Request.Context(), body)
	if err != nil {
		h.logger.Error().Err(err).Uint("id", id).Msg("Error guardando los flujos de la plantilla")
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": flows})
}
