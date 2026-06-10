package usecasequotes

import (
	"context"
	"errors"

	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/domain"
)

var (
	ErrQuoteNotFound      = errors.New("cotizacion no encontrada")
	ErrQuoteBusinessMatch = errors.New("la cotizacion no pertenece a este negocio")
	ErrQuoteAlreadyLinked = errors.New("la cotizacion ya esta asociada a una orden")
)

func (uc *UseCaseQuotes) Associate(ctx context.Context, in domain.AssociateQuoteInput) (*domain.SavedQuoteResponse, error) {
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
	if quote.OrderUUID != nil && *quote.OrderUUID != "" && *quote.OrderUUID != in.OrderUUID {
		return nil, ErrQuoteAlreadyLinked
	}

	quote.OrderUUID = &in.OrderUUID
	if in.SelectedCarrier != "" {
		quote.SelectedCarrier = in.SelectedCarrier
	}
	if in.SelectedIDRate != nil {
		quote.SelectedIDRate = in.SelectedIDRate
	}
	if in.GuideRequested {
		quote.Status = domain.QuoteStatusGuideGenerated
	} else {
		quote.Status = domain.QuoteStatusAssociated
	}

	if err := uc.repo.UpdateSavedQuote(ctx, quote); err != nil {
		return nil, err
	}

	resp := toResponse(quote)
	return &resp, nil
}
