package usecases

import (
	"github.com/secamc93/probability/back/testing/integrations/siigo/internal/domain"
	"github.com/secamc93/probability/back/testing/shared/log"
)

type APISimulator struct {
	logger     log.ILogger
	Repository *domain.Repository
}

func NewAPISimulator(logger log.ILogger) *APISimulator {
	return &APISimulator{
		logger:     logger,
		Repository: domain.NewRepository(),
	}
}

func (s *APISimulator) HandleListWebhooks() []*domain.Webhook {
	return s.Repository.ListWebhooks()
}

func (s *APISimulator) HandleCreateWebhook(applicationID, url, topic string) *domain.Webhook {
	return s.Repository.CreateWebhook(applicationID, url, topic)
}

func (s *APISimulator) HandleDeleteWebhook(id string) bool {
	return s.Repository.DeleteWebhook(id)
}
