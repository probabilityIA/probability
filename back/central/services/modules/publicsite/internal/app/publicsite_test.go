package app

import (
	"context"
	stderrors "errors"
	"strings"
	"testing"

	"github.com/secamc93/probability/back/central/services/modules/publicsite/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/publicsite/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/publicsite/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/publicsite/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newPublicsiteUseCase(repo *mocks.RepositoryMock, bold *mocks.BoldGatewayMock) IUseCase {
	if repo == nil {
		repo = &mocks.RepositoryMock{}
	}
	if bold == nil {
		bold = &mocks.BoldGatewayMock{}
	}
	return New(repo, bold, mocks.NewSilentLogger())
}

func repoSinNegocio() *mocks.RepositoryMock {
	return &mocks.RepositoryMock{
		GetBusinessBySlugFn: func(ctx context.Context, slug string) (*entities.BusinessPage, error) {
			return nil, nil
		},
	}
}

func repoSitioApagado() *mocks.RepositoryMock {
	return &mocks.RepositoryMock{
		IsIntegrationActiveFn: func(ctx context.Context, businessID, typeID uint) (bool, error) {
			return false, nil
		},
	}
}

func carritoValido() *dtos.CreateCheckoutDTO {
	return &dtos.CreateCheckoutDTO{
		Items:        []dtos.CheckoutItemInput{{ProductID: "P-1", Quantity: 2}},
		CustomerName: "Ana", CustomerEmail: "ana@x.com",
	}
}

func TestSlugInexistente_BloqueaTodasLasOperaciones(t *testing.T) {
	uc := newPublicsiteUseCase(repoSinNegocio(), nil)
	ctx := context.Background()

	_, err := uc.GetBusinessPage(ctx, "inexistente")
	assert.ErrorIs(t, err, domainerrors.ErrBusinessNotFound)

	_, _, err = uc.ListCatalog(ctx, "inexistente", dtos.CatalogFilters{})
	assert.ErrorIs(t, err, domainerrors.ErrBusinessNotFound)

	_, err = uc.GetProduct(ctx, "inexistente", "P-1")
	assert.ErrorIs(t, err, domainerrors.ErrBusinessNotFound)

	_, err = uc.GetFeaturedProducts(ctx, "inexistente", 8)
	assert.ErrorIs(t, err, domainerrors.ErrBusinessNotFound)

	_, err = uc.GetCategories(ctx, "inexistente")
	assert.ErrorIs(t, err, domainerrors.ErrBusinessNotFound)

	err = uc.SubmitContact(ctx, "inexistente", &dtos.ContactFormDTO{Name: "a", Message: "m"})
	assert.ErrorIs(t, err, domainerrors.ErrBusinessNotFound)

	_, err = uc.CreateCheckoutSession(ctx, "inexistente", carritoValido(), 0)
	assert.ErrorIs(t, err, domainerrors.ErrBusinessNotFound)

	_, err = uc.GetSession(ctx, "inexistente", 42)
	assert.ErrorIs(t, err, domainerrors.ErrBusinessNotFound)

	_, err = uc.GetMyPrices(ctx, "inexistente", 42)
	assert.ErrorIs(t, err, domainerrors.ErrBusinessNotFound)

	_, _, err = uc.GetMyOrders(ctx, "inexistente", 42, 1, 10)
	assert.ErrorIs(t, err, domainerrors.ErrBusinessNotFound)

	assert.ErrorIs(t, uc.UpdateProfile(ctx, "inexistente", 42, "n", "p", "d"),
		domainerrors.ErrBusinessNotFound)
}

func TestSitioApagado_BloqueaLasOperacionesPublicas(t *testing.T) {
	uc := newPublicsiteUseCase(repoSitioApagado(), nil)
	ctx := context.Background()

	_, err := uc.GetBusinessPage(ctx, "demo")
	assert.ErrorIs(t, err, domainerrors.ErrPublicSiteNotActive)

	_, _, err = uc.ListCatalog(ctx, "demo", dtos.CatalogFilters{})
	assert.ErrorIs(t, err, domainerrors.ErrPublicSiteNotActive)

	_, err = uc.GetProduct(ctx, "demo", "P-1")
	assert.ErrorIs(t, err, domainerrors.ErrPublicSiteNotActive)

	_, err = uc.GetFeaturedProducts(ctx, "demo", 8)
	assert.ErrorIs(t, err, domainerrors.ErrPublicSiteNotActive)

	_, err = uc.GetCategories(ctx, "demo")
	assert.ErrorIs(t, err, domainerrors.ErrPublicSiteNotActive)

	err = uc.SubmitContact(ctx, "demo", &dtos.ContactFormDTO{Name: "a", Message: "m"})
	assert.ErrorIs(t, err, domainerrors.ErrPublicSiteNotActive)
}

func TestSitioApagado_TambienBloqueaElCheckout(t *testing.T) {
	repo := repoSitioApagado()
	bold := &mocks.BoldGatewayMock{}

	_, err := newPublicsiteUseCase(repo, bold).
		CreateCheckoutSession(context.Background(), "demo", carritoValido(), 0)

	assert.ErrorIs(t, err, domainerrors.ErrPublicSiteNotActive)
	assert.Nil(t, repo.CreatedCheckout, "un sitio apagado no debe poder generar ordenes")
	assert.Empty(t, bold.PublishedRefs)
}

func TestElGateUsaElTipoDeIntegracionDeSitioWeb(t *testing.T) {
	repo := &mocks.RepositoryMock{}
	uc := newPublicsiteUseCase(repo, nil)
	ctx := context.Background()

	_, _ = uc.GetBusinessPage(ctx, "demo")
	_, _, _ = uc.ListCatalog(ctx, "demo", dtos.CatalogFilters{})
	_, _ = uc.GetProduct(ctx, "demo", "P-1")

	require.Len(t, repo.GateChecks, 3)
	for i, tipo := range repo.GateChecks {
		assert.Equal(t, uint(tiendaWebIntegrationTypeID), tipo, "consulta %d", i)
	}
}

func TestElGate_ErrorSePropaga(t *testing.T) {
	dbErr := stderrors.New("timeout")
	repo := &mocks.RepositoryMock{
		IsIntegrationActiveFn: func(ctx context.Context, businessID, typeID uint) (bool, error) {
			return false, dbErr
		},
	}
	uc := newPublicsiteUseCase(repo, nil)
	ctx := context.Background()

	_, err := uc.GetBusinessPage(ctx, "demo")
	assert.ErrorIs(t, err, dbErr)

	_, _, err = uc.ListCatalog(ctx, "demo", dtos.CatalogFilters{})
	assert.ErrorIs(t, err, dbErr)

	_, err = uc.GetProduct(ctx, "demo", "P-1")
	assert.ErrorIs(t, err, dbErr)

	_, err = uc.GetFeaturedProducts(ctx, "demo", 8)
	assert.ErrorIs(t, err, dbErr)

	_, err = uc.GetCategories(ctx, "demo")
	assert.ErrorIs(t, err, dbErr)

	assert.ErrorIs(t, uc.SubmitContact(ctx, "demo", &dtos.ContactFormDTO{Name: "a", Message: "m"}), dbErr)

	_, err = uc.CreateCheckoutSession(ctx, "demo", carritoValido(), 0)
	assert.ErrorIs(t, err, dbErr)
}

