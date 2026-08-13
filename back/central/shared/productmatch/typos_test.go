package productmatch

import "testing"

func canal(sku string, qty int) TypoCandidate {
	return TypoCandidate{SKU: sku, Quantity: qty, HasQty: true}
}

func propio(sku string, qty int) TypoCandidate {
	return TypoCandidate{SKU: sku, Quantity: qty, HasQty: true}
}

func unico(t *testing.T, out []TypoSuspect) TypoSuspect {
	t.Helper()
	if len(out) != 1 {
		t.Fatalf("esperaba 1 sugerencia, llegaron %d: %+v", len(out), out)
	}
	return out[0]
}

func TestDetectTyposLetrasTrocadas(t *testing.T) {
	out := DetectTypos(
		[]TypoCandidate{canal("SDH22-S", 35)},
		[]TypoCandidate{propio("SHD22-S", 36)},
		TypoOptions{},
	)
	s := unico(t, out)
	if s.Pattern != PatternTransposed {
		t.Fatalf("patron = %q, esperaba %q", s.Pattern, PatternTransposed)
	}
	if s.ProbabilitySKU != "SHD22-S" {
		t.Fatalf("sugirio %q", s.ProbabilitySKU)
	}
	if s.FixSide != FixInChannel {
		t.Fatalf("lado = %q, esperaba corregir en el canal", s.FixSide)
	}
}

func TestDetectTyposEspacioSobrante(t *testing.T) {
	out := DetectTypos(
		[]TypoCandidate{canal("BN14 -5XL", 10)},
		[]TypoCandidate{propio("BN14-5XL", 10), propio("BN124-5XL", 10)},
		TypoOptions{},
	)
	s := unico(t, out)
	if s.Pattern != PatternSpacing {
		t.Fatalf("patron = %q, esperaba %q", s.Pattern, PatternSpacing)
	}
	if s.ProbabilitySKU != "BN14-5XL" {
		t.Fatalf("sugirio %q, el espacio gana sobre el parecido", s.ProbabilitySKU)
	}
	if s.FixSide != FixInChannel {
		t.Fatalf("lado = %q, el espacio esta en el canal", s.FixSide)
	}
}

func TestDetectTyposEspacioDelLadoDeProbability(t *testing.T) {
	out := DetectTypos(
		[]TypoCandidate{canal("BN14-5XL", 10)},
		[]TypoCandidate{propio("BN14 -5XL", 10)},
		TypoOptions{},
	)
	s := unico(t, out)
	if s.FixSide != FixInProbability {
		t.Fatalf("lado = %q, el espacio esta en Probability", s.FixSide)
	}
}

func TestDetectTyposNuncaCambiaLaTalla(t *testing.T) {
	casos := [][2]string{
		{"TD22-5XL", "TD22-6XL"},
		{"TD250-6XL", "TD250-5XL"},
		{"TD21-2XL", "TD212-XL"},
	}
	for _, c := range casos {
		out := DetectTypos(
			[]TypoCandidate{canal(c[0], 14)},
			[]TypoCandidate{propio(c[1], 14)},
			TypoOptions{},
		)
		if len(out) != 0 {
			t.Fatalf("%s -> %s: sugirio cambiar de talla: %+v", c[0], c[1], out)
		}
	}
}

func TestDetectTyposCallaSiHayVariosCandidatos(t *testing.T) {
	out := DetectTypos(
		[]TypoCandidate{canal("TD212-4XL", 5)},
		[]TypoCandidate{propio("TD214-4XL", 5), propio("TD213-4XL", 5)},
		TypoOptions{},
	)
	if len(out) != 0 {
		t.Fatalf("con dos vecinos a un caracter no puede elegir: %+v", out)
	}
}

func TestDetectTyposCallaSiElEspacioEsAmbiguo(t *testing.T) {
	out := DetectTypos(
		[]TypoCandidate{canal("AB 12", 5)},
		[]TypoCandidate{propio("AB12", 5), propio("A B12", 5)},
		TypoOptions{},
	)
	if len(out) != 0 {
		t.Fatalf("dos SKU colapsan al mismo texto: %+v", out)
	}
}

func TestDetectTyposUnCaracterExigeCantidad(t *testing.T) {
	out := DetectTypos(
		[]TypoCandidate{canal("BT11-5XL", 7)},
		[]TypoCandidate{propio("BN11-5XL", 40)},
		TypoOptions{},
	)
	if len(out) != 0 {
		t.Fatalf("sin cantidad que corrobore no deberia sugerir: %+v", out)
	}
}

func TestDetectTyposDosCerosNoSonEvidencia(t *testing.T) {
	out := DetectTypos(
		[]TypoCandidate{canal("BT11-5XL", 0)},
		[]TypoCandidate{propio("BN11-5XL", 0)},
		TypoOptions{},
	)
	if len(out) != 0 {
		t.Fatalf("dos ceros son el valor por defecto, no una coincidencia: %+v", out)
	}
}

