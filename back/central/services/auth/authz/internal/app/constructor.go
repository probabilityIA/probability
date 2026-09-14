package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/auth/authz/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/auth/authz/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

type IUseCase interface {
	GetEffectiveAccess(ctx context.Context, userID, tokenBusinessID, requestedBusinessID uint) (*entities.Access, error)
	Can(access *entities.Access, resourceCode, actionCode string) bool
	Navigation(access *entities.Access) []entities.NavItem
	SetModuleAccess(modules ports.IModuleAccess)
}

type UseCase struct {
	repo    ports.IRepository
	modules ports.IModuleAccess
	cache   ports.IAccessCache
	log     log.ILogger
}

func New(repo ports.IRepository, modules ports.IModuleAccess, cache ports.IAccessCache, logger log.ILogger) IUseCase {
	return &UseCase{repo: repo, modules: modules, cache: cache, log: logger}
}

func (uc *UseCase) SetModuleAccess(modules ports.IModuleAccess) {
	uc.modules = modules
}
