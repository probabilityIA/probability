package client

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/domain/dtos"
)

type invoicePDFResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Base64   string `json:"base64"`
	FileName string `json:"file_name"`
}

func (c *Client) GetInvoicePDF(ctx context.Context, credentials dtos.Credentials, invoiceID string) ([]byte, string, error) {
	c.log.Info(ctx).Str("siigo_invoice_id", invoiceID).Msg("Getting Siigo invoice PDF")

	token, err := c.authenticate(ctx, credentials.Username, credentials.AccessKey, credentials.AccountID, credentials.PartnerID, credentials.BaseURL)
	if err != nil {
		return nil, "", fmt.Errorf("failed to authenticate with Siigo: %w", err)
	}

	resp, err := c.httpClient.R().
		SetContext(ctx).
		SetAuthToken(token).
		SetHeader("Partner-Id", credentials.PartnerID).
		Get(c.endpointURL(credentials.BaseURL, "/v1/invoices/"+invoiceID+"/pdf"))

	if err != nil {
		c.log.Error(ctx).Err(err).Msg("Siigo invoice pdf request failed - network error")
		return nil, "", fmt.Errorf("error de red al descargar el PDF de la factura en Siigo: %w", err)
	}

	if resp.IsError() {
		c.log.Error(ctx).
			Int("status", resp.StatusCode()).
			Str("body", string(resp.Body())).
			Msg("Siigo invoice pdf request failed")
		if resp.StatusCode() == 404 {
			return nil, "", fmt.Errorf("la factura no existe en Siigo o todavia no tiene PDF disponible")
		}
		return nil, "", errorSiigo(resp.Body(), resp.StatusCode(), "la descarga del PDF")
	}

	var payload invoicePDFResponse
	if err := json.Unmarshal(resp.Body(), &payload); err != nil || payload.Base64 == "" {
		return nil, "", fmt.Errorf("Siigo no devolvio el PDF de la factura")
	}

	contenido, err := base64.StdEncoding.DecodeString(strings.TrimSpace(payload.Base64))
	if err != nil {
		return nil, "", fmt.Errorf("el PDF devuelto por Siigo no se pudo decodificar: %w", err)
	}

	nombre := payload.FileName
	if nombre == "" {
		nombre = payload.Name
	}
	if nombre == "" {
		nombre = "factura.pdf"
	}
	if !strings.HasSuffix(strings.ToLower(nombre), ".pdf") {
		nombre = nombre + ".pdf"
	}

	c.log.Info(ctx).
		Str("siigo_invoice_id", invoiceID).
		Int("bytes", len(contenido)).
		Msg("Siigo invoice PDF downloaded")

	return contenido, nombre, nil
}
