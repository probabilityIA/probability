package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/woocommerce/internal/domain"
)

type WooCommerceClient struct {
	httpClient *http.Client
}

func New() domain.IWooCommerceClient {
	return &WooCommerceClient{
		httpClient: &http.Client{},
	}
}

func friendlyConnError(err error) error {
	var dnsErr *net.DNSError
	if errors.As(err, &dnsErr) {
		return fmt.Errorf("no encontramos ninguna tienda en esa direccion. Revisa que la URL este bien escrita (por ejemplo: https://mitienda.com)")
	}

	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return fmt.Errorf("la tienda tardo demasiado en responder. Verifica que este en linea e intenta de nuevo")
	}

	msg := strings.ToLower(err.Error())
	switch {
	case strings.Contains(msg, "no such host"):
		return fmt.Errorf("no encontramos ninguna tienda en esa direccion. Revisa que la URL este bien escrita (por ejemplo: https://mitienda.com)")
	case strings.Contains(msg, "connection refused"):
		return fmt.Errorf("la tienda rechazo la conexion. Verifica que este en linea y que la URL sea correcta")
	case strings.Contains(msg, "timeout") || strings.Contains(msg, "deadline exceeded"):
		return fmt.Errorf("la tienda tardo demasiado en responder. Verifica que este en linea e intenta de nuevo")
	case strings.Contains(msg, "certificate") || strings.Contains(msg, "x509") || strings.Contains(msg, "tls"):
		return fmt.Errorf("el certificado de seguridad (HTTPS) de la tienda no es valido. Verifica que la tienda tenga un certificado SSL valido")
	case strings.Contains(msg, "unsupported protocol") || strings.Contains(msg, "missing protocol") || strings.Contains(msg, "no scheme"):
		return fmt.Errorf("la URL de la tienda debe empezar con https:// (por ejemplo: https://mitienda.com)")
	default:
		return fmt.Errorf("no pudimos conectarnos con la tienda. Verifica que la URL sea correcta y que el sitio este en linea")
	}
}

func (c *WooCommerceClient) TestConnection(ctx context.Context, storeURL, consumerKey, consumerSecret string) error {
	storeURL = strings.TrimRight(storeURL, "/")
	endpoint := fmt.Sprintf("%s/wp-json/wc/v3/system_status", storeURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return fmt.Errorf("la URL de la tienda no es valida. Revisa que este bien escrita (por ejemplo: https://mitienda.com)")
	}

	req.SetBasicAuth(consumerKey, consumerSecret)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return friendlyConnError(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return domain.ErrInvalidCredentials
	}

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("la tienda respondio pero no encontramos la API de WooCommerce. Verifica que la tienda tenga WooCommerce instalado y los enlaces permanentes activos")
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("la tienda respondio con un estado inesperado (%d). Intenta de nuevo o verifica que la tienda funcione correctamente", resp.StatusCode)
	}

	body, _ := io.ReadAll(io.LimitReader(resp.Body, 2<<20))

	contentType := strings.ToLower(resp.Header.Get("Content-Type"))
	if !strings.Contains(contentType, "application/json") {
		return fmt.Errorf("la tienda respondio con una pagina web en lugar de datos. Verifica que la API REST este habilitada y que la tienda NO este en modo 'Proximamente' o mantenimiento")
	}

	var apiError struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}
	if err := json.Unmarshal(body, &apiError); err == nil && apiError.Code != "" {
		if strings.Contains(apiError.Code, "authentication") || strings.Contains(apiError.Code, "cannot_view") || strings.Contains(apiError.Code, "unauthorized") {
			return domain.ErrInvalidCredentials
		}
		return fmt.Errorf("la tienda respondio con un error: %s", apiError.Message)
	}

	var systemStatus struct {
		Environment json.RawMessage `json:"environment"`
	}
	if err := json.Unmarshal(body, &systemStatus); err != nil || len(systemStatus.Environment) == 0 {
		return domain.ErrInvalidCredentials
	}

	return nil
}
