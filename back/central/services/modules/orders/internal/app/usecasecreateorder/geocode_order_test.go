package usecasecreateorder

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGeoConfidence(t *testing.T) {
	cases := []struct {
		name         string
		locationType string
		partialMatch bool
		want         string
	}{
		{name: "rooftop exacto es alta", locationType: "ROOFTOP", partialMatch: false, want: "high"},
		{name: "rooftop con match parcial es baja", locationType: "ROOFTOP", partialMatch: true, want: "low"},
		{name: "interpolado es media", locationType: "RANGE_INTERPOLATED", partialMatch: false, want: "medium"},
		{name: "centro geometrico es media", locationType: "GEOMETRIC_CENTER", partialMatch: false, want: "medium"},
		{name: "aproximado es media", locationType: "APPROXIMATE", partialMatch: false, want: "medium"},
		{name: "vacio es media", locationType: "", partialMatch: false, want: "medium"},
		{name: "match parcial siempre gana", locationType: "APPROXIMATE", partialMatch: true, want: "low"},
		{name: "minuscula no cuenta como rooftop", locationType: "rooftop", partialMatch: false, want: "medium"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, geoConfidence(tc.locationType, tc.partialMatch))
		})
	}
}

func TestStreetSegment(t *testing.T) {
	cases := []struct {
		name   string
		street string
		index  int
		want   string
	}{
		{name: "primer segmento", street: "Calle 1 # 2-3|apto 101|Chapinero", index: 0, want: "Calle 1 # 2-3"},
		{name: "tercer segmento", street: "Calle 1 # 2-3|apto 101|Chapinero", index: 2, want: "Chapinero"},
		{name: "indice fuera de rango", street: "Calle 1 # 2-3", index: 2, want: ""},
		{name: "sin separadores indice cero", street: "Calle 1 # 2-3", index: 0, want: "Calle 1 # 2-3"},
		{name: "recorta espacios", street: "  Calle 1  |  apto 101  ", index: 1, want: "apto 101"},
		{name: "segmento vacio", street: "Calle 1||Chapinero", index: 1, want: ""},
		{name: "cadena vacia", street: "", index: 0, want: ""},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, streetSegment(tc.street, tc.index))
		})
	}
}

func TestBuildGeocodeQuery_SiempreTerminaEnColombia(t *testing.T) {
	got := buildGeocodeQuery("Calle 1 # 2-3", "Bogota", "Cundinamarca")

	assert.Equal(t, "Calle 1 # 2-3, Bogota, Cundinamarca, Colombia", got)
}

func TestBuildGeocodeQuery_UsaSegmentosCeroYDos(t *testing.T) {
	got := buildGeocodeQuery("Calle 1 # 2-3|apto 101|Chapinero", "Bogota", "Cundinamarca")

	assert.Equal(t, "Calle 1 # 2-3, Chapinero, Bogota, Cundinamarca, Colombia", got)
	assert.NotContains(t, got, "apto 101", "el segmento 1 no se usa para geocodificar")
}

func TestBuildGeocodeQuery_SinDatos_RetornaVacio(t *testing.T) {
	cases := []struct {
		name   string
		street string
		city   string
		state  string
	}{
		{name: "todo vacio", street: "", city: "", state: ""},
		{name: "solo espacios", street: "   ", city: "  ", state: " "},
		{name: "separadores vacios", street: "||", city: "", state: ""},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, "", buildGeocodeQuery(tc.street, tc.city, tc.state))
		})
	}
}

func TestBuildGeocodeQuery_DatosParciales(t *testing.T) {
	cases := []struct {
		name   string
		street string
		city   string
		state  string
		want   string
	}{
		{name: "solo ciudad", street: "", city: "Medellin", state: "", want: "Medellin, Colombia"},
		{name: "solo departamento", street: "", city: "", state: "Antioquia", want: "Antioquia, Colombia"},
		{name: "solo calle", street: "Carrera 9", city: "", state: "", want: "Carrera 9, Colombia"},
		{name: "ciudad y departamento", street: "", city: "Cali", state: "Valle", want: "Cali, Valle, Colombia"},
		{name: "recorta espacios", street: "  Carrera 9  ", city: "  Cali  ", state: "", want: "Carrera 9, Cali, Colombia"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, buildGeocodeQuery(tc.street, tc.city, tc.state))
		})
	}
}
