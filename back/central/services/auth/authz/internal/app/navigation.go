package app

import (
	"github.com/secamc93/probability/back/central/services/auth/authz/internal/domain/entities"
	"github.com/secamc93/probability/back/central/shared/authz"
)

func (uc *UseCase) Can(access *entities.Access, resourceCode, actionCode string) bool {
	if access == nil {
		return false
	}
	if access.IsSuper {
		return true
	}
	want := authz.PermissionCode(resourceCode, actionCode)
	for _, p := range access.Permissions {
		if p == want {
			return true
		}
	}
	return false
}

func (uc *UseCase) Navigation(access *entities.Access) []entities.NavItem {
	items := make([]entities.NavItem, 0, len(authz.Navigation))
	if access == nil {
		return items
	}
	for _, entry := range authz.Navigation {
		if !uc.navVisible(access, entry) {
			continue
		}
		items = append(items, entities.NavItem{
			Key:         entry.Key,
			Label:       entry.Label,
			Route:       entry.Route,
			Section:     entry.Section,
			Description: entry.Description,
		})
	}
	return items
}

func (uc *UseCase) navVisible(access *entities.Access, entry authz.NavEntry) bool {
	if contains(entry.ExcludeRoles, access.RoleCode) {
		return false
	}
	if !access.IsSuper && entry.Module != "" && !contains(access.Modules, entry.Module) {
		return false
	}
	switch entry.Rule {
	case authz.NavRuleAlways:
		return true
	case authz.NavRuleSuperAdminOnly:
		return access.IsSuper
	case authz.NavRuleRole:
		return access.IsSuper || contains(entry.RoleCodes, access.RoleCode)
	case authz.NavRulePermission:
		return uc.Can(access, entry.Resource, authz.ActionRead)
	default:
		return false
	}
}

func contains(list []string, value string) bool {
	for _, v := range list {
		if v == value {
			return true
		}
	}
	return false
}
