package app

const (
	StatusPlanned = "planned"
	StatusActive  = "active"
	StatusClosed  = "closed"
)

var validSprintStatuses = map[string]bool{
	StatusPlanned: true,
	StatusActive:  true,
	StatusClosed:  true,
}
