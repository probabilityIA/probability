package server

import (
	"central/shared/env"
	"central/shared/log"
	"context"
	"fmt"
)

// LogStartupInfo muestra información de inicio del servidor y conexiones
func LogStartupInfo(ctx context.Context, logger log.ILogger, e env.IConfig) {
	port := e.Get("HTTP_PORT")
	serverURL := fmt.Sprintf("http://localhost:%s", port)

	swaggerBaseURL := e.Get("URL_BASE_SWAGGER")
	if swaggerBaseURL == "" {
		swaggerBaseURL = serverURL
	}
	docsURL := fmt.Sprintf("%s/docs/index.html", swaggerBaseURL)

	// Colores ANSI para URLs
	coloredURL := fmt.Sprintf("\033[34;4m%s\033[0m", serverURL) // azul subrayado
	coloredDocs := fmt.Sprintf("\033[33;4m%s\033[0m", docsURL)  // amarillo subrayado

	// Espacio inicial
	logger.Info().Msg(" ")

	// Cabecera
	logger.Info().Msg(" 🚀 Servidor HTTP iniciado correctamente")
	logger.Info().Msgf(" 📍 Disponible en: %s", coloredURL)
	logger.Info().Msgf(" 📖 Documentación: %s", coloredDocs)
	logger.Info().Msg(" ")

	// PostgreSQL (si aplica)
	dbHost := e.Get("DB_HOST")
	dbPort := e.Get("DB_PORT")
	dbName := e.Get("DB_NAME")
	if dbHost != "" && dbPort != "" && dbName != "" {
		dbURL := fmt.Sprintf("postgres://%s:%s/%s", dbHost, dbPort, dbName)
		coloredDB := fmt.Sprintf("\033[36;4m%s\033[0m", dbURL) // cian subrayado
		logger.Info().Msgf(" 🗄️  Conexión PostgreSQL: %s", coloredDB)
		logger.Info().Msg(" ")
	}

	// S3 (si aplica)
	s3Region := e.Get("S3_REGION")
	s3Bucket := e.Get("S3_BUCKET")
	if s3Bucket != "" {
		s3URL := fmt.Sprintf("s3://%s (%s)", s3Bucket, s3Region)
		coloredS3 := fmt.Sprintf("\033[35;4m%s\033[0m", s3URL) // magenta subrayado
		logger.Info().Msgf(" ☁️  AWS S3: %s", coloredS3)
		logger.Info().Msg(" ")
	}

	// Espacio final
	logger.Info().Msg(" ")
}
