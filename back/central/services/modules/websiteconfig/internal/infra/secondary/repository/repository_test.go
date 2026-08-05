package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/secamc93/probability/back/central/shared/testkit"
)

func TestGetConfig_FiltraSiemprePorElNegocio(t *testing.T) {
	base := testkit.NewDB(t)
	base.Mock.ExpectQuery(`SELECT \* FROM "business_website_configs" WHERE business_id = \$1`).
		WithArgs(26, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "business_id", "template"}).
			AddRow(3, 26, "minimal"))

	config, err := New(base).GetConfig(context.Background(), 26)

	require.NoError(t, err)
	assert.EqualValues(t, 26, config.BusinessID)
	base.SinPendientes(t)
}

func TestGetConfig_ElBusinessIDViajaComoArgumentoNoConcatenado(t *testing.T) {
	base := testkit.NewDB(t)
	base.Mock.ExpectQuery(`WHERE business_id = \$1`).
		WithArgs(217, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "business_id"}).AddRow(1, 217))

	_, err := New(base).GetConfig(context.Background(), 217)

	require.NoError(t, err)
	base.SinPendientes(t)
}

func TestGetConfig_SinFilaDevuelveNilSinError(t *testing.T) {
	base := testkit.NewDB(t)
	base.Mock.ExpectQuery(`SELECT \* FROM "business_website_configs"`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))

	config, err := New(base).GetConfig(context.Background(), 26)

	require.NoError(t, err,
		"un negocio sin configuracion no es un error: el caso de uso lo traduce a la plantilla por defecto")
	assert.Nil(t, config)
	base.SinPendientes(t)
}

func TestGetConfig_UnaPlantillaVaciaSeNormalizaADefault(t *testing.T) {
	base := testkit.NewDB(t)
	base.Mock.ExpectQuery(`SELECT \* FROM "business_website_configs"`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "business_id", "template"}).
			AddRow(3, 26, ""))

	config, err := New(base).GetConfig(context.Background(), 26)

	require.NoError(t, err)
	assert.Equal(t, "default", config.Template,
		"filas viejas quedaron con template vacio; se normaliza al leer para que el front nunca reciba cadena vacia")
}

func TestGetConfig_PropagaUnErrorRealDeLaBase(t *testing.T) {
	fallo := errors.New("connection refused")
	base := testkit.NewDB(t)
	base.Mock.ExpectQuery(`SELECT \* FROM "business_website_configs"`).WillReturnError(fallo)

	config, err := New(base).GetConfig(context.Background(), 26)

	require.Error(t, err,
		"solo ErrRecordNotFound se traduce a nil; un fallo de conexion si sube")
	assert.Nil(t, config)
}

func TestGetConfig_TrasladaLosBloquesDeContenidoYLosInterruptores(t *testing.T) {
	base := testkit.NewDB(t)
	base.Mock.ExpectQuery(`SELECT \* FROM "business_website_configs"`).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "business_id", "template", "theme_content", "show_hero", "show_about",
		}).AddRow(3, 26, "minimal", []byte(`{"primary":"#111"}`), true, false))

	config, err := New(base).GetConfig(context.Background(), 26)

	require.NoError(t, err)
	assert.JSONEq(t, `{"primary":"#111"}`, string(config.ThemeContent))
	assert.True(t, config.ShowHero)
	assert.False(t, config.ShowAbout)
}

func TestListProductCategories_FiltraPorNegocioYNoTraeCatalogoAjeno(t *testing.T) {
	base := testkit.NewDB(t)
	base.Mock.ExpectQuery(`business_id = \$1`).
		WithArgs(26).
		WillReturnRows(sqlmock.NewRows([]string{"category", "total"}).
			AddRow("Camisetas", 12).
			AddRow("Pantalones", 4))

	categorias, err := New(base).ListProductCategories(context.Background(), 26)

	require.NoError(t, err)
	require.Len(t, categorias, 2)
	assert.Equal(t, "Camisetas", categorias[0].Name)
	assert.EqualValues(t, 12, categorias[0].ProductCount)
	base.SinPendientes(t)
}

func TestListProductCategories_PropagaElErrorDeLaBase(t *testing.T) {
	base := testkit.NewDB(t)
	base.Mock.ExpectQuery(`business_id`).WillReturnError(errors.New("timeout"))

	_, err := New(base).ListProductCategories(context.Background(), 26)

	assert.Error(t, err)
}
