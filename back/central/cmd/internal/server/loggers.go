package server

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
)

// LogStartupInfo muestra información de inicio del servidor y conexiones
func LogStartupInfo(ctx context.Context, logger log.ILogger, e env.IConfig) {
	// No mostrar nombre de función en logs de startup
	ctx = log.WithSkipFunctionCtx(ctx)

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
	logger.Info(ctx).Msg(" ")

	// Cabecera
	logger.Info(ctx).Msg(" 🚀 Servidor HTTP iniciado correctamente")
	logger.Info(ctx).Msgf(" 📍 Disponible en: %s", coloredURL)
	logger.Info(ctx).Msgf(" 📖 Documentación: %s", coloredDocs)
	logger.Info(ctx).Msg(" ")

	// PostgreSQL (si aplica)
	dbHost := e.Get("DB_HOST")
	dbPort := e.Get("DB_PORT")
	dbName := e.Get("DB_NAME")
	if dbHost != "" && dbPort != "" && dbName != "" {
		dbURL := fmt.Sprintf("postgres://%s:%s/%s", dbHost, dbPort, dbName)
		coloredDB := fmt.Sprintf("\033[36;4m%s\033[0m", dbURL) // cian subrayado
		logger.Info(ctx).Msgf(" 🗄️  Conexión PostgreSQL: %s", coloredDB)
		logger.Info(ctx).Msg(" ")
	}

	// RabbitMQ
	rabbitHost := e.Get("RABBITMQ_HOST")
	rabbitPort := e.Get("RABBITMQ_PORT")
	rabbitVHost := e.Get("RABBITMQ_VHOST")
	if rabbitHost != "" && rabbitPort != "" {
		rabbitURL := fmt.Sprintf("amqp://%s:%s%s", rabbitHost, rabbitPort, rabbitVHost)
		coloredRabbit := fmt.Sprintf("\033[32;4m%s\033[0m", rabbitURL) // verde subrayado
		logger.Info(ctx).Msgf(" 🐰 RabbitMQ: %s", coloredRabbit)
		logger.Info(ctx).Msg(" ")
	}

	// Redis
	redisHost := e.Get("REDIS_HOST")
	redisPort := e.Get("REDIS_PORT")
	if redisHost == "" {
		redisHost = "localhost"
	}
	if redisPort == "" {
		redisPort = "6379"
	}
	if redisHost != "" && redisPort != "" {
		redisURL := fmt.Sprintf("redis://%s:%s", redisHost, redisPort)
		coloredRedis := fmt.Sprintf("\033[31;4m%s\033[0m", redisURL) // rojo subrayado
		logger.Info(ctx).Msgf(" 🔴 Redis: %s", coloredRedis)
		logger.Info(ctx).Msg(" ")
	}

	// S3 (si aplica)
	s3Region := e.Get("S3_REGION")
	s3Bucket := e.Get("S3_BUCKET")
	s3Endpoint := e.Get("S3_ENDPOINT")
	if s3Bucket != "" {
		var s3URL string
		if s3Endpoint != "" {
			// MinIO o S3-compatible
			s3URL = fmt.Sprintf("%s/%s (%s)", s3Endpoint, s3Bucket, s3Region)
		} else {
			// AWS S3 estándar
			s3URL = fmt.Sprintf("s3://%s (%s)", s3Bucket, s3Region)
		}
		coloredS3 := fmt.Sprintf("\033[35;4m%s\033[0m", s3URL) // magenta subrayado
		logger.Info(ctx).Msgf(" ☁️  S3 Storage: %s", coloredS3)
		logger.Info(ctx).Msg(" ")
	}

	// Espacio final
	logger.Info(ctx).Msg(" ")
}
