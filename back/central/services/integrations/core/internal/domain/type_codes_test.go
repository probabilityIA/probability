package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// ============================================
// IntegrationTypeCodeAsInt
// ============================================

func TestIntegrationTypeCodeAsInt_CodigosConocidos(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected int
	}{
		{
			name:     "shopify retorna 1",
			code:     "shopify",
			expected: IntegrationTypeShopify,
		},
		{
			name:     "whatsapp retorna 2",
			code:     "whatsapp",
			expected: IntegrationTypeWhatsApp,
		},
		{
			name:     "whatsap (variante) retorna 2",
			code:     "whatsap",
			expected: IntegrationTypeWhatsApp,
		},
		{
			name:     "whastap (variante) retorna 2",
			code:     "whastap",
			expected: IntegrationTypeWhatsApp,
		},
		{
			name:     "mercado_libre retorna 3",
			code:     "mercado_libre",
			expected: IntegrationTypeMercadoLibre,
		},
		{
			name:     "mercado libre (con espacio) retorna 3",
			code:     "mercado libre",
			expected: IntegrationTypeMercadoLibre,
		},
		{
			name:     "woocommerce retorna 4",
			code:     "woocommerce",
			expected: IntegrationTypeWoocommerce,
		},
		{
			name:     "woocormerce (variante typo) retorna 4",
			code:     "woocormerce",
			expected: IntegrationTypeWoocommerce,
		},
		{
			name:     "softpymes retorna 5",
			code:     "softpymes",
			expected: IntegrationTypeInvoicing,
		},
		{
			name:     "invoicing retorna 5",
			code:     "invoicing",
			expected: IntegrationTypeInvoicing,
		},
		{
			name:     "platform retorna 6",
			code:     "platform",
			expected: IntegrationTypePlatform,
		},
		{
			name:     "factus retorna 7",
			code:     "factus",
			expected: IntegrationTypeFactus,
		},
		{
			name:     "siigo retorna 8",
			code:     "siigo",
			expected: IntegrationTypeSiigo,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resultado := IntegrationTypeCodeAsInt(tt.code)
			assert.Equal(t, tt.expected, resultado)
		})
	}
}

func TestIntegrationTypeCodeAsInt_CodigoDesconocido(t *testing.T) {
	tests := []struct {
		name string
		code string
	}{
		{name: "código vacío retorna 0", code: ""},
		{name: "código inexistente retorna 0", code: "paypal"},
		{name: "código con mayúsculas es case-insensitive", code: "SHOPIFY"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resultado := IntegrationTypeCodeAsInt(tt.code)
			if tt.code == "SHOPIFY" {
				// SHOPIFY en mayúsculas debe resolverse igual que en minúsculas
				assert.Equal(t, IntegrationTypeShopify, resultado)
			} else {
				assert.Equal(t, 0, resultado)
			}
		})
	}
}

func TestIntegrationTypeCodeAsInt_EsCaseInsensitive(t *testing.T) {
	// La función debe comportarse igual con mayúsculas y minúsculas
	assert.Equal(t, IntegrationTypeCodeAsInt("shopify"), IntegrationTypeCodeAsInt("SHOPIFY"))
	assert.Equal(t, IntegrationTypeCodeAsInt("factus"), IntegrationTypeCodeAsInt("Factus"))
	assert.Equal(t, IntegrationTypeCodeAsInt("siigo"), IntegrationTypeCodeAsInt("SIIGO"))
}

// ============================================
// IsValidCategory
// ============================================

func TestIsValidCategory_CategoriasValidas(t *testing.T) {
	tests := []struct {
		name     string
		category string
		expected bool
	}{
		{
			name:     "internal es válida",
			category: IntegrationCategoryInternal,
			expected: true,
		},
		{
			name:     "external es válida",
			category: IntegrationCategoryExternal,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resultado := IsValidCategory(tt.category)
			assert.Equal(t, tt.expected, resultado)
		})
	}
}

func TestIsValidCategory_CategoriasInvalidas(t *testing.T) {
	tests := []struct {
		name     string
		category string
	}{
		{name: "vacía es inválida", category: ""},
		{name: "ecommerce es inválida (no es standard)", category: "ecommerce"},
		{name: "invoicing es inválida", category: "invoicing"},
		{name: "INTERNAL en mayúsculas es inválida (case-sensitive)", category: "INTERNAL"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resultado := IsValidCategory(tt.category)
			assert.False(t, resultado, "La categoría '%s' no debería ser válida", tt.category)
		})
	}
}

// ============================================
// BaseIntegration (implementación por defecto)
// ============================================

func TestBaseIntegration_TodosLosMetodosRetornanErrNotSupported(t *testing.T) {
	// BaseIntegration provee implementación noop que retorna ErrNotSupported
	base := BaseIntegration{}

	t.Run("TestConnection", func(t *testing.T) {
		err := base.TestConnection(nil, nil, nil)
		assert.ErrorIs(t, err, ErrNotSupported)
	})

	t.Run("SyncOrdersByIntegrationID", func(t *testing.T) {
		err := base.SyncOrdersByIntegrationID(nil, "")
		assert.ErrorIs(t, err, ErrNotSupported)
	})

	t.Run("SyncOrdersByIntegrationIDWithParams", func(t *testing.T) {
		err := base.SyncOrdersByIntegrationIDWithParams(nil, "", nil)
		assert.ErrorIs(t, err, ErrNotSupported)
	})

	t.Run("GetWebhookURL", func(t *testing.T) {
		info, err := base.GetWebhookURL(nil, "", 0)
		assert.ErrorIs(t, err, ErrNotSupported)
		assert.Nil(t, info)
	})

	t.Run("ListWebhooks", func(t *testing.T) {
		lista, err := base.ListWebhooks(nil, "")
		assert.ErrorIs(t, err, ErrNotSupported)
		assert.Nil(t, lista)
	})

	t.Run("DeleteWebhook", func(t *testing.T) {
		err := base.DeleteWebhook(nil, "", "")
		assert.ErrorIs(t, err, ErrNotSupported)
	})

	t.Run("VerifyWebhooksByURL", func(t *testing.T) {
		lista, err := base.VerifyWebhooksByURL(nil, "", "")
		assert.ErrorIs(t, err, ErrNotSupported)
		assert.Nil(t, lista)
	})

	t.Run("CreateWebhook", func(t *testing.T) {
		resultado, err := base.CreateWebhook(nil, "", "")
		assert.ErrorIs(t, err, ErrNotSupported)
		assert.Nil(t, resultado)
	})
}
