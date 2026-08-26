package dtos

import "time"

type ReportFilter struct {
	BusinessID uint
	StartDate  time.Time
	EndDate    time.Time
	Carrier    string
	Bucket     string
}

type OrdersFilter struct {
	BusinessID uint
	StartDate  time.Time
	EndDate    time.Time
	Carrier    string
	Status     string
	Collected  *bool
	HasGuide   *bool
	Search     string
	Guides     []string
	Page       int
	PageSize   int
}

type SaveCarrierConfigDTO struct {
	BusinessID         uint
	CarrierName        string
	DiscountPercentage float64
	IsActive           bool
}

type ConfirmCutDTO struct {
	BusinessID  uint
	PeriodStart time.Time
	PeriodEnd   time.Time
	OrderIDs    []string
	UserID      uint
	UserName    string
}

type SelectableOrdersFilter struct {
	BusinessID  uint
	PeriodStart time.Time
	PeriodEnd   time.Time
}

type UpdateCarrierFeeDTO struct {
	BusinessID uint
	ShipmentID uint
	Fee        float64
	UserID     uint
	UserName   string
}

type UpdateCarrierFeeResult struct {
	OrderNumber string
	PreviousFee float64
	NewFee      float64
}

type SendCutEmailDTO struct {
	BusinessID uint
	CutID      uint
	Recipients []string
	UserID     uint
}
