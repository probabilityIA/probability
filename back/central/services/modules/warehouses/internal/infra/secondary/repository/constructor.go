package repository

import (
	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/migration/shared/models"
)

// Repository implementa ports.IRepository
type Repository struct {
	db db.IDatabase
}

// New crea una nueva instancia del repositorio
func New(database db.IDatabase) ports.IRepository {
	return &Repository{db: database}
}

func warehouseModelToEntity(m *models.Warehouse) *entities.Warehouse {
	w := &entities.Warehouse{
		ID:            m.ID,
		BusinessID:    m.BusinessID,
		Name:          m.Name,
		Code:          m.Code,
		Address:       m.Address,
		City:          m.City,
		State:         m.State,
		Country:       m.Country,
		ZipCode:       m.ZipCode,
		Phone:         m.Phone,
		ContactName:   m.ContactName,
		ContactEmail:  m.ContactEmail,
		IsActive:      m.IsActive,
		IsDefault:     m.IsDefault,
		IsFulfillment: m.IsFulfillment,
		Company:       m.Company,
		FirstName:     m.FirstName,
		LastName:      m.LastName,
		Email:         m.Email,
		Suburb:        m.Suburb,
		CityDaneCode:  m.CityDaneCode,
		PostalCode:    m.PostalCode,
		Street:        m.Street,
		Latitude:      m.Latitude,
		Longitude:     m.Longitude,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
	}

	if len(m.Locations) > 0 {
		w.Locations = make([]entities.WarehouseLocation, len(m.Locations))
		for i, loc := range m.Locations {
			w.Locations[i] = *locationModelToEntity(&loc)
		}
	}

	return w
}

func locationModelToEntity(m *models.WarehouseLocation) *entities.WarehouseLocation {
	return &entities.WarehouseLocation{
		ID:            m.ID,
		WarehouseID:   m.WarehouseID,
		Name:          m.Name,
		Code:          m.Code,
		Type:          m.Type,
		IsActive:      m.IsActive,
		IsFulfillment: m.IsFulfillment,
		Capacity:      m.Capacity,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
	}
}
