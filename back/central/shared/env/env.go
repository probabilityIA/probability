package env

import (
	"context"
	"os"
	"path/filepath"
	"reflect"

	"github.com/joho/godotenv"
	"github.com/secamc93/probability/back/central/shared/log"
)

type IConfig interface {
	Get(key string) string
}

type config struct {
	values map[string]string
	logger log.ILogger
}

func loadDotEnv(logger log.ILogger) {
	// Intentar cargar .env desde el directorio actual
	_ = godotenv.Load(".env")

	// Si aún faltan claves, buscar hacia arriba hasta 6 niveles
	cwd, _ := os.Getwd()
	maxLevels := 6
	for i := 0; i < maxLevels; i++ {
		candidate := filepath.Join(cwd, ".env")
		if _, err := os.Stat(candidate); err == nil {
			_ = godotenv.Overload(candidate)
			// quitado: log de archivo .env cargado para evitar ruido en consola
			return
		}
		cwd = filepath.Dir(cwd)
	}
	// Si no se encontró, intentar rutas relativas comunes
	_ = godotenv.Overload("../.env", "../../.env", "../../../.env", "../../../../.env")
}

func New(logger log.ILogger) IConfig {
	loadDotEnv(logger)

	cfg := &Config{}
	missing := []string{}
	values := make(map[string]string)

	v := reflect.ValueOf(cfg).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("env")
		if tag == "" {
			continue
		}
		parts := splitTag(tag)
		key := parts[0]
		required := len(parts) > 1 && parts[1] == "required"
		val := os.Getenv(key)
		if val == "" && required {
			missing = append(missing, key)
		}
		values[key] = val
	}

	if len(missing) > 0 {
		if os.Getenv("RELAX_ENV") == "1" {
			logger.Warn(context.Background()).
				Strs("missing_env_vars", missing).
				Msg("Faltan variables de entorno obligatorias - modo relajado activo (RELAX_ENV=1)")
		} else {
			logger.Fatal(context.Background()).
				Strs("missing_env_vars", missing).
				Msg("Faltan variables de entorno obligatorias - la aplicación no puede continuar")
			// El panic se ejecutará automáticamente después del log fatal
		}
	}

	return &config{values: values, logger: logger}
}

// NewWithLogging crea una nueva configuración con logging automático de errores
func NewWithLogging(logger log.ILogger) IConfig {
	loadDotEnv(logger)

	cfg := &Config{}
	missing := []string{}
	values := make(map[string]string)

	v := reflect.ValueOf(cfg).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("env")
		if tag == "" {
			continue
		}
		parts := splitTag(tag)
		key := parts[0]
		required := len(parts) > 1 && parts[1] == "required"
		val := os.Getenv(key)
		if val == "" && required {
			missing = append(missing, key)
		}
		values[key] = val
	}

	if len(missing) > 0 {
		if os.Getenv("RELAX_ENV") == "1" {
			logger.Warn(context.Background()).
				Strs("missing_env_vars", missing).
				Msg("Faltan variables de entorno obligatorias - modo relajado activo (RELAX_ENV=1)")
		} else {
			logger.Fatal(context.Background()).
				Strs("missing_env_vars", missing).
				Msg("Faltan variables de entorno obligatorias - la aplicación no puede continuar")
			// El panic se ejecutará automáticamente después del log fatal
		}
	}

	return &config{values: values, logger: logger}
}

// Get retorna el valor de una variable de entorno cargada
func (c *config) Get(key string) string {
	return c.values[key]
}

// Config solo se usa internamente para reflexión
// No debe ser accedido directamente fuera de este paquete

type Config struct {
	AppEnv    string `env:"APP_ENV,required"`
	HttpPort  string `env:"HTTP_PORT,required"`
	GrpcPort  string `env:"GRPC_PORT"`
	LogLevel  string `env:"LOG_LEVEL,required"`
	JwtSecret string `env:"JWT_SECRET,required"`
	// NatsHost   string `env:"NATS_HOST,required"`
	// NatsPort   string `env:"NATS_PORT,required"`
	// NatsUser   string `env:"NATS_USER,required"`
	// NatsPass   string `env:"NATS_PASS,required"`
	DbHost         string `env:"DB_HOST,required"`
	DbUser         string `env:"DB_USER,required"`
	DbPass         string `env:"DB_PASS,required"`
	DbPort         string `env:"DB_PORT,required"`
	DbName         string `env:"DB_NAME,required"`
	DbLogLevel     string `env:"DB_LOG_LEVEL,required"`
	PGSSLMODE      string `env:"PGSSLMODE,required"`
	URLBaseSwagger string `env:"URL_BASE_SWAGGER,required"`
	S3Bucket       string `env:"S3_BUCKET,required"`
	S3Region       string `env:"S3_REGION,required"`
	S3AccessKey    string `env:"S3_KEY,required"`
	S3SecretKey    string `env:"S3_SECRET,required"`
	S3Endpoint     string `env:"S3_ENDPOINT"`

	// SMTP/Email
	SMTPHost           string `env:"SMTP_HOST"`
	SMTPPort           string `env:"SMTP_PORT"`
	SMTPUser           string `env:"SMTP_USER"`
	SMTPPass           string `env:"SMTP_PASS"`
	FromEmail          string `env:"FROM_EMAIL"`
	SMTPUseSTARTTLS    string `env:"SMTP_USE_STARTTLS"`
	SMTPUseTLS         string `env:"SMTP_USE_TLS"`
	UrlBaseDomainS3    string `env:"URL_BASE_DOMAIN_S3"`
	WhatsAppURL        string `env:"WHATSAPP_URL,required"`
	WhatsAppToken      string `env:"WHATSAPP_TOKEN,required"`
	WhatsAppPhoneNumID string `env:"WHATSAPP_PHONE_NUMBER_ID,required"`

	// DynamoDB
	DynamoRegion    string `env:"DYNAMO_REGION"`
	DynamoAccessKey string `env:"DYNAMO_ACCESS_KEY"`
	DynamoSecretKey string `env:"DYNAMO_SECRET_KEY"`

	EncryptionKey string `env:"ENCRYPTION_KEY,required"`

	RabbitMQHost  string `env:"RABBITMQ_HOST,required"`
	RabbitMQPort  string `env:"RABBITMQ_PORT,required"`
	RabbitMQUser  string `env:"RABBITMQ_USER,required"`
	RabbitMQPass  string `env:"RABBITMQ_PASS,required"`
	RabbitMQVHost string `env:"RABBITMQ_VHOST,required"`

	RedisOrderEventsChannel string `env:"REDIS_ORDER_EVENTS_CHANNEL,required"`
}

func splitTag(tag string) []string {
	// Usamos SplitN para dividir solo en la primera coma
	parts := make([]string, 0, 2)
	for i, c := range tag {
		if c == ',' {
			parts = append(parts, tag[:i], tag[i+1:])
			return parts
		}
	}
	return []string{tag}
}
