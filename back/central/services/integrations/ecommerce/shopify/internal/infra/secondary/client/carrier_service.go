package client

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

type CreateCarrierServiceResponse struct {
	CarrierService struct {
		ID               int64  `json:"id"`
		Name             string `json:"name"`
		Active           bool   `json:"active"`
		ServiceDiscovery bool   `json:"service_discovery"`
		CallbackURL      string `json:"callback_url"`
		Format           string `json:"format"`
	} `json:"carrier_service"`
}

func (c *shopifyClient) CreateCarrierService(ctx context.Context, storeName, accessToken, callbackURL, name string) (string, error) {
	url := buildURL(storeName, "/admin/api/2024-10/carrier_services.json")

	payload := map[string]interface{}{
		"carrier_service": map[string]interface{}{
			"name":              name,
			"callback_url":      callbackURL,
			"service_discovery": true,
			"format":            "json",
		},
	}

	var result CreateCarrierServiceResponse

	resp, err := c.client.R().
		SetContext(ctx).
		SetHeader("X-Shopify-Access-Token", accessToken).
		SetHeader("Content-Type", "application/json").
		SetBody(payload).
		SetResult(&result).
		Post(url)

	if err != nil {
		return "", fmt.Errorf("no se pudo conectar con la tienda Shopify. Verifica el dominio y la conexion a internet")
	}

	if resp.StatusCode() == http.StatusCreated || resp.StatusCode() == http.StatusOK {
		return fmt.Sprintf("%d", result.CarrierService.ID), nil
	}

	return "", friendlyCarrierServiceError(resp.StatusCode(), resp.String(), storeName)
}

func (c *shopifyClient) DeleteCarrierService(ctx context.Context, storeName, accessToken, carrierServiceID string) error {
	url := buildURL(storeName, fmt.Sprintf("/admin/api/2024-10/carrier_services/%s.json", carrierServiceID))

	resp, err := c.client.R().
		SetContext(ctx).
		SetHeader("X-Shopify-Access-Token", accessToken).
		Delete(url)

	if err != nil {
		return fmt.Errorf("no se pudo conectar con la tienda Shopify. Verifica el dominio y la conexion a internet")
	}

	if resp.StatusCode() == http.StatusOK || resp.StatusCode() == http.StatusNoContent || resp.StatusCode() == http.StatusNotFound {
		return nil
	}

	return friendlyCarrierServiceError(resp.StatusCode(), resp.String(), storeName)
}

func friendlyCarrierServiceError(statusCode int, body, storeName string) error {
	low := strings.ToLower(body)

	switch {
	case strings.Contains(low, "write_shipping") || strings.Contains(low, "merchant approval") || strings.Contains(low, "carrier_calculated") || strings.Contains(low, "carrier service") && strings.Contains(low, "not enabled"):
		return fmt.Errorf("la tienda no tiene habilitada la cotizacion por terceros. Requiere: 1) que la Custom App tenga el permiso 'write_shipping' (reinstalar/autorizar la app con ese scope) y 2) un plan Shopify compatible (Advanced/Plus o el complemento de envio calculado por terceros)")
	case strings.Contains(low, "unavailable shop") || statusCode == http.StatusPaymentRequired:
		return fmt.Errorf("la tienda Shopify no esta disponible (plan inactivo o tienda suspendida). Verifica el estado de la tienda")
	case strings.Contains(low, "already") && strings.Contains(low, "carrier"):
		return fmt.Errorf("ya existe un servicio de transporte registrado para esta tienda")
	}

	switch statusCode {
	case http.StatusUnauthorized:
		return fmt.Errorf("token de acceso invalido o expirado. Vuelve a conectar la tienda Shopify")
	case http.StatusForbidden:
		return fmt.Errorf("acceso denegado por Shopify. La Custom App necesita el permiso 'write_shipping'")
	case http.StatusNotFound:
		return fmt.Errorf("tienda no encontrada en Shopify: %s", storeName)
	case http.StatusUnprocessableEntity:
		return fmt.Errorf("Shopify rechazo el registro del servicio de transporte: %s", strings.TrimSpace(body))
	case http.StatusTooManyRequests:
		return fmt.Errorf("demasiadas solicitudes a Shopify. Intenta nuevamente en unos minutos")
	default:
		return fmt.Errorf("no se pudo registrar el servicio de transporte en Shopify (codigo %d)", statusCode)
	}
}
