package productmatch

import "testing"

func TestResolveFallsBackToDefault(t *testing.T) {
	rules := Resolve(nil, nil)
	if len(rules) != 1 || rules[0].Probability != FieldSKU || rules[0].Channel != FieldSKU {
		t.Fatalf("esperaba el default sku->sku, obtuve %+v", rules)
	}
}

func TestResolvePrefersIntegrationOverType(t *testing.T) {
	integration := []byte(`[{"probability":"barcode","channel":"barcode"}]`)
	typeDefault := []byte(`[{"probability":"sku","channel":"sku"}]`)
	rules := Resolve(integration, typeDefault)
	if len(rules) != 1 || rules[0].Probability != FieldBarcode {
		t.Fatalf("esperaba la regla de la integracion, obtuve %+v", rules)
	}
}

func TestResolveIgnoresInvalidRules(t *testing.T) {
	rules := Resolve([]byte(`[{"probability":"nope","channel":"sku"}]`), nil)
	if len(rules) != 1 || rules[0].Probability != FieldSKU {
		t.Fatalf("esperaba caer al default, obtuve %+v", rules)
	}
}

func TestParseDropsDuplicates(t *testing.T) {
	rules := Parse([]byte(`[{"probability":"sku","channel":"sku"},{"probability":"SKU","channel":"sku"}]`))
	if len(rules) != 1 {
		t.Fatalf("esperaba 1 regla deduplicada, obtuve %+v", rules)
	}
}

func TestReconcileCascadePrefersFirstRule(t *testing.T) {
	rules := []Rule{
		{Probability: FieldSKU, Channel: FieldSKU},
		{Probability: FieldBarcode, Channel: FieldBarcode},
	}
	prob := []Item{
		{SKU: "ABC", Barcode: "7701"},
		{SKU: "ZZZ", Barcode: "7702"},
	}
	channel := []Item{
		{SKU: "abc", Barcode: "9999"},
		{SKU: "", Barcode: "7702"},
	}

	out := Reconcile(rules, prob, channel)
	if len(out.Pairs) != 2 {
		t.Fatalf("esperaba 2 pares, obtuve %+v", out.Pairs)
	}
	if out.Pairs[0].Rule.Channel != FieldSKU {
		t.Fatalf("el primer par debia casar por sku, obtuve %+v", out.Pairs[0].Rule)
	}
	if out.Pairs[1].Rule.Channel != FieldBarcode {
		t.Fatalf("el segundo par debia casar por barcode, obtuve %+v", out.Pairs[1].Rule)
	}
	if len(out.OnlyInProbability) != 0 || len(out.OnlyInChannel) != 0 {
		t.Fatalf("no esperaba sobrantes: %+v %+v", out.OnlyInProbability, out.OnlyInChannel)
	}
}

func TestReconcileReportsLeftovers(t *testing.T) {
	rules := DefaultRules()
	prob := []Item{{SKU: "A"}, {SKU: "B"}}
	channel := []Item{{SKU: "a"}, {SKU: "C"}}

	out := Reconcile(rules, prob, channel)
	if len(out.Pairs) != 1 {
		t.Fatalf("esperaba 1 par, obtuve %+v", out.Pairs)
	}
	if len(out.OnlyInChannel) != 1 || out.OnlyInChannel[0] != 1 {
		t.Fatalf("esperaba el indice 1 solo en canal, obtuve %+v", out.OnlyInChannel)
	}
	if len(out.OnlyInProbability) != 1 || out.OnlyInProbability[0] != 1 {
		t.Fatalf("esperaba el indice 1 solo en probability, obtuve %+v", out.OnlyInProbability)
	}
}

func TestReconcileCountsUnmatchable(t *testing.T) {
	out := Reconcile(DefaultRules(), []Item{{SKU: ""}}, []Item{{SKU: ""}})
	if out.ProbabilityUnmatchable != 1 || out.ChannelUnmatchable != 1 {
		t.Fatalf("esperaba 1 y 1, obtuve %d y %d", out.ProbabilityUnmatchable, out.ChannelUnmatchable)
	}
}

func TestReconcileSkipsAmbiguousKeys(t *testing.T) {
	prob := []Item{{SKU: "DUP"}, {SKU: "dup"}}
	channel := []Item{{SKU: "DUP"}}

	out := Reconcile(DefaultRules(), prob, channel)
	if len(out.Pairs) != 0 {
		t.Fatalf("un sku duplicado en probability no debe casar, obtuve %+v", out.Pairs)
	}
	if len(out.OnlyInChannel) != 1 {
		t.Fatalf("el item de canal debia quedar sin asociar, obtuve %+v", out.OnlyInChannel)
	}
}

func TestAllowedByOptionsFiltersUnsupportedFields(t *testing.T) {
	rules := []Rule{
		{Probability: FieldBarcode, Channel: FieldBarcode},
		{Probability: FieldSKU, Channel: FieldSKU},
	}
	opts := Options{Probability: []string{FieldSKU}, Channel: []string{FieldSKU}}
	filtered := AllowedByOptions(rules, opts)
	if len(filtered) != 1 || filtered[0].Channel != FieldSKU {
		t.Fatalf("esperaba solo sku->sku, obtuve %+v", filtered)
	}
}
