package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/auth/demo/internal/domain"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
)

type IUseCase interface {
	DemoRegister(ctx context.Context, request domain.DemoRegisterRequest) (*domain.DemoRegisterResponse, error)
	VerifyEmail(ctx context.Context, request domain.VerifyEmailRequest) (*domain.VerifyEmailResponse, error)
}

type UseCase struct {
	repository  domain.IDemoRepository
	emailSender domain.IEmailSender
	log         log.ILogger
	env         env.IConfig
}

func New(repository domain.IDemoRepository, emailSender domain.IEmailSender, log log.ILogger, env env.IConfig) IUseCase {
	return &UseCase{
		repository:  repository,
		emailSender: emailSender,
		log:         log,
		env:         env,
	}
}
