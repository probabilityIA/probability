package main

import (
	"os"

	whatsappmock "github.com/secamc93/probability/back/testing/integrations/whatsappmock"
	"github.com/secamc93/probability/back/testing/shared/log"
)

func main() {
	logger := log.New()

	port := os.Getenv("WHATSAPP_MOCK_PORT")
	if port == "" {
		port = "9103"
	}

	centralURL := os.Getenv("WEBHOOK_BASE_URL")
	if centralURL == "" {
		centralURL = "http://localhost:3050"
	}

	secret := os.Getenv("WHATSAPP_WEBHOOK_SECRET")
	if secret == "" {
		logger.Error().Msg("WHATSAPP_WEBHOOK_SECRET es obligatorio: el webhook de central rechaza la firma sin el")
		os.Exit(1)
	}

	if err := whatsappmock.New(logger, port, centralURL, secret).Start(); err != nil {
		logger.Error().Msgf("whatsapp mock: %s", err.Error())
		os.Exit(1)
	}
}
