package ports

import (
	"context"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/orderscompare/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/shared/orderscompare"
)

type IRepository interface {
	GetIntegration(ctx context.Context, integrationID uint) (*Integration, error)
	ListLocalOrders(ctx context.Context, integrationID uint, businessID uint, from, to *time.Time, limit int) ([]dtos.LocalOrder, error)
}

type Integration struct {
	ID              uint
	Name            string
	BusinessID      *uint
	IntegrationType uint
	IsActive        bool
}

type IChannelRegistry interface {
	Supports(integrationTypeID uint) bool
	ListOrders(ctx context.Context, integrationTypeID uint, integrationID string, filters orderscompare.ChannelFilters) ([]orderscompare.ChannelOrder, error)
	ImportOrders(ctx context.Context, integrationTypeID uint, integrationID string, externalIDs []string) (orderscompare.ImportResult, error)
}

type IUseCase interface {
	Compare(ctx context.Context, query dtos.CompareQuery) (*dtos.ComparePage, error)
	Apply(ctx context.Context, cmd dtos.ApplyCommand) (*dtos.ApplyResult, error)
}
