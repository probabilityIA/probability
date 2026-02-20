package domain

import (
	"time"

	"gorm.io/gorm"
)

// OriginAddress representa una dirección de origen física para envíos
type OriginAddress struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	BusinessID   uint   `gorm:"not null;index" json:"business_id"`
	Alias        string `gorm:"size:100;not null" json:"alias"`
	Company      string `gorm:"size:100" json:"company"`
	FirstName    string `gorm:"size:100" json:"first_name"`
	LastName     string `gorm:"size:100" json:"last_name"`
	Email        string `gorm:"size:100" json:"email"`
	Phone        string `gorm:"size:20" json:"phone"`
	Street       string `gorm:"size:255;not null" json:"street"`
	Suburb       string `gorm:"size:100" json:"suburb"`
	CityDaneCode string `gorm:"size:10;not null" json:"city_dane_code"`
	City         string `gorm:"size:100;not null" json:"city"`
	State        string `gorm:"size:100;not null" json:"state"`
	PostalCode   string `gorm:"size:20" json:"postal_code"`
	IsDefault    bool   `gorm:"default:false" json:"is_default"`
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
