package entities

import "time"

type ShippingMargin struct {
	ID              uint
	BusinessID      uint
	CarrierCode     string
	CarrierName     string
	MarginAmount    float64
	InsuranceMargin float64
	IsActive        bool
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
