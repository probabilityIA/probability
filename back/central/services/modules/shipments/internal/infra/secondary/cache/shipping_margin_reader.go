package cache

import (
	"context"
	"encoding/json"
	stderrors "errors"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/domain"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/redis"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
)

const (
	keyPrefix = "shipping_margins:business:"
	ttl       = 24 * time.Hour
)

type cacheValue struct {
	MarginAmount    float64 `json:"margin_amount"`
	InsuranceMargin float64 `json:"insurance_margin"`
}

type ShippingMarginReader struct {
	rdb redis.IRedis
	db  db.IDatabase
	log log.ILogger
}

func NewShippingMarginReader(rdb redis.IRedis, database db.IDatabase, logger log.ILogger) domain.IShippingMarginReader {
	if rdb != nil {
		rdb.RegisterCachePrefix(keyPrefix)
	}
	return &ShippingMarginReader{rdb: rdb, db: database, log: logger}
}

func key(businessID uint) string {
	return fmt.Sprintf("%s%d", keyPrefix, businessID)
}

func (r *ShippingMarginReader) Get(ctx context.Context, businessID uint, carrierCode string) (domain.ShippingMargin, error) {
	if businessID == 0 || carrierCode == "" {
		return domain.ShippingMargin{}, nil
	}

	if r.rdb != nil {
		raw, err := r.rdb.HGet(ctx, key(businessID), carrierCode)
		if err == nil && raw != "" {
			var v cacheValue
			if err := json.Unmarshal([]byte(raw), &v); err == nil {
				return domain.ShippingMargin{MarginAmount: v.MarginAmount, InsuranceMargin: v.InsuranceMargin}, nil
			}
		}
	}

	if r.db == nil {
		return domain.ShippingMargin{}, nil
	}

	var model models.ShippingMargin
	err := r.db.Conn(ctx).
		Where("business_id = ? AND carrier_code = ? AND is_active = ?", businessID, carrierCode, true).
		First(&model).Error
	if err != nil {
		if stderrors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ShippingMargin{}, nil
		}
		return domain.ShippingMargin{}, err
	}

	margin := domain.ShippingMargin{
		MarginAmount:    model.MarginAmount,
		InsuranceMargin: model.InsuranceMargin,
	}

	if r.rdb != nil {
		payload, _ := json.Marshal(cacheValue{
			MarginAmount:    margin.MarginAmount,
			InsuranceMargin: margin.InsuranceMargin,
		})
		k := key(businessID)
		_ = r.rdb.HSet(ctx, k, carrierCode, string(payload))
		_ = r.rdb.Expire(ctx, k, ttl)
	}

	return margin, nil
}
