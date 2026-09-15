package navigation

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/secamc93/probability/back/central/services/modules/ai/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/shared/authz"
)

func entriesFor(keys ...string) []authz.NavEntry {
	var out []authz.NavEntry
	for _, key := range keys {
		for _, entry := range authz.Navigation {
			if entry.Key == key {
				out = append(out, entry)
			}
		}
	}
	return out
}

func TestForUser_IncluyeSubpaginasDeLosModulosVisibles(t *testing.T) {
	catalog := New(func(context.Context, uint, uint, uint) ([]authz.NavEntry, error) {
		return entriesFor("orders", "shipments"), nil
	})

	got, err := catalog.ForUser(context.Background(), dtos.AccessScope{UserID: 1, TokenBusinessID: 26})

	require.NoError(t, err)
	cod, ok := got.Find("shipments.cod")
	require.True(t, ok)
	assert.Equal(t, "/shipments/cod", cod.Route)
	assert.NotEmpty(t, cod.Guide)
	_, ok = got.Find("wallet.finanzas")
	assert.False(t, ok, "una subpagina de un modulo no visible no se ofrece")
}

func TestForUser_TusIntegracionesRequiereOrdenes(t *testing.T) {
	withOrders := New(func(context.Context, uint, uint, uint) ([]authz.NavEntry, error) {
		return entriesFor("orders", "integrations"), nil
	})
	got, err := withOrders.ForUser(context.Background(), dtos.AccessScope{UserID: 1, TokenBusinessID: 26})
	require.NoError(t, err)
	hub, ok := got.Find("integrations.hub.inventory")
	require.True(t, ok)
	assert.Equal(t, "/orders", hub.Route)

	withoutOrders := New(func(context.Context, uint, uint, uint) ([]authz.NavEntry, error) {
		return entriesFor("integrations"), nil
	})
	got, err = withoutOrders.ForUser(context.Background(), dtos.AccessScope{UserID: 1, TokenBusinessID: 26})
	require.NoError(t, err)
	_, ok = got.Find("integrations.hub")
	assert.False(t, ok, "el hub vive en la barra de Ordenes; sin acceso a Ordenes no se ofrece")
}

func TestForUser_NoRevelaModulosDeSuperAdmin(t *testing.T) {
	catalog := New(func(context.Context, uint, uint, uint) ([]authz.NavEntry, error) {
		return entriesFor("home"), nil
	})

	got, err := catalog.ForUser(context.Background(), dtos.AccessScope{UserID: 1, TokenBusinessID: 26})

	require.NoError(t, err)
	assert.Contains(t, got.Denied, "Billetera")
	assert.NotContains(t, got.Denied, "Tickets")
	assert.NotContains(t, got.Denied, "Contabilidad")
}

func TestGuide_TodaEntradaDeNavegacionTieneGuia(t *testing.T) {
	guides := parseGuide(guideSource)
	for _, entry := range authz.Navigation {
		if entry.Rule == authz.NavRuleSuperAdminOnly {
			continue
		}
		assert.NotEmpty(t, guides[entry.Key], "falta la seccion ## %s en guide.md", entry.Key)
	}
	for _, sub := range subpages {
		assert.NotEmpty(t, guides[sub.Key], "falta la seccion ## %s en guide.md", sub.Key)
	}
}
