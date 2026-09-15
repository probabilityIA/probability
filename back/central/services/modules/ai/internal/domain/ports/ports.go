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
