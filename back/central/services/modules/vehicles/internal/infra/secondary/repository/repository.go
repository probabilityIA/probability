package repository

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/vehicles/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/vehicles/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/vehicles/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/vehicles/internal/domain/ports"
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

func (r *Repository) Create(ctx context.Context, vehicle *entities.Vehicle) (*entities.Vehicle, error) {
	model := &models.Vehicle{
		BusinessID:         vehicle.BusinessID,
		Type:               vehicle.Type,
		LicensePlate:       vehicle.LicensePlate,
		Brand:              vehicle.Brand,
		VehicleModel:       vehicle.VehicleModel,
		Year:               vehicle.Year,
		Color:              vehicle.Color,
		Status:             vehicle.Status,
		WeightCapacityKg:   vehicle.WeightCapacityKg,
		VolumeCapacityM3:   vehicle.VolumeCapacityM3,
		PhotoURL:           vehicle.PhotoURL,
		InsuranceExpiry:    vehicle.InsuranceExpiry,
		RegistrationExpiry: vehicle.RegistrationExpiry,
	}

	if err := r.db.Conn(ctx).Create(model).Error; err != nil {
		return nil, err
	}

	vehicle.ID = model.ID
	vehicle.CreatedAt = model.CreatedAt
	vehicle.UpdatedAt = model.UpdatedAt
	return vehicle, nil
}

func (r *Repository) GetByID(ctx context.Context, businessID, vehicleID uint) (*entities.Vehicle, error) {
	var model models.Vehicle
	err := r.db.Conn(ctx).
		Where("id = ? AND business_id = ?", vehicleID, businessID).
		First(&model).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domainerrors.ErrVehicleNotFound
		}
		return nil, err
	}
	return modelToEntity(&model), nil
}

func (r *Repository) List(ctx context.Context, params dtos.ListVehiclesParams) ([]entities.Vehicle, int64, error) {
	var modelsList []models.Vehicle
	var total int64

	query := r.db.Conn(ctx).Model(&models.Vehicle{}).
		Where("business_id = ?", params.BusinessID)

	if params.Search != "" {
		like := "%" + params.Search + "%"
		query = query.Where("license_plate ILIKE ? OR brand ILIKE ? OR vehicle_model ILIKE ? OR color ILIKE ?",
			like, like, like, like)
	}

	if params.Type != "" {
		query = query.Where("type = ?", params.Type)
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

	vehicles := make([]entities.Vehicle, len(modelsList))
	for i, m := range modelsList {
		vehicles[i] = *modelToEntity(&m)
	}
	return vehicles, total, nil
}

func (r *Repository) Update(ctx context.Context, vehicle *entities.Vehicle) (*entities.Vehicle, error) {
	model := &models.Vehicle{
		Model:              gorm.Model{ID: vehicle.ID},
		BusinessID:         vehicle.BusinessID,
		Type:               vehicle.Type,
		LicensePlate:       vehicle.LicensePlate,
		Brand:              vehicle.Brand,
		VehicleModel:       vehicle.VehicleModel,
		Year:               vehicle.Year,
		Color:              vehicle.Color,
		Status:             vehicle.Status,
		WeightCapacityKg:   vehicle.WeightCapacityKg,
		VolumeCapacityM3:   vehicle.VolumeCapacityM3,
		PhotoURL:           vehicle.PhotoURL,
		InsuranceExpiry:    vehicle.InsuranceExpiry,
		RegistrationExpiry: vehicle.RegistrationExpiry,
	}

	if err := r.db.Conn(ctx).Save(model).Error; err != nil {
		return nil, err
	}

	vehicle.UpdatedAt = model.UpdatedAt
	return vehicle, nil
}

func (r *Repository) Delete(ctx context.Context, businessID, vehicleID uint) error {
	result := r.db.Conn(ctx).
		Where("id = ? AND business_id = ?", vehicleID, businessID).
		Delete(&models.Vehicle{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domainerrors.ErrVehicleNotFound
	}
	return nil
}

func (r *Repository) ExistsByLicensePlate(ctx context.Context, businessID uint, plate string, excludeID *uint) (bool, error) {
	var count int64
	query := r.db.Conn(ctx).Model(&models.Vehicle{}).
		Where("business_id = ? AND license_plate = ?", businessID, plate)
	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}
	err := query.Count(&count).Error
	return count > 0, err
}

func modelToEntity(m *models.Vehicle) *entities.Vehicle {
	return &entities.Vehicle{
		ID:                 m.ID,
		BusinessID:         m.BusinessID,
		Type:               m.Type,
		LicensePlate:       m.LicensePlate,
		Brand:              m.Brand,
		VehicleModel:       m.VehicleModel,
		Year:               m.Year,
		Color:              m.Color,
		Status:             m.Status,
		WeightCapacityKg:   m.WeightCapacityKg,
		VolumeCapacityM3:   m.VolumeCapacityM3,
		PhotoURL:           m.PhotoURL,
		InsuranceExpiry:    m.InsuranceExpiry,
		RegistrationExpiry: m.RegistrationExpiry,
		CreatedAt:          m.CreatedAt,
		UpdatedAt:          m.UpdatedAt,
	}
}
