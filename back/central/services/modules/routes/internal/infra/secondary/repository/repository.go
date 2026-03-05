package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/routes/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/routes/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/routes/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/routes/internal/domain/ports"
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

// ============================================
// Route CRUD
// ============================================

func (r *Repository) CreateRoute(ctx context.Context, route *entities.Route, stops []entities.RouteStop) (*entities.Route, error) {
	model := entityToRouteModel(route)

	err := r.db.Conn(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(model).Error; err != nil {
			return err
		}

		for i := range stops {
			stopModel := entityToStopModel(&stops[i])
			stopModel.RouteID = model.ID
			if err := tx.Create(stopModel).Error; err != nil {
				return err
			}
			stops[i].ID = stopModel.ID
			stops[i].RouteID = model.ID
			stops[i].CreatedAt = stopModel.CreatedAt
			stops[i].UpdatedAt = stopModel.UpdatedAt
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	route.ID = model.ID
	route.CreatedAt = model.CreatedAt
	route.UpdatedAt = model.UpdatedAt
	route.Stops = stops

	return route, nil
}

func (r *Repository) GetRouteByID(ctx context.Context, businessID, routeID uint) (*entities.Route, error) {
	var model models.Route
	err := r.db.Conn(ctx).
		Preload("Stops", func(db *gorm.DB) *gorm.DB {
			return db.Order("sequence ASC")
		}).
		Where("id = ? AND business_id = ?", routeID, businessID).
		First(&model).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domainerrors.ErrRouteNotFound
		}
		return nil, err
	}

	route := routeModelToEntity(&model)

	// Denormalize driver name and vehicle plate
	if model.DriverID != nil {
		name, err := r.GetDriverNameByID(ctx, *model.DriverID)
		if err == nil {
			route.DriverName = name
		}
	}
	if model.VehicleID != nil {
		plate, err := r.GetVehiclePlateByID(ctx, *model.VehicleID)
		if err == nil {
			route.VehiclePlate = plate
		}
	}

	return route, nil
}

func (r *Repository) ListRoutes(ctx context.Context, params dtos.ListRoutesParams) ([]entities.Route, int64, error) {
	var modelsList []models.Route
	var total int64

	query := r.db.Conn(ctx).Model(&models.Route{}).
		Where("business_id = ?", params.BusinessID)

	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}

	if params.DriverID != nil {
		query = query.Where("driver_id = ?", *params.DriverID)
	}

	if params.DateFrom != nil {
		query = query.Where("date >= ?", *params.DateFrom)
	}

	if params.DateTo != nil {
		query = query.Where("date <= ?", *params.DateTo)
	}

	if params.Search != "" {
		like := "%" + params.Search + "%"
		query = query.Where("origin_address ILIKE ? OR notes ILIKE ?", like, like)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(params.Offset()).Limit(params.PageSize).
		Order("date DESC, created_at DESC").
		Find(&modelsList).Error; err != nil {
		return nil, 0, err
	}

	routes := make([]entities.Route, len(modelsList))
	for i, m := range modelsList {
		routes[i] = *routeModelToEntity(&m)
		// Denormalize names for list view
		if m.DriverID != nil {
			name, err := r.GetDriverNameByID(ctx, *m.DriverID)
			if err == nil {
				routes[i].DriverName = name
			}
		}
		if m.VehicleID != nil {
			plate, err := r.GetVehiclePlateByID(ctx, *m.VehicleID)
			if err == nil {
				routes[i].VehiclePlate = plate
			}
		}
	}

	return routes, total, nil
}

func (r *Repository) UpdateRoute(ctx context.Context, route *entities.Route) (*entities.Route, error) {
	model := entityToRouteModel(route)
	model.ID = route.ID

	if err := r.db.Conn(ctx).Save(model).Error; err != nil {
		return nil, err
	}

	route.UpdatedAt = model.UpdatedAt
	return route, nil
}

