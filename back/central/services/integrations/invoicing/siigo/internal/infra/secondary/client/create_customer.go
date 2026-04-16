package client

import (
	"context"
	"fmt"
	"strings"

	"github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/infra/secondary/client/mappers"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/infra/secondary/client/request"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/infra/secondary/client/response"
)

// CreateCustomer crea un cliente en Siigo
// Endpoint: POST /v1/customers
func (c *Client) CreateCustomer(ctx context.Context, credentials dtos.Credentials, req *dtos.CreateCustomerRequest) (*dtos.CustomerResult, error) {
	c.log.Info(ctx).
		Str("identification", req.Identification).
		Str("name", req.Name).
		Msg("👤 Creating Siigo customer")

	// Autenticar
	token, err := c.authenticate(ctx, credentials.Username, credentials.AccessKey, credentials.AccountID, credentials.PartnerID, credentials.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to authenticate with Siigo: %w", err)
	}

	// Construir nombre en array
	nameParts := strings.Fields(req.Name)
	if len(nameParts) == 0 {
		nameParts = []string{"Sin Nombre"}
	}

	// Tipo de documento
	idType := req.IDType
	if idType == "" {
		idType = "13" // CC por defecto
	}

	// Tipo de persona
	personType := req.PersonType
	if personType == "" {
		personType = "Person"
	}

	customerBody := struct {
		PersonType     string              `json:"person_type"`
		IDType         request.SiigoIDType `json:"id_type"`
		Identification string              `json:"identification"`
		Name           []string            `json:"name"`
		Address        *request.SiigoAddress  `json:"address,omitempty"`
		Phones         []request.SiigoPhone   `json:"phones,omitempty"`
		Contacts       []request.SiigoContact `json:"contacts,omitempty"`
	}{
		PersonType:     personType,
		IDType:         request.SiigoIDType{Code: idType},
		Identification: req.Identification,
		Name:           nameParts,
	}

	if req.Address != "" {
		customerBody.Address = &request.SiigoAddress{Address: req.Address}
	}

	if req.Phone != "" {
		customerBody.Phones = []request.SiigoPhone{{Number: req.Phone}}
	}

	if req.Email != "" {
		firstName := nameParts[0]
		customerBody.Contacts = []request.SiigoContact{
			{
				FirstName:             firstName,
				Email:                 req.Email,
				SendElectronicInvoice: true,
			},
		}
	}

	var customerResp response.Customer

	resp, err := c.httpClient.R().
		SetContext(ctx).
		SetAuthToken(token).
		SetHeader("Partner-Id", credentials.PartnerID).
		SetBody(customerBody).
		SetResult(&customerResp).
		Post(c.endpointURL(credentials.BaseURL, "/v1/customers"))

	if err != nil {
		c.log.Error(ctx).Err(err).Msg("❌ Siigo create customer request failed - network error")
		return nil, fmt.Errorf("error de red al crear cliente en Siigo: %w", err)
	}

	c.log.Info(ctx).
		Int("status_code", resp.StatusCode()).
		Str("customer_id", customerResp.ID).
		Msg("📥 Siigo create customer response received")

	if resp.IsError() {
		c.log.Error(ctx).
			Int("status", resp.StatusCode()).
			Str("body", string(resp.Body())).
			Msg("❌ Siigo create customer failed")
		return nil, fmt.Errorf("error al crear cliente en Siigo (código %d): %s", resp.StatusCode(), string(resp.Body()))
	}

	result := mappers.CustomerToDTO(&customerResp)
	c.log.Info(ctx).
		Str("customer_id", result.ID).
		Msg("✅ Siigo customer created successfully")

	return result, nil
}
