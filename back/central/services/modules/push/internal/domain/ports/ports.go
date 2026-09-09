package ports

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/push/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/push/internal/domain/entities"
)

type IRepository interface {
	UpsertDeviceToken(ctx context.Context, dto dtos.RegisterDeviceDTO) error
	DeactivateToken(ctx context.Context, userID uint, token string) error
	DeactivateTokens(ctx context.Context, tokens []string) error
	ListActiveTokensByBusiness(ctx context.Context, businessID uint) ([]entities.DeviceToken, error)
	ListDevicesByUser(ctx context.Context, userID uint) ([]entities.DeviceToken, error)
	UserBelongsToBusiness(ctx context.Context, userID, businessID uint) (bool, error)
}

type IPushSender interface {
	Send(ctx context.Context, tokens []string, message entities.PushMessage) (entities.SendResult, error)
	IsConfigured() bool
}
