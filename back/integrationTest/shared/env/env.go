package env

import (
	"os"

	"github.com/joho/godotenv"
)

type IConfig interface {
	Get(key string) string
	GetWithDefault(key, defaultValue string) string
}

type config struct {
	values map[string]string
}

func New() IConfig {
	_ = godotenv.Load()

	c := &config{
		values: make(map[string]string),
	}

	// Cargar variables de entorno comunes
	c.values["WEBHOOK_BASE_URL"] = os.Getenv("WEBHOOK_BASE_URL")
	c.values["SHOPIFY_SHOP_DOMAIN"] = os.Getenv("SHOPIFY_SHOP_DOMAIN")
	c.values["SHOPIFY_API_VERSION"] = os.Getenv("SHOPIFY_API_VERSION")

	return c
}

func (c *config) Get(key string) string {
	if val, ok := c.values[key]; ok {
		return val
	}
	return os.Getenv(key)
}

func (c *config) GetWithDefault(key, defaultValue string) string {
	val := c.Get(key)
	if val == "" {
		return defaultValue
	}
	return val
}