func TestDetectTyposUnCaracterConCantidad(t *testing.T) {
	out := DetectTypos(
		[]TypoCandidate{canal("BT11-5XL", 7)},
		[]TypoCandidate{propio("BN11-5XL", 7)},
		TypoOptions{},
	)
	s := unico(t, out)
	if s.Pattern != PatternOneChar {
		t.Fatalf("patron = %q", s.Pattern)
	}
	if s.ChannelQty != 7 || s.ProbabilityQty != 7 {
		t.Fatalf("no viajo la evidencia de cantidades: %+v", s)
	}
}

func TestDetectTyposCorrigeEnProbabilityCuandoElStockEstaEnElCanal(t *testing.T) {
	out := DetectTypos(
		[]TypoCandidate{canal("SDH22-S", 30)},
		[]TypoCandidate{propio("SHD22-S", 0)},
		TypoOptions{},
	)
	s := unico(t, out)
	if s.FixSide != FixInProbability {
		t.Fatalf("lado = %q, el stock vivo esta en el canal", s.FixSide)
	}
}

func TestDetectTyposSinCantidadConocida(t *testing.T) {
	out := DetectTypos(
		[]TypoCandidate{{SKU: "SDH22-S"}},
		[]TypoCandidate{{SKU: "SHD22-S"}},
		TypoOptions{},
	)
	s := unico(t, out)
	if s.Pattern != PatternTransposed {
		t.Fatalf("la transposicion no necesita cantidades: %+v", s)
	}
	if s.HasChannelQty || s.HasProbQty {
		t.Fatalf("no deberia inventar cantidades: %+v", s)
	}
}

func TestDetectTyposIgnoraVacios(t *testing.T) {
	if out := DetectTypos(nil, []TypoCandidate{propio("A-1", 1)}, TypoOptions{}); len(out) != 0 {
		t.Fatalf("sin canal no hay nada que sugerir: %+v", out)
	}
	if out := DetectTypos([]TypoCandidate{canal("A-1", 1)}, nil, TypoOptions{}); len(out) != 0 {
		t.Fatalf("sin catalogo propio no hay nada que sugerir: %+v", out)
	}
	out := DetectTypos(
		[]TypoCandidate{canal("   ", 1), canal("SDH22-S", 5)},
		[]TypoCandidate{propio("", 1), propio("SHD22-S", 5)},
		TypoOptions{},
	)
	if len(out) != 1 {
		t.Fatalf("los vacios no participan: %+v", out)
	}
}

func TestDetectTyposIdenticosSinEspaciosNoSonHallazgo(t *testing.T) {
	out := DetectTypos(
		[]TypoCandidate{canal("A-1", 1)},
		[]TypoCandidate{propio("A-1", 1)},
		TypoOptions{},
	)
	if len(out) != 0 {
		t.Fatalf("mismo SKU exacto no es un hallazgo: %+v", out)
	}
}

func TestDetectTyposSinGuion(t *testing.T) {
	out := DetectTypos(
		[]TypoCandidate{canal("SDH22", 5)},
		[]TypoCandidate{propio("SHD22", 5)},
		TypoOptions{},
	)
	s := unico(t, out)
	if s.Pattern != PatternTransposed {
		t.Fatalf("patron = %q", s.Pattern)
	}
}

func TestSameQuantity(t *testing.T) {
	if sameQuantity(TypoCandidate{Quantity: 5}, propio("X", 5)) {
		t.Fatal("sin HasQty del lado del canal no hay evidencia")
	}
	if sameQuantity(canal("X", 5), TypoCandidate{Quantity: 5}) {
		t.Fatal("sin HasQty del lado propio no hay evidencia")
	}
	if !sameQuantity(canal("X", 7), propio("Y", 8)) {
		t.Fatal("una unidad de diferencia sigue siendo la misma cantidad")
	}
	if !sameQuantity(canal("X", 8), propio("Y", 7)) {
		t.Fatal("la diferencia se toma en valor absoluto")
	}
	if sameQuantity(canal("X", 5), propio("Y", 9)) {
		t.Fatal("cuatro de diferencia no corrobora nada")
	}
}

func TestIsTransposition(t *testing.T) {
	casos := []struct {
		a, b string
		want bool
	}{
		{"SDH", "SHD", true},
		{"AB", "BA", true},
		{"ABC", "ABC", false},
		{"ABCD", "BADC", false},
		{"ACB", "BCA", false},
		{"ABC", "AB", false},
	}
	for _, c := range casos {
		if got := isTransposition(c.a, c.b); got != c.want {
			t.Fatalf("isTransposition(%q,%q) = %v", c.a, c.b, got)
		}
	}
}

