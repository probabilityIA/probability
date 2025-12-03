package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// Configurar logger
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	log.Info().Msg("Iniciando servicio notify-email...")

	// Crear contexto con cancelación
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Configurar captura de señales para shutdown graceful
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Esperar señal de terminación
	<-sigChan
	log.Info().Ctx(ctx).Msg("Recibida señal de terminación, cerrando servicio...")
	cancel()

	log.Info().Msg("Servicio notify-email finalizado")
}
