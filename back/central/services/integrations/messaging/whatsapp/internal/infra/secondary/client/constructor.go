package client

import (
	"context"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/ports"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/infra/secondary/client/mappers"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/infra/secondary/client/response"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/httpclient"
	"github.com/secamc93/probability/back/central/shared/log"
)

// implementa ports.IWhatsApp
type whatsAppClient struct {
	httpClient *httpclient.Client
	baseURL    string
	logger     log.ILogger
}

// New construye y devuelve un cliente que implementa ports.IWhatsApp.
// La baseURL se toma de las variables de entorno WHATSAPP_URL.
// Utiliza el cliente HTTP compartido de /back/central/shared/httpclient.
func New(config env.IConfig, logger log.ILogger) ports.IWhatsApp {
	baseURL := config.Get("WHATSAPP_URL")

	// Configurar el cliente HTTP compartido
	httpClient := httpclient.New(httpclient.HTTPClientConfig{
		Timeout:    30 * time.Second, // Timeout más largo para WhatsApp API
		BaseURL:    baseURL,
		RetryCount: 2,
		RetryWait:  5 * time.Second,
		Debug:      false,
	}, logger.WithModule("whatsapp-client"))

	// Configurar headers por defecto
	httpClient.SetHeader("Content-Type", "application/json")

	return &whatsAppClient{
		httpClient: httpClient,
		baseURL:    baseURL,
		logger:     logger.WithModule("whatsapp-client"),
	}
}

// SendMessage envía un mensaje (texto o template) vía WhatsApp Cloud API.
// Requiere:
// - phoneNumberID: ID del número de teléfono para construir la URL.
// - msg: DTO de dominio con los datos del mensaje.
// - accessToken: Token de acceso para autenticación con WhatsApp API.
// La URL se construye como: {baseURL}/{phone_number_id}/messages
func (c *whatsAppClient) SendMessage(ctx context.Context, phoneNumberID uint, msg entities.TemplateMessage, accessToken string) (string, error) {
	c.logger.Info().
		Uint("phone_number_id", phoneNumberID).
		Str("to", msg.To).
		Str("type", msg.Type).
		Msg("Sending WhatsApp message")

	payload := mappers.MapDomainToRequest(msg)

	// Construir la URL específica para este phone_number_id
	endpoint := fmt.Sprintf("%d/messages", phoneNumberID)

	var result response.SendMessageResponse
	var errorResponse map[string]interface{}

	resp, err := c.httpClient.R().
		SetContext(ctx).
		SetHeader("Authorization", "Bearer "+accessToken).
		SetBody(payload).
		SetResult(&result).
		SetError(&errorResponse).
		Post(endpoint)

	// Verificar errores de red o HTTP
	if err != nil {
		c.logger.Error().
			Err(err).
			Uint("phone_number_id", phoneNumberID).
			Msg("Network error communicating with WhatsApp API")
		return "", fmt.Errorf("error al comunicarse con WhatsApp API: %w", err)
	}

	// Verificar código de estado HTTP explícitamente
	statusCode := resp.StatusCode()
	if statusCode < 200 || statusCode >= 300 {
		errorBody := resp.String()
		c.logger.Error().
			Int("status_code", statusCode).
			Str("response_body", errorBody).
			Uint("phone_number_id", phoneNumberID).
			Msg("WhatsApp API returned error status")

		// Intentar extraer mensaje de error de la respuesta
		errorMsg := fmt.Sprintf("WhatsApp API retornó error %d", statusCode)
		if errorBody != "" {
			errorMsg = fmt.Sprintf("WhatsApp API error %d: %s", statusCode, errorBody)
		}

		return "", fmt.Errorf("%s", errorMsg)
	}

	// Verificar que la respuesta contiene mensajes
	if len(result.Messages) == 0 {
		errorBody := resp.String()
		c.logger.Error().
			Str("response_body", errorBody).
			Uint("phone_number_id", phoneNumberID).
			Msg("Successful response but no messages in response")
		return "", fmt.Errorf("la respuesta de WhatsApp no contiene mensajes: %s", errorBody)
	}

	messageID := result.Messages[0].ID
	c.logger.Info().
		Str("message_id", messageID).
		Uint("phone_number_id", phoneNumberID).
		Str("to", msg.To).
		Msg("WhatsApp message sent successfully")
	return messageID, nil
}
