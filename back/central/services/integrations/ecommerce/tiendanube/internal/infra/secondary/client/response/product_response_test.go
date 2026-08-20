package response

import (
	"encoding/json"
	"testing"
)

func TestElNombreI18nSeResuelveConPreferenciaDeIdioma(t *testing.T) {
	casos := []struct {
		nombre   string
		json     string
		esperado string
	}{
		{"prefiere espanol", `{"en":"Shirt","es":"Camiseta","pt":"Camisa"}`, "Camiseta"},
		{"cae a portugues si no hay espanol", `{"en":"Shirt","pt":"Camisa"}`, "Camisa"},
		{"cae a ingles", `{"en":"Shirt"}`, "Shirt"},
		{"variante regional de espanol", `{"es_mx":"Playera"}`, "Playera"},
		{"texto plano en vez de objeto", `"Camiseta suelta"`, "Camiseta suelta"},
		{"objeto vacio", `{}`, ""},
		{"nulo", `null`, ""},
		{"todos en blanco", `{"es":"   ","en":""}`, ""},
	}

	for _, c := range casos {
		var l Localized
		if err := json.Unmarshal([]byte(c.json), &l); err != nil {
			t.Fatalf("%s: no se pudo leer %s: %v", c.nombre, c.json, err)
		}
		if got := l.First(); got != c.esperado {
			t.Fatalf("%s: se esperaba %q y llego %q", c.nombre, c.esperado, got)
		}
	}
}

func TestElIdiomaDesconocidoIgualDevuelveAlgo(t *testing.T) {
	var l Localized
	if err := json.Unmarshal([]byte(`{"de":"Hemd"}`), &l); err != nil {
		t.Fatalf("no se pudo leer: %v", err)
	}
	if l.First() != "Hemd" {
		t.Fatalf("con un idioma que no conocemos igual hay que mostrar el nombre, llego %q", l.First())
	}
}

func TestLosPreciosLleganComoTextoYSeParsean(t *testing.T) {
	casos := []struct {
		json     string
		esperado float64
	}{
		{`"25.00"`, 25},
		{`"1234.56"`, 1234.56},
		{`25.5`, 25.5},
		{`""`, 0},
		{`null`, 0},
		{`"no es numero"`, 0},
		{`"  99.9  "`, 99.9},
	}

	for _, c := range casos {
		var n Numeric
		if err := json.Unmarshal([]byte(c.json), &n); err != nil {
			t.Fatalf("no se pudo leer %s: %v", c.json, err)
		}
		if float64(n) != c.esperado {
			t.Fatalf("con %s se esperaba %v y llego %v", c.json, c.esperado, float64(n))
		}
	}
}

func TestElStockVacioSignificaIlimitado(t *testing.T) {
	var ilimitado Stock
	if err := json.Unmarshal([]byte(`""`), &ilimitado); err != nil {
		t.Fatalf("no se pudo leer: %v", err)
	}
	if !ilimitado.Infinite {
		t.Fatal("en Tiendanube un stock vacio es ilimitado: confundirlo con cero deja de vender")
	}

	var nulo Stock
	if err := json.Unmarshal([]byte(`null`), &nulo); err != nil {
		t.Fatalf("no se pudo leer: %v", err)
	}
	if !nulo.Infinite {
		t.Fatal("un stock nulo tambien es ilimitado")
	}

	var cero Stock
	if err := json.Unmarshal([]byte(`0`), &cero); err != nil {
		t.Fatalf("no se pudo leer: %v", err)
	}
	if cero.Infinite || cero.Value != 0 {
		t.Fatalf("cero es cero, no ilimitado: %+v", cero)
	}

	var texto Stock
	if err := json.Unmarshal([]byte(`"12"`), &texto); err != nil {
		t.Fatalf("no se pudo leer: %v", err)
	}
	if texto.Value != 12 || texto.Infinite {
		t.Fatalf("un stock en texto debe parsearse: %+v", texto)
	}

	var invalido Stock
	if err := json.Unmarshal([]byte(`"doce"`), &invalido); err != nil {
		t.Fatalf("no se pudo leer: %v", err)
	}
	if invalido.Value != 0 {
		t.Fatalf("un stock que no es numero queda en cero: %+v", invalido)
	}
}

