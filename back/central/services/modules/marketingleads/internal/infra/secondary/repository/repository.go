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
