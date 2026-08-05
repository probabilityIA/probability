package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/datatypes"
)

func TestCanonicalizeVariantAttributes_Vacio(t *testing.T) {
	got, clave, err := CanonicalizeVariantAttributes(nil)

	require.NoError(t, err)
	assert.Nil(t, got)
	assert.Equal(t, "", clave)
}

func TestCanonicalizeVariantAttributes_ObjetoVacio(t *testing.T) {
	got, clave, err := CanonicalizeVariantAttributes(datatypes.JSON([]byte(`{}`)))

	require.NoError(t, err)
	assert.JSONEq(t, `{}`, string(got))
	assert.Equal(t, "", clave)
}

func TestCanonicalizeVariantAttributes_JSONInvalido_RetornaError(t *testing.T) {
	casos := []string{
		`{roto`,
		`no es json`,
		`[1,2,3]`,
		`"solo un string"`,
	}

	for _, entrada := range casos {
		t.Run(entrada, func(t *testing.T) {
			got, clave, err := CanonicalizeVariantAttributes(datatypes.JSON([]byte(entrada)))

			require.Error(t, err)
			assert.Contains(t, err.Error(), "variant_attributes must be a valid JSON object")
			assert.Nil(t, got)
			assert.Equal(t, "", clave)
		})
	}
}

func TestCanonicalizeVariantAttributes_OrdenDeClavesNoImporta(t *testing.T) {
	_, claveA, err := CanonicalizeVariantAttributes(datatypes.JSON([]byte(`{"talla":"M","color":"rojo"}`)))
	require.NoError(t, err)

	_, claveB, err := CanonicalizeVariantAttributes(datatypes.JSON([]byte(`{"color":"rojo","talla":"M"}`)))
	require.NoError(t, err)

	assert.Equal(t, claveA, claveB,
		"la misma variante con las claves en otro orden produce la misma clave canonica")
}

func TestCanonicalizeVariantAttributes_ClavesSeNormalizanAMinusculas(t *testing.T) {
	_, claveA, err := CanonicalizeVariantAttributes(datatypes.JSON([]byte(`{"Color":"rojo"}`)))
	require.NoError(t, err)

	_, claveB, err := CanonicalizeVariantAttributes(datatypes.JSON([]byte(`{"COLOR":"rojo"}`)))
	require.NoError(t, err)

	_, claveC, err := CanonicalizeVariantAttributes(datatypes.JSON([]byte(`{"color":"rojo"}`)))
	require.NoError(t, err)

	assert.Equal(t, claveA, claveB)
	assert.Equal(t, claveB, claveC)
	assert.Contains(t, claveC, `"color"`)
}

func TestCanonicalizeVariantAttributes_ClavesConEspacios(t *testing.T) {
	_, clave, err := CanonicalizeVariantAttributes(datatypes.JSON([]byte(`{"  color  ":"rojo"}`)))

	require.NoError(t, err)
	assert.Contains(t, clave, `"color"`)
	assert.NotContains(t, clave, `"  color  "`)
}

func TestCanonicalizeVariantAttributes_ValoresConEspaciosSeRecortan(t *testing.T) {
	_, claveA, err := CanonicalizeVariantAttributes(datatypes.JSON([]byte(`{"color":"  rojo  "}`)))
	require.NoError(t, err)

	_, claveB, err := CanonicalizeVariantAttributes(datatypes.JSON([]byte(`{"color":"rojo"}`)))
	require.NoError(t, err)

	assert.Equal(t, claveA, claveB)
}

func TestCanonicalizeVariantAttributes_ValoresNoSeNormalizanAMinusculas(t *testing.T) {
	_, claveA, err := CanonicalizeVariantAttributes(datatypes.JSON([]byte(`{"color":"Rojo"}`)))
	require.NoError(t, err)

	_, claveB, err := CanonicalizeVariantAttributes(datatypes.JSON([]byte(`{"color":"rojo"}`)))
	require.NoError(t, err)

	assert.NotEqual(t, claveA, claveB,
		"solo las claves se pasan a minusculas, los valores conservan su capitalizacion")
}

func TestCanonicalizeVariantAttributes_VariantesDistintasProducenClavesDistintas(t *testing.T) {
	_, claveA, err := CanonicalizeVariantAttributes(datatypes.JSON([]byte(`{"talla":"M","color":"rojo"}`)))
	require.NoError(t, err)

	_, claveB, err := CanonicalizeVariantAttributes(datatypes.JSON([]byte(`{"talla":"L","color":"rojo"}`)))
	require.NoError(t, err)

	assert.NotEqual(t, claveA, claveB)
}

