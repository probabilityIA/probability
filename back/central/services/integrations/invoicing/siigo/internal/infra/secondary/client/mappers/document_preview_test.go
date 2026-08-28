package mappers

import (
	"encoding/json"
	"testing"

	"github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/infra/secondary/client/response"
)

const respuestaRealSiigo = `{
  "id": "0fc3b4fc-08eb-4d13-8b45-8549bb0ca1b9",
  "date": "2026-08-27",
  "mail": {"status": "not_sent", "observations": "The invoice has not been sent by mail"},
  "name": "FV-2-3",
  "items": [
    {"id": "64352a99", "code": "BN14-5XL", "price": 49900, "total": 49900, "quantity": 1, "description": "Buso Cuello Redondo"},
    {"id": "5a0f46e0", "code": "Envio", "price": 20000, "total": 20000, "quantity": 1, "description": "Envio"}
  ],
  "stamp": {"status": "Sending"},
  "total": 69900,
  "number": 3,
  "prefix": "FEVD",
  "seller": 872,
  "balance": 0,
  "customer": {"id": "e3486216", "branch_office": 0, "identification": "1000641221"},
  "document": {"id": 30606},
  "metadata": {"created": "2026-08-28T00:20:58Z"},
  "payments": [{"id": 12401, "name": "Transferencia", "value": 69900}],
  "public_url": "https://documentview.siigo.com/document?data=abc",
  "observations": "order:5333b852 | #VIG-0117"
}`

func TestDocumentoDeFactura(t *testing.T) {
	var r response.CreateInvoiceResponse
	if err := json.Unmarshal([]byte(respuestaRealSiigo), &r); err != nil {
		t.Fatalf("no se pudo parsear la respuesta de Siigo: %v", err)
	}

	doc := DocumentoDeFactura(r, []byte(respuestaRealSiigo))

	if doc["document_number"] != "FV-2-3" {
		t.Errorf("document_number = %v", doc["document_number"])
	}
	if doc["document_prefix"] != "FEVD" {
		t.Errorf("document_prefix = %v", doc["document_prefix"])
	}
	if doc["customer_identification"] != "1000641221" {
		t.Errorf("customer_identification = %v", doc["customer_identification"])
	}
	if doc["total"] != 69900.0 {
		t.Errorf("total = %v", doc["total"])
	}
	if doc["stamp_status"] != "Sending" {
		t.Errorf("stamp_status = %v", doc["stamp_status"])
	}
	if doc["electronic"] != true {
		t.Errorf("electronic = %v", doc["electronic"])
	}
	if doc["public_url"] != "https://documentview.siigo.com/document?data=abc" {
		t.Errorf("public_url = %v", doc["public_url"])
	}

	items, ok := doc["items"].([]map[string]interface{})
	if !ok || len(items) != 2 {
		t.Fatalf("items = %v", doc["items"])
	}
	if items[0]["code"] != "BN14-5XL" || items[0]["total"] != 49900.0 {
		t.Errorf("primer item = %v", items[0])
	}

	pagos, ok := doc["payments"].([]map[string]interface{})
	if !ok || len(pagos) != 1 || pagos[0]["name"] != "Transferencia" {
		t.Errorf("payments = %v", doc["payments"])
	}

	if doc["raw"] == nil {
		t.Error("falta el json crudo de Siigo en raw")
	}
}

func TestDocumentoDeFacturaSinCuerpoCrudo(t *testing.T) {
	doc := DocumentoDeFactura(response.CreateInvoiceResponse{Name: "FV-1-1"}, nil)
	if doc["raw"] != nil {
		t.Error("raw deberia quedar ausente sin cuerpo")
	}
	if doc["preview_version"] != 1 {
		t.Errorf("preview_version = %v", doc["preview_version"])
	}
}
