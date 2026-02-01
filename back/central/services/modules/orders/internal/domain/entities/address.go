package entities

import "time"

// ProbabilityAddress representa una dirección que se guarda en la base de datos
// ✅ ENTIDAD PURA - SIN TAGS
type ProbabilityAddress struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time

	Type    string
	OrderID string

	FirstName string
	LastName  string
	Company   string
	Phone     string

	Street     string
	Street2    string
	City       string
	State      string
	Country    string
	PostalCode string

	Latitude  *float64
	Longitude *float64

	Instructions *string
	IsDefault    bool
	Metadata     []byte
}
