package templates

import (
	"context"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/entities"
)

const (
	maxTemplatesToSync = 100
	pendingSyncGrace   = 5 * time.Minute
)

func (uc *useCase) SyncAllPendingStatuses(ctx context.Context) (int, error) {
	businessIDs, err := uc.repository.ListBusinessesWithPendingTemplates(ctx, time.Now().Add(-pendingSyncGrace))
	if err != nil {
		return 0, err
	}

	synced := 0
	for _, businessID := range businessIDs {
		if _, err := uc.SyncPendingStatuses(ctx, businessID); err != nil {
			uc.logger.Warn().Err(err).Uint("business_id", businessID).
				Msg("No se pudo pedir el estado de las plantillas pendientes")
			continue
		}
		synced++
	}

	return synced, nil
}

func (uc *useCase) SyncPendingStatuses(ctx context.Context, businessID uint) (int, error) {
	if uc.publisher == nil {
		return 0, fmt.Errorf("no hay conexion con Meta para consultar el estado de las plantillas")
	}

	pending, _, err := uc.repository.ListTemplates(ctx, businessID, "", entities.TemplateStatusPending, 1, maxTemplatesToSync)
	if err != nil {
		return 0, err
	}

	orphanCutoff := time.Now().Add(-pendingSyncGrace)

	names := make([]string, 0, len(pending))
	for _, template := range pending {
		if template.Origin != entities.TemplateOriginBusiness || template.BusinessID == nil || *template.BusinessID != businessID {
			continue
		}
		if template.MetaTemplateID == "" {
			uc.markOrphanedSubmissionAsFailed(ctx, template, orphanCutoff)
			continue
		}
		names = append(names, template.Name)
	}

	if len(names) == 0 {
		return 0, nil
	}

	message := dtos.TemplateSubmissionMessage{
		Action:     "sync",
		BusinessID: businessID,
		Name:       "sync",
		Names:      names,
	}

	if err := uc.publisher.PublishTemplateSubmission(ctx, message); err != nil {
		return 0, fmt.Errorf("no se pudo pedir el estado a Meta: %w", err)
	}

	uc.logger.Info().Uint("business_id", businessID).Int("plantillas", len(names)).
		Msg("Consulta de estado de plantillas pendientes encolada")

	return len(names), nil
}

func (uc *useCase) markOrphanedSubmissionAsFailed(ctx context.Context, template entities.WhatsappTemplate, cutoff time.Time) {
	if template.SubmittedAt == nil || template.SubmittedAt.After(cutoff) {
		return
	}

	template.Status = entities.TemplateStatusFailed
	template.RejectedReason = "el envio a Meta se perdio antes de recibir respuesta (nunca se genero meta_template_id); reintentar"
	now := time.Now()
	template.LastSyncedAt = &now

	if err := uc.repository.UpdateTemplate(ctx, &template); err != nil {
		uc.logger.Warn().Err(err).Uint("template_id", template.ID).
			Msg("No se pudo marcar como fallida una plantilla huerfana")
		return
	}

	uc.logger.Warn().Uint("template_id", template.ID).Str("name", template.Name).
		Msg("Plantilla marcada como fallida: se quedo en pending sin respuesta de Meta")
}
