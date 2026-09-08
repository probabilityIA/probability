package usecases

import (
	"context"
	"fmt"
	"testing"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/woocommerce/internal/domain"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/woocommerce/internal/mocks"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/productmatch"
)

func TestBuildVariableAttributes_UnionDeDistintosValores(t *testing.T) {
	products := []domain.ProductForSync{
		{VariantAttributes: map[string]string{"color": "Azul Oscuro", "talla": "3XL"}},
		{VariantAttributes: map[string]string{"color": "Negro", "talla": "3XL"}},
		{VariantAttributes: map[string]string{"color": "Azul Oscuro", "talla": "L"}},
	}
	attrs := buildVariableAttributes(products, []int{0, 1, 2})
	if len(attrs) != 2 {
		t.Fatalf("esperaba 2 atributos (color, talla), obtuve %d", len(attrs))
	}
	byName := map[string][]string{}
	for _, a := range attrs {
		byName[a.Name] = a.Options
	}
	if len(byName["Color"]) != 2 {
		t.Errorf("Color options = %v, esperaba 2 valores distintos", byName["Color"])
	}
	if len(byName["Talla"]) != 2 {
		t.Errorf("Talla options = %v, esperaba 2 valores distintos", byName["Talla"])
	}
}

func TestApplyProductsToWoo_CreaFamiliaVariableConSusVariaciones(t *testing.T) {
	svc := &mocks.IntegrationServiceMock{
		GetIntegrationByIDFn: func(ctx context.Context, integrationID string) (*domain.Integration, error) {
			return &domain.Integration{ID: 221, Config: map[string]interface{}{"store_url": "https://tienda.test"}}, nil
		},
	}

	var createdParentInput domain.CreateVariableProductInput
	var createdVariations []domain.CreateVariationInput
	client := &mocks.WooClientMock{
		GetProductsFn: func(ctx context.Context, storeURL, ck, cs string) ([]domain.WooProduct, error) {
			return nil, nil
		},
		CreateVariableProductFn: func(ctx context.Context, storeURL, ck, cs string, input domain.CreateVariableProductInput) (string, error) {
			createdParentInput = input
			return "500", nil
		},
		CreateProductVariationFn: func(ctx context.Context, storeURL, ck, cs, parentID string, input domain.CreateVariationInput) (string, error) {
			createdVariations = append(createdVariations, input)
			if parentID != "500" {
				t.Errorf("parentID = %s, esperaba 500", parentID)
			}
			return "60" + input.SKU[len(input.SKU)-1:], nil
		},
	}

	var mappedRefs []productmatch.ExternalRefs
	repo := &fakeProductRepo{
		ListProductsByBusinessFn: func(ctx context.Context, businessID uint) ([]domain.ProductForSync, error) {
			return []domain.ProductForSync{
				{ID: "PRD_1", SKU: "BH-1", Name: "Bermuda - Azul - S", FamilyID: "9", FamilyName: "Bermuda Hombre", VariantAttributes: map[string]string{"color": "Azul", "talla": "S"}, TrackInventory: true, StockQuantity: 5, Price: 10},
				{ID: "PRD_2", SKU: "BH-2", Name: "Bermuda - Negro - M", FamilyID: "9", FamilyName: "Bermuda Hombre", VariantAttributes: map[string]string{"color": "Negro", "talla": "M"}, TrackInventory: true, StockQuantity: 8, Price: 10},
			}, nil
		},
	}
	repo.UpsertProductIntegrationMappingFn = func(ctx context.Context, productID string, businessID, integrationID uint, refs productmatch.ExternalRefs) error {
		mappedRefs = append(mappedRefs, refs)
		return nil
	}

	uc := &wooCommerceUseCase{
		client:      client,
		service:     svc,
		productRepo: repo,
		logger:      log.New(),
	}

	if err := uc.ApplyProductsToWoo(context.Background(), "221", 26, "corr-1", "BH-1", "BH-2"); err != nil {
		t.Fatalf("ApplyProductsToWoo error: %v", err)
	}

	if createdParentInput.Name != "Bermuda Hombre" {
		t.Errorf("nombre del padre = %q, esperaba %q", createdParentInput.Name, "Bermuda Hombre")
	}
	if len(createdParentInput.Attributes) != 2 {
		t.Fatalf("atributos del padre = %d, esperaba 2 (color, talla)", len(createdParentInput.Attributes))
	}
	if len(createdVariations) != 2 {
		t.Fatalf("variaciones creadas = %d, esperaba 2", len(createdVariations))
	}
	if len(mappedRefs) != 2 {
		t.Fatalf("mapeos guardados = %d, esperaba 2", len(mappedRefs))
	}
	for _, refs := range mappedRefs {
		if refs.ProductID != "500" {
			t.Errorf("ProductID mapeado = %s, esperaba el padre 500", refs.ProductID)
		}
		if refs.VariantID == "" {
			t.Errorf("VariantID mapeado vacio para sku %s", refs.SKU)
		}
	}
}

