package repository

import (
	"context"
	"errors"

	"github.com/secamc93/probability/back/central/services/modules/paymentstatus/domain"
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

func (r *Repository) GetPaymentStatusByCode(ctx context.Context, code string) (*models.PaymentStatus, error) {
	var status models.PaymentStatus
	err := r.db.Conn(ctx).Where("code = ? AND is_active = ?", code, true).First(&status).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("payment status not found")
		}
		return nil, err
	}
	return &status, nil
}

func (r *Repository) GetPaymentStatusIDByCode(ctx context.Context, code string) (*uint, error) {
	var status models.PaymentStatus
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

func (r *Repository) ListPaymentStatuses(ctx context.Context, isActive *bool) ([]models.PaymentStatus, error) {
	var statuses []models.PaymentStatus
	query := r.db.Conn(ctx).Model(&models.PaymentStatus{})

	if isActive != nil {
		query = query.Where("is_active = ?", *isActive)
	}

	err := query.Order("code ASC").Find(&statuses).Error
	if err != nil {
		return nil, err
	}

	return statuses, nil
}
