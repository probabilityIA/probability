package repository

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/customers/internal/domain/dtos"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/customers/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/customers/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/customers/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
)

type Repository struct {
	db db.IDatabase
}

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

type clientWithOrders struct {
	models.Client
	TotalOrders int64 `gorm:"column:total_orders"`
}

func (r *Repository) List(ctx context.Context, params dtos.ListClientsParams) ([]entities.Client, int64, error) {
	var total int64

	countQuery := r.db.Conn(ctx).Model(&models.Client{}).
		Where("client.business_id = ?", params.BusinessID)

	if params.Search != "" {
		like := "%" + params.Search + "%"
		countQuery = countQuery.Where("client.name ILIKE ? OR client.email ILIKE ? OR client.phone ILIKE ? OR client.dni ILIKE ?", like, like, like, like)
	}
	if params.Email != "" {
		countQuery = countQuery.Where("client.email = ?", params.Email)
	}
	if params.Dni != "" {
		countQuery = countQuery.Where("client.dni = ?", params.Dni)
	}
	if params.Name != "" {
		countQuery = countQuery.Where("client.name ILIKE ?", "%"+params.Name+"%")
	}

	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var rows []clientWithOrders
	query := r.db.Conn(ctx).
		Table("client").
		Select("client.*, COALESCE(cs.total_orders, 0) AS total_orders").
		Joins("LEFT JOIN customer_summary cs ON cs.customer_id = client.id AND cs.business_id = client.business_id AND cs.deleted_at IS NULL").
		Where("client.business_id = ? AND client.deleted_at IS NULL", params.BusinessID)

	if params.Search != "" {
		like := "%" + params.Search + "%"
		query = query.Where("client.name ILIKE ? OR client.email ILIKE ? OR client.phone ILIKE ? OR client.dni ILIKE ?", like, like, like, like)
	}
	if params.Email != "" {
		query = query.Where("client.email = ?", params.Email)
	}
	if params.Dni != "" {
		query = query.Where("client.dni = ?", params.Dni)
	}
	if params.Name != "" {
		query = query.Where("client.name ILIKE ?", "%"+params.Name+"%")
	}

	offset := params.Offset()
	if err := query.Order("total_orders DESC, client.created_at DESC").
		Offset(offset).Limit(params.PageSize).
		Find(&rows).Error; err != nil {
		return nil, 0, err
	}

	clients := make([]entities.Client, len(rows))
	for i, row := range rows {
		c := modelToEntity(&row.Client)
		c.OrderCount = row.TotalOrders
		clients[i] = *c
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
	if email == "" {
		return false, nil
	}
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


func (r *Repository) FindClientByPhone(ctx context.Context, businessID uint, phone string) (*entities.Client, error) {
	var model models.Client
	err := r.db.Conn(ctx).
		Where("business_id = ? AND phone = ?", businessID, phone).
		First(&model).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return modelToEntity(&model), nil
}

func (r *Repository) FindClientByDNI(ctx context.Context, businessID uint, dni string) (*entities.Client, error) {
	var model models.Client
	err := r.db.Conn(ctx).
		Where("business_id = ? AND dni = ?", businessID, dni).
		First(&model).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return modelToEntity(&model), nil
}

func (r *Repository) FindClientByEmail(ctx context.Context, businessID uint, email string) (*entities.Client, error) {
	var model models.Client
	err := r.db.Conn(ctx).
		Where("business_id = ? AND email = ?", businessID, email).
		First(&model).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return modelToEntity(&model), nil
}

func (r *Repository) UpdateClientFields(ctx context.Context, clientID uint, updates map[string]any) error {
	return r.db.Conn(ctx).Model(&models.Client{}).Where("id = ?", clientID).Updates(updates).Error
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
