package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/pay/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

// IWalletHandler define la interfaz del handler de wallet
type IWalletHandler interface {
	RegisterWalletRoutes(router *gin.RouterGroup)
}

// walletHandler implementa IWalletHandler
type walletHandler struct {
	walletUC ports.IWalletUseCase
	log      log.ILogger
}

// NewWalletHandler crea una nueva instancia del handler de wallet
func NewWalletHandler(walletUC ports.IWalletUseCase, logger log.ILogger) IWalletHandler {
	return &walletHandler{
		walletUC: walletUC,
		log:      logger.WithModule("pay.wallet.handler"),
	}
}
