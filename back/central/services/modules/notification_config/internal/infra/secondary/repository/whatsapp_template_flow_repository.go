package repository

import (
	"context"
	"errors"
	"strings"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
)

type templateFlowRepository struct {
	db     db.IDatabase
	logger log.ILogger
}

func NewTemplateFlowRepository(database db.IDatabase, logger log.ILogger) ports.ITemplateFlowRepository {
	return &templateFlowRepository{
		db:     database,
		logger: logger.WithModule("whatsapp_template_flow_repository"),
	}
}

type flowRow struct {
	ID                   uint
	BusinessID           uint
	FlowID               *uint
	SourceTemplateID     uint
	ButtonText           string
	TargetTemplateID     uint
	Enabled              bool
	SourceName           string
	TargetName           string
	TargetLanguage       string
	TargetStatus         string
	TargetHeaderMediaURL string
}

const flowSelect = `
	f.id, f.business_id, f.flow_id, f.source_template_id, f.button_text, f.target_template_id, f.enabled,
	s.name AS source_name,
	t.name AS target_name,
	t.language AS target_language,
	t.status AS target_status,
	t.header_media_url AS target_header_media_url
`

func (r *templateFlowRepository) baseQuery(ctx context.Context) *gorm.DB {
	return r.db.Conn(ctx).
		Table("whatsapp_template_flows AS f").
		Select(flowSelect).
		Joins("JOIN whatsapp_templates AS s ON s.id = f.source_template_id").
		Joins("JOIN whatsapp_templates AS t ON t.id = f.target_template_id").
		Where("f.deleted_at IS NULL")
}

func (r *templateFlowRepository) ListBySource(ctx context.Context, businessID, sourceTemplateID uint) ([]entities.TemplateFlow, error) {
	var rows []flowRow

	if err := r.baseQuery(ctx).
		Where("f.business_id = ? AND f.source_template_id = ?", businessID, sourceTemplateID).
		Order("f.id ASC").
		Scan(&rows).Error; err != nil {
		r.logger.Error().Err(err).Uint("source_template_id", sourceTemplateID).
			Msg("Error listando flujos de plantilla")
		return nil, err
	}

	return toFlowEntities(rows), nil
}

func (r *templateFlowRepository) ListByBusiness(ctx context.Context, businessID uint) ([]entities.TemplateFlow, error) {
	var rows []flowRow

	if err := r.baseQuery(ctx).
		Where("f.business_id = ?", businessID).
		Order("f.source_template_id ASC, f.id ASC").
		Scan(&rows).Error; err != nil {
		r.logger.Error().Err(err).Uint("business_id", businessID).
			Msg("Error listando flujos del negocio")
		return nil, err
	}

	return toFlowEntities(rows), nil
}

func (r *templateFlowRepository) ListByFlow(ctx context.Context, businessID, flowID uint) ([]entities.TemplateFlow, error) {
	var rows []flowRow

	if err := r.baseQuery(ctx).
		Where("f.business_id = ? AND f.flow_id = ?", businessID, flowID).
		Order("f.source_template_id ASC, f.id ASC").
		Scan(&rows).Error; err != nil {
		r.logger.Error().Err(err).Uint("flow_id", flowID).
			Msg("Error listando las transiciones del flujo")
		return nil, err
	}

	return toFlowEntities(rows), nil
}

func (r *templateFlowRepository) Resolve(ctx context.Context, businessID, sourceTemplateID uint, buttonText string) (*entities.TemplateFlow, error) {
	var row flowRow

	err := r.baseQuery(ctx).
		Where("f.business_id = ? AND f.source_template_id = ?", businessID, sourceTemplateID).
		Where("lower(f.button_text) = lower(?)", strings.TrimSpace(buttonText)).
		Where("f.enabled = true").
		Limit(1).
		Scan(&row).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	if row.ID == 0 {
		return nil, nil
	}

	flow := toFlowEntity(row)
	return &flow, nil
}

func (r *templateFlowRepository) ReplaceForSource(ctx context.Context, businessID, sourceTemplateID uint, flows []entities.TemplateFlow) error {
	return r.db.Conn(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Unscoped().
			Where("business_id = ? AND source_template_id = ?", businessID, sourceTemplateID).
			Delete(&models.WhatsappTemplateFlow{}).Error; err != nil {
			return err
		}

		if len(flows) == 0 {
			return nil
		}

		rows := make([]models.WhatsappTemplateFlow, 0, len(flows))
		for _, flow := range flows {
			rows = append(rows, models.WhatsappTemplateFlow{
				BusinessID:       businessID,
				FlowID:           flow.FlowID,
				SourceTemplateID: sourceTemplateID,
				ButtonText:       flow.ButtonText,
				TargetTemplateID: flow.TargetTemplateID,
				Enabled:          flow.Enabled,
			})
		}

		return tx.Create(&rows).Error
	})
}

func (r *templateFlowRepository) TemplateNameByOutboundMessageID(ctx context.Context, messageID string) (string, error) {
	var name string

	err := r.db.Conn(ctx).
		Table("whatsapp_message_logs").
		Select("template_name").
		Where("message_id = ? AND direction = ?", strings.TrimSpace(messageID), "outbound").
		Order("created_at DESC").
		Limit(1).
		Scan(&name).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil
		}
		return "", err
	}

	return name, nil
}

func toFlowEntities(rows []flowRow) []entities.TemplateFlow {
	out := make([]entities.TemplateFlow, 0, len(rows))
	for _, row := range rows {
		out = append(out, toFlowEntity(row))
	}
	return out
}

func toFlowEntity(row flowRow) entities.TemplateFlow {
	return entities.TemplateFlow{
		ID:                   row.ID,
		BusinessID:           row.BusinessID,
		FlowID:               row.FlowID,
		SourceTemplateID:     row.SourceTemplateID,
		ButtonText:           row.ButtonText,
		TargetTemplateID:     row.TargetTemplateID,
		Enabled:              row.Enabled,
		SourceName:           row.SourceName,
		TargetName:           row.TargetName,
		TargetLanguage:       row.TargetLanguage,
		TargetStatus:         row.TargetStatus,
		TargetHeaderMediaURL: row.TargetHeaderMediaURL,
	}
}