func TestErrorAlBuscarElNegocio_SePropaga(t *testing.T) {
	dbErr := stderrors.New("timeout")
	repo := &mocks.RepositoryMock{
		GetBusinessBySlugFn: func(ctx context.Context, slug string) (*entities.BusinessPage, error) {
			return nil, dbErr
		},
	}
	uc := newPublicsiteUseCase(repo, nil)
	ctx := context.Background()

	_, err := uc.GetBusinessPage(ctx, "demo")
	assert.ErrorIs(t, err, dbErr)

	_, _, err = uc.ListCatalog(ctx, "demo", dtos.CatalogFilters{})
	assert.ErrorIs(t, err, dbErr)

	_, err = uc.GetProduct(ctx, "demo", "P-1")
	assert.ErrorIs(t, err, dbErr)

	_, err = uc.GetFeaturedProducts(ctx, "demo", 8)
	assert.ErrorIs(t, err, dbErr)

	_, err = uc.GetCategories(ctx, "demo")
	assert.ErrorIs(t, err, dbErr)

	assert.ErrorIs(t, uc.SubmitContact(ctx, "demo", &dtos.ContactFormDTO{Name: "a", Message: "m"}), dbErr)

	_, err = uc.CreateCheckoutSession(ctx, "demo", carritoValido(), 0)
	assert.ErrorIs(t, err, dbErr)

	_, err = uc.GetSession(ctx, "demo", 42)
	assert.ErrorIs(t, err, dbErr)

	_, err = uc.GetMyPrices(ctx, "demo", 42)
	assert.ErrorIs(t, err, dbErr)

	_, _, err = uc.GetMyOrders(ctx, "demo", 42, 1, 10)
	assert.ErrorIs(t, err, dbErr)

	assert.ErrorIs(t, uc.UpdateProfile(ctx, "demo", 42, "n", "p", "d"), dbErr)
}

func TestCatalogFilters_NormalizeYOffset(t *testing.T) {
	casos := []struct {
		page, pageSize, wantPage, wantSize int
	}{
		{0, 12, 1, 12},
		{-5, 12, 1, 12},
		{1, 0, 1, 12},
		{1, 101, 1, 12},
		{1, 100, 1, 100},
		{3, 24, 3, 24},
	}

	for _, tc := range casos {
		f := dtos.CatalogFilters{Page: tc.page, PageSize: tc.pageSize}
		f.Normalize()
		assert.Equal(t, tc.wantPage, f.Page)
		assert.Equal(t, tc.wantSize, f.PageSize)
	}

	assert.Equal(t, 0, dtos.CatalogFilters{Page: 1, PageSize: 12}.Offset())
	assert.Equal(t, 12, dtos.CatalogFilters{Page: 2, PageSize: 12}.Offset())
	assert.Equal(t, 0, dtos.CatalogFilters{Page: -3, PageSize: 12}.Offset(),
		"una pagina invalida no debe producir offset negativo")
}

func TestListCatalog_NormalizaYFiltraPorNegocio(t *testing.T) {
	var vistoNegocio uint
	var visto dtos.CatalogFilters
	repo := &mocks.RepositoryMock{
		GetBusinessBySlugFn: func(ctx context.Context, slug string) (*entities.BusinessPage, error) {
			return &entities.BusinessPage{ID: 99, Code: slug}, nil
		},
		ListActiveProductsFn: func(ctx context.Context, businessID uint, f dtos.CatalogFilters) ([]entities.PublicProduct, int64, error) {
			vistoNegocio, visto = businessID, f
			return []entities.PublicProduct{{ID: "P-1"}}, 1, nil
		},
	}

	got, total, err := newPublicsiteUseCase(repo, nil).ListCatalog(context.Background(), "demo",
		dtos.CatalogFilters{Page: 0, PageSize: 500, Search: "camisa"})

	require.NoError(t, err)
	assert.Equal(t, uint(99), vistoNegocio,
		"el negocio sale del slug, no de un parametro del cliente")
	assert.Equal(t, 1, visto.Page)
	assert.Equal(t, 12, visto.PageSize)
	assert.Equal(t, "camisa", visto.Search)
	assert.Len(t, got, 1)
	assert.Equal(t, int64(1), total)
}

func TestGetFeaturedProducts_NormalizaElLimite(t *testing.T) {
	casos := []struct {
		limite, want int
	}{
		{0, 8}, {-5, 8}, {21, 8}, {1, 1}, {8, 8}, {20, 20},
	}

	for _, tc := range casos {
		var visto int
		repo := &mocks.RepositoryMock{
			GetFeaturedProductsFn: func(ctx context.Context, businessID uint, limit int) ([]entities.PublicProduct, error) {
				visto = limit
				return nil, nil
			},
		}

		_, err := newPublicsiteUseCase(repo, nil).GetFeaturedProducts(context.Background(), "demo", tc.limite)

		require.NoError(t, err)
		assert.Equal(t, tc.want, visto, "limite=%d", tc.limite)
	}
}

func TestGetProduct_FiltraPorNegocioYPropagaError(t *testing.T) {
	var vistoNegocio uint
	repo := &mocks.RepositoryMock{
		GetProductByIDFn: func(ctx context.Context, businessID uint, productID string) (*entities.PublicProduct, error) {
			vistoNegocio = businessID
			return &entities.PublicProduct{ID: productID}, nil
		},
	}

	got, err := newPublicsiteUseCase(repo, nil).GetProduct(context.Background(), "demo", "P-1")

	require.NoError(t, err)
	assert.Equal(t, uint(26), vistoNegocio)
	assert.Equal(t, "P-1", got.ID)

	repoErr := &mocks.RepositoryMock{
		GetProductByIDFn: func(ctx context.Context, businessID uint, productID string) (*entities.PublicProduct, error) {
			return nil, domainerrors.ErrProductNotFound
		},
	}
	_, err = newPublicsiteUseCase(repoErr, nil).GetProduct(context.Background(), "demo", "P-404")
	assert.ErrorIs(t, err, domainerrors.ErrProductNotFound)
}

func TestGetCategories_PasaLasOcultasAlFiltro(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetHiddenCategoriesFn: func(ctx context.Context, businessID uint) ([]string, error) {
			return []string{"interna", "descontinuada"}, nil
		},
		ListPublicCategoriesFn: func(ctx context.Context, businessID uint, hidden []string) ([]string, error) {
			return []string{"ropa", "calzado"}, nil
		},
	}

	got, err := newPublicsiteUseCase(repo, nil).GetCategories(context.Background(), "demo")

	require.NoError(t, err)
	require.Len(t, repo.HiddenPassed, 1)
	assert.Equal(t, []string{"interna", "descontinuada"}, repo.HiddenPassed[0],
		"las categorias ocultas deben llegar al filtro o se mostrarian al publico")
	assert.Equal(t, []string{"ropa", "calzado"}, got)
}

