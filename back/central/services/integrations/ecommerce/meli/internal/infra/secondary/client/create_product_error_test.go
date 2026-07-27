package client

import "testing"

func TestDescribeMeliErrorTraduceCausaConocida(t *testing.T) {
	body := []byte(`{"message":"Validation error","error":"validation_error","status":400,"cause":[{"department":"structured-data","cause_id":2610,"type":"error","code":"missing.fashion_grid.grid_id.values","references":["item.attributes"],"message":"Attribute GRID_ID is required"}]}`)

	got := describeMeliError(400, body)
	want := "falta la guia de talles que Mercado Libre exige en categorias de moda (atributo GRID_ID)"
	if got != want {
		t.Errorf("describeMeliError() = %q, se esperaba %q", got, want)
	}
}

func TestDescribeMeliErrorUsaElMensajeDeLaCausaDesconocida(t *testing.T) {
	body := []byte(`{"status":400,"cause":[{"code":"item.price.something","message":"El precio debe ser mayor a cero"}]}`)

	if got := describeMeliError(400, body); got != "El precio debe ser mayor a cero" {
		t.Errorf("describeMeliError() = %q", got)
	}
}

func TestDescribeMeliErrorCombinaHastaDosCausas(t *testing.T) {
	body := []byte(`{"status":400,"cause":[{"code":"a","message":"primera"},{"code":"b","message":"segunda"},{"code":"c","message":"tercera"}]}`)

	if got := describeMeliError(400, body); got != "primera | segunda" {
		t.Errorf("describeMeliError() = %q", got)
	}
}

func TestDescribeMeliErrorSinCausaUsaMessage(t *testing.T) {
	body := []byte(`{"message":"Item not published","status":403,"cause":[]}`)

	if got := describeMeliError(403, body); got != "Item not published" {
		t.Errorf("describeMeliError() = %q", got)
	}
}

func TestDescribeMeliErrorConBodyNoJSON(t *testing.T) {
	got := describeMeliError(502, []byte("bad gateway"))
	if got != "Mercado Libre respondio 502: bad gateway" {
		t.Errorf("describeMeliError() = %q", got)
	}
}
