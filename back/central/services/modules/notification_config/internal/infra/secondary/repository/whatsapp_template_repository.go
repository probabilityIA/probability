package repository

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/entities"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type whatsappTemplateRepository struct {
	db     db.IDatabase
	logger log.ILogger
}

func (r *whatsappTemplateRepository) CreateTemplate(ctx context.Context, template *entities.WhatsappTemplate) error {
	model, err := templateToModel(template)
	if err != nil {
		return err
	}

	if err := r.db.Conn(ctx).Create(model).Error; err != nil {
		r.logger.Error().Err(err).Str("name", template.Name).Msg("Error creating whatsapp template")
		return err
	}

	template.ID = model.ID
	template.CreatedAt = model.CreatedAt
	template.UpdatedAt = model.UpdatedAt
	return nil
}

func (r *whatsappTemplateRepository) UpdateTemplate(ctx context.Context, template *entities.WhatsappTemplate) error {
	model, err := templateToModel(template)
	if err != nil {
		return err
	}
	model.ID = template.ID

	if err := r.db.Conn(ctx).Model(&models.WhatsappTemplate{}).
		Where("id = ? AND business_id = ?", template.ID, template.BusinessID).
		Updates(map[string]any{
			"name":             model.Name,
			"language":         model.Language,
			"category":         model.Category,
			"body_text":        model.BodyText,
			"header_text":      model.HeaderText,
			"footer_text":      model.FooterText,
			"variable_mapping": model.VariableMapping,
			"buttons":          model.Buttons,
			"components":       model.Components,
			"waba_id":          model.WABAID,
			"meta_template_id": model.MetaTemplateID,
			"status":           model.Status,
			"rejected_reason":  model.RejectedReason,
			"submitted_at":     model.SubmittedAt,
			"reviewed_at":      model.ReviewedAt,
			"last_synced_at":   model.LastSyncedAt,
		}).Error; err != nil {
		r.logger.Error().Err(err).Uint("id", template.ID).Msg("Error updating whatsapp template")
		return err
	}

	return nil
}

func (r *whatsappTemplateRepository) GetTemplateByID(ctx context.Context, id uint) (*entities.WhatsappTemplate, error) {
	var model models.WhatsappTemplate

	if err := r.db.Conn(ctx).First(&model, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		r.logger.Error().Err(err).Uint("id", id).Msg("Error getting whatsapp template by id")
		return nil, err
	}

	return templateToDomain(&model)
}

func (r *whatsappTemplateRepository) GetTemplateByName(ctx context.Context, businessID uint, name, language string) (*entities.WhatsappTemplate, error) {
	var model models.WhatsappTemplate

	query := r.db.Conn(ctx).Where("name = ? AND (business_id = ? OR business_id IS NULL)", strings.TrimSpace(name), businessID)
	if language != "" {
		query = query.Where("language = ?", language)
	}

	if err := query.First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		r.logger.Error().Err(err).Str("name", name).Msg("Error getting whatsapp template by name")
		return nil, err
	}

	return templateToDomain(&model)
}

func (r *whatsappTemplateRepository) ListTemplates(ctx context.Context, businessID uint, scope, status string, page, pageSize int) ([]entities.WhatsappTemplate, int64, error) {
	var rows []models.WhatsappTemplate
	var total int64

	query := r.db.Conn(ctx).Model(&models.WhatsappTemplate{}).
		Where("scope = ?", scope).
		Where("business_id = ? OR business_id IS NULL", businessID)
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		r.logger.Error().Err(err).Msg("Error counting whatsapp templates")
		return nil, 0, err
	}

	if err := query.Order("origin ASC, created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&rows).Error; err != nil {
		r.logger.Error().Err(err).Msg("Error listing whatsapp templates")
		return nil, 0, err
	}

	out := make([]entities.WhatsappTemplate, 0, len(rows))
	for i := range rows {
		item, err := templateToDomain(&rows[i])
		if err != nil {
			return nil, 0, err
		}
		out = append(out, *item)
	}

	return out, total, nil
}

func (r *whatsappTemplateRepository) DeleteTemplate(ctx context.Context, id uint) error {
	if err := r.db.Conn(ctx).Delete(&models.WhatsappTemplate{}, id).Error; err != nil {
		r.logger.Error().Err(err).Uint("id", id).Msg("Error deleting whatsapp template")
		return err
	}
	return nil
}

