package client

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/entities"
	whaerrors "github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/ports"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/infra/secondary/client/mappers"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/infra/secondary/client/response"
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
// La baseURL se recibe como parámetro (obtenida de platform credentials o .env).
// Utiliza el cliente HTTP compartido de /back/central/shared/httpclient.
func New(baseURL string, logger log.ILogger) ports.IWhatsApp {

	// Limpiar trailing slash de la baseURL para evitar doble slash en endpoints
	baseURL = strings.TrimRight(baseURL, "/")

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

		// Parsear error de Meta Graph API y retornar error del dominio
		metaErr := parseMetaGraphError(errorBody, statusCode, phoneNumberID)
		return "", metaErr
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

// metaGraphErrorResponse estructura JSON de error de la Graph API de Meta (solo para parseo)
type metaGraphErrorResponse struct {
	Error struct {
		Message      string `json:"message"`
		Type         string `json:"type"`
		Code         int    `json:"code"`
		ErrorSubcode int    `json:"error_subcode"`
	} `json:"error"`
}

// parseMetaGraphError parsea el JSON de error de Meta y retorna un error de dominio
func parseMetaGraphError(body string, statusCode int, phoneNumberID uint) *whaerrors.MetaGraphError {
	var raw metaGraphErrorResponse
	if err := json.Unmarshal([]byte(body), &raw); err != nil {
		return whaerrors.NewMetaGraphErrorUnparseable(statusCode)
	}
	return whaerrors.NewMetaGraphError(raw.Error.Code, raw.Error.ErrorSubcode, raw.Error.Message, statusCode, phoneNumberID)
}