func (r *Repository) DeleteRoute(ctx context.Context, businessID, routeID uint) error {
	result := r.db.Conn(ctx).
		Where("id = ? AND business_id = ?", routeID, businessID).
		Delete(&models.Route{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domainerrors.ErrRouteNotFound
	}
	return nil
}

// ============================================
// Route lifecycle
// ============================================

func (r *Repository) UpdateRouteStatus(ctx context.Context, routeID uint, status string) error {
	return r.db.Conn(ctx).Model(&models.Route{}).
		Where("id = ?", routeID).
		Update("status", status).Error
}

func (r *Repository) UpdateRouteCounters(ctx context.Context, routeID uint) error {
	return r.db.Conn(ctx).Exec(`
		UPDATE routes SET
			total_stops = (SELECT COUNT(*) FROM route_stops WHERE route_id = ? AND deleted_at IS NULL),
			completed_stops = (SELECT COUNT(*) FROM route_stops WHERE route_id = ? AND deleted_at IS NULL AND status = 'delivered'),
			failed_stops = (SELECT COUNT(*) FROM route_stops WHERE route_id = ? AND deleted_at IS NULL AND status = 'failed')
		WHERE id = ?
	`, routeID, routeID, routeID, routeID).Error
}

func (r *Repository) SetRouteActualStart(ctx context.Context, routeID uint) error {
	now := time.Now()
	return r.db.Conn(ctx).Model(&models.Route{}).
		Where("id = ?", routeID).
		Update("actual_start_time", now).Error
}

func (r *Repository) SetRouteActualEnd(ctx context.Context, routeID uint) error {
	now := time.Now()
	return r.db.Conn(ctx).Model(&models.Route{}).
		Where("id = ?", routeID).
		Update("actual_end_time", now).Error
}

// ============================================
// Stop CRUD
// ============================================

func (r *Repository) AddStop(ctx context.Context, stop *entities.RouteStop) (*entities.RouteStop, error) {
	model := entityToStopModel(stop)

	if err := r.db.Conn(ctx).Create(model).Error; err != nil {
		return nil, err
	}

	stop.ID = model.ID
	stop.CreatedAt = model.CreatedAt
	stop.UpdatedAt = model.UpdatedAt
	return stop, nil
}

func (r *Repository) UpdateStop(ctx context.Context, stop *entities.RouteStop) (*entities.RouteStop, error) {
	model := entityToStopModel(stop)
	model.ID = stop.ID

	if err := r.db.Conn(ctx).Save(model).Error; err != nil {
		return nil, err
	}

	stop.UpdatedAt = model.UpdatedAt
	return stop, nil
}

func (r *Repository) DeleteStop(ctx context.Context, routeID, stopID uint) error {
	result := r.db.Conn(ctx).
		Where("id = ? AND route_id = ?", stopID, routeID).
		Delete(&models.RouteStop{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domainerrors.ErrStopNotFound
	}
	return nil
}

func (r *Repository) GetStopByID(ctx context.Context, routeID, stopID uint) (*entities.RouteStop, error) {
	var model models.RouteStop
	err := r.db.Conn(ctx).
		Where("id = ? AND route_id = ?", stopID, routeID).
		First(&model).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domainerrors.ErrStopNotFound
		}
		return nil, err
	}
	return stopModelToEntity(&model), nil
}

func (r *Repository) UpdateStopStatus(ctx context.Context, stopID uint, status string, failureReason *string, signatureURL, photoURL string) error {
	updates := map[string]interface{}{
		"status": status,
	}

	if status == "arrived" || status == "delivered" {
		now := time.Now()
		if status == "arrived" {
			updates["actual_arrival"] = now
		} else {
			updates["actual_departure"] = now
		}
	}

	if failureReason != nil {
		updates["failure_reason"] = *failureReason
	}
	if signatureURL != "" {
		updates["signature_url"] = signatureURL
	}
	if photoURL != "" {
		updates["photo_url"] = photoURL
	}

	return r.db.Conn(ctx).Model(&models.RouteStop{}).
		Where("id = ?", stopID).
		Updates(updates).Error
}

func (r *Repository) ReorderStops(ctx context.Context, routeID uint, stopIDs []uint) error {
	return r.db.Conn(ctx).Transaction(func(tx *gorm.DB) error {
		for i, id := range stopIDs {
			if err := tx.Model(&models.RouteStop{}).
				Where("id = ? AND route_id = ?", id, routeID).
				Update("sequence", i+1).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *Repository) GetStopsByRouteID(ctx context.Context, routeID uint) ([]entities.RouteStop, error) {
	var modelsList []models.RouteStop
	err := r.db.Conn(ctx).
		Where("route_id = ?", routeID).
		Order("sequence ASC").
		Find(&modelsList).Error
	if err != nil {
		return nil, err
	}

	stops := make([]entities.RouteStop, len(modelsList))
	for i, m := range modelsList {
		stops[i] = *stopModelToEntity(&m)
	}
	return stops, nil
}

func (r *Repository) SetPendingStopsStatus(ctx context.Context, routeID uint, status string) error {
	return r.db.Conn(ctx).Model(&models.RouteStop{}).
		Where("route_id = ? AND status = 'pending'", routeID).
		Update("status", status).Error
}

// ============================================
// Cross-module queries (replicated locally)
// ============================================

func (r *Repository) GetDriverNameByID(ctx context.Context, driverID uint) (string, error) {
	var result struct {
		FirstName string
		LastName  string
	}
	err := r.db.Conn(ctx).
		Model(&models.Driver{}).
		Select("first_name, last_name").
		Where("id = ? AND deleted_at IS NULL", driverID).
		First(&result).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", domainerrors.ErrDriverNotFound
		}
		return "", err
	}
	return fmt.Sprintf("%s %s", result.FirstName, result.LastName), nil
}

func (r *Repository) UpdateDriverStatus(ctx context.Context, driverID uint, status string) error {
	return r.db.Conn(ctx).Model(&models.Driver{}).
		Where("id = ?", driverID).
		Update("status", status).Error
}

func (r *Repository) GetVehiclePlateByID(ctx context.Context, vehicleID uint) (string, error) {
	var plate string
	err := r.db.Conn(ctx).
		Model(&models.Vehicle{}).
		Select("license_plate").
		Where("id = ? AND deleted_at IS NULL", vehicleID).
		Scan(&plate).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", domainerrors.ErrVehicleNotFound
		}
		return "", err
	}
	return plate, nil
}

