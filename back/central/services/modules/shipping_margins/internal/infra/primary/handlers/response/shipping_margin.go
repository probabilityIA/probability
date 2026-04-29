package response

import (
	"time"

	"github.com/secamc93/probability/back/central/services/modules/shipping_margins/internal/domain/entities"
)

type ShippingMarginResponse struct {
	ID              uint      `json:"id"`
	BusinessID      uint      `json:"business_id"`
	CarrierCode     string    `json:"carrier_code"`
	CarrierName     string    `json:"carrier_name"`
	MarginAmount    float64   `json:"margin_amount"`
	InsuranceMargin float64   `json:"insurance_margin"`
	IsActive        bool      `json:"is_active"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type ShippingMarginsListResponse struct {
	Data       []ShippingMarginResponse `json:"data"`
	Total      int64                    `json:"total"`
	Page       int                      `json:"page"`
	PageSize   int                      `json:"page_size"`
	TotalPages int                      `json:"total_pages"`
}

func FromEntity(m *entities.ShippingMargin) ShippingMarginResponse {
	return ShippingMarginResponse{
		ID:              m.ID,
		BusinessID:      m.BusinessID,
		CarrierCode:     m.CarrierCode,
		CarrierName:     m.CarrierName,
		MarginAmount:    m.MarginAmount,
		InsuranceMargin: m.InsuranceMargin,
		IsActive:        m.IsActive,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
	}
}
