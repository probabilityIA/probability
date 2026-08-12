package productmatch

import "testing"

func canal(sku string, qty int) TypoCandidate {
	return TypoCandidate{SKU: sku, Quantity: qty, HasQty: true}
}

func TestDetectaLetrasTrocadas(t *testing.T) {
	out := DetectTypos(
		[]TypoCandidate{canal("SDH22-S", 35)},
		[]TypoCandidate{canal("SHD22-S", 36)},
	)
	if len(out) != 1 {
		t.Fatalf("esperaba 1 sospechoso, obtuve %d", len(out))
	}
	if out[0].ProbabilitySKU != "SHD22-S" || out[0].Reason != ReasonTransposed {
		t.Fatalf("obtuve %+v", out[0])
	}
}

func TestLasLetrasTrocadasNoNecesitanQueLaCantidadCoincida(t *testing.T) {
	out := DetectTypos(
		[]TypoCandidate{canal("SDH22-S", 0)},
		[]TypoCandidate{canal("SHD22-S", 99)},
	)
	if len(out) != 1 {
		t.Fatal("la transposicion es senal suficiente por si sola")
	}
}

func TestUnCaracterSoloNoAlcanza(t *testing.T) {
	out := DetectTypos(
		[]TypoCandidate{canal("BT11-3XL", 5)},
		[]TypoCandidate{canal("BT11-XL", 40)},
	)
	if len(out) != 0 {
		t.Fatalf("con SKUs cortos casi todo queda a un caracter: sin cantidad que confirme no se sugiere, obtuve %+v", out)
	}
}

func TestUnCaracterConLaCantidadIgualSiSeReporta(t *testing.T) {
	out := DetectTypos(
		[]TypoCandidate{canal("BT11-3XL", 5)},
		[]TypoCandidate{canal("BT11-XL", 5)},
	)
	if len(out) != 1 || out[0].Reason != ReasonOneCharAndQty {
		t.Fatalf("obtuve %+v", out)
	}
}

func TestUnaUnidadDeDiferenciaSigueContandoComoIgual(t *testing.T) {
	out := DetectTypos(
		[]TypoCandidate{canal("BT11-3XL", 5)},
		[]TypoCandidate{canal("BT11-XL", 6)},
	)
	if len(out) != 1 {
		t.Fatal("una venta en curso explica una unidad de diferencia")
	}
}

func TestDosUnidadesDeDiferenciaYaNoConfirman(t *testing.T) {
	out := DetectTypos(
		[]TypoCandidate{canal("BT11-3XL", 5)},
		[]TypoCandidate{canal("BT11-XL", 8)},
	)
	if len(out) != 0 {
		t.Fatalf("obtuve %+v", out)
	}
}

func TestSinCantidadConocidaNoSeUsaEsaSenal(t *testing.T) {
	out := DetectTypos(
		[]TypoCandidate{{SKU: "BT11-3XL"}},
		[]TypoCandidate{{SKU: "BT11-XL"}},
	)
	if len(out) != 0 {
		t.Fatal("sin cantidad de ninguno de los dos lados no hay segunda senal")
	}
}

func TestIgnoraDiferenciasDeEspaciosYMayusculas(t *testing.T) {
	out := DetectTypos(
		[]TypoCandidate{canal(" bh313 - 2xl ", 5)},
		[]TypoCandidate{canal("BH313-2XL", 5)},
	)
	if len(out) != 0 {
		t.Fatal("eso no es un error de digitacion sino un tema de formato: lo resuelve la normalizacion")
	}
}

func TestLaTransposicionGanaSobreOtroCandidato(t *testing.T) {
	out := DetectTypos(
		[]TypoCandidate{canal("SDH22-S", 5)},
		[]TypoCandidate{canal("SDH2-S", 5), canal("SHD22-S", 99)},
	)
	if len(out) != 1 {
		t.Fatalf("esperaba 1, obtuve %d", len(out))
	}
	if out[0].ProbabilitySKU != "SHD22-S" {
		t.Fatalf("la transposicion es mas confiable que un caracter, obtuve %q", out[0].ProbabilitySKU)
	}
}

func TestCadaSKUDelCanalSeReportaUnaSolaVez(t *testing.T) {
	out := DetectTypos(
		[]TypoCandidate{canal("BT11-3XL", 5)},
		[]TypoCandidate{canal("BT11-XL", 5), canal("BT11-4XL", 5)},
	)
	if len(out) != 1 {
		t.Fatalf("un SKU no puede tener dos sugerencias, obtuve %d", len(out))
	}
}

