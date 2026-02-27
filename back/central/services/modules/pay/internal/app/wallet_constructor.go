package app

import (
	"github.com/secamc93/probability/back/central/services/modules/pay/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

// walletUseCase implementa IWalletUseCase
type walletUseCase struct {
	repo           ports.IRepository
	paymentUseCase ports.IUseCase
	log            log.ILogger
}

// NewWalletUseCase crea una nueva instancia del use case de wallet
func NewWalletUseCase(
	repo ports.IRepository,
	paymentUseCase ports.IUseCase,
	logger log.ILogger,
) ports.IWalletUseCase {
	return &walletUseCase{
		repo:           repo,
		paymentUseCase: paymentUseCase,
		log:            logger.WithModule("pay.wallet.usecase"),
	}
}
