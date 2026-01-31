package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/whatsApp/internal/infra/primary/handlers/request"
	"github.com/secamc93/probability/back/central/services/integrations/whatsApp/internal/infra/primary/handlers/response"
)

// SendTemplate maneja el endpoint POST /integrations/whatsapp/send-template
// @Summary Envía una plantilla de WhatsApp
// @Description Envía una plantilla de WhatsApp con variables dinámicas y botones opcionales
// @Tags WhatsApp
// @Accept json
// @Produce json
// @Param request body SendTemplateRequest true "Datos de la plantilla a enviar"
// @Success 200 {object} SendTemplateResponse
// @Failure 400 {object} map[string]interface{} "Error de validación"
// @Failure 500 {object} map[string]interface{} "Error interno del servidor"
// @Router /integrations/whatsapp/send-template [post]
func (h *Handler) SendTemplate(c *gin.Context) {
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

	h.log.Info(ctx).
		Str("template_name", req.TemplateName).
		Str("phone_number", req.PhoneNumber).
		Str("order_number", req.OrderNumber).
		Msg("[Template Handler] - procesando solicitud de envío de plantilla")

	// Inicializar variables si es nil
	if req.Variables == nil {
		req.Variables = make(map[string]string)
	}

	// Enviar plantilla
	messageID, err := h.useCase.SendTemplate(
		ctx,
		req.TemplateName,
		req.PhoneNumber,
		req.Variables,
		req.OrderNumber,
		req.BusinessID,
	)

	if err != nil {
		h.log.Error(ctx).Err(err).
			Str("template_name", req.TemplateName).
			Str("phone_number", req.PhoneNumber).
			Msg("[Template Handler] - error enviando plantilla")

		// Determinar código de error apropiado
		statusCode := http.StatusInternalServerError
		errorType := "internal_error"

		// Errores específicos del dominio
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
