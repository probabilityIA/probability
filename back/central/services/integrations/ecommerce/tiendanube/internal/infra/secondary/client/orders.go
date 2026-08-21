package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/tiendanube/internal/domain"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/tiendanube/internal/infra/secondary/client/response"
)

const ordersPerPage = 50

func (c *TiendanubeClient) GetOrder(ctx context.Context, cred domain.Credential, orderID string) (*domain.TiendanubeOrder, []byte, error) {
	raw, _, err := c.do(ctx, cred, http.MethodGet, "/orders/"+orderID, nil, nil)
	if err != nil {
		return nil, nil, err
	}

	var order response.Order
	if err := json.Unmarshal(raw, &order); err != nil {
		return nil, nil, fmt.Errorf("tiendanube client: parsing order: %w", err)
	}

	return mapOrder(order), raw, nil
}

func (c *TiendanubeClient) GetOrders(ctx context.Context, cred domain.Credential, filters domain.OrderFilters) ([]domain.TiendanubeOrder, error) {
	var collected []domain.TiendanubeOrder

	page := 1
	for {
		query := url.Values{}
		query.Set("page", strconv.Itoa(page))
		query.Set("per_page", strconv.Itoa(ordersPerPage))
		if filters.CreatedAtMin != "" {
			query.Set("created_at_min", filters.CreatedAtMin)
		}
		if filters.CreatedAtMax != "" {
			query.Set("created_at_max", filters.CreatedAtMax)
		}
		if filters.Status != "" {
			query.Set("status", filters.Status)
		}
		if filters.PaymentStatus != "" {
			query.Set("payment_status", filters.PaymentStatus)
		}

		raw, _, err := c.do(ctx, cred, http.MethodGet, "/orders", query, nil)
		if err != nil {
			if errors.Is(err, domain.ErrResourceNotFound) {
				break
			}
			return nil, err
		}

		var orders []response.Order
		if err := json.Unmarshal(raw, &orders); err != nil {
			return nil, fmt.Errorf("tiendanube client: parsing orders: %w", err)
		}
		if len(orders) == 0 {
			break
		}

		for _, order := range orders {
			collected = append(collected, *mapOrder(order))
			if filters.Limit > 0 && len(collected) >= filters.Limit {
				return collected, nil
			}
		}

		if len(orders) < ordersPerPage {
			break
		}
		page++
	}

	return collected, nil
}

func mapOrder(order response.Order) *domain.TiendanubeOrder {
	items := make([]domain.TiendanubeOrderItem, 0, len(order.Products))
	for _, product := range order.Products {
		items = append(items, domain.TiendanubeOrderItem{
			ID:        product.ID.String(),
			ProductID: product.ProductID.String(),
			VariantID: product.VariantID.String(),
			Name:      product.Name,
			SKU:       product.SKU,
			Price:     parseAmount(product.Price),
			Quantity:  parseQuantity(product.Quantity),
			Weight:    parseAmount(product.Weight),
			ImageURL:  product.ImageURL,
		})
	}

	return &domain.TiendanubeOrder{
		ID:                    order.ID.String(),
		Number:                order.Number.String(),
		Currency:              order.Currency,
		Status:                order.Status,
		PaymentStatus:         order.PaymentStatus,
		ShippingStatus:        order.ShippingStatus,
		Subtotal:              parseAmount(order.Subtotal),
		Discount:              parseAmount(order.Discount),
		Total:                 parseAmount(order.Total),
		ShippingCost:          parseAmount(order.ShippingCostCustomer),
		Weight:                parseAmount(order.Weight),
		ContactName:           order.ContactName,
		ContactEmail:          order.ContactEmail,
		ContactPhone:          order.ContactPhone,
		ContactIdentification: order.ContactIdentification,
		Gateway:               order.Gateway,
		GatewayName:           order.GatewayName,
		Note:                  order.Note,
		ShippingOption:        order.ShippingOption,
		TrackingNumber:        order.ShippingTrackingNumber,
		TrackingURL:           order.ShippingTrackingURL,
		BillingAddress: domain.TiendanubeAddress{
			Name:     order.BillingName,
			Phone:    order.BillingPhone,
			Street:   order.BillingAddress,
			Number:   order.BillingNumber,
			Floor:    order.BillingFloor,
			Locality: order.BillingLocality,
			City:     order.BillingCity,
			Province: order.BillingProvince,
			Country:  order.BillingCountry,
			Zipcode:  order.BillingZipcode,
		},
		ShippingAddress: domain.TiendanubeAddress{
			FirstName: order.ShippingAddress.FirstName,
			LastName:  order.ShippingAddress.LastName,
			Name:      order.ShippingAddress.Name,
			Phone:     order.ShippingAddress.Phone,
			Street:    order.ShippingAddress.Address,
			Number:    order.ShippingAddress.Number,
			Floor:     order.ShippingAddress.Floor,
			Locality:  order.ShippingAddress.Locality,
			City:      order.ShippingAddress.City,
			Province:  order.ShippingAddress.Province,
			Country:   order.ShippingAddress.Country,
			Zipcode:   order.ShippingAddress.Zipcode,
		},
		Items:       items,
		CreatedAt:   order.CreatedAt,
		UpdatedAt:   order.UpdatedAt,
		PaidAt:      order.PaidAt,
		CancelledAt: order.CancelledAt,
	}
}

func parseAmount(raw string) float64 {
	if raw == "" {
		return 0
	}
	value, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		return 0
	}
	return value
}

func parseQuantity(raw json.Number) int {
	if raw == "" {
		return 0
	}
	value, err := raw.Int64()
	if err != nil {
		return 0
	}
	return int(value)
}
