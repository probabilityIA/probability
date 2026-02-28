package cache

import "fmt"

// buildCacheKey construye la key de Redis para una configuración
// NUEVA ESTRUCTURA: Format: notification:configs:{integration_id}:{notification_type_id}:{notification_event_type_id}
func buildCacheKey(integrationID uint, notificationTypeID uint, notificationEventTypeID uint) string {
	return fmt.Sprintf("notification:configs:%d:%d:%d", integrationID, notificationTypeID, notificationEventTypeID)
}

// buildIndexKey construye la key del índice inverso
// Format: notification:config:{config_id}:keys
func buildIndexKey(configID uint) string {
	return fmt.Sprintf("notification:config:%d:keys", configID)
}

// buildEventCodeCacheKey construye la key secundaria de Redis para lookup por event code
// Format: notification:configs:evt:{integration_id}:{event_code}
func buildEventCodeCacheKey(integrationID uint, eventCode string) string {
	return fmt.Sprintf("notification:configs:evt:%d:%s", integrationID, eventCode)
}

// boolPtr es un helper para crear punteros a bool
func boolPtr(b bool) *bool {
	return &b
}
