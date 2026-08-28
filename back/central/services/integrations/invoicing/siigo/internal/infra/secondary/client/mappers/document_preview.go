package mappers

import (
	"encoding/json"

	"github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/infra/secondary/client/response"
)

const previewVersion = 1

func DocumentoDeFactura(r response.CreateInvoiceResponse, rawBody []byte) map[string]interface{} {
	items := make([]map[string]interface{}, 0, len(r.Items))
	impuestoTotal := 0.0
	descuentoTotal := 0.0

	for _, it := range r.Items {
		tasa := 0.0
		valorImpuesto := 0.0
		for _, tx := range it.Taxes {
			valorImpuesto += tx.Value
			if tx.Percentage > tasa {
				tasa = tx.Percentage
			}
		}
		impuestoTotal += valorImpuesto
		descuentoTotal += it.Discount

		items = append(items, map[string]interface{}{
			"code":        it.Code,
			"description": it.Description,
			"quantity":    it.Quantity,
			"unit_price":  it.Price,
			"discount":    it.Discount,
			"tax_rate":    tasa,
			"tax_value":   valorImpuesto,
			"total":       it.Total,
		})
	}

	pagos := make([]map[string]interface{}, 0, len(r.Payments))
	for _, p := range r.Payments {
		pagos = append(pagos, map[string]interface{}{
			"name":  p.Name,
			"value": p.Value,
		})
	}

	if r.TotalTax > 0 {
		impuestoTotal = r.TotalTax
	}
	if r.Discount > 0 {
		descuentoTotal = r.Discount
	}

	cufe := r.Stamp.CUFE
	if cufe == "" {
		cufe = r.Metadata.CUFE
	}

	doc := map[string]interface{}{
		"preview_version":         previewVersion,
		"provider":                "siigo",
		"document_number":         r.Name,
		"document_prefix":         r.Prefix,
		"document_type_id":        r.Document.ID,
		"document_date":           r.Date,
		"customer_name":           r.Customer.Name,
		"customer_identification": r.Customer.Identification,
		"total":                   r.Total,
		"tax":                     impuestoTotal,
		"discount":                descuentoTotal,
		"balance":                 r.Balance,
		"items":                   items,
		"payments":                pagos,
		"notes":                   r.Observations,
		"external_id":             r.ID,
	}

	if r.PublicURL != "" {
		doc["public_url"] = r.PublicURL
	}
	if cufe != "" {
		doc["cufe"] = cufe
		doc["electronic"] = true
	}
	if r.Stamp.Status != "" {
		doc["stamp_status"] = r.Stamp.Status
		doc["electronic"] = true
	}
	if r.Stamp.Observations != "" {
		doc["stamp_observations"] = r.Stamp.Observations
	}
	if r.Mail.Status != "" {
		doc["mail_status"] = r.Mail.Status
	}

	if crudo := crudoComoMapa(rawBody); crudo != nil {
		doc["raw"] = crudo
	}

	return doc
}

func crudoComoMapa(rawBody []byte) map[string]interface{} {
	if len(rawBody) == 0 {
		return nil
	}
	var m map[string]interface{}
	if err := json.Unmarshal(rawBody, &m); err != nil {
		return nil
	}
	return m
}
