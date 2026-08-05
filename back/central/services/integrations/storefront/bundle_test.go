package storefront

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/secamc93/probability/back/central/services/integrations/core"
)

func TestNew_CumpleElContratoDeIntegracion(t *testing.T) {
	var _ core.IIntegrationContract = New()
}

func TestTestConnection_SiempreAprueba(t *testing.T) {
	p := New()

	assert.NoError(t, p.TestConnection(context.Background(), nil, nil),
		"no hay API externa que probar: el proveedor es solo un interruptor on/off, por eso siempre aprueba")
	assert.NoError(t, p.TestConnection(context.Background(),
		map[string]interface{}{"lo que sea": 1},
		map[string]interface{}{"credencial": "invalida"}),
		"ni siquiera mira config ni credenciales")
}

func TestElRestoDelContratoDevuelveNoSoportado(t *testing.T) {
	p := New()
	ctx := context.Background()

	assert.ErrorIs(t, p.SyncOrdersByIntegrationID(ctx, "1"), core.ErrNotSupported)
	assert.ErrorIs(t, p.SyncOrdersByIntegrationIDWithParams(ctx, "1", nil), core.ErrNotSupported)
	assert.ErrorIs(t, p.DeleteWebhook(ctx, "1", "w1"), core.ErrNotSupported)
	assert.ErrorIs(t, p.UpdateInventory(ctx, "1", "SKU-1", 5), core.ErrNotSupported)

	_, err := p.GetWebhookURL(ctx, "https://api.probability.com", 1)
	assert.ErrorIs(t, err, core.ErrNotSupported)
	_, err = p.ListWebhooks(ctx, "1")
	assert.ErrorIs(t, err, core.ErrNotSupported)
	_, err = p.VerifyWebhooksByURL(ctx, "1", "https://api.probability.com")
	assert.ErrorIs(t, err, core.ErrNotSupported)
	_, err = p.CreateWebhook(ctx, "1", "https://api.probability.com")
	require.ErrorIs(t, err, core.ErrNotSupported,
		"BaseIntegration cubre todo lo que este proveedor no hace, en vez de dejar metodos sin implementar")
}