func TestGetCategories_FalloAlLeerOcultas_NoAborta(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetHiddenCategoriesFn: func(ctx context.Context, businessID uint) ([]string, error) {
			return nil, stderrors.New("timeout")
		},
	}

	_, err := newPublicsiteUseCase(repo, nil).GetCategories(context.Background(), "demo")

	require.NoError(t, err)
	require.Len(t, repo.HiddenPassed, 1)
	assert.Nil(t, repo.HiddenPassed[0])
}

func TestGetCategories_ErrorAlListar_SePropaga(t *testing.T) {
	dbErr := stderrors.New("timeout")
	repo := &mocks.RepositoryMock{
		ListPublicCategoriesFn: func(ctx context.Context, businessID uint, hidden []string) ([]string, error) {
			return nil, dbErr
		},
	}

	_, err := newPublicsiteUseCase(repo, nil).GetCategories(context.Background(), "demo")

	assert.ErrorIs(t, err, dbErr)
}

func TestSubmitContact_NombreYMensajeObligatorios(t *testing.T) {
	casos := []struct {
		nombre string
		dto    *dtos.ContactFormDTO
	}{
		{"sin nombre", &dtos.ContactFormDTO{Message: "hola"}},
		{"sin mensaje", &dtos.ContactFormDTO{Name: "Ana"}},
		{"sin ninguno", &dtos.ContactFormDTO{}},
	}

	for _, tc := range casos {
		t.Run(tc.nombre, func(t *testing.T) {
			consultado := false
			repo := &mocks.RepositoryMock{
				GetBusinessBySlugFn: func(ctx context.Context, slug string) (*entities.BusinessPage, error) {
					consultado = true
					return &entities.BusinessPage{ID: 26}, nil
				},
			}

			err := newPublicsiteUseCase(repo, nil).SubmitContact(context.Background(), "demo", tc.dto)

			assert.ErrorIs(t, err, domainerrors.ErrInvalidContact)
			assert.False(t, consultado, "se valida antes de tocar la base")
		})
	}
}

func TestSubmitContact_Valido_GuardaEnElNegocioDelSlug(t *testing.T) {
	var vistoNegocio uint
	var recibido *dtos.ContactFormDTO
	repo := &mocks.RepositoryMock{
		GetBusinessBySlugFn: func(ctx context.Context, slug string) (*entities.BusinessPage, error) {
			return &entities.BusinessPage{ID: 99}, nil
		},
		SaveContactSubmissionFn: func(ctx context.Context, businessID uint, dto *dtos.ContactFormDTO) error {
			vistoNegocio, recibido = businessID, dto
			return nil
		},
	}
	dto := &dtos.ContactFormDTO{Name: "Ana", Message: "quiero cotizar"}

	err := newPublicsiteUseCase(repo, nil).SubmitContact(context.Background(), "demo", dto)

	require.NoError(t, err)
	assert.Equal(t, uint(99), vistoNegocio)
	assert.Same(t, dto, recibido)
}

func TestSubmitContact_ErrorAlGuardar_SePropaga(t *testing.T) {
	dbErr := stderrors.New("db caida")
	repo := &mocks.RepositoryMock{
		SaveContactSubmissionFn: func(ctx context.Context, businessID uint, dto *dtos.ContactFormDTO) error {
			return dbErr
		},
	}

	err := newPublicsiteUseCase(repo, nil).SubmitContact(context.Background(), "demo",
		&dtos.ContactFormDTO{Name: "Ana", Message: "hola"})

	assert.ErrorIs(t, err, dbErr)
}

func TestCreateCheckout_CarritoVacio_Rechaza(t *testing.T) {
	repo := &mocks.RepositoryMock{}
	bold := &mocks.BoldGatewayMock{}

	_, err := newPublicsiteUseCase(repo, bold).CreateCheckoutSession(context.Background(), "demo",
		&dtos.CreateCheckoutDTO{Items: nil}, 0)

	assert.ErrorIs(t, err, domainerrors.ErrEmptyCart)
	assert.Nil(t, repo.CreatedCheckout)
	assert.Empty(t, bold.PublishedRefs)
}

func TestCreateCheckout_CantidadInvalida_Rechaza(t *testing.T) {
	for _, cantidad := range []int{0, -1, -50} {
		repo := &mocks.RepositoryMock{}

		_, err := newPublicsiteUseCase(repo, nil).CreateCheckoutSession(context.Background(), "demo",
			&dtos.CreateCheckoutDTO{Items: []dtos.CheckoutItemInput{{ProductID: "P-1", Quantity: cantidad}}}, 0)

		assert.ErrorIs(t, err, domainerrors.ErrInvalidCartItem, "cantidad=%d", cantidad)
		assert.Nil(t, repo.CreatedCheckout)
		assert.Empty(t, repo.GateChecks, "se valida antes de consultar nada")
	}
}

func TestCreateCheckout_PagoEnLinea_TodaviaNoDisponible(t *testing.T) {
	repo := &mocks.RepositoryMock{}
	bold := &mocks.BoldGatewayMock{}
	dto := carritoValido()
	dto.PaymentMethod = "online"

	_, err := newPublicsiteUseCase(repo, bold).
		CreateCheckoutSession(context.Background(), "demo", dto, 0)

	assert.ErrorIs(t, err, domainerrors.ErrOnlinePayNotReady)
	assert.Nil(t, repo.CreatedCheckout)
	assert.Empty(t, bold.PublishedRefs)
}

func TestCreateCheckout_ElPrecioSaleDelCatalogo(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetProductByIDFn: func(ctx context.Context, businessID uint, productID string) (*entities.PublicProduct, error) {
			return &entities.PublicProduct{ID: productID, SKU: "SKU-1", Name: "Camisa", Price: 50000}, nil
		},
	}

	got, err := newPublicsiteUseCase(repo, nil).CreateCheckoutSession(context.Background(), "demo",
		&dtos.CreateCheckoutDTO{Items: []dtos.CheckoutItemInput{{ProductID: "P-1", Quantity: 3}}}, 0)

	require.NoError(t, err)
	assert.InDelta(t, 150000.0, got.Amount, 0.001,
		"el total lo calcula el backend, el cliente no lo envia")
	require.NotNil(t, repo.CreatedCheckout)
	require.Len(t, repo.CreatedCheckout.Items, 1)
	assert.InDelta(t, 50000.0, repo.CreatedCheckout.Items[0].UnitPrice, 0.001)
	assert.Equal(t, "SKU-1", repo.CreatedCheckout.Items[0].SKU)
}

func TestCreateCheckout_SinUsuario_NoAplicaPreciosPersonalizados(t *testing.T) {
	consultado := false
	repo := &mocks.RepositoryMock{
		GetCustomerPricesFn: func(ctx context.Context, businessID, clientID uint) (map[string]float64, error) {
			consultado = true
			return map[string]float64{"P-1": 1}, nil
		},
		GetProductByIDFn: func(ctx context.Context, businessID uint, productID string) (*entities.PublicProduct, error) {
			return &entities.PublicProduct{ID: productID, Price: 50000}, nil
		},
	}

	got, err := newPublicsiteUseCase(repo, nil).
		CreateCheckoutSession(context.Background(), "demo", carritoValido(), 0)

	require.NoError(t, err)
	assert.False(t, consultado, "un visitante anonimo no tiene precios de cliente")
	assert.InDelta(t, 100000.0, got.Amount, 0.001)
}

