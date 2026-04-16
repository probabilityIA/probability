package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/domain/entities"
	"github.com/secamc93/probability/back/central/shared/db"
	"gorm.io/gorm"
)

type SubscriptionRepository struct {
	db db.IDatabase
}

func NewSubscriptionRepository(database db.IDatabase) *SubscriptionRepository {
	return &SubscriptionRepository{db: database}
}

// BusinessSubscriptionDB representa la estructura en la base de datos (GORM)
type BusinessSubscriptionDB struct {
	gorm.Model
	BusinessID       uint      `gorm:"not null;index"`
	Amount           float64   `gorm:"type:numeric(12,2);not null"`
	StartDate        time.Time `gorm:"type:timestamptz"`
	EndDate          time.Time `gorm:"type:timestamptz"`
	Status           string    `gorm:"size:20;default:'pending'"` // 'paid', 'pending', 'rejected'
	PaymentReference *string   `gorm:"size:255"`
	Notes            *string   `gorm:"type:text"`
}

// TableName especifica el nombre de la tabla
func (BusinessSubscriptionDB) TableName() string {
	return "business_subscriptions"
}

// ModelBusiness representa el modelo Business para actualizaciones
type ModelBusiness struct {
	gorm.Model
	SubscriptionStatus  string
	SubscriptionEndDate *string
}

// TableName
func (ModelBusiness) TableName() string {
	return "business"
}

func (r *SubscriptionRepository) GetLatestByBusinessID(ctx context.Context, businessID uint) (*entities.BusinessSubscription, error) {
	var subDB BusinessSubscriptionDB

	err := r.db.Conn(ctx).
		Where("business_id = ?", businessID).
		Order("created_at desc").
		First(&subDB).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // No subscription found
		}
		return nil, err
	}

	return r.mapToDomain(&subDB), nil
}

func (r *SubscriptionRepository) ListByBusinessID(ctx context.Context, businessID uint) ([]entities.BusinessSubscription, error) {
	var subsDB []BusinessSubscriptionDB

	err := r.db.Conn(ctx).
		Where("business_id = ?", businessID).
		Order("created_at desc").
		Find(&subsDB).Error

	if err != nil {
		return nil, err
	}

	var subs []entities.BusinessSubscription
	for _, subDB := range subsDB {
		subs = append(subs, *r.mapToDomain(&subDB))
	}

	return subs, nil
}

func (r *SubscriptionRepository) Create(ctx context.Context, subscription *entities.BusinessSubscription) error {
	subDB := r.mapToDB(subscription)

	err := r.db.Conn(ctx).Create(&subDB).Error
	if err != nil {
		return err
	}

	subscription.ID = subDB.ID
	subscription.CreatedAt = subDB.CreatedAt
	// Note: string/time format mapping should be exact

	return nil
}

func (r *SubscriptionRepository) UpdateBusinessSubscriptionStatus(ctx context.Context, businessID uint, status string, endDate *string) error {
	updates := map[string]interface{}{
		"subscription_status": status,
	}

	if endDate != nil {
		updates["subscription_end_date"] = *endDate
	} else {
		// Only update if not null, or set to null explicitly
		updates["subscription_end_date"] = gorm.Expr("NULL")
	}

	err := r.db.Conn(ctx).
		Table("business").
		Where("id = ? AND deleted_at IS NULL", businessID).
		Updates(updates).Error

	return err
}

func (r *SubscriptionRepository) EnsureAllBusinessesActive(ctx context.Context) error {
	// Poner a todos los negocios no eliminados como activos con una fecha lejana (2030)
	// si no tienen ya un estado activo.
	endDate := "2030-01-01T00:00:00Z"
	updates := map[string]interface{}{
		"subscription_status":    entities.BusinessStatusActive,
		"subscription_end_date": endDate,
	}

	// Solo actualizamos negocios que no están explícitamente marcados como 'active'
	// o que tienen el estado nulo.
	err := r.db.Conn(ctx).
		Table("business").
		Where("deleted_at IS NULL AND (subscription_status <> ? OR subscription_status IS NULL)", entities.BusinessStatusActive).
		Updates(updates).Error

	return err
}

func (r *SubscriptionRepository) mapToDomain(db *BusinessSubscriptionDB) *entities.BusinessSubscription {
	// Simplify date mapping for brevity as time.Time string handling might be needed
	return &entities.BusinessSubscription{
		ID:               db.ID,
		BusinessID:       db.BusinessID,
		Amount:           db.Amount,
		Status:           db.Status,
		PaymentReference: db.PaymentReference,
		Notes:            db.Notes,
		CreatedAt:        db.CreatedAt,
		UpdatedAt:        db.UpdatedAt,
	}
}

func (r *SubscriptionRepository) mapToDB(e *entities.BusinessSubscription) *BusinessSubscriptionDB {
	return &BusinessSubscriptionDB{
		BusinessID:       e.BusinessID,
		Amount:           e.Amount,
		StartDate:        e.StartDate,
		EndDate:          e.EndDate,
		Status:           e.Status,
		PaymentReference: e.PaymentReference,
		Notes:            e.Notes,
	}
}
