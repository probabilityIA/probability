package repository

import (
	"context"
	"errors"

	"github.com/secamc93/probability/back/central/services/modules/fulfillmentstatus/domain"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
)

// Repository implementa domain.IRepository
type Repository struct {
	db     db.IDatabase
	logger log.ILogger
}

// New crea una nueva instancia del repositorio
func New(database db.IDatabase, logger log.ILogger) domain.IRepository {
	return &Repository{
		db:     database,
		logger: logger,
	}
}

func (r *Repository) GetFulfillmentStatusByCode(ctx context.Context, code string) (*models.FulfillmentStatus, error) {
	var status models.FulfillmentStatus
	err := r.db.Conn(ctx).Where("code = ? AND is_active = ?", code, true).First(&status).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("fulfillment status not found")
		}
		return nil, err
	}
	return &status, nil
}

func (r *Repository) GetFulfillmentStatusIDByCode(ctx context.Context, code string) (*uint, error) {
	var status models.FulfillmentStatus
	err := r.db.Conn(ctx).Where("code = ? AND is_active = ?", code, true).First(&status).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// No se encontró, retornar nil sin error
			return nil, nil
		}
		return nil, err
	}
	return &status.ID, nil
}

func (r *Repository) ListFulfillmentStatuses(ctx context.Context, isActive *bool) ([]models.FulfillmentStatus, error) {
	var statuses []models.FulfillmentStatus
	query := r.db.Conn(ctx).Model(&models.FulfillmentStatus{})

	if isActive != nil {
		query = query.Where("is_active = ?", *isActive)
	}

	err := query.Order("code ASC").Find(&statuses).Error
	if err != nil {
		return nil, err
	}

	return statuses, nil
}
