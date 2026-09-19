package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/pay/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/pay/internal/domain/entities"
	models "github.com/secamc93/probability/back/migration/shared/models"
)

func parseSpendDateRange(startDate, endDate string) (time.Time, time.Time, error) {
	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid start_date format: %w", err)
	}
	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid end_date format: %w", err)
	}
	end = end.Add(24 * time.Hour)
	return start, end, nil
}

// GetSpendSummary suma lo que un negocio ha gastado de su propia billetera,
// agrupado por concepto (GUIDE, SUBSCRIPTION, etc.). Solo movimientos type=USAGE
// (debitos reales); RECHARGE/REFUND/ADJUSTMENT son entradas de plata, no gasto.
func (r *Repository) GetSpendSummary(ctx context.Context, dto *dtos.SpendSummaryDTO) ([]dtos.ConceptTotal, error) {
	start, end, err := parseSpendDateRange(dto.StartDate, dto.EndDate)
	if err != nil {
		return nil, err
	}

	var rows []dtos.ConceptTotal
	err = r.db.Conn(ctx).
		Table("transaction").
		Select("concept, COUNT(*) AS count, COALESCE(SUM(amount),0) AS amount").
		Where("business_id = ? AND type = ? AND status = ? AND created_at BETWEEN ? AND ?",
			dto.BusinessID, entities.WalletTxTypeUsage, entities.WalletTxStatusCompleted, start, end).
		Group("concept").
		Order("amount DESC").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// GetGuideStatusSummary desglosa el gasto en guias (concept=GUIDE) por el
// estado real del envio, uniendo con shipments por shipment_id. Una guia sin
// shipment_id (dato viejo o borrado) cae en el estado "desconocido".
func (r *Repository) GetGuideStatusSummary(ctx context.Context, dto *dtos.SpendSummaryDTO) ([]dtos.GuideStatusTotal, error) {
	start, end, err := parseSpendDateRange(dto.StartDate, dto.EndDate)
	if err != nil {
		return nil, err
	}

	var rows []dtos.GuideStatusTotal
	err = r.db.Conn(ctx).
		Table("transaction AS t").
		Select("COALESCE(NULLIF(s.status, ''), 'desconocido') AS status, COUNT(*) AS count, COALESCE(SUM(t.amount),0) AS amount").
		Joins("LEFT JOIN shipments s ON s.id = t.shipment_id").
		Where("t.business_id = ? AND t.type = ? AND t.status = ? AND t.concept = ? AND t.created_at BETWEEN ? AND ?",
			dto.BusinessID, entities.WalletTxTypeUsage, entities.WalletTxStatusCompleted, entities.WalletTxConceptGuide, start, end).
		Group("COALESCE(NULLIF(s.status, ''), 'desconocido')").
		Order("count DESC").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// ListTransactionsFiltered lista paginada de los debitos (type=USAGE) de un
// negocio, opcionalmente filtrada por concepto.
func (r *Repository) ListTransactionsFiltered(ctx context.Context, dto *dtos.TransactionFilterDTO) ([]*entities.WalletTransaction, int64, error) {
	start, end, err := parseSpendDateRange(dto.StartDate, dto.EndDate)
	if err != nil {
		return nil, 0, err
	}

	whereClause := "business_id = ? AND type = ? AND status = ? AND created_at BETWEEN ? AND ?"
	args := []interface{}{dto.BusinessID, entities.WalletTxTypeUsage, entities.WalletTxStatusCompleted, start, end}
	if dto.Concept != "" {
		whereClause += " AND concept = ?"
		args = append(args, dto.Concept)
	}

	var total int64
	if err := r.db.Conn(ctx).Table("transaction").Where(whereClause, args...).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	page := dto.Page
	if page < 1 {
		page = 1
	}
	pageSize := dto.PageSize
	if pageSize < 1 {
		pageSize = 15
	}

	var rows []models.WalletTransaction
	if err := r.db.Conn(ctx).Table("transaction").Where(whereClause, args...).
		Order("created_at DESC").
		Limit(pageSize).
		Offset((page - 1) * pageSize).
		Find(&rows).Error; err != nil {
		return nil, 0, err
	}

	return walletTxListToDomain(rows), total, nil
}