func TestCanonicalizeVariantAttributes_ObjetosAnidados(t *testing.T) {
	_, claveA, err := CanonicalizeVariantAttributes(datatypes.JSON([]byte(
		`{"medidas":{"ancho":"10","alto":"20"},"color":"rojo"}`)))
	require.NoError(t, err)

	_, claveB, err := CanonicalizeVariantAttributes(datatypes.JSON([]byte(
		`{"color":"rojo","medidas":{"alto":"20","ancho":"10"}}`)))
	require.NoError(t, err)

	assert.Equal(t, claveA, claveB, "los objetos anidados tambien se ordenan por clave")
}

func TestCanonicalizeVariantAttributes_ArreglosConservanElOrden(t *testing.T) {
	_, claveA, err := CanonicalizeVariantAttributes(datatypes.JSON([]byte(`{"tags":["a","b"]}`)))
	require.NoError(t, err)

	_, claveB, err := CanonicalizeVariantAttributes(datatypes.JSON([]byte(`{"tags":["b","a"]}`)))
	require.NoError(t, err)

	assert.NotEqual(t, claveA, claveB,
		"el orden de un arreglo si es significativo, no se ordena")
}

func TestCanonicalizeVariantAttributes_ArregloDeObjetosSeNormaliza(t *testing.T) {
	_, claveA, err := CanonicalizeVariantAttributes(datatypes.JSON([]byte(
		`{"opciones":[{"Nombre":"  talla  ","Valor":"M"}]}`)))
	require.NoError(t, err)

	_, claveB, err := CanonicalizeVariantAttributes(datatypes.JSON([]byte(
		`{"opciones":[{"valor":"M","nombre":"talla"}]}`)))
	require.NoError(t, err)

	assert.Equal(t, claveA, claveB)
}

func TestCanonicalizeVariantAttributes_NumerosNoPierdenPrecision(t *testing.T) {
	_, clave, err := CanonicalizeVariantAttributes(datatypes.JSON([]byte(
		`{"peso":12345678901234567890}`)))

	require.NoError(t, err)
	assert.Contains(t, clave, "12345678901234567890",
		"los numeros grandes no se degradan a notacion cientifica")
}

func TestCanonicalizeVariantAttributes_DecimalesSeConservan(t *testing.T) {
	_, clave, err := CanonicalizeVariantAttributes(datatypes.JSON([]byte(`{"peso":1.50}`)))

	require.NoError(t, err)
	assert.Contains(t, clave, "1.50")
}

func TestCanonicalizeVariantAttributes_TiposMixtos(t *testing.T) {
	_, clave, err := CanonicalizeVariantAttributes(datatypes.JSON([]byte(
		`{"activo":true,"stock":null,"nombre":"camisa","precio":25000}`)))

	require.NoError(t, err)
	assert.Contains(t, clave, `"activo":true`)
	assert.Contains(t, clave, `"stock":null`)
	assert.Contains(t, clave, `"nombre":"camisa"`)
	assert.Contains(t, clave, `"precio":25000`)
}

func TestCanonicalizeVariantAttributes_LaClaveCoincideConElJSONDevuelto(t *testing.T) {
	got, clave, err := CanonicalizeVariantAttributes(datatypes.JSON([]byte(`{"talla":"M","color":"rojo"}`)))

	require.NoError(t, err)
	assert.Equal(t, clave, string(got),
		"la clave canonica es exactamente el JSON normalizado")
}

func TestCanonicalizeVariantAttributes_EsIdempotente(t *testing.T) {
	primera, claveA, err := CanonicalizeVariantAttributes(datatypes.JSON([]byte(`{"Talla":"  M  ","color":"rojo"}`)))
	require.NoError(t, err)

	_, claveB, err := CanonicalizeVariantAttributes(primera)
	require.NoError(t, err)

	assert.Equal(t, claveA, claveB, "canonicalizar dos veces da el mismo resultado")
}

func TestNormalizeValue_TiposBasicos(t *testing.T) {
	assert.Equal(t, "texto", normalizeValue("  texto  "))
	assert.Equal(t, true, normalizeValue(true))
	assert.Nil(t, normalizeValue(nil))

	anidado := normalizeValue(map[string]interface{}{"  A  ": "  b  "})
	assert.Equal(t, map[string]interface{}{"a": "b"}, anidado)

	arreglo := normalizeValue([]interface{}{"  x  ", "  y  "})
	assert.Equal(t, []interface{}{"x", "y"}, arreglo)
}

func TestNormalizeMap_OrdenaYNormalizaClaves(t *testing.T) {
	got := normalizeMap(map[string]interface{}{
		"  Zeta ": "  ultimo  ",
		"Alfa":    "primero",
	})

	assert.Equal(t, map[string]interface{}{"zeta": "ultimo", "alfa": "primero"}, got)
}
