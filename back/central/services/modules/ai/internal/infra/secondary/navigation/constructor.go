package navigation

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/ai/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/authz"
)

type VisibleNavigationFunc = func(ctx context.Context, userID, tokenBusinessID, requestedBusinessID uint) ([]authz.NavEntry, error)

type Catalog struct {
	visible VisibleNavigationFunc
	guides  map[string]string
}

func New(visible VisibleNavigationFunc) ports.INavigationCatalog {
	return &Catalog{visible: visible, guides: parseGuide(guideSource)}
}
