package cache

import (
	"context"
	"fmt"

	domain "github.com/secamc93/probability/back/central/services/modules/ai_sales/internal/domain"
	redisclient "github.com/secamc93/probability/back/central/shared/redis"
)

const aiPausedKeyPrefix = "whatsapp:ai_paused:"

type pauseChecker struct {
	redis redisclient.IRedis
}

// NewPauseChecker crea un checker que lee el flag de pausa del módulo WhatsApp.
// Usa la misma clave Redis: whatsapp:ai_paused:{phoneNumber}
func NewPauseChecker(redis redisclient.IRedis) domain.IAIPauseChecker {
	return &pauseChecker{redis: redis}
}

func (p *pauseChecker) IsAIPaused(ctx context.Context, phoneNumber string) bool {
	key := fmt.Sprintf("%s%s", aiPausedKeyPrefix, phoneNumber)
	_, err := p.redis.Get(ctx, key)
	return err == nil
}

var _ domain.IAIPauseChecker = (*pauseChecker)(nil)
