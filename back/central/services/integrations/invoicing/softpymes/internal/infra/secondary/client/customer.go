package client

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes/internal/domain/dtos"
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
func (c *Client) ensureCustomerExists(ctx context.Context, token, referer, customerNit string, customer *dtos.CustomerData, config map[string]interface{}) error {
	// Buscar si el cliente ya existe en Softpymes
	exists, err := c.customerExists(ctx, token, referer, customerNit)
	if err != nil {
		return fmt.Errorf("error checking if customer exists: %w", err)
	}

	if exists {
		c.log.Info(ctx).
			Str("customer_nit", customerNit).
			Msg("Customer already exists in Softpymes")
		return nil
	}

	// El cliente no existe, crearlo
	c.log.Info(ctx).
		Str("customer_nit", customerNit).
		Msg("Customer does not exist in Softpymes, creating...")

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

// createCustomer crea un nuevo tercero y cliente en Softpymes
// Usa el endpoint POST /app/integration/customer (guardar tercero y cliente)
// Docs: https://api-integracion.softpymes.com.co/doc/#api-Clientes-SaveCustomers
func (c *Client) createCustomer(ctx context.Context, token, referer, customerNit string, customer *dtos.CustomerData, config map[string]interface{}) error {
	// Extraer datos del customer tipado
	email := "noreply@probability.com"
	if customer != nil && customer.Email != "" {
		email = customer.Email
	}

	address := "COLOMBIA"
	if customer != nil && customer.Address != "" {
		address = customer.Address
	}

	phone := ""
	if customer != nil && customer.Phone != "" {
		phone = customer.Phone
	}

	customerName := ""
	if customer != nil && customer.Name != "" {
		customerName = customer.Name
	}

	// Obtener companyNit del config
	companyNit := referer
	if cn, ok := config["company_nit"].(string); ok && cn != "" {
		companyNit = cn
	}

	// Separar nombre en partes para persona natural
	firstName, lastName := splitCustomerName(customerName)

	// Construir request - Persona Natural (thirdType = "N")
	// maidenName y otherName son requeridos por Softpymes para persona natural
	customerReq := map[string]interface{}{
		"identificationNumber": customerNit,
		"identificationType":   "13",
		"thirdType":            "N",
		"firstName":            firstName,
		"lastName":             lastName,
		"maidenName":           ".",
		"otherName":            ".",
		"phone":                phone,
		"cellPhone":            phone,
		"branchName":           "PRINCIPAL",
		"billAddress":          address,
		"email":                email,
		"cityCode":             "001",
		"departmentCode":       "11",
		"companyNit":           companyNit,
	}

	c.log.Info(ctx).
		Str("customer_nit", customerNit).
		Str("first_name", firstName).
		Str("last_name", lastName).
		Str("email", email).
		Str("phone", phone).
		Str("company_nit", companyNit).
		Msg("Creating customer (tercero + cliente) in Softpymes")

	resp, err := c.httpClient.R().
		SetContext(ctx).
		SetAuthToken(token).
		SetHeader("Referer", referer).
		SetBody(customerReq).
		SetDebug(true).
		Post("/app/integration/customer")

	if err != nil {
		return fmt.Errorf("customer creation request failed: %w", err)
	}

	if resp.IsError() {
		c.log.Error(ctx).
			Int("status", resp.StatusCode()).
			Str("response", string(resp.Body())).
			Msg("Customer creation failed in Softpymes")
		return fmt.Errorf("customer creation failed with status %d: %s", resp.StatusCode(), string(resp.Body()))
	}

	// Softpymes puede retornar errores dentro de un HTTP 200
	// Formato: [[{"message":"error msg","type":"ERROR_TYPE"}], 400]
	if err := c.checkResponseForErrors(resp.Body(), "customer creation"); err != nil {
		c.log.Error(ctx).
			Str("customer_nit", customerNit).
			Str("response", string(resp.Body())).
			Msg("Customer creation returned error in response body")
		return err
	}

	c.log.Info(ctx).
		Str("customer_nit", customerNit).
		Str("name", customerName).
		Str("response", string(resp.Body())).
		Msg("Customer created successfully in Softpymes")

	return nil
}

// checkResponseForErrors detecta errores de Softpymes embebidos en respuestas HTTP 200
// Softpymes retorna errores en formato: [[{"message":"...","type":"..."},...], statusCode]
func (c *Client) checkResponseForErrors(body []byte, operation string) error {
	// Intentar parsear como array: [[errors], statusCode]
	var rawArray []json.RawMessage
	if err := json.Unmarshal(body, &rawArray); err != nil {
		// No es un array, probablemente es una respuesta exitosa (objeto JSON)
		return nil
	}

	// Si tiene exactamente 2 elementos y el segundo es un número >= 400, es un error
	if len(rawArray) != 2 {
		return nil
	}

	var statusCode int
	if err := json.Unmarshal(rawArray[1], &statusCode); err != nil {
		return nil
	}

	if statusCode < 400 {
		return nil
	}

	// Es un error - extraer mensajes
	var errors []map[string]interface{}
	if err := json.Unmarshal(rawArray[0], &errors); err != nil {
		return fmt.Errorf("%s failed (status %d): %s", operation, statusCode, string(rawArray[0]))
	}

	var messages []string
	for _, e := range errors {
		if msg, ok := e["message"].(string); ok {
			messages = append(messages, msg)
		}
	}

	return fmt.Errorf("%s failed (status %d): %s", operation, statusCode, strings.Join(messages, "; "))
}

// splitCustomerName separa un nombre completo en firstName y lastName
// Ejemplos:
//
//	"Sebastian Camacho"       → "Sebastian", "Camacho"
//	"Juan Carlos Pérez López" → "Juan Carlos", "Pérez López"
//	"Sebastian"               → "Sebastian", "."
//	""                        → "Cliente", "Probability"
func splitCustomerName(fullName string) (string, string) {
	fullName = strings.TrimSpace(fullName)
	if fullName == "" {
		return "Cliente", "Probability"
	}

	parts := strings.Fields(fullName)
	switch len(parts) {
	case 1:
		return parts[0], "."
	case 2:
		return parts[0], parts[1]
	case 3:
		// "Juan Carlos Pérez" → firstName="Juan Carlos", lastName="Pérez"
		return parts[0] + " " + parts[1], parts[2]
	default:
		// "Juan Carlos Pérez López" → firstName="Juan Carlos", lastName="Pérez López"
		mid := len(parts) / 2
		return strings.Join(parts[:mid], " "), strings.Join(parts[mid:], " ")
	}
}
