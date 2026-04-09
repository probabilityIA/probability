package entities

import "time"

type Client struct {
	ID         uint
	BusinessID uint
	Name       string
	Email      *string
	Phone      string
	Dni        *string
	CreatedAt  time.Time
	UpdatedAt  time.Time

	OrderCount  int64
	TotalSpent  float64
	LastOrderAt *time.Time
}
