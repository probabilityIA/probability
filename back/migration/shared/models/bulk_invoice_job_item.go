package models

import (
	"time"

	"github.com/google/uuid"
)

// BulkInvoiceJobItem representa un item individual de un job de facturación masiva
type BulkInvoiceJobItem struct {
	ID           uint       `gorm:"primary_key;autoIncrement"`
	JobID        uuid.UUID  `gorm:"type:uuid;not null;index:idx_bulk_job_items_job;uniqueIndex:uq_job_order"`
	OrderID      string     `gorm:"type:varchar(255);not null;index:idx_bulk_job_items_order;uniqueIndex:uq_job_order"`
	InvoiceID    *uint      `gorm:"index"`
	Status       string     `gorm:"type:varchar(20);not null;default:'pending';index:idx_bulk_job_items_status;check:status IN ('pending','processing','success','failed')"`
	ErrorMessage *string    `gorm:"type:text"`
	ProcessedAt  *time.Time
	CreatedAt    time.Time  `gorm:"not null;default:now()"`

	// Relaciones
	Job     *BulkInvoiceJob `gorm:"foreignKey:JobID;constraint:OnDelete:CASCADE"`
	Invoice *Invoice        `gorm:"foreignKey:InvoiceID;constraint:OnDelete:SET NULL"`
}

// TableName especifica el nombre de la tabla
func (BulkInvoiceJobItem) TableName() string {
	return "bulk_invoice_job_items"
}
