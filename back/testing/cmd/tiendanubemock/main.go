package main

import (
	"os"

	"github.com/secamc93/probability/back/testing/integrations/tiendanube"
	"github.com/secamc93/probability/back/testing/shared/log"
)

func main() {
	logger := log.New()
	port := os.Getenv("TIENDANUBE_MOCK_PORT")
	if port == "" {
		port = "9102"
	}
	if err := tiendanube.New(logger, port).Start(); err != nil {
		logger.Error().Msgf("tiendanube mock: %s", err.Error())
		os.Exit(1)
	}
}
