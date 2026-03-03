package types

// WebhookPayload represents a ready-to-send HTTP request payload.
// The backend builds these and the frontend sends them,
// so the frontend can display request/response details.
type WebhookPayload struct {
	URL     string            `json:"url"`
	Method  string            `json:"method"`
	Headers map[string]string `json:"headers"`
	Body    interface{}       `json:"body"`               // parsed body for display
	RawBody string            `json:"raw_body,omitempty"` // exact JSON bytes to send (preserves HMAC signature)
}
