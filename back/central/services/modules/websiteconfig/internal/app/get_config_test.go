package app

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/secamc93/probability/back/central/services/modules/websiteconfig/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/websiteconfig/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/websiteconfig/internal/mocks"
)

func nuevoCaso(repo *mocks.RepositoryMock) IUseCase {
	if repo == nil {
		repo = &mocks.RepositoryMock{}
	}
	return New(repo, mocks.NewSilentLogger())
}

func TestGetConfig_SinConfiguracionGuardadaDevuelveUnaPlantillaPorDefecto(t *testing.T) {
	repo := &mocks.RepositoryMock{}

	config, err := nuevoCaso(repo).GetConfig(context.Background(), 26)

	require.NoError(t, err)
	require.NotNil(t, config,
		"un negocio recien creado nunca tiene fila en website_configs: la respuesta por defecto evita que el sitio salga en blanco")
	assert.EqualValues(t, 26, config.BusinessID)
	assert.Equal(t, "default", config.Template)
	assert.Zero(t, config.ID, "el default no esta persistido, por eso el ID es cero")
}

func TestGetConfig_ElDefaultSoloEnciendeLasSeccionesMinimasDeUnaTienda(t *testing.T) {
	config, err := nuevoCaso(nil).GetConfig(context.Background(), 26)

	require.NoError(t, err)
	assert.True(t, config.ShowHero)
	assert.True(t, config.ShowFeaturedProducts)
	assert.True(t, config.ShowFullCatalog)
	assert.True(t, config.ShowContact)

	assert.False(t, config.ShowAbout,
		"las secciones que necesitan contenido escrito por el negocio arrancan apagadas para no mostrar plantillas vacias")
	assert.False(t, config.ShowTestimonials)
	assert.False(t, config.ShowLocation)
	assert.False(t, config.ShowSocialMedia)
	assert.False(t, config.ShowWhatsApp)
}

func TestGetConfig_DevuelveLaConfiguracionGuardadaSinTocarla(t *testing.T) {
	guardada := &entities.WebsiteConfig{
		ID: 3, BusinessID: 26, Template: "minimal",
		ShowHero: false, ShowTestimonials: true,
		ThemeContent: []byte(`{"primary":"#111"}`),
	}
	repo := &mocks.RepositoryMock{
		GetConfigFn: func(ctx context.Context, businessID uint) (*entities.WebsiteConfig, error) {
			return guardada, nil
		},
	}

	config, err := nuevoCaso(repo).GetConfig(context.Background(), 26)

	require.NoError(t, err)
	assert.Equal(t, "minimal", config.Template)
	assert.False(t, config.ShowHero,
		"si el negocio apago el hero a proposito, el default no debe volver a encenderlo")
	assert.True(t, config.ShowTestimonials)
	assert.JSONEq(t, `{"primary":"#111"}`, string(config.ThemeContent))
}

func TestGetConfig_ConsultaConElNegocioQueRecibe(t *testing.T) {
	repo := &mocks.RepositoryMock{}

	_, err := nuevoCaso(repo).GetConfig(context.Background(), 217)

	require.NoError(t, err)
	assert.Equal(t, []uint{217}, repo.ConsultasConfig,
		"el business_id resuelto en el handler es el unico filtro: no hay forma de leer la config de otro negocio desde aqui")
}

func TestGetConfig_PropagaElErrorDelRepositorio(t *testing.T) {
	fallo := errors.New("timeout")
	repo := &mocks.RepositoryMock{
		GetConfigFn: func(ctx context.Context, businessID uint) (*entities.WebsiteConfig, error) {
			return nil, fallo
		},
	}

	config, err := nuevoCaso(repo).GetConfig(context.Background(), 26)

	assert.ErrorIs(t, err, fallo)
	assert.Nil(t, config,
		"un fallo de base no se disfraza de configuracion por defecto: el sitio no debe repintarse solo porque postgres tardo")
}

