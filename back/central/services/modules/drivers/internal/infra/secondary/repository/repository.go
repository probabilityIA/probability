package repository

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/drivers/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/drivers/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/drivers/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/drivers/internal/domain/ports"
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

func (r *Repository) Create(ctx context.Context, driver *entities.Driver) (*entities.Driver, error) {
	model := &models.Driver{
		BusinessID:     driver.BusinessID,
		FirstName:      driver.FirstName,
		LastName:       driver.LastName,
		Email:          driver.Email,
		Phone:          driver.Phone,
		Identification: driver.Identification,
		Status:         driver.Status,
		PhotoURL:       driver.PhotoURL,
		LicenseType:    driver.LicenseType,
		LicenseExpiry:  driver.LicenseExpiry,
		WarehouseID:    driver.WarehouseID,
		Notes:          driver.Notes,
	}

	if err := r.db.Conn(ctx).Create(model).Error; err != nil {
		return nil, err
	}

	driver.ID = model.ID
	driver.CreatedAt = model.CreatedAt
	driver.UpdatedAt = model.UpdatedAt
	return driver, nil
}

func (r *Repository) GetByID(ctx context.Context, businessID, driverID uint) (*entities.Driver, error) {
	var model models.Driver
	err := r.db.Conn(ctx).
		Where("id = ? AND business_id = ?", driverID, businessID).
		First(&model).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domainerrors.ErrDriverNotFound
		}
		return nil, err
	}
	return modelToEntity(&model), nil
}

func (r *Repository) List(ctx context.Context, params dtos.ListDriversParams) ([]entities.Driver, int64, error) {
	var modelsList []models.Driver
	var total int64

	query := r.db.Conn(ctx).Model(&models.Driver{}).
		Where("business_id = ?", params.BusinessID)

	if params.Search != "" {
		like := "%" + params.Search + "%"
		query = query.Where("first_name ILIKE ? OR last_name ILIKE ? OR email ILIKE ? OR phone ILIKE ? OR identification ILIKE ?",
			like, like, like, like, like)
	}

	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := params.Offset()
	if err := query.Offset(offset).Limit(params.PageSize).
		Order("created_at DESC").
		Find(&modelsList).Error; err != nil {
		return nil, 0, err
	}

	drivers := make([]entities.Driver, len(modelsList))
	for i, m := range modelsList {
		drivers[i] = *modelToEntity(&m)
	}
	return drivers, total, nil
}

func (r *Repository) Update(ctx context.Context, driver *entities.Driver) (*entities.Driver, error) {
	model := &models.Driver{
		Model:          gorm.Model{ID: driver.ID},
		BusinessID:     driver.BusinessID,
		FirstName:      driver.FirstName,
		LastName:       driver.LastName,
		Email:          driver.Email,
		Phone:          driver.Phone,
		Identification: driver.Identification,
		Status:         driver.Status,
		PhotoURL:       driver.PhotoURL,
		LicenseType:    driver.LicenseType,
		LicenseExpiry:  driver.LicenseExpiry,
		WarehouseID:    driver.WarehouseID,
		Notes:          driver.Notes,
	}

	if err := r.db.Conn(ctx).Save(model).Error; err != nil {
		return nil, err
	}

	driver.UpdatedAt = model.UpdatedAt
	return driver, nil
}

func (r *Repository) Delete(ctx context.Context, businessID, driverID uint) error {
	result := r.db.Conn(ctx).
		Where("id = ? AND business_id = ?", driverID, businessID).
		Delete(&models.Driver{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domainerrors.ErrDriverNotFound
	}
	return nil
}

func (r *Repository) ExistsByIdentification(ctx context.Context, businessID uint, identification string, excludeID *uint) (bool, error) {
	var count int64
	query := r.db.Conn(ctx).Model(&models.Driver{}).
		Where("business_id = ? AND identification = ?", businessID, identification)
	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}
	err := query.Count(&count).Error
	return count > 0, err
}

func modelToEntity(m *models.Driver) *entities.Driver {
	return &entities.Driver{
		ID:             m.ID,
		BusinessID:     m.BusinessID,
		FirstName:      m.FirstName,
		LastName:       m.LastName,
		Email:          m.Email,
		Phone:          m.Phone,
		Identification: m.Identification,
		Status:         m.Status,
		PhotoURL:       m.PhotoURL,
		LicenseType:    m.LicenseType,
		LicenseExpiry:  m.LicenseExpiry,
		WarehouseID:    m.WarehouseID,
		Notes:          m.Notes,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
	}
}
