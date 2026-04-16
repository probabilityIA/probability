package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// BulkInvoiceJob representa un job de facturación masiva
type BulkInvoiceJob struct {
	ID           uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	BusinessID   uint       `gorm:"not null;index:idx_bulk_jobs_business"`
	CreatedByUserID *uint   `gorm:"column:created_by_user_id"`
	TotalOrders  int        `gorm:"not null;check:total_orders > 0"`
	Processed    int        `gorm:"not null;default:0;check:processed >= 0"`
	Successful   int        `gorm:"not null;default:0;check:successful >= 0"`
	Failed       int        `gorm:"not null;default:0;check:failed >= 0"`
	Status       string     `gorm:"type:varchar(20);not null;default:'pending';index:idx_bulk_jobs_status;check:status IN ('pending','processing','completed','failed')"`
	StartedAt    *time.Time
	CompletedAt  *time.Time
	ErrorMessage *string    `gorm:"type:text"`
	CreatedAt    time.Time  `gorm:"not null;default:now();index:idx_bulk_jobs_created_at,sort:desc"`
	UpdatedAt    time.Time  `gorm:"not null;default:now()"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`

	// Relaciones
	Business *Business `gorm:"foreignKey:BusinessID;constraint:OnDelete:CASCADE"`
	Items    []BulkInvoiceJobItem `gorm:"foreignKey:JobID;constraint:OnDelete:CASCADE"`
}

// TableName especifica el nombre de la tabla
func (BulkInvoiceJob) TableName() string {
	return "bulk_invoice_jobs"
}

// BeforeCreate hook para generar UUID si no existe
func (j *BulkInvoiceJob) BeforeCreate(tx *gorm.DB) error {
	if j.ID == uuid.Nil {
		j.ID = uuid.New()
	}
	return nil
}
