package client

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	bcErrors "github.com/secamc93/probability/back/central/services/integrations/pay/bancolombia/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/integrations/pay/bancolombia/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

const (
	defaultSandboxBaseURL   = "https://sb-api.bancolombia.com"
	defaultTokenPath        = "/security/oauth-provider/oauth2/token"
	defaultQRGeneratePath   = "/v1/qr-management/codes"
	defaultScope            = "QR-Management:write:app"
	bancolombiaMediaType    = "application/vnd.bancolombia.v4+json"
	tokenExpirySafetyMargin = 60 * time.Second
)

type cachedToken struct {
	accessToken string
	expiresAt   time.Time
}

type Client struct {
	log    log.ILogger
	mu     sync.Mutex
	tokens map[string]cachedToken
}

func New(logger log.ILogger) ports.IBancolombiaClient {
	return &Client{
		log:    logger.WithModule("bancolombia.client"),
		tokens: make(map[string]cachedToken),
	}
}

func resolveBaseURL(config *ports.BancolombiaConfig) string {
	if config.BaseURL != "" {
		return strings.TrimRight(config.BaseURL, "/")
	}
	return defaultSandboxBaseURL
}

func resolvePath(configured, fallback string) string {
	if strings.TrimSpace(configured) != "" {
		return "/" + strings.TrimLeft(strings.TrimSpace(configured), "/")
	}
	return fallback
}

func (c *Client) tokenCacheKey(config *ports.BancolombiaConfig) string {
	return config.Environment + "|" + config.ClientID
}

func (c *Client) getToken(ctx context.Context, config *ports.BancolombiaConfig) (string, error) {
	key := c.tokenCacheKey(config)

	c.mu.Lock()
	if cached, ok := c.tokens[key]; ok && time.Now().Before(cached.expiresAt) {
		token := cached.accessToken
		c.mu.Unlock()
		return token, nil
	}
	c.mu.Unlock()

	scope := config.Scope
	if scope == "" {
		scope = defaultScope
	}

	basic := base64.StdEncoding.EncodeToString([]byte(config.ClientID + ":" + config.ClientSecret))

	var result struct {
		TokenType   string `json:"token_type"`
		AccessToken string `json:"access_token"`
		Scope       string `json:"scope"`
		ExpiresIn   int    `json:"expires_in"`
	}

	tokenURL := resolveBaseURL(config) + resolvePath(config.TokenPath, defaultTokenPath)
	resp, err := resty.New().R().
		SetContext(ctx).
		SetHeader("Authorization", "Basic "+basic).
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetHeader("Accept", bancolombiaMediaType).
		SetFormData(map[string]string{
			"grant_type": "client_credentials",
			"scope":      scope,
		}).
		SetResult(&result).
		Post(tokenURL)

	if err != nil {
		c.log.Error(ctx).Err(err).Str("token_url", tokenURL).Msg("Bancolombia token request failed")
		return "", fmt.Errorf("%w: %v", bcErrors.ErrAuthFailed, err)
	}
	if resp.IsError() {
		c.log.Error(ctx).Str("status", resp.Status()).Str("body", resp.String()).Msg("Bancolombia token request rejected")
		return "", fmt.Errorf("%w: status=%s body=%s", bcErrors.ErrAuthFailed, resp.Status(), resp.String())
	}
	if result.AccessToken == "" {
		return "", fmt.Errorf("%w: empty access_token", bcErrors.ErrAuthFailed)
	}

	expiresIn := result.ExpiresIn
	if expiresIn <= 0 {
		expiresIn = 1200
	}
	expiry := time.Now().Add(time.Duration(expiresIn)*time.Second - tokenExpirySafetyMargin)

	c.mu.Lock()
	c.tokens[key] = cachedToken{accessToken: result.AccessToken, expiresAt: expiry}
	c.mu.Unlock()

	c.log.Info(ctx).
		Str("environment", config.Environment).
		Int("expires_in", expiresIn).
		Msg("Bancolombia OAuth token obtained")

	return result.AccessToken, nil
}

