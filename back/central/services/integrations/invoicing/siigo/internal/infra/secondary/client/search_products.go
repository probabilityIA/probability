package client

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/domain/dtos"
)

const (
	busquedaPageSize     = 100
	busquedaMaxPages     = 40
	vistaPreviaMaxPages  = 2
	vistaPreviaResultado = 10
)

func normalizar(texto string) string {
	return strings.ToLower(strings.TrimSpace(texto))
}

func coincide(code, name, termino string) bool {
	if termino == "" {
		return true
	}
	return strings.Contains(normalizar(code), termino) || strings.Contains(normalizar(name), termino)
}

func (c *Client) SearchProducts(ctx context.Context, credentials dtos.Credentials, termino string, limite int) ([]dtos.ProductItem, error) {
	if limite <= 0 || limite > 100 {
		limite = 50
	}

	token, err := c.authenticate(ctx, credentials.Username, credentials.AccessKey, credentials.AccountID, credentials.PartnerID, credentials.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to authenticate with Siigo: %w", err)
	}

	buscado := normalizar(termino)
	maxPages := busquedaMaxPages
	if buscado == "" {
		maxPages = vistaPreviaMaxPages
		if limite > vistaPreviaResultado {
			limite = vistaPreviaResultado
		}
	}

	servicios := make([]dtos.ProductItem, 0, limite)
	productos := make([]dtos.ProductItem, 0, limite)
	revisados := 0
	total := 0

	for page := 1; page <= maxPages; page++ {
		var listResp listProductsResponse

		resp, err := c.httpClient.R().
			SetContext(ctx).
			SetAuthToken(token).
			SetHeader("Partner-Id", credentials.PartnerID).
			SetQueryParam("page", strconv.Itoa(page)).
			SetQueryParam("page_size", strconv.Itoa(busquedaPageSize)).
			SetResult(&listResp).
			Get(c.endpointURL(credentials.BaseURL, "/v1/products"))

		if err != nil {
			return nil, fmt.Errorf("error de red al buscar productos en Siigo: %w", err)
		}
		if resp.IsError() {
			return nil, fmt.Errorf("error al buscar productos en Siigo (codigo %d)", resp.StatusCode())
		}

		total = listResp.Pagination.TotalResults

		for _, r := range listResp.Results {
			revisados++
			if strings.TrimSpace(r.Code) == "" || !coincide(r.Code, r.Name, buscado) {
				continue
			}
			item := dtos.ProductItem{ID: r.ID, Code: r.Code, Type: r.Type, Name: r.Name}
			if strings.EqualFold(strings.TrimSpace(r.Type), "Service") {
				servicios = append(servicios, item)
				continue
			}
			productos = append(productos, item)
		}

		if len(servicios) >= limite {
			break
		}
		if buscado == "" && len(servicios)+len(productos) >= limite {
			break
		}
		if len(listResp.Results) < busquedaPageSize || (total > 0 && revisados >= total) {
			break
		}
	}

	resultado := append(servicios, productos...)
	if len(resultado) > limite {
		resultado = resultado[:limite]
	}

	c.log.Info(ctx).
		Str("termino", termino).
		Int("encontrados", len(resultado)).
		Int("servicios", len(servicios)).
		Int("revisados", revisados).
		Int("total_catalogo", total).
		Msg("Busqueda de productos en Siigo")

	return resultado, nil
}
