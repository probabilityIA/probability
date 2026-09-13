package handlers

import (
	"time"

	"github.com/secamc93/probability/back/testing/integrations/whatsappmock/internal/store"
	"github.com/secamc93/probability/back/testing/integrations/whatsappmock/internal/webhook"
	"github.com/secamc93/probability/back/testing/shared/log"
)

type Handler struct {
	store        *store.Store
	webhook      *webhook.Client
	logger       log.ILogger
	replyDelay   time.Duration
	approveDelay time.Duration
	displayPhone string
}

func New(store *store.Store, webhookClient *webhook.Client, logger log.ILogger) *Handler {
	return &Handler{
		store:        store,
		webhook:      webhookClient,
		logger:       logger,
		replyDelay:   700 * time.Millisecond,
		approveDelay: 3 * time.Second,
		displayPhone: "+573001234599",
	}
}
