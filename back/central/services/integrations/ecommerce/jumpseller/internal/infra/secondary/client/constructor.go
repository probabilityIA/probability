package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/jumpseller/internal/domain"
)

const (
	requestTimeout      = 30 * time.Second
	maxBodyBytes        = 8 << 20
	maxRateLimitRetries = 4
)

type JumpsellerClient struct {
	httpClient *http.Client
	pacers     *pacerRegistry
}

func New() domain.IJumpsellerClient {
	return &JumpsellerClient{
		httpClient: &http.Client{Timeout: requestTimeout},
		pacers:     newPacerRegistry(),
	}
}

func baseURL(cred domain.Credential) (string, error) {
	if cred.BaseURL == "" {
		return "", domain.ErrMissingBaseURL
	}
	return strings.TrimRight(cred.BaseURL, "/"), nil
}

func friendlyConnError(err error) error {
	var dnsErr *net.DNSError
	if errors.As(err, &dnsErr) {
		return fmt.Errorf("no pudimos resolver la direccion de Jumpseller. Revisa tu conexion a internet")
	}

	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return fmt.Errorf("Jumpseller tardo demasiado en responder. Intenta de nuevo en unos minutos")
	}

	msg := strings.ToLower(err.Error())
	switch {
	case strings.Contains(msg, "no such host"):
		return fmt.Errorf("no pudimos resolver la direccion de Jumpseller. Revisa tu conexion a internet")
	case strings.Contains(msg, "connection refused"):
		return fmt.Errorf("Jumpseller rechazo la conexion. Intenta de nuevo en unos minutos")
	case strings.Contains(msg, "timeout") || strings.Contains(msg, "deadline exceeded"):
		return fmt.Errorf("Jumpseller tardo demasiado en responder. Intenta de nuevo en unos minutos")
	case strings.Contains(msg, "certificate") || strings.Contains(msg, "x509") || strings.Contains(msg, "tls"):
		return fmt.Errorf("el certificado de seguridad de Jumpseller no es valido")
	default:
		return fmt.Errorf("no pudimos conectarnos con Jumpseller. Intenta de nuevo en unos minutos")
	}
}

func (c *JumpsellerClient) do(ctx context.Context, cred domain.Credential, method, path string, query url.Values, body interface{}) ([]byte, error) {
	root, err := baseURL(cred)
	if err != nil {
		return nil, err
	}

	endpoint := root + path
	if len(query) > 0 {
		endpoint = endpoint + "?" + query.Encode()
	}

	pacer := c.pacers.forStore(cred.APIKey)

	var lastErr error
	for attempt := 0; attempt <= maxRateLimitRetries; attempt++ {
		if werr := pacer.wait(ctx); werr != nil {
			return nil, werr
		}

		raw, status, headers, err := c.attempt(ctx, cred, method, endpoint, body)
		if err != nil {
			return nil, err
		}

		pacer.observeLimit(headers.Get(domain.RateLimitHeader))

		if isRateLimited(status, headers) {
			pacer.backOff()
			lastErr = domain.ErrRateLimited
			continue
		}

		if status >= 500 {
			lastErr = fmt.Errorf("jumpseller client: la tienda respondio %d", status)
			continue
		}

		return c.interpret(raw, status)
	}

	return nil, fmt.Errorf("jumpseller client: se agotaron los reintentos: %w", lastErr)
}

func isRateLimited(status int, headers http.Header) bool {
	if status == http.StatusTooManyRequests {
		return true
	}
	return status == http.StatusForbidden && headers.Get(domain.RateLimitHeader) != ""
}

func (c *JumpsellerClient) attempt(ctx context.Context, cred domain.Credential, method, endpoint string, body interface{}) ([]byte, int, http.Header, error) {
	var reader io.Reader
	if body != nil {
		payload, err := json.Marshal(body)
		if err != nil {
			return nil, 0, nil, fmt.Errorf("jumpseller client: marshaling body: %w", err)
		}
		reader = bytes.NewReader(payload)
	}

	req, err := http.NewRequestWithContext(ctx, method, endpoint, reader)
	if err != nil {
		return nil, 0, nil, fmt.Errorf("jumpseller client: creating request: %w", err)
	}

	req.SetBasicAuth(cred.APIKey, cred.APISecret)
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, 0, nil, friendlyConnError(err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(io.LimitReader(resp.Body, maxBodyBytes))
	if err != nil {
		return nil, 0, nil, fmt.Errorf("jumpseller client: reading response: %w", err)
	}

	return raw, resp.StatusCode, resp.Header, nil
}

func (c *JumpsellerClient) interpret(raw []byte, status int) ([]byte, error) {
	switch {
	case status == http.StatusUnauthorized, status == http.StatusForbidden:
		return nil, domain.ErrInvalidCredentials
	case status == http.StatusNotFound:
		return nil, domain.ErrNoOrdersFound
	case status >= 400:
		return nil, fmt.Errorf("jumpseller client: unexpected status %d: %s", status, string(raw))
	}
	return raw, nil
}
