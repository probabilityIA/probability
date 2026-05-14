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
