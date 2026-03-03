package mappers

import (
	"github.com/secamc93/probability/back/testing/modules/orders/internal/domain/entities"
	"github.com/secamc93/probability/back/testing/modules/orders/internal/infra/primary/handlers/response"
)

func ReferenceDataToResponse(data *entities.ReferenceData) *response.ReferenceData {
	resp := &response.ReferenceData{
		Products:       make([]response.Product, len(data.Products)),
		Integrations:   make([]response.Integration, len(data.Integrations)),
		PaymentMethods: make([]response.PaymentMethod, len(data.PaymentMethods)),
		OrderStatuses:  make([]response.OrderStatus, len(data.OrderStatuses)),
		WebhookTopics:  data.WebhookTopics,
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
			ID:                  ig.ID,
			Name:                ig.Name,
			Code:                ig.Code,
			Category:            ig.Category,
			CategoryID:          ig.CategoryID,
			IntegrationTypeID:   ig.IntegrationTypeID,
			IntegrationTypeCode: ig.IntegrationTypeCode,
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
		Total: result.Total,
	}

	if len(result.Payloads) > 0 {
		resp.Payloads = make([]response.WebhookPayload, len(result.Payloads))
		for i, p := range result.Payloads {
			resp.Payloads[i] = response.WebhookPayload{
				URL:     p.URL,
				Method:  p.Method,
				Headers: p.Headers,
				Body:    p.Body,
				RawBody: p.RawBody,
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

	return resp
}
