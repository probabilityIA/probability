package entities

import "time"

type JobStatus string

const (
	JobStatusRunning   JobStatus = "running"
	JobStatusCompleted JobStatus = "completed"
	JobStatusFailed    JobStatus = "failed"
)

type JobState struct {
	ID            string
	EventCode     string
	BusinessID    *uint
	Status        JobStatus
	TotalEligible int
	Sent          int
	Skipped       int
	Failed        int
	DryRun        bool
	StartedAt     time.Time
	FinishedAt    *time.Time
	ErrorMessage  string
	CreatedBy     uint
}

type Candidate struct {
	OrderID        string
	OrderNumber    string
	BusinessID     uint
	CustomerPhone  string
	TrackingNumber string
	Status         string
	Carrier        string
	CarrierLogoURL string
}
