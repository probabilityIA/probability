package repository

import (
	"context"
	stderrors "errors"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/accounting/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/accounting/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/accounting/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/accounting/internal/infra/secondary/repository/mappers"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

func (r *Repository) CreateInvoice(ctx context.Context, inv *entities.Invoice) (*entities.Invoice, error) {
	model := &models.AccountingInvoice{
		BusinessID:       inv.BusinessID,
		ConceptID:        inv.ConceptID,
		IssueDate:        inv.IssueDate,
		DueDate:          inv.DueDate,
		Status:           inv.Status,
		Notes:            inv.Notes,
		Subtotal:         inv.Subtotal,
		TaxTotal:         inv.TaxTotal,
		Total:            inv.Total,
		TaxDetail:        datatypes.JSON(mappers.TaxLinesToJSON(inv.TaxDetail)),
		EmailTo:          inv.EmailTo,
		CustomerDocument: inv.CustomerDocument,
		CustomerPhone:    inv.CustomerPhone,
		CustomerAddress:  inv.CustomerAddress,
		DianStatus:       entities.DianStatusNone,
		IsTest:           inv.IsTest,
	}
	err := r.db.Conn(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(model).Error; err != nil {
			return err
		}
		for _, item := range inv.Items {
			if err := tx.Create(&models.AccountingInvoiceItem{
				InvoiceID:   model.ID,
				Description: item.Description,
				Quantity:    item.Quantity,
				UnitPrice:   item.UnitPrice,
				Amount:      item.Amount,
			}).Error; err != nil {
				return err
			}
		}
		number := fmt.Sprintf("FP-%d-%05d", model.CreatedAt.Year(), model.ID)
		return tx.Model(&models.AccountingInvoice{}).Where("id = ?", model.ID).Update("number", number).Error
	})
	if err != nil {
		return nil, err
	}
	inv.ID = model.ID
	return inv, nil
}

