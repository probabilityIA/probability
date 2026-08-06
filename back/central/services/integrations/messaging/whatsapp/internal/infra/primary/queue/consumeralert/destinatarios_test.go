package consumeralert

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResolveAlertPhones_SinConfiguracionUsaLosDosNumerosPorDefecto(t *testing.T) {
	phones := resolveAlertPhones("")

	assert.Equal(t, []string{"573023406789", "573192611891"}, phones,
		"si la variable de entorno no esta seteada en prod, las alertas igual llegan a los dos telefonos conocidos")
}

func TestResolveAlertPhones_LeAgregaElIndicativoALosNumerosColombianos(t *testing.T) {
	phones := resolveAlertPhones("3192611891")

	assert.Equal(t, []string{"573192611891"}, phones,
		"un movil colombiano de 10 digitos que empieza en 3 se normaliza con el 57; sin eso Meta rechaza el envio")
}

func TestResolveAlertPhones_AceptaVariosSeparadosPorComa(t *testing.T) {
	phones := resolveAlertPhones("3192611891,573023406789")

	assert.Equal(t, []string{"573192611891", "573023406789"}, phones)
}

func TestResolveAlertPhones_ToleraEspaciosYFormatosSucios(t *testing.T) {
	phones := resolveAlertPhones("  +57 319 261 1891 , 573023406789  ")

	assert.Equal(t, []string{"573192611891", "573023406789"}, phones,
		"la variable la escribe una persona a mano; espacios, signos y el + no deben romper el envio")
}

func TestResolveAlertPhones_DescartaRepetidosAunqueVenganEscritosDistinto(t *testing.T) {
	phones := resolveAlertPhones("3192611891,+573192611891,573192611891")

	assert.Equal(t, []string{"573192611891"}, phones,
		"sin deduplicar, la misma persona recibiria tres WhatsApp por cada alerta")
}

func TestResolveAlertPhones_IgnoraEntradasVacias(t *testing.T) {
	phones := resolveAlertPhones("3192611891,,  ,573023406789")

	assert.Equal(t, []string{"573192611891", "573023406789"}, phones,
		"una coma de mas al final de la variable no debe generar un destinatario vacio")
}

func TestResolveAlertPhones_UnaConfiguracionInservibleCaeALosPorDefecto(t *testing.T) {
	casos := []string{",", "   ", ",,,"}

	for _, raw := range casos {
		phones := resolveAlertPhones(raw)

		assert.Equal(t, defaultAlertPhones, phones,
			"si alguien deja la variable mal escrita, es preferible alertar a los numeros por defecto que quedarse sin alertas")
	}
}

func TestResolveAlertPhones_RespetaElOrdenEnQueSeEscribieron(t *testing.T) {
	phones := resolveAlertPhones("573023406789,3192611891")

	assert.Equal(t, []string{"573023406789", "573192611891"}, phones)
}
