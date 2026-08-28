package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/log"
)

type InvoiceReadRepository struct {
	db  db.IDatabase
	log log.ILogger
}

func NewInvoiceReadRepository(database db.IDatabase, logger log.ILogger) *InvoiceReadRepository {
	return &InvoiceReadRepository{db: database, log: logger}
}

func (r *InvoiceReadRepository) GetInvoiceRef(ctx context.Context, invoiceID uint) (*dtos.InvoiceRef, error) {
	var fila struct {
		ID                     uint
		BusinessID             uint
		ExternalID             *string
		InvoicingIntegrationID *uint
	}

	err := r.db.Conn(ctx).
		Table("invoices").
		Select("id", "business_id", "external_id", "invoicing_integration_id").
		Where("id = ? AND deleted_at IS NULL", invoiceID).
		Limit(1).
		Scan(&fila).Error
	if err != nil {
		return nil, err
	}
	if fila.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	ref := &dtos.InvoiceRef{
		ID:         fila.ID,
		BusinessID: fila.BusinessID,
	}
	if fila.ExternalID != nil {
		ref.ExternalID = *fila.ExternalID
	}
	if fila.InvoicingIntegrationID != nil {
		ref.IntegrationID = *fila.InvoicingIntegrationID
	}
	return ref, nil
}