func TestNoSeSugiereASiMismo(t *testing.T) {
	out := DetectTypos(
		[]TypoCandidate{canal("SHD22-S", 5)},
		[]TypoCandidate{canal("shd22-s", 5)},
	)
	if len(out) != 0 {
		t.Fatal("es el mismo SKU, no un error")
	}
}

func TestListasVaciasNoRevientan(t *testing.T) {
	if out := DetectTypos(nil, []TypoCandidate{canal("A", 1)}); len(out) != 0 {
		t.Fatal("sin canal no hay nada que sugerir")
	}
	if out := DetectTypos([]TypoCandidate{canal("A", 1)}, nil); len(out) != 0 {
		t.Fatal("sin productos no hay contra que comparar")
	}
	if out := DetectTypos(nil, nil); out == nil {
		t.Fatal("debe devolver lista vacia, no nil")
	}
}

func TestIgnoraSKUVacio(t *testing.T) {
	out := DetectTypos(
		[]TypoCandidate{{SKU: "   ", HasQty: true}},
		[]TypoCandidate{canal("A-1", 1)},
	)
	if len(out) != 0 {
		t.Fatal("un SKU vacio no es un error de digitacion")
	}
}

func TestSKUsMuyDistintosNoSeSugieren(t *testing.T) {
	out := DetectTypos(
		[]TypoCandidate{canal("SDH22-S", 5)},
		[]TypoCandidate{canal("CAMISETA-NEGRA-XL", 5)},
	)
	if len(out) != 0 {
		t.Fatalf("obtuve %+v", out)
	}
}

func TestConservaLaPublicacionYLaVarianteParaMostrarlas(t *testing.T) {
	out := DetectTypos(
		[]TypoCandidate{{SKU: "SDH22-S", Ref: "MCO123", Label: "Verde Militar / S", Name: "Shorts", Quantity: 35, HasQty: true}},
		[]TypoCandidate{canal("SHD22-S", 36)},
	)
	if len(out) != 1 {
		t.Fatalf("esperaba 1, obtuve %d", len(out))
	}
	s := out[0]
	if s.ChannelRef != "MCO123" || s.ChannelLabel != "Verde Militar / S" || s.ChannelName != "Shorts" {
		t.Fatalf("debia conservar el contexto para mostrarlo: %+v", s)
	}
	if s.ChannelQty != 35 || s.ProbabilityQty != 36 {
		t.Fatalf("las cantidades son lo que le permite al negocio confirmarlo: %+v", s)
	}
}

func TestTransposicionRequiereCaracteresContiguos(t *testing.T) {
	if isTransposition("ABCD", "DBCA") {
		t.Fatal("intercambiar el primero con el ultimo no es un error de digitacion tipico")
	}
	if !isTransposition("ABCD", "BACD") {
		t.Fatal("dos contiguos intercambiados si lo es")
	}
}

func TestDistanciaDeEdicionCortaAlSuperarElTope(t *testing.T) {
	if d := editDistance("ABCDEFGH", "ZZZZZZZZ", 1); d <= 1 {
		t.Fatalf("obtuve %d", d)
	}
	if d := editDistance("ABC", "ABC", 1); d != 0 {
		t.Fatalf("obtuve %d", d)
	}
	if d := editDistance("ABC", "ABD", 1); d != 1 {
		t.Fatalf("obtuve %d", d)
	}
}

func TestTransposicionDescartaLargosDistintos(t *testing.T) {
	if isTransposition("ABC", "ABCD") {
		t.Fatal("largos distintos no pueden ser una transposicion")
	}
	if isTransposition("ABC", "ABC") {
		t.Fatal("identicos no son un error")
	}
	if isTransposition("ABCD", "AXYD") {
		t.Fatal("dos cambios contiguos que no son intercambio no cuentan")
	}
	if isTransposition("ABCDE", "BACED") {
		t.Fatal("mas de un par intercambiado no es un error de digitacion simple")
	}
}

func TestDosCerosNoConfirmanNada(t *testing.T) {
	out := DetectTypos(
		[]TypoCandidate{canal("TD212-3XL", 0)},
		[]TypoCandidate{canal("TD212-6XL", 0)},
	)
	if len(out) != 0 {
		t.Fatalf("cero contra cero es el valor por defecto de todo lo que nadie toco, no evidencia: %+v", out)
	}
}

func TestUnCeroContraStockRealTampocoConfirma(t *testing.T) {
	out := DetectTypos(
		[]TypoCandidate{canal("TD212-3XL", 0)},
		[]TypoCandidate{canal("TD212-6XL", 9)},
	)
	if len(out) != 0 {
		t.Fatalf("obtuve %+v", out)
	}
}
