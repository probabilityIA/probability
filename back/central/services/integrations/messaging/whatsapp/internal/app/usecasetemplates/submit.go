package usecasetemplates

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/ports"
)

const maxHeaderImageBytes = 5 * 1024 * 1024

func (u *usecase) resolveHeaderHandle(ctx context.Context, submission CustomTemplateSubmission, baseURL, token string) (string, error) {
	imageURL := strings.TrimSpace(submission.HeaderMediaURL)
	if imageURL == "" {
		return "", nil
	}

	platform, err := u.credentialsCache.GetWhatsAppDefaultConfig(ctx)
	if err != nil || platform == nil {
		return "", fmt.Errorf("no se pudieron leer las credenciales de la plataforma para subir la imagen")
	}
	if platform.AppID == "" {
		return "", fmt.Errorf("la plataforma no tiene app_id configurado: sin el Meta no acepta la imagen")
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, imageURL, nil)
	if err != nil {
		return "", fmt.Errorf("url de imagen invalida: %w", err)
	}

	httpClient := &http.Client{Timeout: 30 * time.Second}
	response, err := httpClient.Do(request)
	if err != nil {
		return "", fmt.Errorf("no se pudo descargar la imagen del encabezado: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return "", fmt.Errorf("la imagen del encabezado respondio %d", response.StatusCode)
	}

	data, err := io.ReadAll(io.LimitReader(response.Body, maxHeaderImageBytes+1))
	if err != nil {
		return "", fmt.Errorf("no se pudo leer la imagen del encabezado: %w", err)
	}
	if len(data) > maxHeaderImageBytes {
		return "", fmt.Errorf("la imagen del encabezado supera los 5 MB que acepta Meta")
	}

	contentType := response.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "image/jpeg"
	}

	return u.apiFactory(baseURL).UploadMedia(ctx, platform.AppID, token, contentType, data)
}

func applyHeaderHandle(components []map[string]any, handle string) {
	if handle == "" {
		return
	}
	for _, component := range components {
		componentType, _ := component["type"].(string)
		format, _ := component["format"].(string)
		if strings.EqualFold(componentType, "HEADER") && strings.EqualFold(format, "IMAGE") {
			component["example"] = map[string]any{"header_handle": []string{handle}}
		}
	}
}

type CustomTemplateSubmission struct {
	Action         string           `json:"action"`
	TemplateID     uint             `json:"template_id"`
	BusinessID     uint             `json:"business_id"`
	Name           string           `json:"name"`
	Language       string           `json:"language"`
	Category       string           `json:"category"`
	MetaTemplateID string           `json:"meta_template_id"`
	HeaderMediaURL string           `json:"header_media_url"`
	Components     []map[string]any `json:"components"`
}

type CustomTemplateResult struct {
	TemplateID        uint   `json:"template_id"`
	BusinessID        uint   `json:"business_id"`
	MetaTemplateID    string `json:"meta_template_id"`
	HeaderMediaHandle string `json:"header_media_handle"`
	WABAID            string `json:"waba_id"`
	Name              string `json:"name"`
	Language          string `json:"language"`
	Status            string `json:"status"`
	Reason            string `json:"reason"`
	ErrorMessage      string `json:"error_message"`
}

func (u *usecase) SubmitCustom(ctx context.Context, submission CustomTemplateSubmission) CustomTemplateResult {
	result := CustomTemplateResult{
		TemplateID: submission.TemplateID,
		BusinessID: submission.BusinessID,
	}

	wabaID, token, baseURL, err := u.resolveSubmissionTarget(ctx, submission.BusinessID)
	if err != nil {
		result.ErrorMessage = err.Error()
		return result
	}

	result.WABAID = wabaID

	if len(submission.Components) == 0 {
		result.ErrorMessage = "la plantilla llego sin componentes (cuerpo, botones): no se envia a Meta"
		return result
	}

	handle, err := u.resolveHeaderHandle(ctx, submission, baseURL, token)
	if err != nil {
		result.ErrorMessage = err.Error()
		return result
	}
	applyHeaderHandle(submission.Components, handle)
	result.HeaderMediaHandle = handle

	definition := ports.TemplateDefinitionRemote{
		Name:       strings.TrimSpace(submission.Name),
		Language:   strings.TrimSpace(submission.Language),
		Category:   strings.ToUpper(strings.TrimSpace(submission.Category)),
		Components: submission.Components,
	}

	metaID, err := u.apiFactory(baseURL).CreateTemplate(ctx, wabaID, token, definition)
	if err != nil {
		u.log.Error(ctx).Err(err).
			Uint("template_id", submission.TemplateID).
			Str("name", submission.Name).
			Str("waba_id", wabaID).
			Msg("Meta rechazo la creacion de la plantilla")
		result.ErrorMessage = err.Error()
		return result
	}

	result.MetaTemplateID = metaID
	result.Name = definition.Name
	result.Language = definition.Language
	result.Status = "PENDING"

	u.log.Info(ctx).
		Uint("template_id", submission.TemplateID).
		Str("name", submission.Name).
		Str("meta_template_id", metaID).
		Str("waba_id", wabaID).
		Msg("Plantilla creada en Meta, queda a la espera de revision")

	return result
}

