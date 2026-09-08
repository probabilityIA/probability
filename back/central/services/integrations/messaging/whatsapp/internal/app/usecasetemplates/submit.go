package usecasetemplates

import (
	"context"
	"fmt"
	"strings"

	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/ports"
)

type CustomTemplateSubmission struct {
	Action         string           `json:"action"`
	TemplateID     uint             `json:"template_id"`
	BusinessID     uint             `json:"business_id"`
	Name           string           `json:"name"`
	Language       string           `json:"language"`
	Category       string           `json:"category"`
	MetaTemplateID string           `json:"meta_template_id"`
	Components     []map[string]any `json:"components"`
}

type CustomTemplateResult struct {
	TemplateID     uint   `json:"template_id"`
	BusinessID     uint   `json:"business_id"`
	MetaTemplateID string `json:"meta_template_id"`
	WABAID         string `json:"waba_id"`
	Name           string `json:"name"`
	Language       string `json:"language"`
	Status         string `json:"status"`
	Reason         string `json:"reason"`
	ErrorMessage   string `json:"error_message"`
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
