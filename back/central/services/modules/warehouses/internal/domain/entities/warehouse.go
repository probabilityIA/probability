package entities

import "time"

// Warehouse representa una bodega en el dominio
type Warehouse struct {
	ID            uint
	BusinessID    uint
	Name          string
	Code          string
	Address       string
	City          string
	State         string
	Country       string
	ZipCode       string
	Phone         string
	ContactName   string
	ContactEmail  string
	IsActive      bool
	IsDefault     bool
	IsFulfillment bool
	StructureType string
	ZoneCount     int
	AisleCount    int
	RackCount     int
	LevelCount    int
	PositionCount int
	Company       string
	FirstName     string
	LastName      string
	Email         string
	Suburb        string
	CityDaneCode  string
	PostalCode    string
	Street        string
	Latitude      *float64
	Longitude     *float64

	ShippingPackageStrategy string
	StandardBoxes           []StandardBox

	Locations []WarehouseLocation
	CreatedAt time.Time
	UpdatedAt time.Time
}

// StandardBox es un tamano de caja fijo que la bodega puede usar para cotizar
// envios, con una capacidad maxima de unidades para elegir cual caja aplica.
type StandardBox struct {
	Name     string
	Weight   *float64
	Length   *float64
	Width    *float64
	Height   *float64
	MaxItems int
}