func (u *usecase) UpdateCustom(ctx context.Context, submission CustomTemplateSubmission) CustomTemplateResult {
	result := CustomTemplateResult{
		TemplateID:     submission.TemplateID,
		BusinessID:     submission.BusinessID,
		MetaTemplateID: submission.MetaTemplateID,
	}

	if strings.TrimSpace(submission.MetaTemplateID) == "" {
		result.ErrorMessage = "la plantilla no existe todavia en Meta: no se puede editar"
		return result
	}

	wabaID, token, baseURL, err := u.resolveSubmissionTarget(ctx, submission.BusinessID)
	if err != nil {
		result.ErrorMessage = err.Error()
		return result
	}

	result.WABAID = wabaID

	handle, err := u.resolveHeaderHandle(ctx, submission, baseURL, token)
	if err != nil {
		result.ErrorMessage = err.Error()
		return result
	}
	applyHeaderHandle(submission.Components, handle)
	if handle != "" {
		result.HeaderMediaHandle = handle
	}

	definition := ports.TemplateDefinitionRemote{
		Name:       strings.TrimSpace(submission.Name),
		Language:   strings.TrimSpace(submission.Language),
		Category:   strings.ToUpper(strings.TrimSpace(submission.Category)),
		Components: submission.Components,
	}

	if err := u.apiFactory(baseURL).UpdateTemplate(ctx, token, submission.MetaTemplateID, definition); err != nil {
		u.log.Error(ctx).Err(err).
			Uint("template_id", submission.TemplateID).
			Str("name", submission.Name).
			Str("meta_template_id", submission.MetaTemplateID).
			Msg("Meta rechazo la edicion de la plantilla")
		result.ErrorMessage = err.Error()
		return result
	}

	result.Name = definition.Name
	result.Language = definition.Language
	result.Status = "PENDING"

	u.log.Info(ctx).
		Uint("template_id", submission.TemplateID).
		Str("name", submission.Name).
		Str("meta_template_id", submission.MetaTemplateID).
		Msg("Plantilla editada en Meta, vuelve a quedar en revision")

	return result
}

func (u *usecase) resolveSubmissionTarget(ctx context.Context, businessID uint) (string, string, string, error) {
	platform, platformErr := u.credentialsCache.GetWhatsAppDefaultConfig(ctx)

	config, err := u.credentialsCache.GetWhatsAppConfig(ctx, businessID)
	if err == nil && config != nil && config.OwnNumber && config.WABAID != "" {
		baseURL := config.WhatsAppURL
		if platformErr == nil && platform != nil {
			baseURL = pickBaseURL(config.WhatsAppURL, platform.WhatsAppURL)
		}
		return config.WABAID, config.AccessToken, baseURL, nil
	}

	if platformErr != nil || platform == nil {
		return "", "", "", fmt.Errorf("no hay credenciales de plataforma para crear la plantilla")
	}
	if platform.WABAID == "" {
		return "", "", "", fmt.Errorf("las credenciales de plataforma no tienen waba_id")
	}

	return platform.WABAID, platform.AccessToken, platform.WhatsAppURL, nil
}

func (u *usecase) DeleteCustom(ctx context.Context, submission CustomTemplateSubmission) error {
	wabaID, token, baseURL, err := u.resolveSubmissionTarget(ctx, submission.BusinessID)
	if err != nil {
		return err
	}

	if err := u.apiFactory(baseURL).DeleteTemplate(ctx, wabaID, token, submission.Name, submission.MetaTemplateID); err != nil {
		u.log.Error(ctx).Err(err).
			Str("name", submission.Name).
			Str("waba_id", wabaID).
			Msg("error borrando la plantilla en Meta")
		return err
	}

	u.log.Info(ctx).
		Str("name", submission.Name).
		Str("waba_id", wabaID).
		Msg("plantilla borrada en Meta")

	return nil
}