func (r *Repository) UpdateOrderDriverInfo(ctx context.Context, orderID string, driverID *uint, driverName string, isLastMile bool) error {
	return r.db.Conn(ctx).Model(&models.Order{}).
		Where("id = ?", orderID).
		Updates(map[string]interface{}{
			"driver_id":    driverID,
			"driver_name":  driverName,
			"is_last_mile": isLastMile,
		}).Error
}

func (r *Repository) ClearOrderDriverInfo(ctx context.Context, orderID string) error {
	return r.db.Conn(ctx).Model(&models.Order{}).
		Where("id = ?", orderID).
		Updates(map[string]interface{}{
			"driver_id":    nil,
			"driver_name":  "",
			"is_last_mile": false,
		}).Error
}

// ============================================
// Mappers
// ============================================

func entityToRouteModel(e *entities.Route) *models.Route {
	return &models.Route{
		BusinessID:         e.BusinessID,
		DriverID:           e.DriverID,
		VehicleID:          e.VehicleID,
		Status:             e.Status,
		Date:               e.Date,
		StartTime:          e.StartTime,
		EndTime:            e.EndTime,
		ActualStartTime:    e.ActualStartTime,
		ActualEndTime:      e.ActualEndTime,
		OriginWarehouseID:  e.OriginWarehouseID,
		OriginAddress:      e.OriginAddress,
		OriginLat:          e.OriginLat,
		OriginLng:          e.OriginLng,
		TotalStops:         e.TotalStops,
		CompletedStops:     e.CompletedStops,
		FailedStops:        e.FailedStops,
		TotalDistanceKm:    e.TotalDistanceKm,
		TotalDurationMin:   e.TotalDurationMin,
		Notes:              e.Notes,
	}
}

func routeModelToEntity(m *models.Route) *entities.Route {
	route := &entities.Route{
		ID:                m.ID,
		BusinessID:        m.BusinessID,
		DriverID:          m.DriverID,
		VehicleID:         m.VehicleID,
		Status:            m.Status,
		Date:              m.Date,
		StartTime:         m.StartTime,
		EndTime:           m.EndTime,
		ActualStartTime:   m.ActualStartTime,
		ActualEndTime:     m.ActualEndTime,
		OriginWarehouseID: m.OriginWarehouseID,
		OriginAddress:     m.OriginAddress,
		OriginLat:         m.OriginLat,
		OriginLng:         m.OriginLng,
		TotalStops:        m.TotalStops,
		CompletedStops:    m.CompletedStops,
		FailedStops:       m.FailedStops,
		TotalDistanceKm:   m.TotalDistanceKm,
		TotalDurationMin:  m.TotalDurationMin,
		Notes:             m.Notes,
		CreatedAt:         m.CreatedAt,
		UpdatedAt:         m.UpdatedAt,
	}

	if len(m.Stops) > 0 {
		route.Stops = make([]entities.RouteStop, len(m.Stops))
		for i, s := range m.Stops {
			route.Stops[i] = *stopModelToEntity(&s)
		}
	}

	return route
}

func entityToStopModel(e *entities.RouteStop) *models.RouteStop {
	return &models.RouteStop{
		RouteID:          e.RouteID,
		OrderID:          e.OrderID,
		Sequence:         e.Sequence,
		Status:           e.Status,
		Address:          e.Address,
		City:             e.City,
		Lat:              e.Lat,
		Lng:              e.Lng,
		CustomerName:     e.CustomerName,
		CustomerPhone:    e.CustomerPhone,
		EstimatedArrival: e.EstimatedArrival,
		ActualArrival:    e.ActualArrival,
		ActualDeparture:  e.ActualDeparture,
		SignatureURL:     e.SignatureURL,
		PhotoURL:         e.PhotoURL,
		DeliveryNotes:    e.DeliveryNotes,
		FailureReason:    e.FailureReason,
	}
}

func stopModelToEntity(m *models.RouteStop) *entities.RouteStop {
	return &entities.RouteStop{
		ID:               m.ID,
		RouteID:          m.RouteID,
		OrderID:          m.OrderID,
		Sequence:         m.Sequence,
		Status:           m.Status,
		Address:          m.Address,
		City:             m.City,
		Lat:              m.Lat,
		Lng:              m.Lng,
		CustomerName:     m.CustomerName,
		CustomerPhone:    m.CustomerPhone,
		EstimatedArrival: m.EstimatedArrival,
		ActualArrival:    m.ActualArrival,
		ActualDeparture:  m.ActualDeparture,
		SignatureURL:     m.SignatureURL,
		PhotoURL:         m.PhotoURL,
		DeliveryNotes:    m.DeliveryNotes,
		FailureReason:    m.FailureReason,
		CreatedAt:        m.CreatedAt,
		UpdatedAt:        m.UpdatedAt,
	}
}
