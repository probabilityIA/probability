package usecaseshipment

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const guiaCoordinadoraTexto = "COORDINADORA MERCANTIL\n" +
	"AS 4521\n" +
	"PAQ 12-34\n" +
	"UNIDAD: 1/2\n" +
	"Origen\n" +
	"05001\n" +
	"MED\n" +
	"BOG\n" +
	"Destino\n" +
	"Zona Hub\n" +
	"7\n" +
	"Equipo Reparto\n" +
	"215\n" +
	"ref:\n" +
	"ORDEN\n" +
	"MYS-0302\n" +
	"GUIA: 123.456.789\n" +
	"Postal: 110-01\n" +
	"Observaciones Cliente:\n" +
	"Dejar en porteria\n" +
	"con el vigilante\n\n" +
	"FIN\n"

func TestExtractCoordinadoraFields_GuiaCompleta(t *testing.T) {
	got := extractCoordinadoraFields(guiaCoordinadoraTexto)

	assert.Equal(t, "4521", got["as_code"])
	assert.Equal(t, "12-34", got["paq"])
	assert.Equal(t, "1/2", got["unidad"])
	assert.Equal(t, "05001 MED", got["origen"])
	assert.Equal(t, "BOG", got["destino"])
	assert.Equal(t, "7", got["zona_hub"])
	assert.Equal(t, "215", got["equipo_reparto"])
	assert.Equal(t, "MYS-0302", got["ref"])
	assert.Equal(t, "123.456.789", got["guia"])
	assert.Equal(t, "110-01", got["postal_origen"])
	assert.Equal(t, "Dejar en porteria con el vigilante", got["observaciones"])
}

func TestExtractCoordinadoraFields_TextoVacio_RetornaMapaVacio(t *testing.T) {
	got := extractCoordinadoraFields("")

	assert.NotNil(t, got)
	assert.Empty(t, got)
}

func TestExtractCoordinadoraFields_TextoSinCampos_NoInventaValores(t *testing.T) {
	got := extractCoordinadoraFields("una guia de otra transportadora\nsin ninguno de los campos\n")

	assert.Empty(t, got)
}

func TestExtractCoordinadoraFields_CamposParciales_SoloDevuelveLosPresentes(t *testing.T) {
	got := extractCoordinadoraFields("GUIA: 999.111\nPostal: 76-001\n")

	assert.Len(t, got, 2)
	assert.Equal(t, "999.111", got["guia"])
	assert.Equal(t, "76-001", got["postal_origen"])
	_, tieneOrigen := got["origen"]
	assert.False(t, tieneOrigen)
}

func TestExtractCoordinadoraFields_RefSinEtiquetaOrden(t *testing.T) {
	got := extractCoordinadoraFields("ref:\nABC-9090\n")

	assert.Equal(t, "ABC-9090", got["ref"])
}

func TestExtractCoordinadoraFields_RefEsCaseInsensitive(t *testing.T) {
	got := extractCoordinadoraFields("REF:\nXYZ-1\n")

	assert.Equal(t, "XYZ-1", got["ref"])
}

func TestNormalizeExtractedField(t *testing.T) {
	cases := map[string]string{
		"05001\nMED":              "05001 MED",
		"  espacios   sobrantes ": "espacios sobrantes",
		"linea1\n\nlinea2":        "linea1 linea2",
		"simple":                  "simple",
		"":                        "",
		"   ":                     "",
		"\t\ttab\tseparado":       "tab separado",
	}

	for input, want := range cases {
		assert.Equal(t, want, normalizeExtractedField(input))
	}
}

func TestExtractCoordinadoraMetadata_URLVacia_RetornaError(t *testing.T) {
	got, err := ExtractCoordinadoraMetadata(context.Background(), "")

	require.Error(t, err)
	assert.Nil(t, got)
	assert.Contains(t, err.Error(), "label URL is empty")
}

func TestExtractTextFromPDF_BytesInvalidos_RetornaError(t *testing.T) {
	got, err := extractTextFromPDF([]byte("esto no es un pdf"))

	require.Error(t, err)
	assert.Equal(t, "", got)
	assert.Contains(t, err.Error(), "error reading PDF")
}
