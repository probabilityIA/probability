package encryption

import (
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/domain"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
)

// New crea una nueva instancia del servicio de encriptación
func New(config env.IConfig, logger log.ILogger) domain.IEncryptionService {
	return newEncryptionService(config, logger)
}