func (r *Repository) UpdateInvoice(ctx context.Context, inv *entities.Invoice) (*entities.Invoice, error) {
	err := r.db.Conn(ctx).Transaction(func(tx *gorm.DB) error {
		updates := map[string]interface{}{
			"concept_id":        inv.ConceptID,
			"issue_date":        inv.IssueDate,
			"due_date":          inv.DueDate,
			"notes":             inv.Notes,
			"email_to":          inv.EmailTo,
			"subtotal":          inv.Subtotal,
			"tax_total":         inv.TaxTotal,
			"total":             inv.Total,
			"tax_detail":        datatypes.JSON(mappers.TaxLinesToJSON(inv.TaxDetail)),
			"customer_document": inv.CustomerDocument,
			"customer_phone":    inv.CustomerPhone,
			"customer_address":  inv.CustomerAddress,
			"is_test":           inv.IsTest,
		}
		res := tx.Model(&models.AccountingInvoice{}).Where("id = ?", inv.ID).Updates(updates)
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return domainerrors.ErrInvoiceNotFound
		}
		if err := tx.Where("invoice_id = ?", inv.ID).Delete(&models.AccountingInvoiceItem{}).Error; err != nil {
			return err
		}
		for _, item := range inv.Items {
			if err := tx.Create(&models.AccountingInvoiceItem{
				InvoiceID:   inv.ID,
				Description: item.Description,
				Quantity:    item.Quantity,
				UnitPrice:   item.UnitPrice,
				Amount:      item.Amount,
			}).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return inv, nil
}

type invoiceRow struct {
	models.AccountingInvoice
	BusinessName string
	ConceptName  string
}

func (r *Repository) invoiceToEntity(row *invoiceRow, items []models.AccountingInvoiceItem) *entities.Invoice {
	inv := &entities.Invoice{
		ID:               row.ID,
		Number:           row.Number,
		BusinessID:       row.BusinessID,
		BusinessName:     row.BusinessName,
		ConceptID:        row.ConceptID,
		ConceptName:      row.ConceptName,
		IssueDate:        row.IssueDate,
		DueDate:          row.DueDate,
		Status:           row.Status,
		Notes:            row.Notes,
		Subtotal:         row.Subtotal,
		TaxTotal:         row.TaxTotal,
		Total:            row.Total,
		EmailTo:          row.EmailTo,
		SentAt:           row.SentAt,
		PaidAt:           row.PaidAt,
		CustomerDocument: row.CustomerDocument,
		CustomerPhone:    row.CustomerPhone,
		CustomerAddress:  row.CustomerAddress,
		DianStatus:       row.DianStatus,
		Cufe:             row.Cufe,
		DianNumber:       row.DianNumber,
		DianExternalID:   row.DianExternalID,
		DianQR:           row.DianQR,
		DianEmittedAt:    row.DianEmittedAt,
		IsTest:           row.IsTest,
		CreatedAt:        row.CreatedAt,
		UpdatedAt:        row.UpdatedAt,
	}
	entry := models.AccountingEntry{TaxDetail: row.TaxDetail}
	inv.TaxDetail = mappers.EntryToEntity(&entry, "", "", "").TaxDetail
	for _, item := range items {
		inv.Items = append(inv.Items, entities.InvoiceItem{
			ID:          item.ID,
			Description: item.Description,
			Quantity:    item.Quantity,
			UnitPrice:   item.UnitPrice,
			Amount:      item.Amount,
		})
	}
	return inv
}

func (r *Repository) GetInvoiceByID(ctx context.Context, id uint) (*entities.Invoice, error) {
	var row invoiceRow
	err := r.db.Conn(ctx).Model(&models.AccountingInvoice{}).
		Select("accounting_invoices.*, COALESCE(b.name, '') AS business_name, COALESCE(c.name, '') AS concept_name").
		Joins("LEFT JOIN business b ON b.id = accounting_invoices.business_id").
		Joins("LEFT JOIN accounting_concepts c ON c.id = accounting_invoices.concept_id").
		Where("accounting_invoices.id = ? AND accounting_invoices.deleted_at IS NULL", id).
		Scan(&row).Error
	if err != nil {
		return nil, err
	}
	if row.ID == 0 {
		return nil, domainerrors.ErrInvoiceNotFound
	}
	var items []models.AccountingInvoiceItem
	if err := r.db.Conn(ctx).Where("invoice_id = ?", id).Order("id ASC").Find(&items).Error; err != nil {
		return nil, err
	}
	return r.invoiceToEntity(&row, items), nil
}

func (r *Repository) ListInvoices(ctx context.Context, params dtos.ListInvoicesParams) ([]entities.Invoice, int64, error) {
	base := func() *gorm.DB {
		q := r.db.Conn(ctx).Model(&models.AccountingInvoice{}).
			Where("accounting_invoices.deleted_at IS NULL")
		if params.BusinessID != nil {
			q = q.Where("accounting_invoices.business_id = ?", *params.BusinessID)
		}
		if params.Status != "" {
			q = q.Where("accounting_invoices.status = ?", params.Status)
		}
		if params.IsTest != nil {
			q = q.Where("accounting_invoices.is_test = ?", *params.IsTest)
		}
		return q
	}
	var total int64
	if err := base().Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []invoiceRow
	err := base().
		Select("accounting_invoices.*, COALESCE(b.name, '') AS business_name, COALESCE(c.name, '') AS concept_name").
		Joins("LEFT JOIN business b ON b.id = accounting_invoices.business_id").
		Joins("LEFT JOIN accounting_concepts c ON c.id = accounting_invoices.concept_id").
		Order("accounting_invoices.id DESC").
		Offset(params.Offset()).Limit(params.PageSize).
		Scan(&rows).Error
	if err != nil {
		return nil, 0, err
	}
	out := make([]entities.Invoice, len(rows))
	for i := range rows {
		out[i] = *r.invoiceToEntity(&rows[i], nil)
	}
	return out, total, nil
}

func (r *Repository) DeleteInvoice(ctx context.Context, id uint) error {
	return r.db.Conn(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("invoice_id = ?", id).Delete(&models.AccountingInvoiceItem{}).Error; err != nil {
			return err
		}
		res := tx.Delete(&models.AccountingInvoice{}, id)
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return domainerrors.ErrInvoiceNotFound
		}
		return nil
	})
}

