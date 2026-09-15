package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/infra/primary/handlers/request"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/infra/primary/handlers/response"
)

func (h *handler) SendTemplate(c *gin.Context) {
	ctx := c.Request.Context()

	var req request.SendTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Error(ctx).Err(err).Msg("[Template Handler] - error validando request")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request",
			"message": "Los datos de entrada son inválidos",
			"details": err.Error(),
		})
		return
	}

	businessID, ok := resolveBusinessID(c, req.BusinessID)
	if !ok {
		return
	}

	h.log.Info(ctx).
		Str("template_name", req.TemplateName).
		Str("phone_number", req.PhoneNumber).
		Str("order_number", req.OrderNumber).
		Uint("business_id", businessID).
		Msg("[Template Handler] - procesando solicitud de envío de plantilla")

	if req.Variables == nil {
		req.Variables = make(map[string]string)
	}

	messageID, err := h.useCase.SendTemplate(
		ctx,
		req.TemplateName,
		req.PhoneNumber,
		req.Variables,
		req.OrderNumber,
		businessID,
	)

	if err != nil {
		h.log.Error(ctx).Err(err).
			Str("template_name", req.TemplateName).
			Str("phone_number", req.PhoneNumber).
			Msg("[Template Handler] - error enviando plantilla")

		statusCode := http.StatusInternalServerError
		errorType := "internal_error"

		errorMsg := err.Error()
		if strings.Contains(errorMsg, "plantilla no encontrada") {
			statusCode = http.StatusBadRequest
			errorType = "template_not_found"
		} else if strings.Contains(errorMsg, "variable") && strings.Contains(errorMsg, "faltante") {
			statusCode = http.StatusBadRequest
			errorType = "missing_variable"
		} else if strings.Contains(errorMsg, "número de teléfono inválido") {
			statusCode = http.StatusBadRequest
			errorType = "invalid_phone_number"
		}

		c.JSON(statusCode, gin.H{
			"error":   errorType,
			"message": "Error al enviar plantilla de WhatsApp",
			"details": err.Error(),
		})
		return
	}

	h.log.Info(ctx).
		Str("message_id", messageID).
		Str("template_name", req.TemplateName).
		Msg("[Template Handler] - plantilla enviada exitosamente")

	c.JSON(http.StatusOK, response.SendTemplateResponse{
		MessageID: messageID,
		Status:    "sent",
	})
}
