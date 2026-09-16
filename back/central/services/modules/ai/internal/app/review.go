package app

import (
	"context"
	"strings"

	"github.com/secamc93/probability/back/central/services/modules/ai/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/ai/internal/domain/entities"
)

const maxReviewSearchRunes = 200

func (uc *UseCase) ListReviewMessages(ctx context.Context, filter dtos.ReviewFilter) (*dtos.PaginatedResponse[entities.ReviewMessage], error) {
	filter = normalizeReviewFilter(filter)
	result := &dtos.PaginatedResponse[entities.ReviewMessage]{
		Data:       []entities.ReviewMessage{},
		Page:       filter.Page,
		PageSize:   filter.PageSize,
		TotalPages: 1,
	}
	if uc.conversations == nil {
		return result, nil
	}

	items, total, err := uc.conversations.ListMessages(ctx, filter)
	if err != nil {
		return nil, err
	}
	result.Data = items
	result.Total = total
	if pages := int((total + int64(filter.PageSize) - 1) / int64(filter.PageSize)); pages > 1 {
		result.TotalPages = pages
	}
	return result, nil
}

func (uc *UseCase) GetReviewSummary(ctx context.Context, filter dtos.ReviewFilter) (*entities.ReviewSummary, error) {
	filter = normalizeReviewFilter(filter)
	if uc.conversations == nil {
		return &entities.ReviewSummary{TopDestinations: []entities.DestinationCount{}}, nil
	}

	summary, err := uc.conversations.Summary(ctx, filter)
	if err != nil {
		return nil, err
	}
	summary.EstimatedCostUSD = EstimateCostUSD(summary.InputTokens, summary.OutputTokens)
	return summary, nil
}

func EstimateCostUSD(inputTokens, outputTokens int64) float64 {
	return float64(inputTokens)*InputCostPerMillionUSD/1_000_000 + float64(outputTokens)*OutputCostPerMillionUSD/1_000_000
}

func normalizeReviewFilter(filter dtos.ReviewFilter) dtos.ReviewFilter {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 {
		filter.PageSize = DefaultReviewPageSize
	}
	if filter.PageSize > MaxReviewPageSize {
		filter.PageSize = MaxReviewPageSize
	}
	switch filter.Kind {
	case dtos.ReviewNoDestination, dtos.ReviewNegative, dtos.ReviewPositive, dtos.ReviewErrors, dtos.ReviewNotClicked:
	default:
		filter.Kind = dtos.ReviewAll
	}
	filter.Search = truncateRunes(strings.TrimSpace(filter.Search), maxReviewSearchRunes)
	return filter
}
