package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/woocommerce/internal/domain"
)

func (c *WooCommerceClient) UpdateProductStock(ctx context.Context, storeURL, consumerKey, consumerSecret, productExternalID string, quantity int) error {
	storeURL = strings.TrimRight(storeURL, "/")
	endpoint := fmt.Sprintf("%s/wp-json/wc/v3/products/%s", storeURL, productExternalID)

	payload := map[string]interface{}{
		"manage_stock":   true,
		"stock_quantity": quantity,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("woocommerce client: marshaling stock payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, endpoint, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("woocommerce client: creating request: %w", err)
	}
	req.SetBasicAuth(consumerKey, consumerSecret)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("woocommerce client: stock update request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return domain.ErrInvalidCredentials
	}

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("woocommerce client: producto %s no encontrado en la tienda", productExternalID)
	}

	if resp.StatusCode != http.StatusOK {
		raw, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return fmt.Errorf("woocommerce client: estado inesperado %d al actualizar stock: %s", resp.StatusCode, string(raw))
	}

	return nil
}
