package email

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
)

const resendEndpoint = "https://api.resend.com/emails"

type IEmailService interface {
	SendHTML(ctx context.Context, to, subject, html string) error
}

type EmailService struct {
	apiKey    string
	fromEmail string
	client    *http.Client
	logger    log.ILogger
}

func New(cfg env.IConfig, logger log.ILogger) IEmailService {
	apiKey := cfg.Get("RESEND_API_KEY")
	fromEmail := cfg.Get("FROM_EMAIL")

	if apiKey == "" || fromEmail == "" {
		logger.Fatal(context.Background()).
			Bool("has_api_key", apiKey != "").
			Bool("has_from_email", fromEmail != "").
			Msg("Configuracion de Resend incompleta - verifica RESEND_API_KEY y FROM_EMAIL")
		panic("configuracion de Resend incompleta")
	}

	logger.Info(context.Background()).
		Str("from_email", fromEmail).
		Msg("Resend inicializado correctamente")

	return &EmailService{
		apiKey:    apiKey,
		fromEmail: fromEmail,
		client:    &http.Client{Timeout: 15 * time.Second},
		logger:    logger,
	}
}

type resendRequest struct {
	From    string   `json:"from"`
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	Html    string   `json:"html"`
}

type resendResponse struct {
	ID string `json:"id"`
}

func (e *EmailService) SendHTML(ctx context.Context, to, subject, html string) error {
	payload := resendRequest{
		From:    e.fromEmail,
		To:      []string{to},
		Subject: subject,
		Html:    html,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error serializando payload de Resend: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, resendEndpoint, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("error creando request de Resend: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+e.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := e.client.Do(req)
	if err != nil {
		e.logger.Error(ctx).
			Err(err).
			Str("to", to).
			Str("subject", subject).
			Str("from", e.fromEmail).
			Msg("Error enviando email via Resend")
		return fmt.Errorf("error enviando email via Resend: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		e.logger.Error(ctx).
			Int("status_code", resp.StatusCode).
			Str("to", to).
			Str("subject", subject).
			Str("from", e.fromEmail).
			Str("response", string(respBody)).
			Msg("Resend respondio con error")
		return fmt.Errorf("resend respondio status %d: %s", resp.StatusCode, string(respBody))
	}

	var parsed resendResponse
	_ = json.Unmarshal(respBody, &parsed)

	e.logger.Info(ctx).
		Str("to", to).
		Str("subject", subject).
		Str("message_id", parsed.ID).
		Msg("Email enviado exitosamente via Resend")

	return nil
}
