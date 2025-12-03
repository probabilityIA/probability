package dynamo

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
)

// IDynamoDB define la interfaz para la conexión a DynamoDB
type IDynamoDB interface {
	Connect(ctx context.Context) error
	Close() error
	GetClient() *dynamodb.Client
	WithContext(ctx context.Context) *dynamodb.Client
	HealthCheck(ctx context.Context) error
}

// dynamo implementa la interfaz IDynamoDB
type dynamo struct {
	client *dynamodb.Client
	config aws.Config
	log    log.ILogger
	env    env.IConfig
}

// New crea una nueva instancia de DynamoDB y retorna la interfaz
func New(logger log.ILogger, config env.IConfig) IDynamoDB {
	d := &dynamo{
		log: logger,
		env: config,
	}

	// Conectar automáticamente a DynamoDB
	if err := d.Connect(context.Background()); err != nil {
		logger.Fatal(context.Background()).
			Err(err).
			Msg("Error al conectar a DynamoDB - la aplicación no puede continuar")
	}

	return d
}

// Connect establece la conexión con DynamoDB
func (d *dynamo) Connect(ctx context.Context) error {
	region := d.env.Get("DYNAMO_REGION")
	if region == "" {
		region = "us-east-1" // Valor por defecto
	}

	accessKey := d.env.Get("DYNAMO_ACCESS_KEY")
	secretKey := d.env.Get("DYNAMO_SECRET_KEY")

	var cfg aws.Config
	var err error

	// Si se proporcionan credenciales, usarlas; de lo contrario, usar credenciales por defecto de AWS
	if accessKey != "" && secretKey != "" {
		// Usar credenciales estáticas
		cfg, err = config.LoadDefaultConfig(ctx,
			config.WithRegion(region),
			config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
		)
	} else {
		// Usar credenciales por defecto de AWS (IAM role, environment variables, etc.)
		cfg, err = config.LoadDefaultConfig(ctx,
			config.WithRegion(region),
		)
	}

	if err != nil {
		d.log.Error(ctx).
			Err(err).
			Str("region", region).
			Msg("Error al cargar configuración de AWS para DynamoDB")
		return fmt.Errorf("error al cargar configuración de AWS: %w", err)
	}

	// Crear cliente de DynamoDB
	d.client = dynamodb.NewFromConfig(cfg)
	d.config = cfg

	d.log.Info(ctx).
		Str("region", region).
		Msg("Conexión a DynamoDB establecida correctamente")

	return nil
}

// Close cierra la conexión con DynamoDB
// Nota: DynamoDB client no requiere cierre explícito, pero implementamos el método para mantener la interfaz consistente
func (d *dynamo) Close() error {
	// El cliente de DynamoDB no requiere cierre explícito
	// pero podemos limpiar la referencia
	d.client = nil
	d.log.Info(context.Background()).Msg("Conexión a DynamoDB cerrada")
	return nil
}

// GetClient retorna el cliente de DynamoDB
func (d *dynamo) GetClient() *dynamodb.Client {
	return d.client
}

// WithContext retorna el cliente de DynamoDB (ya incluye contexto en las operaciones)
func (d *dynamo) WithContext(ctx context.Context) *dynamodb.Client {
	// El cliente de DynamoDB ya maneja el contexto en cada operación
	// Retornamos el mismo cliente ya que el contexto se pasa en cada llamada
	return d.client
}

// HealthCheck verifica que la conexión a DynamoDB esté funcionando
func (d *dynamo) HealthCheck(ctx context.Context) error {
	if d.client == nil {
		return fmt.Errorf("cliente de DynamoDB no inicializado")
	}

	// Intentar listar tablas como verificación de salud
	// Usamos un límite de 1 para minimizar el costo
	_, err := d.client.ListTables(ctx, &dynamodb.ListTablesInput{
		Limit: aws.Int32(1),
	})

	if err != nil {
		d.log.Error(ctx).
			Err(err).
			Msg("Error en health check de DynamoDB")
		return fmt.Errorf("error en health check de DynamoDB: %w", err)
	}

	d.log.Debug(ctx).Msg("Health check de DynamoDB exitoso")
	return nil
}
