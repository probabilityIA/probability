package entities

import "time"

type TourProgress struct {
	ID          uint
	UserID      uint
	BusinessID  uint
	TourKey     string
	Version     int
	Status      string
	StepIndex   int
	CompletedAt *time.Time
	UpdatedAt   time.Time
}

const (
	StatusPending    = "pending"
	StatusInProgress = "in_progress"
	StatusCompleted  = "completed"
	StatusSkipped    = "skipped"
)

func IsValidStatus(status string) bool {
	switch status {
	case StatusPending, StatusInProgress, StatusCompleted, StatusSkipped:
		return true
	}
	return false
}
