package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/ai/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/ai/internal/infra/primary/handlers/request"
)

func (h *Handlers) SubmitAssistantFeedback(c *gin.Context) {
	scope, ok := accessScope(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "code": "unauthorized", "message": "Tu sesi\u00f3n no es v\u00e1lida. Vuelve a iniciar sesi\u00f3n."})
		return
	}

	var req request.AssistantFeedback
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "code": "invalid_feedback", "message": "La calificaci\u00f3n no es v\u00e1lida."})
		return
	}

	if err := h.uc.SubmitFeedback(c.Request.Context(), scope.UserID, c.Param("id"), *req.Value); err != nil {
		h.respondRecordError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handlers) respondRecordError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domainerrors.ErrInvalidMessageID):
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "code": "invalid_message_id", "message": "El mensaje no es v\u00e1lido."})
	case errors.Is(err, domainerrors.ErrInvalidFeedback):
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "code": "invalid_feedback", "message": "La calificaci\u00f3n no es v\u00e1lida."})
	case errors.Is(err, domainerrors.ErrMessageNotFound):
		c.JSON(http.StatusNotFound, gin.H{"success": false, "code": "message_not_found", "message": "No encontramos ese mensaje."})
	default:
		h.log.Error(c.Request.Context()).Err(err).Msg("[ai.assistant] error registrando evento del mensaje")
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "code": "internal_error", "message": "No se pudo guardar."})
	}
}
