package fcm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/secamc93/probability/back/central/services/modules/push/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/push/internal/domain/errors"
)

type fcmRequest struct {
	Message fcmMessage `json:"message"`
}

type fcmMessage struct {
	Token        string            `json:"token"`
	Notification fcmNotification   `json:"notification"`
	Data         map[string]string `json:"data,omitempty"`
	Android      *fcmAndroid       `json:"android,omitempty"`
	APNS         *fcmAPNS          `json:"apns,omitempty"`
}

type fcmNotification struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

type fcmAndroid struct {
	Priority string `json:"priority"`
}

type fcmAPNS struct {
	Headers map[string]string `json:"headers"`
}

type fcmError struct {
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Status  string `json:"status"`
	} `json:"error"`
}

func (c *client) Send(ctx context.Context, tokens []string, message entities.PushMessage) (entities.SendResult, error) {
	result := entities.SendResult{}

	if !c.IsConfigured() {
		return result, domainerrors.ErrFCMNotConfigured
	}

	token, err := c.tokenSrc.Token()
	if err != nil {
		return result, fmt.Errorf("obteniendo token de acceso de FCM: %w", err)
	}

	url := fmt.Sprintf("https://fcm.googleapis.com/v1/projects/%s/messages:send", c.projectID)

	var lastTransient error
	for _, deviceToken := range tokens {
		body := fcmRequest{Message: fcmMessage{
			Token:        deviceToken,
			Notification: fcmNotification{Title: message.Title, Body: message.Body},
			Data:         message.Data,
			Android:      &fcmAndroid{Priority: "high"},
			APNS:         &fcmAPNS{Headers: map[string]string{"apns-priority": "10"}},
		}}

		payload, err := json.Marshal(body)
		if err != nil {
			return result, fmt.Errorf("serializando mensaje de FCM: %w", err)
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payload))
		if err != nil {
			return result, fmt.Errorf("armando request a FCM: %w", err)
		}
		req.Header.Set("Authorization", "Bearer "+token.AccessToken)
		req.Header.Set("Content-Type", "application/json")

		resp, err := c.httpClient.Do(req)
		if err != nil {
			lastTransient = fmt.Errorf("llamando a FCM: %w", err)
			continue
		}

		respBody, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		switch {
		case resp.StatusCode >= 200 && resp.StatusCode < 300:
			result.Sent++

		case isInvalidToken(resp.StatusCode, respBody):
			result.InvalidTokens = append(result.InvalidTokens, deviceToken)
			c.log.Warn(ctx).
				Str("motivo", string(respBody)).
				Msg("token de dispositivo invalido, se desactiva")

		default:
			lastTransient = fmt.Errorf("FCM respondio %d: %s", resp.StatusCode, string(respBody))
			c.log.Error(ctx).
				Int("status", resp.StatusCode).
				Str("respuesta", string(respBody)).
				Msg("error enviando push")
		}
	}

	result.TransientError = lastTransient
	return result, lastTransient
}

func isInvalidToken(status int, body []byte) bool {
	if status != http.StatusNotFound && status != http.StatusBadRequest && status != http.StatusForbidden {
		return false
	}

	var parsed fcmError
	if err := json.Unmarshal(body, &parsed); err != nil {
		return status == http.StatusNotFound
	}

	msg := strings.ToUpper(parsed.Error.Status + " " + parsed.Error.Message)
	return strings.Contains(msg, "UNREGISTERED") ||
		strings.Contains(msg, "NOT_FOUND") ||
		strings.Contains(msg, "INVALID_ARGUMENT")
}
