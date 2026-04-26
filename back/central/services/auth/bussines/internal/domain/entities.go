package domain

import "time"

type BusinessType struct {
	ID          uint
	Name        string
	Code        string
	Description string
	Icon        string
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}

// Business representa un negocio en el sistema
type Business struct {
	ID             uint
	Name           string
	Code           string
	OrderPrefix    string
	BusinessTypeID uint
	BusinessType   *BusinessType // Relación con tipo de business
	Timezone       string
	Address        string
	Description    string

	// Configuración de marca blanca
	LogoURL         string
	PrimaryColor    string
	SecondaryColor  string
	TertiaryColor   string
	QuaternaryColor string
	NavbarImageURL  string
	CustomDomain    string
	IsActive        bool

	// Configuración de funcionalidades
	EnableDelivery     bool
	EnablePickup       bool
	EnableReservations bool

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

// BusinessResourceConfigured representa la configuración de recursos para un negocio
type BusinessResourceConfigured struct {
	ResourceID   uint
	ResourceName string
	IsActive     bool
}

// BusinessTypeResourcePermitted representa un recurso permitido para un tipo de negocio
type BusinessTypeResourcePermitted struct {
	ID             uint
	BusinessTypeID uint
	ResourceID     uint
	ResourceName   string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// Resource representa un recurso del sistema
type Resource struct {
	ID        uint
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
