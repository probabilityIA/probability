package bedrock

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
)

// IBedrock define la interfaz para la conexion a Amazon Bedrock
type IBedrock interface {
	GetClient() *bedrockruntime.Client
	Converse(ctx context.Context, input *bedrockruntime.ConverseInput) (*bedrockruntime.ConverseOutput, error)
}

type bedrock struct {
	client *bedrockruntime.Client
	log    log.ILogger
	env    env.IConfig
}

// New crea una nueva instancia de Bedrock y retorna la interfaz
func New(logger log.ILogger, cfg env.IConfig) IBedrock {
	b := &bedrock{
		log: logger,
		env: cfg,
	}

	if err := b.connect(context.Background()); err != nil {
		logger.Warn(context.Background()).
			Err(err).
			Msg("No se pudo conectar a Bedrock - el servicio de AI no estara disponible")
		return b
	}

	return b
}

func (b *bedrock) connect(ctx context.Context) error {
	region := b.env.Get("BEDROCK_REGION")
	if region == "" {
		region = "us-east-1"
	}

	accessKey := b.env.Get("BEDROCK_ACCESS_KEY")
	secretKey := b.env.Get("BEDROCK_SECRET_KEY")

	opts := []func(*config.LoadOptions) error{
		config.WithRegion(region),
	}

	if accessKey != "" && secretKey != "" {
		opts = append(opts, config.WithCredentialsProvider(
			aws.NewCredentialsCache(
				credentials.NewStaticCredentialsProvider(accessKey, secretKey, ""),
			),
		))
	}

	cfg, err := config.LoadDefaultConfig(ctx, opts...)
	if err != nil {
		b.log.Error(ctx).
			Err(err).
			Str("region", region).
			Msg("Error al cargar configuracion de AWS para Bedrock")
		return fmt.Errorf("error al cargar configuracion de AWS: %w", err)
	}

	b.client = bedrockruntime.NewFromConfig(cfg)

	b.log.Info(ctx).
		Str("region", region).
		Msg("Conexion a Amazon Bedrock establecida correctamente")

	return nil
}

// GetClient retorna el cliente de Bedrock
func (b *bedrock) GetClient() *bedrockruntime.Client {
	return b.client
}

// Converse ejecuta una llamada al API Converse de Bedrock
func (b *bedrock) Converse(ctx context.Context, input *bedrockruntime.ConverseInput) (*bedrockruntime.ConverseOutput, error) {
	if b.client == nil {
		return nil, fmt.Errorf("cliente de Bedrock no inicializado")
	}
	return b.client.Converse(ctx, input)
}
