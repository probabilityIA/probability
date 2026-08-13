package app

import (
	"context"
	"errors"

	"github.com/secamc93/probability/back/central/services/integrations/syncruns/internal/domain"
)

func (uc *useCase) MatchMatrix(ctx context.Context, query domain.MatrixQuery) (*domain.MatrixPage, error) {
	if query.BusinessID == 0 {
		return nil, errors.New("business_id es requerido")
	}
	query.Normalize()

	return uc.repo.MatchMatrix(ctx, query)
}
