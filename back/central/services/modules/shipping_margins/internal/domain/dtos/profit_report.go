package dtos

import "time"

type ProfitReportParams struct {
	BusinessID uint
	From       *time.Time
	To         *time.Time
	Carrier    string
}

type ProfitReportRow struct {
	Carrier             string  `json:"carrier"`
	CarrierCode         string  `json:"carrier_code"`
	Shipments           int     `json:"shipments"`
	CarrierCostTotal    float64 `json:"carrier_cost_total"`
	CustomerChargeTotal float64 `json:"customer_charge_total"`
	ProfitTotal         float64 `json:"profit_total"`
}

type ProfitReportResponse struct {
	Rows   []ProfitReportRow `json:"rows"`
	Totals ProfitReportRow   `json:"totals"`
}

type ProfitReportDetailParams struct {
	BusinessID uint
	From       *time.Time
	To         *time.Time
	Carrier    string
	Page       int
	PageSize   int
}

type ProfitReportDetailRow struct {
	ShipmentID      uint      `json:"shipment_id"`
	OrderNumber     string    `json:"order_number"`
	TrackingNumber  string    `json:"tracking_number"`
	Carrier         string    `json:"carrier"`
	CustomerCharge  float64   `json:"customer_charge"`
	CarrierCost     float64   `json:"carrier_cost"`
	Profit          float64   `json:"profit"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"created_at"`
}

type ProfitReportDetailResponse struct {
	Data       []ProfitReportDetailRow `json:"data"`
	Total      int64                   `json:"total"`
	Page       int                     `json:"page"`
	PageSize   int                     `json:"page_size"`
	TotalPages int                     `json:"total_pages"`
}
