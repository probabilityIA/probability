package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handlers) MarkAssistantIntroSeen(c *gin.Context) {
	scope, ok := accessScope(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "code": "unauthorized", "message": "Tu sesi\u00f3n no es v\u00e1lida. Vuelve a iniciar sesi\u00f3n."})
		return
	}

	if err := h.uc.MarkIntroSeen(c.Request.Context(), scope.UserID); err != nil {
		h.log.Error(c.Request.Context()).Err(err).Uint("user_id", scope.UserID).Msg("[ai.assistant] error guardando la presentacion vista")
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "code": "internal_error", "message": "No se pudo guardar la preferencia."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}
