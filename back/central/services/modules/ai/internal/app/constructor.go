package app

import (
	"context"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/ai/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/ai/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/ai/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

const (
	AssistantMessageLimit = 30
	AssistantWindow       = time.Hour
	MaxHistoryMessages    = 12
	MaxUserMessageRunes   = 1000
	MaxAssistantRunes     = 2000
)

type IUseCase interface {
	GetRecommendation(ctx context.Context, origin, destination string) (*entities.Recommendation, error)
	Chat(ctx context.Context, input dtos.ChatInput) (*entities.AssistantReply, error)
	GetAssistantState(ctx context.Context, userID uint) (*entities.AssistantState, error)
	MarkIntroSeen(ctx context.Context, userID uint) error
}

type UseCase struct {
	recommendations ports.IRecommendationProvider
	model           ports.IAssistantModel
	navigation      ports.INavigationCatalog
	store           ports.IAssistantStore
	log             log.ILogger
}

func New(
	recommendations ports.IRecommendationProvider,
	model ports.IAssistantModel,
	navigation ports.INavigationCatalog,
	store ports.IAssistantStore,
	logger log.ILogger,
) IUseCase {
	return &UseCase{
		recommendations: recommendations,
		model:           model,
		navigation:      navigation,
		store:           store,
		log:             logger,
	}
}