func TestCreateCheckout_ConUsuarioYPrecioPersonalizado_LoAplica(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetClientSessionFn: func(ctx context.Context, businessID, userID uint) (*entities.TiendaSession, error) {
			return &entities.TiendaSession{CustomerID: 7}, nil
		},
		GetCustomerPricesFn: func(ctx context.Context, businessID, clientID uint) (map[string]float64, error) {
			return map[string]float64{"P-1": 30000}, nil
		},
		GetProductByIDFn: func(ctx context.Context, businessID uint, productID string) (*entities.PublicProduct, error) {
			return &entities.PublicProduct{ID: productID, Price: 50000}, nil
		},
	}

	got, err := newPublicsiteUseCase(repo, nil).CreateCheckoutSession(context.Background(), "demo",
		&dtos.CreateCheckoutDTO{Items: []dtos.CheckoutItemInput{{ProductID: "P-1", Quantity: 2}}}, 42)

	require.NoError(t, err)
	assert.InDelta(t, 60000.0, got.Amount, 0.001,
		"el precio acordado con el cliente gana sobre el de lista")
	assert.InDelta(t, 30000.0, repo.CreatedCheckout.Items[0].UnitPrice, 0.001)
}

func TestCreateCheckout_ProductoSinPrecioPersonalizado_UsaElDeLista(t *testing.T) {
	precios := map[string]float64{"P-1": 50000, "P-2": 20000}
	repo := &mocks.RepositoryMock{
		GetClientSessionFn: func(ctx context.Context, businessID, userID uint) (*entities.TiendaSession, error) {
			return &entities.TiendaSession{CustomerID: 7}, nil
		},
		GetCustomerPricesFn: func(ctx context.Context, businessID, clientID uint) (map[string]float64, error) {
			return map[string]float64{"P-1": 30000}, nil
		},
		GetProductByIDFn: func(ctx context.Context, businessID uint, productID string) (*entities.PublicProduct, error) {
			return &entities.PublicProduct{ID: productID, Price: precios[productID]}, nil
		},
	}

	got, err := newPublicsiteUseCase(repo, nil).CreateCheckoutSession(context.Background(), "demo",
		&dtos.CreateCheckoutDTO{Items: []dtos.CheckoutItemInput{
			{ProductID: "P-1", Quantity: 1},
			{ProductID: "P-2", Quantity: 1},
		}}, 42)

	require.NoError(t, err)
	assert.InDelta(t, 50000.0, got.Amount, 0.001, "30000 personalizado + 20000 de lista")
}

func TestCreateCheckout_UsuarioSinClienteAsociado_UsaPreciosDeLista(t *testing.T) {
	consultado := false
	repo := &mocks.RepositoryMock{
		GetClientSessionFn: func(ctx context.Context, businessID, userID uint) (*entities.TiendaSession, error) {
			return &entities.TiendaSession{CustomerID: 0}, nil
		},
		GetCustomerPricesFn: func(ctx context.Context, businessID, clientID uint) (map[string]float64, error) {
			consultado = true
			return nil, nil
		},
	}

	_, err := newPublicsiteUseCase(repo, nil).
		CreateCheckoutSession(context.Background(), "demo", carritoValido(), 42)

	require.NoError(t, err)
	assert.False(t, consultado)
}

func TestCreateCheckout_FalloAlLeerPreciosPersonalizados_NoAborta(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetClientSessionFn: func(ctx context.Context, businessID, userID uint) (*entities.TiendaSession, error) {
			return &entities.TiendaSession{CustomerID: 7}, nil
		},
		GetCustomerPricesFn: func(ctx context.Context, businessID, clientID uint) (map[string]float64, error) {
			return nil, stderrors.New("timeout")
		},
		GetProductByIDFn: func(ctx context.Context, businessID uint, productID string) (*entities.PublicProduct, error) {
			return &entities.PublicProduct{ID: productID, Price: 50000}, nil
		},
	}

	got, err := newPublicsiteUseCase(repo, nil).
		CreateCheckoutSession(context.Background(), "demo", carritoValido(), 42)

	require.NoError(t, err, "si no se pueden leer los precios acordados se cobra el de lista")
	assert.InDelta(t, 100000.0, got.Amount, 0.001)
}

func TestCreateCheckout_ProductoInexistente_NoCrea(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetProductByIDFn: func(ctx context.Context, businessID uint, productID string) (*entities.PublicProduct, error) {
			return nil, domainerrors.ErrProductNotFound
		},
	}
	bold := &mocks.BoldGatewayMock{}

	_, err := newPublicsiteUseCase(repo, bold).
		CreateCheckoutSession(context.Background(), "demo", carritoValido(), 0)

	assert.ErrorIs(t, err, domainerrors.ErrProductNotFound)
	assert.Nil(t, repo.CreatedCheckout)
	assert.Empty(t, bold.PublishedRefs)
}

func TestCreateCheckout_NaceComoAcordadoEnPesos(t *testing.T) {
	repo := &mocks.RepositoryMock{}

	got, err := newPublicsiteUseCase(repo, nil).
		CreateCheckoutSession(context.Background(), "demo", carritoValido(), 0)

	require.NoError(t, err)
	assert.Equal(t, "agree", got.PaymentMethod)
	assert.Equal(t, "COP", got.Currency)
	require.NotNil(t, repo.CreatedCheckout)
	assert.Equal(t, entities.CheckoutStatusAgreed, repo.CreatedCheckout.Status)
	assert.Equal(t, "COP", repo.CreatedCheckout.Currency)
}

func TestCreateCheckout_LaReferenciaEsUnicaYConPrefijo(t *testing.T) {
	vistas := map[string]bool{}
	uc := newPublicsiteUseCase(nil, nil)

	for i := 0; i < 30; i++ {
		got, err := uc.CreateCheckoutSession(context.Background(), "demo", carritoValido(), 0)
		require.NoError(t, err)
		require.False(t, vistas[got.Reference], "referencia repetida: %q", got.Reference)
		vistas[got.Reference] = true
		assert.True(t, strings.HasPrefix(got.Reference, agreedReferencePrefix),
			"referencia sin prefijo: %q", got.Reference)
	}
	assert.Len(t, vistas, 30)
}

