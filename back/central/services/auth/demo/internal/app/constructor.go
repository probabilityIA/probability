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
	DemoVerifyOTP(ctx context.Context, request domain.DemoVerifyOTPRequest) (*domain.DemoVerifyOTPResponse, error)
}

type UseCase struct {
	repository   domain.IDemoRepository
	emailSender  domain.IEmailSender
	otpPublisher domain.IDemoOTPPublisher
	log          log.ILogger
	env          env.IConfig
}

func New(repository domain.IDemoRepository, emailSender domain.IEmailSender, otpPublisher domain.IDemoOTPPublisher, log log.ILogger, env env.IConfig) IUseCase {
	return &UseCase{
		repository:   repository,
		emailSender:  emailSender,
		otpPublisher: otpPublisher,
		log:          log,
		env:          env,
	}
}
