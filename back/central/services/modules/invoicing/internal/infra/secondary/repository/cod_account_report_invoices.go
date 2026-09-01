package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/infra/secondary/repository/mappers"
	"github.com/secamc93/probability/back/migration/shared/models"
)

func (r *Repository) GetCODAccountReportInvoices(ctx context.Context, businessID uint, startDate, endDate, accountNumber string, isCOD bool, page, pageSize int) ([]*entities.Invoice, int64, error) {
	whereClause := "invoices.business_id = ? AND invoices.deleted_at IS NULL AND invoices.status != 'cancelled' AND orders.is_cod = ?" +
		" AND invoices.provider_response->'cash_receipt'->'request_body'->'payment'->0->>'accountNumber' = ?"
	args := []interface{}{businessID, isCOD, accountNumber}

	if startDate != "" {
		whereClause += " AND COALESCE(invoices.issued_at, invoices.created_at) >= ?"
		args = append(args, startDate)
	}
	if endDate != "" {
		whereClause += " AND COALESCE(invoices.issued_at, invoices.created_at) < ?::date + interval '1 day'"
		args = append(args, endDate)
	}

	query := r.db.Conn(ctx).
		Model(&models.Invoice{}).
		Joins("JOIN orders ON orders.id = invoices.order_id AND orders.deleted_at IS NULL").
		Where(whereClause, args...)

	var total int64
	if err := query.Count(&total).Error; err != nil {
		r.log.Error(ctx).Err(err).Msg("Failed to count COD account report invoices")
		return nil, 0, fmt.Errorf("failed to count invoices: %w", err)
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	var modelsList []*models.Invoice
	if err := query.
		Select("invoices.*").
		Order("COALESCE(invoices.issued_at, invoices.created_at) DESC").
		Offset(offset).Limit(pageSize).
		Find(&modelsList).Error; err != nil {
		r.log.Error(ctx).Err(err).Msg("Failed to list COD account report invoices")
		return nil, 0, fmt.Errorf("failed to list invoices: %w", err)
	}

	domainList := mappers.InvoiceListToDomain(modelsList)

	if len(modelsList) > 0 {
		orderIDs := make([]string, 0, len(modelsList))
		for _, m := range modelsList {
			orderIDs = append(orderIDs, m.OrderID)
		}

		var orderRows []struct {
			ID          string
			OrderNumber string
		}
		if err := r.db.Conn(ctx).
			Table("orders").
			Select("id, order_number").
			Where("id IN (?)", orderIDs).
			Where("deleted_at IS NULL").
			Scan(&orderRows).Error; err == nil {
			orderNumberMap := make(map[string]string, len(orderRows))
			for _, row := range orderRows {
				orderNumberMap[row.ID] = row.OrderNumber
			}
			for _, inv := range domainList {
				if num, ok := orderNumberMap[inv.OrderID]; ok {
					inv.OrderNumber = num
				}
			}
		}
	}

	return domainList, total, nil
}
