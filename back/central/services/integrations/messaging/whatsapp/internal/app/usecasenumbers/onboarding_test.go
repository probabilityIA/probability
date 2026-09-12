package usecasenumbers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidarNumeroAceptaUnCelularColombiano(t *testing.T) {
	assert.NoError(t, validarNumero("57", "3001234567"))
}

func TestValidarNumeroRechazaElCelularColombianoConUnDigitoDeMas(t *testing.T) {
	err := validarNumero("57", "30279400815")
	require.Error(t, err, "es el numero que Meta acepto y luego no pudo verificar")
	assert.Contains(t, err.Error(), "10")
}

func TestValidarNumeroRechazaElCelularColombianoCorto(t *testing.T) {
	require.Error(t, validarNumero("57", "300123456"))
}

func TestValidarNumeroRechazaElFijoColombianoSinIndicativoMovil(t *testing.T) {
	require.Error(t, validarNumero("57", "6014567890"))
}

func TestValidarNumeroExigeIndicativoYNumero(t *testing.T) {
	require.Error(t, validarNumero("", "3001234567"))
	require.Error(t, validarNumero("57", ""))
}

func TestValidarNumeroRechazaMasDeQuinceDigitos(t *testing.T) {
	err := validarNumero("52", "12345678901234")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "15")
}

func TestValidarNumeroNoImponeLaReglaColombianaAOtrosPaises(t *testing.T) {
	assert.NoError(t, validarNumero("52", "5512345678"))
	assert.NoError(t, validarNumero("1", "4155551234"))
}
