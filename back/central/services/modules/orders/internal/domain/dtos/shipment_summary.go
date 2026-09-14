package dtos

type ShipmentSummary struct {
	ID                  uint
	Carrier             *string
	TrackingNumber      *string
	GuideURL            *string
	Status              string
	CarrierStatus       *string
	CarrierStatusDetail *string
	TotalCost           *float64
	CodCarrierFee       *float64
	CodCollectAmount    *float64
}
