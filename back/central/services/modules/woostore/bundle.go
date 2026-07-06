package woostore

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/woostore/internal/app"
	"github.com/secamc93/probability/back/central/services/modules/woostore/internal/infra/primary/handlers"
	"github.com/secamc93/probability/back/central/services/modules/woostore/internal/infra/secondary/awsec2"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
)

func New(router *gin.RouterGroup, environment env.IConfig, logger log.ILogger) {
	logger = logger.WithModule("woostore")

	region := environment.Get("WOO_STORE_AWS_REGION")
	key := environment.Get("WOO_STORE_AWS_KEY")
	secret := environment.Get("WOO_STORE_AWS_SECRET")
	instanceID := environment.Get("WOO_STORE_INSTANCE_ID")
	storeURL := environment.Get("WOO_STORE_URL")

	if instanceID == "" || key == "" || secret == "" {
		logger.Warn(context.Background()).Msg("woo-store power deshabilitado (falta WOO_STORE_INSTANCE_ID/KEY/SECRET)")
		return
	}

	client, err := awsec2.New(region, key, secret, instanceID, storeURL)
	if err != nil {
		logger.Error(context.Background()).Err(err).Msg("no se pudo inicializar el cliente EC2 de woo-store")
		return
	}

	uc := app.New(client)
	h := handlers.New(uc, logger)
	h.RegisterRoutes(router)
}
