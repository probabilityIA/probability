package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/wallet/app/usecases"
)

type IWalletHandlers interface {
	GetBalance(c *gin.Context)
	RechargeWallet(c *gin.Context)
	GetAllWallets(c *gin.Context)
	RegisterRoutes(router *gin.RouterGroup)
}

type WalletHandlers struct {
	uc *usecases.WalletUsecases
}

func New(uc *usecases.WalletUsecases) IWalletHandlers {
	return &WalletHandlers{
		uc: uc,
	}
}
