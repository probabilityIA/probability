package statusmapper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetIntegrationTypeID(t *testing.T) {
	cases := map[string]uint{
		"shopify":       1,
		"whatsapp":      2,
		"mercado_libre": 3,
		"mercadolibre":  3,
		"woocommerce":   4,
		"platform":      6,
		"jumpseller":    33,
	}

	for entrada, esperado := range cases {
		t.Run(entrada, func(t *testing.T) {
			assert.Equal(t, esperado, GetIntegrationTypeID(entrada))
		})
	}
}

func TestGetIntegrationTypeID_AliasConErroresDeEscritura(t *testing.T) {
	cases := map[string]uint{
		"whatsap":     2,
		"whastap":     2,
		"woocormerce": 4,
	}

	for entrada, esperado := range cases {
		t.Run(entrada, func(t *testing.T) {
			assert.Equal(t, esperado, GetIntegrationTypeID(entrada),
				"alias tolerado por datos historicos con el nombre mal escrito")
		})
	}
}

func TestGetIntegrationTypeID_DesconocidoRetornaCero(t *testing.T) {
	desconocidos := []string{
		"",
		"amazon",
		"vtex",
		"tiktok",
		"Shopify",
		"SHOPIFY",
		" shopify ",
		"mercado libre",
	}

	for _, entrada := range desconocidos {
		t.Run(entrada, func(t *testing.T) {
			assert.Equal(t, uint(0), GetIntegrationTypeID(entrada))
		})
	}
}

func TestGetIntegrationTypeID_EsCaseSensitive(t *testing.T) {
	assert.NotEqual(t, GetIntegrationTypeID("shopify"), GetIntegrationTypeID("Shopify"),
		"el mapeo distingue mayusculas, un canal con otra capitalizacion no mapea")
}
