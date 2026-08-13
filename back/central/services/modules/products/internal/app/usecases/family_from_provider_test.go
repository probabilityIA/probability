package usecases

import (
	"context"
	"testing"

	"github.com/secamc93/probability/back/central/services/modules/products/internal/domain"
	"github.com/secamc93/probability/back/central/services/modules/products/internal/mocks"
)

func mensaje(atributos map[string]string) *domain.ProductProviderUpsertDTO {
	return &domain.ProductProviderUpsertDTO{
		BusinessID:        46,
		SKU:               "JH348-4XL",
		Name:              "Jogger - Gris / 4XL",
		FamilyName:        "Jogger Hombre Confort Plus",
		VariantLabel:      "Gris / 4XL",
		VariantAttributes: atributos,
	}
}

func TestSinAtributosNoSeAsignaFamilia(t *testing.T) {
	uc := New(&mocks.RepositoryMock{})

	variante := uc.resolverVariante(context.Background(), mensaje(nil))
	if variante.FamilyID != nil {
		t.Error("sin color ni talla no hay firma que distinga la variante, no debe entrar a una familia")
	}
}

func TestCreaLaFamiliaCuandoNoExiste(t *testing.T) {
	creadas := 0
	repo := &mocks.RepositoryMock{
		GetProductFamilyByNameFn: func(context.Context, uint, string) (*domain.ProductFamily, error) {
			return nil, nil
		},
		CreateProductFamilyFn: func(_ context.Context, f *domain.ProductFamily) error {
			creadas++
			f.ID = 77
			return nil
		},
	}

	variante := New(repo).resolverVariante(context.Background(), mensaje(map[string]string{"color": "Gris", "talla": "4XL"}))

	if creadas != 1 {
		t.Fatalf("esperaba crear 1 familia, creo %d", creadas)
	}
	if variante.FamilyID == nil || *variante.FamilyID != 77 {
		t.Fatalf("el producto debe quedar en la familia creada, quedo en %v", variante.FamilyID)
	}
	if variante.Label != "Gris / 4XL" {
		t.Errorf("label = %q", variante.Label)
	}
	if len(variante.Attributes) == 0 {
		t.Error("los atributos son los que generan la firma, no pueden ir vacios")
	}
}

func TestReusaLaFamiliaExistente(t *testing.T) {
	creadas := 0
	repo := &mocks.RepositoryMock{
		GetProductFamilyByNameFn: func(context.Context, uint, string) (*domain.ProductFamily, error) {
			return &domain.ProductFamily{ID: 12, Name: "Jogger Hombre Confort Plus"}, nil
		},
		CreateProductFamilyFn: func(context.Context, *domain.ProductFamily) error {
			creadas++
			return nil
		},
	}

	variante := New(repo).resolverVariante(context.Background(), mensaje(map[string]string{"color": "Gris", "talla": "4XL"}))

	if creadas != 0 {
		t.Error("dos publicaciones con el mismo nombre son la misma familia, no se debe duplicar")
	}
	if variante.FamilyID == nil || *variante.FamilyID != 12 {
		t.Fatalf("esperaba la familia 12, obtuve %v", variante.FamilyID)
	}
}

func TestFirmaOcupadaNoRompeLaCreacion(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetProductFamilyByNameFn: func(context.Context, uint, string) (*domain.ProductFamily, error) {
			return &domain.ProductFamily{ID: 12}, nil
		},
		VariantExistsInFamilyFn: func(context.Context, uint, uint, string, *string) (bool, error) {
			return true, nil
		},
	}

	variante := New(repo).resolverVariante(context.Background(), mensaje(map[string]string{"color": "Gris", "talla": "4XL"}))

	if variante.FamilyID != nil {
		t.Error("si esa combinacion ya existe en la familia, el producto entra sin familia en vez de fallar")
	}
	if variante.Label == "" {
		t.Error("la etiqueta de variante se conserva aunque no haya familia")
	}
}

func TestNoLeQuitaLaFamiliaAlQueYaTieneUna(t *testing.T) {
	familia := uint(9)
	consultada := false
	repo := &mocks.RepositoryMock{
		GetProductBySKUFn: func(context.Context, uint, string) (*domain.Product, error) {
			return &domain.Product{ID: "PRD_1", BusinessID: 46, SKU: "JH348-4XL", FamilyID: &familia}, nil
		},
		GetProductFamilyByNameFn: func(context.Context, uint, string) (*domain.ProductFamily, error) {
			consultada = true
			return &domain.ProductFamily{ID: 12}, nil
		},
	}

	dto := mensaje(map[string]string{"color": "Gris", "talla": "4XL"})
	_ = New(repo).UpsertFromProvider(context.Background(), dto)

	if consultada {
		t.Error("si el producto ya tiene familia no se toca: la del canal no manda sobre la que ya existe")
	}
}
