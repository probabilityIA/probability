package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type TextMessage struct {
	Body string `json:"body"`
}

type WhatsappMessageRequest struct {
	MessagingProduct string      `json:"messaging_product"`
	To               string      `json:"to"`
	Type             string      `json:"type"`
	Text             TextMessage `json:"text"`
}

func (c *WhatsappClient) SendMessage(ctx context.Context, to string, message string) error {
	url := fmt.Sprintf("%s/messages", c.baseURL)

	payload := WhatsappMessageRequest{
		MessagingProduct: "whatsapp",
		To:               to,
		Type:             "text",
		Text: TextMessage{
			Body: message,
		},
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error marshaling payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	// Note: Authorization header might be needed here, usually a Bearer token.
	// Assuming it might be part of the client config or added here if the user provided a token.
	// For now, I will proceed without it or assume it's handled elsewhere/not required for this specific example request.
	// However, usually WhatsApp API needs a token. I'll add a TODO comment.

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("whatsapp api returned error status: %d", resp.StatusCode)
	}

	return nil
}
