package response

import "testing"

func TestLocationDecodificaCiudadEscapada(t *testing.T) {
	env := LocationEnvelope{}
	env.Location.ID = 327870
	env.Location.Name = "sebas y corotos"
	env.Location.Main = true
	env.Location.IsStockOrigin = true
	env.Location.LocationAddress.City = "Bogot%C3%A1"
	env.Location.LocationAddress.Country = "CO"

	got := env.ToDomain()

	if got.City != "Bogot\u00e1" {
		t.Fatalf("city = %q: Jumpseller la manda URL-encodeada y hay que decodificarla", got.City)
	}
	if got.Country != "CO" {
		t.Fatalf("country = %q", got.Country)
	}
	if !got.Main || !got.IsStockOrigin {
		t.Fatal("se perdieron los flags main / is_stock_origin")
	}
}

func TestLocationSinEscapesQuedaIgual(t *testing.T) {
	env := LocationEnvelope{}
	env.Location.Name = "Bodega Principal"
	env.Location.LocationAddress.City = "Medellin"

	got := env.ToDomain()

	if got.Name != "Bodega Principal" || got.City != "Medellin" {
		t.Fatalf("un valor sin escapes no debe alterarse: %q / %q", got.Name, got.City)
	}
}
