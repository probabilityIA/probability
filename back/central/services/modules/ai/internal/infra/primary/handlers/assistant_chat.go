package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/ai/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/ai/internal/infra/primary/handlers/mappers"
	"github.com/secamc93/probability/back/central/services/modules/ai/internal/infra/primary/handlers/request"
)

func (h *Handlers) AssistantChat(c *gin.Context) {
	scope, ok := accessScope(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "code": "unauthorized", "message": "Tu sesi\u00f3n no es v\u00e1lida. Vuelve a iniciar sesi\u00f3n."})
		return
	}

	var req request.AssistantChat
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "code": "invalid_message", "message": "El mensaje no tiene un formato v\u00e1lido."})
		return
	}

	reply, err := h.uc.Chat(c.Request.Context(), mappers.ToChatInput(scope, req))
	if err != nil {
		h.respondChatError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": mappers.FromReply(reply)})
}

func (h *Handlers) respondChatError(c *gin.Context, err error) {
	var rateLimited *domainerrors.RateLimitedError
	switch {
	case errors.As(err, &rateLimited):
		c.JSON(http.StatusTooManyRequests, gin.H{
			"success":  false,
			"code":     "rate_limited",
			"message":  "Llegaste al l\u00edmite de mensajes por hora.",
			"limit":    rateLimited.Limit,
			"reset_at": rateLimited.ResetAt,
		})
	case errors.Is(err, domainerrors.ErrMessageTooLong):
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "code": "message_too_long", "message": "El mensaje es muy largo. Escr\u00edbelo en menos de 1.000 caracteres."})
	case errors.Is(err, domainerrors.ErrEmptyConversation):
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "code": "invalid_message", "message": "Escribe una pregunta para V\u00eda."})
	case errors.Is(err, domainerrors.ErrModelUnavailable):
		c.JSON(http.StatusServiceUnavailable, gin.H{"success": false, "code": "assistant_unavailable", "message": "No pude conectarme para responder. Intenta de nuevo en un momento."})
	default:
		h.log.Error(c.Request.Context()).Err(err).Msg("[ai.assistant] error atendiendo el chat")
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "code": "internal_error", "message": "No se pudo procesar el mensaje."})
	}
}
