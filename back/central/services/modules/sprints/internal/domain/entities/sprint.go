package entities

import "time"

type Sprint struct {
	ID            uint
	Name          string
	Goal          string
	StartDate     time.Time
	EndDate       time.Time
	Status        string
	CreatedByID   uint
	CreatedByName string
	TicketCount   int64
	DoneCount     int64
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type SprintCounts struct {
	TicketCount int64
	DoneCount   int64
}
