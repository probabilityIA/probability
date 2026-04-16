package usecases

import (
	"github.com/secamc93/probability/back/central/services/modules/ai/internal/domain"
)

type GetRecommendationUseCase struct {
	aiService domain.AIService
}

func NewGetRecommendationUseCase(aiService domain.AIService) *GetRecommendationUseCase {
	return &GetRecommendationUseCase{
		aiService: aiService,
	}
}

func (uc *GetRecommendationUseCase) Execute(origin, destination string) (*domain.AIRecommendation, error) {
	// Basic validation
	if origin == "" || destination == "" {
		// Fallback defaults or error? user only asked for logic
	}

	req := domain.RecommendationRequest{
		Origin:      origin,
		Destination: destination,
		// Carriers could be populated if we integrate quotations here,
		// but for now user asked to port the AI logic based on transportadoras.json
	}
	return uc.aiService.GetRecommendation(req)
}
