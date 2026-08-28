package ports

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/tours/internal/domain/entities"
)

type IRepository interface {
	ListByUser(ctx context.Context, userID, businessID uint) ([]entities.TourProgress, error)
	Upsert(ctx context.Context, progress entities.TourProgress) (*entities.TourProgress, error)
	UpsertMany(ctx context.Context, items []entities.TourProgress) error
	Delete(ctx context.Context, userID, businessID uint, tourKey string) error
	DeleteAll(ctx context.Context, userID, businessID uint) error
}
