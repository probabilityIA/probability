package envioclick

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
)

const (
	BaseURL = "https://api.envioclickpro.com.co/api/v2"
	// TODO: Move to config
	APIKey = "80f036b8-6ced-4982-aa2b-b32f1acbb7ab" // Using example key for now as requested/implied, USER MUST PROVIDE REAL KEY
)

type Client struct {
	logger log.ILogger
	client *http.Client
}

func New(logger log.ILogger) *Client {
	return &Client{
		logger: logger,
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

func (c *Client) Quote(req domain.EnvioClickQuoteRequest) (*domain.EnvioClickQuoteResponse, error) {
	url := fmt.Sprintf("%s/quotation", BaseURL)
	return c.doRequest(url, req)
}

func (c *Client) Generate(req domain.EnvioClickQuoteRequest) (*domain.EnvioClickGenerateResponse, error) {
	// Generate guide usually takes the idRate selected
	url := fmt.Sprintf("%s/shipment", BaseURL)

	// The request body for generation typically includes idRate and sender/recipient details if strictly needed,
	// but based on doc, usually it's about confirming the rate.
	// For now we assume req is the correct payload structure expected by EnvioClick for generation.

	// NOTE: The user provided documentation for Quotation, but cut off for Generation.
	// We will assume a generic structure or that the user will provide the body payload in the service.
	// We'll return a generic response structure.

	// Adapting for generic response parsing
	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	respBytes, err := c.doRawRequest("POST", url, bodyBytes)
	if err != nil {
		return nil, err
	}

	var resp domain.EnvioClickGenerateResponse
	if err := json.Unmarshal(respBytes, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c *Client) doRequest(url string, payload interface{}) (*domain.EnvioClickQuoteResponse, error) {
	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	respBytes, err := c.doRawRequest("POST", url, bodyBytes)
	if err != nil {
		return nil, err
	}

	var resp domain.EnvioClickQuoteResponse
	if err := json.Unmarshal(respBytes, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

type EnvioClickErrorResponse struct {
	StatusMessages []struct {
		Error []string `json:"error"`
	} `json:"status_messages"`
}

func (c *Client) doRawRequest(method, url string, body []byte) ([]byte, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", APIKey)
	req.Header.Set("Content-Type", "application/json")

	fmt.Printf("DEBUG ENVIOCLICK BODY: %s\n", string(body)) // Temporary log

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		c.logger.Error().Str("body", string(respBody)).Int("status", resp.StatusCode).Msg("EnvioClick API Error")

		var errorResp EnvioClickErrorResponse
		if err := json.Unmarshal(respBody, &errorResp); err == nil && len(errorResp.StatusMessages) > 0 {
			for _, msg := range errorResp.StatusMessages {
				if len(msg.Error) > 0 {
					// Join all error messages found
					fullError := ""
					for _, e := range msg.Error {
						fullError += e + " "
					}
					return nil, fmt.Errorf("%s", mapEnvioClickError(fullError))
				}
			}
		}

		return nil, fmt.Errorf("Error de EnvioClick: %s", mapEnvioClickError(string(respBody)))
	}

	return respBody, nil
}

func mapEnvioClickError(originalErr string) string {
	lowerErr := strings.ToLower(originalErr)

	if (strings.Contains(lowerErr, "destination") || strings.Contains(lowerErr, "destino")) && strings.Contains(lowerErr, "dane") {
		return "error: el codigo dane del destino no existe o no es valido"
	}
	if (strings.Contains(lowerErr, "origin") || strings.Contains(lowerErr, "origen")) && strings.Contains(lowerErr, "dane") {
		return "error: el codigo dane de origen no existe o no es valido"
	}
	if strings.Contains(lowerErr, "contentvalue") || strings.Contains(lowerErr, "declared value") {
		return "El valor declarado es inválido o está fuera de rango"
	}
	if strings.Contains(lowerErr, "weight") || strings.Contains(lowerErr, "peso") {
		return "El peso del paquete es inválido"
	}
	if strings.Contains(lowerErr, "dimensions") || strings.Contains(lowerErr, "height") || strings.Contains(lowerErr, "width") || strings.Contains(lowerErr, "length") {
		return "Las dimensiones del paquete son inválidas"
	}
	if strings.Contains(lowerErr, "missing") || strings.Contains(lowerErr, "falta") {
		return "Faltan datos obligatorios para generar la guía"
	}
	if strings.Contains(lowerErr, "phone") || strings.Contains(lowerErr, "telefóno") {
		return "El número de teléfono es inválido"
	}

	return originalErr
}
