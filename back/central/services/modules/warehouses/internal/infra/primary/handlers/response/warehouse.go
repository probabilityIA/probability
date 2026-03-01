package response

import (
	"time"

	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/entities"
)

// WarehouseResponse respuesta básica de bodega (para listado)
type WarehouseResponse struct {
	ID            uint      `json:"id"`
	BusinessID    uint      `json:"business_id"`
	Name          string    `json:"name"`
	Code          string    `json:"code"`
	Address       string    `json:"address"`
	City          string    `json:"city"`
	State         string    `json:"state"`
	Country       string    `json:"country"`
	ZipCode       string    `json:"zip_code"`
	Phone         string    `json:"phone"`
	ContactName   string    `json:"contact_name"`
	ContactEmail  string    `json:"contact_email"`
	IsActive      bool      `json:"is_active"`
	IsDefault     bool      `json:"is_default"`
	IsFulfillment bool      `json:"is_fulfillment"`
	Company       string    `json:"company"`
	FirstName     string    `json:"first_name"`
	LastName      string    `json:"last_name"`
	Email         string    `json:"email"`
	Suburb        string    `json:"suburb"`
	CityDaneCode  string    `json:"city_dane_code"`
	PostalCode    string    `json:"postal_code"`
	Street        string    `json:"street"`
	Latitude      *float64  `json:"latitude"`
	Longitude     *float64  `json:"longitude"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// WarehouseDetailResponse respuesta con ubicaciones
type WarehouseDetailResponse struct {
	WarehouseResponse
	Locations []LocationResponse `json:"locations"`
}

// LocationResponse respuesta de ubicación
type LocationResponse struct {
	ID            uint      `json:"id"`
	WarehouseID   uint      `json:"warehouse_id"`
	Name          string    `json:"name"`
	Code          string    `json:"code"`
	Type          string    `json:"type"`
	IsActive      bool      `json:"is_active"`
	IsFulfillment bool      `json:"is_fulfillment"`
	Capacity      *int      `json:"capacity"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// WarehouseListResponse respuesta paginada
type WarehouseListResponse struct {
	Data       []WarehouseResponse `json:"data"`
	Total      int64               `json:"total"`
	Page       int                 `json:"page"`
	PageSize   int                 `json:"page_size"`
	TotalPages int                 `json:"total_pages"`
}

// FromEntity convierte una entidad a WarehouseResponse
func FromEntity(w *entities.Warehouse) WarehouseResponse {
	return WarehouseResponse{
		ID:            w.ID,
		BusinessID:    w.BusinessID,
		Name:          w.Name,
		Code:          w.Code,
		Address:       w.Address,
		City:          w.City,
		State:         w.State,
		Country:       w.Country,
		ZipCode:       w.ZipCode,
		Phone:         w.Phone,
		ContactName:   w.ContactName,
		ContactEmail:  w.ContactEmail,
		IsActive:      w.IsActive,
		IsDefault:     w.IsDefault,
		IsFulfillment: w.IsFulfillment,
		Company:       w.Company,
		FirstName:     w.FirstName,
		LastName:      w.LastName,
		Email:         w.Email,
		Suburb:        w.Suburb,
		CityDaneCode:  w.CityDaneCode,
		PostalCode:    w.PostalCode,
		Street:        w.Street,
		Latitude:      w.Latitude,
		Longitude:     w.Longitude,
		CreatedAt:     w.CreatedAt,
		UpdatedAt:     w.UpdatedAt,
	}
}

// DetailFromEntity convierte una entidad con ubicaciones a WarehouseDetailResponse
func DetailFromEntity(w *entities.Warehouse) WarehouseDetailResponse {
	locs := make([]LocationResponse, len(w.Locations))
	for i, loc := range w.Locations {
		locs[i] = LocationFromEntity(&loc)
	}
	return WarehouseDetailResponse{
		WarehouseResponse: FromEntity(w),
		Locations:         locs,
	}
}

// LocationFromEntity convierte una entidad de ubicación a LocationResponse
func LocationFromEntity(l *entities.WarehouseLocation) LocationResponse {
	return LocationResponse{
		ID:            l.ID,
		WarehouseID:   l.WarehouseID,
		Name:          l.Name,
		Code:          l.Code,
		Type:          l.Type,
		IsActive:      l.IsActive,
		IsFulfillment: l.IsFulfillment,
		Capacity:      l.Capacity,
		CreatedAt:     l.CreatedAt,
		UpdatedAt:     l.UpdatedAt,
	}
}
