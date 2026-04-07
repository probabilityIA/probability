package cache

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/redis"
)

const aiPausedPrefix = "whatsapp:ai_paused:"

type aiPauseChecker struct {
	redis redis.IRedis
}

// NewAIPauseChecker crea un checker que lee la clave Redis whatsapp:ai_paused:{phone}
// (gestionada por el módulo whatsapp) para saber si la IA está pausada.
func NewAIPauseChecker(redis redis.IRedis) ports.IAIPauseChecker {
	return &aiPauseChecker{redis: redis}
}

func (c *aiPauseChecker) IsAIPaused(ctx context.Context, phoneNumber string) bool {
	count, err := c.redis.Exists(ctx, aiPausedPrefix+phoneNumber)
	return err == nil && count > 0
}
