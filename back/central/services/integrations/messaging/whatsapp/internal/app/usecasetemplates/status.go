package usecasetemplates

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/ports"
)

func (u *usecase) GetStatus(ctx context.Context, businessID uint, refresh bool) (*ports.WABATemplatesSnapshot, error) {
	config, err := u.credentialsCache.GetWhatsAppConfig(ctx, businessID)
	if err != nil {
		return nil, err
	}

	if !config.OwnNumber || config.WABAID == "" {
		return &ports.WABATemplatesSnapshot{
			IntegrationID: config.IntegrationID,
			BusinessID:    businessID,
			Templates:     []ports.TemplateStatus{},
			RefreshedAt:   time.Now(),
		}, nil
	}

	if !refresh {
		if snapshot, cacheErr := u.templatesCache.Get(ctx, config.IntegrationID); cacheErr == nil && snapshot != nil {
			return snapshot, nil
		}
	}

	api := u.apiFactory(config.WhatsAppURL)

	remote, err := api.ListTemplates(ctx, config.WABAID, config.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("error consultando las plantillas del negocio: %w", err)
	}

	statuses := make([]ports.TemplateStatus, 0, len(remote))
	for _, tpl := range remote {
		statuses = append(statuses, ports.TemplateStatus{
			Name:        tpl.Name,
			Language:    tpl.Language,
			Status:      tpl.Status,
			Category:    tpl.Category,
			MetaID:      tpl.ID,
			Reason:      tpl.RejectedReason,
			UpdatedAt:   time.Now(),
			Provisioned: true,
		})
	}

	snapshot := &ports.WABATemplatesSnapshot{
		IntegrationID: config.IntegrationID,
		BusinessID:    businessID,
		WABAID:        config.WABAID,
		Templates:     statuses,
		RefreshedAt:   time.Now(),
	}

	if err := u.templatesCache.Save(ctx, snapshot); err != nil {
		u.log.Warn(ctx).Err(err).Msg("no se pudo guardar el estado de plantillas en cache")
	}

	return snapshot, nil
}

func (u *usecase) HandleStatusUpdate(ctx context.Context, wabaID, name, language, event, reason string) error {
	wabaID = strings.TrimSpace(wabaID)
	name = strings.TrimSpace(name)

	if wabaID == "" || name == "" {
		return nil
	}

	status := normalizeTemplateEvent(event)

	if err := u.templatesCache.UpdateStatusByWABA(ctx, wabaID, name, language, status, reason); err != nil {
		u.log.Warn(ctx).Err(err).
			Str("waba_id", wabaID).
			Str("template", name).
			Msg("no se pudo actualizar el estado de la plantilla en cache")
		return nil
	}

	u.log.Info(ctx).
		Str("waba_id", wabaID).
		Str("template", name).
		Str("status", status).
		Msg("estado de plantilla actualizado desde el webhook")

	return nil
}

func normalizeTemplateEvent(event string) string {
	switch strings.ToUpper(strings.TrimSpace(event)) {
	case "APPROVED":
		return "APPROVED"
	case "REJECTED":
		return "REJECTED"
	case "PENDING", "PENDING_DELETION":
		return "PENDING"
	case "DISABLED", "PAUSED":
		return "DISABLED"
	case "":
		return "UNKNOWN"
	default:
		return strings.ToUpper(strings.TrimSpace(event))
	}
}
