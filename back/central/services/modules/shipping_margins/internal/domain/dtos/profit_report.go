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
