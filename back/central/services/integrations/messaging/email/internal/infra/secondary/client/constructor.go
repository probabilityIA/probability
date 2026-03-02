package client

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/messaging/email/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/email"
)

// adapter envuelve el servicio SES compartido satisfaciendo ports.IEmailClient
type adapter struct {
	ses email.IEmailService
}

// New crea un adaptador que conecta shared/email.IEmailService con el port local IEmailClient
func New(ses email.IEmailService) ports.IEmailClient {
	return &adapter{ses: ses}
}

func (a *adapter) SendHTML(ctx context.Context, to, subject, html string) error {
	return a.ses.SendHTML(ctx, to, subject, html)
}
