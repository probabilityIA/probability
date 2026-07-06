package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/domain"
)

type wooResolved struct {
	Found      bool                  `json:"found"`
	Salt       string                `json:"salt"`
	Revoked    bool                  `json:"revoked"`
	BusinessID uint                  `json:"business_id"`
	Carrier    *domain.CarrierInfo   `json:"carrier"`
	Origin     *domain.OriginAddress `json:"origin"`
}

func wooResKey(integrationID uint) string {
	return fmt.Sprintf("woores:%d", integrationID)
}

func (h *Handlers) resolveWoo(ctx context.Context, integrationID uint) (*wooResolved, error) {
	key := wooResKey(integrationID)

	if h.redisClient != nil {
		if v, err := h.redisClient.Get(ctx, key); err == nil && v != "" {
			var r wooResolved
			if json.Unmarshal([]byte(v), &r) == nil {
				return &r, nil
			}
		}
	}

	salt, revoked, found, err := h.uc.Repo().GetWooShippingToken(ctx, integrationID)
	if err != nil {
		return nil, err
	}

	r := &wooResolved{Found: found, Salt: salt, Revoked: revoked}

	if bid, berr := h.uc.Repo().GetIntegrationBusinessID(ctx, integrationID); berr == nil && bid > 0 {
		r.BusinessID = bid
		if carrier, cerr := h.carrierResolver.GetActiveShippingCarrier(ctx, bid); cerr == nil {
			r.Carrier = carrier
		}
		if origin, oerr := h.uc.Repo().GetDefaultOriginAddress(ctx, bid); oerr == nil {
			r.Origin = origin
		}
	}

	if h.redisClient != nil {
		if b, mErr := json.Marshal(r); mErr == nil {
			_ = h.redisClient.Set(ctx, key, string(b), 60*time.Second)
		}
	}

	return r, nil
}

func (h *Handlers) bustWooCache(ctx context.Context, integrationID uint) {
	if h.redisClient != nil {
		_ = h.redisClient.Delete(ctx, wooResKey(integrationID))
	}
}

func (h *Handlers) daneCached(ctx context.Context, city, province string) string {
	if city == "" {
		return ""
	}
	key := "woodane:" + province + ":" + city

	if h.redisClient != nil {
		if v, err := h.redisClient.Get(ctx, key); err == nil && v != "" {
			return v
		}
	}

	dane, err := h.uc.Repo().GetCityDaneByName(ctx, city, province)
	if err != nil || dane == "" {
		return ""
	}

	if h.redisClient != nil {
		_ = h.redisClient.Set(ctx, key, dane, time.Hour)
	}
	return dane
}
