package inventorycompare

import "testing"

func TestCuandoElCanalTieneMenosStockLaFilaSeActualiza(t *testing.T) {
	fila := BuildRow(Item{
		SKU: "BC112-L", ProbabilityQty: 12, HasStock: true,
		ChannelQty: 3, ChannelFound: true, ChannelManaged: true,
	})

	if fila.Action != ActionUpdate {
		t.Fatalf("se esperaba %q, llego %q", ActionUpdate, fila.Action)
	}
	if fila.Delta == nil || *fila.Delta != 9 {
		t.Fatalf("el delta deberia ser 9, llego %v", fila.Delta)
	}
}

func TestCuandoLosDosLadosCoincidenNoSeToca(t *testing.T) {
	fila := BuildRow(Item{
		SKU: "BC112-M", ProbabilityQty: 5, HasStock: true,
		ChannelQty: 5, ChannelFound: true, ChannelManaged: true,
	})

	if fila.Action != ActionUnchanged {
		t.Fatalf("se esperaba %q, llego %q", ActionUnchanged, fila.Action)
	}
	if fila.Delta == nil || *fila.Delta != 0 {
		t.Fatalf("el delta deberia ser 0, llego %v", fila.Delta)
	}
}

func TestSinRegistroDeInventarioSeCorrigeElCanalACero(t *testing.T) {
	fila := BuildRow(Item{SKU: "X", ChannelQty: 7, ChannelFound: true, ChannelManaged: true})

	if fila.Action != ActionUpdate {
		t.Fatalf("se esperaba %q, llego %q", ActionUpdate, fila.Action)
	}
	if fila.ProbabilityQty == nil || *fila.ProbabilityQty != 0 {
		t.Fatalf("sin registro se cuenta como cero, llego %v", fila.ProbabilityQty)
	}
	if fila.Reason != ReasonSinInventario {
		t.Fatalf("deberia avisar que no hay registro, llego %q", fila.Reason)
	}
}

func TestSinRegistroYCanalEnCeroNoEsUnCambio(t *testing.T) {
	fila := BuildRow(Item{SKU: "X", ChannelQty: 0, ChannelFound: true, ChannelManaged: true})

	if fila.Action != ActionUnchanged {
		t.Fatalf("se esperaba %q, llego %q", ActionUnchanged, fila.Action)
	}
}

func TestPublicacionQueNoManejaStockSeSalta(t *testing.T) {
	fila := BuildRow(Item{
		SKU: "Y", ProbabilityQty: 4, HasStock: true,
		ChannelQty: 0, ChannelFound: true, ChannelManaged: false,
	})

	if fila.Action != ActionSkip || fila.Reason != ReasonSinManejoStock {
		t.Fatalf("se esperaba saltar por manejo de stock, llego %q / %q", fila.Action, fila.Reason)
	}
}

func TestLaRazonPropiaDelCanalGanaSobreLasDemas(t *testing.T) {
	fila := BuildRow(Item{
		SKU: "Z", ProbabilityQty: 4, HasStock: true,
		ChannelFound: false, SkipReason: "publicacion de fulfillment",
	})

	if fila.Reason != "publicacion de fulfillment" {
		t.Fatalf("razon inesperada: %q", fila.Reason)
	}
}

func TestElResumenCuentaCadaAccion(t *testing.T) {
	filas := Build([]Item{
		{SKU: "b", ProbabilityQty: 1, HasStock: true, ChannelQty: 2, ChannelFound: true, ChannelManaged: true},
		{SKU: "a", ProbabilityQty: 3, HasStock: true, ChannelQty: 3, ChannelFound: true, ChannelManaged: true},
		{SKU: "c", SkipReason: "fulfillment"},
	})

	if filas[0].SKU != "a" {
		t.Fatalf("las filas deberian venir ordenadas por SKU, llego %q", filas[0].SKU)
	}

	totales := Summarize(filas)
	if totales.Total != 3 || totales.ToUpdate != 1 || totales.Unchanged != 1 || totales.Skipped != 1 {
		t.Fatalf("totales inesperados: %+v", totales)
	}
}
