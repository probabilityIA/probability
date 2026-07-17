package usecases

import (
	"context"
	"testing"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/jumpseller/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
)

func nuevoUseCaseParaMapeo() *jumpsellerUseCase {
	return &jumpsellerUseCase{logger: log.New()}
}

func TestUpsertMsgTraeDimensionesYConviertePeso(t *testing.T) {
	uc := nuevoUseCaseParaMapeo()

	j := jumpsellerSKU{
		SKU:        "JS-001",
		Name:       "Camiseta",
		Price:      75000,
		ExternalID: "100",
		Weight:     1500,
		Height:     12,
		Width:      20,
		Length:     30,
	}

	msg := uc.upsertMsgFromJumpseller(context.Background(), 26, 42, j, "g")

	if msg.Weight == nil || !casiIgual(*msg.Weight, 1.5) {
		t.Fatalf("weight = %v, se esperaban 1.5 kg (la tienda reporta gramos)", msg.Weight)
	}
	if msg.WeightUnit != "kg" {
		t.Fatalf("weight_unit = %q, Probability guarda en kg", msg.WeightUnit)
	}
	if msg.Height == nil || *msg.Height != 12 || msg.Width == nil || *msg.Width != 20 || msg.Length == nil || *msg.Length != 30 {
		t.Fatalf("dimensiones mal mapeadas: h=%v w=%v l=%v", msg.Height, msg.Width, msg.Length)
	}
	if msg.DimensionUnit != "cm" {
		t.Fatalf("dimension_unit = %q, se esperaba cm", msg.DimensionUnit)
	}
}

func TestUpsertMsgSinUnidadDeTiendaNoImportaPeso(t *testing.T) {
	uc := nuevoUseCaseParaMapeo()

	j := jumpsellerSKU{SKU: "JS-001", Name: "Camiseta", Weight: 1500, Height: 12}

	msg := uc.upsertMsgFromJumpseller(context.Background(), 26, 42, j, "")

	if msg.Weight != nil {
		t.Fatalf("weight = %v: sin unidad conocida NO se debe adivinar el peso", *msg.Weight)
	}
	if msg.WeightUnit != "" {
		t.Fatalf("weight_unit = %q, se esperaba vacio", msg.WeightUnit)
	}
	if msg.Height == nil || *msg.Height != 12 {
		t.Fatal("las dimensiones si deben viajar aunque el peso se descarte")
	}
}

func TestUpsertMsgSinDimensionesNoMandaUnidad(t *testing.T) {
	uc := nuevoUseCaseParaMapeo()

	j := jumpsellerSKU{SKU: "JS-001", Name: "Camiseta", Price: 100}

	msg := uc.upsertMsgFromJumpseller(context.Background(), 26, 42, j, "kg")

	if msg.Height != nil || msg.Width != nil || msg.Length != nil {
		t.Fatal("un producto sin dimensiones no debe mandar ceros")
	}
	if msg.DimensionUnit != "" {
		t.Fatalf("dimension_unit = %q: sin dimensiones no se manda unidad", msg.DimensionUnit)
	}
}

func TestPesoDeProbabilityHaciaJumpsellerConvierteAUnidadDeLaTienda(t *testing.T) {
	uc := nuevoUseCaseParaMapeo()
	peso := 2.0

	p := domain.ProductForSync{SKU: "P-1", Weight: &peso, WeightUnit: "kg"}

	enGramos := uc.probabilityWeightForStore(context.Background(), p, "g")
	if enGramos == nil || !casiIgual(*enGramos, 2000) {
		t.Fatalf("2 kg hacia una tienda en gramos = %v, se esperaba 2000", enGramos)
	}

	enKg := uc.probabilityWeightForStore(context.Background(), p, "kg")
	if enKg == nil || !casiIgual(*enKg, 2) {
		t.Fatalf("2 kg hacia una tienda en kg = %v, se esperaba 2", enKg)
	}

	if desconocida := uc.probabilityWeightForStore(context.Background(), p, "piedras"); desconocida != nil {
		t.Fatal("unidad de tienda desconocida: no se debe enviar peso")
	}
}

func TestVariantesHeredanDimensionesDelPadre(t *testing.T) {
	productos := []domain.JumpsellerProduct{{
		ID:     100,
		Name:   "Camiseta",
		SKU:    "PADRE",
		Weight: 1.2,
		Height: 5,
		Width:  10,
		Length: 15,
		Variants: []domain.ProductVariant{
			{ID: 11, SKU: "VAR-S", Price: 100},
			{ID: 12, SKU: "VAR-M", Price: 200},
		},
	}}

	flat := flattenProductSKUs(productos)

	if len(flat) != 2 {
		t.Fatalf("se esperaban 2 SKU (uno por variante), hubo %d", len(flat))
	}
	for _, v := range flat {
		if v.Weight != 1.2 || v.Height != 5 || v.Width != 10 || v.Length != 15 {
			t.Fatalf("la variante %s no heredo las dimensiones del padre: %+v", v.SKU, v)
		}
		if v.ProductID != 100 {
			t.Fatalf("la variante %s perdio el id del padre, no se podria actualizar", v.SKU)
		}
	}
	if flat[0].ExternalID != "100:11" || flat[1].ExternalID != "100:12" {
		t.Fatalf("external_id de variantes mal armado: %s / %s", flat[0].ExternalID, flat[1].ExternalID)
	}
}
