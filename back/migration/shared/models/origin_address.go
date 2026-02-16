package models

import (
	"gorm.io/gorm"
)

// OriginAddress representa una dirección de origen configurada por el comercio
type OriginAddress struct {
	gorm.Model
	BusinessID uint   `gorm:"not null;index" json:"business_id"`
	Alias      string `gorm:"size:100;not null" json:"alias"` // Nombre descriptivo (e.g., "Bodega Principal")

	// Datos de contacto (requeridos por EnvioClick y otros transportadores)
	Company   string `gorm:"size:100" json:"company"`
	FirstName string `gorm:"size:100" json:"first_name"`
	LastName  string `gorm:"size:100" json:"last_name"`
	Email     string `gorm:"size:100" json:"email"`
	Phone     string `gorm:"size:20" json:"phone"`

	// Datos geográficos
	Street       string `gorm:"size:255;not null" json:"street"`
	Suburb       string `gorm:"size:100" json:"suburb"` // Colonia / Barrio
	CityDaneCode string `gorm:"size:10;not null" json:"city_dane_code"`
	City         string `gorm:"size:100;not null" json:"city"`
	State        string `gorm:"size:100;not null" json:"state"`
	PostalCode   string `gorm:"size:20" json:"postal_code"`

	IsDefault bool `gorm:"default:false" json:"is_default"`
}
