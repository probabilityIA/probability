package orderscompare

import (
	"testing"
	"time"
)

func TestBuildClasificaOrdenes(t *testing.T) {
	now := time.Now()
	channel := []ChannelOrder{
		{ExternalID: "1001", Number: "#1001", Status: "pending", Total: 50000, CreatedAt: now},
		{ExternalID: "1002", Number: "#1002", Status: "completed", Total: 80000, CreatedAt: now.Add(-time.Hour)},
		{ExternalID: "1003", Number: "#1003", Status: "paid", Total: 30000, CreatedAt: now.Add(-2 * time.Hour)},
	}
	locals := []LocalOrder{
		{OrderID: "uuid-3", OrderNumber: "#1003", ExternalID: "1003", Status: "paid", Total: 30000},
		{OrderID: "uuid-9", OrderNumber: "#9999", ExternalID: "9999", Status: "pending", Total: 1000},
	}

	result := Build(channel, locals)

	if result.Totals.ToCreate != 2 {
		t.Fatalf("esperaba 2 por crear, obtuve %d", result.Totals.ToCreate)
	}
	if result.Totals.InSync != 1 {
		t.Fatalf("esperaba 1 en sincronia, obtuve %d", result.Totals.InSync)
	}
	if result.Totals.OnlyInProbability != 1 {
		t.Fatalf("esperaba 1 solo en Probability, obtuve %d", result.Totals.OnlyInProbability)
	}
	if result.Totals.WithoutInventory != 1 {
		t.Fatalf("esperaba 1 sin movimiento de inventario, obtuve %d", result.Totals.WithoutInventory)
	}
	if result.Totals.Total != 4 {
		t.Fatalf("esperaba 4 filas, obtuve %d", result.Totals.Total)
	}
}

func TestBuildMarcaDiferenciaDeEstadoYTotal(t *testing.T) {
	channel := []ChannelOrder{{ExternalID: "A1", Status: "paid", Total: 100}}
	locals := []LocalOrder{{OrderID: "u1", ExternalID: "a1", Status: "pending", Total: 90}}

	result := Build(channel, locals)

	row := result.Rows[0]
	if row.Action != ActionInSync {
		t.Fatalf("esperaba in_sync, obtuve %s", row.Action)
	}
	if !row.StatusMismatch || !row.TotalMismatch {
		t.Fatalf("esperaba diferencia de estado y total, obtuve %+v", row)
	}
}

func TestSkipsInventory(t *testing.T) {
	casos := map[string]bool{
		"pending":   false,
		"paid":      false,
		"completed": true,
		"delivered": true,
		"shipped":   true,
		"cancelled": true,
		"refunded":  true,
		"abandoned": true,
		"":          false,
	}

	for status, esperado := range casos {
		skips, reason := SkipsInventory(status)
		if skips != esperado {
			t.Fatalf("estado %q: esperaba skip=%v, obtuve %v", status, esperado, skips)
		}
		if skips && reason == "" {
			t.Fatalf("estado %q: se salta el inventario pero no explica por que", status)
		}
	}
}

func TestPaginate(t *testing.T) {
	rows := make([]Row, 25)
	page, total := Paginate(rows, 3, 10)
	if total != 25 {
		t.Fatalf("esperaba total 25, obtuve %d", total)
	}
	if len(page) != 5 {
		t.Fatalf("esperaba 5 filas en la pagina 3, obtuve %d", len(page))
	}

	vacia, _ := Paginate(rows, 9, 10)
	if len(vacia) != 0 {
		t.Fatalf("esperaba pagina vacia, obtuve %d filas", len(vacia))
	}
}

func TestSkipsInventoryForUsaElEstadoDeEntrega(t *testing.T) {
	skips, reason := SkipsInventoryFor("paid", "delivered")
	if !skips {
		t.Fatal("una orden pagada pero ya entregada no puede reservar stock")
	}
	if reason == "" {
		t.Fatal("falta el motivo")
	}

	if skips, _ := SkipsInventoryFor("paid", ""); skips {
		t.Fatal("una orden pagada sin entregar si debe reservar stock")
	}
}

func TestBuildMarcaSinInventarioPorEntrega(t *testing.T) {
	channel := []ChannelOrder{{ExternalID: "P1", Status: "paid", FulfillmentStatus: "delivered", Total: 100}}

	result := Build(channel, nil)

	if result.Totals.WithoutInventory != 1 {
		t.Fatalf("esperaba 1 sin inventario, obtuve %d", result.Totals.WithoutInventory)
	}
	if result.Rows[0].MovesInventory {
		t.Fatal("la fila entregada no deberia mover inventario")
	}
}