func TestCreateCheckout_GuardaLosDatosDelCliente(t *testing.T) {
	repo := &mocks.RepositoryMock{}
	dto := carritoValido()
	dto.CustomerPhone = "3001234567"
	dto.CustomerDni = "1020304050"
	dto.Address = &dtos.CheckoutAddressInput{
		Street: "Cra 1 #2-3", City: "Bogota", State: "Cundinamarca",
		Country: "CO", PostalCode: "110111", Instructions: "timbre 2",
	}

	_, err := newPublicsiteUseCase(repo, nil).
		CreateCheckoutSession(context.Background(), "demo", dto, 0)

	require.NoError(t, err)
	c := repo.CreatedCheckout
	assert.Equal(t, "Ana", c.CustomerName)
	assert.Equal(t, "ana@x.com", c.CustomerEmail)
	assert.Equal(t, "3001234567", c.CustomerPhone)
	assert.Equal(t, "1020304050", c.CustomerDni)
	require.NotNil(t, c.ShippingAddress)
	assert.Equal(t, "Cra 1 #2-3", c.ShippingAddress.Street)
	assert.Equal(t, "timbre 2", c.ShippingAddress.Instructions)
}

func TestCreateCheckout_SinDireccion_QuedaNil(t *testing.T) {
	repo := &mocks.RepositoryMock{}

	_, err := newPublicsiteUseCase(repo, nil).
		CreateCheckoutSession(context.Background(), "demo", carritoValido(), 0)

	require.NoError(t, err)
	assert.Nil(t, repo.CreatedCheckout.ShippingAddress)
}

func TestCreateCheckout_PublicaLaOrdenConLaMismaReferencia(t *testing.T) {
	repo := &mocks.RepositoryMock{}
	bold := &mocks.BoldGatewayMock{}

	got, err := newPublicsiteUseCase(repo, bold).
		CreateCheckoutSession(context.Background(), "demo", carritoValido(), 0)

	require.NoError(t, err)
	require.Len(t, bold.PublishedRefs, 1)
	assert.Equal(t, got.Reference, bold.PublishedRefs[0])
	assert.Equal(t, got.Reference, repo.CreatedCheckout.Reference)
	assert.Equal(t, got.Reference, repo.CreatedCheckout.BoldOrderID)
}

func TestCreateCheckout_FallaAlGuardar_NoPublica(t *testing.T) {
	dbErr := stderrors.New("db caida")
	bold := &mocks.BoldGatewayMock{}
	repo := &mocks.RepositoryMock{
		CreateCheckoutFn: func(ctx context.Context, c *entities.PublicCheckout) error { return dbErr },
	}

	_, err := newPublicsiteUseCase(repo, bold).
		CreateCheckoutSession(context.Background(), "demo", carritoValido(), 0)

	assert.ErrorIs(t, err, dbErr)
	assert.Empty(t, bold.PublishedRefs,
		"no se publica una orden cuyo checkout no quedo guardado")
}

func TestCreateCheckout_FallaAlPublicar_SePropaga(t *testing.T) {
	pubErr := stderrors.New("rabbit caido")
	bold := &mocks.BoldGatewayMock{
		PublishAgreedStorefrontOrderFn: func(ctx context.Context, reference string) error { return pubErr },
	}

	_, err := newPublicsiteUseCase(nil, bold).
		CreateCheckoutSession(context.Background(), "demo", carritoValido(), 0)

	assert.ErrorIs(t, err, pubErr,
		"si la orden no llega a la cola el cliente debe enterarse")
}

func TestCreateCheckout_SinIntegracionDeTienda_SePropaga(t *testing.T) {
	dbErr := stderrors.New("no existe")
	repo := &mocks.RepositoryMock{
		GetTiendaIntegrationIDFn: func(ctx context.Context, businessID uint) (uint, error) {
			return 0, dbErr
		},
	}

	_, err := newPublicsiteUseCase(repo, nil).
		CreateCheckoutSession(context.Background(), "demo", carritoValido(), 0)

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, repo.CreatedCheckout)
}

func TestGetCheckoutStatus_DevuelveElEstado(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetCheckoutByReferenceFn: func(ctx context.Context, reference string) (*entities.PublicCheckout, error) {
			return &entities.PublicCheckout{Reference: reference, Status: entities.CheckoutStatusPaid}, nil
		},
	}

	got, err := newPublicsiteUseCase(repo, nil).GetCheckoutStatus(context.Background(), "SFA123")

	require.NoError(t, err)
	assert.Equal(t, entities.CheckoutStatusPaid, got)
}

func TestGetCheckoutStatus_Inexistente_Rechaza(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetCheckoutByReferenceFn: func(ctx context.Context, reference string) (*entities.PublicCheckout, error) {
			return nil, nil
		},
	}

	got, err := newPublicsiteUseCase(repo, nil).GetCheckoutStatus(context.Background(), "SFA404")

	assert.ErrorIs(t, err, domainerrors.ErrCheckoutNotFound)
	assert.Empty(t, got)
}

func TestGetCheckoutStatus_ErrorDelRepo_SePropaga(t *testing.T) {
	dbErr := stderrors.New("timeout")
	repo := &mocks.RepositoryMock{
		GetCheckoutByReferenceFn: func(ctx context.Context, reference string) (*entities.PublicCheckout, error) {
			return nil, dbErr
		},
	}

	_, err := newPublicsiteUseCase(repo, nil).GetCheckoutStatus(context.Background(), "SFA1")

	assert.ErrorIs(t, err, dbErr)
}

func TestGetSession_UsuarioDeOtroNegocio_Rechaza(t *testing.T) {
	repo := &mocks.RepositoryMock{
		UserBelongsToBusinessFn: func(ctx context.Context, businessID, userID uint) (bool, error) {
			return false, nil
		},
	}

	_, err := newPublicsiteUseCase(repo, nil).GetSession(context.Background(), "demo", 42)

	assert.ErrorIs(t, err, domainerrors.ErrBusinessNotFound,
		"un usuario de otro negocio no debe obtener sesion en esta tienda")
}

func TestGetSession_ConClienteAsociado_DevuelveEsaSesionConAvatar(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetUserAvatarFn: func(ctx context.Context, userID uint) (string, error) {
			return "s3://avatar.png", nil
		},
		GetClientSessionFn: func(ctx context.Context, businessID, userID uint) (*entities.TiendaSession, error) {
			return &entities.TiendaSession{Name: "Ana Cliente", CustomerID: 7}, nil
		},
	}

	got, err := newPublicsiteUseCase(repo, nil).GetSession(context.Background(), "demo", 42)

	require.NoError(t, err)
	assert.Equal(t, "Ana Cliente", got.Name)
	assert.Equal(t, uint(7), got.CustomerID)
	assert.Equal(t, "s3://avatar.png", got.AvatarURL)
}

func TestGetSession_SinClienteAsociado_CaeALosDatosBasicos(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetClientSessionFn: func(ctx context.Context, businessID, userID uint) (*entities.TiendaSession, error) {
			return nil, nil
		},
		GetUserBasicFn: func(ctx context.Context, userID uint) (*entities.TiendaSession, error) {
			return &entities.TiendaSession{Name: "Usuario Basico"}, nil
		},
		GetUserAvatarFn: func(ctx context.Context, userID uint) (string, error) {
			return "s3://a.png", nil
		},
	}

	got, err := newPublicsiteUseCase(repo, nil).GetSession(context.Background(), "demo", 42)

	require.NoError(t, err)
	assert.Equal(t, "Usuario Basico", got.Name)
	assert.Zero(t, got.CustomerID)
	assert.Equal(t, "s3://a.png", got.AvatarURL)
}

