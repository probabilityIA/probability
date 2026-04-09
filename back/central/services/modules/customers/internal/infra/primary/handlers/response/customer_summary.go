package response

import (
	"time"

	"github.com/secamc93/probability/back/central/services/modules/customers/internal/domain/entities"
)

type CustomerSummaryResponse struct {
	ID                uint       `json:"id"`
	CustomerID        uint       `json:"customer_id"`
	BusinessID        uint       `json:"business_id"`
	TotalOrders       int        `json:"total_orders"`
	DeliveredOrders   int        `json:"delivered_orders"`
	CancelledOrders   int        `json:"cancelled_orders"`
	InProgressOrders  int        `json:"in_progress_orders"`
	TotalSpent        float64    `json:"total_spent"`
	AvgTicket         float64    `json:"avg_ticket"`
	TotalPaidOrders   int        `json:"total_paid_orders"`
	AvgDeliveryScore  float64    `json:"avg_delivery_score"`
	FirstOrderAt      *time.Time `json:"first_order_at"`
	LastOrderAt       *time.Time `json:"last_order_at"`
	PreferredPlatform string     `json:"preferred_platform"`
	LastUpdatedAt     time.Time  `json:"last_updated_at"`
}

func SummaryFromEntity(s *entities.CustomerSummary) CustomerSummaryResponse {
	return CustomerSummaryResponse{
		ID:                s.ID,
		CustomerID:        s.CustomerID,
		BusinessID:        s.BusinessID,
		TotalOrders:       s.TotalOrders,
		DeliveredOrders:   s.DeliveredOrders,
		CancelledOrders:   s.CancelledOrders,
		InProgressOrders:  s.InProgressOrders,
		TotalSpent:        s.TotalSpent,
		AvgTicket:         s.AvgTicket,
		TotalPaidOrders:   s.TotalPaidOrders,
		AvgDeliveryScore:  s.AvgDeliveryScore,
		FirstOrderAt:      s.FirstOrderAt,
		LastOrderAt:       s.LastOrderAt,
		PreferredPlatform: s.PreferredPlatform,
		LastUpdatedAt:     s.LastUpdatedAt,
	}
}
