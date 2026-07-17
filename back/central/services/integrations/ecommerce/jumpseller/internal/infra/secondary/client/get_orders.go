package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/jumpseller/internal/domain"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/jumpseller/internal/infra/secondary/client/response"
)

const jumpsellerDateLayout = "2006-01-02"

func (c *JumpsellerClient) GetOrders(ctx context.Context, cred domain.Credential, params *domain.GetOrdersParams) (*domain.GetOrdersResult, [][]byte, error) {
	query := url.Values{}

	if params != nil {
		if params.After != nil && params.Before != nil {
			query.Set("dateFilter", "customDate")
			query.Set("initialDate", params.After.Format(jumpsellerDateLayout))
			query.Set("finalDate", params.Before.Format(jumpsellerDateLayout))
		}
		for _, status := range params.Statuses {
			query.Add("status_filters[]", status)
		}
		if params.Page > 0 {
			query.Set("page", strconv.Itoa(params.Page))
		}
		if params.PerPage > 0 {
			query.Set("limit", strconv.Itoa(params.PerPage))
		}
	}

	raw, err := c.do(ctx, cred, http.MethodGet, "/orders.json", query, nil)
	if err != nil {
		return nil, nil, err
	}

	var envelopes []json.RawMessage
	if err := json.Unmarshal(raw, &envelopes); err != nil {
		return nil, nil, fmt.Errorf("jumpseller client: parsing orders: %w", err)
	}

	orders := make([]domain.JumpsellerOrder, 0, len(envelopes))
	rawBytes := make([][]byte, 0, len(envelopes))
	for _, item := range envelopes {
		var envelope response.OrderEnvelope
		if err := json.Unmarshal(item, &envelope); err != nil {
			continue
		}
		orders = append(orders, envelope.Order.ToDomain())
		rawBytes = append(rawBytes, []byte(item))
	}

	return &domain.GetOrdersResult{Orders: orders, Count: len(orders)}, rawBytes, nil
}

func (c *JumpsellerClient) GetOrder(ctx context.Context, cred domain.Credential, orderID int64) (*domain.JumpsellerOrder, []byte, error) {
	raw, err := c.do(ctx, cred, http.MethodGet, fmt.Sprintf("/orders/%d.json", orderID), nil, nil)
	if err != nil {
		return nil, nil, err
	}

	var envelope response.OrderEnvelope
	if err := json.Unmarshal(raw, &envelope); err != nil {
		return nil, nil, fmt.Errorf("jumpseller client: parsing order: %w", err)
	}

	order := envelope.Order.ToDomain()
	return &order, raw, nil
}

func (c *JumpsellerClient) UpdateOrder(ctx context.Context, cred domain.Credential, orderID int64, fields domain.UpdateOrderFields) error {
	body := response.UpdateOrderRequest{
		Order: response.UpdateOrderFields{
			Status:          fields.Status,
			ShipmentStatus:  fields.ShipmentStatus,
			TrackingNumber:  fields.TrackingNumber,
			TrackingCompany: fields.TrackingCompany,
			AdditionalInfo:  fields.AdditionalInfo,
		},
	}
	_, err := c.do(ctx, cred, http.MethodPut, fmt.Sprintf("/orders/%d.json", orderID), nil, body)
	return err
}
