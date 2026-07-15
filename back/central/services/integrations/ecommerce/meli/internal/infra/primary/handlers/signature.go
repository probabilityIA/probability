package handlers

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"strings"

	"github.com/gin-gonic/gin"
)

func (h *meliHandler) verifyNotificationSignature(ctx context.Context, c *gin.Context, resource string) {
	secret := h.config.Get("MELI_WEBHOOK_SECRET")
	if secret == "" {
		return
	}

	sigHeader := c.GetHeader("x-signature")
	if sigHeader == "" {
		h.logger.Warn(ctx).Msg("MercadoLibre notification without x-signature header")
		return
	}

	ts, v0 := parseSignatureHeader(sigHeader)
	if ts == "" || v0 == "" {
		h.logger.Warn(ctx).Str("x_signature", sigHeader).Msg("Malformed x-signature header")
		return
	}

	requestID := c.GetHeader("x-request-id")
	id := lastPathSegment(resource)

	manifest := "id:" + id + ";request-id:" + requestID + ";ts:" + ts + ";"
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(manifest))
	expected := hex.EncodeToString(mac.Sum(nil))

	if !hmac.Equal([]byte(expected), []byte(v0)) {
		h.logger.Warn(ctx).
			Str("resource", resource).
			Msg("MercadoLibre notification signature mismatch")
	}
}

func parseSignatureHeader(header string) (ts string, v0 string) {
	for _, part := range strings.Split(header, ",") {
		kv := strings.SplitN(strings.TrimSpace(part), "=", 2)
		if len(kv) != 2 {
			continue
		}
		switch kv[0] {
		case "ts":
			ts = kv[1]
		case "v0":
			v0 = kv[1]
		}
	}
	return ts, v0
}

func lastPathSegment(resource string) string {
	trimmed := strings.Trim(resource, "/")
	if idx := strings.IndexByte(trimmed, '?'); idx >= 0 {
		trimmed = trimmed[:idx]
	}
	parts := strings.Split(trimmed, "/")
	return parts[len(parts)-1]
}