func (c *Client) invalidateToken(config *ports.BancolombiaConfig) {
	c.mu.Lock()
	delete(c.tokens, c.tokenCacheKey(config))
	c.mu.Unlock()
}

func (c *Client) GenerateQR(ctx context.Context, config *ports.BancolombiaConfig, amount float64, currency, reference, description string) (*ports.QRGenerationResult, error) {
	if currency == "" {
		currency = "COP"
	}

	reqBody := map[string]interface{}{
		"data": []map[string]interface{}{
			{
				"merchantId":   config.MerchantID,
				"merchantName": config.MerchantName,
				"reference":    reference,
				"amount":       amount,
				"currency":     currency,
				"description":  description,
			},
		},
	}

	result, resp, err := c.postQR(ctx, config, reqBody)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() == http.StatusUnauthorized {
		c.invalidateToken(config)
		result, resp, err = c.postQR(ctx, config, reqBody)
		if err != nil {
			return nil, err
		}
	}
	if resp.IsError() {
		c.log.Error(ctx).Str("status", resp.Status()).Str("body", resp.String()).Msg("Bancolombia QR API returned error")
		return nil, fmt.Errorf("%w: status=%s body=%s", bcErrors.ErrAPIError, resp.Status(), resp.String())
	}

	payload := extractFirstDataObject(result)
	qrValue := pickString(payload, "qrValue", "qr_value", "emvCode", "emv_code", "codeValue", "qrCode", "code", "value")
	qrImage := pickString(payload, "qrImage", "qr_image", "image", "qrBase64", "imageBase64", "base64Image")
	externalID := pickString(payload, "codeId", "qrId", "id", "transactionId", "transferCode")
	if externalID == "" {
		externalID = reference
	}

	if qrValue == "" && qrImage == "" {
		c.log.Error(ctx).Interface("response", result).Msg("Bancolombia QR response without recognizable QR payload")
		return nil, fmt.Errorf("%w: response without qr value", bcErrors.ErrAPIError)
	}

	return &ports.QRGenerationResult{
		ExternalID:    externalID,
		QRValue:       qrValue,
		QRImageBase64: qrImage,
		RawResponse:   result,
	}, nil
}

func (c *Client) postQR(ctx context.Context, config *ports.BancolombiaConfig, reqBody map[string]interface{}) (map[string]interface{}, *resty.Response, error) {
	token, err := c.getToken(ctx, config)
	if err != nil {
		return nil, nil, err
	}

	qrURL := resolveBaseURL(config) + resolvePath(config.QRGeneratePath, defaultQRGeneratePath)
	result := map[string]interface{}{}
	resp, err := resty.New().R().
		SetContext(ctx).
		SetHeader("Authorization", "Bearer "+token).
		SetHeader("Content-Type", bancolombiaMediaType).
		SetHeader("Accept", bancolombiaMediaType).
		SetHeader("message-id", uuid.New().String()).
		SetBody(reqBody).
		SetResult(&result).
		Post(qrURL)

	if err != nil {
		c.log.Error(ctx).Err(err).Str("qr_url", qrURL).Msg("Error calling Bancolombia QR API")
		return nil, nil, fmt.Errorf("bancolombia api call failed: %w", err)
	}
	return result, resp, nil
}

func extractFirstDataObject(result map[string]interface{}) map[string]interface{} {
	if result == nil {
		return map[string]interface{}{}
	}
	switch data := result["data"].(type) {
	case []interface{}:
		if len(data) > 0 {
			if obj, ok := data[0].(map[string]interface{}); ok {
				return obj
			}
		}
	case map[string]interface{}:
		return data
	}
	return result
}

func pickString(m map[string]interface{}, keys ...string) string {
	for _, k := range keys {
		if v, ok := m[k].(string); ok && strings.TrimSpace(v) != "" {
			return strings.TrimSpace(v)
		}
	}
	return ""
}
