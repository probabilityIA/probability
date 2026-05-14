package cache

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"io"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/geozones/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/redis"
)

const (
	keyPrefix = "geozones:display:"
	ttl       = 30 * 24 * time.Hour
)

type Cache struct {
	rdb redis.IRedis
	log log.ILogger
}

func New(rdb redis.IRedis, logger log.ILogger) ports.IDisplayCache {
	if rdb != nil {
		rdb.RegisterCachePrefix(keyPrefix)
	}
	return &Cache{rdb: rdb, log: logger}
}

func (c *Cache) Get(ctx context.Context, key string) ([]byte, bool) {
	if c.rdb == nil {
		return nil, false
	}
	val, err := c.rdb.Get(ctx, keyPrefix+key)
	if err != nil || val == "" {
		return nil, false
	}
	raw, err := base64.StdEncoding.DecodeString(val)
	if err != nil {
		return nil, false
	}
	gz, err := gzip.NewReader(bytes.NewReader(raw))
	if err != nil {
		return nil, false
	}
	defer gz.Close()
	out, err := io.ReadAll(gz)
	if err != nil {
		return nil, false
	}
	return out, true
}

func (c *Cache) Set(ctx context.Context, key string, value []byte) error {
	if c.rdb == nil {
		return nil
	}
	var buf bytes.Buffer
	gz, _ := gzip.NewWriterLevel(&buf, gzip.BestCompression)
	if _, err := gz.Write(value); err != nil {
		return err
	}
	if err := gz.Close(); err != nil {
		return err
	}
	encoded := base64.StdEncoding.EncodeToString(buf.Bytes())
	if err := c.rdb.Set(ctx, keyPrefix+key, encoded, ttl); err != nil {
		c.log.Warn(ctx).Err(err).Str("key", key).Msg("geozones display cache set failed")
		return err
	}
	return nil
}

func (c *Cache) FlushAll(ctx context.Context) error {
	if c.rdb == nil {
		return nil
	}
	keys, err := c.rdb.Keys(ctx, keyPrefix+"*")
	if err != nil {
		return err
	}
	if len(keys) == 0 {
		return nil
	}
	return c.rdb.Delete(ctx, keys...)
}
