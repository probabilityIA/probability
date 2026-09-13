package repository

import (
	"context"
	"strings"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/migration/shared/models"
)

type flowRepository struct {
	db     db.IDatabase
	logger log.ILogger
}

func NewFlowRepository(database db.IDatabase, logger log.ILogger) ports.IFlowRepository {
	return &flowRepository{
		db:     database,
		logger: logger.WithModule("whatsapp_flow_repository"),
	}
}

type flowListRow struct {
	ID                 uint
	BusinessID         uint
	Name               string
	Description        string
	RootTemplateID     *uint
	Enabled            bool
	RootTemplateName   string
	RootTemplateStatus string
	StepCount          int
	PendingCount       int
}

const flowListQuery = `
SELECT
	f.id,
	f.business_id,
	f.name,
	f.description,
	f.root_template_id,
	f.enabled,
	COALESCE(t.name, '') AS root_template_name,
	COALESCE(t.status, '') AS root_template_status,
	(
		SELECT count(*) FROM whatsapp_template_flows tf
		WHERE tf.flow_id = f.id AND tf.deleted_at IS NULL
	) AS step_count,
	(
		SELECT count(*) FROM whatsapp_template_flows tf
		JOIN whatsapp_templates tt ON tt.id = tf.target_template_id
		WHERE tf.flow_id = f.id AND tf.deleted_at IS NULL AND tt.status <> 'approved'
	) AS pending_count
FROM whatsapp_flows f
LEFT JOIN whatsapp_templates t ON t.id = f.root_template_id
WHERE f.business_id = @business_id AND f.deleted_at IS NULL
ORDER BY f.created_at DESC
`

func (r *flowRepository) Create(ctx context.Context, flow *entities.Flow) error {
	model := &models.WhatsappFlow{
		BusinessID:     flow.BusinessID,
		Name:           flow.Name,
		Description:    flow.Description,
		RootTemplateID: flow.RootTemplateID,
		Enabled:        flow.Enabled,
	}

	if err := r.db.Conn(ctx).Create(model).Error; err != nil {
		r.logger.Error().Err(err).Str("name", flow.Name).Msg("Error creando el flujo")
		return err
	}

	flow.ID = model.ID
	flow.CreatedAt = model.CreatedAt
	flow.UpdatedAt = model.UpdatedAt

	return nil
}

func (r *flowRepository) Update(ctx context.Context, flow *entities.Flow) error {
	updates := map[string]any{
		"name":             flow.Name,
		"description":      flow.Description,
		"root_template_id": flow.RootTemplateID,
		"enabled":          flow.Enabled,
	}

	if err := r.db.Conn(ctx).Model(&models.WhatsappFlow{}).
		Where("id = ? AND business_id = ?", flow.ID, flow.BusinessID).
		Updates(updates).Error; err != nil {
		r.logger.Error().Err(err).Uint("id", flow.ID).Msg("Error actualizando el flujo")
		return err
	}

	return nil
}

func (r *flowRepository) GetByID(ctx context.Context, id, businessID uint) (*entities.Flow, error) {
	flows, err := r.List(ctx, businessID)
	if err != nil {
		return nil, err
	}

	for i := range flows {
		if flows[i].ID == id {
			return &flows[i], nil
		}
	}

	return nil, nil
}

func (r *flowRepository) List(ctx context.Context, businessID uint) ([]entities.Flow, error) {
	var rows []flowListRow

	if err := r.db.Conn(ctx).Raw(flowListQuery, map[string]any{
		"business_id": businessID,
	}).Scan(&rows).Error; err != nil {
		r.logger.Error().Err(err).Uint("business_id", businessID).Msg("Error listando los flujos")
		return nil, err
	}

	out := make([]entities.Flow, 0, len(rows))
	for _, row := range rows {
		out = append(out, entities.Flow{
			ID:                 row.ID,
			BusinessID:         row.BusinessID,
			Name:               row.Name,
			Description:        row.Description,
			RootTemplateID:     row.RootTemplateID,
			Enabled:            row.Enabled,
			RootTemplateName:   row.RootTemplateName,
			RootTemplateStatus: row.RootTemplateStatus,
			StepCount:          row.StepCount,
			PendingCount:       row.PendingCount,
		})
	}

	return out, nil
}

func (r *flowRepository) Delete(ctx context.Context, id, businessID uint) error {
	if err := r.db.Conn(ctx).
		Where("id = ? AND business_id = ?", id, businessID).
		Delete(&models.WhatsappFlow{}).Error; err != nil {
		r.logger.Error().Err(err).Uint("id", id).Msg("Error eliminando el flujo")
		return err
	}

	return nil
}

func (r *flowRepository) ExistsByName(ctx context.Context, businessID uint, name string, excludeID uint) (bool, error) {
	var total int64

	query := r.db.Conn(ctx).Model(&models.WhatsappFlow{}).
		Where("business_id = ?", businessID).
		Where("lower(name) = lower(?)", strings.TrimSpace(name))

	if excludeID > 0 {
		query = query.Where("id <> ?", excludeID)
	}

	if err := query.Count(&total).Error; err != nil {
		return false, err
	}

	return total > 0, nil
}
