package httpclient

import (
	"time"

	"github.com/secamc93/probability/back/central/shared/log"

	"github.com/go-resty/resty/v2"
)

// HTTPClientConfig contiene la configuración del cliente HTTP
type HTTPClientConfig struct {
	Timeout    time.Duration
	BaseURL    string
	RetryCount int
	RetryWait  time.Duration
	Debug      bool
}

// Client envuelve resty.Client con logging
type Client struct {
	rest   *resty.Client
	logger log.ILogger
}

// New crea un nuevo cliente HTTP con resty
func New(cfg HTTPClientConfig, logger log.ILogger) *Client {
	if cfg.Timeout == 0 {
		cfg.Timeout = 30 * time.Second
	}
	if cfg.RetryCount == 0 {
		cfg.RetryCount = 2
	}
	if cfg.RetryWait == 0 {
		cfg.RetryWait = 5 * time.Second
	}

	client := resty.New().
		SetTimeout(cfg.Timeout).
		SetRetryCount(cfg.RetryCount).
		SetRetryWaitTime(cfg.RetryWait).
		AddRetryCondition(func(r *resty.Response, err error) bool {
			return r.StatusCode() == 429
		})

	if cfg.BaseURL != "" {
		client.SetBaseURL(cfg.BaseURL)
	}

	if cfg.Debug {
		client.SetDebug(true)
	}

	return &Client{
		rest:   client,
		logger: logger,
	}
}

// R retorna una nueva instancia de Request para construir peticiones
func (c *Client) R() *resty.Request {
	return c.rest.R()
}

// SetDebug activa o desactiva el modo debug
func (c *Client) SetDebug(enable bool) *Client {
	c.rest.SetDebug(enable)
	return c
}

// SetBaseURL establece la URL base
func (c *Client) SetBaseURL(url string) *Client {
	c.rest.SetBaseURL(url)
	return c
}

// SetTimeout establece el timeout
func (c *Client) SetTimeout(timeout time.Duration) *Client {
	c.rest.SetTimeout(timeout)
	return c
}

// SetHeader establece un header global
func (c *Client) SetHeader(key, value string) *Client {
	c.rest.SetHeader(key, value)
	return c
}

// SetHeaders establece múltiples headers globales
func (c *Client) SetHeaders(headers map[string]string) *Client {
	c.rest.SetHeaders(headers)
	return c
}

// GetRestyClient retorna el cliente resty subyacente
func (c *Client) GetRestyClient() *resty.Client {
	return c.rest
}
