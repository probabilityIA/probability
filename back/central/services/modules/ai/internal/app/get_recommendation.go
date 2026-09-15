package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/ai/internal/domain/entities"
)

func (uc *UseCase) GetRecommendation(ctx context.Context, origin, destination string) (*entities.Recommendation, error) {
	return uc.recommendations.GetRecommendation(ctx, origin, destination)
}
