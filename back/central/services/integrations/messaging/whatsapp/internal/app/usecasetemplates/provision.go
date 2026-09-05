package usecasetemplates

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/ports"
)

func (u *usecase) Provision(ctx context.Context, businessID uint) (*ProvisionResult, error) {
	target, err := u.credentialsCache.GetWhatsAppConfig(ctx, businessID)
	if err != nil {
		return nil, err
	}

	if !target.OwnNumber {
		return nil, fmt.Errorf("el negocio %d usa el número de la plataforma: no hay plantillas propias que aprovisionar", businessID)
	}
	if target.WABAID == "" {
		return nil, fmt.Errorf("la integración del negocio %d no tiene waba_id configurado", businessID)
	}

	platform, err := u.credentialsCache.GetWhatsAppDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("credenciales de plataforma no disponibles: %w", err)
	}
	if platform.WABAID == "" {
		return nil, fmt.Errorf("las credenciales de plataforma no tienen waba_id: no se puede leer el catálogo de plantillas origen")
	}

	sourceAPI := u.apiFactory(pickBaseURL(platform.WhatsAppURL, target.WhatsAppURL))
	targetAPI := u.apiFactory(pickBaseURL(target.WhatsAppURL, platform.WhatsAppURL))

	sourceTemplates, err := sourceAPI.ListTemplates(ctx, platform.WABAID, platform.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("error leyendo las plantillas de la plataforma: %w", err)
	}

	existing, err := targetAPI.ListTemplates(ctx, target.WABAID, target.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("error leyendo las plantillas del negocio: %w", err)
	}

	present := make(map[string]ports.TemplateDefinitionRemote, len(existing))
	for _, tpl := range existing {
		present[templateKey(tpl.Name, tpl.Language)] = tpl
	}

	result := &ProvisionResult{
		IntegrationID: target.IntegrationID,
		BusinessID:    businessID,
		WABAID:        target.WABAID,
		Created:       []string{},
		AlreadyExists: []string{},
		Failed:        map[string]string{},
	}

	statuses := make([]ports.TemplateStatus, 0, len(sourceTemplates))

	for _, tpl := range sourceTemplates {
		key := templateKey(tpl.Name, tpl.Language)

		if found, ok := present[key]; ok {
			result.AlreadyExists = append(result.AlreadyExists, tpl.Name)
			statuses = append(statuses, ports.TemplateStatus{
				Name:        found.Name,
				Language:    found.Language,
				Status:      found.Status,
				Category:    found.Category,
				MetaID:      found.ID,
				Reason:      found.RejectedReason,
				UpdatedAt:   time.Now(),
				Provisioned: true,
			})
			continue
		}

		metaID, createErr := targetAPI.CreateTemplate(ctx, target.WABAID, target.AccessToken, tpl)
		if createErr != nil {
			u.log.Warn(ctx).Err(createErr).
				Str("template", tpl.Name).
				Str("waba_id", target.WABAID).
				Msg("no se pudo crear la plantilla en el WABA del negocio")
			result.Failed[tpl.Name] = createErr.Error()
			statuses = append(statuses, ports.TemplateStatus{
				Name:        tpl.Name,
				Language:    tpl.Language,
				Status:      "ERROR",
				Category:    tpl.Category,
				Reason:      createErr.Error(),
				UpdatedAt:   time.Now(),
				Provisioned: false,
			})
			continue
		}

		result.Created = append(result.Created, tpl.Name)
		statuses = append(statuses, ports.TemplateStatus{
			Name:        tpl.Name,
			Language:    tpl.Language,
			Status:      "PENDING",
			Category:    tpl.Category,
			MetaID:      metaID,
			UpdatedAt:   time.Now(),
			Provisioned: true,
		})
	}

	snapshot := &ports.WABATemplatesSnapshot{
		IntegrationID: target.IntegrationID,
		BusinessID:    businessID,
		WABAID:        target.WABAID,
		Templates:     statuses,
		RefreshedAt:   time.Now(),
	}
	if err := u.templatesCache.Save(ctx, snapshot); err != nil {
		u.log.Warn(ctx).Err(err).Msg("no se pudo guardar el estado de plantillas en cache")
	}

	result.Templates = statuses

	u.log.Info(ctx).
		Uint("business_id", businessID).
		Str("waba_id", target.WABAID).
		Int("creadas", len(result.Created)).
		Int("existentes", len(result.AlreadyExists)).
		Int("fallidas", len(result.Failed)).
		Msg("aprovisionamiento de plantillas terminado")

	return result, nil
}

func templateKey(name, language string) string {
	return strings.ToLower(name) + "|" + strings.ToLower(language)
}

func pickBaseURL(primary, fallback string) string {
	if strings.TrimSpace(primary) != "" {
		return primary
	}
	return fallback
}
