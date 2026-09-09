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

type businessProfileClient struct {
	httpClient *httpclient.Client
	logger     log.ILogger
}

func NewBusinessProfileClient(baseURL string, logger log.ILogger) ports.IBusinessProfileAPI {
	baseURL = strings.TrimRight(baseURL, "/")

	client := httpclient.New(httpclient.HTTPClientConfig{
		Timeout:    30 * time.Second,
		BaseURL:    baseURL,
		RetryCount: 1,
		RetryWait:  3 * time.Second,
		Debug:      false,
	}, logger.WithModule("whatsapp-business-profile-client"))

	client.SetHeader("Content-Type", "application/json")

	return &businessProfileClient{
		httpClient: client,
		logger:     logger.WithModule("whatsapp-business-profile-client"),
	}
}

type businessProfileResponse struct {
	Data []struct {
		About             string   `json:"about"`
		Address           string   `json:"address"`
		Description       string   `json:"description"`
		Email             string   `json:"email"`
		ProfilePictureURL string   `json:"profile_picture_url"`
		Websites          []string `json:"websites"`
		Vertical          string   `json:"vertical"`
	} `json:"data"`
}

const profileFields = "about,address,description,email,profile_picture_url,websites,vertical"

func (c *businessProfileClient) GetProfile(ctx context.Context, phoneNumberID, accessToken string) (*ports.WhatsAppBusinessProfile, error) {
	if phoneNumberID == "" {
		return nil, fmt.Errorf("phone_number_id no configurado")
	}

	var result businessProfileResponse

	resp, err := c.httpClient.R().
		SetContext(ctx).
		SetHeader("Authorization", "Bearer "+accessToken).
		SetQueryParam("fields", profileFields).
		SetResult(&result).
		Get(fmt.Sprintf("%s/whatsapp_business_profile", phoneNumberID))
	if err != nil {
		return nil, fmt.Errorf("error consultando el perfil de empresa: %w", err)
	}

	if resp.StatusCode() < 200 || resp.StatusCode() >= 300 {
		return nil, parseMetaGraphError(resp.String(), resp.StatusCode(), 0)
	}

	perfil := &ports.WhatsAppBusinessProfile{}
	if len(result.Data) > 0 {
		d := result.Data[0]
		perfil.About = d.About
		perfil.Address = d.Address
		perfil.Description = d.Description
		perfil.Email = d.Email
		perfil.ProfilePictureURL = d.ProfilePictureURL
		perfil.Websites = d.Websites
		perfil.Vertical = d.Vertical
	}

	return perfil, nil
}

func (c *businessProfileClient) UpdateProfile(ctx context.Context, phoneNumberID, accessToken string, perfil ports.WhatsAppBusinessProfileUpdate) error {
	if phoneNumberID == "" {
		return fmt.Errorf("phone_number_id no configurado")
	}

	body := map[string]any{"messaging_product": "whatsapp"}

	if perfil.About != nil && *perfil.About != "" {
		body["about"] = *perfil.About
	}
	if perfil.Address != nil && *perfil.Address != "" {
		body["address"] = *perfil.Address
	}
	if perfil.Description != nil && *perfil.Description != "" {
		body["description"] = *perfil.Description
	}
	if perfil.Email != nil && *perfil.Email != "" {
		body["email"] = *perfil.Email
	}
	if perfil.Vertical != nil && *perfil.Vertical != "" && *perfil.Vertical != "UNDEFINED" {
		body["vertical"] = *perfil.Vertical
	}
	if perfil.Websites != nil {
		body["websites"] = *perfil.Websites
	}
	if perfil.ProfilePictureHandle != nil {
		body["profile_picture_handle"] = *perfil.ProfilePictureHandle
	}

	if len(body) == 1 {
		return nil
	}

	resp, err := c.httpClient.R().
		SetContext(ctx).
		SetHeader("Authorization", "Bearer "+accessToken).
		SetBody(body).
		Post(fmt.Sprintf("%s/whatsapp_business_profile", phoneNumberID))
	if err != nil {
		return fmt.Errorf("error actualizando el perfil de empresa: %w", err)
	}

	if resp.StatusCode() < 200 || resp.StatusCode() >= 300 {
		return parseMetaGraphError(resp.String(), resp.StatusCode(), 0)
	}

	return nil
}

type uploadSessionResponse struct {
	ID string `json:"id"`
}

type uploadFinishResponse struct {
	Handle string `json:"h"`
}

func (c *businessProfileClient) UploadProfilePhoto(ctx context.Context, appID, accessToken, contentType string, data []byte) (string, error) {
	if appID == "" {
		return "", fmt.Errorf("app_id no configurado: sin el no se puede subir la foto a Meta")
	}
	if len(data) == 0 {
		return "", fmt.Errorf("la imagen esta vacia")
	}

	var sesion uploadSessionResponse

	resp, err := c.httpClient.R().
		SetContext(ctx).
		SetQueryParams(map[string]string{
			"file_length":  fmt.Sprintf("%d", len(data)),
			"file_type":    contentType,
			"access_token": accessToken,
		}).
		SetResult(&sesion).
		Post(fmt.Sprintf("%s/uploads", appID))
	if err != nil {
		return "", fmt.Errorf("error abriendo la sesion de subida: %w", err)
	}
	if resp.StatusCode() < 200 || resp.StatusCode() >= 300 {
		return "", parseMetaGraphError(resp.String(), resp.StatusCode(), 0)
	}
	if sesion.ID == "" {
		return "", fmt.Errorf("Meta no devolvio el id de la sesion de subida")
	}

	var subida uploadFinishResponse

	resp, err = c.httpClient.R().
		SetContext(ctx).
		SetHeader("Authorization", "OAuth "+accessToken).
		SetHeader("file_offset", "0").
		SetHeader("Content-Type", contentType).
		SetBody(data).
		SetResult(&subida).
		Post(sesion.ID)
	if err != nil {
		return "", fmt.Errorf("error subiendo la imagen: %w", err)
	}
	if resp.StatusCode() < 200 || resp.StatusCode() >= 300 {
		return "", parseMetaGraphError(resp.String(), resp.StatusCode(), 0)
	}
	if subida.Handle == "" {
		return "", fmt.Errorf("Meta no devolvio el identificador de la imagen subida")
	}

	return subida.Handle, nil
}