func TestGetSession_SinClienteNiUsuario_Rechaza(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetClientSessionFn: func(ctx context.Context, businessID, userID uint) (*entities.TiendaSession, error) {
			return nil, nil
		},
		GetUserBasicFn: func(ctx context.Context, userID uint) (*entities.TiendaSession, error) {
			return nil, nil
		},
	}

	_, err := newPublicsiteUseCase(repo, nil).GetSession(context.Background(), "demo", 42)

	assert.ErrorIs(t, err, domainerrors.ErrBusinessNotFound)
}

func TestGetSession_FalloAlLeerAvatar_NoAborta(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetUserAvatarFn: func(ctx context.Context, userID uint) (string, error) {
			return "", stderrors.New("s3 caido")
		},
	}

	got, err := newPublicsiteUseCase(repo, nil).GetSession(context.Background(), "demo", 42)

	require.NoError(t, err, "el avatar es cosmetico, no bloquea la sesion")
	assert.Empty(t, got.AvatarURL)
}

func TestGetSession_ErroresDelRepo_SePropagan(t *testing.T) {
	dbErr := stderrors.New("db caida")

	t.Run("al verificar pertenencia", func(t *testing.T) {
		repo := &mocks.RepositoryMock{
			UserBelongsToBusinessFn: func(ctx context.Context, businessID, userID uint) (bool, error) {
				return false, dbErr
			},
		}
		_, err := newPublicsiteUseCase(repo, nil).GetSession(context.Background(), "demo", 42)
		assert.ErrorIs(t, err, dbErr)
	})

	t.Run("al leer la sesion de cliente", func(t *testing.T) {
		repo := &mocks.RepositoryMock{
			GetClientSessionFn: func(ctx context.Context, businessID, userID uint) (*entities.TiendaSession, error) {
				return nil, dbErr
			},
		}
		_, err := newPublicsiteUseCase(repo, nil).GetSession(context.Background(), "demo", 42)
		assert.ErrorIs(t, err, dbErr)
	})

	t.Run("al leer los datos basicos", func(t *testing.T) {
		repo := &mocks.RepositoryMock{
			GetClientSessionFn: func(ctx context.Context, businessID, userID uint) (*entities.TiendaSession, error) {
				return nil, nil
			},
			GetUserBasicFn: func(ctx context.Context, userID uint) (*entities.TiendaSession, error) {
				return nil, dbErr
			},
		}
		_, err := newPublicsiteUseCase(repo, nil).GetSession(context.Background(), "demo", 42)
		assert.ErrorIs(t, err, dbErr)
	})
}

func TestGetMyPrices_SinClienteAsociado_DevuelveMapaVacio(t *testing.T) {
	casos := []struct {
		nombre  string
		session *entities.TiendaSession
	}{
		{"sin sesion", nil},
		{"sesion sin cliente", &entities.TiendaSession{CustomerID: 0}},
	}

	for _, tc := range casos {
		t.Run(tc.nombre, func(t *testing.T) {
			consultado := false
			repo := &mocks.RepositoryMock{
				GetClientSessionFn: func(ctx context.Context, businessID, userID uint) (*entities.TiendaSession, error) {
					return tc.session, nil
				},
				GetCustomerPricesFn: func(ctx context.Context, businessID, clientID uint) (map[string]float64, error) {
					consultado = true
					return nil, nil
				},
			}

			got, err := newPublicsiteUseCase(repo, nil).GetMyPrices(context.Background(), "demo", 42)

			require.NoError(t, err)
			require.NotNil(t, got)
			assert.Empty(t, got)
			assert.False(t, consultado)
		})
	}
}

func TestGetMyPrices_ConCliente_DevuelveSusPrecios(t *testing.T) {
	var vistoCliente uint
	repo := &mocks.RepositoryMock{
		GetClientSessionFn: func(ctx context.Context, businessID, userID uint) (*entities.TiendaSession, error) {
			return &entities.TiendaSession{CustomerID: 7}, nil
		},
		GetCustomerPricesFn: func(ctx context.Context, businessID, clientID uint) (map[string]float64, error) {
			vistoCliente = clientID
			return map[string]float64{"P-1": 30000}, nil
		},
	}

	got, err := newPublicsiteUseCase(repo, nil).GetMyPrices(context.Background(), "demo", 42)

	require.NoError(t, err)
	assert.Equal(t, uint(7), vistoCliente)
	assert.InDelta(t, 30000.0, got["P-1"], 0.001)
}

func TestGetMyPrices_ErroresDelRepo_SePropagan(t *testing.T) {
	dbErr := stderrors.New("db caida")

	repo := &mocks.RepositoryMock{
		GetClientSessionFn: func(ctx context.Context, businessID, userID uint) (*entities.TiendaSession, error) {
			return nil, dbErr
		},
	}
	_, err := newPublicsiteUseCase(repo, nil).GetMyPrices(context.Background(), "demo", 42)
	assert.ErrorIs(t, err, dbErr)

	repo2 := &mocks.RepositoryMock{
		GetCustomerPricesFn: func(ctx context.Context, businessID, clientID uint) (map[string]float64, error) {
			return nil, dbErr
		},
	}
	_, err = newPublicsiteUseCase(repo2, nil).GetMyPrices(context.Background(), "demo", 42)
	assert.ErrorIs(t, err, dbErr)
}

func TestGetMyOrders_NormalizaPaginacion(t *testing.T) {
	casos := []struct {
		page, pageSize, wantPage, wantSize int
	}{
		{0, 10, 1, 10},
		{-5, 10, 1, 10},
		{1, 0, 1, 10},
		{1, 51, 1, 10},
		{1, 50, 1, 50},
		{2, 25, 2, 25},
	}

	for _, tc := range casos {
		var vistoPage, vistoSize int
		repo := &mocks.RepositoryMock{
			ListCustomerOrdersFn: func(ctx context.Context, businessID, userID, customerID uint, page, pageSize int) ([]entities.TiendaOrder, int64, error) {
				vistoPage, vistoSize = page, pageSize
				return nil, 0, nil
			},
		}

		_, _, err := newPublicsiteUseCase(repo, nil).
			GetMyOrders(context.Background(), "demo", 42, tc.page, tc.pageSize)

		require.NoError(t, err)
		assert.Equal(t, tc.wantPage, vistoPage)
		assert.Equal(t, tc.wantSize, vistoSize)
	}
}

func TestGetMyOrders_FiltraPorNegocioUsuarioYCliente(t *testing.T) {
	var vistoNegocio, vistoUsuario, vistoCliente uint
	repo := &mocks.RepositoryMock{
		GetBusinessBySlugFn: func(ctx context.Context, slug string) (*entities.BusinessPage, error) {
			return &entities.BusinessPage{ID: 99}, nil
		},
		GetClientSessionFn: func(ctx context.Context, businessID, userID uint) (*entities.TiendaSession, error) {
			return &entities.TiendaSession{CustomerID: 7}, nil
		},
		ListCustomerOrdersFn: func(ctx context.Context, businessID, userID, customerID uint, page, pageSize int) ([]entities.TiendaOrder, int64, error) {
			vistoNegocio, vistoUsuario, vistoCliente = businessID, userID, customerID
			return []entities.TiendaOrder{{ID: "O-1"}}, 1, nil
		},
	}

	got, total, err := newPublicsiteUseCase(repo, nil).
		GetMyOrders(context.Background(), "demo", 42, 1, 10)

	require.NoError(t, err)
	assert.Equal(t, uint(99), vistoNegocio)
	assert.Equal(t, uint(42), vistoUsuario)
	assert.Equal(t, uint(7), vistoCliente)
	assert.Len(t, got, 1)
	assert.Equal(t, int64(1), total)
}