func TestApplyProductsToWoo_UsaImagenDeLaFamiliaMadreParaElProductoVariable(t *testing.T) {
	svc := &mocks.IntegrationServiceMock{
		GetIntegrationByIDFn: func(ctx context.Context, integrationID string) (*domain.Integration, error) {
			return &domain.Integration{ID: 221, Config: map[string]interface{}{"store_url": "https://tienda.test"}}, nil
		},
	}

	var createdParentInput domain.CreateVariableProductInput
	client := &mocks.WooClientMock{
		GetProductsFn: func(ctx context.Context, storeURL, ck, cs string) ([]domain.WooProduct, error) { return nil, nil },
		CreateVariableProductFn: func(ctx context.Context, storeURL, ck, cs string, input domain.CreateVariableProductInput) (string, error) {
			createdParentInput = input
			return "500", nil
		},
		CreateProductVariationFn: func(ctx context.Context, storeURL, ck, cs, parentID string, input domain.CreateVariationInput) (string, error) {
			return "601", nil
		},
	}

	repo := &fakeProductRepo{
		ListProductsByBusinessFn: func(ctx context.Context, businessID uint) ([]domain.ProductForSync, error) {
			return []domain.ProductForSync{
				{ID: "PRD_1", SKU: "ARG-TIT-S", Name: "Camiseta Titular - S", FamilyID: "6", FamilyName: "Camisas Seleccion Argentina", FamilyImageURL: "https://tienda.test/madre.jpg", ImageURL: "https://tienda.test/subfamilia.jpg", VariantAttributes: map[string]string{"talla": "S"}, TrackInventory: true, StockQuantity: 5, Price: 10},
			}, nil
		},
	}
	repo.UpsertProductIntegrationMappingFn = func(ctx context.Context, productID string, businessID, integrationID uint, refs productmatch.ExternalRefs) error {
		return nil
	}

	uc := &wooCommerceUseCase{client: client, service: svc, productRepo: repo, logger: log.New()}

	if err := uc.ApplyProductsToWoo(context.Background(), "221", 26, "corr-1", "ARG-TIT-S"); err != nil {
		t.Fatalf("ApplyProductsToWoo error: %v", err)
	}

	if createdParentInput.ImageURL != "https://tienda.test/madre.jpg" {
		t.Errorf("ImageURL del padre = %q, esperaba la imagen de la familia madre", createdParentInput.ImageURL)
	}
}

func TestApplyProductsToWoo_SinImagenDeMadreUsaLaDelPrimerProducto(t *testing.T) {
	svc := &mocks.IntegrationServiceMock{
		GetIntegrationByIDFn: func(ctx context.Context, integrationID string) (*domain.Integration, error) {
			return &domain.Integration{ID: 221, Config: map[string]interface{}{"store_url": "https://tienda.test"}}, nil
		},
	}

	var createdParentInput domain.CreateVariableProductInput
	client := &mocks.WooClientMock{
		GetProductsFn: func(ctx context.Context, storeURL, ck, cs string) ([]domain.WooProduct, error) { return nil, nil },
		CreateVariableProductFn: func(ctx context.Context, storeURL, ck, cs string, input domain.CreateVariableProductInput) (string, error) {
			createdParentInput = input
			return "500", nil
		},
		CreateProductVariationFn: func(ctx context.Context, storeURL, ck, cs, parentID string, input domain.CreateVariationInput) (string, error) {
			return "601", nil
		},
	}

	repo := &fakeProductRepo{
		ListProductsByBusinessFn: func(ctx context.Context, businessID uint) ([]domain.ProductForSync, error) {
			return []domain.ProductForSync{
				{ID: "PRD_1", SKU: "ARG-TIT-S", Name: "Camiseta Titular - S", FamilyID: "6", FamilyName: "Camisas Seleccion Argentina", FamilyImageURL: "", ImageURL: "https://tienda.test/subfamilia.jpg", VariantAttributes: map[string]string{"talla": "S"}, TrackInventory: true, StockQuantity: 5, Price: 10},
			}, nil
		},
	}
	repo.UpsertProductIntegrationMappingFn = func(ctx context.Context, productID string, businessID, integrationID uint, refs productmatch.ExternalRefs) error {
		return nil
	}

	uc := &wooCommerceUseCase{client: client, service: svc, productRepo: repo, logger: log.New()}

	if err := uc.ApplyProductsToWoo(context.Background(), "221", 26, "corr-1", "ARG-TIT-S"); err != nil {
		t.Fatalf("ApplyProductsToWoo error: %v", err)
	}

	if createdParentInput.ImageURL != "https://tienda.test/subfamilia.jpg" {
		t.Errorf("ImageURL del padre = %q, esperaba el fallback a la imagen del producto", createdParentInput.ImageURL)
	}
}

