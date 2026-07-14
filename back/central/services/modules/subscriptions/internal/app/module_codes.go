package app

import "github.com/secamc93/probability/back/central/shared/moduleregistry"

func (uc *UseCase) GetModuleCodes() []string {
	codes := make([]string, 0, len(moduleregistry.All))
	for _, m := range moduleregistry.All {
		codes = append(codes, string(m))
	}
	return codes
}
