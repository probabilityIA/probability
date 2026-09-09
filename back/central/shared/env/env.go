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

	_ = godotenv.Load(".env")

	cwd, _ := os.Getwd()
	maxLevels := 6
	for i := 0; i < maxLevels; i++ {
		candidate := filepath.Join(cwd, ".env")
		if _, err := os.Stat(candidate); err == nil {
			_ = godotenv.Overload(candidate)

			return
		}
		cwd = filepath.Dir(cwd)
	}

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

		}
	}

	return &config{values: values, logger: logger}
}

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

		}
	}

	return &config{values: values, logger: logger}
}

func (c *config) Get(key string) string {
	return c.values[key]
}

type Config struct {
	AppEnv    string `env:"APP_ENV,required"`
	HttpPort  string `env:"HTTP_PORT,required"`
	GrpcPort  string `env:"GRPC_PORT"`
	LogLevel  string `env:"LOG_LEVEL,required"`
	JwtSecret string `env:"JWT_SECRET,required"`

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

	RedisHost               string `env:"REDIS_HOST"`
	RedisPort               string `env:"REDIS_PORT"`
	RedisPassword           string `env:"REDIS_PASSWORD"`
	RedisOrderEventsChannel string `env:"REDIS_ORDER_EVENTS_CHANNEL,required"`

	ResendAPIKey       string `env:"RESEND_API_KEY"`
	FromEmail          string `env:"FROM_EMAIL"`
	FrontendBaseURL    string `env:"FRONTEND_BASE_URL"`
	UrlBaseDomainS3    string `env:"URL_BASE_DOMAIN_S3"`
	WhatsAppURL        string `env:"WHATSAPP_URL,required"`
	WhatsAppToken      string `env:"WHATSAPP_TOKEN,required"`
	WhatsAppPhoneNumID string `env:"WHATSAPP_PHONE_NUMBER_ID,required"`

	DynamoRegion    string `env:"DYNAMO_REGION"`
	DynamoAccessKey string `env:"DYNAMO_ACCESS_KEY"`
	DynamoSecretKey string `env:"DYNAMO_SECRET_KEY"`

	EncryptionKey string `env:"ENCRYPTION_KEY,required"`

	RabbitMQHost        string `env:"RABBITMQ_HOST,required"`
	RabbitMQPort        string `env:"RABBITMQ_PORT,required"`
	RabbitMQUser        string `env:"RABBITMQ_USER,required"`
	RabbitMQPass        string `env:"RABBITMQ_PASS,required"`
	RabbitMQVHost       string `env:"RABBITMQ_VHOST,required"`
	RabbitMQOrdersQueue string `env:"RABBITMQ_ORDERS_CREATE,required"`

	ShopifyClientID     string `env:"SHOPIFY_CLIENT_ID"`
	ShopifyClientSecret string `env:"SHOPIFY_CLIENT_SECRET"`
	ShopifyRedirectURI  string `env:"SHOPIFY_REDIRECT_URI"`
	ShopifyScopes       string `env:"SHOPIFY_SCOPES"`
	ShopifyShopDomain   string `env:"SHOPIFY_SHOP_DOMAIN"`
	ShopifyAPIVersion   string `env:"SHOPIFY_API_VERSION"`

	WebhookBaseURL string `env:"WEBHOOK_BASE_URL"`

	WooStoreAWSRegion  string `env:"WOO_STORE_AWS_REGION"`
	WooStoreAWSKey     string `env:"WOO_STORE_AWS_KEY"`
	WooStoreAWSSecret  string `env:"WOO_STORE_AWS_SECRET"`
	WooStoreInstanceID string `env:"WOO_STORE_INSTANCE_ID"`
	WooStoreURL        string `env:"WOO_STORE_URL"`

	SoftpymesAPIURL string `env:"SOFTPYMES_API_URL"`

	GoogleMapsAPIKey string `env:"GOOGLE_MAPS_API_KEY"`

	BoldIdentityKey string `env:"BOLD_IDENTITY_KEY"`
	BoldSecretKey   string `env:"BOLD_SECRET_KEY"`

	GoogleOAuthClientID     string `env:"GOOGLE_OAUTH_CLIENT_ID"`
	GoogleOAuthClientSecret string `env:"GOOGLE_OAUTH_CLIENT_SECRET"`
	GoogleOAuthRedirectURI  string `env:"GOOGLE_OAUTH_REDIRECT_URI"`

	FCMProjectID       string `env:"FCM_PROJECT_ID"`
	FCMCredentialsFile string `env:"FCM_CREDENTIALS_FILE"`
	FCMCredentialsJSON string `env:"FCM_CREDENTIALS_JSON"`

	BedrockAccessKey string `env:"BEDROCK_ACCESS_KEY"`
	BedrockSecretKey string `env:"BEDROCK_SECRET_KEY"`
	BedrockRegion    string `env:"BEDROCK_REGION"`
}

func splitTag(tag string) []string {

	parts := make([]string, 0, 2)
	for i, c := range tag {
		if c == ',' {
			parts = append(parts, tag[:i], tag[i+1:])
			return parts
		}
	}
	return []string{tag}
}