func TestGetMyOrders_FalloAlLeerLaSesion_ListaSinCliente(t *testing.T) {
	var vistoCliente uint
	repo := &mocks.RepositoryMock{
		GetClientSessionFn: func(ctx context.Context, businessID, userID uint) (*entities.TiendaSession, error) {
			return nil, stderrors.New("timeout")
		},
		ListCustomerOrdersFn: func(ctx context.Context, businessID, userID, customerID uint, page, pageSize int) ([]entities.TiendaOrder, int64, error) {
			vistoCliente = customerID
			return nil, 0, nil
		},
	}

	_, _, err := newPublicsiteUseCase(repo, nil).GetMyOrders(context.Background(), "demo", 42, 1, 10)

	require.NoError(t, err)
	assert.Zero(t, vistoCliente, "sin cliente resuelto se listan solo las del usuario")
}

func TestResolveClient_SinClienteAsociado_BloqueaElPerfil(t *testing.T) {
	casos := []struct {
		nombre  string
		session *entities.TiendaSession
	}{
		{"sin sesion", nil},
		{"sesion sin cliente", &entities.TiendaSession{CustomerID: 0}},
	}

	for _, tc := range casos {
		t.Run(tc.nombre, func(t *testing.T) {
			repo := &mocks.RepositoryMock{
				GetClientSessionFn: func(ctx context.Context, businessID, userID uint) (*entities.TiendaSession, error) {
					return tc.session, nil
				},
			}
			uc := newPublicsiteUseCase(repo, nil)
			ctx := context.Background()

			assert.ErrorIs(t, uc.UpdateProfile(ctx, "demo", 42, "n", "p", "d"),
				domainerrors.ErrBusinessNotFound)
			assert.ErrorIs(t, uc.SaveAvatar(ctx, "demo", 42, "url"),
				domainerrors.ErrBusinessNotFound)

			_, err := uc.ListAddresses(ctx, "demo", 42)
			assert.ErrorIs(t, err, domainerrors.ErrBusinessNotFound)

			assert.ErrorIs(t, uc.AddAddress(ctx, "demo", 42, &entities.CustomerAddress{}),
				domainerrors.ErrBusinessNotFound)
			assert.ErrorIs(t, uc.SetPrimaryAddress(ctx, "demo", 42, 1),
				domainerrors.ErrBusinessNotFound)
			assert.ErrorIs(t, uc.DeleteAddress(ctx, "demo", 42, 1),
				domainerrors.ErrBusinessNotFound)
		})
	}
}

func TestUpdateProfile_PasaNegocioYUsuario(t *testing.T) {
	var vistoNegocio, vistoUsuario uint
	var vistoNombre, vistoTel, vistoDni string
	repo := &mocks.RepositoryMock{
		GetBusinessBySlugFn: func(ctx context.Context, slug string) (*entities.BusinessPage, error) {
			return &entities.BusinessPage{ID: 99}, nil
		},
		UpdateClientProfileFn: func(ctx context.Context, businessID, userID uint, name, phone, dni string) error {
			vistoNegocio, vistoUsuario = businessID, userID
			vistoNombre, vistoTel, vistoDni = name, phone, dni
			return nil
		},
	}

	err := newPublicsiteUseCase(repo, nil).
		UpdateProfile(context.Background(), "demo", 42, "Ana", "3001234567", "1020304050")

	require.NoError(t, err)
	assert.Equal(t, uint(99), vistoNegocio)
	assert.Equal(t, uint(42), vistoUsuario)
	assert.Equal(t, "Ana", vistoNombre)
	assert.Equal(t, "3001234567", vistoTel)
	assert.Equal(t, "1020304050", vistoDni)
}

func TestSaveAvatar_GuardaContraElUsuario(t *testing.T) {
	var vistoUsuario uint
	var vistaURL string
	repo := &mocks.RepositoryMock{
		UpdateUserAvatarFn: func(ctx context.Context, userID uint, avatarURL string) error {
			vistoUsuario, vistaURL = userID, avatarURL
			return nil
		},
	}

	err := newPublicsiteUseCase(repo, nil).
		SaveAvatar(context.Background(), "demo", 42, "s3://nuevo.png")

	require.NoError(t, err)
	assert.Equal(t, uint(42), vistoUsuario)
	assert.Equal(t, "s3://nuevo.png", vistaURL)
}

func TestAddAddress_LaPrimeraQuedaComoPrincipal(t *testing.T) {
	repo := &mocks.RepositoryMock{
		ListCustomerAddressesFn: func(ctx context.Context, businessID, customerID uint) ([]entities.CustomerAddress, error) {
			return nil, nil
		},
	}
	addr := &entities.CustomerAddress{Street: "Cra 1", IsPrimary: false}

	err := newPublicsiteUseCase(repo, nil).AddAddress(context.Background(), "demo", 42, addr)

	require.NoError(t, err)
	assert.True(t, addr.IsPrimary,
		"si es la unica direccion debe quedar como principal aunque no lo pidan")
}

func TestAddAddress_ConDireccionesPrevias_RespetaElFlag(t *testing.T) {
	repo := &mocks.RepositoryMock{
		ListCustomerAddressesFn: func(ctx context.Context, businessID, customerID uint) ([]entities.CustomerAddress, error) {
			return []entities.CustomerAddress{{ID: 1, IsPrimary: true}}, nil
		},
	}
	addr := &entities.CustomerAddress{Street: "Cra 2", IsPrimary: false}

	err := newPublicsiteUseCase(repo, nil).AddAddress(context.Background(), "demo", 42, addr)

	require.NoError(t, err)
	assert.False(t, addr.IsPrimary,
		"agregar una segunda direccion no debe desplazar a la principal existente")
}

func TestAddAddress_ErrorAlListar_NoCrea(t *testing.T) {
	dbErr := stderrors.New("timeout")
	repo := &mocks.RepositoryMock{
		ListCustomerAddressesFn: func(ctx context.Context, businessID, customerID uint) ([]entities.CustomerAddress, error) {
			return nil, dbErr
		},
	}

	err := newPublicsiteUseCase(repo, nil).
		AddAddress(context.Background(), "demo", 42, &entities.CustomerAddress{})

	assert.ErrorIs(t, err, dbErr)
	assert.Empty(t, repo.CreatedAddresses)
}

