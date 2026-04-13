package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/announcements/internal/infra/primary/handlers/mappers"
)

func (h *handler) Get(c *gin.Context) {
	id, ok := h.parseIDParam(c)
	if !ok {
		return
	}

	result, err := h.uc.GetAnnouncement(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": err.Error()})
		return
	}

	resp := mappers.EntityToResponse(result)
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}
