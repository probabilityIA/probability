package email

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	sestypes "github.com/aws/aws-sdk-go-v2/service/sesv2/types"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
)

// IEmailService interfaz genérica de envío de correo
type IEmailService interface {
	SendHTML(ctx context.Context, to, subject, html string) error
}

type EmailService struct {
	client    *sesv2.Client
	fromEmail string
	logger    log.ILogger
}

// New crea una nueva instancia del servicio de email usando Amazon SES.
// Lee las variables SES_REGION, SES_ACCESS_KEY, SES_SECRET_KEY y FROM_EMAIL.
func New(cfg env.IConfig, logger log.ILogger) IEmailService {
	region := cfg.Get("SES_REGION")
	accessKey := cfg.Get("SES_ACCESS_KEY")
	secretKey := cfg.Get("SES_SECRET_KEY")
	fromEmail := cfg.Get("FROM_EMAIL")

	if region == "" || accessKey == "" || secretKey == "" || fromEmail == "" {
		logger.Fatal(context.Background()).
			Bool("has_region", region != "").
			Bool("has_access_key", accessKey != "").
			Bool("has_secret_key", secretKey != "").
			Bool("has_from_email", fromEmail != "").
			Msg("❌ Configuración de Amazon SES incompleta — verifica SES_REGION, SES_ACCESS_KEY, SES_SECRET_KEY y FROM_EMAIL")
		panic("configuración de Amazon SES incompleta")
	}

	awsCfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(region),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(accessKey, secretKey, ""),
		),
	)
	if err != nil {
		logger.Fatal(context.Background()).
			Err(err).
			Msg("❌ Error cargando configuración AWS para SES")
		panic("error cargando configuración AWS para SES: " + err.Error())
	}

	logger.Info(context.Background()).
		Str("region", region).
		Str("from_email", fromEmail).
		Msg("✅ Amazon SES inicializado correctamente")

	return &EmailService{
		client:    sesv2.NewFromConfig(awsCfg),
		fromEmail: fromEmail,
		logger:    logger,
	}
}

// SendHTML envía un correo electrónico con contenido HTML a través de Amazon SES.
func (e *EmailService) SendHTML(ctx context.Context, to, subject, html string) error {
	input := &sesv2.SendEmailInput{
		FromEmailAddress: aws.String(e.fromEmail),
		Destination: &sestypes.Destination{
			ToAddresses: []string{to},
		},
		Content: &sestypes.EmailContent{
			Simple: &sestypes.Message{
				Subject: &sestypes.Content{
					Data:    aws.String(subject),
					Charset: aws.String("UTF-8"),
				},
				Body: &sestypes.Body{
					Html: &sestypes.Content{
						Data:    aws.String(html),
						Charset: aws.String("UTF-8"),
					},
				},
			},
		},
	}

	_, err := e.client.SendEmail(ctx, input)
	if err != nil {
		e.logger.Error(ctx).
			Err(err).
			Str("to", to).
			Str("subject", subject).
			Str("from", e.fromEmail).
			Msg("Error enviando email via Amazon SES")
		return fmt.Errorf("error enviando email via SES: %w", err)
	}

	e.logger.Info(ctx).
		Str("to", to).
		Str("subject", subject).
		Msg("Email enviado exitosamente via Amazon SES")

	return nil
}
