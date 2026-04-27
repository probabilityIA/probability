package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/core/internal/domain"
	"github.com/secamc93/probability/back/migration/shared/models"
)

func (r *Repository) RecordCredentialReveal(ctx context.Context, audit *domain.CredentialRevealAudit) error {
	row := models.CredentialRevealAudit{
		UserID:            audit.UserID,
		BusinessID:        audit.BusinessID,
		IntegrationTypeID: audit.IntegrationTypeID,
		IntegrationCode:   audit.IntegrationCode,
		IPAddress:         audit.IPAddress,
		UserAgent:         audit.UserAgent,
	}
	if err := r.db.Conn(ctx).Create(&row).Error; err != nil {
		return fmt.Errorf("create credential_reveal_audit: %w", err)
	}
	return nil
}