func (r *Repository) SetInvoiceStatus(ctx context.Context, id uint, status string, sentAt, paidAt *time.Time, emailTo string) error {
	updates := map[string]interface{}{
		"status":   status,
		"sent_at":  sentAt,
		"paid_at":  paidAt,
		"email_to": emailTo,
	}
	res := r.db.Conn(ctx).Model(&models.AccountingInvoice{}).Where("id = ?", id).Updates(updates)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return domainerrors.ErrInvoiceNotFound
	}
	return nil
}

func (r *Repository) SetInvoiceCustomer(ctx context.Context, id uint, document, phone, address string) error {
	res := r.db.Conn(ctx).Model(&models.AccountingInvoice{}).Where("id = ?", id).Updates(map[string]interface{}{
		"customer_document": document,
		"customer_phone":    phone,
		"customer_address":  address,
	})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return domainerrors.ErrInvoiceNotFound
	}
	return nil
}

func (r *Repository) SetInvoiceDian(ctx context.Context, id uint, result *entities.DianEmitResult) error {
	now := time.Now()
	res := r.db.Conn(ctx).Model(&models.AccountingInvoice{}).Where("id = ?", id).Updates(map[string]interface{}{
		"dian_status":      entities.DianStatusValidated,
		"cufe":             result.CUFE,
		"dian_number":      result.Number,
		"dian_external_id": result.ExternalID,
		"dian_qr":          result.QRCode,
		"dian_emitted_at":  &now,
	})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return domainerrors.ErrInvoiceNotFound
	}
	return nil
}

func (r *Repository) DeleteEntryBySource(ctx context.Context, sourceType, sourceID string) error {
	return r.db.Conn(ctx).
		Where("source_type = ? AND source_id = ?", sourceType, sourceID).
		Delete(&models.AccountingEntry{}).Error
}

func (r *Repository) GetTaxesByIDs(ctx context.Context, ids []uint) ([]entities.Tax, error) {
	var taxModels []models.AccountingTax
	if err := r.db.Conn(ctx).Where("id IN ?", ids).Find(&taxModels).Error; err != nil {
		return nil, err
	}
	out := make([]entities.Tax, len(taxModels))
	for i := range taxModels {
		out[i] = *mappers.TaxToEntity(&taxModels[i])
	}
	return out, nil
}

func (r *Repository) GetBusinessName(ctx context.Context, id uint) (string, error) {
	var result struct{ Name string }
	err := r.db.Conn(ctx).Table("business").Select("name").
		Where("id = ? AND deleted_at IS NULL", id).Limit(1).First(&result).Error
	if err != nil {
		if stderrors.Is(err, gorm.ErrRecordNotFound) {
			return "", domainerrors.ErrBusinessNotFound
		}
		return "", err
	}
	return result.Name, nil
}

func (r *Repository) GetBusinessDefaultEmail(ctx context.Context, id uint) (string, error) {
	var result struct{ Email string }
	err := r.db.Conn(ctx).Table("user_businesses ub").
		Select("u.email").
		Joins(`JOIN "user" u ON u.id = ub.user_id AND u.deleted_at IS NULL AND u.is_active`).
		Where("ub.business_id = ?", id).
		Order("u.id ASC").Limit(1).First(&result).Error
	if err != nil {
		if stderrors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil
		}
		return "", err
	}
	return result.Email, nil
}
