package repository

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/paymentstatus/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/paymentstatus/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/paymentstatus/internal/domain/ports"
	"github.com/secamc93/probability/back/central/services/modules/paymentstatus/internal/infra/secondary/repository/mappers"
	"github.com/secamc93/probability/back/central/services/modules/paymentstatus/internal/infra/secondary/repository/models"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/log"
	"gorm.io/gorm"
)

// Repository implementa ports.IRepository
type Repository struct {
	db     db.IDatabase
	logger log.ILogger
}

// New crea una nueva instancia del repositorio
func New(database db.IDatabase, logger log.ILogger) ports.IRepository {
	return &Repository{
		db:     database,
		logger: logger,
	}
}

func (r *Repository) GetPaymentStatusByCode(ctx context.Context, code string) (*entities.PaymentStatus, error) {
	var status models.PaymentStatus

	// ✅ GORM infiere la tabla desde PaymentStatus.TableName()
	err := r.db.Conn(ctx).Where("code = ? AND is_active = ?", code, true).First(&status).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domainerrors.ErrPaymentStatusNotFound
		}
		return nil, err
	}

	// ✅ Convertir a dominio usando mapper
	domain := mappers.ToDomain(&status)
	return &domain, nil
}

func (r *Repository) GetPaymentStatusIDByCode(ctx context.Context, code string) (*uint, error) {
	var status models.PaymentStatus

	// ✅ GORM infiere la tabla desde PaymentStatus.TableName()
	err := r.db.Conn(ctx).Where("code = ? AND is_active = ?", code, true).First(&status).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// No se encontró, retornar nil sin error
			return nil, nil
		}
		return nil, err
	}

	return &status.ID, nil
}

func (r *Repository) ListPaymentStatuses(ctx context.Context, isActive *bool) ([]entities.PaymentStatus, error) {
	var statuses []models.PaymentStatus

	// ✅ GORM infiere la tabla desde PaymentStatus.TableName()
	query := r.db.Conn(ctx)

	if isActive != nil {
		query = query.Where("is_active = ?", *isActive)
	}

	err := query.Order("code ASC").Find(&statuses).Error
	if err != nil {
		return nil, err
	}

	// ✅ Convertir a dominio usando mapper
	return mappers.ToDomainList(statuses), nil
}
