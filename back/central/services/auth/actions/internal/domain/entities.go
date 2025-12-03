package domain

import "time"

// Action representa una acción del sistema
type Action struct {
	ID          uint
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}
