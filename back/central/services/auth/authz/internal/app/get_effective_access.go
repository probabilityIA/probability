package app

import (
	"context"
	"sort"

	"github.com/secamc93/probability/back/central/services/auth/authz/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/auth/authz/internal/domain/errors"
	"github.com/secamc93/probability/back/central/shared/authz"
	"github.com/secamc93/probability/back/central/shared/moduleregistry"
)

const platformScopeCode = "platform"

func (uc *UseCase) GetEffectiveAccess(ctx context.Context, userID, tokenBusinessID, requestedBusinessID uint) (*entities.Access, error) {
	if tokenBusinessID == 0 {
		return uc.superAdminAccess(ctx, userID, requestedBusinessID)
	}

	if uc.cache != nil {
		if cached, ok := uc.cache.Get(ctx, userID, tokenBusinessID); ok {
			return cached, nil
		}
	}

	businessID := tokenBusinessID
	staff, err := uc.repo.GetStaffRole(ctx, userID, &businessID)
	if err != nil {
		return nil, err
	}
	if staff == nil {
		return nil, domainerrors.ErrNoBusinessRelation
	}

	if staff.RoleScopeCode == platformScopeCode {
		return uc.superAdminAccess(ctx, userID, businessID)
	}

	access := &entities.Access{
		UserID:             userID,
		BusinessID:         businessID,
		BusinessName:       staff.BusinessName,
		SubscriptionStatus: staff.SubscriptionStatus,
		RoleName:           staff.RoleName,
		RoleCode:           authz.RoleCode(staff.RoleName),
	}
	if staff.RoleID != nil {
		access.RoleID = *staff.RoleID
	}

	modules, err := uc.modules.GetAccessibleModules(ctx, businessID)
	if err != nil {
		return nil, err
	}
	access.Modules = modules

	permissions, err := uc.resolvePermissions(ctx, staff, businessID, modules)
	if err != nil {
		return nil, err
	}
	access.Permissions = permissions

	if uc.cache != nil {
		uc.cache.Set(ctx, access)
	}
	return access, nil
}

func (uc *UseCase) resolvePermissions(ctx context.Context, staff *entities.StaffRole, businessID uint, modules []string) ([]string, error) {
	if staff.RoleID == nil {
		return []string{}, nil
	}

	rolePermissions, err := uc.repo.GetRolePermissions(ctx, *staff.RoleID)
	if err != nil {
		return nil, err
	}

	configs, err := uc.repo.GetBusinessResourceConfig(ctx, businessID)
	if err != nil {
		return nil, err
	}
	activeResources := make(map[string]bool, len(configs))
	for _, cfg := range configs {
		if resource, ok := authz.ResourceByLegacyName(cfg.ResourceName); ok && cfg.Active {
			activeResources[resource.Code] = true
		}
	}
	restrictByConfig := len(configs) > 0

	allowedModules := make(map[string]bool, len(modules))
	for _, m := range modules {
		allowedModules[m] = true
	}

	set := make(map[string]bool)
	for _, rp := range rolePermissions {
		resource, ok := authz.ResourceByLegacyName(rp.ResourceName)
		if !ok {
			uc.log.Warn(ctx).Str("resource", rp.ResourceName).Msg("[authz] recurso de BD sin codigo en el catalogo, se ignora")
			continue
		}
		action, ok := authz.ActionCode(rp.ActionName)
		if !ok {
			uc.log.Warn(ctx).Str("action", rp.ActionName).Msg("[authz] accion de BD sin codigo en el catalogo, se ignora")
			continue
		}
		if resource.SuperAdminOnly {
			continue
		}
		if restrictByConfig && !activeResources[resource.Code] {
			continue
		}
		if resource.Module != "" && !allowedModules[resource.Module] {
			continue
		}
		set[authz.PermissionCode(resource.Code, action)] = true
	}

	permissions := make([]string, 0, len(set))
	for code := range set {
		permissions = append(permissions, code)
	}
	sort.Strings(permissions)
	return permissions, nil
}

func (uc *UseCase) superAdminAccess(ctx context.Context, userID, requestedBusinessID uint) (*entities.Access, error) {
	staff, err := uc.repo.GetStaffRole(ctx, userID, nil)
	if err != nil {
		return nil, err
	}
	if staff == nil {
		return nil, domainerrors.ErrNoBusinessRelation
	}

	access := &entities.Access{
		UserID:             userID,
		BusinessID:         requestedBusinessID,
		IsSuper:            true,
		RoleName:           staff.RoleName,
		RoleCode:           authz.RoleCode(staff.RoleName),
		SubscriptionStatus: "active",
		Modules:            allModuleCodes(),
		Permissions:        []string{},
	}
	if staff.RoleID != nil {
		access.RoleID = *staff.RoleID
	}
	if requestedBusinessID > 0 {
		name, err := uc.repo.GetBusinessName(ctx, requestedBusinessID)
		if err != nil {
			return nil, err
		}
		access.BusinessName = name
	}
	return access, nil
}

func allModuleCodes() []string {
	codes := make([]string, 0, len(moduleregistry.All))
	for _, m := range moduleregistry.All {
		codes = append(codes, string(m))
	}
	return codes
}
