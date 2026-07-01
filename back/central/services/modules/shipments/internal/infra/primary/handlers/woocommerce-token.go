package handlers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

func deriveWooToken(secret string, integrationID uint, salt string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(fmt.Sprintf("wc-shipping-token:v1:%d:%s", integrationID, salt)))
	return hex.EncodeToString(mac.Sum(nil))[:32]
}

func (h *Handlers) wooTokenMatches(integrationID uint, salt, provided string) bool {
	if provided == "" || h.tokenSecret == "" {
		return false
	}
	expected := deriveWooToken(h.tokenSecret, integrationID, salt)
	return hmac.Equal([]byte(provided), []byte(expected))
}

func buildWooConnectionKey(baseURL string, integrationID uint, token string) string {
	payload := map[string]interface{}{
		"url":            baseURL,
		"integration_id": integrationID,
		"token":          token,
	}
	b, _ := json.Marshal(payload)
	return base64.RawURLEncoding.EncodeToString(b)
}
