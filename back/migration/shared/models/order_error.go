package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// OrderError registra los errores ocurridos durante el procesamiento de órdenes
type OrderError struct {
	gorm.Model

	// Contexto del error
	ExternalID      string `gorm:"size:255;index"` // ID en plataforma externa (si se pudo extraer)
	IntegrationID   uint   `gorm:"index"`          // ID de la integración (si se conoce)
	BusinessID      *uint  `gorm:"index"`          // ID del negocio (si se conoce)
	IntegrationType string `gorm:"size:50;index"`  // "shopify", "whatsapp", etc.
	Platform        string `gorm:"size:50;index"`  // Plataforma origen

	// Detalles del error
	ErrorType    string         `gorm:"size:100;index"`     // Tipo de error (ej: "validation_error", "database_error")
	ErrorMessage string         `gorm:"type:text;not null"` // Mensaje de error detallado
	ErrorStack   *string        `gorm:"type:text"`          // Stack trace (opcional)
	RawData      datatypes.JSON `gorm:"type:jsonb"`         // Payload original que causó el error

	// Estado
	Status     string     `gorm:"size:50;default:'new';index"` // "new", "resolved", "ignored"
	ResolvedAt *time.Time // Cuándo se resolvió
	ResolvedBy *uint      `gorm:"index"`     // Usuario que resolvió
	Resolution *string    `gorm:"type:text"` // Notas de resolución

	// Relaciones (opcionales, solo si existen)
	Integration *Integration `gorm:"foreignKey:IntegrationID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Business    *Business    `gorm:"foreignKey:BusinessID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

// TableName especifica el nombre de la tabla para OrderError
func (OrderError) TableName() string {
	return "order_errors"
}