func (r *whatsappTemplateRepository) UpdateTemplateStatusByMeta(ctx context.Context, wabaID, name, language, status, reason string) error {
	updates := map[string]any{
		"status":      status,
		"reviewed_at": time.Now(),
	}
	if status == entities.TemplateStatusRejected {
		updates["rejected_reason"] = reason
	} else {
		updates["rejected_reason"] = ""
	}

	query := r.db.Conn(ctx).Model(&models.WhatsappTemplate{}).
		Where("waba_id = ? AND name = ?", strings.TrimSpace(wabaID), strings.TrimSpace(name))
	if language != "" {
		query = query.Where("language = ?", language)
	}

	result := query.Updates(updates)
	if result.Error != nil {
		r.logger.Error().Err(result.Error).Str("waba_id", wabaID).Str("name", name).
			Msg("Error updating whatsapp template status from webhook")
		return result.Error
	}

	if result.RowsAffected == 0 {
		r.logger.Warn().Str("waba_id", wabaID).Str("name", name).
			Msg("Webhook de plantilla sin fila local: la plantilla no se creo desde Probability")
	}

	return nil
}

func templateToModel(template *entities.WhatsappTemplate) (*models.WhatsappTemplate, error) {
	variables, err := json.Marshal(template.Variables)
	if err != nil {
		return nil, err
	}

	buttons, err := json.Marshal(template.Buttons)
	if err != nil {
		return nil, err
	}

	components, err := json.Marshal(template.Components)
	if err != nil {
		return nil, err
	}

	origin := template.Origin
	if origin == "" {
		origin = entities.TemplateOriginBusiness
	}

	scope := template.Scope
	if scope == "" {
		scope = entities.TemplateScopeScheduled
	}

	return &models.WhatsappTemplate{
		BusinessID:      template.BusinessID,
		Origin:          origin,
		Scope:           scope,
		IntegrationID:   template.IntegrationID,
		Name:            template.Name,
		Language:        template.Language,
		Category:        template.Category,
		BodyText:        template.BodyText,
		HeaderText:      template.HeaderText,
		FooterText:      template.FooterText,
		VariableMapping: datatypes.JSON(variables),
		Buttons:         datatypes.JSON(buttons),
		Components:      datatypes.JSON(components),
		WABAID:          template.WABAID,
		MetaTemplateID:  template.MetaTemplateID,
		Status:          template.Status,
		RejectedReason:  template.RejectedReason,
		SubmittedAt:     template.SubmittedAt,
		ReviewedAt:      template.ReviewedAt,
		LastSyncedAt:    template.LastSyncedAt,
		CreatedByID:     template.CreatedByID,
	}, nil
}

func templateToDomain(model *models.WhatsappTemplate) (*entities.WhatsappTemplate, error) {
	template := &entities.WhatsappTemplate{
		ID:             model.ID,
		BusinessID:     model.BusinessID,
		Origin:         model.Origin,
		Scope:          model.Scope,
		IntegrationID:  model.IntegrationID,
		Name:           model.Name,
		Language:       model.Language,
		Category:       model.Category,
		BodyText:       model.BodyText,
		HeaderText:     model.HeaderText,
		FooterText:     model.FooterText,
		WABAID:         model.WABAID,
		MetaTemplateID: model.MetaTemplateID,
		Status:         model.Status,
		RejectedReason: model.RejectedReason,
		SubmittedAt:    model.SubmittedAt,
		ReviewedAt:     model.ReviewedAt,
		LastSyncedAt:   model.LastSyncedAt,
		CreatedByID:    model.CreatedByID,
		CreatedAt:      model.CreatedAt,
		UpdatedAt:      model.UpdatedAt,
	}

	if len(model.VariableMapping) > 0 {
		if err := json.Unmarshal(model.VariableMapping, &template.Variables); err != nil {
			return nil, err
		}
	}

	if len(model.Buttons) > 0 {
		if err := json.Unmarshal(model.Buttons, &template.Buttons); err != nil {
			return nil, err
		}
	}

	return template, nil
}
