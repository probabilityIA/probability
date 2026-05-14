package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/geozones/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/geozones/internal/domain/ports"
)

type Resolver struct {
	repo ports.IRepository
}

func NewResolver(repo ports.IRepository) ports.IResolver {
	return &Resolver{repo: repo}
}

func (r *Resolver) Resolve(ctx context.Context, lat, lng float64, businessID uint) (*entities.GeozoneAncestors, error) {
	return r.repo.ResolveAncestors(ctx, lat, lng, businessID)
}
