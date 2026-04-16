package statusmapper

import "testing"

func TestEqualJSON_MismoContenido_RetornaTrue(t *testing.T) {
	// Arrange
	a := []byte(`{"key":"value","num":42}`)
	b := []byte(`{"num":42,"key":"value"}`) // Mismo contenido, diferente orden

	// Act
	result := EqualJSON(a, b)

	// Assert
	if !result {
		t.Error("se esperaba true para JSONs con mismo contenido pero diferente orden")
	}
}

func TestEqualJSON_ContenidoDiferente_RetornaFalse(t *testing.T) {
	// Arrange
	a := []byte(`{"key":"value1"}`)
	b := []byte(`{"key":"value2"}`)

	// Act
	result := EqualJSON(a, b)

	// Assert
	if result {
		t.Error("se esperaba false para JSONs con contenido diferente")
	}
}

func TestEqualJSON_JSONInvalido_RetornaFalse(t *testing.T) {
	// Arrange
	a := []byte(`{invalid json}`)
	b := []byte(`{"key":"value"}`)

	// Act
	result := EqualJSON(a, b)

	// Assert
	if result {
		t.Error("se esperaba false cuando el primer JSON es inválido")
	}
}
