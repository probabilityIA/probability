package fcm

import (
	"context"
	"net/http"
	"os"
	"strings"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/secamc93/probability/back/central/services/modules/push/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
)

const fcmScope = "https://www.googleapis.com/auth/firebase.messaging"

type client struct {
	projectID  string
	tokenSrc   oauth2.TokenSource
	httpClient *http.Client
	log        log.ILogger
}

func New(cfg env.IConfig, logger log.ILogger) ports.IPushSender {
	c := &client{
		projectID:  strings.TrimSpace(cfg.Get("FCM_PROJECT_ID")),
		httpClient: &http.Client{Timeout: 15 * time.Second},
		log:        logger.WithModule("push-fcm"),
	}

	creds := credentialsJSON(cfg)
	if len(creds) == 0 || c.projectID == "" {
		c.log.Warn(context.Background()).
			Msg("FCM sin credenciales o sin project id, las notificaciones push quedan desactivadas")
		return c
	}

	jwtConfig, err := google.JWTConfigFromJSON(creds, fcmScope)
	if err != nil {
		c.log.Error(context.Background()).Err(err).Msg("credenciales de FCM invalidas")
		return c
	}

	c.tokenSrc = jwtConfig.TokenSource(context.Background())
	c.log.Info(context.Background()).Str("project_id", c.projectID).Msg("FCM configurado")
	return c
}

func (c *client) IsConfigured() bool {
	return c.tokenSrc != nil && c.projectID != ""
}

func credentialsJSON(cfg env.IConfig) []byte {
	if raw := strings.TrimSpace(cfg.Get("FCM_CREDENTIALS_JSON")); raw != "" {
		return []byte(raw)
	}

	path := strings.TrimSpace(cfg.Get("FCM_CREDENTIALS_FILE"))
	if path == "" {
		return nil
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	return content
}
