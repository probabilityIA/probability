package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/secamc93/probability/back/testing/modules/orders/internal/domain/entities"
	"github.com/secamc93/probability/back/testing/modules/orders/internal/domain/ports"
)

type CentralClient struct {
	baseURL    string
	httpClient *http.Client
}

func New(centralAPIURL string) ports.ICentralClient {
	return &CentralClient{
		baseURL: centralAPIURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *CentralClient) GetBaseURL() string {
	return c.baseURL
}

func (c *CentralClient) CreateOrder(ctx context.Context, token string, orderPayload map[string]interface{}) (*entities.CreatedOrder, *entities.APICallLog, error) {
	url := c.baseURL + "/api/v1/orders"

	body, err := json.Marshal(orderPayload)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal order payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	start := time.Now()
	resp, err := c.httpClient.Do(req)
	durationMs := time.Since(start).Milliseconds()

	if err != nil {
		apiLog := &entities.APICallLog{
			Success:    false,
			Timestamp:  start.Format(time.RFC3339),
			DurationMs: durationMs,
			Request: entities.APIRequest{
				Method: http.MethodPost,
				URL:    url,
				Body:   orderPayload,
			},
			Response: entities.APIResponse{
				StatusCode: 0,
				Body:       err.Error(),
			},
		}
		return nil, apiLog, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read response: %w", err)
	}

	apiLog := &entities.APICallLog{
		Timestamp:  start.Format(time.RFC3339),
		DurationMs: durationMs,
		Request: entities.APIRequest{
			Method: http.MethodPost,
			URL:    url,
			Body:   orderPayload,
		},
		Response: entities.APIResponse{
			StatusCode: resp.StatusCode,
			Body:       string(respBody),
		},
	}

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		apiLog.Success = false
		return nil, apiLog, fmt.Errorf("central API returned %d: %s", resp.StatusCode, string(respBody))
	}

	var result struct {
		Success bool `json:"success"`
		Data    struct {
			ID           string  `json:"ID"`
			OrderNumber  string  `json:"OrderNumber"`
			TotalAmount  float64 `json:"TotalAmount"`
			CustomerName string  `json:"CustomerName"`
		} `json:"data"`
		Message string `json:"message"`
	}

	if err := json.Unmarshal(respBody, &result); err != nil {
		apiLog.Success = false
		return nil, apiLog, fmt.Errorf("failed to parse response: %w", err)
	}

	apiLog.Success = true

	return &entities.CreatedOrder{
		ID:           result.Data.ID,
		OrderNumber:  result.Data.OrderNumber,
		Total:        result.Data.TotalAmount,
		CustomerName: result.Data.CustomerName,
	}, apiLog, nil
}