func TestUpdateConfig_DelegaElUpsertConElNegocioYElDTO(t *testing.T) {
	repo := &mocks.RepositoryMock{}
	plantilla := "minimal"
	visible := true
	dto := &dtos.UpdateConfigDTO{Template: &plantilla, ShowTestimonials: &visible}

	config, err := nuevoCaso(repo).UpdateConfig(context.Background(), 26, dto)

	require.NoError(t, err)
	require.Len(t, repo.Upserts, 1)
	assert.EqualValues(t, 26, repo.Upserts[0].BusinessID)
	assert.Same(t, dto, repo.Upserts[0].DTO)
	assert.EqualValues(t, 26, config.BusinessID)
}

func TestUpdateConfig_NoValidaNadaDelDTOAntesDeEscribir(t *testing.T) {
	repo := &mocks.RepositoryMock{}
	plantilla := "plantilla-que-no-existe"

	_, err := nuevoCaso(repo).UpdateConfig(context.Background(), 26, &dtos.UpdateConfigDTO{
		Template:    &plantilla,
		HeroContent: []byte("esto no es json"),
	})

	require.NoError(t, err,
		"el caso de uso no valida el nombre de plantilla ni que los bloques de contenido sean JSON: eso queda en el repositorio y en el front")
	assert.Len(t, repo.Upserts, 1)
}

func TestUpdateConfig_UnDTOVacioLlegaIgualAlRepositorio(t *testing.T) {
	repo := &mocks.RepositoryMock{}

	_, err := nuevoCaso(repo).UpdateConfig(context.Background(), 26, &dtos.UpdateConfigDTO{})

	require.NoError(t, err)
	require.Len(t, repo.Upserts, 1)
	assert.Nil(t, repo.Upserts[0].DTO.Template,
		"los punteros nil son la senal de 'no cambiar este campo'; el upsert parcial se resuelve en el repositorio")
}

func TestUpdateConfig_PropagaElErrorDelRepositorio(t *testing.T) {
	fallo := errors.New("insert fallido")
	repo := &mocks.RepositoryMock{
		UpsertConfigFn: func(ctx context.Context, businessID uint, dto *dtos.UpdateConfigDTO) (*entities.WebsiteConfig, error) {
			return nil, fallo
		},
	}

	config, err := nuevoCaso(repo).UpdateConfig(context.Background(), 26, &dtos.UpdateConfigDTO{})

	assert.ErrorIs(t, err, fallo)
	assert.Nil(t, config)
}

func TestGetProductCategories_DevuelveLasCategoriasConSuConteo(t *testing.T) {
	repo := &mocks.RepositoryMock{
		ListProductCategoriesFn: func(ctx context.Context, businessID uint) ([]entities.ProductCategory, error) {
			return []entities.ProductCategory{
				{Name: "Camisetas", ProductCount: 12},
				{Name: "Pantalones", ProductCount: 4},
			}, nil
		},
	}

	categorias, err := nuevoCaso(repo).GetProductCategories(context.Background(), 26)

	require.NoError(t, err)
	require.Len(t, categorias, 2)
	assert.Equal(t, "Camisetas", categorias[0].Name)
	assert.EqualValues(t, 12, categorias[0].ProductCount,
		"el conteo alimenta el selector de categorias ocultas en el editor del sitio")
}

func TestGetProductCategories_ConsultaSoloElNegocioRecibido(t *testing.T) {
	repo := &mocks.RepositoryMock{}

	_, err := nuevoCaso(repo).GetProductCategories(context.Background(), 217)

	require.NoError(t, err)
	assert.Equal(t, []uint{217}, repo.ConsultasCategorias)
}

func TestGetProductCategories_UnCatalogoVacioNoEsError(t *testing.T) {
	categorias, err := nuevoCaso(nil).GetProductCategories(context.Background(), 26)

	require.NoError(t, err)
	assert.Empty(t, categorias)
}

func TestGetProductCategories_PropagaElErrorDelRepositorio(t *testing.T) {
	fallo := errors.New("timeout")
	repo := &mocks.RepositoryMock{
		ListProductCategoriesFn: func(ctx context.Context, businessID uint) ([]entities.ProductCategory, error) {
			return nil, fallo
		},
	}

	_, err := nuevoCaso(repo).GetProductCategories(context.Background(), 26)

	assert.ErrorIs(t, err, fallo)
}
