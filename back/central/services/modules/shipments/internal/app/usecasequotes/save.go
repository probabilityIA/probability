package usecasequotes

import (
	"context"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/domain"
)

const defaultQuoteTTL = 24 * time.Hour

func (uc *UseCaseQuotes) SaveQuote(ctx context.Context, in domain.SaveQuoteInput) (*domain.SavedQuote, error) {
	ttl := in.TTL
	if ttl <= 0 {
		ttl = defaultQuoteTTL
	}
	expiresAt := time.Now().Add(ttl)

	quote := &domain.SavedQuote{
		BusinessID:       in.BusinessID,
		IntegrationID:    in.IntegrationID,
		Source:           in.Source,
		CorrelationID:    in.CorrelationID,
		OrderUUID:        in.OrderUUID,
		ExternalOrderRef: in.ExternalOrderRef,
		RequestPayload:   in.RequestPayload,
		Rates:            in.Rates,
		Status:           domain.QuoteStatusCreated,
		ExpiresAt:        &expiresAt,
	}

	if err := uc.repo.CreateSavedQuote(ctx, quote); err != nil {
		return nil, err
	}
	return quote, nil
}
