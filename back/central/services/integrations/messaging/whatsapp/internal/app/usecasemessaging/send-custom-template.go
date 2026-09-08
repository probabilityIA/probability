package usecasemessaging

import (
	"context"
	"fmt"
	"strings"

	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/entities"
)

func (u *usecases) SendCustomTemplate(
	ctx context.Context,
	templateName string,
	language string,
	phoneNumber string,
	parameters []string,
	businessID uint,
) (string, error) {
	templateName = strings.TrimSpace(templateName)
	if templateName == "" {
		return "", fmt.Errorf("el nombre de la plantilla es obligatorio")
	}

	if strings.TrimSpace(language) == "" {
		language = "es"
	}

	phoneNumber = NormalizePhoneNumber(phoneNumber)
	if err := ValidatePhoneNumber(phoneNumber); err != nil {
		return "", fmt.Errorf("numero de telefono invalido: %w", err)
	}

	whatsappConfig, err := u.credentialsCache.GetWhatsAppConfig(ctx, businessID)
	if err != nil {
		return "", fmt.Errorf("error obteniendo configuracion de WhatsApp: %w", err)
	}

	components := []entities.TemplateComponent{}
	if len(parameters) > 0 {
		bodyParams := make([]entities.TemplateParameter, 0, len(parameters))
		for _, value := range parameters {
			bodyParams = append(bodyParams, entities.TemplateParameter{
				Type: "text",
				Text: value,
			})
		}
		components = append(components, entities.TemplateComponent{
			Type:       "body",
			Parameters: bodyParams,
		})
	}

	msg := entities.TemplateMessage{
		MessagingProduct: "whatsapp",
		RecipientType:    "individual",
		To:               phoneNumber,
		Type:             "template",
		Template: entities.TemplateData{
			Name:       templateName,
			Language:   entities.TemplateLanguage{Code: language},
			Components: components,
		},
	}

	waClient := u.whatsApp
	if whatsappConfig.WhatsAppURL != "" && u.clientFactory != nil {
		waClient = u.clientFactory(whatsappConfig.WhatsAppURL)
	}

	messageID, err := waClient.SendMessage(ctx, whatsappConfig.PhoneNumberID, msg, whatsappConfig.AccessToken)
	if err != nil {
		u.log.Error(ctx).Err(err).
			Str("template_name", templateName).
			Str("phone_number", phoneNumber).
			Uint("business_id", businessID).
			Msg("[WhatsApp UseCase] - error enviando plantilla propia del negocio")
		return "", err
	}

	u.log.Info(ctx).
		Str("template_name", templateName).
		Str("phone_number", phoneNumber).
		Str("message_id", messageID).
		Msg("[WhatsApp UseCase] - plantilla propia del negocio enviada")

	return messageID, nil
}
