package client

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/httpclient"
	"github.com/secamc93/probability/back/central/shared/log"
)

type phoneNumbersClient struct {
	httpClient *httpclient.Client
	logger     log.ILogger
}

func NewPhoneNumbersClient(baseURL string, logger log.ILogger) ports.IPhoneNumbersAPI {
	baseURL = strings.TrimRight(baseURL, "/")

	client := httpclient.New(httpclient.HTTPClientConfig{
		Timeout:    30 * time.Second,
		BaseURL:    baseURL,
		RetryCount: 1,
		RetryWait:  3 * time.Second,
		Debug:      false,
	}, logger.WithModule("whatsapp-phone-numbers-client"))

	client.SetHeader("Content-Type", "application/json")

	return &phoneNumbersClient{
		httpClient: client,
		logger:     logger.WithModule("whatsapp-phone-numbers-client"),
	}
}

type addPhoneNumberResponse struct {
	ID string `json:"id"`
}

type phoneNumberDetailResponse struct {
	ID                 string `json:"id"`
	DisplayPhoneNumber string `json:"display_phone_number"`
	VerifiedName       string `json:"verified_name"`
	QualityRating      string `json:"quality_rating"`
	CodeVerification   string `json:"code_verification_status"`
	NameStatus         string `json:"name_status"`
	Status             string `json:"status"`
}

func (c *phoneNumbersClient) AddPhoneNumber(ctx context.Context, wabaID, accessToken, countryCode, phoneNumber, verifiedName string) (string, error) {
	if wabaID == "" {
		return "", fmt.Errorf("waba_id no configurado")
	}

	var result addPhoneNumberResponse

	resp, err := c.httpClient.R().
		SetContext(ctx).
		SetHeader("Authorization", "Bearer "+accessToken).
		SetBody(map[string]any{
			"cc":            countryCode,
			"phone_number":  phoneNumber,
			"verified_name": verifiedName,
		}).
		SetResult(&result).
		Post(fmt.Sprintf("%s/phone_numbers", wabaID))
	if err != nil {
		return "", fmt.Errorf("error agregando el número al WABA %s: %w", wabaID, err)
	}

	if resp.StatusCode() < 200 || resp.StatusCode() >= 300 {
		return "", parseMetaGraphError(resp.String(), resp.StatusCode(), 0)
	}

	if result.ID == "" {
		return "", fmt.Errorf("Meta no devolvió el phone_number_id del número agregado")
	}

	return result.ID, nil
}

func (c *phoneNumbersClient) RequestCode(ctx context.Context, phoneNumberID, accessToken, method, language string) error {
	if method == "" {
		method = "SMS"
	}
	if language == "" {
		language = "es"
	}

	resp, err := c.httpClient.R().
		SetContext(ctx).
		SetHeader("Authorization", "Bearer "+accessToken).
		SetBody(map[string]any{
			"code_method": strings.ToUpper(method),
			"language":    language,
		}).
		Post(fmt.Sprintf("%s/request_code", phoneNumberID))
	if err != nil {
		return fmt.Errorf("error pidiendo el código de verificación: %w", err)
	}

	if resp.StatusCode() < 200 || resp.StatusCode() >= 300 {
		return parseMetaGraphError(resp.String(), resp.StatusCode(), 0)
	}

	return nil
}

func (c *phoneNumbersClient) VerifyCode(ctx context.Context, phoneNumberID, accessToken, code string) error {
	resp, err := c.httpClient.R().
		SetContext(ctx).
		SetHeader("Authorization", "Bearer "+accessToken).
		SetBody(map[string]any{"code": code}).
		Post(fmt.Sprintf("%s/verify_code", phoneNumberID))
	if err != nil {
		return fmt.Errorf("error verificando el código: %w", err)
	}

	if resp.StatusCode() < 200 || resp.StatusCode() >= 300 {
		return parseMetaGraphError(resp.String(), resp.StatusCode(), 0)
	}

	return nil
}

func (c *phoneNumbersClient) Register(ctx context.Context, phoneNumberID, accessToken, pin string) error {
	resp, err := c.httpClient.R().
		SetContext(ctx).
		SetHeader("Authorization", "Bearer "+accessToken).
		SetBody(map[string]any{
			"messaging_product": "whatsapp",
			"pin":               pin,
		}).
		Post(fmt.Sprintf("%s/register", phoneNumberID))
	if err != nil {
		return fmt.Errorf("error registrando el número en la Cloud API: %w", err)
	}

	if resp.StatusCode() < 200 || resp.StatusCode() >= 300 {
		return parseMetaGraphError(resp.String(), resp.StatusCode(), 0)
	}

	return nil
}

func (c *phoneNumbersClient) GetPhoneNumber(ctx context.Context, phoneNumberID, accessToken string) (*ports.WABAPhoneNumber, error) {
	var result phoneNumberDetailResponse

	resp, err := c.httpClient.R().
		SetContext(ctx).
		SetHeader("Authorization", "Bearer "+accessToken).
		SetQueryParam("fields", "id,display_phone_number,verified_name,quality_rating,code_verification_status,name_status,status").
		SetResult(&result).
		Get(phoneNumberID)
	if err != nil {
		return nil, fmt.Errorf("error consultando el número %s: %w", phoneNumberID, err)
	}

	if resp.StatusCode() < 200 || resp.StatusCode() >= 300 {
		return nil, parseMetaGraphError(resp.String(), resp.StatusCode(), 0)
	}

	return &ports.WABAPhoneNumber{
		ID:                     result.ID,
		DisplayPhoneNumber:     result.DisplayPhoneNumber,
		VerifiedName:           result.VerifiedName,
		QualityRating:          result.QualityRating,
		CodeVerificationStatus: result.CodeVerification,
		NameStatus:             result.NameStatus,
		Status:                 result.Status,
	}, nil
}

func (c *phoneNumbersClient) DeletePhoneNumber(ctx context.Context, phoneNumberID, accessToken string) error {
	resp, err := c.httpClient.R().
		SetContext(ctx).
		SetHeader("Authorization", "Bearer "+accessToken).
		Delete(phoneNumberID)
	if err != nil {
		return fmt.Errorf("error borrando el número %s: %w", phoneNumberID, err)
	}

	if resp.StatusCode() < 200 || resp.StatusCode() >= 300 {
		return parseMetaGraphError(resp.String(), resp.StatusCode(), 0)
	}

	return nil
}

func (c *phoneNumbersClient) DeregisterPhoneNumber(ctx context.Context, phoneNumberID, accessToken string) error {
	resp, err := c.httpClient.R().
		SetContext(ctx).
		SetHeader("Authorization", "Bearer "+accessToken).
		Post(fmt.Sprintf("%s/deregister", phoneNumberID))
	if err != nil {
		return fmt.Errorf("error dando de baja el número %s de la Cloud API: %w", phoneNumberID, err)
	}

	if resp.StatusCode() < 200 || resp.StatusCode() >= 300 {
		return parseMetaGraphError(resp.String(), resp.StatusCode(), 0)
	}

	return nil
}
