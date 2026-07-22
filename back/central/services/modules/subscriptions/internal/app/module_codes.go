package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/shared/moduleregistry"
)

func (uc *UseCase) GetModuleCodes() []string {
	codes := make([]string, 0, len(moduleregistry.All))
	for _, m := range moduleregistry.All {
		codes = append(codes, string(m))
	}
	return codes
}

func (uc *UseCase) GetModuleCatalog() []dtos.ModuleInfo {
	catalog := make([]dtos.ModuleInfo, 0, len(moduleregistry.All))
	for _, m := range moduleregistry.All {
		catalog = append(catalog, dtos.ModuleInfo{
			Code: string(m),
			Name: moduleregistry.DisplayName(string(m)),
		})
	}
	return catalog
}

func (uc *UseCase) GetAccessibleModules(ctx context.Context, businessID uint) ([]string, error) {
	accessible := make([]string, 0, len(moduleregistry.All))
	for _, m := range moduleregistry.All {
		ok, err := uc.HasModuleAccess(ctx, businessID, string(m))
		if err != nil {
			return nil, err
		}
		if ok {
			accessible = append(accessible, string(m))
		}
	}
	return accessible, nil
}
