package client

import (
	"context"
	"os"
	"strconv"
	"testing"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/tiendanube/internal/domain"
)

func mockCredential(t *testing.T) domain.Credential {
	t.Helper()
	base := os.Getenv("TIENDANUBE_MOCK_URL")
	if base == "" {
		t.Skip("TIENDANUBE_MOCK_URL no definida, se omite la prueba contra el mock")
	}
	return domain.Credential{
		AccessToken: "mock-tiendanube-token",
		StoreID:     "1234567",
		BaseURL:     base,
		UserAgent:   "Probability Test (qa@probabilityia.com.co)",
	}
}

func TestElClienteLeeLaTiendaDelMock(t *testing.T) {
	c := New()
	cred := mockCredential(t)

	store, err := c.GetStoreInfo(context.Background(), cred)
	if err != nil {
		t.Fatalf("GetStoreInfo fallo: %v", err)
	}
	if store.Name == "" {
		t.Fatalf("el nombre de la tienda llego vacio, el objeto i18n no se resolvio")
	}
	if store.Currency != "COP" {
		t.Fatalf("moneda inesperada: %q", store.Currency)
	}
}

func TestElClienteAplanaVariantesYPreciosEnTexto(t *testing.T) {
	c := New()
	cred := mockCredential(t)

	products, err := c.GetProducts(context.Background(), cred)
	if err != nil {
		t.Fatalf("GetProducts fallo: %v", err)
	}
	if len(products) == 0 {
		t.Fatal("el mock no devolvio productos")
	}

	var encontrado bool
	for _, p := range products {
		if p.Name == "" {
			t.Fatalf("producto %d sin nombre resuelto", p.ID)
		}
		for _, v := range p.Variants {
			if v.SKU == "8013XL" {
				encontrado = true
				if v.Price != 89000 {
					t.Fatalf("el precio en texto no se parseo: %v", v.Price)
				}
				if v.Stock != 12 {
					t.Fatalf("stock inesperado: %d", v.Stock)
				}
			}
		}
	}
	if !encontrado {
		t.Fatal("no se encontro el SKU sembrado 8013XL")
	}
}

func TestElClienteCreaActualizaYEscribeStock(t *testing.T) {
	c := New()
	cred := mockCredential(t)
	ctx := context.Background()

	sku := "TN-TEST-" + strconv.FormatInt(int64(os.Getpid()), 10)
	peso := 1.25

	productID, variantID, err := c.CreateProduct(ctx, cred, domain.CreateProductInput{
		Name:          "Producto creado por Probability",
		SKU:           sku,
		Price:         123456.78,
		Description:   "creado en la prueba",
		StockQuantity: 9,
		ManageStock:   true,
		Weight:        &peso,
	})
	if err != nil {
		t.Fatalf("CreateProduct fallo: %v", err)
	}
	if productID == 0 || variantID == 0 {
		t.Fatalf("Tiendanube no devolvio ids: product=%d variant=%d", productID, variantID)
	}

	if err := c.UpdateProduct(ctx, cred, productID, domain.UpdateProductInput{Name: "Nombre actualizado"}); err != nil {
		t.Fatalf("UpdateProduct fallo: %v", err)
	}

	nuevoPrecio := 99999.0
	if err := c.UpdateVariant(ctx, cred, productID, variantID, domain.UpdateVariantInput{Price: &nuevoPrecio}); err != nil {
		t.Fatalf("UpdateVariant fallo: %v", err)
	}

	if err := c.SetVariantStock(ctx, cred, productID, variantID, 42); err != nil {
		t.Fatalf("SetVariantStock fallo: %v", err)
	}

	target, err := c.ResolveStockTarget(ctx, cred, sku)
	if err != nil {
		t.Fatalf("ResolveStockTarget fallo: %v", err)
	}
	if !target.Found || target.ProductID != productID || target.VariantID != variantID {
		t.Fatalf("ResolveStockTarget no encontro la variante: %+v", target)
	}

	externalID := strconv.FormatInt(productID, 10) + ":" + strconv.FormatInt(variantID, 10)
	stock, err := c.GetProductsStock(ctx, cred, []string{externalID})
	if err != nil {
		t.Fatalf("GetProductsStock fallo: %v", err)
	}
	if len(stock) != 1 {
		t.Fatalf("se esperaba una fila de stock, llegaron %d", len(stock))
	}
	if !stock[0].Found || stock[0].Quantity != 42 || !stock[0].ManageStock {
		t.Fatalf("el stock del canal no refleja la escritura: %+v", stock[0])
	}

	products, err := c.GetProducts(ctx, cred)
	if err != nil {
		t.Fatalf("GetProducts tras crear fallo: %v", err)
	}
	for _, p := range products {
		if p.ID != productID {
			continue
		}
		if p.Name != "Nombre actualizado" {
			t.Fatalf("el nombre no se actualizo: %q", p.Name)
		}
		if len(p.Variants) != 1 || p.Variants[0].Price != nuevoPrecio {
			t.Fatalf("el precio de la variante no se actualizo: %+v", p.Variants)
		}
		return
	}
	t.Fatalf("el producto creado %d no aparece en el listado", productID)
}

func TestElStockDesconocidoSeReportaComoNoEncontrado(t *testing.T) {
	c := New()
	cred := mockCredential(t)

	stock, err := c.GetProductsStock(context.Background(), cred, []string{"999999:888888"})
	if err != nil {
		t.Fatalf("GetProductsStock fallo: %v", err)
	}
	if len(stock) != 1 || stock[0].Found {
		t.Fatalf("una publicacion inexistente deberia venir como no encontrada: %+v", stock)
	}
}
