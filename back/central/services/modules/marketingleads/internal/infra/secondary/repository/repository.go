package repository

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/marketingleads/internal/domain/entities"
	"github.com/secamc93/probability/back/migration/shared/models"
)

func leadToModel(l *entities.MarketingLead) *models.MarketingLead {
	return &models.MarketingLead{
		Name:       l.Name,
		Email:      l.Email,
		Phone:      l.Phone,
		SurveySlug: l.SurveySlug,
		ScoreTotal: l.ScoreTotal,
		ScoreMax:   l.ScoreMax,
		Level:      l.Level,
	}
}

func (r *Repository) CreateLead(ctx context.Context, lead *entities.MarketingLead) error {
	model := leadToModel(lead)
	if err := r.db.Conn(ctx).Create(model).Error; err != nil {
		return err
	}
	lead.ID = model.ID
	lead.CreatedAt = model.CreatedAt
	return nil
}

func (r *Repository) SetWhatsAppMessageID(ctx context.Context, leadID uint, messageID string) error {
	return r.db.Conn(ctx).
		Model(&models.MarketingLead{}).
		Where("id = ?", leadID).
		Update("whats_app_message_id", messageID).Error
}

func modelToLead(m *models.MarketingLead) entities.MarketingLead {
	return entities.MarketingLead{
		ID:                m.ID,
		Name:              m.Name,
		Email:             m.Email,
		Phone:             m.Phone,
		SurveySlug:        m.SurveySlug,
		ScoreTotal:        m.ScoreTotal,
		ScoreMax:          m.ScoreMax,
		Level:             m.Level,
		WhatsAppMessageID: m.WhatsAppMessageID,
		CreatedAt:         m.CreatedAt,
	}
}

func (r *Repository) ListLeads(ctx context.Context, page, pageSize int) ([]entities.MarketingLead, int64, error) {
	var models_ []models.MarketingLead
	var total int64

	if err := r.db.Conn(ctx).Model(&models.MarketingLead{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := r.db.Conn(ctx).
		Order("id DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&models_).Error; err != nil {
		return nil, 0, err
	}

	leads := make([]entities.MarketingLead, 0, len(models_))
	for i := range models_ {
		leads = append(leads, modelToLead(&models_[i]))
	}
	return leads, total, nil
}
