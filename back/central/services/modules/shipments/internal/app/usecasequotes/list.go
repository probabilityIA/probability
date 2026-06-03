package usecasequotes

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/domain"
)

func (uc *UseCaseQuotes) List(ctx context.Context, filter domain.SavedQuoteFilter) (*domain.SavedQuotesListResponse, error) {
	rows, total, err := uc.repo.ListSavedQuotes(ctx, filter)
	if err != nil {
		return nil, err
	}

	page := filter.Page
	if page < 1 {
		page = 1
	}
	pageSize := filter.PageSize
	if pageSize < 1 {
		pageSize = 10
	}

	data := make([]domain.SavedQuoteResponse, len(rows))
	for i := range rows {
		data[i] = toResponse(&rows[i])
	}

	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))

	return &domain.SavedQuotesListResponse{
		Data:       data,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

func (uc *UseCaseQuotes) GetByID(ctx context.Context, id uint) (*domain.SavedQuoteResponse, error) {
	q, err := uc.repo.GetSavedQuoteByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if q == nil {
		return nil, nil
	}
	resp := toResponse(q)
	return &resp, nil
}

func toResponse(q *domain.SavedQuote) domain.SavedQuoteResponse {
	return domain.SavedQuoteResponse{
		ID:                  q.ID,
		BusinessID:          q.BusinessID,
		IntegrationID:       q.IntegrationID,
		Source:              q.Source,
		CorrelationID:       q.CorrelationID,
		OrderUUID:           q.OrderUUID,
		ExternalOrderRef:    q.ExternalOrderRef,
		Rates:               q.Rates,
		SelectedCarrier:     q.SelectedCarrier,
		SelectedServiceCode: q.SelectedServiceCode,
		Status:              q.Status,
		ExpiresAt:           q.ExpiresAt,
		CreatedAt:           q.CreatedAt,
	}
}
