package mappers

import (
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/infra/secondary/client/request"
)

func BuildCreateInvoiceRequest(req *dtos.CreateInvoiceRequest, customerID string) request.SiigoInvoice {
	config := req.Config

	documentID := getIntFromConfig(config, "document_id", 0)
	paymentMethodID := getIntFromConfig(config, "payment_method_id", 0)
	taxID := getIntFromConfig(config, "tax_id", 0)
	sellerID := getIntFromConfig(config, "seller_id", 0)
	costCenterID := getIntFromConfig(config, "cost_center_id", 0)
	stampSend := getBoolFromConfig(config, "stamp_send", true)
	mailSend := getBoolFromConfig(config, "mail_send", false)

	items := make([]request.SiigoItem, 0, len(req.Items))
	for _, item := range req.Items {
		siigoItem := request.SiigoItem{
			Code:        item.SKU,
			Description: item.Name,
			Quantity:    float64(item.Quantity),
			Price:       item.UnitPrice,
			Discount:    item.Discount,
		}

		if taxID > 0 && item.Tax > 0 {
			siigoItem.Taxes = []request.SiigoTax{{ID: taxID}}
		}

		items = append(items, siigoItem)
	}

	customerIdentification := req.Customer.DNI
	if customerIdentification == "" {
		customerIdentification = customerID
	}

	currency := req.Currency
	if currency == "" {
		currency = "COP"
	}

	var payments []request.SiigoPayment
	if paymentMethodID > 0 {
		payments = []request.SiigoPayment{
			{
				ID:    paymentMethodID,
				Value: req.Total,
			},
		}
	}

	invoice := request.SiigoInvoice{
		Document: request.SiigoDocument{ID: documentID},
		Date:     time.Now().Format("2006-01-02"),
		Customer: request.SiigoCustomerRef{
			Identification: customerIdentification,
			BranchOffice:   0,
		},
		Seller:   sellerID,
		Items:    items,
		Payments: payments,
		Stamp:    &request.SiigoStamp{Send: stampSend},
		Mail:     &request.SiigoMail{Send: mailSend},
	}

	if costCenterID > 0 {
		invoice.CostCenter = costCenterID
	}

	if currency != "COP" {
		invoice.Currency = &request.SiigoCurrency{Code: currency}
	}

	return invoice
}

func getStringFromConfig(config map[string]interface{}, key, defaultVal string) string {
	if v, ok := config[key]; ok {
		if s, ok := v.(string); ok && s != "" {
			return s
		}
	}
	return defaultVal
}

func getIntFromConfig(config map[string]interface{}, key string, defaultVal int) int {
	if v, ok := config[key]; ok {
		switch val := v.(type) {
		case int:
			return val
		case float64:
			return int(val)
		case int64:
			return int(val)
		}
	}
	return defaultVal
}

func getBoolFromConfig(config map[string]interface{}, key string, defaultVal bool) bool {
	if v, ok := config[key]; ok {
		if b, ok := v.(bool); ok {
			return b
		}
	}
	return defaultVal
}
