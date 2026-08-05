package app

import (
	"github.com/secamc93/probability/back/central/services/modules/pay/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

type walletUseCase struct {
	repo           ports.IRepository
	paymentUseCase ports.IUseCase
	rabbit         rabbitmq.IQueue
	config         env.IConfig
	log            log.ILogger
}

func NewWalletUseCase(
	repo ports.IRepository,
	paymentUseCase ports.IUseCase,
	rabbit rabbitmq.IQueue,
	config env.IConfig,
	logger log.ILogger,
) ports.IWalletUseCase {
	return &walletUseCase{
		repo:           repo,
		paymentUseCase: paymentUseCase,
		rabbit:         rabbit,
		config:         config,
		log:            logger.WithModule("pay.wallet.usecase"),
	}
}
