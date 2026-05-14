package app

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/geozones/internal/domain/dtos"
)

const (
	displayPayloadVersion = "v7"
)

func toleranceForZoom(zoom int) float64 {
	switch {
	case zoom <= 7:
		return 0.0005
	case zoom <= 9:
		return 0.0001
	}
	return 0
}

func zoomBucket(zoom int) string {
	switch {
	case zoom <= 7:
		return "z7"
	case zoom <= 9:
		return "z9"
	}
	return "zfull"
}

func (u *UseCase) GetForDisplay(ctx context.Context, geozoneType string, zoom int, bbox *dtos.Bbox, parentID *uint) ([]byte, string, error) {
	t := geozoneType
	if t == "" {
		t = "all"
	}
	bucket := zoomBucket(zoom)
	tolerance := toleranceForZoom(zoom)

	useCache := bbox == nil
	parentSeg := "p0"
	if parentID != nil {
		parentSeg = fmt.Sprintf("p%d", *parentID)
	}
	cacheKey := fmt.Sprintf("%s:%s:%s:%s", t, parentSeg, bucket, displayPayloadVersion)
	etag := fmt.Sprintf(`"%s-%s-%s-%s"`, t, parentSeg, bucket, displayPayloadVersion)

	if useCache && u.cache != nil {
		if cached, ok := u.cache.Get(ctx, cacheKey); ok {
			return cached, etag, nil
		}
	}

	features, err := u.repo.GetForDisplay(ctx, dtos.DisplayParams{
		Type:      geozoneType,
		Tolerance: tolerance,
		Bbox:      bbox,
		ParentID:  parentID,
	})
	if err != nil {
		return nil, "", err
	}

	fc := dtos.DisplayFeatureCollection{
		Type:     "FeatureCollection",
		Features: features,
	}

	payload, err := json.Marshal(fc)
	if err != nil {
		return nil, "", err
	}

	if useCache && u.cache != nil {
		_ = u.cache.Set(ctx, cacheKey, payload)
	}

	return payload, etag, nil
}

func (u *UseCase) FlushDisplayCache(ctx context.Context) error {
	if u.cache == nil {
		return nil
	}
	return u.cache.FlushAll(ctx)
}
