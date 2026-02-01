package entities

import "time"

// OrderError representa un error ocurrido durante el procesamiento de una orden
// ✅ ENTIDAD PURA - SIN TAGS
type OrderError struct {
	ID              uint
	ExternalID      string
	IntegrationID   uint
	BusinessID      *uint
	IntegrationType string
	Platform        string
	ErrorType       string
	ErrorMessage    string
	ErrorStack      *string
	RawData         []byte
	Status          string
	ResolvedAt      *time.Time
	ResolvedBy      *uint
	Resolution      *string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
