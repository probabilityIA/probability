package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strings"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

type mediaClient struct {
	baseURL string
	http    *http.Client
	logger  log.ILogger
}

func NewMediaClient(baseURL string, logger log.ILogger) ports.IMediaAPI {
	return &mediaClient{
		baseURL: strings.TrimRight(baseURL, "/"),
		http:    &http.Client{Timeout: 60 * time.Second},
		logger:  logger.WithModule("whatsapp-media-client"),
	}
}

func (c *mediaClient) UploadMedia(ctx context.Context, phoneNumberID uint, accessToken, filename, mimeType string, data []byte) (string, error) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	if err := writer.WriteField("messaging_product", "whatsapp"); err != nil {
		return "", err
	}
	if err := writer.WriteField("type", mimeType); err != nil {
		return "", err
	}

	header := make(textproto.MIMEHeader)
	header.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file"; filename="%s"`, strings.ReplaceAll(filename, `"`, "")))
	header.Set("Content-Type", mimeType)
	part, err := writer.CreatePart(header)
	if err != nil {
		return "", err
	}
	if _, err := part.Write(data); err != nil {
		return "", err
	}
	if err := writer.Close(); err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/%d/media", c.baseURL, phoneNumberID), &body)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	var result struct {
		ID string `json:"id"`
	}
	if err := c.doJSON(req, phoneNumberID, &result); err != nil {
		return "", err
	}
	if result.ID == "" {
		return "", fmt.Errorf("meta no devolvio el id del archivo subido")
	}
	return result.ID, nil
}

func (c *mediaClient) SendMediaMessage(ctx context.Context, phoneNumberID uint, accessToken, to, mediaType, mediaID, caption, filename string) (string, error) {
	media := map[string]string{"id": mediaID}
	if caption != "" && mediaType != entities.MediaTypeAudio && mediaType != entities.MediaTypeSticker {
		media["caption"] = caption
	}
	if mediaType == entities.MediaTypeDocument && filename != "" {
		media["filename"] = filename
	}

	payload := map[string]any{
		"messaging_product": "whatsapp",
		"recipient_type":    "individual",
		"to":                to,
		"type":              mediaType,
		mediaType:           media,
	}
	raw, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/%d/messages", c.baseURL, phoneNumberID), bytes.NewReader(raw))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	var result struct {
		Messages []struct {
			ID string `json:"id"`
		} `json:"messages"`
	}
	if err := c.doJSON(req, phoneNumberID, &result); err != nil {
		return "", err
	}
	if len(result.Messages) == 0 || result.Messages[0].ID == "" {
		return "", fmt.Errorf("meta no devolvio el id del mensaje")
	}
	return result.Messages[0].ID, nil
}

func (c *mediaClient) GetMediaInfo(ctx context.Context, mediaID, accessToken string) (*ports.MediaInfo, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/%s", c.baseURL, mediaID), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	var result struct {
		URL      string `json:"url"`
		MimeType string `json:"mime_type"`
		FileSize int64  `json:"file_size"`
	}
	if err := c.doJSON(req, 0, &result); err != nil {
		return nil, err
	}
	if result.URL == "" {
		return nil, fmt.Errorf("meta no devolvio la url del archivo %s", mediaID)
	}
	return &ports.MediaInfo{URL: result.URL, MimeType: result.MimeType, FileSize: result.FileSize}, nil
}

func (c *mediaClient) DownloadMedia(ctx context.Context, url, accessToken string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error descargando archivo de meta: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("meta respondio %d al descargar el archivo", resp.StatusCode)
	}

	data, err := io.ReadAll(io.LimitReader(resp.Body, 100<<20))
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (c *mediaClient) doJSON(req *http.Request, phoneNumberID uint, target any) error {
	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("error al comunicarse con WhatsApp API: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		c.logger.Error().Int("status_code", resp.StatusCode).Str("response_body", string(body)).Msg("WhatsApp media API returned error status")
		return parseMetaGraphError(string(body), resp.StatusCode, phoneNumberID)
	}
	return json.Unmarshal(body, target)
}
