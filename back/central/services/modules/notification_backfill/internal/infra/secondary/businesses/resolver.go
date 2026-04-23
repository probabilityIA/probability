package businesses

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/notification_backfill/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/db"
)

type resolver struct {
	db db.IDatabase
}

func NewResolver(database db.IDatabase) ports.IBusinessNameResolver {
	return &resolver{db: database}
}

func (r *resolver) ResolveNames(ctx context.Context, ids []uint) (map[uint]string, error) {
	out := make(map[uint]string, len(ids))
	if len(ids) == 0 {
		return out, nil
	}

	type row struct {
		ID   uint
		Name string
	}

	var rows []row
	if err := r.db.Conn(ctx).
		Table("business").
		Select("id, name").
		Where("id IN ?", ids).
		Scan(&rows).Error; err != nil {
		return out, err
	}

	for _, r := range rows {
		out[r.ID] = r.Name
	}
	return out, nil
}
