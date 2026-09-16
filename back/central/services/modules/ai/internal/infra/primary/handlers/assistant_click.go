package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handlers) MarkAssistantClick(c *gin.Context) {
	scope, ok := accessScope(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "code": "unauthorized", "message": "Tu sesi\u00f3n no es v\u00e1lida. Vuelve a iniciar sesi\u00f3n."})
		return
	}

	if err := h.uc.MarkDestinationClicked(c.Request.Context(), scope.UserID, c.Param("id")); err != nil {
		h.respondRecordError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}
