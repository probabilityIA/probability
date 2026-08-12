package entities

import "time"

type SiigoReferral struct {
	ID         uint
	Name       string
	Email      string
	Phone      string
	OrderRange string
	CreatedAt  time.Time
}
