package handlers

import (
	"time"

	"github.com/secamc93/probability/back/central/services/modules/codreport/internal/domain/entities"
)

type carrierAggregateResponse struct {
	Carrier        string  `json:"carrier"`
	OrdersCount    int     `json:"orders_count"`
	TotalCollected float64 `json:"total_collected"`
	DiscountPct    float64 `json:"discount_pct"`
	TotalDiscount  float64 `json:"total_discount"`
	TotalNet       float64 `json:"total_net"`
}

type monthlyPointResponse struct {
	Month     string  `json:"month"`
	Label     string  `json:"label"`
	Orders    int     `json:"orders"`
	Collected float64 `json:"collected"`
	Discount  float64 `json:"discount"`
	Net       float64 `json:"net"`
}

type carrierDetailResponse struct {
	Carrier         string  `json:"carrier"`
	Orders          int     `json:"orders"`
	EnCurso         float64 `json:"en_curso"`
	EnCursoOrders   int     `json:"en_curso_orders"`
	Entregado       float64 `json:"entregado"`
	EntregadoOrders int     `json:"entregado_orders"`
	PorPagar        float64 `json:"por_pagar"`
	Recaudado       float64 `json:"recaudado"`
	Cargo           float64 `json:"cargo"`
	Total           float64 `json:"total"`
}

type historyPointResponse struct {
	Label     string  `json:"label"`
	Entregado float64 `json:"entregado"`
	EnCurso   float64 `json:"en_curso"`
}

type summaryResponse struct {
	TotalCollected  float64                    `json:"total_collected"`
	TotalPending    float64                    `json:"total_pending"`
	TotalDiscount   float64                    `json:"total_discount"`
	TotalNet        float64                    `json:"total_net"`
	OrdersCollected int                        `json:"orders_collected"`
	OrdersPending   int                        `json:"orders_pending"`
	ByCarrier       []carrierAggregateResponse `json:"by_carrier"`
	Monthly         []monthlyPointResponse     `json:"monthly"`
	EnCursoTotal    float64                    `json:"en_curso_total"`
	EnCursoOrders   int                        `json:"en_curso_orders"`
	EntregadoTotal  float64                    `json:"entregado_total"`
	EntregadoOrders int                        `json:"entregado_orders"`
	CarrierDetail   []carrierDetailResponse    `json:"carrier_detail"`
	History         []historyPointResponse     `json:"history"`
}

type codOrderResponse struct {
	OrderID       string     `json:"order_id"`
	OrderNumber   string     `json:"order_number"`
	ShipmentID    uint       `json:"shipment_id"`
	HasGuide      bool       `json:"has_guide"`
	CustomerName  string     `json:"customer_name"`
	Carrier       string     `json:"carrier"`
	CodTotal      float64    `json:"cod_total"`
	CodCarrierFee float64    `json:"cod_carrier_fee"`
	ShippingCost  float64    `json:"shipping_cost"`
	DiscountPct   float64    `json:"discount_pct"`
	Discount      float64    `json:"discount"`
	Net           float64    `json:"net"`
	Currency      string     `json:"currency"`
	Status        string     `json:"status"`
	Collected     bool       `json:"collected"`
	Paid          bool       `json:"paid"`
	CodState      string     `json:"cod_state"`
	CutStatus     string     `json:"cut_status"`
	CreatedAt     time.Time  `json:"created_at"`
	DeliveredAt   *time.Time `json:"delivered_at"`
}

type paymentCutResponse struct {
	ID              uint                       `json:"id"`
	PeriodStart     time.Time                  `json:"period_start"`
	PeriodEnd       time.Time                  `json:"period_end"`
	Status          string                     `json:"status"`
	OrdersCount     int                        `json:"orders_count"`
	TotalCollected  float64                    `json:"total_collected"`
	TotalDiscount   float64                    `json:"total_discount"`
	TotalNet        float64                    `json:"total_net"`
	ByCarrier       []carrierAggregateResponse `json:"by_carrier"`
	ConfirmedBy       uint                     `json:"confirmed_by"`
	ConfirmedByName   string                   `json:"confirmed_by_name"`
	ConfirmedByAvatar string                   `json:"confirmed_by_avatar"`
	ConfirmedAt       *time.Time               `json:"confirmed_at"`
}

type carrierConfigResponse struct {
	ID                 uint    `json:"id"`
	CarrierName        string  `json:"carrier_name"`
	DiscountPercentage float64 `json:"discount_percentage"`
	IsActive           bool    `json:"is_active"`
}

func mapCarrierAggregates(in []entities.CarrierAggregate) []carrierAggregateResponse {
	out := make([]carrierAggregateResponse, len(in))
	for i := range in {
		out[i] = carrierAggregateResponse{
			Carrier:        in[i].Carrier,
			OrdersCount:    in[i].OrdersCount,
			TotalCollected: in[i].TotalCollected,
			DiscountPct:    in[i].DiscountPct,
			TotalDiscount:  in[i].TotalDiscount,
			TotalNet:       in[i].TotalNet,
		}
	}
	return out
}

