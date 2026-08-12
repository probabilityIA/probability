package usecases

import (
	"testing"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/meli/internal/domain"
)

func TestDetectSKUChangesAvisaCuandoElCanalCambioElSKU(t *testing.T) {
	mapped := []domain.MappedItem{
		{SKU: "LD26-2XL", ExternalItemID: "MCO1", ExternalVariantID: "10", ExternalSKU: "LD26-2XL"},
	}
	channel := []domain.MeliProduct{
		{ID: "MCO1", VariationID: "10", SKU: "LD26-3XL", ParentName: "Leguins", VariantLabel: "Negro / 3XL"},
	}

	out := detectSKUChanges(mapped, channel)
	if len(out) != 1 {
		t.Fatalf("esperaba 1 cambio detectado, obtuve %d", len(out))
	}
	if out[0].MatchedBy != "LD26-2XL" || out[0].MatchedValue != "LD26-3XL" {
		t.Fatalf("debia reportar el SKU viejo y el nuevo, obtuve %q -> %q", out[0].MatchedBy, out[0].MatchedValue)
	}
	if out[0].ParentRef != "MCO1" {
		t.Fatalf("debia traer la publicacion padre, obtuve %q", out[0].ParentRef)
	}
}

func TestDetectSKUChangesNoAvisaSiSigueIgual(t *testing.T) {
	mapped := []domain.MappedItem{
		{SKU: "LD26-2XL", ExternalItemID: "MCO1", ExternalVariantID: "10", ExternalSKU: "LD26-2XL"},
	}
	channel := []domain.MeliProduct{
		{ID: "MCO1", VariationID: "10", SKU: "LD26-2XL"},
	}
	if out := detectSKUChanges(mapped, channel); len(out) != 0 {
		t.Fatalf("no debia reportar nada, obtuve %d", len(out))
	}
}

func TestDetectSKUChangesIgnoraDiferenciasDeFormato(t *testing.T) {
	mapped := []domain.MappedItem{
		{SKU: "LD26-2XL", ExternalItemID: "MCO1", ExternalVariantID: "10", ExternalSKU: "LD26-2XL"},
	}
	channel := []domain.MeliProduct{
		{ID: "MCO1", VariationID: "10", SKU: " ld26-2xl "},
	}
	if out := detectSKUChanges(mapped, channel); len(out) != 0 {
		t.Fatalf("mayusculas y espacios no son un cambio real, obtuve %d", len(out))
	}
}

func TestDetectSKUChangesIgnoraLaVarianteQueYaNoEsta(t *testing.T) {
	mapped := []domain.MappedItem{
		{SKU: "LD26-2XL", ExternalItemID: "MCO1", ExternalVariantID: "999", ExternalSKU: "LD26-2XL"},
	}
	channel := []domain.MeliProduct{
		{ID: "MCO1", VariationID: "10", SKU: "LD26-3XL"},
	}
	if out := detectSKUChanges(mapped, channel); len(out) != 0 {
		t.Fatalf("una variante ausente no es un cambio de SKU, obtuve %d", len(out))
	}
}

func TestDetectSKUChangesNoConfundePublicacionesDistintas(t *testing.T) {
	mapped := []domain.MappedItem{
		{SKU: "LD26-2XL", ExternalItemID: "MCO1", ExternalVariantID: "10", ExternalSKU: "LD26-2XL"},
	}
	channel := []domain.MeliProduct{
		{ID: "MCO2", VariationID: "10", SKU: "OTRO-SKU"},
	}
	if out := detectSKUChanges(mapped, channel); len(out) != 0 {
		t.Fatalf("mismo variation_id en otra publicacion no es la misma variante, obtuve %d", len(out))
	}
}

func TestDetectSKUChangesIgnoraElMapeoSinSKUGuardado(t *testing.T) {
	mapped := []domain.MappedItem{
		{SKU: "LD26-2XL", ExternalItemID: "MCO1", ExternalVariantID: "10", ExternalSKU: ""},
	}
	channel := []domain.MeliProduct{{ID: "MCO1", VariationID: "10", SKU: "LD26-3XL"}}

	if out := detectSKUChanges(mapped, channel); len(out) != 0 {
		t.Fatalf("sin copia guardada no hay contra que comparar, obtuve %d", len(out))
	}
}

func TestDetectSKUChangesIgnoraLaVarianteQueQuedoSinSKUEnElCanal(t *testing.T) {
	mapped := []domain.MappedItem{
		{SKU: "LD26-2XL", ExternalItemID: "MCO1", ExternalVariantID: "10", ExternalSKU: "LD26-2XL"},
	}
	channel := []domain.MeliProduct{{ID: "MCO1", VariationID: "10", SKU: ""}}

	if out := detectSKUChanges(mapped, channel); len(out) != 0 {
		t.Fatalf("le quitaron el SKU: eso lo reporta el grupo sin-SKU, no este, obtuve %d", len(out))
	}
}

func TestParentRefSoloAplicaAVariantes(t *testing.T) {
	if got := parentRef(domain.MeliProduct{ID: "MCO1", VariationID: "10"}); got != "MCO1" {
		t.Fatalf("una variante debe reportar su publicacion padre, obtuve %q", got)
	}
	if got := parentRef(domain.MeliProduct{ID: "MCO1"}); got != "" {
		t.Fatalf("una publicacion simple no tiene padre, obtuve %q", got)
	}
}
