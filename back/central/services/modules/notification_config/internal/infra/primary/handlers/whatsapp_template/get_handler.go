package whatsapp_template

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *handler) GetByID(c *gin.Context) {
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

	template, err := h.useCase.GetByID(c.Request.Context(), id, businessID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "error": err.Error()})
		return
	}
	if template == nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "plantilla no encontrada"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": template})
}
