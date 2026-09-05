package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
	redisclient "github.com/secamc93/probability/back/central/shared/redis"
)

const (
	templatesSnapshotTTL = 6 * time.Hour
	wabaIndexTTL         = 0
)

type templatesCache struct {
	redis redisclient.IRedis
	log   log.ILogger
}

func newTemplatesCache(redis redisclient.IRedis, logger log.ILogger) ports.ITemplateStatusCache {
	return &templatesCache{
		redis: redis,
		log:   logger.WithModule("whatsapp-templates-cache"),
	}
}

func templatesSnapshotKey(integrationID uint) string {
	return fmt.Sprintf("whatsapp:templates:%d", integrationID)
}

func wabaIndexKey(wabaID string) string {
	return "whatsapp:waba:" + wabaID
}

func (c *templatesCache) Get(ctx context.Context, integrationID uint) (*ports.WABATemplatesSnapshot, error) {
	if c.redis == nil {
		return nil, fmt.Errorf("redis no disponible")
	}

	raw, err := c.redis.Get(ctx, templatesSnapshotKey(integrationID))
	if err != nil {
		return nil, err
	}

	var snapshot ports.WABATemplatesSnapshot
	if err := json.Unmarshal([]byte(raw), &snapshot); err != nil {
		return nil, err
	}

	return &snapshot, nil
}

func (c *templatesCache) Save(ctx context.Context, snapshot *ports.WABATemplatesSnapshot) error {
	if c.redis == nil || snapshot == nil {
		return nil
	}

	data, err := json.Marshal(snapshot)
	if err != nil {
		return err
	}

	if err := c.redis.Set(ctx, templatesSnapshotKey(snapshot.IntegrationID), string(data), templatesSnapshotTTL); err != nil {
		return err
	}

	if snapshot.WABAID != "" {
		if err := c.IndexWABA(ctx, snapshot.WABAID, snapshot.IntegrationID); err != nil {
			c.log.Warn(ctx).Err(err).Str("waba_id", snapshot.WABAID).Msg("no se pudo indexar el WABA")
		}
	}

	return nil
}

func (c *templatesCache) IndexWABA(ctx context.Context, wabaID string, integrationID uint) error {
	if c.redis == nil {
		return nil
	}
	wabaID = strings.TrimSpace(wabaID)
	if wabaID == "" || integrationID == 0 {
		return nil
	}
	return c.redis.Set(ctx, wabaIndexKey(wabaID), strconv.Itoa(int(integrationID)), wabaIndexTTL)
}

func (c *templatesCache) IntegrationByWABA(ctx context.Context, wabaID string) (uint, error) {
	if c.redis == nil {
		return 0, fmt.Errorf("redis no disponible")
	}
	wabaID = strings.TrimSpace(wabaID)
	if wabaID == "" {
		return 0, nil
	}

	raw, err := c.redis.Get(ctx, wabaIndexKey(wabaID))
	if err != nil {
		return 0, err
	}

	id, err := strconv.ParseUint(raw, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}

func (c *templatesCache) UpdateStatusByWABA(ctx context.Context, wabaID, name, language, status, reason string) error {
	integrationID, err := c.IntegrationByWABA(ctx, wabaID)
	if err != nil {
		return err
	}
	if integrationID == 0 {
		return nil
	}

	snapshot, err := c.Get(ctx, integrationID)
	if err != nil || snapshot == nil {
		return err
	}

	updated := false
	for i := range snapshot.Templates {
		if snapshot.Templates[i].Name != name {
			continue
		}
		if language != "" && snapshot.Templates[i].Language != language {
			continue
		}
		snapshot.Templates[i].Status = status
		snapshot.Templates[i].Reason = reason
		snapshot.Templates[i].UpdatedAt = time.Now()
		updated = true
	}

	if !updated {
		snapshot.Templates = append(snapshot.Templates, ports.TemplateStatus{
			Name:      name,
			Language:  language,
			Status:    status,
			Reason:    reason,
			UpdatedAt: time.Now(),
		})
	}

	return c.Save(ctx, snapshot)
}
