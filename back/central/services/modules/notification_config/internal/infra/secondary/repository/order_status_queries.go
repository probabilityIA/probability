package repository

import (
	"context"

	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/log"
)

// orderStatusQuerier consulta order_statuses (tabla de otro módulo)
// Replicado localmente para evitar compartir repositorios entre módulos
type orderStatusQuerier struct {
	db     db.IDatabase
	logger log.ILogger
}

// NewOrderStatusQuerier crea una instancia del querier de order statuses
func NewOrderStatusQuerier(database db.IDatabase, logger log.ILogger) *orderStatusQuerier {
	return &orderStatusQuerier{
		db:     database,
		logger: logger.WithModule("order_status_querier"),
	}
}

// GetOrderStatusCodesByIDs retorna un map de id→code para los IDs dados
// Consulta la tabla order_statuses (gestionada por módulo orderstatus)
func (q *orderStatusQuerier) GetOrderStatusCodesByIDs(ctx context.Context, ids []uint) (map[uint]string, error) {
	if len(ids) == 0 {
		return make(map[uint]string), nil
	}

	var results []struct {
		ID   uint
		Code string
	}

	err := q.db.Conn(ctx).
		Table("order_statuses").
		Select("id, code").
		Where("id IN ? AND deleted_at IS NULL", ids).
		Find(&results).Error

	if err != nil {
		q.logger.Error().
			Err(err).
			Interface("ids", ids).
			Msg("Error querying order status codes")
		return nil, err
	}

	codeMap := make(map[uint]string, len(results))
	for _, r := range results {
		codeMap[r.ID] = r.Code
	}

	return codeMap, nil
}
