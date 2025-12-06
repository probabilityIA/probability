package app

import (
	"context"
	"fmt"
	"strconv"

	"github.com/secamc93/probability/back/central/services/integrations/whatsApp/internal/domain"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
)

type IUseCaseSendMessage interface {
	SendMessage(ctx context.Context, req domain.SendMessageRequest) (string, error)
}

type SendMessageUsecase struct {
	whatsApp domain.IWhatsApp
	log      log.ILogger
	config   env.IConfig
}

func New(whatsApp domain.IWhatsApp, logger log.ILogger, config env.IConfig) *SendMessageUsecase {
	return &SendMessageUsecase{
		whatsApp: whatsApp,
		log:      logger,
		config:   config,
	}
}

// SendMessage envía un mensaje de WhatsApp con el número de orden
func (u *SendMessageUsecase) SendMessage(ctx context.Context, req domain.SendMessageRequest) (string, error) {
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

	// Construir mensaje tipo plantilla
	msg := domain.TemplateMessage{
		MessagingProduct: "whatsapp",
		RecipientType:    "individual",
		To:               req.PhoneNumber,
		Type:             "template",
		Template: domain.TemplateData{
			Name:     "order_status_9",
			Language: domain.TemplateLanguage{Code: "es"},
			Components: []domain.TemplateComponent{
				{
					Type: "body",
					Parameters: []domain.TemplateParameter{
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

	messageID, err := u.whatsApp.SendMessage(ctx, uint(phoneNumberID), msg)
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
