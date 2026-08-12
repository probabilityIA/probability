package client

import "testing"

func itemSinVariantes() meliItemDetail {
	return meliItemDetail{
		ID:                "MCO1",
		Title:             "Camiseta",
		Price:             39900,
		AvailableQuantity: 12,
		SellerCustomField: "CAM-1",
		Thumbnail:         "http://thumb",
	}
}

func itemConDosVariantes() meliItemDetail {
	d := meliItemDetail{
		ID:    "MCO2",
		Title: "Leguins Tallas Plus",
		Price: 59900,
	}
	d.Shipping.LogisticType = "xd_drop_off"
	d.Variations = []meliItemVariation{
		{
			ID:                10,
			AvailableQuantity: 5,
			Attributes:        []meliItemAttribute{{ID: "SELLER_SKU", ValueName: "LD26-2XL"}},
			AttributeCombinations: []struct {
				Name      string `json:"name"`
				ValueName string `json:"value_name"`
			}{{Name: "Color", ValueName: "Negro"}, {Name: "Talla", ValueName: "2XL"}},
		},
		{
			ID:                11,
			AvailableQuantity: 3,
			SellerCustomField: "LD26-3XL",
			AttributeCombinations: []struct {
				Name      string `json:"name"`
				ValueName string `json:"value_name"`
			}{{Name: "Color", ValueName: "Negro"}, {Name: "Talla", ValueName: "3XL"}},
		},
	}
	return d
}

func TestItemSinVariantesDaUnaSolaFila(t *testing.T) {
	out := itemToProducts(itemSinVariantes())
	if len(out) != 1 {
		t.Fatalf("esperaba 1 fila, obtuve %d", len(out))
	}
	p := out[0]
	if p.SKU != "CAM-1" || p.VariationID != "" || p.StockQuantity != 12 {
		t.Fatalf("fila incorrecta: %+v", p)
	}
	if p.Name != "Camiseta" {
		t.Fatalf("el nombre debia ser el titulo, obtuve %q", p.Name)
	}
}

func TestItemConVariantesDaUnaFilaPorVariante(t *testing.T) {
	out := itemToProducts(itemConDosVariantes())
	if len(out) != 2 {
		t.Fatalf("esperaba 2 filas, obtuve %d", len(out))
	}
	for _, p := range out {
		if p.ID != "MCO2" {
			t.Fatalf("todas comparten la publicacion padre, obtuve %q", p.ID)
		}
		if p.ParentName != "Leguins Tallas Plus" {
			t.Fatalf("debia traer el nombre del padre, obtuve %q", p.ParentName)
		}
		if p.LogisticType != "xd_drop_off" {
			t.Fatalf("debia traer el tipo logistico, obtuve %q", p.LogisticType)
		}
	}
}

func TestElSKUDeLaVarianteSaleDeSusAtributos(t *testing.T) {
	out := itemToProducts(itemConDosVariantes())
	if out[0].SKU != "LD26-2XL" {
		t.Fatalf("debia leer SELLER_SKU de los atributos de la variante, obtuve %q", out[0].SKU)
	}
}

func TestElSellerCustomFieldDeLaVarianteTienePrioridad(t *testing.T) {
	out := itemToProducts(itemConDosVariantes())
	if out[1].SKU != "LD26-3XL" {
		t.Fatalf("debia usar seller_custom_field de la variante, obtuve %q", out[1].SKU)
	}
}

func TestCadaVarianteLlevaSuPropiaCantidad(t *testing.T) {
	out := itemToProducts(itemConDosVariantes())
	if out[0].StockQuantity != 5 || out[1].StockQuantity != 3 {
		t.Fatalf("las cantidades no son por publicacion sino por variante, obtuve %d y %d",
			out[0].StockQuantity, out[1].StockQuantity)
	}
}

func TestElNombreDeLaVarianteIncluyeLaCombinacion(t *testing.T) {
	out := itemToProducts(itemConDosVariantes())
	if out[0].Name != "Leguins Tallas Plus - Negro / 2XL" {
		t.Fatalf("obtuve %q", out[0].Name)
	}
	if out[0].VariantLabel != "Negro / 2XL" {
		t.Fatalf("la etiqueta de variante debia ser solo la combinacion, obtuve %q", out[0].VariantLabel)
	}
}