func TestEditDistanceCorta(t *testing.T) {
	if d := editDistance("ABCDEF", "A", 1); d <= 1 {
		t.Fatalf("deberia cortar por diferencia de largo, dio %d", d)
	}
	if d := editDistance("ABC", "ABC", 1); d != 0 {
		t.Fatalf("iguales = %d", d)
	}
	if d := editDistance("ABC", "ABD", 1); d != 1 {
		t.Fatalf("un caracter = %d", d)
	}
	if d := editDistance("ABC", "AXY", 1); d <= 1 {
		t.Fatalf("dos caracteres deberia superar el tope, dio %d", d)
	}
}

func TestSplitSKU(t *testing.T) {
	base, suf := splitSKU("TD212-4XL")
	if base != "TD212" || suf != "4XL" {
		t.Fatalf("split = %q / %q", base, suf)
	}
	base, suf = splitSKU("TD212")
	if base != "TD212" || suf != "" {
		t.Fatalf("sin guion = %q / %q", base, suf)
	}
	base, suf = splitSKU("A-B-C")
	if base != "A-B" || suf != "C" {
		t.Fatalf("parte por el ultimo guion, dio %q / %q", base, suf)
	}
}

func TestAutoridadDelERPMandaSobreElStock(t *testing.T) {
	canalConStock := []TypoCandidate{canal("SDH22-S", 30)}
	propioSinStock := []TypoCandidate{propio("SHD22-S", 0)}

	s := unico(t, DetectTypos(canalConStock, propioSinStock, TypoOptions{Authority: AuthorityProbability}))
	if s.FixSide != FixInChannel {
		t.Fatalf("con el ERP alimentando a Probability se corrige en el canal, dio %q", s.FixSide)
	}

	s = unico(t, DetectTypos(canalConStock, propioSinStock, TypoOptions{Authority: AuthorityChannel}))
	if s.FixSide != FixInProbability {
		t.Fatalf("si el canal es la fuente de verdad se corrige en Probability, dio %q", s.FixSide)
	}
}

func TestAutoridadNoTocaLosEspacios(t *testing.T) {
	s := unico(t, DetectTypos(
		[]TypoCandidate{canal("BN14 -5XL", 10)},
		[]TypoCandidate{propio("BN14-5XL", 10)},
		TypoOptions{Authority: AuthorityChannel},
	))
	if s.FixSide != FixInChannel {
		t.Fatalf("el espacio se corrige donde esta, no donde diga la autoridad: %q", s.FixSide)
	}
}

func TestTypoDetailsSeparaEspaciosDeDigitacion(t *testing.T) {
	suspects := []TypoSuspect{
		{
			ChannelSKU: "A -1", ProbabilitySKU: "A-1", Pattern: PatternSpacing, FixSide: FixInChannel,
			ChannelQty: 3, HasChannelQty: true, ProbabilityQty: 9, HasProbQty: true,
		},
		{ChannelSKU: "SDH-1", ProbabilitySKU: "SHD-1", Pattern: PatternTransposed, FixSide: FixInChannel},
	}
	out := TypoDetails(suspects, "Siigo")
	if out[0].Group != GroupSKUSpacing || out[0].Tone != "info" {
		t.Fatalf("el espacio va en su propio grupo: %+v", out[0])
	}
	if out[1].Group != GroupSKUTypo || out[1].Tone != "warn" {
		t.Fatalf("la digitacion va como advertencia: %+v", out[1])
	}
	if out[1].Label != "letras trocadas en Siigo: deberia decir SHD-1" {
		t.Fatalf("instruccion inesperada: %q", out[1].Label)
	}
	if out[0].ChannelQty == nil || *out[0].ChannelQty != 3 || out[0].OwnQty == nil || *out[0].OwnQty != 9 {
		t.Fatalf("las cantidades son la evidencia, tienen que viajar: %+v", out[0])
	}
	if CountPattern(suspects, PatternSpacing) != 1 {
		t.Fatal("conteo por patron incorrecto")
	}
}

func TestFixInstructionApuntaAProbability(t *testing.T) {
	got := FixInstruction(TypoSuspect{
		ChannelSKU: "AB-1", ProbabilitySKU: "AC-1",
		Pattern: PatternOneChar, FixSide: FixInProbability,
	}, "Siigo")
	if got != "un caracter de diferencia y existencias parecidas; en Probability deberia decir AB-1" {
		t.Fatalf("instruccion inesperada: %q", got)
	}
}

func TestTypoDetailsSinCantidadesConocidas(t *testing.T) {
	out := TypoDetails([]TypoSuspect{{
		ChannelSKU: "A -1", ProbabilitySKU: "A-1", Pattern: PatternSpacing, FixSide: FixInChannel,
	}}, "Siigo")
	if out[0].ChannelQty != nil || out[0].OwnQty != nil {
		t.Fatalf("sin HasQty las cantidades viajan vacias: %+v", out[0])
	}
}
