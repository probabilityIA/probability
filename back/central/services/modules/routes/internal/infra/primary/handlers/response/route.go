package response

import (
	"time"

	"github.com/secamc93/probability/back/central/services/modules/routes/internal/domain/entities"
)

type RouteResponse struct {
	ID                uint       `json:"id"`
	BusinessID        uint       `json:"business_id"`
	DriverID          *uint      `json:"driver_id"`
	VehicleID         *uint      `json:"vehicle_id"`
	DriverName        string     `json:"driver_name"`
	VehiclePlate      string     `json:"vehicle_plate"`
	Status            string     `json:"status"`
	Date              time.Time  `json:"date"`
	StartTime         *time.Time `json:"start_time"`
	EndTime           *time.Time `json:"end_time"`
	ActualStartTime   *time.Time `json:"actual_start_time"`
	ActualEndTime     *time.Time `json:"actual_end_time"`
	OriginWarehouseID *uint      `json:"origin_warehouse_id"`
	OriginAddress     string     `json:"origin_address"`
	OriginLat         *float64   `json:"origin_lat"`
	OriginLng         *float64   `json:"origin_lng"`
	TotalStops        int        `json:"total_stops"`
	CompletedStops    int        `json:"completed_stops"`
	FailedStops       int        `json:"failed_stops"`
	TotalDistanceKm   *float64   `json:"total_distance_km"`
	TotalDurationMin  *int       `json:"total_duration_min"`
	Notes             *string    `json:"notes"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

type RouteDetailResponse struct {
	RouteResponse
	Stops []RouteStopResponse `json:"stops"`
}

type RoutesListResponse struct {
	Data       []RouteResponse `json:"data"`
	Total      int64           `json:"total"`
	Page       int             `json:"page"`
	PageSize   int             `json:"page_size"`
	TotalPages int             `json:"total_pages"`
}

func FromEntity(r *entities.Route) RouteResponse {
	return RouteResponse{
		ID:                r.ID,
		BusinessID:        r.BusinessID,
		DriverID:          r.DriverID,
		VehicleID:         r.VehicleID,
		DriverName:        r.DriverName,
		VehiclePlate:      r.VehiclePlate,
		Status:            r.Status,
		Date:              r.Date,
		StartTime:         r.StartTime,
		EndTime:           r.EndTime,
		ActualStartTime:   r.ActualStartTime,
		ActualEndTime:     r.ActualEndTime,
		OriginWarehouseID: r.OriginWarehouseID,
		OriginAddress:     r.OriginAddress,
		OriginLat:         r.OriginLat,
		OriginLng:         r.OriginLng,
		TotalStops:        r.TotalStops,
		CompletedStops:    r.CompletedStops,
		FailedStops:       r.FailedStops,
		TotalDistanceKm:   r.TotalDistanceKm,
		TotalDurationMin:  r.TotalDurationMin,
		Notes:             r.Notes,
		CreatedAt:         r.CreatedAt,
		UpdatedAt:         r.UpdatedAt,
	}
}

func DetailFromEntity(r *entities.Route) RouteDetailResponse {
	stops := make([]RouteStopResponse, len(r.Stops))
	for i, s := range r.Stops {
		stops[i] = StopFromEntity(&s)
	}
	return RouteDetailResponse{
		RouteResponse: FromEntity(r),
		Stops:         stops,
	}
}
