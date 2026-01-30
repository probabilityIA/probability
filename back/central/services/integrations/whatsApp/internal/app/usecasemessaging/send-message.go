package usecasemessaging

import (
	"context"
	"fmt"
	"strconv"

	"github.com/secamc93/probability/back/central/services/integrations/whatsApp/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/integrations/whatsApp/internal/domain/entities"
)

// SendMessage envía un mensaje de WhatsApp con el número de orden
func (u *Usecases) SendMessage(ctx context.Context, req dtos.SendMessageRequest) (string, error) {
	// Validar número de teléfono
	if err := ValidatePhoneNumber(req.PhoneNumber); err != nil {
		u.log.Error(ctx).Err(err).
			Str("phone_number", req.PhoneNumber).
			Str("order_number", req.OrderNumber).
			Msg("[WhatsApp] - número de teléfono inválido")
		return "", fmt.Errorf("número de teléfono inválido: %w", err)
	}

	// Obtener phone_number_id de variable de entorno
	phoneNumberIDStr := u.config.Get("WHATSAPP_PHONE_NUMBER_ID")
	if phoneNumberIDStr == "" {
		u.log.Error(ctx).Msg("[WhatsApp] - WHATSAPP_PHONE_NUMBER_ID no configurado")
		return "", fmt.Errorf("WHATSAPP_PHONE_NUMBER_ID no configurado")
	}

	phoneNumberID, err := strconv.ParseUint(phoneNumberIDStr, 10, 32)
	if err != nil {
		u.log.Error(ctx).Err(err).Str("phone_number_id", phoneNumberIDStr).Msg("[WhatsApp] - WHATSAPP_PHONE_NUMBER_ID inválido")
		return "", fmt.Errorf("WHATSAPP_PHONE_NUMBER_ID inválido: %w", err)
	}

	// Obtener access_token de variable de entorno (legacy - deprecado)
	accessToken := u.config.Get("WHATSAPP_TOKEN")
	if accessToken == "" {
		u.log.Error(ctx).Msg("[WhatsApp] - WHATSAPP_TOKEN no configurado")
		return "", fmt.Errorf("WHATSAPP_TOKEN no configurado")
	}

	// Construir mensaje tipo plantilla
	msg := entities.TemplateMessage{
		MessagingProduct: "whatsapp",
		RecipientType:    "individual",
		To:               req.PhoneNumber,
		Type:             "template",
		Template: entities.TemplateData{
			Name:     "order_status_9",
			Language: entities.TemplateLanguage{Code: "es"},
			Components: []entities.TemplateComponent{
				{
					Type: "body",
					Parameters: []entities.TemplateParameter{
						{
							Type:          "text",
							ParameterName: "pedido_id",
							Text:          req.OrderNumber,
						},
						{
							Type:          "text",
							ParameterName: "estado",
							Text:          "Actualizado", // Estado genérico ya que no tenemos más lógica de estados
						},
					},
				},
			},
		},
	}

	u.log.Info(ctx).
		Str("to", msg.To).
		Str("order_number", req.OrderNumber).
		Uint("phone_number_id", uint(phoneNumberID)).
		Str("template_name", msg.Template.Name).
		Msg("[WhatsApp] - enviando mensaje")

	messageID, err := u.whatsApp.SendMessage(ctx, uint(phoneNumberID), msg, accessToken)
	if err != nil {
		u.log.Error(ctx).Err(err).
			Str("phone_number", req.PhoneNumber).
			Str("order_number", req.OrderNumber).
			Uint("phone_number_id", uint(phoneNumberID)).
			Str("template_name", msg.Template.Name).
			Msg("[WhatsApp] - error enviando mensaje")
		return "", fmt.Errorf("error al enviar mensaje de WhatsApp: %w", err)
	}

	u.log.Info(ctx).
		Str("message_id", messageID).
		Str("order_number", req.OrderNumber).
		Str("phone_number", req.PhoneNumber).
		Msg("[WhatsApp] - mensaje enviado correctamente")

	return messageID, nil
}
