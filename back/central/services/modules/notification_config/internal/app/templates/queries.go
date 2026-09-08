package templates

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/entities"
)

func (uc *useCase) GetByID(ctx context.Context, id, businessID uint) (*entities.WhatsappTemplate, error) {
	template, err := uc.repository.GetTemplateByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if template == nil {
		return nil, nil
	}
	if !template.BelongsTo(businessID) {
		return nil, fmt.Errorf("la plantilla no pertenece al negocio")
	}
	if template.Origin == entities.TemplateOriginInternal {
		return nil, fmt.Errorf("plantilla de uso interno")
	}
	return template, nil
}

func (uc *useCase) List(ctx context.Context, businessID uint, scope, status string, page, pageSize int) ([]entities.WhatsappTemplate, int64, error) {
	scope = strings.TrimSpace(scope)
	if !entities.IsAllowedScope(scope) {
		scope = entities.TemplateScopeScheduled
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return uc.repository.ListTemplates(ctx, businessID, scope, strings.TrimSpace(status), page, pageSize)
}

func (uc *useCase) Update(ctx context.Context, dto dtos.UpdateTemplateDTO) (*entities.WhatsappTemplate, error) {
	current, err := uc.GetByID(ctx, dto.ID, dto.BusinessID)
	if err != nil {
		return nil, err
	}
	if current == nil {
		return nil, fmt.Errorf("plantilla no encontrada")
	}

	if current.Status == entities.TemplateStatusPending {
		return nil, fmt.Errorf("la plantilla esta en revision de Meta: no se puede editar hasta que responda")
	}

	if !current.IsEditable() {
		return nil, fmt.Errorf("las plantillas predeterminadas del sistema no se pueden editar")
	}

	createDTO := dtos.CreateTemplateDTO{
		BusinessID: dto.BusinessID,
		Name:       current.Name,
		Language:   current.Language,
		Category:   dto.Category,
		HeaderText: dto.HeaderText,
		BodyText:   dto.BodyText,
		FooterText: dto.FooterText,
		Variables:  dto.Variables,
		Buttons:    dto.Buttons,
	}

	rebuilt, err := buildTemplate(createDTO)
	if err != nil {
		return nil, err
	}

	rebuilt.ID = current.ID
	rebuilt.CreatedByID = current.CreatedByID
	rebuilt.WABAID = current.WABAID
	rebuilt.MetaTemplateID = ""
	rebuilt.Status = entities.TemplateStatusDraft
	rebuilt.RejectedReason = ""

	if err := uc.repository.UpdateTemplate(ctx, rebuilt); err != nil {
		return nil, err
	}

	if err := uc.submit(ctx, rebuilt); err != nil {
		return rebuilt, err
	}

	return rebuilt, nil
}

func (uc *useCase) Delete(ctx context.Context, id, businessID uint) error {
	template, err := uc.GetByID(ctx, id, businessID)
	if err != nil {
		return err
	}
	if template == nil {
		return fmt.Errorf("plantilla no encontrada")
	}
	if !template.IsEditable() {
		return fmt.Errorf("las plantillas predeterminadas del sistema no se pueden eliminar")
	}

	if uc.publisher != nil && template.MetaTemplateID != "" {
		message := dtos.TemplateSubmissionMessage{
			Action:         "delete",
			TemplateID:     template.ID,
			BusinessID:     businessID,
			Name:           template.Name,
			Language:       template.Language,
			MetaTemplateID: template.MetaTemplateID,
		}
		if err := uc.publisher.PublishTemplateSubmission(ctx, message); err != nil {
			return fmt.Errorf("no se pudo encolar el borrado en Meta: %w", err)
		}
	}

	return uc.repository.DeleteTemplate(ctx, id)
}

func (uc *useCase) Resubmit(ctx context.Context, id, businessID uint) (*entities.WhatsappTemplate, error) {
	template, err := uc.GetByID(ctx, id, businessID)
	if err != nil {
		return nil, err
	}
	if template == nil {
		return nil, fmt.Errorf("plantilla no encontrada")
	}
	if template.Status == entities.TemplateStatusPending {
		return nil, fmt.Errorf("la plantilla ya esta en revision de Meta")
	}
	if template.Status == entities.TemplateStatusApproved {
		return nil, fmt.Errorf("la plantilla ya esta aprobada")
	}

	if err := uc.submit(ctx, template); err != nil {
		return template, err
	}

	return template, nil
}

func (uc *useCase) ApplySubmissionResult(ctx context.Context, result dtos.TemplateSubmissionResult) error {
	template, err := uc.repository.GetTemplateByID(ctx, result.TemplateID)
	if err != nil {
		return err
	}
	if template == nil {
		uc.logger.Warn().Uint("template_id", result.TemplateID).
			Msg("Resultado de envio para una plantilla que ya no existe")
		return nil
	}

	template.MetaTemplateID = result.MetaTemplateID
	template.WABAID = result.WABAID

	if result.ErrorMessage != "" {
		template.Status = entities.TemplateStatusFailed
		template.RejectedReason = result.ErrorMessage
	} else {
		template.Status = normalizeMetaEvent(result.Status)
		if template.Status == "" {
			template.Status = entities.TemplateStatusPending
		}
		template.RejectedReason = ""
	}

	now := time.Now()
	template.LastSyncedAt = &now

	return uc.repository.UpdateTemplate(ctx, template)
}

func (uc *useCase) ApplyMetaStatus(ctx context.Context, wabaID, name, language, event, reason string) error {
	status := normalizeMetaEvent(event)
	if status == "" {
		return nil
	}
	return uc.repository.UpdateTemplateStatusByMeta(ctx, wabaID, name, language, status, reason)
}

func (uc *useCase) VariableCatalog() map[string]string {
	return entities.AllowedVariableSources()
}

func normalizeMetaEvent(event string) string {
	switch strings.ToUpper(strings.TrimSpace(event)) {
	case "APPROVED":
		return entities.TemplateStatusApproved
	case "REJECTED":
		return entities.TemplateStatusRejected
	case "PENDING", "PENDING_DELETION", "IN_APPEAL":
		return entities.TemplateStatusPending
	case "PAUSED":
		return entities.TemplateStatusPaused
	case "DISABLED":
		return entities.TemplateStatusDisabled
	default:
		return ""
	}
}
