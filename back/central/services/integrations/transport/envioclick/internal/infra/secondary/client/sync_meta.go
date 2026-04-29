package client

import (
	"encoding/json"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/secamc93/probability/back/central/services/integrations/transport/envioclick/internal/domain"
)

func captureMeta(meta *domain.SyncMeta, method, url string, requestBody any, started time.Time, resp *resty.Response, reqErr error) {
	if meta == nil {
		return
	}
	meta.Method = method
	meta.URL = url
	meta.StartedAt = started
	meta.CompletedAt = time.Now()
	meta.DurationMs = int(meta.CompletedAt.Sub(started) / time.Millisecond)

	if requestBody != nil {
		if b, err := json.Marshal(requestBody); err == nil {
			meta.RequestBody = b
		}
	}

	if resp != nil {
		meta.ResponseStatus = resp.StatusCode()
		if body := resp.Body(); len(body) > 0 {
			meta.ResponseBody = append(meta.ResponseBody[:0], body...)
		}
	}

	if reqErr != nil && len(meta.ResponseBody) == 0 {
		errBody, _ := json.Marshal(map[string]string{"error": reqErr.Error()})
		meta.ResponseBody = errBody
	}
}
