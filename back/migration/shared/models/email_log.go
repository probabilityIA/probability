package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

//
//	EMAIL LOGS - Logs de envío de emails de notificación
//

// EmailLog registra cada intento de envío de email de notificación.
// Útil para auditoría, debugging y reintentos futuros.
type EmailLog struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	CreatedAt time.Time `gorm:"not null;index"`

	// Contexto del negocio
	BusinessID    uint `gorm:"not null;index"`
	IntegrationID uint `gorm:"not null;index"`
	ConfigID      uint `gorm:"not null;index"`

	// Destinatario y contenido
	To        string `gorm:"size:255;not null;index"`
	Subject   string `gorm:"size:512;not null"`
	EventType string `gorm:"size:128;not null;index"`

	// Estado del envío: "sent" o "failed"
	Status       string  `gorm:"size:32;not null;index"`
	ErrorMessage *string `gorm:"type:text"`
}

// TableName especifica el nombre de la tabla para EmailLog
func (EmailLog) TableName() string {
	return "email_logs"
}

// BeforeCreate genera un UUID si no se ha asignado
func (e *EmailLog) BeforeCreate(tx *gorm.DB) error {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	return nil
}
