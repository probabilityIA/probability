package client

import (
	"context"
	"encoding/json"
	"fmt"
)

// CustomerSearchResponse representa la respuesta de búsqueda de cliente en Softpymes
type CustomerSearchResponse struct {
	Message        string `json:"message,omitempty"`        // "No existe cliente con identificación: XXX"
	Identification string `json:"identification,omitempty"` // NIT del cliente
	Name           string `json:"name,omitempty"`           // Nombre del cliente
	BranchCode     string `json:"branchCode,omitempty"`     // Código de sucursal del cliente
	Email          string `json:"email,omitempty"`
	BillAddress    string `json:"billAddress,omitempty"`
}

// ensureCustomerExists verifica que el cliente exista en Softpymes y lo crea si no existe
// Endpoint de búsqueda: GET /app/integration/customer?identification=XXX
// Endpoint de creación: POST /app/integration/customer_new/
func (c *Client) ensureCustomerExists(ctx context.Context, token, referer, customerNit string, customer map[string]interface{}, config map[string]interface{}) error {
	// Buscar si el cliente ya existe en Softpymes
	exists, err := c.customerExists(ctx, token, referer, customerNit)
	if err != nil {
		return fmt.Errorf("error checking if customer exists: %w", err)
	}

	if exists {
		c.log.Info(ctx).
			Str("customer_nit", customerNit).
			Msg("✅ Customer already exists in Softpymes")
		return nil
	}

	// El cliente no existe, crearlo
	c.log.Info(ctx).
		Str("customer_nit", customerNit).
		Msg("📝 Customer does not exist in Softpymes, creating...")

	return c.createCustomer(ctx, token, referer, customerNit, customer, config)
}

// customerExists verifica si un cliente existe en Softpymes por su NIT
func (c *Client) customerExists(ctx context.Context, token, referer, customerNit string) (bool, error) {
	resp, err := c.httpClient.R().
		SetContext(ctx).
		SetAuthToken(token).
		SetHeader("Referer", referer).
		SetQueryParam("identification", customerNit).
		Get("/app/integration/customer")

	if err != nil {
		return false, fmt.Errorf("customer search request failed: %w", err)
	}

	if resp.IsError() {
		return false, fmt.Errorf("customer search failed with status %d", resp.StatusCode())
	}

	// Softpymes retorna 200 con un objeto que tiene "message" si no existe
	// o un objeto con "identification" si existe
	var result map[string]interface{}
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return false, fmt.Errorf("error parsing customer search response: %w", err)
	}

	// Si tiene "identification", el cliente existe
	if _, ok := result["identification"]; ok {
		return true, nil
	}

	// Si solo tiene "message", el cliente no existe
	return false, nil
}

// createCustomer crea un nuevo cliente en Softpymes
// Usa el endpoint POST /app/integration/customer_new/ (agregar cliente a tercero)
func (c *Client) createCustomer(ctx context.Context, token, referer, customerNit string, customer map[string]interface{}, config map[string]interface{}) error {
	// Extraer email del customer (si existe)
	email := "noreply@probability.com"
	if customer != nil {
		if e, ok := customer["email"].(string); ok && e != "" {
			email = e
		}
	}

	// Extraer dirección del customer (si existe)
	address := "COLOMBIA"
	if customer != nil {
		if addr, ok := customer["address"].(string); ok && addr != "" {
			address = addr
		}
	}

	// Obtener companyNit del config (referer sin dígito de verificación)
	companyNit := referer
	if cn, ok := config["company_nit"].(string); ok && cn != "" {
		companyNit = cn
	}

	// Construir request de creación de cliente
	// Docs: https://api-integracion.softpymes.com.co/doc/#api-Clientes-AddCustomers
	customerReq := map[string]interface{}{
		"identificationNumber": customerNit,
		"phone":                "",
		"cellPhone":            "",
		"branchName":           "PRINCIPAL",
		"billAddress":          address,
		"email":                email,
		"cityCode":             "001",          // Bogotá por defecto
		"departmentCode":       "11",           // D.C. por defecto
		"companyNit":           companyNit,
	}

	c.log.Info(ctx).
		Str("customer_nit", customerNit).
		Str("email", email).
		Str("company_nit", companyNit).
		Msg("📤 Creating customer in Softpymes")

	resp, err := c.httpClient.R().
		SetContext(ctx).
		SetAuthToken(token).
		SetHeader("Referer", referer).
		SetBody(customerReq).
		Post("/app/integration/customer_new/")

	if err != nil {
		return fmt.Errorf("customer creation request failed: %w", err)
	}

	if resp.IsError() {
		c.log.Error(ctx).
			Int("status", resp.StatusCode()).
			Str("response", string(resp.Body())).
			Msg("❌ Customer creation failed in Softpymes")
		return fmt.Errorf("customer creation failed with status %d: %s", resp.StatusCode(), string(resp.Body()))
	}

	c.log.Info(ctx).
		Str("customer_nit", customerNit).
		Str("response", string(resp.Body())).
		Msg("✅ Customer created successfully in Softpymes")

	return nil
}
