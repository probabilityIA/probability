package response

import (
	"time"

	"github.com/secamc93/probability/back/central/services/modules/routes/internal/domain/entities"
)

type RouteStopResponse struct {
	ID               uint       `json:"id"`
	RouteID          uint       `json:"route_id"`
	OrderID          *string    `json:"order_id"`
	Sequence         int        `json:"sequence"`
	Status           string     `json:"status"`
	Address          string     `json:"address"`
	City             string     `json:"city"`
	Lat              *float64   `json:"lat"`
	Lng              *float64   `json:"lng"`
	CustomerName     string     `json:"customer_name"`
	CustomerPhone    string     `json:"customer_phone"`
	EstimatedArrival *time.Time `json:"estimated_arrival"`
	ActualArrival    *time.Time `json:"actual_arrival"`
	ActualDeparture  *time.Time `json:"actual_departure"`
	SignatureURL     string     `json:"signature_url"`
	PhotoURL         string     `json:"photo_url"`
	DeliveryNotes    *string    `json:"delivery_notes"`
	FailureReason    *string    `json:"failure_reason"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

func StopFromEntity(s *entities.RouteStop) RouteStopResponse {
	return RouteStopResponse{
		ID:               s.ID,
		RouteID:          s.RouteID,
		OrderID:          s.OrderID,
		Sequence:         s.Sequence,
		Status:           s.Status,
		Address:          s.Address,
		City:             s.City,
		Lat:              s.Lat,
		Lng:              s.Lng,
		CustomerName:     s.CustomerName,
		CustomerPhone:    s.CustomerPhone,
		EstimatedArrival: s.EstimatedArrival,
		ActualArrival:    s.ActualArrival,
		ActualDeparture:  s.ActualDeparture,
		SignatureURL:     s.SignatureURL,
		PhotoURL:         s.PhotoURL,
		DeliveryNotes:    s.DeliveryNotes,
		FailureReason:    s.FailureReason,
		CreatedAt:        s.CreatedAt,
		UpdatedAt:        s.UpdatedAt,
	}
}
