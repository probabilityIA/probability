package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

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

// WarehouseLocation representa una ubicación dentro de una bodega
type WarehouseLocation struct {
	gorm.Model
	WarehouseID   uint   `gorm:"not null;index;uniqueIndex:idx_location_warehouse_code,priority:1"` // ID de la bodega a la que pertenece esta ubicación
	Name          string `gorm:"size:255;not null"`                                                  // Nombre de la ubicación (ej: "Estante A-01")
	Code          string `gorm:"size:50;not null;uniqueIndex:idx_location_warehouse_code,priority:2"` // Código único dentro de la bodega (ej: "A-01")
	Type          string `gorm:"size:50;default:'storage'"`                                          // Tipo de ubicación: storage, picking, packing, receiving, shipping
	IsActive      bool   `gorm:"default:true;index"`                                                 // Indica si la ubicación está activa
	IsFulfillment bool   `gorm:"default:false;index"`                                                // Indica si la ubicación maneja fulfillment
	Capacity      *int   //                                                                            Capacidad máxima de la ubicación (nil = sin límite)

	// Relación
	Warehouse Warehouse `gorm:"foreignKey:WarehouseID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"` // Bodega a la que pertenece
}

// TableName especifica el nombre de la tabla
func (WarehouseLocation) TableName() string {
	return "warehouse_locations"
}
