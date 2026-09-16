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

	MaxToolRounds       = 4
	MaxToolCallsPerTurn = 3
	MaxListedOrders     = 10
	DefaultSummaryDays  = 7

	ConversationRetention = 365 * 24 * time.Hour

	InputCostPerMillionUSD  = 0.15
	OutputCostPerMillionUSD = 1.20

	DefaultReviewPageSize = 20
	MaxReviewPageSize     = 100
)

var colombia = time.FixedZone("COT", -5*60*60)

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

	IngestAlertEvent(ctx context.Context, event dtos.AlertEvent) error
	ListAlerts(ctx context.Context, query dtos.AlertQuery) (*dtos.PaginatedResponse[entities.Alert], error)
	GetAlertsUnread(ctx context.Context, businessID, userID uint) (*entities.AlertsUnread, error)
	MarkAlertsSeen(ctx context.Context, businessID, userID uint) error
	PurgeExpiredAlerts(ctx context.Context) (int64, error)
}

type UseCase struct {
	recommendations ports.IRecommendationProvider
	model           ports.IAssistantModel
	navigation      ports.INavigationCatalog
	store           ports.IAssistantStore
	recorder        ports.IConversationRecorder
	conversations   ports.IConversationRepository
	businessData    ports.IBusinessDataReader
	alerts          ports.IAlertRepository
	log             log.ILogger
	now             func() time.Time
}

func New(
	recommendations ports.IRecommendationProvider,
	model ports.IAssistantModel,
	navigation ports.INavigationCatalog,
	store ports.IAssistantStore,
	recorder ports.IConversationRecorder,
	conversations ports.IConversationRepository,
	businessData ports.IBusinessDataReader,
	alerts ports.IAlertRepository,
	logger log.ILogger,
) IUseCase {
	return &UseCase{
		recommendations: recommendations,
		model:           model,
		navigation:      navigation,
		store:           store,
		recorder:        recorder,
		conversations:   conversations,
		businessData:    businessData,
		alerts:          alerts,
		log:             logger,
		now:             time.Now,
	}
}
