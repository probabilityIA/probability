package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/shipping_margins/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/shipping_margins/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/redis"
)

const (
	keyPrefix = "shipping_margins:business:"
	ttl       = 24 * time.Hour
)

type cacheValue struct {
	MarginAmount    float64 `json:"margin_amount"`
	InsuranceMargin float64 `json:"insurance_margin"`
}

type Writer struct {
	rdb redis.IRedis
	log log.ILogger
}

func New(rdb redis.IRedis, logger log.ILogger) ports.ICacheWriter {
	if rdb != nil {
		rdb.RegisterCachePrefix(keyPrefix)
	}
	return &Writer{rdb: rdb, log: logger}
}

func key(businessID uint) string {
	return fmt.Sprintf("%s%d", keyPrefix, businessID)
}

func (w *Writer) Upsert(ctx context.Context, m *entities.ShippingMargin) error {
	if w.rdb == nil || m == nil {
		return nil
	}
	payload, err := json.Marshal(cacheValue{
		MarginAmount:    m.MarginAmount,
		InsuranceMargin: m.InsuranceMargin,
	})
	if err != nil {
		return err
	}
	k := key(m.BusinessID)
	if err := w.rdb.HSet(ctx, k, m.CarrierCode, string(payload)); err != nil {
		w.log.Warn(ctx).Err(err).Str("key", k).Str("carrier", m.CarrierCode).Msg("shipping_margins cache upsert failed")
		return err
	}
	if err := w.rdb.Expire(ctx, k, ttl); err != nil {
		w.log.Warn(ctx).Err(err).Str("key", k).Msg("shipping_margins cache expire failed")
	}
	return nil
}

func (w *Writer) Delete(ctx context.Context, businessID uint, carrierCode string) error {
	if w.rdb == nil {
		return nil
	}
	k := key(businessID)
	if err := w.rdb.HDel(ctx, k, carrierCode); err != nil {
		w.log.Warn(ctx).Err(err).Str("key", k).Str("carrier", carrierCode).Msg("shipping_margins cache delete failed")
		return err
	}
	return nil
}
