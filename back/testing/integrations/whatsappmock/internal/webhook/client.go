package webhook

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/secamc93/probability/back/testing/shared/log"
)

type Payload struct {
	Object string  `json:"object"`
	Entry  []Entry `json:"entry"`
}

type Entry struct {
	ID      string   `json:"id"`
	Changes []Change `json:"changes"`
}

type Change struct {
	Value Value  `json:"value"`
	Field string `json:"field"`
}

type Value struct {
	MessagingProduct string    `json:"messaging_product,omitempty"`
	Metadata         *Metadata `json:"metadata,omitempty"`
	Contacts         []Contact `json:"contacts,omitempty"`
	Messages         []Message `json:"messages,omitempty"`
	Statuses         []Status  `json:"statuses,omitempty"`

	Event                   string `json:"event,omitempty"`
	MessageTemplateName     string `json:"message_template_name,omitempty"`
	MessageTemplateLanguage string `json:"message_template_language,omitempty"`
	Reason                  string `json:"reason,omitempty"`
}

type Metadata struct {
	DisplayPhoneNumber string `json:"display_phone_number"`
	PhoneNumberID      string `json:"phone_number_id"`
}

type Contact struct {
	Profile Profile `json:"profile"`
	WaID    string  `json:"wa_id"`
}

type Profile struct {
	Name string `json:"name"`
}

type Message struct {
	From      string   `json:"from"`
	ID        string   `json:"id"`
	Timestamp string   `json:"timestamp"`
	Type      string   `json:"type"`
	Button    *Button  `json:"button,omitempty"`
	Context   *Context `json:"context,omitempty"`
	Text      *Text    `json:"text,omitempty"`
	Image     *Media   `json:"image,omitempty"`
	Document  *Media   `json:"document,omitempty"`
	Audio     *Media   `json:"audio,omitempty"`
	Video     *Media   `json:"video,omitempty"`
}

type Text struct {
	Body string `json:"body"`
}

type Media struct {
	ID       string `json:"id"`
	MimeType string `json:"mime_type"`
	SHA256   string `json:"sha256"`
	Caption  string `json:"caption,omitempty"`
	Filename string `json:"filename,omitempty"`
}

type Button struct {
	Payload string `json:"payload"`
	Text    string `json:"text"`
}

type Context struct {
	From string `json:"from"`
	ID   string `json:"id"`
}

type Status struct {
	ID          string `json:"id"`
	Status      string `json:"status"`
	Timestamp   string `json:"timestamp"`
	RecipientID string `json:"recipient_id"`
}

type Client struct {
	baseURL    string
	secret     string
	logger     log.ILogger
	httpClient *http.Client
}

func New(baseURL, secret string, logger log.ILogger) *Client {
	return &Client{
		baseURL:    strings.TrimRight(baseURL, "/"),
		secret:     secret,
		logger:     logger,
		httpClient: &http.Client{Timeout: 20 * time.Second},
	}
}

func (c *Client) Send(payload Payload) error {
	url := c.baseURL + "/api/v1/integrations/whatsapp/webhook"

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error serializando el webhook: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("error armando el request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Hub-Signature-256", "sha256="+c.sign(body))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error enviando el webhook: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("el webhook respondio %d", resp.StatusCode)
	}

	return nil
}

func (c *Client) sign(payload []byte) string {
	mac := hmac.New(sha256.New, []byte(c.secret))
	mac.Write(payload)
	return hex.EncodeToString(mac.Sum(nil))
}

func Now() string {
	return fmt.Sprintf("%d", time.Now().Unix())
}
