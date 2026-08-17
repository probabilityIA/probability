package usecasequotes

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/domain"
)

func (uc *UseCaseQuotes) SetSelection(ctx context.Context, in domain.QuoteSelectionInput) (*domain.SavedQuoteResponse, error) {
	quote, err := uc.repo.GetSavedQuoteByID(ctx, in.QuoteID)
	if err != nil {
		return nil, err
	}
	if quote == nil {
		return nil, ErrQuoteNotFound
	}
	if in.BusinessID > 0 && quote.BusinessID != in.BusinessID {
		return nil, ErrQuoteBusinessMatch
	}

	quote.SelectedCarrier = in.SelectedCarrier
	quote.SelectedIDRate = in.SelectedIDRate

	if err := uc.repo.UpdateSavedQuote(ctx, quote); err != nil {
		return nil, err
	}

	resp := toResponse(quote)
	return &resp, nil
}
