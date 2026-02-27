package repository

import (
	"context"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/customers/internal/domain/dtos"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/customers/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/customers/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/customers/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
)

// Repository implementa ports.IRepository
type Repository struct {
	db db.IDatabase
}

// New crea una nueva instancia del repositorio
func New(database db.IDatabase) ports.IRepository {
	return &Repository{db: database}
}

func (r *Repository) Create(ctx context.Context, client *entities.Client) (*entities.Client, error) {
	model := &models.Client{
		BusinessID: client.BusinessID,
		Name:       client.Name,
		Email:      client.Email,
		Phone:      client.Phone,
		Dni:        client.Dni,
	}

	if err := r.db.Conn(ctx).Create(model).Error; err != nil {
		return nil, err
	}

	client.ID = model.ID
	client.CreatedAt = model.CreatedAt
	client.UpdatedAt = model.UpdatedAt
	return client, nil
}

func (r *Repository) GetByID(ctx context.Context, businessID, clientID uint) (*entities.Client, error) {
	var model models.Client
	err := r.db.Conn(ctx).
		Where("id = ? AND business_id = ?", clientID, businessID).
		First(&model).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domainerrors.ErrClientNotFound
		}
		return nil, err
	}
	return modelToEntity(&model), nil
}

func (r *Repository) List(ctx context.Context, params dtos.ListClientsParams) ([]entities.Client, int64, error) {
	var modelsList []models.Client
	var total int64

	query := r.db.Conn(ctx).Model(&models.Client{}).
		Where("business_id = ?", params.BusinessID)

	if params.Search != "" {
		like := "%" + params.Search + "%"
		query = query.Where("name ILIKE ? OR email ILIKE ? OR phone ILIKE ?", like, like, like)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := params.Offset()
	if err := query.Offset(offset).Limit(params.PageSize).Order("created_at DESC").Find(&modelsList).Error; err != nil {
		return nil, 0, err
	}

	clients := make([]entities.Client, len(modelsList))
	for i, m := range modelsList {
		clients[i] = *modelToEntity(&m)
	}
	return clients, total, nil
}

func (r *Repository) Update(ctx context.Context, client *entities.Client) (*entities.Client, error) {
	model := &models.Client{
		Model:      gorm.Model{ID: client.ID},
		BusinessID: client.BusinessID,
		Name:       client.Name,
		Email:      client.Email,
		Phone:      client.Phone,
		Dni:        client.Dni,
	}

	if err := r.db.Conn(ctx).Save(model).Error; err != nil {
		return nil, err
	}

	client.UpdatedAt = model.UpdatedAt
	return client, nil
}

func (r *Repository) Delete(ctx context.Context, businessID, clientID uint) error {
	result := r.db.Conn(ctx).
		Where("id = ? AND business_id = ?", clientID, businessID).
		Delete(&models.Client{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domainerrors.ErrClientNotFound
	}
	return nil
}

func (r *Repository) ExistsByEmail(ctx context.Context, businessID uint, email string, excludeID *uint) (bool, error) {
	var count int64
	query := r.db.Conn(ctx).Model(&models.Client{}).
		Where("business_id = ? AND email = ?", businessID, email)
	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}
	err := query.Count(&count).Error
	return count > 0, err
}

func (r *Repository) ExistsByDni(ctx context.Context, businessID uint, dni string, excludeID *uint) (bool, error) {
	var count int64
	query := r.db.Conn(ctx).Model(&models.Client{}).
		Where("business_id = ? AND dni = ?", businessID, dni)
	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}
	err := query.Count(&count).Error
	return count > 0, err
}

// GetOrderStats consulta la tabla orders directamente.
// Replicado localmente para evitar compartir repositorios entre módulos.
func (r *Repository) GetOrderStats(ctx context.Context, clientID uint) (int64, float64, *time.Time, error) {
	type statsResult struct {
		OrderCount  int64
		TotalSpent  float64
		LastOrderAt *time.Time
	}

	var result statsResult
	err := r.db.Conn(ctx).
		Table("orders").
		Select("COUNT(*) AS order_count, COALESCE(SUM(total), 0) AS total_spent, MAX(created_at) AS last_order_at").
		Where("customer_id = ? AND deleted_at IS NULL", clientID).
		Scan(&result).Error
	if err != nil {
		return 0, 0, nil, err
	}

	return result.OrderCount, result.TotalSpent, result.LastOrderAt, nil
}

func modelToEntity(m *models.Client) *entities.Client {
	return &entities.Client{
		ID:         m.ID,
		BusinessID: m.BusinessID,
		Name:       m.Name,
		Email:      m.Email,
		Phone:      m.Phone,
		Dni:        m.Dni,
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
	}
}
