package entities

import "time"

type Flow struct {
	ID             uint
	BusinessID     uint
	Name           string
	Description    string
	RootTemplateID *uint
	Enabled        bool
	CreatedAt      time.Time
	UpdatedAt      time.Time

	RootTemplateName   string
	RootTemplateStatus string
	StepCount          int
	PendingCount       int
}
