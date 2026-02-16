package domain

import (
	"gorm.io/gorm"
)

// OriginAddress representa una dirección de origen física para envíos
type OriginAddress struct {
	gorm.Model
	BusinessID   uint   `gorm:"not null;index"`
	Alias        string `gorm:"size:100;not null"`
	Company      string `gorm:"size:100"`
	FirstName    string `gorm:"size:100"`
	LastName     string `gorm:"size:100"`
	Email        string `gorm:"size:100"`
	Phone        string `gorm:"size:20"`
	Street       string `gorm:"size:255;not null"`
	Suburb       string `gorm:"size:100"`
	CityDaneCode string `gorm:"size:10;not null"`
	City         string `gorm:"size:100;not null"`
	State        string `gorm:"size:100;not null"`
	PostalCode   string `gorm:"size:20"`
	IsDefault    bool   `gorm:"default:false"`
}

// CreateOriginAddressRequest DTO para crear dirección
type CreateOriginAddressRequest struct {
	Alias        string `json:"alias" binding:"required"`
	Company      string `json:"company" binding:"required"`
	FirstName    string `json:"first_name" binding:"required"`
	LastName     string `json:"last_name" binding:"required"`
	Email        string `json:"email" binding:"required,email"`
	Phone        string `json:"phone" binding:"required"`
	Street       string `json:"street" binding:"required"`
	Suburb       string `json:"suburb"`
	CityDaneCode string `json:"city_dane_code" binding:"required"`
	City         string `json:"city" binding:"required"`
	State        string `json:"state" binding:"required"`
	PostalCode   string `json:"postal_code"`
	IsDefault    bool   `json:"is_default"`
}

// UpdateOriginAddressRequest DTO para actualizar dirección
type UpdateOriginAddressRequest struct {
	Alias        *string `json:"alias"`
	Company      *string `json:"company"`
	FirstName    *string `json:"first_name"`
	LastName     *string `json:"last_name"`
	Email        *string `json:"email" binding:"omitempty,email"`
	Phone        *string `json:"phone"`
	Street       *string `json:"street"`
	Suburb       *string `json:"suburb"`
	CityDaneCode *string `json:"city_dane_code"`
	City         *string `json:"city"`
	State        *string `json:"state"`
	PostalCode   *string `json:"postal_code"`
	IsDefault    *bool   `json:"is_default"`
}
