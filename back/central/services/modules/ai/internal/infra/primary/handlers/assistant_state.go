package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/ai/internal/infra/primary/handlers/mappers"
)

func (h *Handlers) GetAssistantState(c *gin.Context) {
	scope, ok := accessScope(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "code": "unauthorized", "message": "Tu sesi\u00f3n no es v\u00e1lida. Vuelve a iniciar sesi\u00f3n."})
		return
	}

	state, err := h.uc.GetAssistantState(c.Request.Context(), scope.UserID)
	if err != nil {
		h.log.Error(c.Request.Context()).Err(err).Uint("user_id", scope.UserID).Msg("[ai.assistant] error leyendo el estado")
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "code": "internal_error", "message": "No se pudo leer el estado del asistente."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": mappers.FromState(state)})
}
