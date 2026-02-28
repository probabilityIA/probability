package cache

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/redis"
)

// IOrderStatusQuerier interfaz simple para consultar order_statuses
// Se replica localmente para evitar compartir repositorios entre módulos
type IOrderStatusQuerier interface {
	GetOrderStatusCodesByIDs(ctx context.Context, ids []uint) (map[uint]string, error)
}

// cacheManager implementa la interfaz ICacheManager
type cacheManager struct {
	redis              redis.IRedis
	repo               ports.IRepository
	orderStatusQuerier IOrderStatusQuerier
	logger             log.ILogger
}
