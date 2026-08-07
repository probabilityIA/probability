package ports

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/commercial/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/commercial/internal/domain/entities"
)

type IRepository interface {
	ListProspects(ctx context.Context, filters dtos.ListProspectsFilters) ([]entities.Prospect, int64, error)
	GetStats(ctx context.Context) (*entities.ProspectStats, error)
}