func TestLaVarianteSinSKUNoSeDescarta(t *testing.T) {
	d := itemConDosVariantes()
	d.Variations[0].Attributes = nil
	d.Variations[0].SellerCustomField = ""

	out := itemToProducts(d)
	if len(out) != 2 {
		t.Fatalf("una variante sin SKU debe seguir apareciendo para poder reportarla, obtuve %d filas", len(out))
	}
	if out[0].SKU != "" {
		t.Fatalf("debia quedar con SKU vacio, obtuve %q", out[0].SKU)
	}
}

func TestVariacionSinCombinacionUsaElTitulo(t *testing.T) {
	d := itemConDosVariantes()
	d.Variations[0].AttributeCombinations = nil

	out := itemToProducts(d)
	if out[0].Name != "Leguins Tallas Plus" {
		t.Fatalf("sin combinacion el nombre es el titulo, obtuve %q", out[0].Name)
	}
	if out[0].VariantLabel != "" {
		t.Fatalf("sin combinacion no hay etiqueta, obtuve %q", out[0].VariantLabel)
	}
}

func TestLaCombinacionIgnoraValoresVacios(t *testing.T) {
	d := itemConDosVariantes()
	d.Variations[0].AttributeCombinations[0].ValueName = "   "

	if got := variationLabel(d.Variations[0]); got != "2XL" {
		t.Fatalf("los valores vacios no deben ensuciar la etiqueta, obtuve %q", got)
	}
}

func TestExtractSKUUsaElAtributoCuandoNoHayCustomField(t *testing.T) {
	d := itemSinVariantes()
	d.SellerCustomField = ""
	d.Attributes = []meliItemAttribute{{ID: "SELLER_SKU", ValueName: "DEL-ATRIBUTO"}}

	if got := extractSKU(d); got != "DEL-ATRIBUTO" {
		t.Fatalf("obtuve %q", got)
	}
}

func TestExtractSKUVacioSiNoHayNinguno(t *testing.T) {
	d := itemSinVariantes()
	d.SellerCustomField = "  "
	d.Attributes = []meliItemAttribute{{ID: "BRAND", ValueName: "Viga"}}

	if got := extractSKU(d); got != "" {
		t.Fatalf("obtuve %q", got)
	}
}

func TestExtractBarcodeReconoceLosTresFormatos(t *testing.T) {
	for _, id := range []string{"GTIN", "EAN", "UPC"} {
		d := itemSinVariantes()
		d.Attributes = []meliItemAttribute{{ID: id, ValueName: "7701234567890"}}
		if got := extractBarcode(d); got != "7701234567890" {
			t.Fatalf("%s: obtuve %q", id, got)
		}
	}
}

func TestExtractBarcodeIgnoraValoresVacios(t *testing.T) {
	d := itemSinVariantes()
	d.Attributes = []meliItemAttribute{{ID: "GTIN", ValueName: "  "}}
	if got := extractBarcode(d); got != "" {
		t.Fatalf("obtuve %q", got)
	}
}

func TestImagenPrefiereLaSecureURL(t *testing.T) {
	d := itemSinVariantes()
	d.Pictures = []struct {
		SecureURL string `json:"secure_url"`
		URL       string `json:"url"`
	}{{SecureURL: "https://segura", URL: "http://insegura"}}

	if got := itemImageURL(d); got != "https://segura" {
		t.Fatalf("obtuve %q", got)
	}
}

func TestImagenCaeAlThumbnailSiNoHayFotos(t *testing.T) {
	if got := itemImageURL(itemSinVariantes()); got != "http://thumb" {
		t.Fatalf("obtuve %q", got)
	}
}

func TestImagenUsaLaURLPlanaSiNoHaySecure(t *testing.T) {
	d := itemSinVariantes()
	d.Pictures = []struct {
		SecureURL string `json:"secure_url"`
		URL       string `json:"url"`
	}{{SecureURL: "", URL: "http://plana"}}

	if got := itemImageURL(d); got != "http://plana" {
		t.Fatalf("obtuve %q", got)
	}
}

func TestImagenCaeAlThumbnailSiLasFotosNoTienenURL(t *testing.T) {
	d := itemSinVariantes()
	d.Pictures = []struct {
		SecureURL string `json:"secure_url"`
		URL       string `json:"url"`
	}{{}}

	if got := itemImageURL(d); got != "http://thumb" {
		t.Fatalf("obtuve %q", got)
	}
}
