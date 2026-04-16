package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/infra/primary/handlers/request"
)

// SendManualReply maneja el endpoint POST /whatsapp/conversations/:id/reply
// Permite a un agente responder manualmente a un cliente de WhatsApp desde el dashboard.
// Requiere ventana de servicio activa (cliente escribió en las últimas 24h).
func (h *handler) SendManualReply(c *gin.Context) {
	ctx := c.Request.Context()
	conversationID := c.Param("id")

	var req request.ManualReplyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Error(ctx).Err(err).Msg("[ManualReply Handler] - request inválido")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request",
			"message": "Los datos de entrada son inválidos",
			"details": err.Error(),
		})
		return
	}

	// Extraer user_id del JWT para auditoría
	sentBy, _ := c.Get("user_id")
	sentByStr, _ := sentBy.(string)

	h.log.Info(ctx).
		Str("conversation_id", conversationID).
		Str("phone_number", req.PhoneNumber).
		Uint("business_id", req.BusinessID).
		Str("sent_by", sentByStr).
		Msg("[ManualReply Handler] - enviando reply manual")

	messageID, err := h.useCase.SendManualReply(
		ctx,
		conversationID,
		req.PhoneNumber,
		req.BusinessID,
		req.Text,
		sentByStr,
	)
	if err != nil {
		h.log.Error(ctx).Err(err).
			Str("conversation_id", conversationID).
			Msg("[ManualReply Handler] - error enviando reply manual")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "send_failed",
			"message": "Error al enviar el mensaje",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message_id": messageID,
		"status":     "sent",
	})
}
