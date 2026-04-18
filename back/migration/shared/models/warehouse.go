package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type WarehouseLocationFlags struct {
	IsPicking    bool `json:"is_picking"`
	IsBulk       bool `json:"is_bulk"`
	IsQuarantine bool `json:"is_quarantine"`
	IsDamaged    bool `json:"is_damaged"`
	IsReturns    bool `json:"is_returns"`
	IsCrossDock  bool `json:"is_cross_dock"`
	IsHazmat     bool `json:"is_hazmat"`
}

// Warehouse representa una bodega o almacén del negocio
type Warehouse struct {
	gorm.Model
	BusinessID    uint           `gorm:"not null;index;uniqueIndex:idx_warehouse_business_code,priority:1"` // ID del negocio al que pertenece la bodega
	Name          string         `gorm:"size:255;not null"`                                                 // Nombre de la bodega (ej: "Bodega Principal")
	Code          string         `gorm:"size:50;not null;uniqueIndex:idx_warehouse_business_code,priority:2"` // Código único dentro del negocio (ej: "BOD-001")
	Address       string         `gorm:"size:500"`                                                          // Dirección física de la bodega
	City          string         `gorm:"size:100"`                                                          // Ciudad donde está ubicada la bodega
	State         string         `gorm:"size:100"`                                                          // Departamento o estado donde está la bodega
	Country       string         `gorm:"size:50;default:'CO'"`                                              // País donde está la bodega (ISO 3166-1 alpha-2)
	ZipCode       string         `gorm:"size:20"`                                                           // Código postal de la bodega
	Phone         string         `gorm:"size:50"`                                                           // Teléfono de contacto de la bodega
	ContactName   string         `gorm:"size:255"`                                                          // Nombre del responsable o encargado de la bodega
	ContactEmail  string         `gorm:"size:255"`                                                          // Email del responsable de la bodega
	IsActive      bool           `gorm:"default:true;index"`                                                // Indica si la bodega está activa y operativa
	IsDefault     bool           `gorm:"default:false;index"`                                               // Indica si es la bodega principal/por defecto del negocio
	IsFulfillment bool           `gorm:"default:false;index"`                                               // Indica si la bodega maneja fulfillment (envíos directos)
	Metadata      datatypes.JSON `gorm:"type:jsonb"`                                                        // Metadatos adicionales en formato JSON (configuración personalizada)

	// Campos de contacto para carrier (requeridos por APIs de transportadoras)
	Company      string `gorm:"size:100"`  // Nombre de la empresa remitente
	FirstName    string `gorm:"size:100"`  // Nombre del contacto
	LastName     string `gorm:"size:100"`  // Apellido del contacto
	Email        string `gorm:"size:100"`  // Email del contacto (carrier)
	Suburb       string `gorm:"size:100"`  // Barrio / Colonia
	CityDaneCode string `gorm:"size:10"`   // Código DANE de la ciudad
	PostalCode   string `gorm:"size:20"`   // Código postal
	Street       string `gorm:"size:255"`  // Dirección exacta de la calle (carrier format)
	Latitude     *float64                  // Latitud GPS de la bodega
	Longitude    *float64                  // Longitud GPS de la bodega

	// Relaciones
	Business  Business            `gorm:"foreignKey:BusinessID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"` // Negocio al que pertenece
	Locations []WarehouseLocation `gorm:"foreignKey:WarehouseID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"` // Ubicaciones dentro de la bodega
}

// TableName especifica el nombre de la tabla
func (Warehouse) TableName() string {
	return "warehouses"
}

type WarehouseLocation struct {
	gorm.Model
	WarehouseID   uint   `gorm:"not null;index;uniqueIndex:idx_location_warehouse_code,priority:1"`
	LevelID       *uint  `gorm:"index"`
	Name          string `gorm:"size:255;not null"`
	Code          string `gorm:"size:50;not null;uniqueIndex:idx_location_warehouse_code,priority:2"`
	Type          string `gorm:"size:50;default:'storage'"`
	IsActive      bool   `gorm:"default:true;index"`
	IsFulfillment bool   `gorm:"default:false;index"`
	Capacity      *int

	MaxWeightKg  *float64
	MaxVolumeCm3 *float64
	LengthCm     *float64
	WidthCm      *float64
	HeightCm     *float64
	Priority     int             `gorm:"default:0;index"`
	Flags        datatypes.JSON  `gorm:"type:jsonb"`

	Warehouse Warehouse `gorm:"foreignKey:WarehouseID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (WarehouseLocation) TableName() string {
	return "warehouse_locations"
}
