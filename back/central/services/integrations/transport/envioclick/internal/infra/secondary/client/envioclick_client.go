package client

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/transport/envioclick/internal/domain"
	"github.com/secamc93/probability/back/central/shared/httpclient"
	"github.com/secamc93/probability/back/central/shared/log"
)

const (
	DefaultBaseURL = "https://api.envioclickpro.com.co/api/v2"
)

// Client implements domain.IEnvioClickClient using the shared httpclient
type Client struct {
	httpClient *httpclient.Client
	log        log.ILogger
}

// New creates a new EnvioClick HTTP client.
// The baseURL is no longer fixed at construction — each method receives it dynamically.
func New(logger log.ILogger) domain.IEnvioClickClient {
	logger.Info(context.Background()).Msg("🔍 Creating EnvioClick HTTP client")

	httpConfig := httpclient.HTTPClientConfig{
		BaseURL:    DefaultBaseURL, // placeholder; overridden per-request
		Timeout:    30 * time.Second,
		RetryCount: 2,
		RetryWait:  3 * time.Second,
		Debug:      true,
	}

	httpClient := httpclient.New(httpConfig, logger)
	httpClient.SetHeader("Accept", "application/json")
	httpClient.SetHeader("Content-Type", "application/json")

	return &Client{
		httpClient: httpClient,
		log:        logger.WithModule("envioclick.client"),
	}
}

// envioClickErrorResponse represents the structured error format from EnvioClick
type envioClickErrorResponse struct {
	StatusMessages []struct {
		Error []string `json:"error"`
	} `json:"status_messages"`
}

// parseEnvioClickError extracts a human-readable error from an EnvioClick error response body
func parseEnvioClickError(body []byte) string {
	var errorResp envioClickErrorResponse
	if err := json.Unmarshal(body, &errorResp); err == nil && len(errorResp.StatusMessages) > 0 {
		for _, msg := range errorResp.StatusMessages {
			if len(msg.Error) > 0 {
				fullError := strings.Join(msg.Error, " ")
				return mapEnvioClickError(fullError)
			}
		}
	}
	return mapEnvioClickError(string(body))
}

// mapEnvioClickError translates EnvioClick error messages to user-friendly Spanish
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
	if strings.Contains(lowerErr, "missing") || strings.Contains(lowerErr, "falta") || strings.Contains(lowerErr, "unprocessed entity") {
		return "Error de validación: Faltan datos obligatorios o hay datos inválidos en la solicitud"
	}
	if strings.Contains(lowerErr, "phone") || strings.Contains(lowerErr, "telefóno") {
		return "error: el número de teléfono es inválido o tiene formato incorrecto"
	}

	return originalErr
}
