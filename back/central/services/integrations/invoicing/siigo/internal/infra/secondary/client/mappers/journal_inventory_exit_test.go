package mappers

import (
	"testing"

	"github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/domain/dtos"
)

func inventoryExitConfig() map[string]interface{} {
	return map[string]interface{}{
		"inventory_exit_account_code":        "14350501",
		"inventory_exit_offset_account_code": "61350501",
		"inventory_exit_document_id":         float64(9911),
		"inventory_exit_warehouse_id":        float64(7),
	}
}

func inventoryExitRequest() *dtos.CreateJournalRequest {
	return &dtos.CreateJournalRequest{
		Currency: "COP",
		Date:     "2026-08-18",
		Items: []dtos.JournalItemData{
			{SKU: "CAM-001", Name: "Camiseta", Quantity: 2, TotalPrice: 60000},
			{SKU: "PAN-002", Name: "Pantaloneta", Quantity: 1, TotalPrice: 40000},
		},
	}
}

func TestSalidaDeInventarioGeneraDosLineasPorItem(t *testing.T) {
	req := inventoryExitRequest()
	req.Config = inventoryExitConfig()

	journal := BuildCreateJournalRequest(req)

	if len(journal.Items) != 4 {
		t.Fatalf("lineas = %d, se esperaban 4 (una de inventario y una de contrapartida por cada item)", len(journal.Items))
	}
	if journal.Document.ID != 9911 {
		t.Errorf("document id = %d, se esperaba el inventory_exit_document_id 9911", journal.Document.ID)
	}
}

func TestSalidaDeInventarioQuedaBalanceada(t *testing.T) {
	req := inventoryExitRequest()
	req.Config = inventoryExitConfig()

	journal := BuildCreateJournalRequest(req)

	var debit, credit float64
	for _, item := range journal.Items {
		switch item.Account.Movement {
		case "Debit":
			debit += item.Value
		case "Credit":
			credit += item.Value
		default:
			t.Fatalf("movimiento inesperado: %q", item.Account.Movement)
		}
	}

	if debit != credit {
		t.Errorf("debito %.2f != credito %.2f: Siigo rechaza un comprobante descuadrado", debit, credit)
	}
	if credit != 100000 {
		t.Errorf("credito = %.2f, se esperaba 100000 (60000 + 40000)", credit)
	}
}

func TestSoloLaLineaDeInventarioLlevaProductoYBodega(t *testing.T) {
	req := inventoryExitRequest()
	req.Config = inventoryExitConfig()

	journal := BuildCreateJournalRequest(req)

	conProducto := 0
	for _, item := range journal.Items {
		if item.Product == nil {
			if item.Account.Movement != "Debit" {
				t.Errorf("la linea sin producto deberia ser el debito de contrapartida, es %q", item.Account.Movement)
			}
			continue
		}
		conProducto++
		if item.Account.Movement != "Credit" {
			t.Errorf("la linea con producto debe acreditar inventario, movimiento = %q", item.Account.Movement)
		}
		if item.Account.Code != "14350501" {
			t.Errorf("cuenta de inventario = %q, se esperaba 14350501", item.Account.Code)
		}
		if item.Product.Warehouse != 7 {
			t.Errorf("bodega = %d, se esperaba 7: sin bodega Siigo no descarga el kardex", item.Product.Warehouse)
		}
		if item.Product.Quantity == 0 {
			t.Error("cantidad en cero: no descargaria inventario")
		}
	}

	if conProducto != 2 {
		t.Errorf("lineas con producto = %d, se esperaban 2", conProducto)
	}
}

func TestSinCuentasConfiguradasUsaElMapeoNormal(t *testing.T) {
	req := inventoryExitRequest()
	req.Config = map[string]interface{}{
		"journal_document_id":  float64(100),
		"default_account_code": "13050501",
	}

	journal := BuildCreateJournalRequest(req)

	if len(journal.Items) != 2 {
		t.Fatalf("lineas = %d, se esperaban 2: sin cuentas de salida no se duplica", len(journal.Items))
	}
	if journal.Document.ID != 100 {
		t.Errorf("document id = %d, se esperaba 100", journal.Document.ID)
	}
}
