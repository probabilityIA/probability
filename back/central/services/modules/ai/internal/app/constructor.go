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
	MaxPathnameRunes      = 255

	ConversationRetention = 365 * 24 * time.Hour

	InputCostPerMillionUSD  = 0.15
	OutputCostPerMillionUSD = 1.20

	DefaultReviewPageSize = 20
	MaxReviewPageSize     = 100
)

type IUseCase interface {
	GetRecommendation(ctx context.Context, origin, destination string) (*entities.Recommendation, error)

	Chat(ctx context.Context, input dtos.ChatInput) (*entities.AssistantReply, error)
	GetAssistantState(ctx context.Context, userID uint) (*entities.AssistantState, error)
	MarkIntroSeen(ctx context.Context, userID uint) error

	SubmitFeedback(ctx context.Context, userID uint, messageID string, value int) error
	MarkDestinationClicked(ctx context.Context, userID uint, messageID string) error

	PersistMessage(ctx context.Context, record entities.MessageRecord) error
	PersistFeedback(ctx context.Context, userID uint, messageID string, value int, at time.Time) error
	PersistClick(ctx context.Context, userID uint, messageID string, at time.Time) error

	ListReviewMessages(ctx context.Context, filter dtos.ReviewFilter) (*dtos.PaginatedResponse[entities.ReviewMessage], error)
	GetReviewSummary(ctx context.Context, filter dtos.ReviewFilter) (*entities.ReviewSummary, error)
	PurgeExpiredConversations(ctx context.Context) (int64, error)
}

type UseCase struct {
	recommendations ports.IRecommendationProvider
	model           ports.IAssistantModel
	navigation      ports.INavigationCatalog
	store           ports.IAssistantStore
	recorder        ports.IConversationRecorder
	conversations   ports.IConversationRepository
	log             log.ILogger
}

func New(
	recommendations ports.IRecommendationProvider,
	model ports.IAssistantModel,
	navigation ports.INavigationCatalog,
	store ports.IAssistantStore,
	recorder ports.IConversationRecorder,
	conversations ports.IConversationRepository,
	logger log.ILogger,
) IUseCase {
	return &UseCase{
		recommendations: recommendations,
		model:           model,
		navigation:      navigation,
		store:           store,
		recorder:        recorder,
		conversations:   conversations,
		log:             logger,
	}
}
