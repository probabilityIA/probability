package consumeralert

import (
	"strings"

	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/app/usecasemessaging"
)

const envAlertPhones = "MONITORING_ALERT_PHONES"

var defaultAlertPhones = []string{
	"573023406789",
	"573192611891",
}

func resolveAlertPhones(raw string) []string {
	candidatos := defaultAlertPhones
	if strings.TrimSpace(raw) != "" {
		candidatos = strings.Split(raw, ",")
	}

	vistos := make(map[string]struct{}, len(candidatos))
	phones := make([]string, 0, len(candidatos))
	for _, candidato := range candidatos {
		phone := usecasemessaging.NormalizePhoneNumber(candidato)
		if phone == "" {
			continue
		}
		if _, repetido := vistos[phone]; repetido {
			continue
		}
		vistos[phone] = struct{}{}
		phones = append(phones, phone)
	}

	if len(phones) == 0 {
		return defaultAlertPhones
	}
	return phones
}
