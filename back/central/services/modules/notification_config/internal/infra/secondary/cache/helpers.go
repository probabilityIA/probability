package cache

import (
	"context"
	"fmt"
)

func buildCacheKey(integrationID uint, notificationTypeID uint, notificationEventTypeID uint) string {
	return fmt.Sprintf("notification:configs:%d:%d:%d", integrationID, notificationTypeID, notificationEventTypeID)
}

func buildIndexKey(configID uint) string {
	return fmt.Sprintf("notification:config:%d:keys", configID)
}

func buildEventCodeCacheKey(integrationID uint, eventCode string) string {
	return fmt.Sprintf("notification:configs:evt:%d:%s", integrationID, eventCode)
}

func boolPtr(b bool) *bool {
	return &b
}

const assistantNotificationTypeID uint = 6

func buildBusinessEventCacheKey(businessID uint, eventCode string) string {
	return fmt.Sprintf("notification:configs:biz:%d:%s", businessID, eventCode)
}

func (c *cacheManager) cacheBusinessWide(ctx context.Context, notificationTypeID uint, businessID *uint, eventCode string, configID uint, configJSON string) {
	if notificationTypeID != assistantNotificationTypeID || businessID == nil || *businessID == 0 || eventCode == "" {
		return
	}
	key := buildBusinessEventCacheKey(*businessID, eventCode)
	if err := c.redis.HSet(ctx, key, fmt.Sprintf("%d", configID), configJSON); err != nil {
		c.logger.Warn(ctx).Err(err).Str("key", key).Uint("config_id", configID).Msg("Error cacheando regla del asistente por negocio")
		return
	}
	if err := c.redis.HSet(ctx, buildIndexKey(configID), key, "1"); err != nil {
		c.logger.Warn(ctx).Err(err).Uint("config_id", configID).Msg("Error actualizando indice inverso de la regla del asistente")
	}
}