func TestDireccionesYPerfil_PasanElClienteResuelto(t *testing.T) {
	var vistoNegocio, vistoCliente, vistaDireccion uint
	repo := &mocks.RepositoryMock{
		GetBusinessBySlugFn: func(ctx context.Context, slug string) (*entities.BusinessPage, error) {
			return &entities.BusinessPage{ID: 99}, nil
		},
		GetClientSessionFn: func(ctx context.Context, businessID, userID uint) (*entities.TiendaSession, error) {
			return &entities.TiendaSession{CustomerID: 7}, nil
		},
		ListCustomerAddressesFn: func(ctx context.Context, businessID, customerID uint) ([]entities.CustomerAddress, error) {
			vistoNegocio, vistoCliente = businessID, customerID
			return []entities.CustomerAddress{{ID: 1}}, nil
		},
		SetPrimaryAddressFn: func(ctx context.Context, businessID, customerID, addressID uint) error {
			vistaDireccion = addressID
			return nil
		},
	}
	uc := newPublicsiteUseCase(repo, nil)

	got, err := uc.ListAddresses(context.Background(), "demo", 42)
	require.NoError(t, err)
	assert.Equal(t, uint(99), vistoNegocio)
	assert.Equal(t, uint(7), vistoCliente,
		"las direcciones son del cliente resuelto, no de un ID que mande el front")
	assert.Len(t, got, 1)

	require.NoError(t, uc.SetPrimaryAddress(context.Background(), "demo", 42, 3))
	assert.Equal(t, uint(3), vistaDireccion)
}

func TestDeleteAddress_PasaElClienteResuelto(t *testing.T) {
	var vistoNegocio, vistoCliente, vistaDireccion uint
	repo := &mocks.RepositoryMock{
		GetClientSessionFn: func(ctx context.Context, businessID, userID uint) (*entities.TiendaSession, error) {
			return &entities.TiendaSession{CustomerID: 7}, nil
		},
		DeleteCustomerAddressFn: func(ctx context.Context, businessID, customerID, addressID uint) error {
			vistoNegocio, vistoCliente, vistaDireccion = businessID, customerID, addressID
			return nil
		},
	}

	err := newPublicsiteUseCase(repo, nil).DeleteAddress(context.Background(), "demo", 42, 3)

	require.NoError(t, err)
	assert.Equal(t, uint(26), vistoNegocio)
	assert.Equal(t, uint(7), vistoCliente,
		"no se puede borrar la direccion de otro cliente pasando su ID")
	assert.Equal(t, uint(3), vistaDireccion)
}

func TestPerfilYDirecciones_ErroresDelRepo_SePropagan(t *testing.T) {
	dbErr := stderrors.New("db caida")
	repo := &mocks.RepositoryMock{
		UpdateClientProfileFn: func(ctx context.Context, businessID, userID uint, n, p, d string) error {
			return dbErr
		},
		UpdateUserAvatarFn: func(ctx context.Context, userID uint, avatarURL string) error { return dbErr },
		CreateCustomerAddressFn: func(ctx context.Context, businessID, customerID uint, a *entities.CustomerAddress) error {
			return dbErr
		},
		SetPrimaryAddressFn: func(ctx context.Context, businessID, customerID, addressID uint) error {
			return dbErr
		},
		DeleteCustomerAddressFn: func(ctx context.Context, businessID, customerID, addressID uint) error {
			return dbErr
		},
	}
	uc := newPublicsiteUseCase(repo, nil)
	ctx := context.Background()

	assert.ErrorIs(t, uc.UpdateProfile(ctx, "demo", 42, "n", "p", "d"), dbErr)
	assert.ErrorIs(t, uc.SaveAvatar(ctx, "demo", 42, "url"), dbErr)
	assert.ErrorIs(t, uc.AddAddress(ctx, "demo", 42, &entities.CustomerAddress{}), dbErr)
	assert.ErrorIs(t, uc.SetPrimaryAddress(ctx, "demo", 42, 1), dbErr)
	assert.ErrorIs(t, uc.DeleteAddress(ctx, "demo", 42, 1), dbErr)

	repoList := &mocks.RepositoryMock{
		ListCustomerAddressesFn: func(ctx context.Context, businessID, customerID uint) ([]entities.CustomerAddress, error) {
			return nil, dbErr
		},
	}
	_, err := newPublicsiteUseCase(repoList, nil).ListAddresses(ctx, "demo", 42)
	assert.ErrorIs(t, err, dbErr)
}

func TestResolveClient_ErrorAlLeerLaSesion_SePropaga(t *testing.T) {
	dbErr := stderrors.New("timeout")
	repo := &mocks.RepositoryMock{
		GetClientSessionFn: func(ctx context.Context, businessID, userID uint) (*entities.TiendaSession, error) {
			return nil, dbErr
		},
	}

	assert.ErrorIs(t, newPublicsiteUseCase(repo, nil).
		UpdateProfile(context.Background(), "demo", 42, "n", "p", "d"), dbErr)
}

func TestGetBusinessPage_DevuelveLaPagina(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetBusinessBySlugFn: func(ctx context.Context, slug string) (*entities.BusinessPage, error) {
			return &entities.BusinessPage{ID: 26, Name: "Demo", Code: slug, PrimaryColor: "#111"}, nil
		},
	}

	got, err := newPublicsiteUseCase(repo, nil).GetBusinessPage(context.Background(), "demo")

	require.NoError(t, err)
	assert.Equal(t, "Demo", got.Name)
	assert.Equal(t, "demo", got.Code)
	assert.Equal(t, "#111", got.PrimaryColor)
}

func TestListCatalog_ErrorDelRepo_SePropaga(t *testing.T) {
	dbErr := stderrors.New("timeout")
	repo := &mocks.RepositoryMock{
		ListActiveProductsFn: func(ctx context.Context, businessID uint, f dtos.CatalogFilters) ([]entities.PublicProduct, int64, error) {
			return nil, 0, dbErr
		},
	}

	_, _, err := newPublicsiteUseCase(repo, nil).
		ListCatalog(context.Background(), "demo", dtos.CatalogFilters{})

	assert.ErrorIs(t, err, dbErr)
}

func TestGetFeaturedProducts_ErrorDelRepo_SePropaga(t *testing.T) {
	dbErr := stderrors.New("timeout")
	repo := &mocks.RepositoryMock{
		GetFeaturedProductsFn: func(ctx context.Context, businessID uint, limit int) ([]entities.PublicProduct, error) {
			return nil, dbErr
		},
	}

	_, err := newPublicsiteUseCase(repo, nil).GetFeaturedProducts(context.Background(), "demo", 8)

	assert.ErrorIs(t, err, dbErr)
}

func TestGetMyOrders_ErrorDelRepo_SePropaga(t *testing.T) {
	dbErr := stderrors.New("timeout")
	repo := &mocks.RepositoryMock{
		ListCustomerOrdersFn: func(ctx context.Context, businessID, userID, customerID uint, page, pageSize int) ([]entities.TiendaOrder, int64, error) {
			return nil, 0, dbErr
		},
	}

	_, _, err := newPublicsiteUseCase(repo, nil).GetMyOrders(context.Background(), "demo", 42, 1, 10)

	assert.ErrorIs(t, err, dbErr)
}
