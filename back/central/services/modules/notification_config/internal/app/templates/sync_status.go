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

	names := make([]string, 0, len(pending))
	for _, template := range pending {
		if template.Origin != entities.TemplateOriginBusiness || template.BusinessID == nil || *template.BusinessID != businessID {
			continue
		}
		if template.MetaTemplateID == "" {
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
