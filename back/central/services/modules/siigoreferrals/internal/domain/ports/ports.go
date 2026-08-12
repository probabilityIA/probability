package ports

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/siigoreferrals/internal/domain/entities"
)

type IRepository interface {
	CreateReferral(ctx context.Context, referral *entities.SiigoReferral) error
	ListReferrals(ctx context.Context, page, pageSize int) ([]entities.SiigoReferral, int64, error)
}
