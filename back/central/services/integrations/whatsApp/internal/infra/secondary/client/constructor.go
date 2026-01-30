package client

import (
	"context"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/secamc93/probability/back/central/services/integrations/whatsApp/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/integrations/whatsApp/internal/domain/ports"
	"github.com/secamc93/probability/back/central/services/integrations/whatsApp/internal/infra/secondary/client/mappers"
	"github.com/secamc93/probability/back/central/services/integrations/whatsApp/internal/infra/secondary/client/response"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/httpclient"
)

// implementa ports.IWhatsApp
type whatsAppClient struct {
	rest        *resty.Client
	accessToken string
}

// New construye y devuelve un cliente que implementa ports.IWhatsApp.
// La baseURL y el token se toman de las variables de entorno WHATSAPP_URL y WHATSAPP_TOKEN.
// Utiliza el cliente HTTP compartido para reutilizar la configuración común.
func New(config env.IConfig) ports.IWhatsApp {
	baseURL := config.Get("WHATSAPP_URL")
	accessToken := config.Get("WHATSAPP_TOKEN")

	// Usar el cliente HTTP compartido
	httpClient := httpclient.NewHTTPClient(httpclient.HTTPClientConfig{
		Timeout:         30 * time.Second, // Timeout más largo para WhatsApp API
		MaxIdleConns:    100,
		IdleConnTimeout: 90 * time.Second,
	})

	// Configurar resty con el cliente HTTP compartido
	rest := resty.NewWithClient(httpClient).
		SetBaseURL(baseURL).
		AddRetryCondition(func(r *resty.Response, err error) bool {
			return r.StatusCode() == 429
		}).
		SetRetryCount(2).
		SetRetryWaitTime(5 * time.Second).
		OnAfterResponse(func(c *resty.Client, r *resty.Response) error {
			// Solo loguear, no devolver error aquí para que SendMessage pueda manejar el error
			// con más contexto (status code, body, etc.)
			if r.StatusCode() >= 400 {
				fmt.Printf("⚠️ [WhatsAppClient] OnAfterResponse detectó error HTTP %d: %s\n", r.StatusCode(), r.String())
			}
			return nil
		})

	// Headers por defecto
	rest.SetHeader("Content-Type", "application/json")

	return &whatsAppClient{
		rest:        rest,
		accessToken: accessToken,
	}
}

// SendMessage envía un mensaje (texto o template) vía WhatsApp Cloud API.
// Requiere:
// - phoneNumberID: ID del número de teléfono para construir la URL.
// - msg: DTO de dominio con los datos del mensaje.
// - accessToken: Token de acceso para autenticación con WhatsApp API.
// La URL se construye como: {baseURL}{phone_number_id}/messages
func (c *whatsAppClient) SendMessage(ctx context.Context, phoneNumberID uint, msg entities.TemplateMessage, accessToken string) (string, error) {
	fmt.Printf("🚀 [WhatsAppClient] SendMessage called for PhoneID: %d\n", phoneNumberID)
	payload := mappers.MapDomainToRequest(msg)

	// Construir la URL específica para este phone_number_id
	endpoint := fmt.Sprintf("%d/messages", phoneNumberID)

	var result response.SendMessageResponse
	var errorResponse map[string]interface{}

	resp, err := c.rest.R().
		SetContext(ctx).
		SetHeader("Authorization", "Bearer "+accessToken).
		SetBody(payload).
		SetResult(&result).
		SetError(&errorResponse).
		SetDebug(true).
		Post(endpoint)

	// Verificar errores de red o HTTP
	if err != nil {
		fmt.Printf("❌ [WhatsAppClient] Error de red/HTTP: %v\n", err)
		return "", fmt.Errorf("error al comunicarse con WhatsApp API: %w", err)
	}

	// Verificar código de estado HTTP explícitamente
	statusCode := resp.StatusCode()
	if statusCode < 200 || statusCode >= 300 {
		errorBody := resp.String()
		fmt.Printf("❌ [WhatsAppClient] Error HTTP %d: %s\n", statusCode, errorBody)

		// Intentar extraer mensaje de error de la respuesta
		errorMsg := fmt.Sprintf("WhatsApp API retornó error %d", statusCode)
		if errorBody != "" {
			errorMsg = fmt.Sprintf("WhatsApp API error %d: %s", statusCode, errorBody)
		}

		return "", fmt.Errorf(errorMsg)
	}

	// Verificar que la respuesta contiene mensajes
	if len(result.Messages) == 0 {
		errorBody := resp.String()
		fmt.Printf("❌ [WhatsAppClient] Respuesta exitosa pero sin mensajes: %s\n", errorBody)
		return "", fmt.Errorf("la respuesta de WhatsApp no contiene mensajes: %s", errorBody)
	}

	messageID := result.Messages[0].ID
	fmt.Printf("✅ [WhatsAppClient] Mensaje enviado exitosamente. MessageID: %s\n", messageID)
	return messageID, nil
}
