package cache

import "fmt"

// buildCacheKey construye la key de Redis para una configuración
// Format: notification:configs:{integration_id}:{trigger}
func buildCacheKey(integrationID uint, trigger string) string {
	return fmt.Sprintf("notification:configs:%d:%s", integrationID, trigger)
}

// buildIndexKey construye la key del índice inverso
// Format: notification:config:{config_id}:keys
func buildIndexKey(configID uint) string {
	return fmt.Sprintf("notification:config:%d:keys", configID)
}

// boolPtr es un helper para crear punteros a bool
func boolPtr(b bool) *bool {
	return &b
}
