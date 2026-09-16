package ports

import (
	"context"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/ai/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/ai/internal/domain/entities"
)

type IRecommendationProvider interface {
	GetRecommendation(ctx context.Context, origin, destination string) (*entities.Recommendation, error)
}

type IAssistantModel interface {
	Reply(ctx context.Context, req dtos.ModelRequest) (*dtos.ModelReply, error)
}

type INavigationCatalog interface {
	ForUser(ctx context.Context, scope dtos.AccessScope) (*entities.NavigationCatalog, error)
}

type IAssistantStore interface {
	ConsumeMessage(ctx context.Context, userID uint, window time.Duration) (*entities.Usage, error)
	GetUsage(ctx context.Context, userID uint) (*entities.Usage, error)
	IsIntroSeen(ctx context.Context, userID uint) (bool, error)
	MarkIntroSeen(ctx context.Context, userID uint) error
}

type IBusinessDataReader interface {
	FindOrders(ctx context.Context, businessID uint, number string) ([]entities.OrderInfo, error)
	FindShipments(ctx context.Context, businessID uint, trackingNumber string) ([]entities.ShipmentInfo, error)
	ListOrders(ctx context.Context, businessID uint, query dtos.OrderQuery) ([]entities.OrderSummary, int64, error)
	SummarizeOrders(ctx context.Context, businessID uint, from, to time.Time) (*entities.OrdersOverview, error)
	DescribeIdentity(ctx context.Context, userID uint, businessID *uint) (*entities.ChatIdentity, error)
	CountUnreadWhatsAppChats(ctx context.Context) ([]entities.UnreadChats, error)
	FindCustomerNameByPhone(ctx context.Context, businessID uint, phone string) (string, error)
}

type IConversationRecorder interface {
	RecordMessage(ctx context.Context, record entities.MessageRecord) error
	RecordFeedback(ctx context.Context, userID uint, messageID string, value int, at time.Time) error
	RecordClick(ctx context.Context, userID uint, messageID string, at time.Time) error
}

type IConversationRepository interface {
	SaveMessage(ctx context.Context, record entities.MessageRecord) error
	ApplyFeedback(ctx context.Context, userID uint, messageID string, value int, at time.Time) (bool, error)
	ApplyClick(ctx context.Context, userID uint, messageID string, at time.Time) (bool, error)
	ListMessages(ctx context.Context, filter dtos.ReviewFilter) ([]entities.ReviewMessage, int64, error)
	Summary(ctx context.Context, filter dtos.ReviewFilter) (*entities.ReviewSummary, error)
	DeleteOlderThan(ctx context.Context, cutoff time.Time) (int64, error)
}

type IAlertRepository interface {
	SaveAlert(ctx context.Context, alert entities.Alert) (bool, error)
	ListAlerts(ctx context.Context, query dtos.AlertQuery) ([]entities.Alert, int64, error)
	CountUnread(ctx context.Context, businessID, userID uint) (int64, *entities.Alert, error)
	MarkSeen(ctx context.Context, businessID, userID uint, at time.Time) error
	RecentAlerts(ctx context.Context, businessID uint, limit int) ([]entities.Alert, error)
	DeleteAlertsOlderThan(ctx context.Context, cutoff time.Time) (int64, error)
	LastAlertOfType(ctx context.Context, businessID uint, eventType string) (*entities.Alert, error)
}
