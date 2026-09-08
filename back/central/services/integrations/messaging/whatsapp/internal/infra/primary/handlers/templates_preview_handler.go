package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *handler) PreviewTemplates(c *gin.Context) {
	businessID := c.GetUint("business_id")
	if businessID == 0 {
		if param := c.Query("business_id"); param != "" {
			if id, err := strconv.ParseUint(param, 10, 64); err == nil && id > 0 {
				businessID = uint(id)
			}
		}
	}

	if businessID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "business_id es requerido para super admin"})
		return
	}

	eventCode := c.Query("event_code")
	if eventCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "event_code es requerido"})
		return
	}

	previews, err := h.templatesUseCase.PreviewByEvent(c.Request.Context(), businessID, eventCode)
	if err != nil {
		h.log.Error().Err(err).Str("event_code", eventCode).
			Msg("Error consultando la vista previa de plantillas")
		c.JSON(http.StatusBadGateway, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": previews})
}
