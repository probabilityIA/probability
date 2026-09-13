package whatsapp_template

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/dtos"
)

func (h *handler) ListFlowGroups(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "business_id es requerido para super admin"})
		return
	}

	flows, err := h.useCase.ListFlowGroups(c.Request.Context(), businessID)
	if err != nil {
		h.logger.Error().Err(err).Uint("business_id", businessID).Msg("Error listando los flujos")
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": flows})
}

func (h *handler) CreateFlowGroup(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "business_id es requerido para super admin"})
		return
	}

	var body dtos.CreateFlowDTO
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "cuerpo invalido"})
		return
	}

	body.BusinessID = businessID

	flow, err := h.useCase.CreateFlowGroup(c.Request.Context(), body)
	if err != nil {
		h.logger.Error().Err(err).Msg("Error creando el flujo")
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": flow})
}

func (h *handler) UpdateFlowGroup(c *gin.Context) {
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

	var body dtos.UpdateFlowDTO
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "cuerpo invalido"})
		return
	}

	body.ID = id
	body.BusinessID = businessID

	flow, err := h.useCase.UpdateFlowGroup(c.Request.Context(), body)
	if err != nil {
		h.logger.Error().Err(err).Uint("id", id).Msg("Error actualizando el flujo")
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": flow})
}

func (h *handler) DeleteFlowGroup(c *gin.Context) {
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

	if err := h.useCase.DeleteFlowGroup(c.Request.Context(), id, businessID); err != nil {
		h.logger.Error().Err(err).Uint("id", id).Msg("Error eliminando el flujo")
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *handler) ListFlowTransitions(c *gin.Context) {
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

	transitions, err := h.useCase.ListFlowTransitions(c.Request.Context(), id, businessID)
	if err != nil {
		h.logger.Error().Err(err).Uint("id", id).Msg("Error listando las transiciones del flujo")
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": transitions})
}
