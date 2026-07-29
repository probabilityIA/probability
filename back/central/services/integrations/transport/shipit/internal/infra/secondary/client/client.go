package client

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/transport/shipit/internal/domain"
	"github.com/secamc93/probability/back/central/shared/httpclient"
	"github.com/secamc93/probability/back/central/shared/log"
)

const (
	DefaultBaseURL         = "https://api.shipit.cl"
	DefaultTrackingBaseURL = "https://courierstatus.shipit.cl"
	AcceptHeader           = "application/vnd.shipit.v4"
)

type Client struct {
	httpClient *httpclient.Client
	log        log.ILogger
}

func New(logger log.ILogger) domain.IShipitClient {
	httpConfig := httpclient.HTTPClientConfig{
		BaseURL:    DefaultBaseURL,
		Timeout:    30 * time.Second,
		RetryCount: 2,
		RetryWait:  3 * time.Second,
		Debug:      true,
	}

	httpClient := httpclient.New(httpConfig, logger)
	httpClient.SetHeader("Accept", AcceptHeader)
	httpClient.SetHeader("Content-Type", "application/json")

	return &Client{
		httpClient: httpClient,
		log:        logger.WithModule("shipit.client"),
	}
}

func resolveBaseURL(baseURL string) string {
	if strings.TrimSpace(baseURL) == "" {
		return DefaultBaseURL
	}
	return strings.TrimRight(baseURL, "/")
}

func resolveTrackingBaseURL(baseURL string) string {
	resolved := resolveBaseURL(baseURL)
	if strings.Contains(resolved, "api.shipit.cl") {
		return DefaultTrackingBaseURL
	}
	return resolved
}

func authHeaders(creds domain.Credentials) map[string]string {
	return map[string]string{
		"X-Shipit-Email":        creds.Email,
		"X-Shipit-Access-Token": creds.AccessToken,
		"Accept":                AcceptHeader,
		"Content-Type":          "application/json",
	}
}

type shipitErrorResponse struct {
	State   string `json:"state"`
	Message string `json:"message"`
	Error   string `json:"error"`
	Errors  any    `json:"errors"`
}

func parseShipitError(statusCode int, body []byte) string {
	if statusCode == 401 || statusCode == 403 {
		return "Credenciales de Shipit invalidas: verifica el email y el access token"
	}

	var errResp shipitErrorResponse
	if err := json.Unmarshal(body, &errResp); err == nil {
		if errResp.Message != "" {
			return "Error de Shipit: " + errResp.Message
		}
		if errResp.Error != "" {
			return "Error de Shipit: " + errResp.Error
		}
	}

	preview := string(body)
	if len(preview) > 300 {
		preview = preview[:300]
	}
	return "Error de Shipit: " + preview
}

func (c *Client) TestAuthentication(baseURL string, creds domain.Credentials) error {
	ctx := context.Background()
	url := resolveBaseURL(baseURL) + "/v/webhooks/package"

	resp, err := c.httpClient.R().
		SetContext(ctx).
		SetHeaders(authHeaders(creds)).
		Get(url)
	if err != nil {
		return err
	}

	if resp.StatusCode() == 401 || resp.StatusCode() == 403 {
		return domain.ErrInvalidCredentials
	}
	if resp.StatusCode() >= 500 {
		return domain.ErrEmptyResponse
	}

	c.log.Info(ctx).Int("status", resp.StatusCode()).Msg("Shipit authentication verified")
	return nil
}
