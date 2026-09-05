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

type templatesClient struct {
	httpClient *httpclient.Client
	logger     log.ILogger
}

func NewTemplatesClient(baseURL string, logger log.ILogger) ports.ITemplateAPI {
	baseURL = strings.TrimRight(baseURL, "/")

	client := httpclient.New(httpclient.HTTPClientConfig{
		Timeout:    30 * time.Second,
		BaseURL:    baseURL,
		RetryCount: 2,
		RetryWait:  5 * time.Second,
		Debug:      false,
	}, logger.WithModule("whatsapp-templates-client"))

	client.SetHeader("Content-Type", "application/json")

	return &templatesClient{
		httpClient: client,
		logger:     logger.WithModule("whatsapp-templates-client"),
	}
}

type templateListResponse struct {
	Data []struct {
		ID              string           `json:"id"`
		Name            string           `json:"name"`
		Language        string           `json:"language"`
		Category        string           `json:"category"`
		Status          string           `json:"status"`
		RejectedReason  string           `json:"rejected_reason"`
		ParameterFormat string           `json:"parameter_format"`
		Components      []map[string]any `json:"components"`
	} `json:"data"`
	Paging struct {
		Next string `json:"next"`
	} `json:"paging"`
}

type templateCreateResponse struct {
	ID       string `json:"id"`
	Status   string `json:"status"`
	Category string `json:"category"`
}

func (c *templatesClient) ListTemplates(ctx context.Context, wabaID, accessToken string) ([]ports.TemplateDefinitionRemote, error) {
	if wabaID == "" {
		return nil, fmt.Errorf("waba_id no configurado")
	}

	var result templateListResponse

	resp, err := c.httpClient.R().
		SetContext(ctx).
		SetHeader("Authorization", "Bearer "+accessToken).
		SetQueryParam("limit", "250").
		SetQueryParam("fields", "id,name,language,category,status,rejected_reason,parameter_format,components").
		SetResult(&result).
		Get(fmt.Sprintf("%s/message_templates", wabaID))
	if err != nil {
		return nil, fmt.Errorf("error consultando plantillas del WABA %s: %w", wabaID, err)
	}

	if resp.StatusCode() < 200 || resp.StatusCode() >= 300 {
		return nil, parseMetaGraphError(resp.String(), resp.StatusCode(), 0)
	}

	templates := make([]ports.TemplateDefinitionRemote, 0, len(result.Data))
	for _, item := range result.Data {
		templates = append(templates, ports.TemplateDefinitionRemote{
			ID:              item.ID,
			Name:            item.Name,
			Language:        item.Language,
			Category:        item.Category,
			Status:          item.Status,
			RejectedReason:  item.RejectedReason,
			ParameterFormat: item.ParameterFormat,
			Components:      item.Components,
		})
	}

	return templates, nil
}

func (c *templatesClient) CreateTemplate(ctx context.Context, wabaID, accessToken string, template ports.TemplateDefinitionRemote) (string, error) {
	if wabaID == "" {
		return "", fmt.Errorf("waba_id no configurado")
	}
	if template.Name == "" || template.Language == "" {
		return "", fmt.Errorf("la plantilla necesita nombre e idioma")
	}

	payload := map[string]any{
		"name":       template.Name,
		"language":   template.Language,
		"category":   template.Category,
		"components": template.Components,
	}
	if template.ParameterFormat != "" {
		payload["parameter_format"] = template.ParameterFormat
	}

	var result templateCreateResponse

	resp, err := c.httpClient.R().
		SetContext(ctx).
		SetHeader("Authorization", "Bearer "+accessToken).
		SetBody(payload).
		SetResult(&result).
		Post(fmt.Sprintf("%s/message_templates", wabaID))
	if err != nil {
		return "", fmt.Errorf("error creando plantilla %s en el WABA %s: %w", template.Name, wabaID, err)
	}

	if resp.StatusCode() < 200 || resp.StatusCode() >= 300 {
		return "", parseMetaGraphError(resp.String(), resp.StatusCode(), 0)
	}

	return result.ID, nil
}
