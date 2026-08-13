package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/woocommerce/internal/domain"
)

type wooStockResponse struct {
	ID            int64 `json:"id"`
	StockQuantity *int  `json:"stock_quantity"`
	ManageStock   bool  `json:"manage_stock"`
}

func (c *WooCommerceClient) GetProductsStock(ctx context.Context, storeURL, consumerKey, consumerSecret string, externalIDs []string) ([]domain.ChannelStock, error) {
	storeURL = strings.TrimRight(storeURL, "/")

	simples := make([]string, 0, len(externalIDs))
	porPadre := make(map[string][]string)
	vistos := make(map[string]bool, len(externalIDs))

	for _, ref := range externalIDs {
		limpio := strings.TrimSpace(ref)
		if limpio == "" || vistos[limpio] {
			continue
		}
		vistos[limpio] = true
		if padre, variacion, ok := splitVariationRef(limpio); ok {
			porPadre[padre] = append(porPadre[padre], variacion)
			continue
		}
		simples = append(simples, limpio)
	}

	stocks := make([]domain.ChannelStock, 0, len(vistos))

	for inicio := 0; inicio < len(simples); inicio += 100 {
		fin := inicio + 100
		if fin > len(simples) {
			fin = len(simples)
		}
		endpoint := fmt.Sprintf("%s/wp-json/wc/v3/products?per_page=100&status=any&include=%s",
			storeURL, strings.Join(simples[inicio:fin], ","))

		filas, err := c.consultarStock(ctx, endpoint, consumerKey, consumerSecret)
		if err != nil {
			return nil, err
		}
		for _, fila := range filas {
			stocks = append(stocks, aStock(fmt.Sprintf("%d", fila.ID), fila))
		}
	}

	for padre, variaciones := range porPadre {
		for inicio := 0; inicio < len(variaciones); inicio += 100 {
			fin := inicio + 100
			if fin > len(variaciones) {
				fin = len(variaciones)
			}
			endpoint := fmt.Sprintf("%s/wp-json/wc/v3/products/%s/variations?per_page=100&include=%s",
				storeURL, padre, strings.Join(variaciones[inicio:fin], ","))

			filas, err := c.consultarStock(ctx, endpoint, consumerKey, consumerSecret)
			if err != nil {
				return nil, err
			}
			for _, fila := range filas {
				stocks = append(stocks, aStock(fmt.Sprintf("%s:%d", padre, fila.ID), fila))
			}
		}
	}

	return stocks, nil
}

func aStock(externalID string, fila wooStockResponse) domain.ChannelStock {
	cantidad := 0
	if fila.StockQuantity != nil {
		cantidad = *fila.StockQuantity
	}
	return domain.ChannelStock{
		ExternalID:        externalID,
		AvailableQuantity: cantidad,
		ManageStock:       fila.ManageStock,
		Found:             true,
	}
}

func (c *WooCommerceClient) consultarStock(ctx context.Context, endpoint, consumerKey, consumerSecret string) ([]wooStockResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("woocommerce client: creating stock request: %w", err)
	}
	req.SetBasicAuth(consumerKey, consumerSecret)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("woocommerce client: stock request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, domain.ErrInvalidCredentials
	}
	if resp.StatusCode != http.StatusOK {
		raw, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return nil, fmt.Errorf("woocommerce client: unexpected status %d reading stock: %s", resp.StatusCode, string(raw))
	}

	var filas []wooStockResponse
	if err := json.NewDecoder(resp.Body).Decode(&filas); err != nil {
		return nil, fmt.Errorf("woocommerce client: decoding stock response: %w", err)
	}
	return filas, nil
}
