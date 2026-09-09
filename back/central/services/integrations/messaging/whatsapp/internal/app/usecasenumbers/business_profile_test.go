package usecasenumbers

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func texto(v string) *string { return &v }

func TestValidarPerfilAceptaUnPerfilCompleto(t *testing.T) {
	sitios := []string{"https://probabilityia.com.co"}
	err := validarPerfil(BusinessProfileInput{
		About:       texto("Envios a todo el pais"),
		Description: texto("Tienda de accesorios"),
		Email:       texto("hola@tienda.com"),
		Vertical:    texto("RETAIL"),
		Websites:    &sitios,
	})
	assert.NoError(t, err)
}

func TestValidarPerfilRechazaElAboutMuyLargo(t *testing.T) {
	err := validarPerfil(BusinessProfileInput{About: texto(strings.Repeat("a", maxAbout+1))})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "139")
}

func TestValidarPerfilRechazaLaDescripcionMuyLarga(t *testing.T) {
	err := validarPerfil(BusinessProfileInput{Description: texto(strings.Repeat("a", maxDescription+1))})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "512")
}

func TestValidarPerfilCuentaCaracteresNoBytes(t *testing.T) {
	err := validarPerfil(BusinessProfileInput{About: texto(strings.Repeat("ñ", maxAbout))})
	assert.NoError(t, err, "una ñ ocupa dos bytes: contar bytes recortaria textos validos")
}

func TestValidarPerfilRechazaMasDeDosSitios(t *testing.T) {
	sitios := []string{"https://a.com", "https://b.com", "https://c.com"}
	err := validarPerfil(BusinessProfileInput{Websites: &sitios})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "2")
}

func TestValidarPerfilExigeEsquemaEnElSitio(t *testing.T) {
	sitios := []string{"probabilityia.com.co"}
	err := validarPerfil(BusinessProfileInput{Websites: &sitios})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "http")
}

func TestValidarPerfilIgnoraSitiosVacios(t *testing.T) {
	sitios := []string{"https://a.com", ""}
	assert.NoError(t, validarPerfil(BusinessProfileInput{Websites: &sitios}))
}

func TestValidarPerfilRechazaCorreoSinArroba(t *testing.T) {
	err := validarPerfil(BusinessProfileInput{Email: texto("hola.tienda.com")})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "correo")
}

func TestValidarPerfilRechazaCategoriaInventada(t *testing.T) {
	err := validarPerfil(BusinessProfileInput{Vertical: texto("FERRETERIA")})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "categoria")
}

func TestValidarPerfilAceptaCategoriaEnMinusculas(t *testing.T) {
	assert.NoError(t, validarPerfil(BusinessProfileInput{Vertical: texto("retail")}))
}
