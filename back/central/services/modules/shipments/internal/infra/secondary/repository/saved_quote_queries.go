package repository

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/domain"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

func (r *Repository) CreateSavedQuote(ctx context.Context, quote *domain.SavedQuote) error {
	payloadJSON, _ := json.Marshal(quote.RequestPayload)
	ratesJSON, _ := json.Marshal(quote.Rates)

	m := &models.ShippingQuote{
		BusinessID:          quote.BusinessID,
		IntegrationID:       quote.IntegrationID,
		Source:              quote.Source,
		CorrelationID:       quote.CorrelationID,
		OrderUUID:           quote.OrderUUID,
		ExternalOrderRef:    quote.ExternalOrderRef,
		RequestPayload:      datatypes.JSON(payloadJSON),
		Rates:               datatypes.JSON(ratesJSON),
		SelectedCarrier:     quote.SelectedCarrier,
		SelectedServiceCode: quote.SelectedServiceCode,
		SelectedIDRate:      quote.SelectedIDRate,
		Status:              quote.Status,
		ExpiresAt:           quote.ExpiresAt,
	}
	if err := r.db.Conn(ctx).Create(m).Error; err != nil {
		return err
	}
	quote.ID = m.ID
	quote.CreatedAt = m.CreatedAt
	quote.UpdatedAt = m.UpdatedAt
	return nil
}

func (r *Repository) GetSavedQuoteByID(ctx context.Context, id uint) (*domain.SavedQuote, error) {
	var m models.ShippingQuote
	if err := r.db.Conn(ctx).First(&m, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	q := savedQuoteModelToDomain(&m)
	r.fillQuoteOrderNumbers(ctx, []*domain.SavedQuote{q})
	return q, nil
}

func (r *Repository) ListSavedQuotes(ctx context.Context, filter domain.SavedQuoteFilter) ([]domain.SavedQuote, int64, error) {
	q := r.db.Conn(ctx).Model(&models.ShippingQuote{}).Where("business_id = ?", filter.BusinessID)
	if filter.Source != "" {
		q = q.Where("source = ?", filter.Source)
	}
	if filter.Status != "" {
		q = q.Where("status = ?", filter.Status)
	}
	if filter.OrderUUID != "" {
		q = q.Where("order_uuid = ?", filter.OrderUUID)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	page := filter.Page
	if page < 1 {
		page = 1
	}
	pageSize := filter.PageSize
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	var rows []models.ShippingQuote
	if err := q.Order("created_at DESC").
		Limit(pageSize).
		Offset((page - 1) * pageSize).
		Find(&rows).Error; err != nil {
		return nil, 0, err
	}

	out := make([]domain.SavedQuote, len(rows))
	refs := make([]*domain.SavedQuote, len(rows))
	for i := range rows {
		out[i] = *savedQuoteModelToDomain(&rows[i])
		refs[i] = &out[i]
	}
	r.fillQuoteOrderNumbers(ctx, refs)
	return out, total, nil
}

func (r *Repository) fillQuoteOrderNumbers(ctx context.Context, quotes []*domain.SavedQuote) {
	ids := make([]string, 0, len(quotes))
	for _, q := range quotes {
		if q != nil && q.OrderUUID != nil && *q.OrderUUID != "" {
			ids = append(ids, *q.OrderUUID)
		}
	}
	if len(ids) == 0 {
		return
	}

	var rows []struct {
		ID          string
		OrderNumber string
	}
	if err := r.db.Conn(ctx).
		Table("orders").
		Select("id, order_number").
		Where("id IN ? AND deleted_at IS NULL", ids).
		Scan(&rows).Error; err != nil {
		return
	}

	numbers := make(map[string]string, len(rows))
	for _, row := range rows {
		numbers[row.ID] = row.OrderNumber
	}
	for _, q := range quotes {
		if q != nil && q.OrderUUID != nil {
			q.OrderNumber = numbers[*q.OrderUUID]
		}
	}
}

func (r *Repository) UpdateSavedQuote(ctx context.Context, quote *domain.SavedQuote) error {
	updates := map[string]interface{}{
		"order_uuid":            quote.OrderUUID,
		"selected_carrier":      quote.SelectedCarrier,
		"selected_service_code": quote.SelectedServiceCode,
		"selected_id_rate":      quote.SelectedIDRate,
		"status":                quote.Status,
	}
	return r.db.Conn(ctx).
		Model(&models.ShippingQuote{}).
		Where("id = ?", quote.ID).
		Updates(updates).Error
}

func (r *Repository) GetOrderSelectedShipping(ctx context.Context, orderUUID string) (*domain.OrderSelectedShipping, error) {
	query := `
		SELECT sl->>'code' AS code, sl->>'title' AS title, sl->>'source' AS source, sl->>'price' AS price
		FROM orders o
		CROSS JOIN LATERAL jsonb_array_elements(COALESCE(o.shipping_details->'shipping_lines', '[]'::jsonb)) sl
		WHERE o.id = @order
		ORDER BY (lower(COALESCE(sl->>'source', '')) LIKE 'probability%') DESC
		LIMIT 1`

	var result struct {
		Code   *string
		Title  *string
		Source *string
		Price  *string
	}
	err := r.db.Conn(ctx).
		Raw(query, map[string]interface{}{"order": orderUUID}).
		Scan(&result).Error
	if err != nil {
		return nil, err
	}
	if result.Code == nil && result.Title == nil {
		return nil, nil
	}
	return &domain.OrderSelectedShipping{
		Code:   strDeref(result.Code),
		Title:  strDeref(result.Title),
		Source: strDeref(result.Source),
		Price:  strDeref(result.Price),
	}, nil
}

func (r *Repository) GetIntegrationConfigFlag(ctx context.Context, integrationID uint, key string) (bool, error) {
	var val string
	err := r.db.Conn(ctx).
		Table("integrations").
		Select("COALESCE(config->>?, '')", key).
		Where("id = ? AND deleted_at IS NULL", integrationID).
		Limit(1).
		Scan(&val).Error
	if err != nil {
		return false, err
	}
	return strings.EqualFold(strings.TrimSpace(val), "true"), nil
}

func (r *Repository) GetIntegrationConfigValue(ctx context.Context, integrationID uint, key string) (string, error) {
	var val string
	err := r.db.Conn(ctx).
		Table("integrations").
		Select("COALESCE(config->>?, '')", key).
		Where("id = ? AND deleted_at IS NULL", integrationID).
		Limit(1).
		Scan(&val).Error
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(val), nil
}

func savedQuoteModelToDomain(m *models.ShippingQuote) *domain.SavedQuote {
	d := &domain.SavedQuote{
		ID:                  m.ID,
		BusinessID:          m.BusinessID,
		IntegrationID:       m.IntegrationID,
		Source:              m.Source,
		CorrelationID:       m.CorrelationID,
		OrderUUID:           m.OrderUUID,
		ExternalOrderRef:    m.ExternalOrderRef,
		SelectedCarrier:     m.SelectedCarrier,
		SelectedServiceCode: m.SelectedServiceCode,
		SelectedIDRate:      m.SelectedIDRate,
		Status:              m.Status,
		ExpiresAt:           m.ExpiresAt,
		CreatedAt:           m.CreatedAt,
		UpdatedAt:           m.UpdatedAt,
	}
	if len(m.RequestPayload) > 0 {
		_ = json.Unmarshal(m.RequestPayload, &d.RequestPayload)
	}
	if len(m.Rates) > 0 {
		_ = json.Unmarshal(m.Rates, &d.Rates)
	}
	return d
}

func strDeref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
