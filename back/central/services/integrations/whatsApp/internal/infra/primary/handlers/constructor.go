package handlers

import (
	"github.com/secamc93/probability/back/central/services/integrations/whatsApp/internal/app"
	"github.com/secamc93/probability/back/central/shared/log"
)

type WhatsAppHandler struct {
	useCase app.IUseCaseSendMessage
	logger  log.ILogger
}

func New(useCase app.IUseCaseSendMessage, logger log.ILogger) *WhatsAppHandler {
	// El logger ya viene con service="integrations" desde el bundle
	// Solo agregamos el módulo específico
	contextualLogger := logger.WithModule("whatsapp")

	return &WhatsAppHandler{
		useCase: useCase,
		logger:  contextualLogger,
	}
}
