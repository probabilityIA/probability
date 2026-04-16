package mappers

import (
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/invoicing/factus/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/factus/internal/infra/secondary/client/request"
)

// BuildCreateBillRequest construye el body del POST /v1/bills/validate desde el DTO de dominio
func BuildCreateBillRequest(req *dtos.CreateInvoiceRequest) request.CreateBillBody {
	cfg := req.Config

	items := make([]request.CreateBillItem, 0, len(req.Items))
	for _, item := range req.Items {
		taxRate := GetConfigString(cfg, "default_tax_rate", "19.00")
		if item.TaxRate != nil {
			taxRate = fmt.Sprintf("%.2f", *item.TaxRate)
		}

		codeRef := item.SKU
		if item.ProductID != nil && *item.ProductID != "" {
			codeRef = *item.ProductID
		}

		items = append(items, request.CreateBillItem{
			SchemeID:       GetConfigString(cfg, "item_scheme_id", "1"),
			CodeReference:  codeRef,
			Name:           item.Name,
			Quantity:       item.Quantity,
			Price:          item.UnitPrice,
			TaxRate:        taxRate,
			UnitMeasureID:  GetConfigInt(cfg, "unit_measure_id", 70),
			StandardCodeID: GetConfigInt(cfg, "standard_code_id", 1),
			IsExcluded:     0,
			TributeID:      GetConfigInt(cfg, "item_tribute_id", 1),
		})
	}

	refCode := GetConfigString(cfg, "reference_code", req.OrderID)
	if refCode == "" {
		refCode = req.OrderID
	}

	return request.CreateBillBody{
		NumberingRangeID:  GetConfigInt(cfg, "numbering_range_id", 0),
		ReferenceCode:     refCode,
		PaymentForm:       GetConfigString(cfg, "payment_form", "1"),
		PaymentMethodCode: GetConfigString(cfg, "payment_method_code", "10"),
		OperationType:     GetConfigInt(cfg, "operation_type", 10),
		SendEmail:         false,
		Document:          GetConfigString(cfg, "document", "01"),
		Customer: request.CreateBillCustomer{
			Identification:           req.Customer.DNI,
			DV:                       GetConfigString(cfg, "customer_dv", ""),
			Names:                    req.Customer.Name,
			Address:                  req.Customer.Address,
			Email:                    req.Customer.Email,
			Phone:                    req.Customer.Phone,
			LegalOrganizationID:      GetConfigString(cfg, "legal_organization_id", "2"),
			TributeID:                GetConfigString(cfg, "tribute_id", "21"),
			IdentificationDocumentID: GetConfigString(cfg, "identification_document_id", "3"),
			MunicipalityID:           GetConfigString(cfg, "municipality_id", ""),
		},
		Items: items,
	}
}