func TestElProductoSeTraduceAlDominioConSusVariantes(t *testing.T) {
	crudo := `{
		"id": 2001,
		"name": {"es": "Camiseta"},
		"description": {"es": "Una camiseta"},
		"published": true,
		"variants": [
			{"id": 9001, "sku": " SKU-1 ", "barcode": " 777 ", "price": "100.50", "stock": 5, "stock_management": true, "weight": "1.5", "depth": "30", "width": "20", "height": "10"}
		],
		"images": [{"id": 1, "src": "", "position": 1}, {"id": 2, "src": "https://x.com/a.jpg", "position": 2}]
	}`

	var p Product
	if err := json.Unmarshal([]byte(crudo), &p); err != nil {
		t.Fatalf("no se pudo leer el producto: %v", err)
	}

	d := p.ToDomain()
	if d.Name != "Camiseta" || d.Description != "Una camiseta" {
		t.Fatalf("nombre o descripcion mal resueltos: %+v", d)
	}
	if d.ImageURL != "https://x.com/a.jpg" {
		t.Fatalf("debe tomarse la primera imagen con src real, llego %q", d.ImageURL)
	}
	if len(d.Variants) != 1 {
		t.Fatalf("se perdio la variante: %+v", d.Variants)
	}

	v := d.Variants[0]
	if v.SKU != "SKU-1" || v.Barcode != "777" {
		t.Fatalf("sku y barcode deben venir recortados: %q %q", v.SKU, v.Barcode)
	}
	if v.Price != 100.5 || v.Stock != 5 || !v.StockManagement {
		t.Fatalf("precio o stock mal traducidos: %+v", v)
	}
	if v.Weight != 1.5 || v.Depth != 30 || v.Width != 20 || v.Height != 10 {
		t.Fatalf("medidas mal traducidas: %+v", v)
	}
	if v.ProductID != 2001 {
		t.Fatalf("la variante debe saber a que producto pertenece: %d", v.ProductID)
	}
}

func TestUnProductoSinImagenesNoRompe(t *testing.T) {
	var p Product
	if err := json.Unmarshal([]byte(`{"id":1,"name":{"es":"X"},"variants":[]}`), &p); err != nil {
		t.Fatalf("no se pudo leer: %v", err)
	}
	d := p.ToDomain()
	if d.ImageURL != "" || len(d.Variants) != 0 {
		t.Fatalf("sin imagenes ni variantes debe quedar vacio, no fallar: %+v", d)
	}
}

func TestLaTiendaResuelveSuNombreI18n(t *testing.T) {
	var s Store
	if err := json.Unmarshal([]byte(`{"id":7,"name":{"es":"Mi tienda"},"url":"https://x.com","country":"CO","currency":"COP","main_language":"es"}`), &s); err != nil {
		t.Fatalf("no se pudo leer la tienda: %v", err)
	}
	if s.Name.First() != "Mi tienda" || s.Currency != "COP" || s.ID != 7 {
		t.Fatalf("tienda mal leida: %+v", s)
	}
}

func TestUnJSONCorruptoSeReportaComoError(t *testing.T) {
	var l Localized
	if err := json.Unmarshal([]byte(`{"es": 5}`), &l); err == nil {
		t.Fatal("un nombre que no es texto debe reportarse, no ignorarse")
	}
	var n Numeric
	if err := json.Unmarshal([]byte(`{}`), &n); err == nil {
		t.Fatal("un precio que es objeto debe reportarse")
	}
	var s Stock
	if err := json.Unmarshal([]byte(`{}`), &s); err == nil {
		t.Fatal("un stock que es objeto debe reportarse")
	}
}
