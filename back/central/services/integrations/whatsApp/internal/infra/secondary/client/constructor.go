package client

import (
	"context"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/secamc93/probability/back/central/services/integrations/whatsApp/internal/domain"
	"github.com/secamc93/probability/back/central/services/integrations/whatsApp/internal/infra/secondary/client/mappers"
	"github.com/secamc93/probability/back/central/services/integrations/whatsApp/internal/infra/secondary/client/response"
	httpclient "github.com/secamc93/probability/back/central/shared/client"
	"github.com/secamc93/probability/back/central/shared/env"
)

// implementa domain.IWhatsApp
type whatsAppClient struct {
	rest        *resty.Client
	accessToken string
}

// New construye y devuelve un cliente que implementa domain.IWhatsApp.
// La baseURL y el token se toman de las variables de entorno WHATSAPP_URL y WHATSAPP_TOKEN.
// Utiliza el cliente HTTP compartido para reutilizar la configuración común.
func New(config env.IConfig) domain.IWhatsApp {
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
			if r.StatusCode() >= 400 {
				return fmt.Errorf("whatsapp_request_failed: status_code=%d, body=%s", r.StatusCode(), r.String())
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
// La URL se construye como: {baseURL}{phone_number_id}/messages
// El token se obtiene de la variable de entorno WHATSAPP_TOKEN.
func (c *whatsAppClient) SendMessage(ctx context.Context, phoneNumberID uint, msg domain.TemplateMessage) (string, error) {
	payload := mappers.MapDomainToRequest(msg)

	// Construir la URL específica para este phone_number_id
	endpoint := fmt.Sprintf("%d/messages", phoneNumberID)

	var result response.SendMessageResponse
	resp, err := c.rest.R().
		SetContext(ctx).
		SetHeader("Authorization", "Bearer "+c.accessToken).
		SetBody(payload).
		SetResult(&result).
		SetDebug(true).
		Post(endpoint)
	if err != nil {
		return "", err
	}

	if len(result.Messages) > 0 {
		return result.Messages[0].ID, nil
	}

	return resp.Status(), nil
}
