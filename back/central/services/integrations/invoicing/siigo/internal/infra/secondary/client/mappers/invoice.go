package mappers

import (
	"strings"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/infra/secondary/client/request"
)

// BuildCreateInvoiceRequest construye el body de la request de Siigo para crear una factura
func BuildCreateInvoiceRequest(req *dtos.CreateInvoiceRequest, customerID string) request.SiigoInvoice {
	config := req.Config

	// Document ID (tipo de documento FV desde la config de la integración)
	documentID := getIntFromConfig(config, "document_id", 0)

	// Payment method ID
	paymentMethodID := getIntFromConfig(config, "payment_method_id", 0)

	// Tax ID (IVA)
	taxID := getIntFromConfig(config, "tax_id", 0)

	// ID type del cliente (código de tipo de documento)
	customerIDType := getStringFromConfig(config, "customer_id_type", "13") // 13 = CC por defecto

	// Person type
	personType := getStringFromConfig(config, "person_type", "Person")

	// Construir items
	items := make([]request.SiigoItem, 0, len(req.Items))
	for _, item := range req.Items {
		siigoItem := request.SiigoItem{
			Code:        request.SiigoProductCode{Code: item.SKU},
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

	// Construir nombre del cliente
	customerName := splitCustomerName(req.Customer.Name)

	// Construir contactos (email)
	var contacts []request.SiigoContact
	if req.Customer.Email != "" {
		contacts = []request.SiigoContact{
			{
				FirstName:             customerName[0],
				Email:                 req.Customer.Email,
				SendElectronicInvoice: true,
			},
		}
	}

	// Construir teléfonos
	var phones []request.SiigoPhone
	if req.Customer.Phone != "" {
		phones = []request.SiigoPhone{
			{Number: req.Customer.Phone},
		}
	}

	// Dirección
	var address *request.SiigoAddress
	if req.Customer.Address != "" {
		address = &request.SiigoAddress{Address: req.Customer.Address}
	}

	// Identificación del cliente
	customerIdentification := req.Customer.DNI
	if customerIdentification == "" {
		customerIdentification = customerID
	}

	// Moneda
	currency := req.Currency
	if currency == "" {
		currency = "COP"
	}

	// Pagos
	var payments []request.SiigoPayment
	if paymentMethodID > 0 {
		payments = []request.SiigoPayment{
			{
				ID:    paymentMethodID,
				Value: req.Total,
			},
		}
	}

	return request.SiigoInvoice{
		Document: request.SiigoDocument{ID: documentID},
		Date:     time.Now().Format("2006-01-02"),
		Currency: request.SiigoCurrency{Code: currency},
		Customer: request.SiigoCustomerRef{
			PersonType:     personType,
			IDType:         request.SiigoIDType{Code: customerIDType},
			Identification: customerIdentification,
			Name:           customerName,
			Address:        address,
			Phones:         phones,
			Contacts:       contacts,
		},
		Items:    items,
		Payments: payments,
	}
}

// splitCustomerName divide el nombre completo en [first, last] o [company]
func splitCustomerName(fullName string) []string {
	parts := strings.Fields(fullName)
	if len(parts) == 0 {
		return []string{"Sin Nombre"}
	}
	if len(parts) == 1 {
		return []string{parts[0]}
	}
	// first_name = primera palabra, last_name = resto
	firstName := parts[0]
	lastName := strings.Join(parts[1:], " ")
	return []string{firstName, lastName}
}

// getStringFromConfig obtiene un string de la config o retorna el default
func getStringFromConfig(config map[string]interface{}, key, defaultVal string) string {
	if v, ok := config[key]; ok {
		if s, ok := v.(string); ok && s != "" {
			return s
		}
	}
	return defaultVal
}

// getIntFromConfig obtiene un int de la config o retorna el default
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
