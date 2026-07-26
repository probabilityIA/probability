package dtos

import "github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/entities"

// CreateNotificationConfigDTO representa los datos para crear una configuración
// NUEVA ESTRUCTURA: Usa IDs de tablas
type CreateNotificationConfigDTO struct {
	BusinessID              *uint  // Opcional - auto-asignado por middleware
	IntegrationID           uint   // Integración origen
	NotificationTypeID      uint   // Canal de salida (ID de notification_types)
	NotificationEventTypeID uint   // Tipo de evento (ID de notification_event_types)
	Enabled                 bool   // Estado de la configuración
	Description             string // Descripción opcional
	OrderStatusIDs          []uint // Filtro opcional de estados
	CODOnly                 bool   // Solo ordenes contra entrega
}

// ============================================
// LEGACY: Estructura anterior (DEPRECADA)
// Se mantiene temporalmente para compatibilidad
// ============================================

type CreateNotificationConfigDTOLegacy struct {
	IntegrationID    uint
	NotificationType string
	IsActive         bool
	Conditions       entities.NotificationConditions
	Config           entities.NotificationConfig
	Description      string
	Priority         int
}