func TestApplyProductsToWoo_ReutilizaParentSiUnHermanoYaEstaAsociado(t *testing.T) {
	svc := &mocks.IntegrationServiceMock{
		GetIntegrationByIDFn: func(ctx context.Context, integrationID string) (*domain.Integration, error) {
			return &domain.Integration{ID: 221, Config: map[string]interface{}{"store_url": "https://tienda.test"}}, nil
		},
	}

	createVariableCalls := 0
	var createdVariations []domain.CreateVariationInput
	client := &mocks.WooClientMock{
		GetProductsFn: func(ctx context.Context, storeURL, ck, cs string) ([]domain.WooProduct, error) {
			return []domain.WooProduct{{ID: "700", SKU: "ARG-TIT-S", Barcode: ""}}, nil
		},
		CreateVariableProductFn: func(ctx context.Context, storeURL, ck, cs string, input domain.CreateVariableProductInput) (string, error) {
			createVariableCalls++
			return "999", nil
		},
		CreateProductVariationFn: func(ctx context.Context, storeURL, ck, cs, parentID string, input domain.CreateVariationInput) (string, error) {
			createdVariations = append(createdVariations, input)
			return "701", nil
		},
	}

	repo := &fakeProductRepo{
		ListProductsByBusinessFn: func(ctx context.Context, businessID uint) ([]domain.ProductForSync, error) {
			return []domain.ProductForSync{
				{ID: "PRD_1", SKU: "ARG-TIT-S", Name: "Camiseta Titular - S", FamilyID: "6", FamilyName: "Camisas Seleccion Argentina", VariantAttributes: map[string]string{"talla": "S"}, TrackInventory: true, StockQuantity: 5, Price: 10},
				{ID: "PRD_2", SKU: "ARG-TIT-M", Name: "Camiseta Titular - M", FamilyID: "6", FamilyName: "Camisas Seleccion Argentina", VariantAttributes: map[string]string{"talla": "M"}, TrackInventory: true, StockQuantity: 5, Price: 10},
			}, nil
		},
	}
	repo.UpsertProductIntegrationMappingFn = func(ctx context.Context, productID string, businessID, integrationID uint, refs productmatch.ExternalRefs) error {
		return nil
	}

	uc := &wooCommerceUseCase{client: client, service: svc, productRepo: repo, logger: log.New()}

	if err := uc.ApplyProductsToWoo(context.Background(), "221", 26, "corr-1", "ARG-TIT-S", "ARG-TIT-M"); err != nil {
		t.Fatalf("ApplyProductsToWoo error: %v", err)
	}

	if createVariableCalls != 0 {
		t.Errorf("CreateVariableProduct se llamo %d veces, esperaba 0 (debia reutilizar el parent 700 ya asociado)", createVariableCalls)
	}
	if len(createdVariations) != 1 {
		t.Fatalf("variaciones creadas = %d, esperaba 1 (solo el sku nuevo ARG-TIT-M)", len(createdVariations))
	}
	if createdVariations[0].SKU != "ARG-TIT-M" {
		t.Errorf("SKU de la variacion creada = %s, esperaba ARG-TIT-M", createdVariations[0].SKU)
	}
}

func TestApplyProductsToWoo_FalloAlCrearElProductoVariableNoIntentaCrearVariaciones(t *testing.T) {
	svc := &mocks.IntegrationServiceMock{
		GetIntegrationByIDFn: func(ctx context.Context, integrationID string) (*domain.Integration, error) {
			return &domain.Integration{ID: 221, Config: map[string]interface{}{"store_url": "https://tienda.test"}}, nil
		},
	}

	variationCalls := 0
	client := &mocks.WooClientMock{
		GetProductsFn: func(ctx context.Context, storeURL, ck, cs string) ([]domain.WooProduct, error) { return nil, nil },
		CreateVariableProductFn: func(ctx context.Context, storeURL, ck, cs string, input domain.CreateVariableProductInput) (string, error) {
			return "", fmt.Errorf("woocommerce caido")
		},
		CreateProductVariationFn: func(ctx context.Context, storeURL, ck, cs, parentID string, input domain.CreateVariationInput) (string, error) {
			variationCalls++
			return "701", nil
		},
	}

	repo := &fakeProductRepo{
		ListProductsByBusinessFn: func(ctx context.Context, businessID uint) ([]domain.ProductForSync, error) {
			return []domain.ProductForSync{
				{ID: "PRD_1", SKU: "ARG-TIT-S", Name: "Camiseta Titular - S", FamilyID: "6", FamilyName: "Camisas Seleccion Argentina", VariantAttributes: map[string]string{"talla": "S"}, TrackInventory: true, StockQuantity: 5, Price: 10},
			}, nil
		},
	}

	uc := &wooCommerceUseCase{client: client, service: svc, productRepo: repo, logger: log.New()}

	if err := uc.ApplyProductsToWoo(context.Background(), "221", 26, "corr-1", "ARG-TIT-S"); err != nil {
		t.Fatalf("ApplyProductsToWoo error: %v", err)
	}

	if variationCalls != 0 {
		t.Errorf("CreateProductVariation se llamo %d veces, esperaba 0 porque el padre fallo", variationCalls)
	}
}
