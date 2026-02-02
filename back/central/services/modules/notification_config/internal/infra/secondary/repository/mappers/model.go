package mappers

import (
	"github.com/secamc93/probability/back/migration/shared/models"
)

// IntegrationNotificationConfigModel usa el modelo de BD desde migrations
// NUEVA ESTRUCTURA: Usa BusinessNotificationConfig con IDs de tablas normalizadas
type IntegrationNotificationConfigModel = models.BusinessNotificationConfig
