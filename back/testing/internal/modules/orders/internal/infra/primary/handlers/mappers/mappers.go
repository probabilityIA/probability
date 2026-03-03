package mappers

import (
	"github.com/secamc93/probability/back/testing/internal/modules/orders/internal/domain/entities"
	"github.com/secamc93/probability/back/testing/internal/modules/orders/internal/infra/primary/handlers/response"
)

func ReferenceDataToResponse(data *entities.ReferenceData) *response.ReferenceData {
	resp := &response.ReferenceData{
		Products:       make([]response.Product, len(data.Products)),
		Integrations:   make([]response.Integration, len(data.Integrations)),
		PaymentMethods: make([]response.PaymentMethod, len(data.PaymentMethods)),
		OrderStatuses:  make([]response.OrderStatus, len(data.OrderStatuses)),
	}

	for i, p := range data.Products {
		resp.Products[i] = response.Product{
			ID:       p.ID,
			Name:     p.Name,
			SKU:      p.SKU,
			Price:    p.Price,
			Currency: p.Currency,
		}
	}

	for i, ig := range data.Integrations {
		resp.Integrations[i] = response.Integration{
			ID:                ig.ID,
			Name:              ig.Name,
			Code:              ig.Code,
			Category:          ig.Category,
			IntegrationTypeID: ig.IntegrationTypeID,
		}
	}

	for i, pm := range data.PaymentMethods {
		resp.PaymentMethods[i] = response.PaymentMethod{
			ID:   pm.ID,
			Code: pm.Code,
			Name: pm.Name,
		}
	}

	for i, os := range data.OrderStatuses {
		resp.OrderStatuses[i] = response.OrderStatus{
			ID:   os.ID,
			Code: os.Code,
			Name: os.Name,
		}
	}

	return resp
}

func GenerateResultToResponse(result *entities.GenerateResult) *response.GenerateResult {
	resp := &response.GenerateResult{
		Total:   result.Total,
		Created: result.Created,
		Failed:  result.Failed,
	}

	if len(result.Orders) > 0 {
		resp.Orders = make([]response.CreatedOrder, len(result.Orders))
		for i, o := range result.Orders {
			resp.Orders[i] = response.CreatedOrder{
				ID:           o.ID,
				OrderNumber:  o.OrderNumber,
				Total:        o.Total,
				CustomerName: o.CustomerName,
			}
		}
	}

	if len(result.Errors) > 0 {
		resp.Errors = make([]response.OrderError, len(result.Errors))
		for i, e := range result.Errors {
			resp.Errors[i] = response.OrderError{
				Index:   e.Index,
				Message: e.Message,
			}
		}
	}

	if len(result.APILogs) > 0 {
		resp.APILogs = make([]response.APICallLog, len(result.APILogs))
		for i, l := range result.APILogs {
			resp.APILogs[i] = response.APICallLog{
				Index:      l.Index,
				Success:    l.Success,
				Timestamp:  l.Timestamp,
				DurationMs: l.DurationMs,
				Request: response.APIRequest{
					Method: l.Request.Method,
					URL:    l.Request.URL,
					Body:   l.Request.Body,
				},
				Response: response.APIResponse{
					StatusCode: l.Response.StatusCode,
					Body:       l.Response.Body,
				},
			}
		}
	}

	return resp
}
