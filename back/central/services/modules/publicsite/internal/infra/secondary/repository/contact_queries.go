package repository

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/publicsite/internal/domain/dtos"
	"github.com/secamc93/probability/back/migration/shared/models"
)

func (r *Repository) SaveContactSubmission(ctx context.Context, businessID uint, dto *dtos.ContactFormDTO) error {
	submission := models.ContactSubmission{
		BusinessID: businessID,
		Name:       dto.Name,
		Email:      dto.Email,
		Phone:      dto.Phone,
		Message:    dto.Message,
		Source:     "website",
		Status:     "new",
	}
	return r.db.Conn(ctx).Create(&submission).Error
}
