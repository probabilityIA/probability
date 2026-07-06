package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/woostore/internal/domain"
)

type Iapp interface {
	Start(ctx context.Context) (*domain.PowerState, error)
	Stop(ctx context.Context) (*domain.PowerState, error)
	Status(ctx context.Context) (*domain.PowerState, error)
}

type UseCase struct {
	mgr domain.IPowerManager
}

func New(mgr domain.IPowerManager) Iapp {
	return &UseCase{mgr: mgr}
}

func (u *UseCase) Start(ctx context.Context) (*domain.PowerState, error) {
	return u.mgr.Start(ctx)
}

func (u *UseCase) Stop(ctx context.Context) (*domain.PowerState, error) {
	return u.mgr.Stop(ctx)
}

func (u *UseCase) Status(ctx context.Context) (*domain.PowerState, error) {
	return u.mgr.Status(ctx)
}
