package dtos

import "time"

type BackfillFilter struct {
	EventCode  string
	BusinessID *uint
	Days       int
	Limit      int
}

type OrderCandidateResponse struct {
	OrderID        string `json:"order_id"`
	OrderNumber    string `json:"order_number"`
	CustomerPhone  string `json:"customer_phone"`
	TrackingNumber string `json:"tracking_number,omitempty"`
	Status         string `json:"status"`
	Carrier        string `json:"carrier,omitempty"`
	CarrierLogoURL string `json:"carrier_logo_url,omitempty"`
}

type BusinessGroup struct {
	BusinessID   uint                     `json:"business_id"`
	BusinessName string                   `json:"business_name"`
	Count        int                      `json:"count"`
	Orders       []OrderCandidateResponse `json:"orders"`
}

type PreviewResponse struct {
	EventCode     string          `json:"event_code"`
	TotalEligible int             `json:"total_eligible"`
	Businesses    []BusinessGroup `json:"businesses"`
}

type RunRequest struct {
	EventCode  string `json:"event_code" binding:"required"`
	BusinessID *uint  `json:"business_id"`
	Days       int    `json:"days"`
	Limit      int    `json:"limit"`
}

type RunResponse struct {
	JobID string `json:"job_id"`
}

type JobResponse struct {
	ID            string     `json:"id"`
	EventCode     string     `json:"event_code"`
	BusinessID    *uint      `json:"business_id,omitempty"`
	Status        string     `json:"status"`
	TotalEligible int        `json:"total_eligible"`
	Sent          int        `json:"sent"`
	Skipped       int        `json:"skipped"`
	Failed        int        `json:"failed"`
	StartedAt     time.Time  `json:"started_at"`
	FinishedAt    *time.Time `json:"finished_at,omitempty"`
	ErrorMessage  string     `json:"error_message,omitempty"`
}

type RegisteredEventResponse struct {
	EventCode string `json:"event_code"`
	EventName string `json:"event_name"`
	Channel   string `json:"channel"`
}