func mapSummary(s *entities.CodSummary) summaryResponse {
	monthly := make([]monthlyPointResponse, len(s.Monthly))
	for i := range s.Monthly {
		monthly[i] = monthlyPointResponse{
			Month:     s.Monthly[i].Month,
			Label:     s.Monthly[i].Label,
			Orders:    s.Monthly[i].Orders,
			Collected: s.Monthly[i].Collected,
			Discount:  s.Monthly[i].Discount,
			Net:       s.Monthly[i].Net,
		}
	}
	detail := make([]carrierDetailResponse, len(s.CarrierDetail))
	for i := range s.CarrierDetail {
		d := s.CarrierDetail[i]
		detail[i] = carrierDetailResponse{
			Carrier:         d.Carrier,
			Orders:          d.Orders,
			EnCurso:         d.EnCurso,
			EnCursoOrders:   d.EnCursoOrders,
			Entregado:       d.Entregado,
			EntregadoOrders: d.EntregadoOrders,
			PorPagar:        d.PorPagar,
			Recaudado:       d.Recaudado,
			Cargo:           d.Cargo,
			Total:           d.Total,
		}
	}
	history := make([]historyPointResponse, len(s.History))
	for i := range s.History {
		history[i] = historyPointResponse{
			Label:     s.History[i].Label,
			Entregado: s.History[i].Entregado,
			EnCurso:   s.History[i].EnCurso,
		}
	}
	return summaryResponse{
		TotalCollected:  s.TotalCollected,
		TotalPending:    s.TotalPending,
		TotalDiscount:   s.TotalDiscount,
		TotalNet:        s.TotalNet,
		OrdersCollected: s.OrdersCollected,
		OrdersPending:   s.OrdersPending,
		ByCarrier:       mapCarrierAggregates(s.ByCarrier),
		Monthly:         monthly,
		EnCursoTotal:    s.EnCursoTotal,
		EnCursoOrders:   s.EnCursoOrders,
		EntregadoTotal:  s.EntregadoTotal,
		EntregadoOrders: s.EntregadoOrders,
		CarrierDetail:   detail,
		History:         history,
	}
}

func mapOrders(in []entities.CodOrder) []codOrderResponse {
	out := make([]codOrderResponse, len(in))
	for i := range in {
		out[i] = codOrderResponse{
			OrderID:       in[i].OrderID,
			OrderNumber:   in[i].OrderNumber,
			ShipmentID:    in[i].ShipmentID,
			HasGuide:      in[i].HasGuide,
			CustomerName:  in[i].CustomerName,
			Carrier:       in[i].Carrier,
			CodTotal:      in[i].CodTotal,
			CodCarrierFee: in[i].CodCarrierFee,
			ShippingCost:  in[i].ShippingCost,
			DiscountPct:   in[i].DiscountPct,
			Discount:      in[i].Discount,
			Net:           in[i].Net,
			Currency:      in[i].Currency,
			Status:        in[i].Status,
			Collected:     in[i].Collected,
			Paid:          in[i].Paid,
			CodState:      in[i].CodState,
			CutStatus:     in[i].CutStatus,
			CreatedAt:     in[i].CreatedAt,
			DeliveredAt:   in[i].DeliveredAt,
		}
	}
	return out
}

func mapCut(c *entities.PaymentCut) paymentCutResponse {
	return paymentCutResponse{
		ID:              c.ID,
		PeriodStart:     c.PeriodStart,
		PeriodEnd:       c.PeriodEnd,
		Status:          c.Status,
		OrdersCount:     c.OrdersCount,
		TotalCollected:  c.TotalCollected,
		TotalDiscount:   c.TotalDiscount,
		TotalNet:        c.TotalNet,
		ByCarrier:         mapCarrierAggregates(c.ByCarrier),
		ConfirmedBy:       c.ConfirmedBy,
		ConfirmedByName:   c.ConfirmedByName,
		ConfirmedByAvatar: c.ConfirmedByAvatar,
		ConfirmedAt:       c.ConfirmedAt,
	}
}

func mapCuts(in []entities.PaymentCut) []paymentCutResponse {
	out := make([]paymentCutResponse, len(in))
	for i := range in {
		out[i] = mapCut(&in[i])
	}
	return out
}

func mapCarrierConfigs(in []entities.CarrierConfig) []carrierConfigResponse {
	out := make([]carrierConfigResponse, len(in))
	for i := range in {
		out[i] = carrierConfigResponse{
			ID:                 in[i].ID,
			CarrierName:        in[i].CarrierName,
			DiscountPercentage: in[i].DiscountPercentage,
			IsActive:           in[i].IsActive,
		}
	}
	return out
}
