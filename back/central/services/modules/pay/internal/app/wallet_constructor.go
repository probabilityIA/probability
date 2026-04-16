package app

import (
	"github.com/secamc93/probability/back/central/services/modules/pay/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
)

// walletUseCase implementa IWalletUseCase
type walletUseCase struct {
	repo           ports.IRepository
	paymentUseCase ports.IUseCase
	config         env.IConfig
	log            log.ILogger
}

// NewWalletUseCase crea una nueva instancia del use case de wallet
func NewWalletUseCase(
	repo ports.IRepository,
	paymentUseCase ports.IUseCase,
	config env.IConfig,
	logger log.ILogger,
) ports.IWalletUseCase {
	return &walletUseCase{
		repo:           repo,
		paymentUseCase: paymentUseCase,
		config:         config,
		log:            logger.WithModule("pay.wallet.usecase"),
	}
}
