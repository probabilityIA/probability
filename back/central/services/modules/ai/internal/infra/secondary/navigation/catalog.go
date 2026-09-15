package navigation

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/ai/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/ai/internal/domain/entities"
	"github.com/secamc93/probability/back/central/shared/authz"
)

func (c *Catalog) ForUser(ctx context.Context, scope dtos.AccessScope) (*entities.NavigationCatalog, error) {
	if c.visible == nil {
		return nil, fmt.Errorf("navegacion no disponible")
	}

	entries, err := c.visible(ctx, scope.UserID, scope.TokenBusinessID, scope.RequestedBusinessID)
	if err != nil {
		return nil, fmt.Errorf("resolver navegacion visible: %w", err)
	}

	visible := make(map[string]bool, len(entries))
	for _, entry := range entries {
		visible[entry.Key] = true
	}

	catalog := &entities.NavigationCatalog{}
	for _, entry := range entries {
		catalog.Allowed = append(catalog.Allowed, entities.Destination{
			Key:         entry.Key,
			Label:       entry.Label,
			Route:       entry.Route,
			Description: entry.Description,
			Guide:       c.guides[entry.Key],
		})
		for _, sub := range subpages {
			if sub.Parent != entry.Key || (sub.AlsoRequires != "" && !visible[sub.AlsoRequires]) {
				continue
			}
			destination := entities.Destination{
				Key:         sub.Key,
				Label:       sub.Label,
				Route:       sub.Route,
				Description: sub.Description,
				Guide:       c.guides[sub.Key],
			}
			if sub.Highlight != nil {
				destination.Highlight = &entities.Highlight{
					Target: sub.Highlight.Target,
					Title:  sub.Highlight.Title,
					Hint:   sub.Highlight.Hint,
				}
			}
			catalog.Allowed = append(catalog.Allowed, destination)
		}
	}

	for _, entry := range authz.Navigation {
		if visible[entry.Key] || entry.Rule == authz.NavRuleSuperAdminOnly {
			continue
		}
		catalog.Denied = append(catalog.Denied, entry.Label)
	}
	return catalog, nil
}
