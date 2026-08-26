package domain

import "strings"

type CarrierSetting struct {
	Code         string `json:"code"`
	Enabled      bool   `json:"enabled"`
	AllowCOD     bool   `json:"allow_cod"`
	AllowPrepaid bool   `json:"allow_prepaid"`
}

func NormalizeCarrierCode(name string) string {
	return strings.ToUpper(strings.TrimSpace(name))
}

func CarrierAllowed(settings []CarrierSetting, carrierName string, isCOD bool) bool {
	if len(settings) == 0 {
		return true
	}
	code := NormalizeCarrierCode(carrierName)
	if code == "" {
		return true
	}
	for _, s := range settings {
		if NormalizeCarrierCode(s.Code) != code {
			continue
		}
		if !s.Enabled {
			return false
		}
		if isCOD {
			return s.AllowCOD
		}
		return s.AllowPrepaid
	}
	return true
}
